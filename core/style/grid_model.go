//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply CSS Grid to a style.
//
// Basic Usage:
//
//	s := style.New().
//	  Display(style.DisplayGrid).
//	  GridTemplateColumns(style.Fr(1), style.Fr(2), style.PxTrack(200)).
//	  GridTemplateRows(style.MinMax(style.Px(48), style.Auto)).
//	  ColumnGap(style.Px(16)).
//	  RowGap(style.Px(8)).
//	  GridAutoFlow(style.AutoFlowRowDense)
//
// For more information, see the style package documentation.
package style

import (
	"fmt"
	"log"
	"strings"
)

// GridTrack represents a single grid track sizing value as a CSS string.
// Examples: "1fr", "200px", "minmax(48px, auto)", "repeat(3, 1fr)"
type GridTrack string

func (g GridTrack) String() string { return string(g) }
func (g GridTrack) Validate() error {
	// Keep permissive for forward-compat; we already validate LengthValue elsewhere.
	// Could add stricter validation later if needed.
	return nil
}

// Constructors for common track types

// Fr returns a fractional track (e.g. 1fr, 2fr).
func Fr(value float64) GridTrack { return GridTrack(fmt.Sprintf("%.2ffr", value)) }

// PxTrack returns a pixel-based track (e.g. 200px).
func PxTrack(px float64) GridTrack { return GridTrack(Px(px).String()) }

// EmTrack returns an em-based track.
func EmTrack(em float64) GridTrack { return GridTrack(Em(em).String()) }

// RemTrack returns a rem-based track.
func RemTrack(rem float64) GridTrack { return GridTrack(Rem(rem).String()) }

// PercentTrack returns a percent-based track (e.g. 50%).
func PercentTrack(p float64) GridTrack { return GridTrack(Percent(p).String()) }

// MinMax returns a minmax(...) track from two CSS values.
func MinMax(min CSSValue, max CSSValue) GridTrack {
	return GridTrack(fmt.Sprintf("minmax(%s, %s)", min.String(), max.String()))
}

// Repeat returns a repeat(...) track list, e.g. repeat(3, 1fr) or repeat(12, minmax(0, 1fr))
func Repeat(count int, track GridTrack) GridTrack {
	return GridTrack(fmt.Sprintf("repeat(%d, %s)", count, track.String()))
}

type AreaName string

func (a AreaName) String() string {
	return string(a)
}

func (a AreaName) Validate() error {
	return nil
}

type GridTemplateBuilder struct {
	nrows       int
	ncols       int
	areas       [][]AreaName
	columnsSize []GridTrack
	rowsSize    []GridTrack
}

func NewGridTemplateBuilder(nrows int, ncols int) *GridTemplateBuilder {
	b := &GridTemplateBuilder{
		nrows:       nrows,
		ncols:       ncols,
		areas:       make([][]AreaName, nrows),
		columnsSize: make([]GridTrack, ncols),
		rowsSize:    make([]GridTrack, nrows),
	}
	for r := 0; r < nrows; r++ {
		b.areas[r] = make([]AreaName, ncols)
	}
	return b
}

func (b *GridTemplateBuilder) Area(row int, col int, area string) *GridTemplateBuilder {
	if row < 0 || row >= b.nrows || col < 0 || col >= b.ncols {
		fmt.Printf("GridTemplateBuilder: Area: row %d or col %d out of bounds\n", row, col)
		return b
	}
	b.areas[row][col] = AreaName(area)
	return b
}

func (b *GridTemplateBuilder) ColumnSize(col int, size GridTrack) *GridTemplateBuilder {
	if col < 0 || col >= b.ncols {
		fmt.Printf("GridTemplateBuilder: ColumnSize: col %d out of bounds\n", col)
		return b
	}
	b.columnsSize[col] = size
	return b
}

func (b *GridTemplateBuilder) RowSize(row int, size GridTrack) *GridTemplateBuilder {
	if row < 0 || row >= b.nrows {
		fmt.Printf("GridTemplateBuilder: RowSize: row %d out of bounds\n", row)
		return b
	}
	b.rowsSize[row] = size
	return b
}

