//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply the background and appearance to a style.
// It is used to apply the background and appearance to a style.
// API fluent like the other style functions
//
// Basic Usage:
//
//   style := style.New().BackgroundColor(style.RGB(0, 0, 0))
//
// For more information, see the style package documentation
//
// All properties available are (table with css equivalent)
// | Property | CSS Equivalent |
// | -------- | -------------- |
// | Background Repeat | background-repeat |
// | Background Color | background-color |
// | Opacity | opacity |
// | Background Gradient | background |
// | Background Image | background-image |
// | Background Size | background-size |
// | Background Position | background-position |
// | Backdrop Filter | backdrop-filter |
//
// For more information, see the style package documentation

package style

import (
	"fmt"
)

// background and appearance
type BackgroundRepeatValue string

func (b BackgroundRepeatValue) String() string {
	return string(b)
}

func (b BackgroundRepeatValue) Validate() error {
	return ValidateCSS("background-repeat", string(b))
}

const (
	BackgroundRepeatRepeat   BackgroundRepeatValue = "repeat"
	BackgroundRepeatNoRepeat BackgroundRepeatValue = "no-repeat"
	BackgroundRepeatSpace    BackgroundRepeatValue = "space"
	BackgroundRepeatRound    BackgroundRepeatValue = "round"
	BackgroundRepeatY        BackgroundRepeatValue = "no-repeat repeat"
	BackgroundRepeatX        BackgroundRepeatValue = "repeat no-repeat"
	BackgroundRepeatInherit  BackgroundRepeatValue = "inherit"
	BackgroundRepeatInitial  BackgroundRepeatValue = "initial"
	BackgroundRepeatRevert   BackgroundRepeatValue = "revert"
	BackgroundRepeatUnset    BackgroundRepeatValue = "unset"
)

// Usage examples :
//
//	style.BackgroundRepeat2Axes(style.BackgroundRepeatRepeat, style.BackgroundRepeatNoRepeat)
func BackgroundRepeat2Axes(valueX, valueY BackgroundRepeatValue) BackgroundRepeatValue {
	outValueX := valueX
	if !(valueX == BackgroundRepeatRepeat || valueX == BackgroundRepeatNoRepeat || valueX == BackgroundRepeatSpace || valueX == BackgroundRepeatRound) {
		fmt.Println("BackgroundRepeat2Axes: valueX is not a valid BackgroundRepeatValue for X-axis")
		outValueX = BackgroundRepeatNoRepeat
	}
	outValueY := valueY
	if !(valueY == BackgroundRepeatRepeat || valueY == BackgroundRepeatNoRepeat || valueY == BackgroundRepeatSpace || valueY == BackgroundRepeatRound) {
		fmt.Println("BackgroundRepeat2Axes: valueY is not a valid BackgroundRepeatValue for Y-axis")
		outValueY = BackgroundRepeatNoRepeat
	}
	return BackgroundRepeatValue(fmt.Sprintf("%s %s", outValueX, outValueY))
}

// Usage examples :
//
//	style.BackgroundRepeat(style.BackgroundRepeatRepeat)
//	style.BackgroundRepeat(style.BackgroundRepeatNoRepeat)
//	style.BackgroundRepeat(style.BackgroundRepeatSpace)
//
//	style.BackgroundRepeat2Axes(style.BackgroundRepeatSpace, style.BackgroundRepeatNoRepeat)
func (s *Style) BackgroundRepeat(value BackgroundRepeatValue) *Style {
	s.Base["background-repeat"] = string(value)
	return s
}

// Usage examples :
//
//	style.BackgroundColor(style.RGB(0, 0, 0))
//	style.BackgroundColor(style.Hex("#000000"))
//	style.BackgroundColor(style.HSL(0, 0, 0))
func (s *Style) BackgroundColor(value ColorValue) *Style {
	validateCSSValue("background-color", value)
	s.Base["background-color"] = value.String()
	return s
}

// Usage examples :
//
//	style.Opacity(style.Opacity(0.5))
//	style.Opacity(style.Opacity(1))
//	style.Opacity(style.Opacity(0))
func (s *Style) Opacity(value float64) *Style {
	s.Base["opacity"] = fmt.Sprintf("%f", value)
	return s
}

// BackgroundGradient is a function that applies a gradient to the background
// check the Gradient type for more information in gradient.go
// Usage examples :
//
//	style.BackgroundGradient(style.NewGradient(style.GradientTypeLinear, style.GradientDirectionToTop, style.ColorValue("#000000"), style.ColorValue("#FFFFFF")))
func (s *Style) BackgroundGradient(value Gradient) *Style {
	s.Base["background"] = value.String()
	return s
}

// ImageValue is a type that represents an image
// Placeholder for the Image type
// TODO : Add the Image type and the way image are managed in the style package
type ImageValue string

func (i ImageValue) String() string {
	return string(i)
}

func (i ImageValue) Validate() error {
	return ValidateCSS("background-image", string(i))
}

// Placeholder for type safe image value
// TODO : Add the Image type and the way image are managed in the style package
// can be a correct URL or a base64 string
func NewImageValueFromURL(url string) ImageValue {
	return ImageValue(url)
}

