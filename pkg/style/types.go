//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the core interfaces and types for the style package.
//
// In Vortex for absolute safety, each CSS property is checked and validated at compile time.
// This is done to avoid runtime errors and to ensure that the CSS properties are valid.
// This is done to ensure that the CSS properties are valid and to avoid runtime errors.
// Pros :
// - Absolute safety
// - Compile time errors
// - Type safety
// Cons :
// - More verbose
// - More complex
// - More difficult to write
// - More difficult to read
//
// This file includes:
// - Core interfaces for CSS values
// - Length value type with validation
// - Utility functions for CSS property management

// Basic Usage:
// Check the other files for more information and examples
//

package style

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// === CORE INTERFACES ===

// CSSValue represents a CSS value that can be applied to a CSS property
// It is used to validate the CSS value and to get the CSS value as a string
type CSSValue interface {
	// String returns the CSS value as a string
	String() string
	// Validate checks if the CSS value is valid
	Validate() error
}

// ValidationError represents a CSS value validation error
type ValidationError struct {
	Property string
	Value    string
	Reason   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid CSS %s value '%s': %s", e.Property, e.Value, e.Reason)
}

// Bacth validation returning errors
func BatchValidateWithErrors(propertyName string, values ...CSSValue) error {
	for _, value := range values {
		validateCSSValue(propertyName, value)
	}
	return &ValidationError{
		Property: propertyName,
		Value:    "",
		Reason:   "invalid CSS value",
	}
}

// Bacth validation without returning errors
func BatchValidate(propertyName string, values ...CSSValue) {
	for _, value := range values {
		if err := value.Validate(); err != nil {
			log.Printf("CSS validation warning for %s: %v", propertyName, err)
		}
	}
}

// Validation Helper function to avoid duplicate code
func validateCSSValue(propertyName string, value CSSValue) {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning for %s: %v", propertyName, err)
	}
}

// Cast a list of CSSValues to a list of strings
func CSSValuesToString[T CSSValue](values ...T) []string {
	valuesStrings := make([]string, len(values))
	for i, value := range values {
		valuesStrings[i] = value.String()
	}
	return valuesStrings
}

// ===== KEYWORD VALUES =====
type KeywordValue struct {
	value string
}

func (k KeywordValue) String() string  { return k.value }
func (k KeywordValue) Validate() error { return nil }

// Keyword constants
var (
	Auto    = KeywordValue{value: "auto"}
	Inherit = KeywordValue{value: "inherit"}
	Initial = KeywordValue{value: "initial"}
	Revert  = KeywordValue{value: "revert"}
	Unset   = KeywordValue{value: "unset"}
)

// ===== LENGTH VALUES =====

type LengthUnit string

const (
	// Absolute units
	UnitPx LengthUnit = "px"
	UnitPt LengthUnit = "pt"
	UnitPc LengthUnit = "pc"
	UnitIn LengthUnit = "in"
	UnitMm LengthUnit = "mm"
	UnitCm LengthUnit = "cm"

	// Relative units
	UnitEm  LengthUnit = "em"
	UnitRem LengthUnit = "rem"
	UnitEx  LengthUnit = "ex"
	UnitCh  LengthUnit = "ch"

	// Viewport units
	UnitVw   LengthUnit = "vw"
	UnitVh   LengthUnit = "vh"
	UnitVmin LengthUnit = "vmin"
	UnitVmax LengthUnit = "vmax"

	// Percentage
	UnitPercent LengthUnit = "%"
)

type LengthValue struct {
	Value float64
	Unit  LengthUnit
}

func (l LengthValue) String() string {
	return fmt.Sprintf("%.2f%s", l.Value, l.Unit)
}

func (l LengthValue) Validate() error {
	var validUnits = []LengthUnit{
		UnitPx, UnitPt, UnitPc, UnitIn, UnitMm, UnitCm,
		UnitEm, UnitRem, UnitEx, UnitCh,
		UnitVw, UnitVh, UnitVmin, UnitVmax,
		UnitPercent,
	}

	for _, unit := range validUnits {
		if l.Unit == unit {
			return nil
		}
	}

	return &ValidationError{
		Property: "length",
		Value:    l.String(),
		Reason:   fmt.Sprintf("invalid unit '%s'", l.Unit),
	}
}

