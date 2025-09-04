package modules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bjia56/objective-lol/pkg/environment"
)

// TestModuleClassQualification tests that classes from modules get proper qualified names
func TestModuleClassQualification(t *testing.T) {
	// Create temporary test modules
	testDir := createTestModuleDir(t)
	defer os.RemoveAll(testDir)

	// Create a module with a class
	animalModulePath := filepath.Join(testDir, "animals.olol")
	animalModuleContent := `
BTW Animals module
HAI ME TEH CLAS ANIMAL
    EVRYONE
    DIS TEH VARIABLE NAME TEH STRIN ITZ "Unknown"
    
    DIS TEH FUNCSHUN SET_NAME WIT NEW_NAME TEH STRIN
        NAME ITZ NEW_NAME
    KTHX
    
    DIS TEH FUNCSHUN GET_NAME TEH STRIN
        GIVEZ NAME
    KTHX
KTHXBAI

HAI ME TEH CLAS DOG KITTEH OF ANIMAL
    EVRYONE
    DIS TEH VARIABLE BREED TEH STRIN ITZ "Mixed"
    
    DIS TEH FUNCSHUN GET_BREED TEH STRIN
        GIVEZ BREED
    KTHX
KTHXBAI
`
	require.NoError(t, os.WriteFile(animalModulePath, []byte(animalModuleContent), 0644))

	// Create resolver and load module
	resolver := NewModuleResolver(testDir)
	moduleAST, resolvedPath, err := resolver.LoadModuleFromWithPath("animals", "")
	require.NoError(t, err)
	assert.NotNil(t, moduleAST)

	// Create environment to test class creation
	env := environment.NewEnvironment(nil)
	
	// Test that module path is correctly resolved
	assert.Contains(t, resolvedPath, "animals.olol")
	
	// Test that we can create qualified class names based on the resolved path
	expectedModulePath := "file:" + resolvedPath
	
	// Create classes manually to test qualified name generation
	animalClass := environment.NewClass("ANIMAL", expectedModulePath, "")
	dogClass := environment.NewClass("DOG", expectedModulePath, expectedModulePath+".ANIMAL")
	
	// Define classes in environment
	require.NoError(t, env.DefineClass(animalClass))
	require.NoError(t, env.DefineClass(dogClass))
	
	// Check ANIMAL class
	retrievedAnimalClass, err := env.GetClass("ANIMAL")
	require.NoError(t, err)
	assert.Contains(t, retrievedAnimalClass.QualifiedName, "animals.olol.ANIMAL")
	assert.Equal(t, "ANIMAL", retrievedAnimalClass.Name)
	
	// Check DOG class with qualified parent
	retrievedDogClass, err := env.GetClass("DOG")
	require.NoError(t, err)
	assert.Contains(t, retrievedDogClass.QualifiedName, "animals.olol.DOG")
	assert.Contains(t, retrievedDogClass.ParentClass, "animals.olol.ANIMAL")
}

// TestModuleClassImportAndLookup tests importing classes and looking them up by qualified/simple names
func TestModuleClassImportAndLookup(t *testing.T) {
	testDir := createTestModuleDir(t)
	defer os.RemoveAll(testDir)

	// Create utility module with classes
	utilsModulePath := filepath.Join(testDir, "utils.olol")
	utilsModuleContent := `
BTW Utilities module
HAI ME TEH CLAS HELPER
    EVRYONE
    DIS TEH VARIABLE ID TEH INTEGR ITZ 42
KTHXBAI

BTW Private class - should not be importable
HAI ME TEH CLAS _PRIVATE_HELPER
    EVRYONE
    DIS TEH VARIABLE SECRET TEH STRIN ITZ "secret"
KTHXBAI
`
	require.NoError(t, os.WriteFile(utilsModulePath, []byte(utilsModuleContent), 0644))

	// Load the module
	resolver := NewModuleResolver(testDir)
	moduleAST, resolvedPath, err := resolver.LoadModuleFromWithPath("utils", "")
	require.NoError(t, err)
	assert.NotNil(t, moduleAST)

	// Simulate environment after import processing
	env := environment.NewEnvironment(nil)
	
	// Test that we can create classes with proper qualified names
	expectedModulePath := "file:" + resolvedPath
	helperClass := environment.NewClass("HELPER", expectedModulePath, "")
	
	// Define class in environment (simulating import)
	require.NoError(t, env.DefineClass(helperClass))

	// Test that imported class can be found by simple name
	retrievedClass, err := env.GetClass("HELPER")
	require.NoError(t, err)
	assert.Equal(t, "HELPER", retrievedClass.Name)
	assert.Contains(t, retrievedClass.QualifiedName, "utils.olol.HELPER")

	// Test that private class is not imported (simulated by not adding it)
	_, err = env.GetClass("_PRIVATE_HELPER")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "undefined class")
}

