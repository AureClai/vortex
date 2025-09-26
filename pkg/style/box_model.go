//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply the box model to a style.
// It is used to apply the box model to a style.
//
// Basic Usage:
//
//   style := style.New(
//       style.Display(style.DisplayFlex),
//       style.WidthPx(200),
//       style.HeightAuto(),
//       style.Margin(MarginAll, style.Px(10)), // 10px margin on all sides
//   )
//
// For more information, see the style package documentation
//
// All properties available are (table with css equivalent)
// | Property | CSS Equivalent |
// | -------- | -------------- |
// | Display  | display        |
// | Width    | width          |
// | Height   | height         |
// | Margin   | margin         |
// | Padding  | padding        |
// | Border   | border         |
// | Box Sizing | box-sizing     |
// | Box Shadow | box-shadow     |
//
// Note :
// Multiple shadows compatibility
//
//	style.BoxShadow(
//		style.BoxShadowValue{OffsetX: style.Px(10), OffsetY: style.Px(10), BlurRadius: style.Px(10), SpreadRadius: style.Px(10), Color: style.RGB(0, 0, 0), IsInset: false},
//		style.BoxShadowValue{OffsetX: style.Px(20), OffsetY: style.Px(20), BlurRadius: style.Px(20), SpreadRadius: style.Px(20), Color: style.RGB(0, 0, 0), IsInset: false},
//	)
//
// For more information, see the style package documentation

package style

import (
	"fmt"
	"log"
	"strings"
)

// 1. Box Model
type DisplayValue string

const (
	DisplayBlock       DisplayValue = "block"
	DisplayInline      DisplayValue = "inline"
	DisplayInlineBlock DisplayValue = "inline-block"
	DisplayNone        DisplayValue = "none"
	DisplayFlex        DisplayValue = "flex"
	DisplayGrid        DisplayValue = "grid"
)

func (d DisplayValue) String() string  { return string(d) }
func (d DisplayValue) Validate() error { return ValidateCSS("display", string(d)) }

func Display(value DisplayValue) StyleOption {
	return func(s *Style) {
		if err := value.Validate(); err != nil {
			log.Printf("CSS validation warning: %v", err)
		}
		s.Base["display"] = value.String()
	}
}

// Usage examples :
//
//	style.Width(style.Px(10))
//	style.Width(style.Cm(20))
//	style.Width(style.Em(30))
//	style.Width(style.Auto)
func Width(value CSSValue) StyleOption {
	// if value is a LengthValue, we need to validate it
	validateCSSValue("width", value)
	return func(s *Style) {
		s.Base["width"] = value.String()
	}
}

// Typed Width
func WidthPx(value float64) StyleOption      { return Width(Px(value)) }
func WidthEm(value float64) StyleOption      { return Width(Em(value)) }
func WidthRem(value float64) StyleOption     { return Width(Rem(value)) }
func WidthVw(value float64) StyleOption      { return Width(Vw(value)) }
func WidthVh(value float64) StyleOption      { return Width(Vh(value)) }
func WidthVmin(value float64) StyleOption    { return Width(Vmin(value)) }
func WidthVmax(value float64) StyleOption    { return Width(Vmax(value)) }
func WidthPercent(value float64) StyleOption { return Width(Percent(value)) }

func WidthAuto() StyleOption    { return Width(Auto) }
func WidthInherit() StyleOption { return Width(Inherit) }
func WidthInitial() StyleOption { return Width(Initial) }
func WidthRevert() StyleOption  { return Width(Revert) }
func WidthUnset() StyleOption   { return Width(Unset) }

// Usage examples :
//
//	style.Height(style.Px(10))
//	style.Height(style.Cm(20))
//	style.Height(style.Em(30))
//	style.Height(style.Auto)
func Height(value CSSValue) StyleOption {
	validateCSSValue("height", value)
	return func(s *Style) {
		s.Base["height"] = value.String()
	}
}

