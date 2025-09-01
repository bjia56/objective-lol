package modules

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModuleResolver(t *testing.T) {
	baseDir := "/test/base"
	resolver := NewModuleResolver(baseDir)

	assert.NotNil(t, resolver)
	// The actual fields are private, so we can't test them directly
	// Just verify the resolver was created successfully
}

func TestModuleResolver_ResolvePath(t *testing.T) {
	baseDir := "/home/user/project"
	resolver := NewModuleResolver(baseDir)

	tests := []struct {
		name           string
		modulePath     string
		expectedSuffix string // Just check the suffix since base paths vary
	}{
		{
			"Simple module",
			"helper",
			"helper.olol",
		},
		{
			"Module with subdirectory",
			"utils/math",
			filepath.Join("utils", "math.olol"),
		},
		{
			"Module with extension already",
			"helper.olol",
			"helper.olol",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resolved, err := resolver.ResolvePath(test.modulePath)
			require.NoError(t, err)

			// Check that the path ends with the expected suffix
			assert.True(t, filepath.IsAbs(resolved), "Should return absolute path")
			assert.True(t,
				strings.HasSuffix(resolved, test.expectedSuffix) ||
					strings.HasSuffix(resolved, filepath.FromSlash(test.expectedSuffix)),
				"Path %s should end with %s", resolved, test.expectedSuffix)
		})
	}
}

func TestModuleResolver_LoadModule(t *testing.T) {
	resolver := NewModuleResolver("/nonexistent")

	// Test loading a non-existent module
	_, err := resolver.LoadModule("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "module file not found")
}

func TestModuleResolver_PathHandling(t *testing.T) {
	resolver := NewModuleResolver("/base")

	// Test that paths with .olol extension are handled correctly
	resolved1, err := resolver.ResolvePath("test")
	require.NoError(t, err)

	resolved2, err := resolver.ResolvePath("test.olol")
	require.NoError(t, err)

	// Both should resolve to the same path (with .olol extension)
	assert.Equal(t, resolved1, resolved2)
	assert.True(t, strings.HasSuffix(resolved1, ".olol"))
}

func TestModuleResolver_AbsolutePaths(t *testing.T) {
	resolver := NewModuleResolver("/base")

	tests := []string{
		"relative/path",
		"/absolute/path",
		"simple",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			resolved, err := resolver.ResolvePath(test)
			require.NoError(t, err)
			assert.True(t, filepath.IsAbs(resolved), "Should always return absolute paths")
			assert.True(t, strings.HasSuffix(resolved, ".olol"), "Should always have .olol extension")
		})
	}
}
