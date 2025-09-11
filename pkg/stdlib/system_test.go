package stdlib

import (
	"os"
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/stretchr/testify/assert"
)

func TestSystemModuleRegistration(t *testing.T) {
	env := environment.NewEnvironment(nil)
	
	// Register SYSTEM module
	err := RegisterSYSTEMInEnv(env)
	assert.NoError(t, err)
	
	// Verify ENVBASKIT class is registered
	class, err := env.GetClass("ENVBASKIT")
	assert.NoError(t, err)
	assert.Equal(t, "ENVBASKIT", class.Name)
	
	// Verify ENV global variable exists
	envVariable, err := env.GetVariable("ENV")
	assert.NoError(t, err)
	assert.NotNil(t, envVariable)
	assert.Equal(t, "ENV", envVariable.Name)
	assert.Equal(t, "ENVBASKIT", envVariable.Type)
	
	// Verify ENV is an ENVBASKIT instance
	envInstance, ok := envVariable.Value.(*environment.ObjectInstance)
	assert.True(t, ok)
	assert.Equal(t, "ENVBASKIT", envInstance.Class.Name)
}

func TestENVBasicOperations(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSYSTEMInEnv(env)
	
	envVariable, _ := env.GetVariable("ENV")
	envInstance := envVariable.Value.(*environment.ObjectInstance)
	
	// Test PUT - should set both in BASKIT and actual environment
	putFunc := envInstance.Class.PublicFunctions["PUT"]
	key := environment.StringValue("TEST_KEY")
	value := environment.StringValue("TEST_VALUE")
	
	result, err := putFunc.NativeImpl(nil, envInstance, []environment.Value{key, value})
	assert.NoError(t, err)
	assert.Equal(t, environment.NOTHIN, result)
	
	// Verify environment variable was actually set
	actualValue := os.Getenv("TEST_KEY")
	assert.Equal(t, "TEST_VALUE", actualValue)
	
	// Test GET - should retrieve from BASKIT/environment
	getFunc := envInstance.Class.PublicFunctions["GET"]
	result, err = getFunc.NativeImpl(nil, envInstance, []environment.Value{key})
	assert.NoError(t, err)
	assert.Equal(t, "TEST_VALUE", result.String())
	
	// Clean up
	os.Unsetenv("TEST_KEY")
}

func TestENVContains(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSYSTEMInEnv(env)
	
	envVariable, _ := env.GetVariable("ENV")
	envInstance := envVariable.Value.(*environment.ObjectInstance)
	
	// Set a test environment variable directly
	os.Setenv("CONTAINS_TEST", "value")
	defer os.Unsetenv("CONTAINS_TEST")
	
	// Test CONTAINS - should find environment variable
	containsFunc := envInstance.Class.PublicFunctions["CONTAINS"]
	key := environment.StringValue("CONTAINS_TEST")
	
	result, err := containsFunc.NativeImpl(nil, envInstance, []environment.Value{key})
	assert.NoError(t, err)
	assert.Equal(t, environment.YEZ, result)
	
	// Test with non-existent key
	nonExistentKey := environment.StringValue("NON_EXISTENT_KEY")
	result, err = containsFunc.NativeImpl(nil, envInstance, []environment.Value{nonExistentKey})
	assert.NoError(t, err)
	assert.Equal(t, environment.NO, result)
}

func TestENVRemove(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSYSTEMInEnv(env)
	
	envVariable, _ := env.GetVariable("ENV")
	envInstance := envVariable.Value.(*environment.ObjectInstance)
	
	// Set a test environment variable
	testKey := "REMOVE_TEST"
	testValue := "remove_value"
	os.Setenv(testKey, testValue)
	
	// First ensure it's loaded into BASKIT by calling GET
	getFunc := envInstance.Class.PublicFunctions["GET"]
	key := environment.StringValue(testKey)
	getFunc.NativeImpl(nil, envInstance, []environment.Value{key})
	
	// Test REMOVE - should remove from both BASKIT and environment
	removeFunc := envInstance.Class.PublicFunctions["REMOVE"]
	result, err := removeFunc.NativeImpl(nil, envInstance, []environment.Value{key})
	assert.NoError(t, err)
	assert.Equal(t, testValue, result.String())
	
	// Verify environment variable was actually unset
	actualValue := os.Getenv(testKey)
	assert.Equal(t, "", actualValue)
	
	// Test removing non-existent key should return error
	nonExistentKey := environment.StringValue("NON_EXISTENT_REMOVE")
	_, err = removeFunc.NativeImpl(nil, envInstance, []environment.Value{nonExistentKey})
	assert.Error(t, err)
}

