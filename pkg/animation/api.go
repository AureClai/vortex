//go:build js && wasm

package animation

import (
	"fmt"
	"strings"
	"syscall/js"
	"time"

	"github.com/AureClai/vortex/pkg/vdom"
)

// Engine manages all animated components using the Graph system
type Engine struct {
	graphs      map[*vdom.VNode]*Graph
	isRunning   bool
	requestID   js.Value
	frameCount  int64
	lastLogTime int64
}

// GraphBuilder is used to build a new graph (state machine) for a given VNode
type GraphBuilder struct {
	componentVNode *vdom.VNode
	graph          *Graph
}

// NodeBuilder is used to build a new node for a given graph
type NodeBuilder struct {
	graphBuilder *GraphBuilder
	nodeName     string
	clip         *Clip
}

// ClipBuilder is used to build a new clip for a given node
type ClipBuilder struct {
	nodeBuilder *NodeBuilder
	clip        *Clip
}

// ScaleKeyframe represents a scale animation keyframe
type ScaleKeyframe struct {
	Time  time.Duration
	Value float32
}

// RotationKeyframe represents a rotation animation keyframe
type RotationKeyframe struct {
	Time  time.Duration
	Value float32
}

var globalEngine *Engine

func init() {
	globalEngine = &Engine{
		graphs: make(map[*vdom.VNode]*Graph),
	}
	globalEngine.Start()
	// fmt.Println("🎬 Animation engine initialized") // Only this one
}

// Animate creates a new graph (state machine) for a given VNode
func Animate(component *vdom.VNode) *GraphBuilder {
	fmt.Println("🎬 Animation.Animate() called for component")
	return &GraphBuilder{
		componentVNode: component,
		graph:          NewGraph(),
	}
}

// AddNode adds a new node to the graph
func (gb *GraphBuilder) AddNode(nodeName string) *NodeBuilder {
	return &NodeBuilder{
		graphBuilder: gb,
		nodeName:     nodeName,
	}
}

// WithClip defines the animation clip for this state
func (nb *NodeBuilder) WithClip(duration time.Duration) *ClipBuilder {
	clip := &Clip{
		Name:      nb.nodeName,
		Duration:  duration,
		Tracks:    make(map[string]*Track),
		IsLooping: true,              // Default to looping
		FrameRate: 60.0,              // Default to 60fps
		Events:    make([]*Event, 0), // Default to no events
	}

	nb.clip = clip

	return &ClipBuilder{
		nodeBuilder: nb,
		clip:        clip,
	}
}

// Scale adds scale animation keyframes
func (cb *ClipBuilder) Scale(keyframes ...ScaleKeyframe) *ClipBuilder {
	track := &Track{
		Property:          "scale",
		Keyframes:         make([]*Keyframe, len(keyframes)),
		InterpolationType: InterpolationTypeBezier,
	}

	for i, kf := range keyframes {
		track.Keyframes[i] = &Keyframe{
			Time:  kf.Time,
			Value: kf.Value,
		}
	}

	cb.clip.Tracks[track.Property] = track
	return cb
}

// Rotation adds rotation animation keyframes
func (cb *ClipBuilder) Rotation(keyframes ...RotationKeyframe) *ClipBuilder {
	track := &Track{
		Property:          "rotate",
		Keyframes:         make([]*Keyframe, len(keyframes)),
		InterpolationType: InterpolationTypeBezier,
	}

	for i, kf := range keyframes {
		track.Keyframes[i] = &Keyframe{
			Time:  kf.Time,
			Value: kf.Value,
		}
	}

	cb.clip.Tracks[track.Property] = track
	return cb
}

// Translation adds translation animation keyframes
func (cb *ClipBuilder) Translation(keyframes ...TranslationKeyframe) *ClipBuilder {
	track := &Track{
		Property:          "translate",
		Keyframes:         make([]*Keyframe, len(keyframes)),
		InterpolationType: InterpolationTypeBezier,
	}

	for i, kf := range keyframes {
		track.Keyframes[i] = &Keyframe{
			Time:  kf.Time,
			Value: kf.Value,
		}
	}

	cb.clip.Tracks[track.Property] = track
	return cb
}

