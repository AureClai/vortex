//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply the flexbox to a style.
// It is used to apply the flexbox to a style.
//
// Basic Usage:
//
//   style := style.New(
//       style.Flex(style.Flex1),
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
func Flex(value FlexValue) StyleOption {
	return func(s *Style) {
		s.Base["flex"] = string(value)
	}
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
func FlexDirection(value FlexDirectionValue) StyleOption {
	return func(s *Style) {
		s.Base["flex-direction"] = string(value)
	}
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
func FlexWrap(value FlexWrapValue) StyleOption {
	return func(s *Style) {
		s.Base["flex-wrap"] = string(value)
	}
}

// FlexGrow is a function that applies a flex grow to the element
// Usage examples :
//
//	style.FlexGrow(style.FlexInt(1))
//	style.FlexGrow(style.FlexAuto)
//	style.FlexGrow(style.FlexNone)
func FlexGrow(value FlexValue) StyleOption {
	return func(s *Style) {
		s.Base["flex-grow"] = string(value)
	}
}

// FlexShrink is a function that applies a flex shrink to the element
// Usage examples :
//
//	style.FlexShrink(style.FlexInt(1))
//	style.FlexShrink(style.FlexAuto)
//	style.FlexShrink(style.FlexNone)
func FlexShrink(value FlexValue) StyleOption {
	return func(s *Style) {
		s.Base["flex-shrink"] = string(value)
	}
}

// FlexBasis is a function that applies a flex basis to the element
// Usage examples :
//
//	style.FlexBasis(style.FlexInt(1))
//	style.FlexBasis(style.FlexAuto)
//	style.FlexBasis(style.FlexNone)
func FlexBasis(value CSSValue) StyleOption { // Accept both LengthValue and FlexValue
	return func(s *Style) {
		s.Base["flex-basis"] = value.String()
	}
}

// Convenience functions
func FlexBasisPx(value float64) StyleOption      { return FlexBasis(Px(value)) }
func FlexBasisPercent(value float64) StyleOption { return FlexBasis(Percent(value)) }
func FlexBasisAuto() StyleOption                 { return FlexBasis(FlexAuto) }

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
func AlignItems(value AlignItemsValue) StyleOption {
	return func(s *Style) {
		s.Base["align-items"] = string(value)
	}
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
func JustifyContent(value JustifyContentValue) StyleOption {
	return func(s *Style) {
		s.Base["justify-content"] = string(value)
	}
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
func Gap(value GapValue) StyleOption {
	return func(s *Style) {
		s.Base["gap"] = string(value)
	}
}
