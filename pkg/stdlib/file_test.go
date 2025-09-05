package stdlib

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/interpreter"
)

func TestDOCUMENTConstructor(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterFILEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register FILE module: %v", err)
	}

	// Get DOCUMENT class
	docClass, err := env.GetClass("DOCUMENT")
	if err != nil {
		t.Fatal("DOCUMENT class not found")
	}

	// Create interpreter context
	ctx := interpreter.NewFunctionContext(nil, env)

	// Test valid constructor
	instance := &environment.ObjectInstance{
		Class:     docClass,
		Variables: make(map[string]*environment.Variable),
	}

	constructor := docClass.PublicFunctions["DOCUMENT"]
	args := []environment.Value{
		environment.StringValue("test.txt"),
		environment.StringValue("W"),
	}

	_, err = constructor.NativeImpl(ctx, instance, args)
	if err != nil {
		t.Errorf("Constructor failed: %v", err)
	}

	// Verify native data was set
	docData, ok := instance.NativeData.(*DocumentData)
	if !ok {
		t.Fatal("Native data not set correctly")
	}

	if docData.FilePath != "test.txt" {
		t.Errorf("Expected FilePath 'test.txt', got '%s'", docData.FilePath)
	}

	if docData.FileMode != "W" {
		t.Errorf("Expected FileMode 'W', got '%s'", docData.FileMode)
	}

	if docData.IsOpen {
		t.Error("Expected IsOpen to be false initially")
	}

	// Test invalid mode
	instance2 := &environment.ObjectInstance{
		Class:     docClass,
		Variables: make(map[string]*environment.Variable),
	}

	invalidArgs := []environment.Value{
		environment.StringValue("test.txt"),
		environment.StringValue("INVALID"),
	}

	_, err = constructor.NativeImpl(ctx, instance2, invalidArgs)
	if err == nil {
		t.Error("Expected constructor to fail with invalid mode")
	}
}

func TestDOCUMENTFileOperations(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "document_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")

	env := environment.NewEnvironment(nil)
	err = RegisterFILEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register FILE module: %v", err)
	}

	docClass, _ := env.GetClass("DOCUMENT")
	ctx := interpreter.NewFunctionContext(nil, env)

	// Create DOCUMENT instance
	instance := &environment.ObjectInstance{
		Class:     docClass,
		Variables: make(map[string]*environment.Variable),
	}

	constructor := docClass.PublicFunctions["DOCUMENT"]
	args := []environment.Value{
		environment.StringValue(testFile),
		environment.StringValue("W"),
	}

	_, err = constructor.NativeImpl(ctx, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Test OPEN
	openFunc := docClass.PublicFunctions["OPEN"]
	_, err = openFunc.NativeImpl(ctx, instance, []environment.Value{})
	if err != nil {
		t.Errorf("OPEN failed: %v", err)
	}

	// Test WRITE
	writeFunc := docClass.PublicFunctions["WRITE"]
	testData := "Hello, World!"
	result, err := writeFunc.NativeImpl(ctx, instance, []environment.Value{
		environment.StringValue(testData),
	})
	if err != nil {
		t.Errorf("WRITE failed: %v", err)
	}

	// Verify bytes written
	if bytesWritten, ok := result.(environment.IntegerValue); ok {
		if int(bytesWritten) != len(testData) {
			t.Errorf("Expected %d bytes written, got %d", len(testData), int(bytesWritten))
		}
	} else {
		t.Error("WRITE should return integer")
	}

	// Test FLUSH
	flushFunc := docClass.PublicFunctions["FLUSH"]
	_, err = flushFunc.NativeImpl(ctx, instance, []environment.Value{})
	if err != nil {
		t.Errorf("FLUSH failed: %v", err)
	}

	// Test SIZE
	sizeFunc := docClass.PublicFunctions["SIZE"]
	sizeResult, err := sizeFunc.NativeImpl(ctx, instance, []environment.Value{})
	if err != nil {
		t.Errorf("SIZE failed: %v", err)
	}

	if size, ok := sizeResult.(environment.IntegerValue); ok {
		if int(size) != len(testData) {
			t.Errorf("Expected size %d, got %d", len(testData), int(size))
		}
	} else {
		t.Error("SIZE should return integer")
	}

	// Test CLOSE
	closeFunc := docClass.PublicFunctions["CLOSE"]
	_, err = closeFunc.NativeImpl(ctx, instance, []environment.Value{})
	if err != nil {
		t.Errorf("CLOSE failed: %v", err)
	}

	// Verify file was written
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("Failed to read test file: %v", err)
	}

	if string(content) != testData {
		t.Errorf("Expected file content '%s', got '%s'", testData, string(content))
	}
}

