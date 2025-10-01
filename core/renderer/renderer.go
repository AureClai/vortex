//go:build js && wasm

package renderer

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/AureClai/vortex/core/component"
)

type Renderer struct {
	container    js.Value         // Container element
	currentVNode *component.VNode // Current virtual node state for patching algorithms

	styleElement    js.Value        // Ref to the <style> balise
	injectedClasses map[string]bool // Keep track from the classes already injected
	animationFrame  js.Value

	// incremental rendering
	dirtySet     map[component.Component]struct{}
	rafScheduled bool
	// index of component -> vnode for current tree
	compIndex map[component.Component]*component.VNode
}

func NewRenderer(containerID string) *Renderer {
	document := js.Global().Get("document")
	container := document.Call("getElementById", containerID)

	// Styling
	head := document.Get("head")
	styleEl := document.Call("createElement", "style")
	styleEl.Set("id", "vortex-styles")
	head.Call("appendChild", styleEl)

	return &Renderer{
		container:    container,
		currentVNode: nil,

		styleElement:    styleEl,
		injectedClasses: make(map[string]bool),
		dirtySet:        make(map[component.Component]struct{}),
		compIndex:       make(map[component.Component]*component.VNode),
	}
}

func (r *Renderer) Render(newVNode *component.VNode) {
	expanded := r.resolveComponents(newVNode) // return pure elements
	r.buildComponentIndex(expanded)
	r.Patch(r.container, r.currentVNode, expanded)
	r.currentVNode = expanded
}

func (r *Renderer) Invalidate(comp component.Component) {
	r.dirtySet[comp] = struct{}{}
	if !r.rafScheduled {
		r.rafScheduled = true
		js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			r.rafScheduled = false
			r.flushDirty()
			return nil
		}))
	}
}

func (r *Renderer) flushDirty() {
	// 1) Collapse dirty set: if an ancestor comp is dirty, drop its descendants
	roots := r.filterDirtyRoots()

	// 2) For eash dirty root component : re-render its subtree and patch at boundary
	for comp := range roots {
		oldSub := r.compIndex[comp]
		if oldSub == nil || !oldSub.Element.Truthy() {
			continue
		}
		parent := oldSub.Element.Get("parentNode")
		newSub := r.resolveComponents(comp.Render())
		// Patch subtree
		r.Patch(parent, oldSub, newSub)
		// replace in current tree and reindex this branch
		r.replaceSubtree(r.currentVNode, oldSub, newSub)
		r.reindexComponentBranch(newSub)
	}

	// clear set
	for k := range r.dirtySet {
		delete(r.dirtySet, k)
	}
}

// resolveComponents expands VNodeComponent nodes into element/text trees,
// preserving EventHandlers/Attrs/etc. It also attaches Component reference to
// the expanded subtree root for indexing.
func (r *Renderer) resolveComponents(v *component.VNode) *component.VNode {
	if v == nil {
		return nil
	}
	if v.Type == component.VNodeComponent && v.Component != nil {
		expanded := v.Component.Render()
		res := r.resolveComponents(expanded)
		// Tag the root with the component (for indexing boundary)
		if res != nil {
			res.Component = v.Component
		}
		return res
	}
	// clone node shallowly, resolve children
	out := *v
	out.Children = make([]*component.VNode, len(v.Children))
	for _, ch := range v.Children {
		out.Children = append(out.Children, r.resolveComponents(ch))
	}
	return &out
}

// buildComponentIndex builds the index of components in the tree
func (r *Renderer) buildComponentIndex(v *component.VNode) {
	r.compIndex = make(map[component.Component]*component.VNode)
	var walk func(*component.VNode)
	walk = func(v *component.VNode) {
		if v == nil {
			return
		}
		if v.Component != nil {
			r.compIndex[v.Component] = v
		}
		for _, ch := range v.Children {
			walk(ch)
		}
	}
	walk(v)
}

// reindexComponentBranch reindexes the component branch
func (r *Renderer) reindexComponentBranch(v *component.VNode) {
	var walk func(*component.VNode)
	walk = func(v *component.VNode) {
		if v == nil {
			return
		}
		if v.Component != nil {
			r.compIndex[v.Component] = v
		}
		for _, ch := range v.Children {
			walk(ch)
		}
	}
	walk(v)
}

// replaceSubtree replaces the subtree with the new subtree
func (r *Renderer) replaceSubtree(root *component.VNode, oldSub *component.VNode, newSub *component.VNode) bool {
	if root == nil {
		return false
	}
	for i, ch := range root.Children {
		if ch == oldSub {
			root.Children[i] = newSub
			return true
		}
		if r.replaceSubtree(ch, oldSub, newSub) {
			return true
		}
	}
	return false
}