func (s *Style) GridTemplate(b *GridTemplateBuilder) *Style {
	// grid-template-areas
	areas := ""
	for r := 0; r < b.nrows; r++ {
		line := ""
		for c := 0; c < b.ncols; c++ {
			name := b.areas[r][c].String()
			if name == "" {
				name = "." // empty cell
			}
			if c > 0 {
				line += " "
			}
			line += name
		}
		areas += fmt.Sprintf("\"%s\"", line)
		if r < b.nrows-1 {
			areas += "\n"
		}
	}

	// grid-template-rows
	rowsSize := ""
	for r := 0; r < b.nrows; r++ {
		track := b.rowsSize[r]
		if track == "" {
			rowsSize += "auto"
		} else {
			rowsSize += track.String()
		}
		if r < b.nrows-1 {
			rowsSize += " "
		}
	}

	// grid-template-columns
	columnsSize := ""
	for c := 0; c < b.ncols; c++ {
		track := b.columnsSize[c]
		if track == "" {
			columnsSize += "auto"
		} else {
			columnsSize += track.String()
		}
		if c < b.ncols-1 {
			columnsSize += " "
		}
	}

	s.Base["grid-template-areas"] = areas
	s.Base["grid-template-rows"] = rowsSize
	s.Base["grid-template-columns"] = columnsSize
	return s
}

func (s *Style) GridTemplateColumns(tracks ...GridTrack) *Style {
	parts := make([]string, len(tracks))
	for i, t := range tracks {
		parts[i] = t.String()
	}
	s.Base["grid-template-columns"] = strings.Join(parts, " ")
	return s
}

func (s *Style) GridTemplateRows(tracks ...GridTrack) *Style {
	parts := make([]string, len(tracks))
	for i, t := range tracks {
		parts[i] = t.String()
	}
	s.Base["grid-template-rows"] = strings.Join(parts, " ")
	return s
}

func (s *Style) GridTemplateAreas(rows ...string) *Style {
	quoted := make([]string, len(rows))
	for i, r := range rows {
		quoted[i] = fmt.Sprintf("\"%s\"", r)
	}
	s.Base["grid-template-areas"] = strings.Join(quoted, "\n")
	return s
}

// Auto placement

type AutoFlowValue string

const (
	AutoFlowRow         AutoFlowValue = "row"
	AutoFlowColumn      AutoFlowValue = "column"
	AutoFlowDense       AutoFlowValue = "dense"
	AutoFlowRowDense    AutoFlowValue = "row dense"
	AutoFlowColumnDense AutoFlowValue = "column dense"
)

func (a AutoFlowValue) String() string { return string(a) }
func (a AutoFlowValue) Validate() error {
	return ValidateCSS("grid-auto-flow", string(a))
}

func (s *Style) GridAutoFlow(value AutoFlowValue) *Style {
	s.Base["grid-auto-flow"] = value.String()
	return s
}

// GridAutoRows sets "grid-auto-rows" (implicit rows sizing)
func (s *Style) GridAutoRows(track GridTrack) *Style {
	s.Base["grid-auto-rows"] = track.String()
	return s
}

// GridAutoColumns sets "grid-auto-columns" (implicit columns sizing)
func (s *Style) GridAutoColumns(track GridTrack) *Style {
	s.Base["grid-auto-columns"] = track.String()
	return s
}

// Gaps

// ColumnGap sets "column-gap"; prefer this when row/column gaps differ.
func (s *Style) ColumnGap(value LengthValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["column-gap"] = value.String()
	return s
}

// RowGap sets "row-gap"; prefer this when row/column gaps differ.
func (s *Style) RowGap(value LengthValue) *Style {
	if err := value.Validate(); err != nil {
		log.Printf("CSS validation warning: %v", err)
	}
	s.Base["row-gap"] = value.String()
	return s
}

// Alignment for grid containers (items/content)
// Note: AlignItems/JustifyContent already exist; they map well to grid.
// Here we add the missing per-grid-axis item alignment and the "place-*" shorthands.

type JustifyItemsValue string
type AlignItemsGridValue string

