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
	"strings"
)

// Gradient is a type that represents a gradient
// Usage examples :
//
//	style.NewGradient(style.GradientTypeLinear, style.GradientDirectionToTop, style.ColorValue("#000000"), style.ColorValue("#FFFFFF"))
type Gradient struct {
	Type      GradientType
	Direction GradientDirection
	Colors    []ColorValue
}

func (g Gradient) String() string {
	return fmt.Sprintf("%s-gradient(%s , %s)", g.Type, g.Direction, strings.Join(CSSValuesToString(g.Colors...), ","))
}

func NewGradient(gType GradientType, gDirection GradientDirection, gColors ...ColorValue) Gradient {
	return Gradient{Type: gType, Direction: gDirection, Colors: gColors}
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
