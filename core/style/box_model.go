//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply the box model to a style.
// It is used to apply the box model to a style.
// API fluent like the other style functions
//
// Basic Usage:
//
//   style := style.New().
// 	Display(style.DisplayFlex).
// 	WidthPx(200).
// 	HeightAuto().
// 	Margin(style.MarginAll, style.Px(10))
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

func (s *Style) Display(value DisplayValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["display"] = value.String()
	return s
}

// Usage examples :
//
//	style.Width(style.Px(10))
//	style.Width(style.Cm(20))
//	style.Width(style.Em(30))
//	style.Width(style.Auto)
func (s *Style) Width(value CSSValue) *Style {
	// if value is a LengthValue, we need to validate it
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["width"] = value.String()
	return s
}

func (s *Style) MinWidth(value CSSValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["min-width"] = value.String()
	return s
}

func (s *Style) MaxWidth(value CSSValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["max-width"] = value.String()
	return s
}

// Typed Width
func (s *Style) WidthPx(value float64) *Style      { return s.Width(Px(value)) }
func (s *Style) WidthEm(value float64) *Style      { return s.Width(Em(value)) }
func (s *Style) WidthRem(value float64) *Style     { return s.Width(Rem(value)) }
func (s *Style) WidthVw(value float64) *Style      { return s.Width(Vw(value)) }
func (s *Style) WidthVh(value float64) *Style      { return s.Width(Vh(value)) }
func (s *Style) WidthVmin(value float64) *Style    { return s.Width(Vmin(value)) }
func (s *Style) WidthVmax(value float64) *Style    { return s.Width(Vmax(value)) }
func (s *Style) WidthPercent(value float64) *Style { return s.Width(Percent(value)) }

func (s *Style) WidthAuto() *Style    { return s.Width(Auto) }
func (s *Style) WidthInherit() *Style { return s.Width(Inherit) }
func (s *Style) WidthInitial() *Style { return s.Width(Initial) }
func (s *Style) WidthRevert() *Style  { return s.Width(Revert) }
func (s *Style) WidthUnset() *Style   { return s.Width(Unset) }

// Usage examples :
//
//	style.Height(style.Px(10))
//	style.Height(style.Cm(20))
//	style.Height(style.Em(30))
//	style.Height(style.Auto)
func (s *Style) Height(value CSSValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["height"] = value.String()
	return s

}

func (s *Style) MinHeight(value CSSValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["min-height"] = value.String()
	return s
}

func (s *Style) MaxHeight(value CSSValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["max-height"] = value.String()
	return s
}

// Typed Height
func (s *Style) HeightPx(value float64) *Style      { return s.Height(Px(value)) }
func (s *Style) HeightEm(value float64) *Style      { return s.Height(Em(value)) }
func (s *Style) HeightRem(value float64) *Style     { return s.Height(Rem(value)) }
func (s *Style) HeightVw(value float64) *Style      { return s.Height(Vw(value)) }
func (s *Style) HeightVh(value float64) *Style      { return s.Height(Vh(value)) }
func (s *Style) HeightVmin(value float64) *Style    { return s.Height(Vmin(value)) }
func (s *Style) HeightVmax(value float64) *Style    { return s.Height(Vmax(value)) }
func (s *Style) HeightPercent(value float64) *Style { return s.Height(Percent(value)) }
func (s *Style) HeightAuto() *Style                 { return s.Height(Auto) }
func (s *Style) HeightInherit() *Style              { return s.Height(Inherit) }
func (s *Style) HeightInitial() *Style              { return s.Height(Initial) }
func (s *Style) HeightRevert() *Style               { return s.Height(Revert) }
func (s *Style) HeightUnset() *Style                { return s.Height(Unset) }

// Position based on the top, right, bottom, left
type PositionDirection string

const (
	PositionTop    PositionDirection = "top"
	PositionRight  PositionDirection = "right"
	PositionBottom PositionDirection = "bottom"
	PositionLeft   PositionDirection = "left"
)