// Opacity adds opacity animation keyframes
func (cb *ClipBuilder) Opacity(keyframes ...OpacityKeyframe) *ClipBuilder {
	track := &Track{
		Property:          "opacity",
		Keyframes:         make([]*Keyframe, len(keyframes)),
		InterpolationType: InterpolationTypeBezier,
	}

	for i, kf := range keyframes {
		track.Keyframes[i] = &Keyframe{
			Time:  kf.Time,
			Value: kf.Value,
		}
	}

	cb.clip.Tracks[track.Property] = track
	return cb
}

// TranslationKeyframe represents a translation keyframe
type TranslationKeyframe struct {
	Time  time.Duration
	Value [2]float32 // X, Y translation
}

// OpacityKeyframe represents an opacity keyframe
type OpacityKeyframe struct {
	Time  time.Duration
	Value float32
}

func (cb *ClipBuilder) Done() *GraphBuilder {
	fmt.Printf("🎬 Creating node '%s' with clip duration: %v\n", cb.nodeBuilder.nodeName, cb.clip.Duration)

	// Create the Node (state) with the clip
	node := &Node{
		Name:  cb.nodeBuilder.nodeName,
		Clip:  cb.clip,
		Loop:  cb.clip.IsLooping,
		Speed: 1.0,
		OnUpdate: func(deltatime time.Duration) {
			// This is where the animation logic is executed
			cb.nodeBuilder.graphBuilder.applyAnimationToComponent()
		},
	}

	// Add the node to the graph
	cb.nodeBuilder.graphBuilder.graph.AddNode(node)
	fmt.Printf("🎬 Node '%s' added to graph. Total nodes: %d\n", node.Name, len(cb.nodeBuilder.graphBuilder.graph.Nodes))

	return cb.nodeBuilder.graphBuilder
}

// AddEdge adds a new edge to the graph
func (gb *GraphBuilder) AddEdge(from, to string, condition EdgeCondition) *GraphBuilder {
	edge := &Edge{
		From:      from,
		To:        to,
		Condition: condition,
		Duration:  time.Millisecond * 200, // Default to 200ms
		BlendType: BlendLinear,
	}

	gb.graph.AddEdge(edge)
	return gb
}

// StartWith sets the initial node and starts the animation graph
func (gb *GraphBuilder) StartWith(initialNode string) {
	fmt.Printf("🎬 StartWith(%s) called\n", initialNode)

	err := gb.graph.SetNode(initialNode)
	if err != nil {
		fmt.Printf("❌ Failed to set initial node: %v\n", err)
		return
	}

	gb.graph.IsPlaying = true
	fmt.Printf("🎬 Graph IsPlaying set to: %t\n", gb.graph.IsPlaying)

	// Register with the global engine
	globalEngine.graphs[gb.componentVNode] = gb.graph
	fmt.Printf("🎬 Animation registered with engine. Total graphs: %d\n", len(globalEngine.graphs))
}

// applyAnimationToComponent applies the animation to the component
func (gb *GraphBuilder) applyAnimationToComponent() {
	if gb.graph.CurrentNode == nil || gb.graph.CurrentNode.Clip == nil {
		return
	}

	// Get current animated property values
	properties := gb.calculateAnimatedProperties()

	// Apply to DOM element
	gb.applyPropertiesToDOM(properties)
}

// calculateAnimatedProperties calculates the animated properties for the current node
func (gb *GraphBuilder) calculateAnimatedProperties() map[string]interface{} {
	properties := make(map[string]interface{})

	node := gb.graph.CurrentNode
	if node == nil {
		return properties
	}

	currentClip := node.Clip
	if currentClip == nil {
		return properties
	}

	localTime := node.LocalTime

	// Handle looping
	if currentClip.IsLooping && localTime > currentClip.Duration {
		localTime = time.Duration(int64(localTime) % int64(currentClip.Duration))
	}

	// Calculate each property's current value
	for propertyName, track := range currentClip.Tracks {
		value := gb.interpolateProperty(propertyName, track, localTime)
		properties[propertyName] = value
	}

	return properties
}

