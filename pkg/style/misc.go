//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply the other properties to a style.
// It is used to apply the other properties to a style.
//
// Basic Usage:
//
//   style := style.New()
//       .Cursor(style.CursorPointer).
//       .ZIndex(1).
//       .Overflow(style.OverflowAxisX, style.OverflowTypeScroll).
//       .Position(style.PositionStatic)
//   )
//
// For more information, see the style package documentation
//
// All properties available are (table with css equivalent)
// | Property | CSS Equivalent |
// | -------- | -------------- |
// | Cursor | cursor        |
// | Z-Index | z-index        |
// | Overflow | overflow        |
// | Position | position        |
//
// For more information, see the style package documentation

package style

import "fmt"

// --- Other

// CursorValue is a type that represents a cursor value
// Usage examples :
//
//	style.Cursor(style.CursorPointer)
type CursorValue string

const (
	CursorPointer    CursorValue = "pointer"
	CursorWait       CursorValue = "wait"
	CursorCrosshair  CursorValue = "crosshair"
	CursorDefault    CursorValue = "default"
	CursorNotAllowed CursorValue = "not-allowed"
	CursorGrab       CursorValue = "grab"
	CursorGrabbing   CursorValue = "grabbing"
	CursorHelp       CursorValue = "help"
	CursorMove       CursorValue = "move"
)

func (c CursorValue) String() string {
	return string(c)
}

func (c CursorValue) Validate() error {
	return ValidateCSS("cursor", string(c))
}

// Cursor is a function that applies a cursor to the element
// Usage examples :
//
//	style.Cursor(style.CursorPointer)
func (s *Style) Cursor(value CursorValue) *Style {
	s.Base["cursor"] = value.String()
	return s
}

// ZIndex is a function that applies a z-index to the element
// Usage examples :
//
//	style.ZIndex(style.ZIndex(1))
func (s *Style) ZIndex(value int) *Style {
	s.Base["z-index"] = fmt.Sprintf("%d", value)
	return s
}

// OverflowAxis is a type that represents an overflow axis
// Usage examples :
//
//	style.OverflowAxisX
type OverflowAxis string

const (
	OverflowAxisX    OverflowAxis = "x"
	OverflowAxisY    OverflowAxis = "y"
	OverflowAxisBoth OverflowAxis = "both"
)

type OverflowValue string

func (o OverflowValue) String() string {
	return string(o)
}

func (o OverflowValue) Validate() error {
	return ValidateCSS("overflow", string(o))
}

// OverflowValue is a type that represents an overflow value
// Usage examples :
//
// style.OverflowVisible
type OverflowType string

const (
	OverflowVisible OverflowType = "visible"
	OverflowHidden  OverflowType = "hidden"
	OverflowScroll  OverflowType = "scroll"
	OverflowAuto    OverflowType = "auto"
)

// Overflow is a function that applies an overflow to the element
// Usage examples :
//
//	style.Overflow(style.OverflowAxisX, style.OverflowTypeScroll)
//	style.Overflow(style.OverflowAxisY, style.OverflowTypeScroll)
//	style.Overflow(style.OverflowAxisBoth, style.OverflowTypeHidden)
func (s *Style) Overflow(axis OverflowAxis, value OverflowType) *Style {
	suffix := ""
	if axis == OverflowAxisX {
		suffix = "-x"
	}
	if axis == OverflowAxisY {
		suffix = "-y"
	}
	s.Base["overflow"+suffix] = string(value)
	return s
}

// PositionValue is a type that represents a position value
// Usage examples :
//
// style.PositionStatic
// style.PositionRelative
// style.PositionAbsolute
type PositionValue string

const (
	PositionStatic   PositionValue = "static"
	PositionRelative PositionValue = "relative"
	PositionAbsolute PositionValue = "absolute"
	PositionFixed    PositionValue = "fixed"
	PositionSticky   PositionValue = "sticky"
)

func (p PositionValue) String() string {
	return string(p)
}

func (p PositionValue) Validate() error {
	return ValidateCSS("position", string(p))
}

// Position is a function that applies a position to the element
// Usage examples :
//
// style.Position(style.PositionStatic)
// style.Position(style.PositionRelative)
// style.Position(style.PositionAbsolute)
func (s *Style) Position(value PositionValue) *Style {
	s.Base["position"] = value.String()
	return s
}