func TestENVRefresh(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSYSTEMInEnv(env)
	
	envVariable, _ := env.GetVariable("ENV")
	envInstance := envVariable.Value.(*environment.ObjectInstance)
	
	// Set some environment variables after ENV was created
	os.Setenv("REFRESH_TEST_1", "value1")
	os.Setenv("REFRESH_TEST_2", "value2")
	defer func() {
		os.Unsetenv("REFRESH_TEST_1")
		os.Unsetenv("REFRESH_TEST_2")
	}()
	
	// Call REFRESH to sync with current environment
	refreshFunc := envInstance.Class.PublicFunctions["REFRESH"]
	result, err := refreshFunc.NativeImpl(nil, envInstance, []environment.Value{})
	assert.NoError(t, err)
	assert.Equal(t, environment.NOTHIN, result)
	
	// Verify the new environment variables are now accessible via GET
	getFunc := envInstance.Class.PublicFunctions["GET"]
	key1 := environment.StringValue("REFRESH_TEST_1")
	result, err = getFunc.NativeImpl(nil, envInstance, []environment.Value{key1})
	assert.NoError(t, err)
	assert.Equal(t, "value1", result.String())
	
	key2 := environment.StringValue("REFRESH_TEST_2")
	result, err = getFunc.NativeImpl(nil, envInstance, []environment.Value{key2})
	assert.NoError(t, err)
	assert.Equal(t, "value2", result.String())
}

func TestENVClear(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSYSTEMInEnv(env)
	
	envVariable, _ := env.GetVariable("ENV")
	envInstance := envVariable.Value.(*environment.ObjectInstance)
	
	// Set some test environment variables via PUT
	putFunc := envInstance.Class.PublicFunctions["PUT"]
	key1 := environment.StringValue("CLEAR_TEST_1")
	value1 := environment.StringValue("value1")
	key2 := environment.StringValue("CLEAR_TEST_2")
	value2 := environment.StringValue("value2")
	
	putFunc.NativeImpl(nil, envInstance, []environment.Value{key1, value1})
	putFunc.NativeImpl(nil, envInstance, []environment.Value{key2, value2})
	
	// Verify they were set
	assert.Equal(t, "value1", os.Getenv("CLEAR_TEST_1"))
	assert.Equal(t, "value2", os.Getenv("CLEAR_TEST_2"))
	
	// Call CLEAR
	clearFunc := envInstance.Class.PublicFunctions["CLEAR"]
	result, err := clearFunc.NativeImpl(nil, envInstance, []environment.Value{})
	assert.NoError(t, err)
	assert.Equal(t, environment.NOTHIN, result)
	
	// Verify environment variables were unset
	assert.Equal(t, "", os.Getenv("CLEAR_TEST_1"))
	assert.Equal(t, "", os.Getenv("CLEAR_TEST_2"))
	
	// Verify BASKIT is empty (check SIZ)
	sizVar := envInstance.Class.PublicVariables["SIZ"]
	size, err := sizVar.NativeGet(envInstance)
	assert.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(0), size)
}

func TestENVCopy(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSYSTEMInEnv(env)
	
	envVariable, _ := env.GetVariable("ENV")
	envInstance := envVariable.Value.(*environment.ObjectInstance)
	
	// Set some test data
	putFunc := envInstance.Class.PublicFunctions["PUT"]
	key := environment.StringValue("COPY_TEST")
	value := environment.StringValue("copy_value")
	putFunc.NativeImpl(nil, envInstance, []environment.Value{key, value})
	
	// Call COPY
	copyFunc := envInstance.Class.PublicFunctions["COPY"]
	result, err := copyFunc.NativeImpl(nil, envInstance, []environment.Value{})
	assert.NoError(t, err)
	
	// Verify result is an ENVBASKIT instance
	copyInstance, ok := result.(*environment.ObjectInstance)
	assert.True(t, ok)
	assert.Equal(t, "ENVBASKIT", copyInstance.Class.Name)
	
	// Verify copied data
	getFunc := copyInstance.Class.PublicFunctions["GET"]
	copiedValue, err := getFunc.NativeImpl(nil, copyInstance, []environment.Value{key})
	assert.NoError(t, err)
	assert.Equal(t, "copy_value", copiedValue.String())
	
	// Clean up
	os.Unsetenv("COPY_TEST")
}

