package docs

import (
	"fmt"
	"strings"
)

// FormatJSDocAsMarkdown formats JSDocInfo into rich markdown for LSP hover
func FormatJSDocAsMarkdown(info JSDocInfo, symbolName, symbolType string, symbolKind string) string {
	if !info.HasJSDocTags() {
		// Fallback to simple description if no JSDoc tags
		if info.Description != "" {
			return info.Description
		}
		return ""
	}

	var sections []string

	// Symbol signature section with syntax highlighting
	signature := buildSymbolSignature(info, symbolName, symbolType, symbolKind)
	if signature != "" {
		sections = append(sections, fmt.Sprintf("```olol\n%s\n```", signature))
	}

	// Description section
	if info.Description != "" {
		sections = append(sections, info.Description)
	}

	// Parameters section
	if len(info.Params) > 0 {
		paramSection := formatParameters(info.Params)
		sections = append(sections, paramSection)
	}

	// Returns section
	if info.Returns != "" {
		sections = append(sections, fmt.Sprintf("**Returns:** %s", info.Returns))
	}

	// Examples section
	if len(info.Examples) > 0 {
		exampleSections := formatExamples(info.Examples)
		sections = append(sections, exampleSections...)
	}

	// Notes section
	if len(info.Notes) > 0 {
		noteSection := formatNotes(info.Notes)
		sections = append(sections, noteSection)
	}

	// See also section
	if len(info.See) > 0 {
		seeSection := formatSeeAlso(info.See)
		sections = append(sections, seeSection)
	}

	// Category information
	if info.Category != "" {
		sections = append(sections, fmt.Sprintf("**Category:** %s", formatCategory(info.Category)))
	}

	return strings.Join(sections, "\n\n")
}

// buildSymbolSignature builds a syntax-highlighted signature from JSDoc info
func buildSymbolSignature(info JSDocInfo, symbolName, symbolType string, symbolKind string) string {
	if info.Syntax != "" {
		return info.Syntax
	}

	// Fallback signature generation
	switch strings.ToLower(symbolKind) {
	case "function":
		signature := fmt.Sprintf("FUNCSHUN %s", symbolName)

		if len(info.Params) > 0 {
			var paramParts []string
			for _, param := range info.Params {
				if param.Type != "" && param.Type != "..." {
					paramParts = append(paramParts, fmt.Sprintf("WIT %s TEH %s", param.Name, param.Type))
				} else {
					paramParts = append(paramParts, fmt.Sprintf("WIT %s", param.Name))
				}
			}
			signature += " " + strings.Join(paramParts, " AN ")
		}

		if symbolType != "" && symbolType != "NOTHIN" {
			signature += fmt.Sprintf(" GIVEZ %s", symbolType)
		}

		return signature

	case "class":
		return fmt.Sprintf("CLAS %s", symbolName)

	case "variable":
		if symbolType != "" {
			return fmt.Sprintf("VARIABLE %s TEH %s", symbolName, symbolType)
		}
		return fmt.Sprintf("VARIABLE %s", symbolName)

	default:
		if symbolType != "" {
			return fmt.Sprintf("%s: %s", symbolName, symbolType)
		}
		return symbolName
	}
}

// formatParameters formats parameter information as a markdown table
func formatParameters(params []JSDocParam) string {
	if len(params) == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, "**Parameters:**")
	lines = append(lines, "")

	for _, param := range params {
		paramLine := fmt.Sprintf("- `%s`", param.Name)

		if param.Type != "" {
			paramLine += fmt.Sprintf(" (%s)", param.Type)
		}

		if param.Description != "" {
			paramLine += fmt.Sprintf(": %s", param.Description)
		}

		lines = append(lines, paramLine)
	}

	return strings.Join(lines, "\n")
}

// formatExamples formats code examples with syntax highlighting
func formatExamples(examples []JSDocExample) []string {
	var sections []string

	for i, example := range examples {
		var exampleLines []string

		title := example.Title
		if title == "" {
			title = fmt.Sprintf("Example %d", i+1)
		}

		exampleLines = append(exampleLines, fmt.Sprintf("**%s:**", title))
		exampleLines = append(exampleLines, "")
		exampleLines = append(exampleLines, "```olol")
		exampleLines = append(exampleLines, example.Code...)
		exampleLines = append(exampleLines, "```")

		if example.Description != "" {
			exampleLines = append(exampleLines, "")
			exampleLines = append(exampleLines, example.Description)
		}

		sections = append(sections, strings.Join(exampleLines, "\n"))
	}

	return sections
}

// formatNotes formats note information
func formatNotes(notes []string) string {
	if len(notes) == 0 {
		return ""
	}

	var lines []string

	if len(notes) == 1 {
		lines = append(lines, fmt.Sprintf("**Note:** %s", notes[0]))
	} else {
		lines = append(lines, "**Notes:**")
		lines = append(lines, "")
		for _, note := range notes {
			lines = append(lines, fmt.Sprintf("- %s", note))
		}
	}

	return strings.Join(lines, "\n")
}

// formatSeeAlso formats see-also references
func formatSeeAlso(seeRefs []string) string {
	if len(seeRefs) == 0 {
		return ""
	}

	// Format as inline links or just comma-separated list
	var formatted []string
	for _, ref := range seeRefs {
		formatted = append(formatted, fmt.Sprintf("`%s`", ref))
	}

	return fmt.Sprintf("**See also:** %s", strings.Join(formatted, ", "))
}

// formatCategory formats category name for display
func formatCategory(category string) string {
	// Convert kebab-case to Title Case
	words := strings.Split(category, "-")
	var titleWords []string
	for _, word := range words {
		if len(word) > 0 {
			titleWords = append(titleWords, strings.ToUpper(string(word[0]))+strings.ToLower(word[1:]))
		}
	}
	return strings.Join(titleWords, " ")
}

// FormatBasicSymbolInfo formats basic symbol information for non-JSDoc symbols
func FormatBasicSymbolInfo(name, symbolType, kind, visibility, parentClass, sourceModule, rawDocs string) string {
	var sections []string

	// Basic signature
	signature := fmt.Sprintf("%s %s", strings.ToUpper(kind), name)
	if symbolType != "" {
		signature += fmt.Sprintf(": %s", symbolType)
	}
	sections = append(sections, fmt.Sprintf("```olol\n%s\n```", signature))

	// Context information
	var contextInfo []string
	if visibility != "" {
		contextInfo = append(contextInfo, fmt.Sprintf("**Visibility:** `%s`", visibility))
	}
	if parentClass != "" {
		contextInfo = append(contextInfo, fmt.Sprintf("**Class:** `%s`", parentClass))
	}
	if sourceModule != "" && sourceModule != "stdlib" {
		contextInfo = append(contextInfo, fmt.Sprintf("**Module:** `%s`", sourceModule))
	}
	if len(contextInfo) > 0 {
		sections = append(sections, strings.Join(contextInfo, " • "))
	}

	// Raw documentation
	if rawDocs != "" {
		sections = append(sections, "---")
		sections = append(sections, rawDocs)
	}

	return strings.Join(sections, "\n\n")
}
