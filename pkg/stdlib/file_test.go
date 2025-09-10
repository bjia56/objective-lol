package stdlib

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
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

	// Test valid constructor
	instance := &environment.ObjectInstance{
		Class:     docClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := docClass.PublicFunctions["DOCUMENT"]
	args := []environment.Value{
		environment.StringValue("test.txt"),
		environment.StringValue("W"),
	}

	_, err = constructor.NativeImpl(nil, instance, args)
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
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance2)

	invalidArgs := []environment.Value{
		environment.StringValue("test.txt"),
		environment.StringValue("INVALID"),
	}

	_, err = constructor.NativeImpl(nil, instance2, invalidArgs)
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

	// Create DOCUMENT instance
	instance := &environment.ObjectInstance{
		Class:     docClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := docClass.PublicFunctions["DOCUMENT"]
	args := []environment.Value{
		environment.StringValue(testFile),
		environment.StringValue("W"),
	}

	_, err = constructor.NativeImpl(nil, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Test OPEN
	openFunc := docClass.PublicFunctions["OPEN"]
	_, err = openFunc.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Errorf("OPEN failed: %v", err)
	}

	// Test WRITE
	writeFunc := docClass.PublicFunctions["WRITE"]
	testData := "Hello, World!"
	result, err := writeFunc.NativeImpl(nil, instance, []environment.Value{
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
	_, err = flushFunc.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Errorf("FLUSH failed: %v", err)
	}

	// Test SIZ property
	sizeMember := docClass.PublicVariables["SIZ"]
	sizeResult, err := sizeMember.NativeGet(instance)
	if err != nil {
		t.Errorf("SIZ property failed: %v", err)
	}

	if size, ok := sizeResult.(environment.IntegerValue); ok {
		if int(size) != len(testData) {
			t.Errorf("Expected size %d, got %d", len(testData), int(size))
		}
	} else {
		t.Error("SIZ property should return integer")
	}

	// Test CLOSE
	closeFunc := docClass.PublicFunctions["CLOSE"]
	_, err = closeFunc.NativeImpl(nil, instance, []environment.Value{})
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

	// Create DOCUMENT instance for reading
	instance := &environment.ObjectInstance{
		Class:     docClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := docClass.PublicFunctions["DOCUMENT"]
	args := []environment.Value{
		environment.StringValue(testFile),
		environment.StringValue("R"),
	}

	_, err = constructor.NativeImpl(nil, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Open file for reading
	openFunc := docClass.PublicFunctions["OPEN"]
	_, err = openFunc.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Errorf("OPEN failed: %v", err)
	}

	// Test READ
	readFunc := docClass.PublicFunctions["READ"]
	readResult, err := readFunc.NativeImpl(nil, instance, []environment.Value{
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
	posResult, err := tellFunc.NativeImpl(nil, instance, []environment.Value{})
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
	_, err = seekFunc.NativeImpl(nil, instance, []environment.Value{
		environment.IntegerValue(0),
	})
	if err != nil {
		t.Errorf("SEEK failed: %v", err)
	}

	// Read entire content
	fullReadResult, err := readFunc.NativeImpl(nil, instance, []environment.Value{
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
	existsResult, err := existsFunc.NativeImpl(nil, instance, []environment.Value{})
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
	_, err = closeFunc.NativeImpl(nil, instance, []environment.Value{})
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

	// Create DOCUMENT instance
	instance := &environment.ObjectInstance{
		Class:     docClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := docClass.PublicFunctions["DOCUMENT"]
	args := []environment.Value{
		environment.StringValue("nonexistent.txt"),
		environment.StringValue("R"),
	}

	_, err = constructor.NativeImpl(nil, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Test operations on closed file should fail
	readFunc := docClass.PublicFunctions["READ"]
	_, err = readFunc.NativeImpl(nil, instance, []environment.Value{
		environment.IntegerValue(10),
	})
	if err == nil {
		t.Error("READ on closed file should fail")
	}

	writeFunc := docClass.PublicFunctions["WRITE"]
	_, err = writeFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("test"),
	})
	if err == nil {
		t.Error("WRITE on closed file should fail")
	}

	// Test opening nonexistent file for reading should fail
	openFunc := docClass.PublicFunctions["OPEN"]
	_, err = openFunc.NativeImpl(nil, instance, []environment.Value{})
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

	// Create DOCUMENT instance
	instance := &environment.ObjectInstance{
		Class:     docClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := docClass.PublicFunctions["DOCUMENT"]
	args := []environment.Value{
		environment.StringValue(testFile),
		environment.StringValue("R"),
	}

	_, err = constructor.NativeImpl(nil, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Test DELETE
	deleteFunc := docClass.PublicFunctions["DELETE"]
	_, err = deleteFunc.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Errorf("DELETE failed: %v", err)
	}

	// Verify file was deleted
	_, err = os.Stat(testFile)
	if err == nil {
		t.Error("File should have been deleted")
	}
}

func TestCABINETConstructor(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterFILEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register FILE module: %v", err)
	}

	// Get CABINET class
	cabinetClass, err := env.GetClass("CABINET")
	if err != nil {
		t.Fatal("CABINET class not found")
	}

	// Test valid constructor
	instance := &environment.ObjectInstance{
		Class:     cabinetClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := cabinetClass.PublicFunctions["CABINET"]
	args := []environment.Value{
		environment.StringValue("/test/dir"),
	}

	_, err = constructor.NativeImpl(nil, instance, args)
	if err != nil {
		t.Errorf("Constructor failed: %v", err)
	}

	// Verify native data was set
	cabinetData, ok := instance.NativeData.(*CabinetData)
	if !ok {
		t.Fatal("Native data not set correctly")
	}

	if cabinetData.DirPath != "/test/dir" {
		t.Errorf("Expected DirPath '/test/dir', got '%s'", cabinetData.DirPath)
	}
}

func TestCABINETDirectoryOperations(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "cabinet_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testSubDir := filepath.Join(tmpDir, "test_subdir")

	env := environment.NewEnvironment(nil)
	err = RegisterFILEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register FILE module: %v", err)
	}

	cabinetClass, _ := env.GetClass("CABINET")

	// Create CABINET instance
	instance := &environment.ObjectInstance{
		Class:     cabinetClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := cabinetClass.PublicFunctions["CABINET"]
	args := []environment.Value{
		environment.StringValue(testSubDir),
	}

	_, err = constructor.NativeImpl(nil, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Test EXISTS - should be false initially
	existsFunc := cabinetClass.PublicFunctions["EXISTS"]
	result, err := existsFunc.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Errorf("EXISTS failed: %v", err)
	}

	if exists, ok := result.(environment.BoolValue); ok {
		if bool(exists) {
			t.Error("Directory should not exist initially")
		}
	} else {
		t.Error("EXISTS should return boolean")
	}

	// Test CREATE
	createFunc := cabinetClass.PublicFunctions["CREATE"]
	_, err = createFunc.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Errorf("CREATE failed: %v", err)
	}

	// Test EXISTS - should be true now
	result, err = existsFunc.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Errorf("EXISTS after create failed: %v", err)
	}

	if exists, ok := result.(environment.BoolValue); ok {
		if !bool(exists) {
			t.Error("Directory should exist after CREATE")
		}
	}

	// Verify directory was actually created
	stat, err := os.Stat(testSubDir)
	if err != nil {
		t.Errorf("Directory was not created: %v", err)
	}
	if !stat.IsDir() {
		t.Error("Created path is not a directory")
	}

	// Test DELETE
	deleteFunc := cabinetClass.PublicFunctions["DELETE"]
	_, err = deleteFunc.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Errorf("DELETE failed: %v", err)
	}

	// Verify directory was deleted
	_, err = os.Stat(testSubDir)
	if err == nil {
		t.Error("Directory should have been deleted")
	}
}

func TestCABINETListOperations(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "cabinet_list_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := []string{"test1.txt", "test2.log", "test3.txt"}
	for _, filename := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		err = os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	env := environment.NewEnvironment(nil)
	err = RegisterFILEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register FILE module: %v", err)
	}

	cabinetClass, _ := env.GetClass("CABINET")

	// Create CABINET instance
	instance := &environment.ObjectInstance{
		Class:     cabinetClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := cabinetClass.PublicFunctions["CABINET"]
	args := []environment.Value{
		environment.StringValue(tmpDir),
	}

	_, err = constructor.NativeImpl(nil, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Test LIST
	listFunc := cabinetClass.PublicFunctions["LIST"]
	result, err := listFunc.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Errorf("LIST failed: %v", err)
	}

	// Verify result is a BUKKIT instance
	if bukkitInstance, ok := result.(*environment.ObjectInstance); ok {
		if bukkitInstance.Class.Name != "BUKKIT" {
			t.Errorf("Expected BUKKIT, got %s", bukkitInstance.Class.Name)
		}

		// Get the bukkit slice
		if bukkitSlice, ok := bukkitInstance.NativeData.(BukkitSlice); ok {
			if len(bukkitSlice) != len(testFiles) {
				t.Errorf("Expected %d files, got %d", len(testFiles), len(bukkitSlice))
			}

			// Verify file names (order may vary)
			foundFiles := make(map[string]bool)
			for _, value := range bukkitSlice {
				if strVal, ok := value.(environment.StringValue); ok {
					foundFiles[string(strVal)] = true
				}
			}

			for _, expectedFile := range testFiles {
				if !foundFiles[expectedFile] {
					t.Errorf("Expected file %s not found in list", expectedFile)
				}
			}
		} else {
			t.Error("BUKKIT native data is not BukkitSlice")
		}
	} else {
		t.Error("LIST should return BUKKIT instance")
	}
}

func TestCABINETFindOperations(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "cabinet_find_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files with different extensions
	testFiles := []string{"test1.txt", "test2.log", "test3.txt", "readme.md"}
	for _, filename := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		err = os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	env := environment.NewEnvironment(nil)
	err = RegisterFILEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register FILE module: %v", err)
	}

	cabinetClass, _ := env.GetClass("CABINET")

	// Create CABINET instance
	instance := &environment.ObjectInstance{
		Class:     cabinetClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := cabinetClass.PublicFunctions["CABINET"]
	args := []environment.Value{
		environment.StringValue(tmpDir),
	}

	_, err = constructor.NativeImpl(nil, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Test FIND with *.txt pattern
	findFunc := cabinetClass.PublicFunctions["FIND"]
	result, err := findFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("*.txt"),
	})
	if err != nil {
		t.Errorf("FIND failed: %v", err)
	}

	// Verify result is a BUKKIT instance with txt files
	if bukkitInstance, ok := result.(*environment.ObjectInstance); ok {
		if bukkitSlice, ok := bukkitInstance.NativeData.(BukkitSlice); ok {
			expectedTxtFiles := []string{"test1.txt", "test3.txt"}
			if len(bukkitSlice) != len(expectedTxtFiles) {
				t.Errorf("Expected %d txt files, got %d", len(expectedTxtFiles), len(bukkitSlice))
			}

			foundFiles := make(map[string]bool)
			for _, value := range bukkitSlice {
				if strVal, ok := value.(environment.StringValue); ok {
					foundFiles[string(strVal)] = true
				}
			}

			for _, expectedFile := range expectedTxtFiles {
				if !foundFiles[expectedFile] {
					t.Errorf("Expected txt file %s not found", expectedFile)
				}
			}
		}
	}

	// Test FIND with *.log pattern
	result, err = findFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("*.log"),
	})
	if err != nil {
		t.Errorf("FIND *.log failed: %v", err)
	}

	if bukkitInstance, ok := result.(*environment.ObjectInstance); ok {
		if bukkitSlice, ok := bukkitInstance.NativeData.(BukkitSlice); ok {
			if len(bukkitSlice) != 1 {
				t.Errorf("Expected 1 log file, got %d", len(bukkitSlice))
			}
			if strVal, ok := bukkitSlice[0].(environment.StringValue); ok {
				if string(strVal) != "test2.log" {
					t.Errorf("Expected test2.log, got %s", string(strVal))
				}
			}
		}
	}

	// Test FIND with pattern that matches nothing
	result, err = findFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("*.nonexistent"),
	})
	if err != nil {
		t.Errorf("FIND *.nonexistent failed: %v", err)
	}

	if bukkitInstance, ok := result.(*environment.ObjectInstance); ok {
		if bukkitSlice, ok := bukkitInstance.NativeData.(BukkitSlice); ok {
			if len(bukkitSlice) != 0 {
				t.Errorf("Expected 0 matches for *.nonexistent, got %d", len(bukkitSlice))
			}
		}
	}
}