func TestENVSizeProperty(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSYSTEMInEnv(env)
	
	envVariable, _ := env.GetVariable("ENV")
	envInstance := envVariable.Value.(*environment.ObjectInstance)
	
	// Get initial size
	sizVar := envInstance.Class.PublicVariables["SIZ"]
	initialSize, err := sizVar.NativeGet(envInstance)
	assert.NoError(t, err)
	
	// Add some variables
	putFunc := envInstance.Class.PublicFunctions["PUT"]
	key1 := environment.StringValue("SIZE_TEST_1")
	value1 := environment.StringValue("value1")
	key2 := environment.StringValue("SIZE_TEST_2")
	value2 := environment.StringValue("value2")
	
	putFunc.NativeImpl(nil, envInstance, []environment.Value{key1, value1})
	putFunc.NativeImpl(nil, envInstance, []environment.Value{key2, value2})
	
	// Check size increased
	newSize, err := sizVar.NativeGet(envInstance)
	assert.NoError(t, err)
	assert.Equal(t, initialSize.(environment.IntegerValue)+2, newSize)
	
	// Clean up
	os.Unsetenv("SIZE_TEST_1")
	os.Unsetenv("SIZE_TEST_2")
}

func TestENVInheritedMethods(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSYSTEMInEnv(env)
	
	envVariable, _ := env.GetVariable("ENV")
	envInstance := envVariable.Value.(*environment.ObjectInstance)
	
	// Add some test data
	putFunc := envInstance.Class.PublicFunctions["PUT"]
	key1 := environment.StringValue("INHERIT_TEST_1")
	value1 := environment.StringValue("value1")
	key2 := environment.StringValue("INHERIT_TEST_2")
	value2 := environment.StringValue("value2")
	
	putFunc.NativeImpl(nil, envInstance, []environment.Value{key1, value1})
	putFunc.NativeImpl(nil, envInstance, []environment.Value{key2, value2})
	
	// Test KEYS method (inherited from BASKIT)
	keysFunc := envInstance.Class.PublicFunctions["KEYS"]
	result, err := keysFunc.NativeImpl(nil, envInstance, []environment.Value{})
	assert.NoError(t, err)
	
	// Should return a BUKKIT instance
	keysInstance, ok := result.(*environment.ObjectInstance)
	assert.True(t, ok)
	assert.Equal(t, "BUKKIT", keysInstance.Class.Name)
	
	// Test VALUES method (inherited from BASKIT)
	valuesFunc := envInstance.Class.PublicFunctions["VALUES"]
	result, err = valuesFunc.NativeImpl(nil, envInstance, []environment.Value{})
	assert.NoError(t, err)
	
	// Should return a BUKKIT instance
	valuesInstance, ok := result.(*environment.ObjectInstance)
	assert.True(t, ok)
	assert.Equal(t, "BUKKIT", valuesInstance.Class.Name)
	
	// Clean up
	os.Unsetenv("INHERIT_TEST_1")
	os.Unsetenv("INHERIT_TEST_2")
}

func TestENVNoPublicConstructor(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSYSTEMInEnv(env)
	
	class, err := env.GetClass("ENVBASKIT")
	assert.NoError(t, err)
	
	// Verify constructor is not in public functions
	_, hasPublicConstructor := class.PublicFunctions["ENVBASKIT"]
	assert.False(t, hasPublicConstructor, "ENVBASKIT should not have a public constructor")
	
	// Verify constructor is in private functions
	_, hasPrivateConstructor := class.PrivateFunctions["ENVBASKIT"]
	assert.True(t, hasPrivateConstructor, "ENVBASKIT should have a private constructor")
}