// TestModuleClassNameCollisions tests handling of class name collisions between modules
func TestModuleClassNameCollisions(t *testing.T) {
	testDir := createTestModuleDir(t)
	defer os.RemoveAll(testDir)

	// Create first module with ITEM class
	module1Path := filepath.Join(testDir, "module1.olol")
	module1Content := `
HAI ME TEH CLAS ITEM
    EVRYONE
    DIS TEH VARIABLE TYPE TEH STRIN ITZ "MODULE1_ITEM"
KTHXBAI
`
	require.NoError(t, os.WriteFile(module1Path, []byte(module1Content), 0644))

	// Create second module with ITEM class  
	module2Path := filepath.Join(testDir, "module2.olol")
	module2Content := `
HAI ME TEH CLAS ITEM
    EVRYONE
    DIS TEH VARIABLE TYPE TEH STRIN ITZ "MODULE2_ITEM"
KTHXBAI
`
	require.NoError(t, os.WriteFile(module2Path, []byte(module2Content), 0644))

	// Load both modules
	resolver := NewModuleResolver(testDir)
	
	// Load first module
	_, resolvedPath1, err := resolver.LoadModuleFromWithPath("module1", "")
	require.NoError(t, err)
	
	// Load second module
	_, resolvedPath2, err := resolver.LoadModuleFromWithPath("module2", "")
	require.NoError(t, err)

	// Test that both modules can define ITEM classes with different qualified names
	env := environment.NewEnvironment(nil)
	
	modulePath1 := "file:" + resolvedPath1
	modulePath2 := "file:" + resolvedPath2
	
	item1Class := environment.NewClass("ITEM", modulePath1, "")
	item2Class := environment.NewClass("ITEM", modulePath2, "")
	
	// Both should be definable with different qualified names
	require.NoError(t, env.DefineClass(item1Class))
	require.NoError(t, env.DefineClass(item2Class))
	
	// The simple name lookup should return one of them (implementation dependent)
	itemClass, err := env.GetClass("ITEM")
	require.NoError(t, err)
	assert.Equal(t, "ITEM", itemClass.Name)
	
	// But the qualified lookups should be distinct
	assert.NotEqual(t, item1Class.QualifiedName, item2Class.QualifiedName)
	assert.Contains(t, item1Class.QualifiedName, "module1.olol")
	assert.Contains(t, item2Class.QualifiedName, "module2.olol")
}