// Typed Height
func HeightPx(value float64) StyleOption      { return Height(Px(value)) }
func HeightEm(value float64) StyleOption      { return Height(Em(value)) }
func HeightRem(value float64) StyleOption     { return Height(Rem(value)) }
func HeightVw(value float64) StyleOption      { return Height(Vw(value)) }
func HeightVh(value float64) StyleOption      { return Height(Vh(value)) }
func HeightVmin(value float64) StyleOption    { return Height(Vmin(value)) }
func HeightVmax(value float64) StyleOption    { return Height(Vmax(value)) }
func HeightPercent(value float64) StyleOption { return Height(Percent(value)) }
func HeightAuto() StyleOption                 { return Height(Auto) }
func HeightInherit() StyleOption              { return Height(Inherit) }
func HeightInitial() StyleOption              { return Height(Initial) }
func HeightRevert() StyleOption               { return Height(Revert) }
func HeightUnset() StyleOption                { return Height(Unset) }

// Position based on the top, right, bottom, left
type PositionDirection string

const (
	PositionTop    PositionDirection = "top"
	PositionRight  PositionDirection = "right"
	PositionBottom PositionDirection = "bottom"
	PositionLeft   PositionDirection = "left"
)

func PositionSide(direction PositionDirection, value LengthValue) StyleOption {
	return func(s *Style) {
		s.Base[string(direction)] = value.String()
	}
}

// margin
type MarginDirection string

const (
	MarginTop    MarginDirection = "top"
	MarginRight  MarginDirection = "right"
	MarginBottom MarginDirection = "bottom"
	MarginLeft   MarginDirection = "left"
	MarginAll    MarginDirection = "all"
	MarginX      MarginDirection = "x"
	MarginY      MarginDirection = "y"
	MarginBlock  MarginDirection = "block"
	MarginInline MarginDirection = "inline"
)

// Usage examples :
//
//	style.Margin(style.MarginTop, style.Px(10))
//	style.Margin(style.MarginRight, style.Cm(20))
//	style.Margin(style.MarginBottom, style.Em(30))
//	style.Margin(style.MarginLeft, style.Px(40))
//	style.Margin(style.MarginAll, style.Px(50))
func Margin(direction MarginDirection, value CSSValue) StyleOption {
	validateCSSValue("margin", value)
	return func(s *Style) {
		switch direction {
		case MarginAll:
			s.Base["margin"] = value.String()
		case MarginX:
			s.Base["margin-left"] = value.String()
			s.Base["margin-right"] = value.String()
		case MarginY:
			s.Base["margin-top"] = value.String()
			s.Base["margin-bottom"] = value.String()
		default:
			s.Base["margin-"+string(direction)] = value.String()
		}
	}
}

// padding
type PaddingDirection string

const (
	PaddingTop    PaddingDirection = "top"
	PaddingRight  PaddingDirection = "right"
	PaddingBottom PaddingDirection = "bottom"
	PaddingLeft   PaddingDirection = "left"
	PaddingAll    PaddingDirection = "all"
	PaddingX      PaddingDirection = "x"
	PaddingY      PaddingDirection = "y"
	PaddingBlock  PaddingDirection = "block"
	PaddingInline PaddingDirection = "inline"
)

// Usage examples :
//
//	style.Padding(style.PaddingTop, style.Px(10))
//	style.Padding(style.PaddingRight, style.Cm(20))
//	style.Padding(style.PaddingBottom, style.Em(30))
//	style.Padding(style.PaddingLeft, style.Px(40))
//	style.Padding(style.PaddingAll, style.Px(50))
func Padding(direction PaddingDirection, value CSSValue) StyleOption {
	validateCSSValue("padding", value)
	return func(s *Style) {
		switch direction {
		case PaddingAll:
			s.Base["padding"] = value.String()
		case PaddingX:
			s.Base["padding-left"] = value.String()
			s.Base["padding-right"] = value.String()
		case PaddingY:
			s.Base["padding-top"] = value.String()
			s.Base["padding-bottom"] = value.String()
		default:
			s.Base["padding-"+string(direction)] = value.String()
		}
	}
}

// border
type BorderStyleType string

const (
	BorderSolid  BorderStyleType = "solid"
	BorderDashed BorderStyleType = "dashed"
	BorderDotted BorderStyleType = "dotted"
	BorderDouble BorderStyleType = "double"
)

// Usage examples :
//
//	style.BorderWidth(style.Px(10))
//	style.BorderWidth(style.Cm(20))
//	style.BorderWidth(style.Em(30))
func BorderWidth(value LengthValue) StyleOption {
	validateCSSValue("border-width", value)
	return func(s *Style) {
		s.Base["border-width"] = value.String()
	}
}