// interpolateProperty interpolates a property value at a given time
func (gb *GraphBuilder) interpolateProperty(propertyName string, track *Track, currentTime time.Duration) interface{} {
	if len(track.Keyframes) == 0 {
		return nil
	}

	// Find surrounding keyframes
	var leftKf, rightKf *Keyframe

	for i, kf := range track.Keyframes {
		if kf.Time <= currentTime {
			leftKf = kf
			if i+1 < len(track.Keyframes) {
				rightKf = track.Keyframes[i+1]
			}
		} else {
			break
		}
	}

	// If before first keyframe
	if leftKf == nil {
		return track.Keyframes[0].Value
	}

	// If after last keyframe
	if rightKf == nil {
		return leftKf.Value
	}

	// Interpolate between keyframes
	t := float32(currentTime-leftKf.Time) / float32(rightKf.Time-leftKf.Time)
	return gb.lerpValue(leftKf.Value, rightKf.Value, t)
}

// lerpValue linearly interpolates between two values
func (gb *GraphBuilder) lerpValue(from, to interface{}, t float32) interface{} {
	switch from.(type) {
	case float32:
		return from.(float32)*(1-t) + to.(float32)*t
	case [2]float32:
		fromVec := from.([2]float32)
		toVec := to.([2]float32)
		return [2]float32{
			fromVec[0]*(1-t) + toVec[0]*t,
			fromVec[1]*(1-t) + toVec[1]*t,
		}
	default:
		return from
	}
}

// applyPropertiesToDOM applies animated properties to the DOM element
func (gb *GraphBuilder) applyPropertiesToDOM(properties map[string]interface{}) {
	if gb.componentVNode == nil || gb.componentVNode.Element.IsUndefined() {
		// Element not yet rendered to DOM, skip animation frame
		return
	}

	element := gb.componentVNode.Element
	if element.IsUndefined() {
		return
	}

	style := element.Get("style")
	if style.IsUndefined() {
		return
	}

	// Batch transform properties into one update
	var transformParts []string

	for propName, value := range properties {
		switch propName {
		case "rotate":
			if rotation, ok := value.(float32); ok {
				transformParts = append(transformParts, fmt.Sprintf("rotate(%.2fdeg)", rotation))
			}
		case "scale":
			if scale, ok := value.(float32); ok {
				transformParts = append(transformParts, fmt.Sprintf("scale(%.6f)", scale))
			}
		case "translate":
			if trans, ok := value.([2]float32); ok {
				transformParts = append(transformParts, fmt.Sprintf("translate(%.2fpx, %.2fpx)", trans[0], trans[1]))
			}
		case "opacity":
			if opacity, ok := value.(float32); ok {
				style.Set("opacity", fmt.Sprintf("%.6f", opacity))
			}
		}
	}

	// Apply all transforms in one DOM update
	if len(transformParts) > 0 {
		style.Set("transform", strings.Join(transformParts, " "))
	}
}

// The engine loop that drives all the graphs
func (e *Engine) Start() {
	if e.isRunning {
		return
	}

	e.isRunning = true
	e.scheduleFrame()
}

func (e *Engine) scheduleFrame() {
	if !e.isRunning {
		return
	}

	e.requestID = js.Global().Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e.update()
		e.scheduleFrame()
		return nil
	}))
}

func (e *Engine) update() {
	deltaTime := time.Millisecond * 16 // ~60fps

	// Early return if no graphs
	if len(e.graphs) == 0 {
		return
	}

	// Performance monitoring every (5 seconds)
	e.frameCount++
	now := time.Now().Unix()
	if now-e.lastLogTime >= 5 {
		fps := float64(e.frameCount) / float64(now-e.lastLogTime)
		fmt.Printf("🎬 Engine update FPS: %.1f, Active graphs: %d\n", fps, len(e.graphs))
		e.frameCount = 0
		e.lastLogTime = now
	}

	// Update all graphs (remove logging)
	for _, graph := range e.graphs {
		if graph.IsPlaying {
			graph.Update(deltaTime)
		}
	}
}

// GetActiveGraphCount returns the number of active graphs
func GetActiveGraphCount() int {
	if globalEngine == nil {
		return 0
	}
	return len(globalEngine.graphs)
}
