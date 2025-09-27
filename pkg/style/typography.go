//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply the typography to a style.
// It is used to apply the typography to a style.
//
// Basic Usage:
//
//   style := style.New().
//       .FontFamily("Arial").
//       .FontSize(style.Px(16)).
//       .FontWeight(style.FontWeightBold).
//       .Color(style.RGB(255, 255, 255)).
//       .TextAlign(style.TextAlignCenter).
//       .LineHeight(style.Px(1.5)).
//       .LetterSpacing(style.Px(1)).
//       .TextDecoration(style.TextDecorationNone)
//   )
//
// For more information, see the style package documentation
//
// All properties available are (table with css equivalent)
// | Property | CSS Equivalent |
// | -------- | -------------- |
// | Font Family | font-family        |
// | Font Size | font-size        |
// | Font Weight | font-weight        |
// | Color | color        |
// | Text Align | text-align        |
// | Line Height | line-height        |
// | Letter Spacing | letter-spacing        |
// | Text Decoration | text-decoration        |
// | Text Shadow | text-shadow        |
// | Word Spacing | word-spacing        |
// | Font Variant | font-variant        |
// | Font Stretch | font-stretch        |
// | Font Style | font-style        |
//
// Note :
// Multiple shadows compatibility
//
//	style.TextShadow(
//		style.TextShadowValue{OffsetX: style.Px(10), OffsetY: style.Px(10), BlurRadius: style.Px(10), Color: style.RGB(0, 0, 0)},
//		style.TextShadowValue{OffsetX: style.Px(20), OffsetY: style.Px(20), BlurRadius: style.Px(20), Color: style.RGB(0, 0, 0)},
//	)
//
// For more information, see the style package documentation

package style

import (
	"fmt"
	"strings"
)

// --- Typography

// TODO : Add a FontFamilyValue type and validate it but it need further research.
//
// Usage examples :
//
//	style.FontFamily("Arial")
//	style.FontFamily("Helvetica")
//	style.FontFamily("Times New Roman")
func (s *Style) FontFamily(value string) *Style {
	s.Base["font-family"] = value
	return s

}

// Usage examples :
//
//	style.FontSize(style.Px(16))
//	style.FontSize(style.Cm(20))
//	style.FontSize(style.Em(30))
func (s *Style) FontSize(value LengthValue) *Style {
	validateCSSValue("font-size", value)
	s.Base["font-size"] = value.String()
	return s
}

// Usage examples :
//
//	style.FontWeight(style.FontWeightNormal)
//	style.FontWeight(style.FontWeightBold)
//	style.FontWeight(style.FontWeight100)
func (s *Style) FontWeight(value FontWeightValue) *Style {
	validateCSSValue("font-weight", value)
	s.Base["font-weight"] = value.String()
	return s
}

// Color is a function that applies a color to the text
// Usage examples :
//
//	style.Color(style.RGB(255, 255, 255))
//	style.Color(style.Hex("#000000"))
//	style.Color(style.HSL(0, 0, 0))
func (s *Style) Color(value ColorValue) *Style {
	validateCSSValue("color", value)
	s.Base["color"] = value.String()
	return s
}

// ColorGradient is a function that applies a gradient to the color
// Usage examples :
//
//	style.ColorGradient(style.NewGradient(style.GradientTypeLinear, style.GradientDirectionToTop, style.ColorValue("#000000"), style.ColorValue("#FFFFFF")))
func (s *Style) ColorGradient(value Gradient) *Style {
	validateCSSValue("color", value)
	s.Base["color"] = value.String()
	return s
}

type TextAlignValue string

const (
	TextAlignLeft   TextAlignValue = "left"
	TextAlignCenter TextAlignValue = "center"
	TextAlignRight  TextAlignValue = "right"
)

// Usage examples :
//
//	style.TextAlign(style.TextAlignLeft)
//	style.TextAlign(style.TextAlignCenter)
//	style.TextAlign(style.TextAlignRight)
func (s *Style) TextAlign(value TextAlignValue) *Style {
	s.Base["text-align"] = string(value)
	return s
}