func TestDOCUMENTReadOperations(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "document_read_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "read_test.txt")
	testContent := "Hello, World!\nThis is a test file."

	// Write test content to file
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	env := environment.NewEnvironment(nil)
	err = RegisterFILEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register FILE module: %v", err)
	}

	docClass, _ := env.GetClass("DOCUMENT")
	ctx := interpreter.NewFunctionContext(nil, env)

	// Create DOCUMENT instance for reading
	instance := &environment.ObjectInstance{
		Class:     docClass,
		Variables: make(map[string]*environment.Variable),
	}

	constructor := docClass.PublicFunctions["DOCUMENT"]
	args := []environment.Value{
		environment.StringValue(testFile),
		environment.StringValue("R"),
	}

	_, err = constructor.NativeImpl(ctx, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Open file for reading
	openFunc := docClass.PublicFunctions["OPEN"]
	_, err = openFunc.NativeImpl(ctx, instance, []environment.Value{})
	if err != nil {
		t.Errorf("OPEN failed: %v", err)
	}

	// Test READ
	readFunc := docClass.PublicFunctions["READ"]
	readResult, err := readFunc.NativeImpl(ctx, instance, []environment.Value{
		environment.IntegerValue(5),
	})
	if err != nil {
		t.Errorf("READ failed: %v", err)
	}

	if content, ok := readResult.(environment.StringValue); ok {
		if string(content) != "Hello" {
			t.Errorf("Expected 'Hello', got '%s'", string(content))
		}
	} else {
		t.Error("READ should return string")
	}

	// Test TELL (current position)
	tellFunc := docClass.PublicFunctions["TELL"]
	posResult, err := tellFunc.NativeImpl(ctx, instance, []environment.Value{})
	if err != nil {
		t.Errorf("TELL failed: %v", err)
	}

	if pos, ok := posResult.(environment.IntegerValue); ok {
		if int(pos) != 5 {
			t.Errorf("Expected position 5, got %d", int(pos))
		}
	} else {
		t.Error("TELL should return integer")
	}

	// Test SEEK
	seekFunc := docClass.PublicFunctions["SEEK"]
	_, err = seekFunc.NativeImpl(ctx, instance, []environment.Value{
		environment.IntegerValue(0),
	})
	if err != nil {
		t.Errorf("SEEK failed: %v", err)
	}

	// Read entire content
	fullReadResult, err := readFunc.NativeImpl(ctx, instance, []environment.Value{
		environment.IntegerValue(100),
	})
	if err != nil {
		t.Errorf("Full READ failed: %v", err)
	}

	if content, ok := fullReadResult.(environment.StringValue); ok {
		if string(content) != testContent {
			t.Errorf("Expected '%s', got '%s'", testContent, string(content))
		}
	}

	// Test EXISTS
	existsFunc := docClass.PublicFunctions["EXISTS"]
	existsResult, err := existsFunc.NativeImpl(ctx, instance, []environment.Value{})
	if err != nil {
		t.Errorf("EXISTS failed: %v", err)
	}

	if exists, ok := existsResult.(environment.BoolValue); ok {
		if !bool(exists) {
			t.Error("EXISTS should return true for existing file")
		}
	} else {
		t.Error("EXISTS should return boolean")
	}

	// Close file
	closeFunc := docClass.PublicFunctions["CLOSE"]
	_, err = closeFunc.NativeImpl(ctx, instance, []environment.Value{})
	if err != nil {
		t.Errorf("CLOSE failed: %v", err)
	}
}

func TestDOCUMENTErrorHandling(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterFILEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register FILE module: %v", err)
	}

	docClass, _ := env.GetClass("DOCUMENT")
	ctx := interpreter.NewFunctionContext(nil, env)

	// Create DOCUMENT instance
	instance := &environment.ObjectInstance{
		Class:     docClass,
		Variables: make(map[string]*environment.Variable),
	}

	constructor := docClass.PublicFunctions["DOCUMENT"]
	args := []environment.Value{
		environment.StringValue("nonexistent.txt"),
		environment.StringValue("R"),
	}

	_, err = constructor.NativeImpl(ctx, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Test operations on closed file should fail
	readFunc := docClass.PublicFunctions["READ"]
	_, err = readFunc.NativeImpl(ctx, instance, []environment.Value{
		environment.IntegerValue(10),
	})
	if err == nil {
		t.Error("READ on closed file should fail")
	}

	writeFunc := docClass.PublicFunctions["WRITE"]
	_, err = writeFunc.NativeImpl(ctx, instance, []environment.Value{
		environment.StringValue("test"),
	})
	if err == nil {
		t.Error("WRITE on closed file should fail")
	}

	// Test opening nonexistent file for reading should fail
	openFunc := docClass.PublicFunctions["OPEN"]
	_, err = openFunc.NativeImpl(ctx, instance, []environment.Value{})
	if err == nil {
		t.Error("OPEN nonexistent file for reading should fail")
	}
}

func TestDOCUMENTDeleteOperation(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "document_delete_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "delete_test.txt")

	// Create test file
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	env := environment.NewEnvironment(nil)
	err = RegisterFILEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register FILE module: %v", err)
	}

	docClass, _ := env.GetClass("DOCUMENT")
	ctx := interpreter.NewFunctionContext(nil, env)

	// Create DOCUMENT instance
	instance := &environment.ObjectInstance{
		Class:     docClass,
		Variables: make(map[string]*environment.Variable),
	}

	constructor := docClass.PublicFunctions["DOCUMENT"]
	args := []environment.Value{
		environment.StringValue(testFile),
		environment.StringValue("R"),
	}

	_, err = constructor.NativeImpl(ctx, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Test DELETE
	deleteFunc := docClass.PublicFunctions["DELETE"]
	_, err = deleteFunc.NativeImpl(ctx, instance, []environment.Value{})
	if err != nil {
		t.Errorf("DELETE failed: %v", err)
	}

	// Verify file was deleted
	_, err = os.Stat(testFile)
	if err == nil {
		t.Error("File should have been deleted")
	}
}