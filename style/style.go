package style

import (
	"fmt"
	"strings"
)

// Property is a simple alias for a CSS property
// Ex: "color" -> "blue"
type Property map[string]string

// Style is the object that contains the full definition of a the style of a component
type Style struct {
	Base         Property
	Pseudos      map[string]Property // ":hover", ":active", ":focus", etc.
	MediaQueries map[string]Property // "screen and (max-width: 768px)", "screen and (min-width: 769px)"

	// For caching
	className string
	css       string
}

// StyleOption is a function that modifies the style
type StyleOption func(*Style)

// New creates a new style object applying the given options
func New(options ...StyleOption) *Style {
	s := &Style{
		Base:         make(Property),
		Pseudos:      make(map[string]Property),
		MediaQueries: make(map[string]Property),
	}
	for _, option := range options {
		s.Update(option)
	}
	return s
}

func (s *Style) List() []StyleOption {
	options := []StyleOption{}
	for _, option := range s.Base {
		options = append(options, CustomStyle(option, s.Base[option]))
	}
	for _, pseudo := range s.Pseudos {
		for key, value := range pseudo {
			options = append(options, CustomStyle(key, value))
		}
	}
	for _, mediaQuery := range s.MediaQueries {
		for key, value := range mediaQuery {
			options = append(options, CustomStyle(key, value))
		}
	}
	return options
}

func Extend(baseStyle *Style, options ...StyleOption) *Style {
	// 1. Create a deep copy of the base style
	s := &Style{
		Base:         make(Property),
		Pseudos:      make(map[string]Property),
		MediaQueries: make(map[string]Property),
	}

	// 2. Deep copy the base properties
	for key, value := range baseStyle.Base {
		s.Base[key] = value
	}

	// 3. Deep copy the pseudo properties
	for pseudo, properties := range baseStyle.Pseudos {
		newPseudoProps := make(Property)
		for key, value := range properties {
			newPseudoProps[key] = value
		}
		s.Pseudos[pseudo] = newPseudoProps
	}

	// 4. Deep copy the media queries
	for mediaQuery, properties := range baseStyle.MediaQueries {
		newMediaProps := make(Property)
		for key, value := range properties {
			newMediaProps[key] = value
		}
		s.MediaQueries[mediaQuery] = newMediaProps
	}

	// 5. Apply the options
	s.Update(options...)
	return s
}

func (s *Style) Update(options ...StyleOption) {
	for _, option := range options {
		option(s)
	}
}

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

func Display(value DisplayValue) StyleOption {
	return func(s *Style) {
		s.Base["display"] = string(value)
	}
}

// width
func Width(value string) StyleOption {
	return func(s *Style) {
		s.Base["width"] = value
	}
}