const (
	JustifyItemsStart   JustifyItemsValue = "start"
	JustifyItemsEnd     JustifyItemsValue = "end"
	JustifyItemsCenter  JustifyItemsValue = "center"
	JustifyItemsStretch JustifyItemsValue = "stretch"

	AlignItemsStartG   AlignItemsGridValue = "start"
	AlignItemsEndG     AlignItemsGridValue = "end"
	AlignItemsCenterG  AlignItemsGridValue = "center"
	AlignItemsStretchG AlignItemsGridValue = "stretch"
)

func (v JustifyItemsValue) String() string    { return string(v) }
func (v AlignItemsGridValue) String() string  { return string(v) }
func (v JustifyItemsValue) Validate() error   { return ValidateCSS("justify-items", string(v)) }
func (v AlignItemsGridValue) Validate() error { return ValidateCSS("align-items", string(v)) }

func (s *Style) JustifyItems(value JustifyItemsValue) *Style {
	s.Base["justify-items"] = value.String()
	return s
}

func (s *Style) AlignItemsGrid(value AlignItemsGridValue) *Style {
	s.Base["align-items"] = value.String()
	return s
}

type PlaceItemsValue struct {
	Align   AlignItemsGridValue
	Justify JustifyItemsValue
}

func (p PlaceItemsValue) String() string {
	// If align==justify we can emit single value; otherwise two values.
	if p.Align.String() == p.Justify.String() {
		return p.Align.String()
	}
	return fmt.Sprintf("%s %s", p.Align.String(), p.Justify.String())
}

func (s *Style) PlaceItems(value PlaceItemsValue) *Style {
	s.Base["place-items"] = value.String()
	return s
}

// Similar shorthands for content alignment

type ContentAxisValue string

const (
	ContentStart        ContentAxisValue = "start"
	ContentEnd          ContentAxisValue = "end"
	ContentCenter       ContentAxisValue = "center"
	ContentStretch      ContentAxisValue = "stretch"
	ContentSpaceBetween ContentAxisValue = "space-between"
	ContentSpaceAround  ContentAxisValue = "space-around"
	ContentSpaceEvenly  ContentAxisValue = "space-evenly"
)

func (v ContentAxisValue) String() string  { return string(v) }
func (v ContentAxisValue) Validate() error { return ValidateCSS("content", string(v)) }

func (s *Style) JustifyContentGrid(value ContentAxisValue) *Style {
	s.Base["justify-content"] = value.String()
	return s
}

func (s *Style) AlignContentGrid(value ContentAxisValue) *Style {
	s.Base["align-content"] = value.String()
	return s
}

type PlaceContentValue struct {
	Align   ContentAxisValue
	Justify ContentAxisValue
}

func (p PlaceContentValue) String() string {
	if p.Align.String() == p.Justify.String() {
		return p.Align.String()
	}
	return fmt.Sprintf("%s %s", p.Align.String(), p.Justify.String())
}

func (s *Style) PlaceContent(value PlaceContentValue) *Style {
	s.Base["place-content"] = value.String()
	return s
}

// Child placement helpers

// GridArea assigns a named area to the item (or shorthand row/column positions when provided as single string).
func (s *Style) GridArea(name string) *Style {
	s.Base["grid-area"] = name
	return s
}

// GridColumn sets "grid-column: start / end".
func (s *Style) GridColumn(start, end string) *Style {
	s.Base["grid-column"] = fmt.Sprintf("%s / %s", start, end)
	return s
}

// GridRow sets "grid-row: start / end".
func (s *Style) GridRow(start, end string) *Style {
	s.Base["grid-row"] = fmt.Sprintf("%s / %s", start, end)
	return s
}

// GridColumnSpan sets "grid-column: span N".
func (s *Style) GridColumnSpan(n int) *Style {
	s.Base["grid-column"] = fmt.Sprintf("span %d", n)
	return s
}

// GridRowSpan sets "grid-row: span N".
func (s *Style) GridRowSpan(n int) *Style {
	s.Base["grid-row"] = fmt.Sprintf("span %d", n)
	return s
}
