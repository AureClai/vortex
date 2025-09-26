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

package style

var (
	flexCenterStyle *Style
	fullWidthStyle  *Style
	hiddenStyle     *Style
)

func init() {
	// Pre-generate the common styles
	flexCenterStyle = New(
		Display(DisplayFlex),
		JustifyContent(JustifyContentCenter),
		AlignItems(AlignItemsCenter),
	)
	fullWidthStyle = New(
		Width(Percent(100)),
	)
	hiddenStyle = New(
		Display(DisplayNone),
	)

	// CSS to cache
	_ = flexCenterStyle.ToCSS()
	_ = fullWidthStyle.ToCSS()
	_ = hiddenStyle.ToCSS()
}

func FlexCenter() *Style { return flexCenterStyle }
func FullWidth() *Style  { return fullWidthStyle }
func Hidden() *Style     { return hiddenStyle }
