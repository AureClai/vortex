//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply the media queries to a style.
// It is used to apply the media queries to a style.
//
// Basic Usage:
//
//   style := style.New(
//       style.MediaQuery(style.MediaQueryTypeMinWidth, "600px", style.BackgroundColor("#f0f0f0")),
//   )
//
// For more information, see the style package documentation
//
// All properties available are (table with css equivalent of the media query type	)
// | Media Query Type | CSS Equivalent |
// | -------- | -------------- |
// | Media Query Type Min Width | @media (min-width: 600px)        |
// | Media Query Type Max Width | @media (max-width: 600px)        |
// | Media Query Type Min Height | @media (min-height: 600px)        |
// | Media Query Type Max Height | @media (max-height: 600px)        |
// | Media Query Type Min Aspect Ratio | @media (min-aspect-ratio: 600px)        |
// | Media Query Type Max Aspect Ratio | @media (max-aspect-ratio: 600px)        |
// | Media Query Type Min Device Aspect Ratio | @media (min-device-aspect-ratio: 600px)        |
// | Media Query Type Max Device Aspect Ratio | @media (max-device-aspect-ratio: 600px)        |
// | Media Query Type Min Device Width | @media (min-device-width: 600px)        |
// | Media Query Type Max Device Width | @media (max-device-width: 600px)        |
// | Media Query Type Min Device Height | @media (min-device-height: 600px)        |
// | Media Query Type Max Device Height | @media (max-device-height: 600px)        |
// | Media Query Type Min Resolution | @media (min-resolution: 600px)        |
// | Media Query Type Max Resolution | @media (max-resolution: 600px)        |
// | Media Query Type Min Device Resolution | @media (min-device-resolution: 600px)        |
// | Media Query Type Max Device Resolution | @media (max-device-resolution: 600px)        |
//
// For more information, see the style package documentation

package style

import "fmt"

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
