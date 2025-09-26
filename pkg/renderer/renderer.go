//go:build js && wasm

package renderer

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/AureClai/vortex/pkg/vdom"
)

type Renderer struct {
	container    js.Value    // Container element
	currentVNode *vdom.VNode // Current virtual node state for patching algorithms

	styleElement    js.Value        // Ref to the <style> balise
	injectedClasses map[string]bool // Keep track from the classes already injected
	animationFrame  js.Value
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
	}
}

func (r *Renderer) Render(newVNode *vdom.VNode) {
	// Clear container and render new tree
	r.Patch(r.container, r.currentVNode, newVNode)
	r.currentVNode = newVNode
}

func logPatch(parent js.Value, currentVNode *vdom.VNode, newVNode *vdom.VNode) {
	fmt.Printf("\n\n")
	fmt.Printf("Parent: %+v\n", parent.Get("tagName").String())
	fmt.Printf("Current VNode: %+v ; pointer %p\n", currentVNode, currentVNode)
	fmt.Printf("New VNode: %+v ; pointer %p\n", newVNode, newVNode)
}

// Patch the DOM from old to new
// This is the main algorithm for the virtual DOM diffing
// Here remain most of the efficiciency for the virtual DOM
// TODO: Implement a better algorithm for keyed lists
func (r *Renderer) Patch(parent js.Value, currentVNode *vdom.VNode, newVNode *vdom.VNode) {
	logPatch(parent, currentVNode, newVNode)
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

	// Cas 3: Remplacement
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
	if newVNode.Type == vdom.VNodeText {
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

func (r *Renderer) updateStyle(oldVNode, newVNode *vdom.VNode) {
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
func (r *Renderer) updateProps(oldVNode, newVNode *vdom.VNode) {
	oldProps := oldVNode.Props
	newProps := newVNode.Props
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
func (r *Renderer) patchChildren(oldVNode, newVNode *vdom.VNode) {
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
		var oldChild, newChild *vdom.VNode

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

func (r *Renderer) createDomNode(vnode *vdom.VNode) js.Value {
	if vnode == nil {
		return js.Null()
	}

	document := js.Global().Get("document")

	switch vnode.Type {
	case vdom.VNodeText:
		textNode := document.Call("createTextNode", vnode.Text)
		vnode.Element = textNode
		return textNode

	case vdom.VNodeElement:
		element := document.Call("createElement", vnode.Tag)
		vnode.Element = element

		// Set properties
		for key, value := range vnode.Props {
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
func (r *Renderer) processStyle(vnode *vdom.VNode) {
	if vnode.Type != vdom.VNodeElement {
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
