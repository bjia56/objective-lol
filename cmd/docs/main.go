package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bjia56/objective-lol/pkg/stdlib"
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

// ModuleDoc represents documentation for a complete module
type ModuleDoc struct {
	Name       string
	Variables  map[string]stdlib.StdlibDefinition
	Functions  map[string]stdlib.StdlibDefinition
	Classes    map[string]stdlib.StdlibDefinition
	Categories []string
	JSDocInfo  map[string]JSDocInfo
	ClassItems map[string]stdlib.StdlibDefinition // CLASS.ITEM -> Definition (includes kind)
}

func main() {
	var outputDir string
	flag.StringVar(&outputDir, "output", "docs/standard-library", "Output directory for generated documentation")
	flag.Parse()

	if err := generateDocs(outputDir); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Documentation generated successfully in %s\n", outputDir)
}

func generateDocs(outputDir string) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Get stdlib initializers
	initializers := stdlib.DefaultStdlibInitializers()
	globalInitializers := stdlib.DefaultGlobalInitializers()

	modules := make(map[string]*ModuleDoc)

	// Process stdlib modules
	for moduleName, initializer := range initializers {
		definitions := stdlib.GetStdlibDefinitions(initializer)
		categories := stdlib.GetModuleCategories(moduleName)

		moduleDoc := &ModuleDoc{
			Name:       moduleName,
			Variables:  make(map[string]stdlib.StdlibDefinition),
			Functions:  make(map[string]stdlib.StdlibDefinition),
			Classes:    make(map[string]stdlib.StdlibDefinition),
			Categories: categories,
			JSDocInfo:  make(map[string]JSDocInfo),
			ClassItems: make(map[string]stdlib.StdlibDefinition),
		}

		for _, def := range definitions {
			jsDoc := parseJSDoc(def.Docs)
			moduleDoc.JSDocInfo[def.Name] = jsDoc

			switch def.Kind {
			case stdlib.StdlibDefinitionKindVariable:
				// Check if this is a class member variable
				if strings.Contains(def.Name, ".") {
					// This is a class member variable, store it in ClassItems
					moduleDoc.ClassItems[def.Name] = def
				} else {
					moduleDoc.Variables[def.Name] = def
				}
			case stdlib.StdlibDefinitionKindFunction:
				// Check if this is a class method
				if strings.Contains(def.Name, ".") {
					// This is a class method, store it in ClassItems
					moduleDoc.ClassItems[def.Name] = def
				} else {
					moduleDoc.Functions[def.Name] = def
				}
			case stdlib.StdlibDefinitionKindClass:
				moduleDoc.Classes[def.Name] = def
			}
		}

		modules[moduleName] = moduleDoc
	}

	// Process global modules (ARRAYS, MAPS)
	for _, initializer := range globalInitializers {
		definitions := stdlib.GetStdlibDefinitions(initializer)

		moduleName := "GLOBAL"
		if len(definitions) > 0 {
			switch definitions[0].Name {
			case "BUKKIT":
				moduleName = "ARRAYS"
			case "BASKIT":
				moduleName = "MAPS"
			}
		}

		categories := stdlib.GetModuleCategories(moduleName)

		moduleDoc := &ModuleDoc{
			Name:       moduleName,
			Variables:  make(map[string]stdlib.StdlibDefinition),
			Functions:  make(map[string]stdlib.StdlibDefinition),
			Classes:    make(map[string]stdlib.StdlibDefinition),
			Categories: categories,
			JSDocInfo:  make(map[string]JSDocInfo),
			ClassItems: make(map[string]stdlib.StdlibDefinition),
		}

		for _, def := range definitions {
			jsDoc := parseJSDoc(def.Docs)
			moduleDoc.JSDocInfo[def.Name] = jsDoc

			switch def.Kind {
			case stdlib.StdlibDefinitionKindVariable:
				// Check if this is a class member variable
				if strings.Contains(def.Name, ".") {
					// This is a class member variable, store it in ClassItems
					moduleDoc.ClassItems[def.Name] = def
				} else {
					moduleDoc.Variables[def.Name] = def
				}
			case stdlib.StdlibDefinitionKindFunction:
				// Check if this is a class method
				if strings.Contains(def.Name, ".") {
					// This is a class method, store it in ClassItems
					moduleDoc.ClassItems[def.Name] = def
				} else {
					moduleDoc.Functions[def.Name] = def
				}
			case stdlib.StdlibDefinitionKindClass:
				moduleDoc.Classes[def.Name] = def
			}
		}

		modules[moduleName] = moduleDoc
	}

	// Generate markdown files for each module
	for moduleName, moduleDoc := range modules {
		filename := strings.ToLower(moduleName) + ".md"
		filepath := filepath.Join(outputDir, filename)

		if err := generateModuleDoc(filepath, moduleDoc); err != nil {
			return fmt.Errorf("failed to generate %s: %v", filename, err)
		}
	}

	// Generate overview file
	overviewPath := filepath.Join(outputDir, "overview.md")
	if err := generateOverview(overviewPath, modules); err != nil {
		return fmt.Errorf("failed to generate overview: %v", err)
	}

	return nil
}