// filterDirtyRoots filters the dirty roots
func (r *Renderer) filterDirtyRoots() map[component.Component]struct{} {
	// Simple version : no explicit parent links; trat all as roots.
	// It the component always use AddChildComponent and we index all boundaries,
	// we can detect ancestors by waking up via compIndex-reverse map if added.
	out := make(map[component.Component]struct{}, len(r.dirtySet))
	for c := range r.dirtySet {
		out[c] = struct{}{}
	}
	return out
}

func logPatch(parent js.Value, currentVNode *component.VNode, newVNode *component.VNode) {
	fmt.Printf("\n\n")
	fmt.Printf("Patching DOM from %p to %p\n", currentVNode, newVNode)
	fmt.Printf("Parent: %+v\n", parent.Get("tagName").String())
	fmt.Printf("Current VNode: %+v ; pointer %p\n", currentVNode, currentVNode)
	fmt.Printf("New VNode: %+v ; pointer %p\n", newVNode, newVNode)
}

// Patch the DOM from old to new
// This is the main algorithm for the virtual DOM diffing
// Here remain most of the efficiciency for the virtual DOM
// TODO: Implement a better algorithm for keyed lists
func (r *Renderer) Patch(parent js.Value, currentVNode *component.VNode, newVNode *component.VNode) {
	//logPatch(parent, currentVNode, newVNode)

	// Cas 1: Création
	if currentVNode == nil && newVNode != nil {
		domNode := r.createDomNode(newVNode)
		parent.Call("appendChild", domNode)
		return
	}

	// Cas 2: Suppression
	if currentVNode != nil && newVNode == nil {
		parent.Call("removeChild", currentVNode.Element)
		return
	}

	// Cas 3: Both nil - nothing to do
	if currentVNode == nil && newVNode == nil {
		return
	}

	// Cas 4: Remplacement
	if currentVNode.Tag != newVNode.Tag || currentVNode.Type != newVNode.Type {
		newDomNode := r.createDomNode(newVNode)
		parent.Call("replaceChild", newDomNode, currentVNode.Element)
		return
	}

	// --- LE PASSAGE DE TÉMOIN CRUCIAL ---
	// Si nous sommes ici, les nœuds sont du même type.
	// On passe la référence DOM de l'ancien nœud au nouveau.
	newVNode.Element = currentVNode.Element

	// Mise à jour d'un nœud texte
	if newVNode.Type == component.VNodeText {
		if currentVNode.Text != newVNode.Text {
			newVNode.Element.Set("textContent", newVNode.Text)
		}
	} else {
		// Mise à jour des propriétés, du style, et des enfants pour un nœud élément
		r.updateProps(currentVNode, newVNode)
		r.updateStyle(currentVNode, newVNode) // LA PIÈCE MANQUANTE
		r.patchChildren(currentVNode, newVNode)
	}
}

func (r *Renderer) updateStyle(oldVNode, newVNode *component.VNode) {
	oldStyle := oldVNode.AppliedStyle
	newStyle := newVNode.AppliedStyle
	element := newVNode.Element

	// Si le style n'a pas changé, on ne fait rien.
	if oldStyle == newStyle {
		return
	}

	// Gérer les classes existantes pour ne pas les effacer
	currentClasses := element.Get("className").String()
	classMap := make(map[string]bool)
	for _, c := range strings.Fields(currentClasses) {
		classMap[c] = true
	}

	// Retirer l'ancienne classe de style si elle existait
	if oldStyle != nil {
		delete(classMap, oldStyle.GetClassName())
	}

	// Ajouter la nouvelle classe de style si elle existe
	if newStyle != nil {
		newClassName := newStyle.GetClassName()
		// Injecter le CSS si nécessaire
		if !r.injectedClasses[newClassName] {
			css := newStyle.ToCSS()
			currentCSS := r.styleElement.Get("innerHTML").String()
			r.styleElement.Set("innerHTML", currentCSS+css)
			r.injectedClasses[newClassName] = true
		}
		classMap[newClassName] = true
	}

	// Reconstruire la chaîne de classes et l'appliquer avec setAttribute
	var finalClasses []string
	for c := range classMap {
		finalClasses = append(finalClasses, c)
	}
	element.Call("setAttribute", "class", strings.Join(finalClasses, " "))
}