func (s *Style) PositionSide(direction PositionDirection, value LengthValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base[string(direction)] = value.String()
	return s
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
func (s *Style) Margin(direction MarginDirection, value CSSValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
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
	return s
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
func (s *Style) Padding(direction PaddingDirection, value CSSValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
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
	return s
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
func (s *Style) BorderWidth(value LengthValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["border-width"] = value.String()
	return s

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
func (s *Style) BorderColor(value ColorValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["border-color"] = value.String()
	return s
}

// Usage examples :
//
//	style.Border(style.Px(10), style.BorderSolid, style.RGB(0, 0, 0))
//	style.Border(style.Cm(20), style.BorderDashed, style.RGB(0, 0, 0))
//	style.Border(style.Em(30), style.BorderDotted, style.RGB(0, 0, 0))
func (s *Style) Border(width LengthValue, style BorderStyleType, color ColorValue) *Style {
	if err := width.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	if err := color.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["border-width"] = width.String()
	s.Base["border-style"] = string(style)
	s.Base["border-color"] = color.String()
	return s

}

// Border Radius
type BorderRadiusValue struct {
	TopLeft     LengthValue
	TopRight    LengthValue
	BottomRight LengthValue
	BottomLeft  LengthValue
}

func BorderRadiusAll(value LengthValue) BorderRadiusValue {
	return BorderRadiusValue{
		TopLeft:     value,
		TopRight:    value,
		BottomRight: value,
		BottomLeft:  value,
	}
}

func NewBorderRadiusValue(topLeft, topRight, bottomRight, bottomLeft LengthValue) BorderRadiusValue {
	return BorderRadiusValue{
		TopLeft:     topLeft,
		TopRight:    topRight,
		BottomRight: bottomRight,
		BottomLeft:  bottomLeft,
	}
}

func (b BorderRadiusValue) String() string {
	return fmt.Sprintf("%s %s %s %s", b.TopLeft.String(), b.TopRight.String(), b.BottomRight.String(), b.BottomLeft.String())
}
func (b BorderRadiusValue) Validate() error {
	return BatchValidateWithErrors("border-radius", b.TopLeft, b.TopRight, b.BottomRight, b.BottomLeft)
}

// Usage examples :
//
//	style.BorderRadius(style.BorderRadiusAll(style.Px(10)))
//	style.BorderRadius(style.NewBorderRadiusValue(style.Px(10), style.Px(10), style.Px(10), style.Px(10)))
func (s *Style) BorderRadius(value BorderRadiusValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["border-radius"] = value.String()
	return s
}

// box-sizing
type BoxSizingType string

const (
	BoxSizingContentBox BoxSizingType = "content-box"
	BoxSizingBorderBox  BoxSizingType = "border-box"
)

func (b BoxSizingType) String() string {
	return string(b)
}
func (b BoxSizingType) Validate() error {
	return ValidateCSS("box-sizing", string(b))
}

// Usage examples :
//
//	style.BoxSizing(style.BoxSizingContentBox)
//	style.BoxSizing(style.BoxSizingBorderBox)
func (s *Style) BoxSizing(value BoxSizingType) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["box-sizing"] = value.String()
	return s

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
func (s *Style) BoxShadow(value ...BoxShadowValue) *Style {
	// Cast to CSSValue
	cssValues := make([]CSSValue, len(value))
	for i, v := range value {
		cssValues[i] = v
	}
	BatchValidate("box-shadow", cssValues...)

	// Join the shadow styles with ,\n to be compatoible with multiple shadows
	shadowStyles := make([]string, len(value))
	for i, v := range value {
		shadowStyles[i] = fmt.Sprintf("%s %s %s %s %s", v.OffsetX.String(), v.OffsetY.String(), v.BlurRadius.String(), v.SpreadRadius.String(), v.Color.String())
	}
	s.Base["box-shadow"] = strings.Join(shadowStyles, ",\n")
	return s
}
