//go:build js && wasm

package component

import (
	"fmt"
	"strings"
)

// Text component is a span stateless component
type TextComponent struct {
	*FunctionalComponent
}

func Text(content string) *TextComponent {
	return &TextComponent{
		FunctionalComponent: NewFunctionalComponent(func() *VNode {
			return &VNode{Type: VNodeText, Text: content}
		}),
	}
}

// Heading component
type HeadingComponent struct {
	*FunctionalComponent
	level int
}

func H(content string, level int) *HeadingComponent {
	if level < 1 || level > 6 {
		level = 1
	}

	return &HeadingComponent{
		FunctionalComponent: NewFunctionalComponent(func() *VNode {
			return &VNode{
				Type: VNodeElement,
				Tag:  fmt.Sprintf("h%d", level),
				Children: []*VNode{
					{
						Type: VNodeText,
						Text: content,
					},
				},
			}
		}),
		level: level,
	}
}

// Paragraph component
type ParagraphComponent struct {
	*FunctionalComponent
}

func P(content ...Component) *ParagraphComponent {
	childrenVNodes := make([]*VNode, len(content))
	for i, child := range content {
		childrenVNodes[i] = child.Render()
	}
	return &ParagraphComponent{
		FunctionalComponent: NewFunctionalComponent(func() *VNode {
			return &VNode{
				Type:     VNodeElement,
				Tag:      "p",
				Children: childrenVNodes,
			}
		}),
	}
}

// Markdown engine
type MarkdownComponent struct {
	*FunctionalComponent
}

// Usage :
// markdown := component.Markdown("**bold** *italic* # Heading 1")
func Markdown(content string) *MarkdownComponent {
	return &MarkdownComponent{
		FunctionalComponent: NewFunctionalComponent(func() *VNode {
			return &VNode{
				Type:     VNodeElement,
				Tag:      "div",
				Children: parseMarkdown(content),
			}
		}),
	}
}

// Actual markdown parser (basic implementation)
func parseMarkdown(content string) []*VNode {
	content = strings.TrimSpace(content)
	if content == "" {
		return []*VNode{}
	}

	// Split into paragraphs (double newlines)
	paragraphs := strings.Split(content, "\n\n")
	nodes := make([]*VNode, 0, len(paragraphs))

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}

		// Parse different markdown elements
		if node := parseMarkdownBlock(paragraph); node != nil {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func parseMarkdownBlock(text string) *VNode {
	text = strings.TrimSpace(text)

	// Headers: # H1, ## H2, ### H3, etc.
	if strings.HasPrefix(text, "#") {
		return parseHeader(text)
	}

	// Regular paragraph with inline formatting
	return &VNode{
		Type:     VNodeElement,
		Tag:      "p",
		Children: parseInlineMarkdown(text),
	}
}

func parseHeader(text string) *VNode {
	level := 0
	for i, char := range text {
		if char == '#' {
			level++
		} else {
			text = strings.TrimSpace(text[i:])
			break
		}
	}

	if level > 6 {
		level = 6
	}
	if level < 1 {
		level = 1
	}

	return &VNode{
		Type:     VNodeElement,
		Tag:      fmt.Sprintf("h%d", level),
		Children: parseInlineMarkdown(text),
	}
}

func parseInlineMarkdown(text string) []*VNode {
	if text == "" {
		return []*VNode{}
	}

	// Simple regex-based parsing for **bold** and *italic*
	result := []*VNode{}
	remaining := text

	for remaining != "" {
		// Find next markdown token
		boldStart := strings.Index(remaining, "**")
		italicStart := strings.Index(remaining, "*")

		// No more markdown - add remaining as text
		if boldStart == -1 && italicStart == -1 {
			if remaining != "" {
				result = append(result, &VNode{
					Type: VNodeText,
					Text: remaining,
				})
			}
			break
		}

		// Handle **bold**
		if boldStart != -1 && (italicStart == -1 || boldStart < italicStart) {
			// Add text before bold
			if boldStart > 0 {
				result = append(result, &VNode{
					Type: VNodeText,
					Text: remaining[:boldStart],
				})
			}

			// Find closing **
			boldEnd := strings.Index(remaining[boldStart+2:], "**")
			if boldEnd != -1 {
				boldText := remaining[boldStart+2 : boldStart+2+boldEnd]
				result = append(result, &VNode{
					Type: VNodeElement,
					Tag:  "strong",
					Children: []*VNode{
						{
							Type: VNodeText,
							Text: boldText,
						},
					},
				})
				remaining = remaining[boldStart+2+boldEnd+2:]
			} else {
				// No closing **, treat as regular text
				result = append(result, &VNode{
					Type: VNodeText,
					Text: remaining[:boldStart+2],
				})
				remaining = remaining[boldStart+2:]
			}
			continue
		}

		// Handle *italic* (similar logic)
		if italicStart != -1 {
			// Add text before italic
			if italicStart > 0 {
				result = append(result, &VNode{
					Type: VNodeText,
					Text: remaining[:italicStart],
				})
			}

			// Find closing *
			italicEnd := strings.Index(remaining[italicStart+1:], "*")
			if italicEnd != -1 {
				italicText := remaining[italicStart+1 : italicStart+1+italicEnd]
				result = append(result, &VNode{
					Type: VNodeElement,
					Tag:  "em",
					Children: []*VNode{
						{
							Type: VNodeText,
							Text: italicText,
						},
					},
				})
				remaining = remaining[italicStart+1+italicEnd+1:]
			} else {
				// No closing *, treat as regular text
				result = append(result, &VNode{
					Type: VNodeText,
					Text: remaining[:italicStart+1],
				})
				remaining = remaining[italicStart+1:]
			}
		}
	}

	return result
}
