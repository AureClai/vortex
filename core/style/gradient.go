//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply the gradient to a style.
// It is used to apply the gradient to a style.
//
// Basic Usage:
//
//   style := style.New(
//       style.BackgroundColorGradient(style.GradientTypeLinear, style.GradientDirectionToTop, style.ColorValue("#000000"), style.ColorValue("#FFFFFF")),
//   )

//
// For more information, see the style package documentation

package style

import (
	"fmt"
	"sort"
	"strings"
)

type GradientStop struct {
	color    ColorValue
	position float64
}

func (c GradientStop) String() string {
	return fmt.Sprintf("%s %0.2f%%", c.color.String(), c.position*100)
}

func (c GradientStop) Validate() error {
	if c.position < 0 || c.position > 1 {
		return &ValidationError{
			Property: "position",
			Value:    fmt.Sprintf("%0.2f", c.position),
			Reason:   "position must be between 0 and 1",
		}
	}
	if err := c.color.Validate(); err != nil {
		return err
	}
	return nil
}

// Gradient is a type that represents a gradient
// Usage examples :
//
//	style.NewGradient(style.GradientTypeLinear, style.GradientDirectionToTop, style.ColorValue("#000000"), style.ColorValue("#FFFFFF"))
type Gradient struct {
	Type      GradientType
	Direction GradientDirection
	Stops     []GradientStop
}

func (g Gradient) String() string {
	return fmt.Sprintf("%s-gradient(%s , %s)", g.Type, g.Direction, strings.Join(CSSValuesToString(g.Stops...), ","))
}

func (g Gradient) Validate() error {
	return ValidateCSS("background", g.String())
}

// Constructor Basic
func NewGradient(gType GradientType, gDirection GradientDirection, gStops ...GradientStop) Gradient {
	return Gradient{Type: gType, Direction: gDirection, Stops: gStops}
}

type GradientBuilder struct {
	gradient Gradient
	offset   float64
}

// Fluent API
func NewGradientBuilder() *GradientBuilder {
	return &GradientBuilder{gradient: Gradient{}, offset: 0}
}

func (b *GradientBuilder) Type(gType GradientType) *GradientBuilder {
	b.gradient.Type = gType
	return b
}

func (b *GradientBuilder) Direction(gDirection GradientDirection) *GradientBuilder {
	b.gradient.Direction = gDirection
	return b
}

func (b *GradientBuilder) AddStop(color ColorValue, position float64) *GradientBuilder {
	b.gradient.Stops = append(b.gradient.Stops, GradientStop{color: color, position: position})
	return b
}

func (b *GradientBuilder) Offset(offset float64) *GradientBuilder {
	b.offset = offset
	return b
}

func (b *GradientBuilder) Build() Gradient {

	if len(b.gradient.Stops) < 2 {
		fmt.Printf("warning : at least 2 stops are required, returning black gradient\n")
		return Gradient{Type: GradientTypeLinear, Direction: GradientDirectionToTop, Stops: []GradientStop{{color: ColorValue{Value: "#000000"}, position: 0}, {color: ColorValue{Value: "#000000"}, position: 1}}}
	}
	if b.offset < 0 || b.offset > 1 {
		fmt.Printf("warning : offset must be between 0 and 1, returning black gradient\n")
		return Gradient{Type: GradientTypeLinear, Direction: GradientDirectionToTop, Stops: []GradientStop{{color: ColorValue{Value: "#000000"}, position: 0}, {color: ColorValue{Value: "#000000"}, position: 1}}}
	}

	if b.offset != 0 {
		var err error
		b.gradient, err = applyOffset(b.gradient, b.offset)
		if err != nil {
			fmt.Printf("warning : failed to apply offset, returning black gradient\n")
			return Gradient{Type: GradientTypeLinear, Direction: GradientDirectionToTop, Stops: []GradientStop{{color: ColorValue{Value: "#000000"}, position: 0}, {color: ColorValue{Value: "#000000"}, position: 1}}}
		}
	}

	if err := b.gradient.Validate(); err != nil {
		fmt.Printf("warning : failed to validate gradient, returning black gradient\n")
		return Gradient{Type: GradientTypeLinear, Direction: GradientDirectionToTop, Stops: []GradientStop{{color: ColorValue{Value: "#000000"}, position: 0}, {color: ColorValue{Value: "#000000"}, position: 1}}}
	}

	return b.gradient
}