// Constructor functions for length value
func Px(value float64) LengthValue      { return LengthValue{Value: value, Unit: UnitPx} }
func Pt(value float64) LengthValue      { return LengthValue{Value: value, Unit: UnitPt} }
func Pc(value float64) LengthValue      { return LengthValue{Value: value, Unit: UnitPc} }
func In(value float64) LengthValue      { return LengthValue{Value: value, Unit: UnitIn} }
func Mm(value float64) LengthValue      { return LengthValue{Value: value, Unit: UnitMm} }
func Cm(value float64) LengthValue      { return LengthValue{Value: value, Unit: UnitCm} }
func Em(value float64) LengthValue      { return LengthValue{Value: value, Unit: UnitEm} }
func Rem(value float64) LengthValue     { return LengthValue{Value: value, Unit: UnitRem} }
func Ex(value float64) LengthValue      { return LengthValue{Value: value, Unit: UnitEx} }
func Ch(value float64) LengthValue      { return LengthValue{Value: value, Unit: UnitCh} }
func Vw(value float64) LengthValue      { return LengthValue{Value: value, Unit: UnitVw} }
func Vh(value float64) LengthValue      { return LengthValue{Value: value, Unit: UnitVh} }
func Vmin(value float64) LengthValue    { return LengthValue{Value: value, Unit: UnitVmin} }
func Vmax(value float64) LengthValue    { return LengthValue{Value: value, Unit: UnitVmax} }
func Percent(value float64) LengthValue { return LengthValue{Value: value, Unit: UnitPercent} }

// ===== COLOR VALUES =====

type ColorValue struct {
	Value string
}

func (c ColorValue) String() string {
	return c.Value
}

func (c ColorValue) Validate() error {
	return validateColorString(c.Value)
}

// Color Constructor functions
func RGB(r, g, b int) ColorValue {
	return ColorValue{Value: fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)}
}
func RGBA(r, g, b int, a float64) ColorValue {
	return ColorValue{Value: fmt.Sprintf("rgba(%d, %d, %d, %0.2f)", r, g, b, a)}
}
func HSL(h, s, l int) ColorValue {
	return ColorValue{Value: fmt.Sprintf("hsl(%d, %d, %d)", h, s, l)}
}

func HEX(color string) ColorValue {
	// Ensure # prefix
	if !strings.HasPrefix(color, "#") {
		color = "#" + color
	}
	return ColorValue{Value: color}
}