func TestENVGetFromActualEnvironment(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSYSTEMInEnv(env)
	
	envVariable, _ := env.GetVariable("ENV")
	envInstance := envVariable.Value.(*environment.ObjectInstance)
	
	// Set an environment variable directly (not through ENV)
	testKey := "DIRECT_ENV_TEST"
	testValue := "direct_value"
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)
	
	// GET should find it in the actual environment and load it into BASKIT
	getFunc := envInstance.Class.PublicFunctions["GET"]
	key := environment.StringValue(testKey)
	result, err := getFunc.NativeImpl(nil, envInstance, []environment.Value{key})
	assert.NoError(t, err)
	assert.Equal(t, testValue, result.String())
	
	// Now it should also be in the BASKIT, so a second GET should work the same way
	result2, err := getFunc.NativeImpl(nil, envInstance, []environment.Value{key})
	assert.NoError(t, err)
	assert.Equal(t, testValue, result2.String())
}

func TestSelectiveImportENVOnly(t *testing.T) {
	env := environment.NewEnvironment(nil)
	
	// Import only ENV variable
	err := RegisterSYSTEMInEnv(env, "ENV")
	assert.NoError(t, err)
	
	// ENV should be available
	envVariable, err := env.GetVariable("ENV")
	assert.NoError(t, err)
	assert.Equal(t, "ENVBASKIT", envVariable.Type)
	
	// ENVBASKIT class should also be available since ENV requires it
	class, err := env.GetClass("ENVBASKIT")
	assert.NoError(t, err, "ENVBASKIT class should be automatically imported when importing ENV")
	assert.Equal(t, "ENVBASKIT", class.Name)
	
	// ENV should work correctly
	envInstance := envVariable.Value.(*environment.ObjectInstance)
	putFunc := envInstance.Class.PublicFunctions["PUT"]
	key := environment.StringValue("DEPENDENCY_TEST")
	value := environment.StringValue("works")
	
	result, err := putFunc.NativeImpl(nil, envInstance, []environment.Value{key, value})
	assert.NoError(t, err)
	assert.Equal(t, environment.NOTHIN, result)
	
	// Clean up
	os.Unsetenv("DEPENDENCY_TEST")
}

func TestSelectiveImportENVBASKITOnly(t *testing.T) {
	env := environment.NewEnvironment(nil)
	
	// Import only ENVBASKIT class
	err := RegisterSYSTEMInEnv(env, "ENVBASKIT")
	assert.NoError(t, err)
	
	// ENVBASKIT class should be available
	class, err := env.GetClass("ENVBASKIT")
	assert.NoError(t, err)
	assert.Equal(t, "ENVBASKIT", class.Name)
	
	// ENV variable should not be available since we didn't import it
	_, err = env.GetVariable("ENV")
	assert.Error(t, err, "ENV variable should not be available without explicit import")
}

func TestSelectiveImportMultiple(t *testing.T) {
	env := environment.NewEnvironment(nil)
	
	// Import both ENV and ENVBASKIT
	err := RegisterSYSTEMInEnv(env, "ENV", "ENVBASKIT")
	assert.NoError(t, err)
	
	// Both should be available
	envVariable, err := env.GetVariable("ENV")
	assert.NoError(t, err)
	assert.Equal(t, "ENVBASKIT", envVariable.Type)
	
	class, err := env.GetClass("ENVBASKIT")
	assert.NoError(t, err)
	assert.Equal(t, "ENVBASKIT", class.Name)
}

func TestSelectiveImportCaseInsensitive(t *testing.T) {
	env := environment.NewEnvironment(nil)
	
	// Import with lowercase - should still work
	err := RegisterSYSTEMInEnv(env, "env", "envbaskit")
	assert.NoError(t, err)
	
	// Both should be available
	envVariable, err := env.GetVariable("ENV")
	assert.NoError(t, err)
	assert.Equal(t, "ENVBASKIT", envVariable.Type)
	
	class, err := env.GetClass("ENVBASKIT")
	assert.NoError(t, err)
	assert.Equal(t, "ENVBASKIT", class.Name)
}

func TestSelectiveImportUnknownDeclaration(t *testing.T) {
	env := environment.NewEnvironment(nil)
	
	// Try to import unknown declaration
	err := RegisterSYSTEMInEnv(env, "UNKNOWN_THING")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown SYSTEM declaration: UNKNOWN_THING")
}