//go:build js && wasm

// common.go is a file that contains all the already generated code for the styling
// it is used to avoid code duplication and to make the code more readable
//
// It offers developper to use common styling without having to write the code themselves
// The covergage is highly incremental and will be updated as the styling is updated
//
// Covered :
// - Flex Center : FlexCenter : display: flex; justify-content: center; align-items: center;
// - Full Width : FullWidth : width: 100%;
// - Hidden : Hidden : display: none;
// - Color Pink Neon Glow : ColorPinkNeonGlow : color: #FF69B4;

package style

var (
	flexCenterStyle        *Style
	fullWidthStyle         *Style
	hiddenStyle            *Style
	colorPinkNeonGlowStyle *Style
)

func initCommonStyles() {
	// Pre-generate the common styles
	flexCenterStyle = New().
		Display(DisplayFlex).
		JustifyContent(JustifyContentCenter).
		AlignItems(AlignItemsCenter).Precompile()
	fullWidthStyle = New().
		Width(Percent(100)).Precompile()
	hiddenStyle = New().
		Display(DisplayNone)
	colorPinkNeonGlowStyle = New().
		Color(HEX("#FFFFFF")).
		TextShadow(
			TextShadowValue{OffsetX: Px(0), OffsetY: Px(0), BlurRadius: Px(10), Color: HEX("#FF69B4")},
			TextShadowValue{OffsetX: Px(0), OffsetY: Px(0), BlurRadius: Px(10), Color: HEX("#FF69B4")},
		).Precompile()
}

// Quick access to common patterns (automatically precompiled)
func FullWidth() *Style         { return fullWidthStyle }
func Hidden() *Style            { return hiddenStyle }
func ColorPinkNeonGlow() *Style { return colorPinkNeonGlowStyle }
func FlexCenter() *Style        { return flexCenterStyle }