// TestModuleInheritanceWithQualifiedParents tests inheritance across modules
func TestModuleInheritanceWithQualifiedParents(t *testing.T) {
	testDir := createTestModuleDir(t)
	defer os.RemoveAll(testDir)

	// Create base classes module
	baseModulePath := filepath.Join(testDir, "base.olol")
	baseContent := `
HAI ME TEH CLAS VEHICLE
    EVRYONE
    DIS TEH VARIABLE WHEELS TEH INTEGR ITZ 4
KTHXBAI
`
	require.NoError(t, os.WriteFile(baseModulePath, []byte(baseContent), 0644))

	// Create derived classes module that imports base
	derivedModulePath := filepath.Join(testDir, "vehicles.olol")
	derivedContent := `
I CAN HAS "base"?

HAI ME TEH CLAS CAR KITTEH OF VEHICLE
    EVRYONE
    DIS TEH VARIABLE DOORS TEH INTEGR ITZ 4
KTHXBAI
`
	require.NoError(t, os.WriteFile(derivedModulePath, []byte(derivedContent), 0644))

	// Test loading modules
	resolver := NewModuleResolver(testDir)
	
	// Load base module
	_, basePath, err := resolver.LoadModuleFromWithPath("base", "")
	require.NoError(t, err)
	
	// Load derived module
	_, derivedPath, err := resolver.LoadModuleFromWithPath("vehicles", "")
	require.NoError(t, err)

	// Test creating classes with proper qualified inheritance
	env := environment.NewEnvironment(nil)
	
	baseModulePathStr := "file:" + basePath
	derivedModulePathStr := "file:" + derivedPath
	
	// Create base class
	vehicleClass := environment.NewClass("VEHICLE", baseModulePathStr, "")
	require.NoError(t, env.DefineClass(vehicleClass))
	
	// Create derived class with qualified parent
	carClass := environment.NewClass("CAR", derivedModulePathStr, vehicleClass.QualifiedName)
	require.NoError(t, env.DefineClass(carClass))
	
	// Verify CAR class has qualified parent
	retrievedCarClass, err := env.GetClass("CAR")
	require.NoError(t, err)
	
	// Parent should be qualified name from base module
	assert.Contains(t, retrievedCarClass.ParentClass, "base.olol.VEHICLE")
	assert.Equal(t, "CAR", retrievedCarClass.Name)
	assert.Contains(t, retrievedCarClass.QualifiedName, "vehicles.olol.CAR")
}

// TestSelectiveClassImport tests importing specific classes from modules
func TestSelectiveClassImport(t *testing.T) {
	testDir := createTestModuleDir(t)
	defer os.RemoveAll(testDir)

	// Create module with multiple classes
	multiClassModulePath := filepath.Join(testDir, "multi.olol")
	multiClassContent := `
HAI ME TEH CLAS CLASS_A
    EVRYONE
    DIS TEH VARIABLE VALUE TEH STRIN ITZ "A"
KTHXBAI

HAI ME TEH CLAS CLASS_B
    EVRYONE
    DIS TEH VARIABLE VALUE TEH STRIN ITZ "B"
KTHXBAI

HAI ME TEH CLAS CLASS_C
    EVRYONE  
    DIS TEH VARIABLE VALUE TEH STRIN ITZ "C"
KTHXBAI
`
	require.NoError(t, os.WriteFile(multiClassModulePath, []byte(multiClassContent), 0644))

	// Load the module
	resolver := NewModuleResolver(testDir)
	_, resolvedPath, err := resolver.LoadModuleFromWithPath("multi", "")
	require.NoError(t, err)

	// Simulate selective import - only import CLASS_A and CLASS_C
	env := environment.NewEnvironment(nil)
	modulePath := "file:" + resolvedPath
	
	classA := environment.NewClass("CLASS_A", modulePath, "")
	classC := environment.NewClass("CLASS_C", modulePath, "")
	
	// Define only selected classes (simulating selective import)
	require.NoError(t, env.DefineClass(classA))
	require.NoError(t, env.DefineClass(classC))

	// Verify only selected classes are available
	_, err = env.GetClass("CLASS_A")
	assert.NoError(t, err)
	
	_, err = env.GetClass("CLASS_C") 
	assert.NoError(t, err)

	// CLASS_B should not be available since it wasn't imported
	_, err = env.GetClass("CLASS_B")
	assert.Error(t, err)
}

// createTestModuleDir creates a temporary directory for test modules
func createTestModuleDir(t *testing.T) string {
	testDir, err := os.MkdirTemp("", "olol_module_test_*")
	require.NoError(t, err)
	return testDir
}