// Usage examples :
//
//	style.BorderStyle(style.BorderSolid)
//	style.BorderStyle(style.BorderDashed)
//	style.BorderStyle(style.BorderDotted)
//	style.BorderStyle(style.BorderDouble)
func BorderStyle(value BorderStyleType) StyleOption {
	return func(s *Style) {
		s.Base["border-style"] = string(value)
	}
}

// Usage examples :
//
//	style.BorderColor(style.RGB(0, 0, 0))
//	style.BorderColor(style.Hex("#000000"))
//	style.BorderColor(style.HSL(0, 0, 0))
func BorderColor(value ColorValue) StyleOption {
	validateCSSValue("border-color", value)
	return func(s *Style) {
		s.Base["border-color"] = value.String()
	}
}

// Usage examples :
//
//	style.Border(style.Px(10), style.BorderSolid, style.RGB(0, 0, 0))
//	style.Border(style.Cm(20), style.BorderDashed, style.RGB(0, 0, 0))
//	style.Border(style.Em(30), style.BorderDotted, style.RGB(0, 0, 0))
func Border(width LengthValue, style BorderStyleType, color ColorValue) StyleOption {
	validateCSSValue("border-width", width)
	validateCSSValue("border-color", color)
	if color.Validate() != nil {
		log.Printf("CSS validation warning: %v", color.Validate())
	}
	return func(s *Style) {
		s.Base["border-width"] = width.String()
		s.Base["border-style"] = string(style)
		s.Base["border-color"] = color.String()
	}
}

// Usage examples :
//
//	style.BorderRadius(style.Px(10))
//	style.BorderRadius(style.Cm(20))
//	style.BorderRadius(style.Em(30))
func BorderRadius(value LengthValue) StyleOption {
	validateCSSValue("border-radius", value)
	return func(s *Style) {
		s.Base["border-radius"] = value.String()
	}
}

// box-sizing
type BoxSizingType string

const (
	BoxSizingContentBox BoxSizingType = "content-box"
	BoxSizingBorderBox  BoxSizingType = "border-box"
)

// Usage examples :
//
//	style.BoxSizing(style.BoxSizingContentBox)
//	style.BoxSizing(style.BoxSizingBorderBox)
func BoxSizing(value BoxSizingType) StyleOption {
	return func(s *Style) {
		s.Base["box-sizing"] = string(value)
	}
}

type BoxShadowValue struct {
	OffsetX      LengthValue
	OffsetY      LengthValue
	BlurRadius   LengthValue
	SpreadRadius LengthValue
	Color        ColorValue
	IsInset      bool
}

func (b BoxShadowValue) String() string {
	return fmt.Sprintf("%s %s %s %s %s", b.OffsetX.String(), b.OffsetY.String(), b.BlurRadius.String(), b.SpreadRadius.String(), b.Color.String())
}
func (b BoxShadowValue) Validate() error {
	return BatchValidateWithErrors("box-shadow", b.OffsetX, b.OffsetY, b.BlurRadius, b.SpreadRadius, b.Color)
}

// box shadow
// Improvement : Multiple shadows compatibility
// Usage :
// style.BoxShadow(
//
//	style.BoxShadowValue{OffsetX: style.Px(10), OffsetY: style.Px(10), BlurRadius: style.Px(10), SpreadRadius: style.Px(10), Color: style.RGB(0, 0, 0), IsInset: false},
//	style.BoxShadowValue{OffsetX: style.Px(20), OffsetY: style.Px(20), BlurRadius: style.Px(20), SpreadRadius: style.Px(20), Color: style.RGB(0, 0, 0), IsInset: false},
//
// )
func BoxShadow(value ...BoxShadowValue) StyleOption {
	// Cast to CSSValue
	cssValues := make([]CSSValue, len(value))
	for i, v := range value {
		cssValues[i] = v
	}
	BatchValidate("box-shadow", cssValues...)

	return func(s *Style) {
		// Join the shadow styles with ,\n to be compatoible with multiple shadows
		shadowStyles := make([]string, len(value))
		for i, v := range value {
			shadowStyles[i] = fmt.Sprintf("%s %s %s %s %s", v.OffsetX.String(), v.OffsetY.String(), v.BlurRadius.String(), v.SpreadRadius.String(), v.Color.String())
		}
		s.Base["box-shadow"] = strings.Join(shadowStyles, ",\n")
	}

}
