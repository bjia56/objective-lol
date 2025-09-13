package stdlib

import (
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
)

func TestRegisterSOCKETInEnv(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test registering all classes
	err := RegisterSOCKETInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register SOCKET classes: %v", err)
	}

	// Check if SOKKIT class is registered
	sokkitClass, err := env.GetClass("SOKKIT")
	if err != nil {
		t.Fatal("SOKKIT class not found after registration")
	}

	// Check if WIRE class is registered
	wireClass, err := env.GetClass("WIRE")
	if err != nil {
		t.Fatal("WIRE class not found after registration")
	}

	// Verify SOKKIT class properties
	if sokkitClass.Name != "SOKKIT" {
		t.Errorf("Expected SOKKIT class name, got %s", sokkitClass.Name)
	}

	// Verify WIRE class properties
	if wireClass.Name != "WIRE" {
		t.Errorf("Expected WIRE class name, got %s", wireClass.Name)
	}
}

func TestRegisterSOCKETInEnvSelective(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test registering only SOKKIT
	err := RegisterSOCKETInEnv(env, "SOKKIT")
	if err != nil {
		t.Fatalf("Failed to register SOKKIT class: %v", err)
	}

	// Check if SOKKIT class is registered
	_, err = env.GetClass("SOKKIT")
	if err != nil {
		t.Fatal("SOKKIT class not found after selective registration")
	}

	// Check if WIRE class is also registered (dependency)
	_, err = env.GetClass("WIRE")
	if err != nil {
		t.Fatal("WIRE class not found after SOKKIT registration (should be auto-imported)")
	}
}

func TestRegisterSOCKETInEnvInvalidClass(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test registering invalid class
	err := RegisterSOCKETInEnv(env, "INVALID_CLASS")
	if err == nil {
		t.Fatal("Expected error when registering invalid class")
	}

	expectedError := "unknown SOCKET declaration: INVALID_CLASS"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSOCKITClassFunctions(t *testing.T) {
	socketClasses := getSocketClasses()
	sokkitClass := socketClasses["SOKKIT"]

	expectedFunctions := []string{
		"SOKKIT", "BIND", "LISTEN", "ACCEPT", "CONNECT",
		"SEND_TO", "RECEIVE_FROM", "CLOSE",
	}

	for _, funcName := range expectedFunctions {
		if _, exists := sokkitClass.PublicFunctions[funcName]; !exists {
			t.Errorf("SOKKIT class missing expected function: %s", funcName)
		}
	}
}

func TestWIREClassFunctions(t *testing.T) {
	socketClasses := getSocketClasses()
	wireClass := socketClasses["WIRE"]

	expectedFunctions := []string{
		"SEND", "RECEIVE", "RECEIVE_ALL", "CLOSE",
	}

	for _, funcName := range expectedFunctions {
		if _, exists := wireClass.PublicFunctions[funcName]; !exists {
			t.Errorf("WIRE class missing expected function: %s", funcName)
		}
	}
}

func TestSOCKITClassVariables(t *testing.T) {
	socketClasses := getSocketClasses()
	sokkitClass := socketClasses["SOKKIT"]

	expectedVariables := map[string]string{
		"PROTOCOL": "STRIN",
		"HOST":     "STRIN",
		"PORT":     "INTEGR",
		"TIMEOUT":  "INTEGR",
	}

	for varName, expectedType := range expectedVariables {
		variable, exists := sokkitClass.PublicVariables[varName]
		if !exists {
			t.Errorf("SOKKIT class missing expected variable: %s", varName)
			continue
		}
		if variable.Type != expectedType {
			t.Errorf("SOKKIT variable %s expected type %s, got %s", varName, expectedType, variable.Type)
		}
	}
}

func TestWIREClassVariables(t *testing.T) {
	socketClasses := getSocketClasses()
	wireClass := socketClasses["WIRE"]

	expectedVariables := map[string]string{
		"REMOTE_HOST":  "STRIN",
		"REMOTE_PORT":  "INTEGR",
		"LOCAL_HOST":   "STRIN",
		"LOCAL_PORT":   "INTEGR",
		"IS_CONNECTED": "BOOL",
	}

	for varName, expectedType := range expectedVariables {
		variable, exists := wireClass.PublicVariables[varName]
		if !exists {
			t.Errorf("WIRE class missing expected variable: %s", varName)
			continue
		}
		if variable.Type != expectedType {
			t.Errorf("WIRE variable %s expected type %s, got %s", varName, expectedType, variable.Type)
		}
	}
}
