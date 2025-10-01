//go:build js && wasm

// precompilation.go is a file that contains the precompilation framework
// it is used to precompile the styles and cache the CSS
// It can significantly improve the performance of the application
// by avoiding the generation of the CSS at runtime

package style

import (
	"fmt"
	"log"
	"time"
)

// PrecompilationEngine manages style precompilation for performance optimization
type PrecompilationEngine struct {
	// Style storage and lookup
	styles         []*Style
	styleMap       map[string]*Style // className -> Style for fast lookup
	compiledStates map[*Style]bool

	// Performance tracking
	precompileTime time.Duration
	totalStyles    int

	// Configuration
	enableMetrics bool
}

// PrecompilationConfig configures the precompilation engine
type PrecompilationConfig struct {
	EnableMetrics   bool // Track performance metrics
	InitialCapacity int  // Initial capacity for style storage
}

// NewPrecompilationEngine creates a new precompilation engine
func NewPrecompilationEngine(config *PrecompilationConfig) *PrecompilationEngine {
	if config == nil {
		config = &PrecompilationConfig{
			EnableMetrics:   true,
			InitialCapacity: 100,
		}
	}

	engine := &PrecompilationEngine{
		styles:         make([]*Style, 0, config.InitialCapacity),
		styleMap:       make(map[string]*Style, config.InitialCapacity),
		compiledStates: make(map[*Style]bool, config.InitialCapacity),
		enableMetrics:  config.EnableMetrics,
	}

	return engine
}

// AddStyle adds a style to the precompilation queue
func (e *PrecompilationEngine) AddStyle(style *Style) {
	if style == nil {
		return
	}

	e.styles = append(e.styles, style)
}

// AddStyles adds multiple styles to the precompilation queue
func (e *PrecompilationEngine) AddStyles(styles ...*Style) {
	for _, style := range styles {
		if style != nil {
			e.AddStyle(style)
		}
	}
}

// RunAllPrecompilation precompiles all registered styles
func (e *PrecompilationEngine) RunAllPrecompilation() {
	if e.enableMetrics {
		start := time.Now()
		defer func() {
			e.precompileTime = time.Since(start)
			e.totalStyles = len(e.styles)
			log.Printf("Precompiled %d styles in %v", e.totalStyles, e.precompileTime)
		}()
	}

	// Precompile all registered styles
	for _, style := range e.styles {
		e.precompileStyle(style)
	}
}

// precompileStyle precompiles a single style
func (e *PrecompilationEngine) precompileStyle(style *Style) {
	if style == nil || e.IsStyleCompiled(style) {
		return
	}

	// Generate CSS (triggers caching)
	css := style.ToCSS()
	className := style.GetClassName()

	// Store in lookup map
	e.styleMap[className] = style
	e.compiledStates[style] = true

	if e.enableMetrics && css != "" {
		// Could add CSS size tracking here
	}
}

// IsStyleCompiled checks if a style has been precompiled
func (e *PrecompilationEngine) IsStyleCompiled(style *Style) bool {
	if style == nil {
		return false
	}
	return e.compiledStates[style]
}

// GetStyleByClassName retrieves a precompiled style by its CSS class name
func (e *PrecompilationEngine) GetStyleByClassName(className string) *Style {
	return e.styleMap[className]
}

// FindStyle searches for a style in the precompiled collection
func (e *PrecompilationEngine) FindStyle(style *Style) *Style {
	if style == nil {
		return nil
	}

	// Check if already compiled
	if e.IsStyleCompiled(style) {
		return style
	}

	// Try to find by className if generated
	className := style.GetClassName()
	if found := e.GetStyleByClassName(className); found != nil {
		return found
	}

	return nil
}

// PrecompileAndRegister precompiles a style and registers it immediately
func (e *PrecompilationEngine) PrecompileAndRegister(style *Style) *Style {
	if style == nil {
		return nil
	}

	e.AddStyle(style)
	e.precompileStyle(style)
	return style
}

// GetMetrics returns performance metrics
type PrecompilationMetrics struct {
	TotalStyles         int
	PrecompileTime      time.Duration
	AverageTimePerStyle time.Duration
	CacheHitRate        float64
}

func (e *PrecompilationEngine) GetMetrics() PrecompilationMetrics {
	var avgTime time.Duration
	if e.totalStyles > 0 {
		avgTime = e.precompileTime / time.Duration(e.totalStyles)
	}

	return PrecompilationMetrics{
		TotalStyles:         e.totalStyles,
		PrecompileTime:      e.precompileTime,
		AverageTimePerStyle: avgTime,
		CacheHitRate:        e.calculateCacheHitRate(),
	}
}

// calculateCacheHitRate calculates cache hit rate (placeholder for now)
func (e *PrecompilationEngine) calculateCacheHitRate() float64 {
	// Could track cache hits vs misses
	return 0.0
}

// Clear removes all precompiled styles
func (e *PrecompilationEngine) Clear() {
	e.styles = e.styles[:0]
	e.styleMap = make(map[string]*Style)
	e.compiledStates = make(map[*Style]bool)
	e.precompileTime = 0
	e.totalStyles = 0
}

// PrintStats prints precompilation statistics
func (e *PrecompilationEngine) PrintStats() {
	metrics := e.GetMetrics()
	fmt.Printf("=== Vortex Precompilation Stats ===\n")
	fmt.Printf("Total Styles: %d\n", metrics.TotalStyles)
	fmt.Printf("Precompile Time: %v\n", metrics.PrecompileTime)
	fmt.Printf("Avg Time/Style: %v\n", metrics.AverageTimePerStyle)
	fmt.Printf("=====================================\n")
}

// Global precompilation engine instance
var globalPrecompiler *PrecompilationEngine

// init initializes the global precompilation engine
func init() {
	globalPrecompiler = NewPrecompilationEngine(&PrecompilationConfig{
		EnableMetrics:   true,
		InitialCapacity: 200,
	})
	initCommonStyles()
}

// Convenience functions for global precompiler

// Precompile adds styles to the global precompiler
func Precompile(styles ...*Style) {
	globalPrecompiler.AddStyles(styles...)
}

// Precompile adds a style to the global precompiler
func (s *Style) Precompile() *Style {
	globalPrecompiler.AddStyle(s)
	return s
}

// RunPrecompilation runs precompilation on the global engine
func RunPrecompilation() {
	globalPrecompiler.RunAllPrecompilation()
}

// IsPrecompiled checks if a style is precompiled in the global engine
func IsPrecompiled(style *Style) bool {
	return globalPrecompiler.IsStyleCompiled(style)
}

// PrecompilationStats prints global precompilation statistics
func PrecompilationStats() {
	globalPrecompiler.PrintStats()
}