// CSS Color Keywords
var (
	ColorTransparent  = ColorValue{Value: "transparent"}
	ColorCurrentColor = ColorValue{Value: "currentColor"}
	ColorInherit      = ColorValue{Value: "inherit"}
	ColorInitial      = ColorValue{Value: "initial"}

	// Common colors
	ColorBlack              = ColorValue{Value: "black"}
	ColorWhite              = ColorValue{Value: "white"}
	ColorRed                = ColorValue{Value: "red"}
	ColorGreen              = ColorValue{Value: "green"}
	ColorBlue               = ColorValue{Value: "blue"}
	ColorYellow             = ColorValue{Value: "yellow"}
	ColorPurple             = ColorValue{Value: "purple"}
	ColorOrange             = ColorValue{Value: "orange"}
	ColorPink               = ColorValue{Value: "pink"}
	ColorBrown              = ColorValue{Value: "brown"}
	ColorGray               = ColorValue{Value: "gray"}
	ColorSilver             = ColorValue{Value: "silver"}
	ColorGold               = ColorValue{Value: "gold"}
	ColorMaroon             = ColorValue{Value: "maroon"}
	ColorNavy               = ColorValue{Value: "navy"}
	ColorTeal               = ColorValue{Value: "teal"}
	ColorOlive              = ColorValue{Value: "olive"}
	ColorLime               = ColorValue{Value: "lime"}
	ColorAqua               = ColorValue{Value: "aqua"}
	ColorFuchsia            = ColorValue{Value: "fuchsia"}
	ColorIndigo             = ColorValue{Value: "indigo"}
	ColorViolet             = ColorValue{Value: "violet"}
	ColorRebeccaPurple      = ColorValue{Value: "rebecca-purple"}
	ColorYellowGreen        = ColorValue{Value: "yellow-green"}
	ColorTurquoise          = ColorValue{Value: "turquoise"}
	ColorSkyBlue            = ColorValue{Value: "sky-blue"}
	ColorLightBlue          = ColorValue{Value: "light-blue"}
	ColorLightGreen         = ColorValue{Value: "light-green"}
	ColorLightRed           = ColorValue{Value: "light-red"}
	ColorLightYellow        = ColorValue{Value: "light-yellow"}
	ColorLightPurple        = ColorValue{Value: "light-purple"}
	ColorLightOrange        = ColorValue{Value: "light-orange"}
	ColorLightPink          = ColorValue{Value: "light-pink"}
	ColorLightBrown         = ColorValue{Value: "light-brown"}
	ColorLightGray          = ColorValue{Value: "light-gray"}
	ColorLightSilver        = ColorValue{Value: "light-silver"}
	ColorLightGold          = ColorValue{Value: "light-gold"}
	ColorLightMaroon        = ColorValue{Value: "light-maroon"}
	ColorLightNavy          = ColorValue{Value: "light-navy"}
	ColorLightTeal          = ColorValue{Value: "light-teal"}
	ColorLightOlive         = ColorValue{Value: "light-olive"}
	ColorLightLime          = ColorValue{Value: "light-lime"}
	ColorLightAqua          = ColorValue{Value: "light-aqua"}
	ColorLightFuchsia       = ColorValue{Value: "light-fuchsia"}
	ColorLightIndigo        = ColorValue{Value: "light-indigo"}
	ColorLightViolet        = ColorValue{Value: "light-violet"}
	ColorLightRebeccaPurple = ColorValue{Value: "light-rebecca-purple"}
	ColorLightYellowGreen   = ColorValue{Value: "light-yellow-green"}
	ColorLightTurquoise     = ColorValue{Value: "light-turquoise"}
	ColorLightSkyBlue       = ColorValue{Value: "light-sky-blue"}
)

func validateColorString(color string) error {
	// Hex colors (#fff, #ffffff)
	hexPattern := regexp.MustCompile(`^#[0-9a-fA-F]{3}([0-9a-fA-F]{3})?$`)
	if hexPattern.MatchString(color) {
		return nil
	}

	// RGB/RGBA colors (rgb(255, 255, 255), rgba(255, 255, 255, 1))
	rgbPattern := regexp.MustCompile(`^rgba?\(\s*\d+\s*,\s*\d+\s*,\s*\d+\s*(,\s*[\d.]+)?\s*\)$`)
	if rgbPattern.MatchString(color) {
		return nil
	}

	// HSL/HSLA colors (hsl(255, 255, 255), hsla(255, 255, 255, 1))
	hslPattern := regexp.MustCompile(`^hsl\((\d{1,3}),\s*(\d{1,3}),\s*(\d{1,3})(,\s*(\d{1,3}))?\)$`)
	if hslPattern.MatchString(color) {
		return nil
	}

	// CSS Color Keywords
	cssColorKeywords := []string{
		"transparent",
		"currentColor",
		"inherit",
		"initial",
		"black",
		"white",
		"red",
		"green",
		"blue",
		"yellow",
		"purple",
		"orange",
		"pink",
		"brown",
		"gray",
		"silver",
		"gold",
		"maroon",
		"navy",
		"teal",
		"olive",
		"lime",
		"aqua",
		"fuchsia",
		"indigo",
		"violet",
		"rebecca-purple",
		"yellow-green",
		"turquoise",
		"sky-blue",
		"light-blue",
		"light-green",
		"light-red",
		"light-yellow",
		"light-purple",
		"light-orange",
		"light-pink",
		"light-brown",
		"light-gray",
		"light-silver",
		"light-gold",
		"light-maroon",
		"light-navy",
		"light-teal",
		"light-olive",
		"light-lime",
		"light-aqua",
		"light-fuchsia",
		"light-indigo",
		"light-violet",
		"light-rebecca-purple",
		"light-yellow-green",
		"light-turquoise",
		"light-sky-blue",
		"light-blue-green",
		"light-blue-green",
	}
	for _, keyword := range cssColorKeywords {
		if keyword == color {
			return nil
		}
	}
	return &ValidationError{
		Property: "color",
		Value:    color,
		Reason:   "invalid color value",
	}
}