// Usage examples :
//
//	style.LineHeight(style.Px(16))
//	style.LineHeight(style.Cm(20))
//	style.LineHeight(style.Em(30))
func (s *Style) LineHeight(value LengthValue) *Style {
	validateCSSValue("line-height", value)
	s.Base["line-height"] = value.String()
	return s
}

// Usage examples :
//
//	style.LetterSpacing(style.Px(16))
//	style.LetterSpacing(style.Cm(20))
//	style.LetterSpacing(style.Em(30))
func (s *Style) LetterSpacing(value LengthValue) *Style {
	validateCSSValue("letter-spacing", value)
	s.Base["letter-spacing"] = value.String()
	return s
}

type TextDecorationValue string

const (
	TextDecorationNone        TextDecorationValue = "none"
	TextDecorationUnderline   TextDecorationValue = "underline"
	TextDecorationOverline    TextDecorationValue = "overline"
	TextDecorationLineThrough TextDecorationValue = "line-through"
)

// Usage examples :
//
//	style.TextDecoration(style.TextDecorationNone)
//	style.TextDecoration(style.TextDecorationUnderline)
//	style.TextDecoration(style.TextDecorationOverline)
//	style.TextDecoration(style.TextDecorationLineThrough)
func (s *Style) TextDecoration(value TextDecorationValue) *Style {
	s.Base["text-decoration"] = string(value)
	return s
}

type TextShadowValue struct {
	OffsetX    LengthValue
	OffsetY    LengthValue
	BlurRadius LengthValue
	Color      ColorValue
}

func (t TextShadowValue) String() string {
	return fmt.Sprintf("%s %s %s %s", t.OffsetX.String(), t.OffsetY.String(), t.BlurRadius.String(), t.Color.String())
}

func (t TextShadowValue) Validate() error {
	return BatchValidateWithErrors("text-shadow", t.OffsetX, t.OffsetY, t.BlurRadius, t.Color)
}

// Improvement : Multiple shadows compatibility
// Usage examples :
// style.TextShadow(
//
//	style.TextShadowValue{OffsetX: style.Px(10), OffsetY: style.Px(10), BlurRadius: style.Px(10), Color: style.RGB(0, 0, 0)},
//	style.TextShadowValue{OffsetX: style.Px(20), OffsetY: style.Px(20), BlurRadius: style.Px(20), Color: style.RGB(0, 0, 0)},
//
// )
func (s *Style) TextShadow(value ...TextShadowValue) *Style {
	// Cast to CSSValue
	cssValues := make([]CSSValue, len(value))
	for i, v := range value {
		cssValues[i] = v
	}
	BatchValidate("text-shadow", cssValues...)

	// Join the shadow styles with ,\n to be compatoible with multiple shadows
	shadowStyles := make([]string, len(value))
	for i, v := range value {
		shadowStyles[i] = fmt.Sprintf("%s %s %s %s", v.OffsetX.String(), v.OffsetY.String(), v.BlurRadius.String(), v.Color.String())
	}
	s.Base["text-shadow"] = strings.Join(shadowStyles, ",\n")
	return s
}

// Usage examples :
//
//	style.WordSpacing(style.Px(10))
//	style.WordSpacing(style.Cm(20))
//	style.WordSpacing(style.Em(30))
func (s *Style) WordSpacing(value LengthValue) *Style {
	validateCSSValue("word-spacing", value)
	s.Base["word-spacing"] = value.String()
	return s
}

type FontVariantValue string

