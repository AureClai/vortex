//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply the flexbox to a style.
// It is used to apply the flexbox to a style.
//
// Basic Usage:
//
//   style := style.New().
//       .Flex(style.Flex1).
//       .FlexDirection(style.FlexDirectionRow).
//       .FlexWrap(style.FlexWrapNowrap).
//       .FlexGrow(style.Flex1).
//       .FlexShrink(style.Flex1).
//       .FlexBasis(style.Flex1).
//       .AlignItems(style.AlignItemsStart).
//       .JustifyContent(style.JustifyContentStart).
//       .Gap(style.NewGapValueFromLengthValue(style.Px(10))).
//   )
//
// For more information, see the style package documentation
//
// All properties available are (table with css equivalent)
// | Property | CSS Equivalent |
// | -------- | -------------- |
// | Flex | flex        |
// | Flex Direction | flex-direction        |
// | Flex Wrap | flex-wrap        |
// | Flex Grow | flex-grow        |
// | Flex Shrink | flex-shrink        |
// | Flex Basis | flex-basis        |
//
// For more information, see the style package documentation

package style

import "fmt"

// FlexValue is a type that represents a flex value
// Usage examples :
//
//	style.Flex(style.FlexAuto)
//	style.Flex(style.FlexNone)
type FlexValue string

const (
	FlexAuto    FlexValue = "auto"    // Auto flex
	FlexNone    FlexValue = "none"    // No flex
	FlexInitial FlexValue = "initial" // Initial flex
	FlexInherit FlexValue = "inherit" // Inherit flex
	FlexUnset   FlexValue = "unset"   // Unset flex
)

func (f FlexValue) String() string {
	return string(f)
}

func (f FlexValue) Validate() error {
	return ValidateCSS("flex", string(f))
}

// FlexInt is a function that applies a flex to the element
// Usage examples :
//
//	style.Flex(style.FlexInt(1))
func FlexInt(value int) FlexValue {
	return FlexValue(fmt.Sprintf("%d", value))
}

// Flex is a function that applies a flex to the element
// Usage examples :
//
//	style.Flex(style.FlexAuto)
//	style.Flex(style.FlexNone)
//	style.Flex(style.FlexInt(1))
func (s *Style) Flex(value FlexValue) *Style {
	s.Base["flex"] = string(value)
	return s

}

type FlexDirectionValue string

const (
	FlexDirectionRow           FlexDirectionValue = "row"
	FlexDirectionColumn        FlexDirectionValue = "column"
	FlexDirectionRowReverse    FlexDirectionValue = "row-reverse"
	FlexDirectionColumnReverse FlexDirectionValue = "column-reverse"
)

// FlexDirection is a function that applies a flex direction to the element
// Usage examples :
//
//	style.FlexDirection(style.FlexDirectionRow)
//	style.FlexDirection(style.FlexDirectionColumn)
//	style.FlexDirection(style.FlexDirectionRowReverse)
//	style.FlexDirection(style.FlexDirectionColumnReverse)
func (s *Style) FlexDirection(value FlexDirectionValue) *Style {
	s.Base["flex-direction"] = string(value)
	return s
}

type FlexWrapValue string

const (
	FlexWrapNowrap      FlexWrapValue = "nowrap"
	FlexWrapWrap        FlexWrapValue = "wrap"
	FlexWrapWrapReverse FlexWrapValue = "wrap-reverse"
)

// FlexWrap is a function that applies a flex wrap to the element
// Usage examples :
//
//	style.FlexWrap(style.FlexWrapNowrap)
//	style.FlexWrap(style.FlexWrapWrap)
//	style.FlexWrap(style.FlexWrapWrapReverse)
func (s *Style) FlexWrap(value FlexWrapValue) *Style {
	s.Base["flex-wrap"] = string(value)
	return s

}

// FlexGrow is a function that applies a flex grow to the element
// Usage examples :
//
//	style.FlexGrow(style.FlexInt(1))
//	style.FlexGrow(style.FlexAuto)
//	style.FlexGrow(style.FlexNone)
func (s *Style) FlexGrow(value FlexValue) *Style {
	s.Base["flex-grow"] = string(value)
	return s

}