// updateProps gère la mise à jour des attributs HTML d'un élément.
// Il ne gère PAS le style, qui est traité séparément.
func (r *Renderer) updateProps(oldVNode, newVNode *component.VNode) {
	oldProps := oldVNode.Attrs
	newProps := newVNode.Attrs
	element := newVNode.Element // On utilise la référence du nouveau VNode

	// 1. Supprimer les anciennes propriétés qui n'existent plus dans les nouvelles.
	for key := range oldProps {
		if _, exists := newProps[key]; !exists {
			element.Call("removeAttribute", key)
		}
	}

	// 2. Ajouter ou modifier les nouvelles propriétés.
	for key, newValue := range newProps {
		oldValue := oldProps[key]

		// On ne touche au DOM que si la valeur a réellement changé.
		if oldValue != newValue {
			// NOTE : Pour une robustesse maximale, vous pourriez ajouter un switch
			// sur le type de `newValue` pour gérer les booléens, les nombres, etc.
			// Pour l'instant, une conversion en string est un bon début.
			element.Call("setAttribute", key, fmt.Sprintf("%v", newValue))
		}
	}
}

// patchChildren gère la réconciliation des enfants d'un élément en se basant sur leur index.
func (r *Renderer) patchChildren(oldVNode, newVNode *component.VNode) {
	oldChildren := oldVNode.Children
	newChildren := newVNode.Children
	parent := newVNode.Element

	// Déterminer la longueur de la plus longue des deux listes d'enfants.
	maxLen := len(oldChildren)
	if len(newChildren) > maxLen {
		maxLen = len(newChildren)
	}

	// Parcourir et "patcher" chaque enfant.
	for i := 0; i < maxLen; i++ {
		var oldChild, newChild *component.VNode

		// Obtenir l'ancien enfant s'il existe à cet index.
		if i < len(oldChildren) {
			oldChild = oldChildren[i]
		}

		// Obtenir le nouvel enfant s'il existe à cet index.
		if i < len(newChildren) {
			newChild = newChildren[i]
		}

		// Appeler Patch récursivement. Cette fonction gérera automatiquement :
		// - La création (si oldChild est nil)
		// - La suppression (si newChild est nil)
		// - La mise à jour (si les deux existent)
		r.Patch(parent, oldChild, newChild)
	}
}

func (r *Renderer) createDomNode(vnode *component.VNode) js.Value {
	if vnode == nil {
		return js.Null()
	}

	document := js.Global().Get("document")

	switch vnode.Type {
	case component.VNodeText:
		textNode := document.Call("createTextNode", vnode.Text)
		vnode.Element = textNode
		return textNode

	case component.VNodeElement:
		element := document.Call("createElement", vnode.Tag)
		vnode.Element = element

		// Set properties
		for key, value := range vnode.Attrs {
			if key == "style" {
				// Handle inline styles
				element.Get("style").Set("cssText", value)
			} else {
				element.Call("setAttribute", key, value)
			}
		}

		// Add event listeners
		for event, handler := range vnode.EventHandlers {
			element.Call("addEventListener", event, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				handler(args[0])
				return nil
			}))
		}

		// Process the CSS-in-Go style
		r.processStyle(vnode)

		// Append children
		for _, child := range vnode.Children {
			childNode := r.createDomNode(child)
			if childNode.Truthy() {
				element.Call("appendChild", childNode)
			}
		}

		return element
	}

	return js.Null()
}

// processStyle inject the styles
func (r *Renderer) processStyle(vnode *component.VNode) {
	if vnode.Type != component.VNodeElement {
		fmt.Printf("WARNING: processStyle called on a non-element vnode %v\n", vnode.Type)
	}

	// If the style is nil, we return
	if vnode.AppliedStyle == nil {
		return
	}

	className := vnode.AppliedStyle.GetClassName()

	// Inject the CSS if it has not already be done
	if !r.injectedClasses[className] {
		css := vnode.AppliedStyle.ToCSS()
		currentCSS := r.styleElement.Get("innerHTML").String()
		r.styleElement.Set("innerHTML", currentCSS+css)
		r.injectedClasses[className] = true
	}

	// Ajoute la class à l'élément
	vnode.Element.Call("setAttribute", "class", className)
}

// RequestFrame requests an animation frame for smooth rendering
func (r *Renderer) RequestFrame() {
	if r.animationFrame.Truthy() {
		js.Global().Call("cancelAnimationFrame", r.animationFrame)
	}

	r.animationFrame = js.Global().Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Frame callback - can be used for post-render operations
		return nil
	}))
}

// GetContainer returns the container element
func (r *Renderer) GetContainer() js.Value {
	return r.container
}