// height
func Height(value string) StyleOption {
	return func(s *Style) {
		s.Base["height"] = value
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

func Margin(direction MarginDirection, value string) StyleOption {
	return func(s *Style) {
		switch direction {
		case MarginAll:
			s.Base["margin"] = value
		case MarginX:
			s.Base["margin-left"] = value
			s.Base["margin-right"] = value
		case MarginY:
			s.Base["margin-top"] = value
			s.Base["margin-bottom"] = value
		default:
			s.Base["margin-"+string(direction)] = value
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

func Padding(direction PaddingDirection, value string) StyleOption {
	return func(s *Style) {
		switch direction {
		case PaddingAll:
			s.Base["padding"] = value
		case PaddingX:
			s.Base["padding-left"] = value
			s.Base["padding-right"] = value
		case PaddingY:
			s.Base["padding-top"] = value
			s.Base["padding-bottom"] = value
		default:
			s.Base["padding-"+string(direction)] = value
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

func BorderWidth(value string) StyleOption {
	return func(s *Style) {
		s.Base["border-width"] = value
	}
}

func BorderStyle(value BorderStyleType) StyleOption {
	return func(s *Style) {
		s.Base["border-style"] = string(value)
	}
}

func BorderColor(value string) StyleOption {
	return func(s *Style) {
		s.Base["border-color"] = value
	}
}

func Border(shortcut string) StyleOption {
	return func(s *Style) {
		s.Base["border"] = shortcut
	}
}

func BorderRadius(value string) StyleOption {
	return func(s *Style) {
		s.Base["border-radius"] = value
	}
}

// box-sizing
type BoxSizingType string

const (
	BoxSizingContentBox BoxSizingType = "content-box"
	BoxSizingBorderBox  BoxSizingType = "border-box"
)

func BoxSizing(value BoxSizingType) StyleOption {
	return func(s *Style) {
		s.Base["box-sizing"] = string(value)
	}
}

// box shadow
func BoxShadow(offsetX, offsetY, blurRadius, spreadRadius, color string, isInset bool) StyleOption {
	return func(s *Style) {
		// On construit la valeur en respectant l'ordre CSS.
		parts := []string{offsetX, offsetY}
		if blurRadius != "" {
			parts = append(parts, blurRadius)
		}
		if spreadRadius != "" {
			parts = append(parts, spreadRadius)
		}
		parts = append(parts, color)

		shadowValue := strings.Join(parts, " ")

		if isInset {
			shadowValue += " inset"
		}
		s.Base["box-shadow"] = shadowValue
	}
}

// 2. Flexbox
type FlexValue string

const (
	FlexAuto FlexValue = "auto" // Auto flex
	FlexNone FlexValue = "none" // No flex
	Flex1    FlexValue = "1"
	Flex2    FlexValue = "2"
	Flex3    FlexValue = "3"
)

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

func FlexWrap(value FlexWrapValue) StyleOption {
	return func(s *Style) {
		s.Base["flex-wrap"] = string(value)
	}
}

func FlexGrow(value FlexValue) StyleOption {
	return func(s *Style) {
		s.Base["flex-grow"] = string(value)
	}
}

func FlexShrink(value FlexValue) StyleOption {
	return func(s *Style) {
		s.Base["flex-shrink"] = string(value)
	}
}

func FlexBasis(value FlexValue) StyleOption {
	return func(s *Style) {
		s.Base["flex-basis"] = string(value)
	}
}

type AlignItemsValue string

const (
	AlignItemsStart    AlignItemsValue = "start"
	AlignItemsEnd      AlignItemsValue = "end"
	AlignItemsCenter   AlignItemsValue = "center"
	AlignItemsBaseline AlignItemsValue = "baseline"
	AlignItemsStretch  AlignItemsValue = "stretch"
)

func AlignItems(value AlignItemsValue) StyleOption {
	return func(s *Style) {
		s.Base["align-items"] = string(value)
	}
}

type JustifyContentValue string

const (
	JustifyContentStart    JustifyContentValue = "start"
	JustifyContentEnd      JustifyContentValue = "end"
	JustifyContentCenter   JustifyContentValue = "center"
	JustifyContentBaseline JustifyContentValue = "baseline"
	JustifyContentStretch  JustifyContentValue = "stretch"
)

func JustifyContent(value JustifyContentValue) StyleOption {
	return func(s *Style) {
		s.Base["justify-content"] = string(value)
	}
}

func Gap(value string) StyleOption {
	return func(s *Style) {
		s.Base["gap"] = value
	}
}

// --- Typography
func FontFamily(value string) StyleOption {
	return func(s *Style) {
		s.Base["font-family"] = value
	}
}

func FontSize(value string) StyleOption {
	return func(s *Style) {
		s.Base["font-size"] = value
	}
}

func FontWeight(value string) StyleOption {
	return func(s *Style) {
		s.Base["font-weight"] = value
	}
}

func Color(value string) StyleOption {
	return func(s *Style) {
		s.Base["color"] = value
	}
}

type TextAlignValue string

const (
	TextAlignLeft   TextAlignValue = "left"
	TextAlignCenter TextAlignValue = "center"
	TextAlignRight  TextAlignValue = "right"
)

func TextAlign(value TextAlignValue) StyleOption {
	return func(s *Style) {
		s.Base["text-align"] = string(value)
	}
}

func LineHeight(value string) StyleOption {
	return func(s *Style) {
		s.Base["line-height"] = value
	}
}

func LetterSpacing(value string) StyleOption {
	return func(s *Style) {
		s.Base["letter-spacing"] = value
	}
}

type TextDecorationValue string

const (
	TextDecorationNone        TextDecorationValue = "none"
	TextDecorationUnderline   TextDecorationValue = "underline"
	TextDecorationOverline    TextDecorationValue = "overline"
	TextDecorationLineThrough TextDecorationValue = "line-through"
)

func TextDecoration(value TextDecorationValue) StyleOption {
	return func(s *Style) {
		s.Base["text-decoration"] = string(value)
	}
}

// background and appearance
type BackgroundRepeatValue string

const (
	BackgroundRepeatRepeat   BackgroundRepeatValue = "repeat"
	BackgroundRepeatNoRepeat BackgroundRepeatValue = "no-repeat"
	BackgroundRepeatSpace    BackgroundRepeatValue = "space"
)

func BackgroundRepeat(value BackgroundRepeatValue) StyleOption {
	return func(s *Style) {
		s.Base["background-repeat"] = string(value)
	}
}

func BackgroundColor(value string) StyleOption {
	return func(s *Style) {
		s.Base["background-color"] = value
	}
}

func Opacity(value string) StyleOption {
	return func(s *Style) {
		s.Base["opacity"] = value
	}
}

// --- Other

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

func Cursor(value CursorValue) StyleOption {
	return func(s *Style) {
		s.Base["cursor"] = string(value)
	}
}

func ZIndex(value string) StyleOption {
	return func(s *Style) {
		s.Base["z-index"] = value
	}
}

type OverflowValue string

const (
	OverflowVisible OverflowValue = "visible"
	OverflowHidden  OverflowValue = "hidden"
	OverflowScroll  OverflowValue = "scroll"
	OverflowAuto    OverflowValue = "auto"
)

func Overflow(value OverflowValue) StyleOption {
	return func(s *Style) {
		s.Base["overflow"] = string(value)
	}
}

type PositionValue string

const (
	PositionStatic   PositionValue = "static"
	PositionRelative PositionValue = "relative"
	PositionAbsolute PositionValue = "absolute"
)

func Position(value PositionValue) StyleOption {
	return func(s *Style) {
		s.Base["position"] = string(value)
	}
}

// Function to apply a style which is not in the function already defined
func CustomStyle(property string, value string) StyleOption {
	return func(s *Style) {
		s.Base[property] = value
	}
}

// --- Pseudo-classes

// OnHiver applies the given styles for the pseudo :hover
func OnHover(properties ...StyleOption) StyleOption {
	return func(s *Style) {
		// CVreate a temporary style object to collect the hover properties
		hoverStyle := New(properties...)
		s.Pseudos[":hover"] = hoverStyle.Base
	}
}

// OnActive applies the given styles for the pseudo :active
func OnActive(properties ...StyleOption) StyleOption {
	return func(s *Style) {
		activeStyle := New(properties...)
		s.Pseudos[":active"] = activeStyle.Base
	}
}

// OnFocus applies the given styles for the pseudo :focus
func OnFocus(properties ...StyleOption) StyleOption {
	return func(s *Style) {
		focusStyle := New(properties...)
		s.Pseudos[":focus"] = focusStyle.Base
	}
}

// OnFocusWithin applies the given styles for the pseudo :focus-within
func OnFocusWithin(properties ...StyleOption) StyleOption {
	return func(s *Style) {
		focusWithinStyle := New(properties...)
		s.Pseudos[":focus-within"] = focusWithinStyle.Base
	}
}

// OnFocusVisible applies the given styles for the pseudo :focus-visible
func OnFocusVisible(properties ...StyleOption) StyleOption {
	return func(s *Style) {
		focusVisibleStyle := New(properties...)
		s.Pseudos[":focus-visible"] = focusVisibleStyle.Base
	}
}

// --- Media Queries

// Type of media query
type MediaQueryType string

const (
	MediaQueryTypeMinWidth             MediaQueryType = "min-width"
	MediaQueryTypeMaxWidth             MediaQueryType = "max-width"
	MediaQueryTypeMinHeight            MediaQueryType = "min-height"
	MediaQueryTypeMaxHeight            MediaQueryType = "max-height"
	MediaQueryTypeMinAspectRatio       MediaQueryType = "min-aspect-ratio"
	MediaQueryTypeMaxAspectRatio       MediaQueryType = "max-aspect-ratio"
	MediaQueryTypeMinDeviceAspectRatio MediaQueryType = "min-device-aspect-ratio"
	MediaQueryTypeMaxDeviceAspectRatio MediaQueryType = "max-device-aspect-ratio"
	MediaQueryTypeMinDeviceWidth       MediaQueryType = "min-device-width"
	MediaQueryTypeMaxDeviceWidth       MediaQueryType = "max-device-width"
	MediaQueryTypeMinDeviceHeight      MediaQueryType = "min-device-height"
	MediaQueryTypeMaxDeviceHeight      MediaQueryType = "max-device-height"
	MediaQueryTypeMinResolution        MediaQueryType = "min-resolution"
	MediaQueryTypeMaxResolution        MediaQueryType = "max-resolution"
	MediaQueryTypeMinDeviceResolution  MediaQueryType = "min-device-resolution"
	MediaQueryTypeMaxDeviceResolution  MediaQueryType = "max-device-resolution"
)

// MediaQuery applies the given styles for the given media query
func MediaQuery(queryType MediaQueryType, queryValue string, properties ...StyleOption) StyleOption {
	return func(s *Style) {
		mediaStyle := New(properties...)

		// Syntaxe CSS correcte, par ex: "@media (min-width: 600px)"
		fullQuery := fmt.Sprintf("@media (%s: %s)", queryType, queryValue)

		// On initialise la map si elle est nil, une bonne pratique
		if s.MediaQueries[fullQuery] == nil {
			s.MediaQueries[fullQuery] = make(Property)
		}

		// On fusionne les propriétés
		for key, value := range mediaStyle.Base {
			s.MediaQueries[fullQuery][key] = value
		}
	}
}