// FlexShrink is a function that applies a flex shrink to the element
// Usage examples :
//
//	style.FlexShrink(style.FlexInt(1))
//	style.FlexShrink(style.FlexAuto)
//	style.FlexShrink(style.FlexNone)
func (s *Style) FlexShrink(value FlexValue) *Style {
	s.Base["flex-shrink"] = string(value)
	return s

}

// FlexBasis is a function that applies a flex basis to the element
// Usage examples :
//
//	style.FlexBasis(style.FlexInt(1))
//	style.FlexBasis(style.FlexAuto)
//	style.FlexBasis(style.FlexNone)
func (s *Style) FlexBasis(value CSSValue) *Style { // Accept both LengthValue and FlexValue
	s.Base["flex-basis"] = value.String()
	return s

}

// Convenience functions
func (s *Style) FlexBasisPx(value float64) *Style      { return s.FlexBasis(Px(value)) }
func (s *Style) FlexBasisPercent(value float64) *Style { return s.FlexBasis(Percent(value)) }
func (s *Style) FlexBasisAuto() *Style                 { return s.FlexBasis(FlexAuto) }

type AlignItemsValue string

const (
	AlignItemsStart    AlignItemsValue = "start"
	AlignItemsEnd      AlignItemsValue = "end"
	AlignItemsCenter   AlignItemsValue = "center"
	AlignItemsBaseline AlignItemsValue = "baseline"
	AlignItemsStretch  AlignItemsValue = "stretch"
)

// AlignItems is a function that applies a align items to the element
// Usage examples :
//
//	style.AlignItems(style.AlignItemsStart)
//	style.AlignItems(style.AlignItemsEnd)
//	style.AlignItems(style.AlignItemsCenter)
//	style.AlignItems(style.AlignItemsBaseline)
//	style.AlignItems(style.AlignItemsStretch)
func (s *Style) AlignItems(value AlignItemsValue) *Style {
	s.Base["align-items"] = string(value)
	return s

}

type JustifyContentValue string

const (
	JustifyContentStart        JustifyContentValue = "start"
	JustifyContentEnd          JustifyContentValue = "end"
	JustifyContentCenter       JustifyContentValue = "center"
	JustifyContentBaseline     JustifyContentValue = "baseline"
	JustifyContentStretch      JustifyContentValue = "stretch"
	JustifyContentSpaceBetween JustifyContentValue = "space-between"
	JustifyContentSpaceAround  JustifyContentValue = "space-around"
	JustifyContentSpaceEvenly  JustifyContentValue = "space-evenly"
)

// JustifyContent is a function that applies a justify content to the element
// Usage examples :
//
//	style.JustifyContent(style.JustifyContentStart)
//	style.JustifyContent(style.JustifyContentEnd)
//	style.JustifyContent(style.JustifyContentCenter)
//	style.JustifyContent(style.JustifyContentBaseline)
//	style.JustifyContent(style.JustifyContentStretch)
func (s *Style) JustifyContent(value JustifyContentValue) *Style {
	s.Base["justify-content"] = string(value)
	return s

}

// Gap is a function that applies a gap to the element
// Usage examples :
//	style.Gap(style.NewGapValueFromLengthValue(style.Px(10)))
//	style.Gap(style.NewGapValueFromTwoLengthValues(style.Px(10), style.Em(10)))

type GapValue string

func NewGapValueFromLengthValue(value LengthValue) GapValue {
	return GapValue(value.String())
}

func NewGapValueFromTwoLengthValues(value1, value2 LengthValue) GapValue {
	return GapValue(fmt.Sprintf("%s %s", value1.String(), value2.String()))
}

// Gap is a function that applies a gap to the element
// Usage examples :
//
//	style.Gap(style.NewGapValueFromLengthValue(style.Px(10)))
//	style.Gap(style.NewGapValueFromTwoLengthValues(style.Px(10), style.Em(10)))
func (s *Style) Gap(value GapValue) *Style {
	s.Base["gap"] = string(value)
	return s

}