func parseJSDoc(docs []string) JSDocInfo {
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

func parseParam(content string) JSDocParam {
	// Parse format: {TYPE} name - description
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

func generateModuleDoc(filepath string, moduleDoc *ModuleDoc) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "# %s Module\n\n", moduleDoc.Name)

	// Write import section
	fmt.Fprintf(file, "## Import\n\n")
	fmt.Fprintf(file, "```lol\n")
	fmt.Fprintf(file, "BTW Full import\n")
	fmt.Fprintf(file, "I CAN HAS %s?\n\n", moduleDoc.Name)
	fmt.Fprintf(file, "BTW Selective import examples\n")

	// Show a few example imports
	count := 0
	for name := range moduleDoc.Functions {
		if count < 2 {
			fmt.Fprintf(file, "I CAN HAS %s FROM %s?\n", name, moduleDoc.Name)
			count++
		}
	}
	for name := range moduleDoc.Variables {
		if count < 3 {
			fmt.Fprintf(file, "I CAN HAS %s FROM %s?\n", name, moduleDoc.Name)
			count++
		}
	}
	fmt.Fprintf(file, "```\n\n")

	// Organize content by categories
	if len(moduleDoc.Categories) > 0 {
		for _, category := range moduleDoc.Categories {
			writeCategory(file, category, moduleDoc)
		}

		// Write uncategorized items
		writeUncategorized(file, moduleDoc)
	} else {
		// No categories, write everything
		if len(moduleDoc.Variables) > 0 {
			writeVariables(file, moduleDoc.Variables, moduleDoc.JSDocInfo)
		}
		if len(moduleDoc.Functions) > 0 {
			writeFunctions(file, moduleDoc.Functions, moduleDoc.JSDocInfo)
		}
		if len(moduleDoc.Classes) > 0 {
			writeClasses(file, moduleDoc.Classes, moduleDoc.JSDocInfo, moduleDoc.ClassItems)
		}
	}

	return nil
}

func writeCategory(file *os.File, category string, moduleDoc *ModuleDoc) {
	// Find items in this category
	variables := make(map[string]stdlib.StdlibDefinition)
	functions := make(map[string]stdlib.StdlibDefinition)
	classes := make(map[string]stdlib.StdlibDefinition)

	for name, def := range moduleDoc.Variables {
		if moduleDoc.JSDocInfo[name].Category == category {
			variables[name] = def
		}
	}
	for name, def := range moduleDoc.Functions {
		if moduleDoc.JSDocInfo[name].Category == category {
			functions[name] = def
		}
	}
	for name, def := range moduleDoc.Classes {
		if moduleDoc.JSDocInfo[name].Category == category {
			classes[name] = def
		}
	}

	if len(variables) == 0 && len(functions) == 0 && len(classes) == 0 {
		return
	}

	// Write category header
	categoryTitle := strings.ReplaceAll(category, "-", " ")
	// Simple title case without deprecated strings.Title
	words := strings.Fields(categoryTitle)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	categoryTitle = strings.Join(words, " ")
	fmt.Fprintf(file, "## %s\n\n", categoryTitle)

	writeVariables(file, variables, moduleDoc.JSDocInfo)
	writeFunctions(file, functions, moduleDoc.JSDocInfo)
	writeClasses(file, classes, moduleDoc.JSDocInfo, moduleDoc.ClassItems)
}

func writeUncategorized(file *os.File, moduleDoc *ModuleDoc) {
	// Find uncategorized items
	variables := make(map[string]stdlib.StdlibDefinition)
	functions := make(map[string]stdlib.StdlibDefinition)
	classes := make(map[string]stdlib.StdlibDefinition)

	for name, def := range moduleDoc.Variables {
		jsDoc := moduleDoc.JSDocInfo[name]
		if jsDoc.Category == "" || !contains(moduleDoc.Categories, jsDoc.Category) {
			variables[name] = def
		}
	}
	for name, def := range moduleDoc.Functions {
		jsDoc := moduleDoc.JSDocInfo[name]
		if jsDoc.Category == "" || !contains(moduleDoc.Categories, jsDoc.Category) {
			functions[name] = def
		}
	}
	for name, def := range moduleDoc.Classes {
		jsDoc := moduleDoc.JSDocInfo[name]
		if jsDoc.Category == "" || !contains(moduleDoc.Categories, jsDoc.Category) {
			classes[name] = def
		}
	}

	if len(variables) > 0 || len(functions) > 0 || len(classes) > 0 {
		fmt.Fprintf(file, "## Miscellaneous\n\n")
		writeVariables(file, variables, moduleDoc.JSDocInfo)
		writeFunctions(file, functions, moduleDoc.JSDocInfo)
		writeClasses(file, classes, moduleDoc.JSDocInfo, moduleDoc.ClassItems)
	}
}

