//go::build js && wasm

package animation

import "time"

// AnimationClip contains keyframe data for a single animation
type Clip struct {
	Name     string
	Duration time.Duration
	Tracks   map[string]*Track
	Events   []*Event

	// Meta data
	FrameRate float32
	IsLooping bool
}

// AnimationTrack represents animation data for a specific property
type Track struct {
	Property  string
	Keyframes []*Keyframe

	// Interpolation settings
	InterpolationType InterpolationType
	PreInfinity       InfinityType
	PostInfinity      InfinityType
}

// Keyframe represents a single animation keyframe
type Keyframe struct {
	Time  time.Duration
	Value interface{} // Could be float32, vec2, vec3, vec4, color, etc.

	// Bezier handles for smooth curves
	InTangent  *TangentHandle
	OutTangent *TangentHandle
}

// TangentHandle represents a control point for Bezier curves
type TangentHandle struct {
	X, Y float32
}

// InterpolationType defines how values are interpolated between keyframes
type InterpolationType int

const (
	InterpolationTypeLinear InterpolationType = iota
	InterpolationTypeStep
	InterpolationTypeBezier
	InterpolationCustom
)

// InfinityType defines how values are extrapolated beyond keyframe times
type InfinityType int

const (
	InfinityTypeConstant InfinityType = iota
	InfinityTypeLinear
	InfinityTypeCycle
	InfinityTypeCycleWithOffset
	InfinityTypeOscillate
)

// AnimationEvent triggers callbacks at specific times
type Event struct {
	Time     time.Duration
	Name     string
	Data     map[string]interface{}
	Callback func(data map[string]interface{})
}