func applyOffset(gradient Gradient, offset float64) (Gradient, error) {
	// Ensure stops are sorted by position, as the logic depends on it.
	sort.Slice(gradient.Stops, func(i, j int) bool {
		return gradient.Stops[i].position < gradient.Stops[j].position
	})

	if len(gradient.Stops) == 0 {
		return Gradient{}, fmt.Errorf("cannot apply offset to a gradient with no stops")
	}

	// ===== STEP 1: Find the color at the precise offset point (with interpolation) =====

	var offsetColor ColorValue
	var err error

	// Case 1: Offset is before the first stop.
	if offset <= gradient.Stops[0].position {
		offsetColor = gradient.Stops[0].color
		// Case 2: Offset is after the last stop.
	} else if offset >= gradient.Stops[len(gradient.Stops)-1].position {
		offsetColor = gradient.Stops[len(gradient.Stops)-1].color
		// Case 3: Offset is between two stops, so we need to interpolate.
	} else {
		var stopIn, stopOut GradientStop
		// Find the two stops that bracket the offset.
		for i := 0; i < len(gradient.Stops)-1; i++ {
			if offset >= gradient.Stops[i].position && offset < gradient.Stops[i+1].position {
				stopIn = gradient.Stops[i]
				stopOut = gradient.Stops[i+1]
				break
			}
		}

		// Calculate how far the offset is between the two stops (a value from 0 to 1).
		relativeOffset := (offset - stopIn.position) / (stopOut.position - stopIn.position)
		offsetColor, err = ColorValueInterpolation(stopIn.color, stopOut.color, relativeOffset)
		if err != nil {
			return Gradient{}, fmt.Errorf("failed to interpolate offset color: %w", err)
		}
	}

	// ===== STEP 2: Perform the cyclic shift using the found offsetColor =====

	var newStops []GradientStop

	// Process stops at or after the offset
	for _, stop := range gradient.Stops {
		if stop.position >= offset {
			newStops = append(newStops, GradientStop{
				color:    stop.color,
				position: stop.position - offset, // Shift left to the beginning
			})
		}
	}

	// Process stops before the offset and wrap them around to the end
	for _, stop := range gradient.Stops {
		if stop.position < offset {
			newStops = append(newStops, GradientStop{
				color:    stop.color,
				position: stop.position + (1.0 - offset), // Shift right to the end
			})
		}
	}

	// Add the new starting stop at 0%
	newStops = append(newStops, GradientStop{color: offsetColor, position: 0.0})
	// Ensure the end color matches the new start color for a seamless loop
	newStops = append(newStops, GradientStop{color: offsetColor, position: 1.0})

	// Sort the re-assembled stops by their new positions
	sort.Slice(newStops, func(i, j int) bool {
		return newStops[i].position < newStops[j].position
	})

	finalStops := mergeDuplicateStops(newStops)

	// NOTE: You might want to merge duplicate positions here for a cleaner gradient.

	// Return a new gradient with the transformed stops
	return Gradient{
		Type:      gradient.Type,
		Direction: gradient.Direction,
		Stops:     finalStops,
	}, nil
}

// mergeDuplicateStops cleans a sorted slice of stops, removing duplicates.
// It keeps the last-seen color for any given position.
func mergeDuplicateStops(stops []GradientStop) []GradientStop {
	if len(stops) < 2 {
		return stops
	}

	cleanStops := make([]GradientStop, 0)
	cleanStops = append(cleanStops, stops[0])

	for i := 1; i < len(stops); i++ {
		lastCleanStop := &cleanStops[len(cleanStops)-1]
		currentStop := stops[i]

		if currentStop.position == lastCleanStop.position {
			// Overwrite the last stop with this new one at the same position
			*lastCleanStop = currentStop
		} else {
			// Append as a new stop
			cleanStops = append(cleanStops, currentStop)
		}
	}
	return cleanStops
}

// GradientType is a type that represents the type of gradient
// Usage examples :
//
//	style.GradientTypeLinear
//	style.GradientTypeRadial
//	style.GradientTypeConic
type GradientType string

const (
	GradientTypeLinear GradientType = "linear"
	GradientTypeRadial GradientType = "radial"
	GradientTypeConic  GradientType = "conic"
)

// GradientDirection is a type that represents the direction of the gradient
// Usage examples :
//
//	style.GradientDirectionToTop
//	style.GradientDirectionToBottom
type GradientDirection string

const (
	GradientDirectionToTop          GradientDirection = "to top"
	GradientDirectionToBottom       GradientDirection = "to bottom"
	GradientDirectionToLeft         GradientDirection = "to left"
	GradientDirectionToRight        GradientDirection = "to right"
	GradientDirectionToTopLeft      GradientDirection = "to top left"
	GradientDirectionToTopRight     GradientDirection = "to top right"
	GradientDirectionToBottomLeft   GradientDirection = "to bottom left"
	GradientDirectionToBottomRight  GradientDirection = "to bottom right"
	GradientDirectionToCenter       GradientDirection = "to center"
	GradientDirectionToCenterLeft   GradientDirection = "to center left"
	GradientDirectionToCenterRight  GradientDirection = "to center right"
	GradientDirectionToCenterTop    GradientDirection = "to center top"
	GradientDirectionToCenterBottom GradientDirection = "to center bottom"
	GradientDirectionToLeftTop      GradientDirection = "to left top"
	GradientDirectionToLeftBottom   GradientDirection = "to left bottom"
	GradientDirectionToRightTop     GradientDirection = "to right top"
	GradientDirectionToRightBottom  GradientDirection = "to right bottom"
)

// Deg is a function that returns a GradientDirection with the given value in degrees
// Usage examples :
//
//	style.Deg(45)
func Deg(value int) GradientDirection {
	return GradientDirection(fmt.Sprintf("%ddeg", value))
}