// validateLengthString validates a length string
func validateLengthString(value string) error {
	// Special keywords
	if value == "auto" || value == "inherit" || value == "initial" || value == "revert" || value == "unset" {
		return nil
	}

	// Length pattern : number + unit
	lengthPattern := regexp.MustCompile(`^\d+(\.\d+)?(px|pt|pc|in|mm|cm|em|rem|ex|ch|vw|vh|vmin|vmax|%)$`)
	if lengthPattern.MatchString(value) {
		return nil
	}

	// Calc() expressions (basic validations)
	if strings.HasPrefix(value, "calc(") && strings.HasSuffix(value, ")") {
		return nil //TODO: More complex validation
	}

	return &ValidationError{
		Property: "length",
		Value:    value,
		Reason:   "invalid length value",
	}
}

// valdidateFontWeightString validates a font weight string
func validateFontWeightString(value string) error {
	// Numeric values
	if weight, err := strconv.Atoi(value); err == nil {
		if weight >= 100 && weight <= 900 && weight%100 == 0 {
			return nil
		}
	}

	// Keyword values
	validKeywords := []string{
		"normal",
		"bold",
		"lighter",
		"bolder",
		"inherit",
		"initial",
		"revert",
		"unset",
	}
	for _, keyword := range validKeywords {
		if keyword == value {
			return nil
		}
	}

	return &ValidationError{
		Property: "font-weight",
		Value:    value,
		Reason:   "must be 100-900 (multiple of 100) or a keyword (normal, bold, lighter, bolder, inherit, initial, revert, unset)",
	}
}

// ===== FONT WEIGHT TYPES ====
type FontWeightValue struct {
	value string
}

func (f FontWeightValue) String() string  { return f.value }
func (f FontWeightValue) Validate() error { return validateFontWeightString(f.value) }

// Font weight constants
var (
	FontWeightNormal  = FontWeightValue{value: "normal"}
	FontWeightBold    = FontWeightValue{value: "bold"}
	FontWeightLighter = FontWeightValue{value: "lighter"}
	FontWeightBolder  = FontWeightValue{value: "bolder"}

	FontWeight100 = FontWeightValue{value: "100"}
	FontWeight200 = FontWeightValue{value: "200"}
	FontWeight300 = FontWeightValue{value: "300"}
	FontWeight400 = FontWeightValue{value: "400"}
	FontWeight500 = FontWeightValue{value: "500"}
	FontWeight600 = FontWeightValue{value: "600"}
	FontWeight700 = FontWeightValue{value: "700"}
	FontWeight800 = FontWeightValue{value: "800"}
	FontWeight900 = FontWeightValue{value: "900"}
)

// ===== VALIDATION UTILITIES =====
func ValidateCSS(property string, value string) error {
	switch property {
	case "width", "height", "margin", "margin-top", "margin-right", "margin-bottom", "margin-left":
	case "padding", "padding-top", "padding-right", "padding-bottom", "padding-left":
	case "border-width", "border-radius":
		return validateLengthString(value)

	case "color", "background-color", "border-color":
		return validateColorString(value)

	case "font-weight":
		return validateFontWeightString(value)

	case "font-size":
		return validateLengthString(value) // Font size follows length rules

	default:
		return nil // Allow unknown properties (future CSS features)
	}
	return nil
}

// ===== HELPER FUNCTIONS =====

// ParseLength attempt to parse a string as a length value
func ParseLength(value string) (CSSValue, error) {
	if err := validateLengthString(value); err != nil {
		return nil, err
	}

	// Special keywords
	switch value {
	case "auto":
		return Auto, nil
	case "inherit":
		return Inherit, nil
	case "initial":
		return Initial, nil
	case "revert":
		return Revert, nil
	case "unset":
		return Unset, nil
	}

	// Return a validated string for now
	// TODO: Parse into proper LengthValue struct
	return KeywordValue{value}, nil
}

// ParseColor attempts to parse a string as a color value
func ParseColor(value string) (CSSValue, error) {
	if err := validateColorString(value); err != nil {
		return ColorValue{}, err
	}
	return ColorValue{value}, nil
}