func writeVariables(file *os.File, variables map[string]stdlib.StdlibDefinition, jsDocInfo map[string]JSDocInfo) {
	if len(variables) == 0 {
		return
	}

	// Sort variables by name
	names := make([]string, 0, len(variables))
	for name := range variables {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		def := variables[name]
		jsDoc := jsDocInfo[name]

		fmt.Fprintf(file, "### %s\n\n", name)

		if jsDoc.Description != "" {
			fmt.Fprintf(file, "%s\n\n", jsDoc.Description)
		}

		fmt.Fprintf(file, "**Type:** %s\n", def.Type)
		if jsDoc.Value != "" {
			fmt.Fprintf(file, "**Value:** %s\n", jsDoc.Value)
		}
		fmt.Fprintf(file, "\n")

		// Write examples
		for _, example := range jsDoc.Examples {
			if example.Title != "" {
				fmt.Fprintf(file, "```lol\n")
				fmt.Fprintf(file, "I CAN HAS %s FROM %s?\n\n", name, "MODULE") // We'll need module context
				for _, codeLine := range example.Code {
					fmt.Fprintf(file, "%s\n", codeLine)
				}
				fmt.Fprintf(file, "```\n\n")
			}
		}
	}
}

func writeFunctions(file *os.File, functions map[string]stdlib.StdlibDefinition, jsDocInfo map[string]JSDocInfo) {
	if len(functions) == 0 {
		return
	}

	// Sort functions by name
	names := make([]string, 0, len(functions))
	for name := range functions {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		def := functions[name]
		jsDoc := jsDocInfo[name]

		fmt.Fprintf(file, "### %s\n\n", name)

		if jsDoc.Description != "" {
			fmt.Fprintf(file, "%s\n\n", jsDoc.Description)
		}

		if jsDoc.Syntax != "" {
			fmt.Fprintf(file, "**Syntax:** `%s`\n", jsDoc.Syntax)
		}
		fmt.Fprintf(file, "**Returns:** %s\n\n", def.Type)

		// Write parameters
		if len(jsDoc.Params) > 0 {
			fmt.Fprintf(file, "**Parameters:**\n")
			for _, param := range jsDoc.Params {
				fmt.Fprintf(file, "- `%s` (%s): %s\n", param.Name, param.Type, param.Description)
			}
			fmt.Fprintf(file, "\n")
		}

		// Write examples
		for _, example := range jsDoc.Examples {
			if example.Title != "" && len(example.Code) > 0 {
				fmt.Fprintf(file, "**Example: %s**\n\n", example.Title)
				fmt.Fprintf(file, "```lol\n")
				for _, codeLine := range example.Code {
					fmt.Fprintf(file, "%s\n", codeLine)
				}
				fmt.Fprintf(file, "```\n\n")
			}
		}

		// Write notes
		for _, note := range jsDoc.Notes {
			fmt.Fprintf(file, "**Note:** %s\n\n", note)
		}

		// Write see also
		if len(jsDoc.See) > 0 {
			fmt.Fprintf(file, "**See also:** %s\n\n", strings.Join(jsDoc.See, ", "))
		}
	}
}