func TestCABINETPathProperty(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterFILEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register FILE module: %v", err)
	}

	cabinetClass, _ := env.GetClass("CABINET")

	// Create CABINET instance
	instance := &environment.ObjectInstance{
		Class:     cabinetClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := cabinetClass.PublicFunctions["CABINET"]
	testPath := "/test/path"
	args := []environment.Value{
		environment.StringValue(testPath),
	}

	_, err = constructor.NativeImpl(nil, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Test PATH property
	pathMember := cabinetClass.PublicVariables["PATH"]
	pathResult, err := pathMember.NativeGet(instance)
	if err != nil {
		t.Errorf("PATH property failed: %v", err)
	}

	if path, ok := pathResult.(environment.StringValue); ok {
		if string(path) != testPath {
			t.Errorf("Expected path '%s', got '%s'", testPath, string(path))
		}
	} else {
		t.Error("PATH property should return string")
	}
}

func TestCABINETErrorHandling(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterFILEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register FILE module: %v", err)
	}

	cabinetClass, _ := env.GetClass("CABINET")

	// Create CABINET instance with nonexistent directory
	instance := &environment.ObjectInstance{
		Class:     cabinetClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := cabinetClass.PublicFunctions["CABINET"]
	args := []environment.Value{
		environment.StringValue("/nonexistent/directory"),
	}

	_, err = constructor.NativeImpl(nil, instance, args)
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Test LIST on nonexistent directory should fail
	listFunc := cabinetClass.PublicFunctions["LIST"]
	_, err = listFunc.NativeImpl(nil, instance, []environment.Value{})
	if err == nil {
		t.Error("LIST on nonexistent directory should fail")
	}

	// Test FIND on nonexistent directory should fail
	findFunc := cabinetClass.PublicFunctions["FIND"]
	_, err = findFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("*.txt"),
	})
	if err == nil {
		t.Error("FIND on nonexistent directory should fail")
	}

	// Test DELETE on nonexistent directory should fail
	deleteFunc := cabinetClass.PublicFunctions["DELETE"]
	_, err = deleteFunc.NativeImpl(nil, instance, []environment.Value{})
	if err == nil {
		t.Error("DELETE on nonexistent directory should fail")
	}

	// Test invalid pattern for FIND
	tmpDir, err := os.MkdirTemp("", "cabinet_error_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temp file
	os.Create(filepath.Join(tmpDir, "file.txt"))

	// Create valid CABINET for invalid pattern test
	validInstance := &environment.ObjectInstance{
		Class:     cabinetClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(validInstance)

	args = []environment.Value{
		environment.StringValue(tmpDir),
	}
	_, err = constructor.NativeImpl(nil, validInstance, args)
	if err != nil {
		t.Fatalf("Valid constructor failed: %v", err)
	}

	// Test FIND with invalid pattern (malformed bracket expression)
	_, err = findFunc.NativeImpl(nil, validInstance, []environment.Value{
		environment.StringValue("[]"),
	})
	if err == nil {
		t.Error("FIND with invalid pattern should fail")
	}
}
