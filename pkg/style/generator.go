//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to generate the CSS content of a style
// and the class name of a style.
// It is used to generate the CSS content of a style and the class name of a style.
//
// Basic Usage:
// Developper should not use this file directly has every functionnality are available
// in the style package through the type-safe api
//
// For example, to apply a background color to a component, you can use the following code:
//
//   style := style.New(
//       style.BackgroundColor("#f0f0f0"),
//   )
//
// For Custom style, use the CustomStyle function
//
//   style := style.New(
//       style.CustomStyle("background-color", "#f0f0f0"),
//   )
//
// For more information, see the style package documentation

package style

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
)

// GetClassName generate and return a class name unique and stable for a style
// It uses an hash of the CSS content to ensure unicity
//
// Should not be used directly, see style package documentation
func (s *Style) GetClassName() string {
	if s.className != "" {
		return s.className // Return from the cache
	}

	// We generate first a CSS brut contents, WITHOUT the class name
	rawCSSContent := s.genereateCSSContent()

	// We hash this content to obtain a stable class name
	h := fnv.New32a()
	h.Write([]byte(rawCSSContent))
	s.className = fmt.Sprintf("vtx-%d", h.Sum32())

	return s.className
}

// ToCSS convert a Style object in its textual CSS representation
// It depends on the GetClassName method to generate the class name
// BUT does not call it recursively
//
// Should not be used directly, see style package documentation
func (s *Style) ToCSS() string {
	if s.css != "" {
		return s.css // Retrun from the cache
	}

	className := s.GetClassName() // We get the class name
	var sb strings.Builder

	// Generate the base styles
	sb.WriteString(fmt.Sprintf(".%s {%s}\n", className, propsToCSS(s.Base)))

	// Generate the pseudo-classes
	for pseudo, props := range s.Pseudos {
		sb.WriteString(fmt.Sprintf(".%s%s {%s}\n", className, pseudo, propsToCSS(props)))
	}

	// Generate the media queries
	for query, props := range s.MediaQueries {
		sb.WriteString(fmt.Sprintf("%s { .%s {%s} }\n", query, className, propsToCSS(props)))
	}

	s.css = sb.String()
	return s.css
}

// genereateCSSContent generate the CSS content of the style WITHOUT the class name
// It is used only used to generate the class name
//
// Should not be used directly, see style package documentation
func (s *Style) genereateCSSContent() string {
	var sb strings.Builder

	// Generate the base styles
	sb.WriteString(propsToCSS(s.Base))

	for pseudo, props := range s.Pseudos {
		sb.WriteString(pseudo + propsToCSS(props))
	}
	for query, props := range s.MediaQueries {
		sb.WriteString(query + propsToCSS(props))
	}
	return sb.String()

}

// propsToCSS is a utilitary function to convert a propriety map to text
//
// Should not be used directly, see style package documentation
func propsToCSS(props Property) string {
	if len(props) == 0 {
		return ""
	}
	// Sort the keys for a stable hash
	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s: %s;", k, props[k]))
	}
	return strings.Join(parts, " ")
}