func writeClasses(file *os.File, classes map[string]stdlib.StdlibDefinition, jsDocInfo map[string]JSDocInfo, classItems map[string]stdlib.StdlibDefinition) {
	if len(classes) == 0 {
		return
	}

	// Sort classes by name
	names := make([]string, 0, len(classes))
	for name := range classes {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		jsDoc := jsDocInfo[name]

		fmt.Fprintf(file, "### %s Class\n\n", name)

		if jsDoc.Description != "" {
			fmt.Fprintf(file, "%s\n\n", jsDoc.Description)
		}

		// Write class methods
		classMethods := getClassMethods(name, classItems)
		if len(classMethods) > 0 {
			fmt.Fprintf(file, "**Methods:**\n\n")
			for _, methodName := range classMethods {
				fullMethodName := name + "." + methodName
				methodDoc := jsDocInfo[fullMethodName]

				fmt.Fprintf(file, "#### %s\n\n", methodName)

				if methodDoc.Description != "" {
					fmt.Fprintf(file, "%s\n\n", methodDoc.Description)
				}

				if methodDoc.Syntax != "" {
					fmt.Fprintf(file, "**Syntax:** `%s`\n", methodDoc.Syntax)
				}

				// Write parameters
				if len(methodDoc.Params) > 0 {
					fmt.Fprintf(file, "**Parameters:**\n")
					for _, param := range methodDoc.Params {
						fmt.Fprintf(file, "- `%s` (%s): %s\n", param.Name, param.Type, param.Description)
					}
					fmt.Fprintf(file, "\n")
				}

				// Write examples
				for _, example := range methodDoc.Examples {
					if example.Title != "" && len(example.Code) > 0 {
						fmt.Fprintf(file, "**Example: %s**\n\n", example.Title)
						fmt.Fprintf(file, "```lol\n")
						for _, codeLine := range example.Code {
							fmt.Fprintf(file, "%s\n", codeLine)
						}
						fmt.Fprintf(file, "```\n\n")
					}
				}

				// Write notes
				for _, note := range methodDoc.Notes {
					fmt.Fprintf(file, "**Note:** %s\n\n", note)
				}
			}
		}

		// Write class member variables
		classMemberVars := getClassMemberVariables(name, classItems)
		if len(classMemberVars) > 0 {
			fmt.Fprintf(file, "**Member Variables:**\n\n")
			for _, varName := range classMemberVars {
				fullVarName := name + "." + varName
				varDoc := jsDocInfo[fullVarName]

				fmt.Fprintf(file, "#### %s\n\n", varName)

				if varDoc.Description != "" {
					fmt.Fprintf(file, "%s\n\n", varDoc.Description)
				}

				if varDoc.Type != "" {
					fmt.Fprintf(file, "**Type:** %s\n", varDoc.Type)
				}

				if varDoc.Value != "" {
					fmt.Fprintf(file, "**Value:** %s\n", varDoc.Value)
				}

				fmt.Fprintf(file, "\n")

				// Write examples
				for _, example := range varDoc.Examples {
					if example.Title != "" && len(example.Code) > 0 {
						fmt.Fprintf(file, "**Example: %s**\n\n", example.Title)
						fmt.Fprintf(file, "```lol\n")
						for _, codeLine := range example.Code {
							fmt.Fprintf(file, "%s\n", codeLine)
						}
						fmt.Fprintf(file, "```\n\n")
					}
				}

				// Write notes
				for _, note := range varDoc.Notes {
					fmt.Fprintf(file, "**Note:** %s\n\n", note)
				}
			}
		}

		// Write class-level examples
		for _, example := range jsDoc.Examples {
			if example.Title != "" && len(example.Code) > 0 {
				fmt.Fprintf(file, "**Example: %s**\n\n", example.Title)
				fmt.Fprintf(file, "```lol\n")
				for _, codeLine := range example.Code {
					fmt.Fprintf(file, "%s\n", codeLine)
				}
				fmt.Fprintf(file, "```\n\n")
			}
		}
	}
}

// getClassMethods extracts method names for a given class from ClassItems
func getClassMethods(className string, classItems map[string]stdlib.StdlibDefinition) []string {
	var methods []string
	prefix := className + "."

	for name, def := range classItems {
		if strings.HasPrefix(name, prefix) && def.Kind == stdlib.StdlibDefinitionKindFunction {
			methodName := strings.TrimPrefix(name, prefix)
			methods = append(methods, methodName)
		}
	}

	sort.Strings(methods)
	return methods
}

// getClassMemberVariables extracts member variable names for a given class from ClassItems
func getClassMemberVariables(className string, classItems map[string]stdlib.StdlibDefinition) []string {
	var variables []string
	prefix := className + "."

	for name, def := range classItems {
		if strings.HasPrefix(name, prefix) && def.Kind == stdlib.StdlibDefinitionKindVariable {
			variableName := strings.TrimPrefix(name, prefix)
			variables = append(variables, variableName)
		}
	}

	sort.Strings(variables)
	return variables
}

func generateOverview(filepath string, modules map[string]*ModuleDoc) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "# Standard Library Overview\n\n")
	fmt.Fprintf(file, "This documentation is automatically generated from JSDoc-style comments in the source code.\n\n")
	fmt.Fprintf(file, "## Available Modules\n\n")

	// Sort modules by name
	names := make([]string, 0, len(modules))
	for name := range modules {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		moduleDoc := modules[name]
		filename := strings.ToLower(name) + ".md"

		totalItems := len(moduleDoc.Variables) + len(moduleDoc.Functions) + len(moduleDoc.Classes)
		fmt.Fprintf(file, "- [%s](%s) - %d items\n", name, filename, totalItems)
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
