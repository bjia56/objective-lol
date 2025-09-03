package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/parser"
)

// ModuleResolver handles cross-platform module path resolution and caching
type ModuleResolver struct {
	astCache       map[string]*ast.ProgramNode         // Cache parsed AST by normalized path
	envCache       map[string]*environment.Environment // Cache executed environments by normalized path
	executingStack []string                            // Stack to detect circular imports during execution
	baseDirectory  string                              // Base directory for relative path resolution
}

// NewModuleResolver creates a new module resolver
func NewModuleResolver(baseDirectory string) *ModuleResolver {
	// Normalize base directory to absolute path
	absBase, err := filepath.Abs(baseDirectory)
	if err != nil {
		absBase = baseDirectory
	}

	return &ModuleResolver{
		astCache:       make(map[string]*ast.ProgramNode),
		envCache:       make(map[string]*environment.Environment),
		executingStack: make([]string, 0),
		baseDirectory:  absBase,
	}
}

// ResolvePath converts a POSIX-style module path to a native OS path with .olol extension
func (r *ModuleResolver) ResolvePath(posixPath string) (string, error) {
	// Add .olol extension if not present
	modulePath := posixPath
	if !strings.HasSuffix(strings.ToLower(modulePath), ".olol") {
		modulePath = modulePath + ".olol"
	}

	// Convert POSIX path separators to native OS separators
	nativePath := filepath.FromSlash(modulePath)

	// Handle absolute vs relative paths
	var resolvedPath string
	if filepath.IsAbs(nativePath) {
		resolvedPath = nativePath
	} else {
		// Resolve relative to base directory
		resolvedPath = filepath.Join(r.baseDirectory, nativePath)
	}

	// Convert to absolute path for consistent caching
	absPath, err := filepath.Abs(resolvedPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for %s: %v", posixPath, err)
	}

	return absPath, nil
}

// resolvePathFrom resolves a path relative to a specific directory
func (r *ModuleResolver) resolvePathFrom(posixPath, baseDir string) (string, error) {
	// Add .olol extension if not present
	modulePath := posixPath
	if !strings.HasSuffix(strings.ToLower(modulePath), ".olol") {
		modulePath = modulePath + ".olol"
	}

	// Convert POSIX path separators to native OS separators
	nativePath := filepath.FromSlash(modulePath)

	// Resolve relative to the specified base directory
	resolvedPath := filepath.Join(baseDir, nativePath)

	// Convert to absolute path for consistent caching
	absPath, err := filepath.Abs(resolvedPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for %s from %s: %v", posixPath, baseDir, err)
	}

	return absPath, nil
}

// LoadModule loads and parses a module file, with caching and circular import detection
func (r *ModuleResolver) LoadModule(posixPath string) (*ast.ProgramNode, error) {
	ast, _, err := r.LoadModuleFromWithPath(posixPath, "")
	return ast, err
}

// LoadModuleFrom loads a module with a specific base directory context (for relative imports within modules)
func (r *ModuleResolver) LoadModuleFrom(posixPath, importingFileDir string) (*ast.ProgramNode, error) {
	ast, _, err := r.LoadModuleFromWithPath(posixPath, importingFileDir)
	return ast, err
}

// LoadModuleFromWithPath loads a module and returns both the AST and the resolved file path
func (r *ModuleResolver) LoadModuleFromWithPath(posixPath, importingFileDir string) (*ast.ProgramNode, string, error) {
	// Resolve the path with appropriate context
	var resolvedPath string
	var err error

	if importingFileDir != "" && !filepath.IsAbs(filepath.FromSlash(posixPath)) {
		// Resolve relative to the importing file's directory
		resolvedPath, err = r.resolvePathFrom(posixPath, importingFileDir)
	} else {
		// Use standard resolution (relative to base directory or absolute)
		resolvedPath, err = r.ResolvePath(posixPath)
	}

	if err != nil {
		return nil, "", err
	}

	// Check if AST is already cached
	if cached, exists := r.astCache[resolvedPath]; exists {
		return cached, resolvedPath, nil
	}

	// Check if file exists
	if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("module file not found: %s (resolved to: %s)", posixPath, resolvedPath)
	}

	// Read and parse the module file
	content, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read module file %s: %v", posixPath, err)
	}

	// Parse the module
	lexer := parser.NewLexer(string(content))
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	// Check for parsing errors
	if errors := p.Errors(); len(errors) > 0 {
		return nil, "", fmt.Errorf("parsing errors in module %s: %v", posixPath, errors)
	}

	// Cache the parsed AST
	r.astCache[resolvedPath] = program

	return program, resolvedPath, nil
}

// IsInExecutingStack checks if a path is currently being executed (circular import detection)
func (r *ModuleResolver) IsInExecutingStack(path string) bool {
	return slices.Contains(r.executingStack, path)
}

// AddToExecutingStack adds a path to the execution stack
func (r *ModuleResolver) AddToExecutingStack(path string) {
	r.executingStack = append(r.executingStack, path)
}

// RemoveFromExecutingStack removes the last path from the execution stack
func (r *ModuleResolver) RemoveFromExecutingStack() {
	if len(r.executingStack) > 0 {
		r.executingStack = r.executingStack[:len(r.executingStack)-1]
	}
}

// SetBaseDirectory updates the base directory for relative path resolution
func (r *ModuleResolver) SetBaseDirectory(baseDir string) error {
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute base directory: %v", err)
	}
	r.baseDirectory = absBase
	return nil
}

// GetCachedModules returns a list of all cached module paths
func (r *ModuleResolver) GetCachedModules() []string {
	paths := make([]string, 0, len(r.astCache))
	for path := range r.astCache {
		paths = append(paths, path)
	}
	return paths
}

// ClearCache clears the module cache
func (r *ModuleResolver) ClearCache() {
	r.astCache = make(map[string]*ast.ProgramNode)
	r.envCache = make(map[string]*environment.Environment)
}

// GetCachedEnvironment returns a cached environment for a resolved path
func (r *ModuleResolver) GetCachedEnvironment(resolvedPath string) (*environment.Environment, bool) {
	env, exists := r.envCache[resolvedPath]
	if !exists {
		return nil, false
	}
	return env, true
}

// CacheEnvironment stores an environment for a resolved path
func (r *ModuleResolver) CacheEnvironment(resolvedPath string, env *environment.Environment) {
	r.envCache[resolvedPath] = env
}