const (
	FontVariantNormal              FontVariantValue = "normal"
	FontVariantSmallCaps           FontVariantValue = "small-caps"
	FontVariantAllSmallCaps        FontVariantValue = "all-small-caps"
	FontVariantPetiteCaps          FontVariantValue = "petite-caps"
	FontVariantAllPetiteCaps       FontVariantValue = "all-petite-caps"
	FontVariantTitlingCaps         FontVariantValue = "titling-caps"
	FontVariantAllTitlingCaps      FontVariantValue = "all-titling-caps"
	FontVariantSlashedZero         FontVariantValue = "slashed-zero"
	FontVariantAllSlashedZero      FontVariantValue = "all-slashed-zero"
	FontVariantNumeric             FontVariantValue = "numeric"
	FontVariantDiagonalFractions   FontVariantValue = "diagonal-fractions"
	FontVariantStackedFractions    FontVariantValue = "stacked-fractions"
	FontVariantOrdinal             FontVariantValue = "ordinal"
	FontVariantRouble              FontVariantValue = "rouble"
	FontVariantOldstyleNumbers     FontVariantValue = "oldstyle-numbers"
	FontVariantLiningNumbers       FontVariantValue = "lining-numbers"
	FontVariantTabularNumbers      FontVariantValue = "tabular-numbers"
	FontVariantProportionalNumbers FontVariantValue = "proportional-numbers"
	FontVariantMonospace           FontVariantValue = "monospace"
	FontVariantKana                FontVariantValue = "kana"
	FontVariantKanaTypeFace        FontVariantValue = "kana-type-face"
	FontVariantProportionalKana    FontVariantValue = "proportional-kana"
	FontVariantInherit             FontVariantValue = "inherit"
	FontVariantInitial             FontVariantValue = "initial"
	FontVariantRevert              FontVariantValue = "revert"
	FontVariantUnset               FontVariantValue = "unset"
)

// Usage examples :
//
//	style.FontVariant(style.FontVariantNormal)
//	style.FontVariant(style.FontVariantSmallCaps)
//	style.FontVariant(style.FontVariantAllSmallCaps)
//	style.FontVariant(style.FontVariantPetiteCaps)
//	style.FontVariant(style.FontVariantAllPetiteCaps)
func (s *Style) FontVariant(value ...FontVariantValue) *Style {
	// Cast to string
	stringValues := make([]string, len(value))
	for i, v := range value {
		stringValues[i] = string(v)
	}
	s.Base["font-variant"] = strings.Join(stringValues, " ")
	return s
}

type FontStretchValue string

const (
	FontStretchNormal         FontStretchValue = "normal"
	FontStretchUltraCondensed FontStretchValue = "ultra-condensed"
	FontStretchExtraCondensed FontStretchValue = "extra-condensed"
	FontStretchCondensed      FontStretchValue = "condensed"
	FontStretchSemiCondensed  FontStretchValue = "semi-condensed"
	FontStretchSemiExpanded   FontStretchValue = "semi-expanded"
	FontStretchExpanded       FontStretchValue = "expanded"
	FontStretchExtraExpanded  FontStretchValue = "extra-expanded"
	FontStretchUltraExpanded  FontStretchValue = "ultra-expanded"
	FontStretchInherit        FontStretchValue = "inherit"
	FontStretchInitial        FontStretchValue = "initial"
	FontStretchRevert         FontStretchValue = "revert"
	FontStretchUnset          FontStretchValue = "unset"
)

func (f FontStretchValue) String() string  { return string(f) }
func (f FontStretchValue) Validate() error { return nil }

func FontStretchPercent(value int) FontStretchValue {
	return FontStretchValue(fmt.Sprintf("%d%%", value))
}

// Usage examples :
//
//	style.FontStretch(style.FontStretchNormal)
//	style.FontStretch(style.FontStretchUltraCondensed)
//	style.FontStretchPercent(20)
func (s *Style) FontStretch(value FontStretchValue) *Style {
	validateCSSValue("font-stretch", value)
	s.Base["font-stretch"] = string(value)
	return s
}

type FontStyleValue string

const (
	FontStyleNormal  FontStyleValue = "normal"
	FontStyleItalic  FontStyleValue = "italic"
	FontStyleOblique FontStyleValue = "oblique"
	FontStyleInherit FontStyleValue = "inherit"
	FontStyleInitial FontStyleValue = "initial"
	FontStyleRevert  FontStyleValue = "revert"
	FontStyleUnset   FontStyleValue = "unset"
)

func (f FontStyleValue) String() string  { return string(f) }
func (f FontStyleValue) Validate() error { return nil }

// Usage examples :
//
//	style.FontStyle(style.FontStyleNormal)
//	style.FontStyle(style.FontStyleItalic)
//	style.FontStyle(style.FontStyleOblique)
func (s *Style) FontStyle(value FontStyleValue) *Style {
	s.Base["font-style"] = string(value)
	return s
}