// BackgroundImage is a function that applies an image to the background
// Usage examples :
//
//	style.BackgroundImage(style.Image("https://example.com/image.png"))
//
// TODO : Add the Image type and the way image are managed in the style package
// From now use a placeholder type ImageValue
// type ImageValue string
func (s *Style) BackgroundImage(value ImageValue) *Style {
	// No need to validate the image as we use a type safe value
	s.Base["background-image"] = value.String()
	return s
}

type BackgroundSizeValue string

func (b BackgroundSizeValue) String() string {
	return string(b)
}

func (b BackgroundSizeValue) Validate() error {
	return ValidateCSS("background-size", string(b))
}

const (
	BackgroundSizeCover   BackgroundSizeValue = "cover"
	BackgroundSizeContain BackgroundSizeValue = "contain"
	BackgroundSizeAuto    BackgroundSizeValue = "auto"
	BackgroundSizeInitial BackgroundSizeValue = "initial"
	BackgroundSizeInherit BackgroundSizeValue = "inherit"
	BackgroundSizeRevert  BackgroundSizeValue = "revert"
	BackgroundSizeUnset   BackgroundSizeValue = "unset"
)

func BackgroundSizePercent(value float64) BackgroundSizeValue {
	return BackgroundSizeValue(fmt.Sprintf("%f%%", value))
}

// BackgroundSize is a function that applies a size to the background
// Usage examples :
//
//	style.BackgroundSize(style.BackgroundSizeCover)
//	style.BackgroundSize(style.BackgroundSizeContain)
//	style.BackgroundSize(style.BackgroundSizeAuto)
func (s *Style) BackgroundSize(value BackgroundSizeValue) *Style {
	s.Base["background-size"] = value.String()
	return s
}

type SideKeyWordX string
type SideKeyWordY string

const (
	SideTop    SideKeyWordY = "top"
	SideBottom SideKeyWordY = "bottom"
	SideLeft   SideKeyWordX = "left"
	SideRight  SideKeyWordX = "right"
)

type BackgroundPositionValue string

func (b BackgroundPositionValue) String() string {
	return string(b)
}

func (b BackgroundPositionValue) Validate() error {
	return ValidateCSS("background-position", string(b))
}

const (
	BackgroundPositionTop         BackgroundPositionValue = "top"
	BackgroundPositionBottom      BackgroundPositionValue = "bottom"
	BackgroundPositionLeft        BackgroundPositionValue = "left"
	BackgroundPositionRight       BackgroundPositionValue = "right"
	BackgroundPositionCenter      BackgroundPositionValue = "center"
	BackgroundPositionInitial     BackgroundPositionValue = "initial"
	BackgroundPositionInherit     BackgroundPositionValue = "inherit"
	BackgroundPositionRevert      BackgroundPositionValue = "revert"
	BackgroundPositionRevertLayer BackgroundPositionValue = "revert-layer"
	BackgroundPositionUnset       BackgroundPositionValue = "unset"
)

// BackgroundPositionFrom is a function that applies a position to the background
// Usage examples :
//
//	style.BackgroundPositionFrom(style.SideTop, style.Px(10), style.SideLeft, style.Px(10))
//	style.BackgroundPositionFrom(style.SideBottom, style.Px(10), style.SideRight, style.Px(10))
//	style.BackgroundPositionFrom(style.SideCenter, style.Px(10), style.SideCenter, style.Px(10))
func BackgroundPositionFrom(x SideKeyWordX, offsetX LengthValue, y SideKeyWordY, offsetY LengthValue) BackgroundPositionValue {
	return BackgroundPositionValue(fmt.Sprintf("%s %s %s %s", x, offsetX.String(), y, offsetY.String()))
}

// BackgroundPositionFromTopLeft is a function that applies a position to the background
// Usage examples :
//
//	style.BackgroundPositionFromTopLeft(style.Ch(10), style.Em(8))
func BackgroundPositionFromTopLeft(offsetX, offsetY LengthValue) BackgroundPositionValue {
	return BackgroundPositionValue(fmt.Sprintf("%s %s", offsetX.String(), offsetY.String()))
}

func (s *Style) BackgroundPosition(value BackgroundPositionValue) *Style {
	s.Base["background-position"] = value.String()
	return s
}

type BackdropFilterValue string

const (
	BackdropFilterNone    BackdropFilterValue = "none"
	BackdropFilterInherit BackdropFilterValue = "inherit"
	BackdropFilterInitial BackdropFilterValue = "initial"
	BackdropFilterRevert  BackdropFilterValue = "revert"
	BackdropFilterUnset   BackdropFilterValue = "unset"
)

// Usage examples :
//
//	style.BackdropFilterBlur(style.Px(20))
//	style.BackdropFilterNone()
//	style.BackdropFilterInherit()
//	style.BackdropFilterInitial()
//	style.BackdropFilterRevert()
//	style.BackdropFilterUnset()
func BackdropFilterBlur(value LengthValue) BackdropFilterValue {
	return BackdropFilterValue(fmt.Sprintf("blur(%s)", value.String()))
}

func (s *Style) BackdropFilter(value BackdropFilterValue) *Style {
	s.Base["backdrop-filter"] = string(value)
	return s
}
