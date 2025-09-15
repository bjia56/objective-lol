package docs

import (
	"strings"
)

// JSDocInfo represents parsed JSDoc-style documentation
type JSDocInfo struct {
	Description string
	Syntax      string
	Params      []JSDocParam
	Returns     string
	Examples    []JSDocExample
	Category    string
	See         []string
	Notes       []string
	Type        string // for variables
	Value       string // for variables
}

// JSDocParam represents a function parameter
type JSDocParam struct {
	Name        string
	Type        string
	Description string
}

// JSDocExample represents a code example
type JSDocExample struct {
	Title       string
	Code        []string
	Description string
}

// ParseJSDoc parses JSDoc-style documentation from a list of strings
func ParseJSDoc(docs []string) JSDocInfo {
	info := JSDocInfo{
		Params:   []JSDocParam{},
		Examples: []JSDocExample{},
		See:      []string{},
		Notes:    []string{},
	}

	var currentExample *JSDocExample
	var descriptionLines []string

	for _, line := range docs {
		line = strings.TrimSpace(line)

		if line == "" {
			if currentExample != nil {
				currentExample.Code = append(currentExample.Code, "")
			} else {
				descriptionLines = append(descriptionLines, "")
			}
			continue
		}

		if strings.HasPrefix(line, "@") {
			// Save current example if any
			if currentExample != nil {
				info.Examples = append(info.Examples, *currentExample)
				currentExample = nil
			}

			parts := strings.SplitN(line, " ", 2)
			tag := parts[0]
			content := ""
			if len(parts) > 1 {
				content = parts[1]
			}

			switch tag {
			case "@syntax":
				info.Syntax = content
			case "@param":
				param := parseParam(content)
				info.Params = append(info.Params, param)
			case "@returns":
				info.Returns = content
			case "@example":
				if currentExample != nil {
					info.Examples = append(info.Examples, *currentExample)
				}
				currentExample = &JSDocExample{
					Title: content,
					Code:  []string{},
				}
			case "@category":
				info.Category = content
			case "@see":
				info.See = strings.Split(content, ", ")
			case "@note":
				info.Notes = append(info.Notes, content)
			case "@type":
				info.Type = content
			case "@value":
				info.Value = content
			}
		} else {
			if currentExample != nil {
				currentExample.Code = append(currentExample.Code, line)
			} else {
				descriptionLines = append(descriptionLines, line)
			}
		}
	}

	// Save final example
	if currentExample != nil {
		info.Examples = append(info.Examples, *currentExample)
	}

	// Clean up description
	info.Description = strings.TrimSpace(strings.Join(descriptionLines, "\n"))

	return info
}

// parseParam parses a parameter specification in the format: {TYPE} name - description
func parseParam(content string) JSDocParam {
	param := JSDocParam{}

	if strings.HasPrefix(content, "{") {
		typeEnd := strings.Index(content, "}")
		if typeEnd > 0 {
			param.Type = strings.TrimSpace(content[1:typeEnd])
			content = strings.TrimSpace(content[typeEnd+1:])
		}
	}

	parts := strings.SplitN(content, " - ", 2)
	if len(parts) >= 1 {
		param.Name = strings.TrimSpace(parts[0])
	}
	if len(parts) >= 2 {
		param.Description = strings.TrimSpace(parts[1])
	}

	return param
}

// IsJSDoc checks if documentation contains JSDoc tags
func IsJSDoc(docs []string) bool {
	for _, line := range docs {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "@") &&
		   (strings.HasPrefix(line, "@param") ||
		    strings.HasPrefix(line, "@returns") ||
		    strings.HasPrefix(line, "@syntax") ||
		    strings.HasPrefix(line, "@example") ||
		    strings.HasPrefix(line, "@note") ||
		    strings.HasPrefix(line, "@see") ||
		    strings.HasPrefix(line, "@category") ||
		    strings.HasPrefix(line, "@type") ||
		    strings.HasPrefix(line, "@value")) {
			return true
		}
	}
	return false
}

// HasJSDocTags checks if parsed JSDoc contains meaningful structured information
func (info JSDocInfo) HasJSDocTags() bool {
	return info.Syntax != "" || len(info.Params) > 0 || info.Returns != "" ||
		   len(info.Examples) > 0 || info.Category != "" || len(info.See) > 0 ||
		   len(info.Notes) > 0 || info.Type != "" || info.Value != ""
}