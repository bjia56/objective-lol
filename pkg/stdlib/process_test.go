package stdlib

import (
	"strings"
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
)

func TestRegisterPROCESSInEnv(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterPROCESSInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register PROCESS in env: %v", err)
	}

	// Check that classes were registered
	classes := env.GetAllClasses()
	if _, exists := classes["PIPE"]; !exists {
		t.Error("PIPE class not registered")
	}
	if _, exists := classes["MINION"]; !exists {
		t.Error("MINION class not registered")
	}
}

func TestRegisterPROCESSInEnvSelective(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterPROCESSInEnv(env, "MINION")
	if err != nil {
		t.Fatalf("Failed to register PROCESS classes selectively: %v", err)
	}

	classes := env.GetAllClasses()
	if _, exists := classes["MINION"]; !exists {
		t.Error("MINION class not registered")
	}
}

func TestMinionConstructor(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)
	RegisterPROCESSInEnv(env)

	// Create a BUKKIT for command line
	bukkit, err := env.NewObjectInstance("BUKKIT")
	if err != nil {
		t.Fatalf("Failed to create BUKKIT: %v", err)
	}
	bukkit.NativeData = BukkitSlice{
		environment.StringValue("echo"),
		environment.StringValue("hello world"),
	}

	// Create MINION
	minionClass := env.GetAllClasses()["MINION"]
	minion := &environment.ObjectInstance{
		Class:     minionClass,
		Variables: make(map[string]*environment.Variable),
	}

	_, err = minionClass.PublicFunctions["MINION"].NativeImpl(nil, minion, []environment.Value{bukkit})
	if err != nil {
		t.Fatalf("MINION constructor failed: %v", err)
	}

	// Check that variables were set
	if cmdlineVar, exists := minion.Variables["CMDLINE"]; !exists {
		t.Error("CMDLINE variable not set")
	} else {
		if cmdlineVar.Value != bukkit {
			t.Error("CMDLINE variable not set correctly")
		}
	}

	if runningVar, exists := minion.Variables["RUNNING"]; !exists {
		t.Error("RUNNING variable not set")
	} else {
		if runningVar.Value != environment.NO {
			t.Error("RUNNING should be NO initially")
		}
	}
}

func TestMinionSetWorkdir(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)
	RegisterPROCESSInEnv(env)

	// Create a BUKKIT for command line
	bukkit, err := env.NewObjectInstance("BUKKIT")
	if err != nil {
		t.Fatalf("Failed to create BUKKIT: %v", err)
	}
	bukkit.NativeData = BukkitSlice{
		environment.StringValue("pwd"),
	}

	// Create MINION
	minionClass := env.GetAllClasses()["MINION"]
	minion := &environment.ObjectInstance{
		Class:     minionClass,
		Variables: make(map[string]*environment.Variable),
	}

	_, err = minionClass.PublicFunctions["MINION"].NativeImpl(nil, minion, []environment.Value{bukkit})
	if err != nil {
		t.Fatalf("MINION constructor failed: %v", err)
	}

	// Test SET_WORKDIR
	_, err = minionClass.PublicFunctions["SET_WORKDIR"].NativeImpl(nil, minion, []environment.Value{
		environment.StringValue("/tmp"),
	})
	if err != nil {
		t.Fatalf("SET_WORKDIR failed: %v", err)
	}

	// Check that working directory was set
	minionData := minion.NativeData.(*MinionData)
	if minionData.WorkDir != "/tmp" {
		t.Errorf("Expected WorkDir to be '/tmp', got '%s'", minionData.WorkDir)
	}
}

func TestMinionAddEnv(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)
	RegisterPROCESSInEnv(env)

	// Create a BUKKIT for command line
	bukkit, err := env.NewObjectInstance("BUKKIT")
	if err != nil {
		t.Fatalf("Failed to create BUKKIT: %v", err)
	}
	bukkit.NativeData = BukkitSlice{
		environment.StringValue("env"),
	}

	// Create MINION
	minionClass := env.GetAllClasses()["MINION"]
	minion := &environment.ObjectInstance{
		Class:     minionClass,
		Variables: make(map[string]*environment.Variable),
	}

	_, err = minionClass.PublicFunctions["MINION"].NativeImpl(nil, minion, []environment.Value{bukkit})
	if err != nil {
		t.Fatalf("MINION constructor failed: %v", err)
	}

	// Test ADD_ENV
	_, err = minionClass.PublicFunctions["ADD_ENV"].NativeImpl(nil, minion, []environment.Value{
		environment.StringValue("TEST_VAR"),
		environment.StringValue("test_value"),
	})
	if err != nil {
		t.Fatalf("ADD_ENV failed: %v", err)
	}

	// Check that environment variable was added
	minionData := minion.NativeData.(*MinionData)
	found := false
	for _, envVar := range minionData.Env {
		if strings.HasPrefix(envVar, "TEST_VAR=test_value") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Environment variable TEST_VAR=test_value not found")
	}
}

func TestMinionBasicExecution(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)
	RegisterPROCESSInEnv(env)

	// Create a BUKKIT for command line
	bukkit, err := env.NewObjectInstance("BUKKIT")
	if err != nil {
		t.Fatalf("Failed to create BUKKIT: %v", err)
	}
	bukkit.NativeData = BukkitSlice{
		environment.StringValue("echo"),
		environment.StringValue("hello world"),
	}

	// Create MINION
	minionClass := env.GetAllClasses()["MINION"]
	minion := &environment.ObjectInstance{
		Class:     minionClass,
		Variables: make(map[string]*environment.Variable),
	}

	_, err = minionClass.PublicFunctions["MINION"].NativeImpl(nil, minion, []environment.Value{bukkit})
	if err != nil {
		t.Fatalf("MINION constructor failed: %v", err)
	}

	// Start the process
	_, err = minionClass.PublicFunctions["START"].NativeImpl(nil, minion, []environment.Value{})
	if err != nil {
		t.Fatalf("START failed: %v", err)
	}

	// Check that process started
	if runningVar, exists := minion.Variables["RUNNING"]; !exists {
		t.Error("RUNNING variable not found")
	} else {
		if runningVar.Value != environment.YEZ {
			t.Error("RUNNING should be YEZ after start")
		}
	}

	// Check that PID was set
	if pidVar, exists := minion.Variables["PID"]; !exists {
		t.Error("PID variable not found")
	} else {
		if pidVal, ok := pidVar.Value.(environment.IntegerValue); !ok || int(pidVal) <= 0 {
			t.Error("PID should be positive after start")
		}
	}

	// Check that pipes were created
	if stdinVar, exists := minion.Variables["STDIN"]; !exists {
		t.Error("STDIN variable not found")
	} else {
		if stdinVar.Value == environment.NOTHIN {
			t.Error("STDIN should be set after start")
		}
	}

	if stdoutVar, exists := minion.Variables["STDOUT"]; !exists {
		t.Error("STDOUT variable not found")
	} else {
		if stdoutVar.Value == environment.NOTHIN {
			t.Error("STDOUT should be set after start")
		}
	}

	// Wait for process to complete
	result, err := minionClass.PublicFunctions["WAIT"].NativeImpl(nil, minion, []environment.Value{})
	if err != nil {
		t.Fatalf("WAIT failed: %v", err)
	}

	// Check exit code
	if exitCode, ok := result.(environment.IntegerValue); !ok || int(exitCode) != 0 {
		t.Errorf("Expected exit code 0, got %v", result)
	}

	// Check that process is finished
	if finishedVar, exists := minion.Variables["FINISHED"]; !exists {
		t.Error("FINISHED variable not found")
	} else {
		if finishedVar.Value != environment.YEZ {
			t.Error("FINISHED should be YEZ after wait")
		}
	}
}

func TestPipeReadWrite(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)
	RegisterPROCESSInEnv(env)

	// Create a BUKKIT for command line (cat command)
	bukkit, err := env.NewObjectInstance("BUKKIT")
	if err != nil {
		t.Fatalf("Failed to create BUKKIT: %v", err)
	}
	bukkit.NativeData = BukkitSlice{
		environment.StringValue("cat"),
	}

	// Create MINION
	minionClass := env.GetAllClasses()["MINION"]
	minion := &environment.ObjectInstance{
		Class:     minionClass,
		Variables: make(map[string]*environment.Variable),
	}

	_, err = minionClass.PublicFunctions["MINION"].NativeImpl(nil, minion, []environment.Value{bukkit})
	if err != nil {
		t.Fatalf("MINION constructor failed: %v", err)
	}

	// Start the process
	_, err = minionClass.PublicFunctions["START"].NativeImpl(nil, minion, []environment.Value{})
	if err != nil {
		t.Fatalf("START failed: %v", err)
	}

	// Get stdin and stdout pipes
	stdinPipe := minion.Variables["STDIN"].Value.(*environment.ObjectInstance)
	stdoutPipe := minion.Variables["STDOUT"].Value.(*environment.ObjectInstance)

	pipeClass := env.GetAllClasses()["PIPE"]

	// Write to stdin
	testData := "hello pipe test\n"
	_, err = pipeClass.PublicFunctions["WRITE"].NativeImpl(nil, stdinPipe, []environment.Value{
		environment.StringValue(testData),
	})
	if err != nil {
		t.Fatalf("WRITE to stdin failed: %v", err)
	}

	// Close stdin to signal EOF to cat
	_, err = pipeClass.PublicFunctions["CLOSE"].NativeImpl(nil, stdinPipe, []environment.Value{})
	if err != nil {
		t.Fatalf("CLOSE stdin failed: %v", err)
	}

	// Read from stdout
	output, err := pipeClass.PublicFunctions["READ"].NativeImpl(nil, stdoutPipe, []environment.Value{
		environment.IntegerValue(100),
	})
	if err != nil {
		t.Fatalf("READ from stdout failed: %v", err)
	}

	// Check output
	if outputStr, ok := output.(environment.StringValue); !ok {
		t.Error("READ should return StringValue")
	} else {
		if string(outputStr) != testData {
			t.Errorf("Expected output '%s', got '%s'", testData, string(outputStr))
		}
	}

	// Wait for process to complete
	_, err = minionClass.PublicFunctions["WAIT"].NativeImpl(nil, minion, []environment.Value{})
	if err != nil {
		t.Fatalf("WAIT failed: %v", err)
	}
}

func TestMinionIsAlive(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)
	RegisterPROCESSInEnv(env)

	// Create a BUKKIT for command line (sleep command)
	bukkit, err := env.NewObjectInstance("BUKKIT")
	if err != nil {
		t.Fatalf("Failed to create BUKKIT: %v", err)
	}
	bukkit.NativeData = BukkitSlice{
		environment.StringValue("sleep"),
		environment.StringValue("0.1"),
	}

	// Create MINION
	minionClass := env.GetAllClasses()["MINION"]
	minion := &environment.ObjectInstance{
		Class:     minionClass,
		Variables: make(map[string]*environment.Variable),
	}

	_, err = minionClass.PublicFunctions["MINION"].NativeImpl(nil, minion, []environment.Value{bukkit})
	if err != nil {
		t.Fatalf("MINION constructor failed: %v", err)
	}

	// Check IS_ALIVE before starting
	result, err := minionClass.PublicFunctions["IS_ALIVE"].NativeImpl(nil, minion, []environment.Value{})
	if err != nil {
		t.Fatalf("IS_ALIVE failed: %v", err)
	}
	if result != environment.NO {
		t.Error("IS_ALIVE should be NO before start")
	}

	// Start the process
	_, err = minionClass.PublicFunctions["START"].NativeImpl(nil, minion, []environment.Value{})
	if err != nil {
		t.Fatalf("START failed: %v", err)
	}

	// Check IS_ALIVE after starting
	result, err = minionClass.PublicFunctions["IS_ALIVE"].NativeImpl(nil, minion, []environment.Value{})
	if err != nil {
		t.Fatalf("IS_ALIVE failed: %v", err)
	}
	if result != environment.YEZ {
		t.Error("IS_ALIVE should be YEZ after start")
	}

	// Wait for process to complete
	_, err = minionClass.PublicFunctions["WAIT"].NativeImpl(nil, minion, []environment.Value{})
	if err != nil {
		t.Fatalf("WAIT failed: %v", err)
	}

	// Check IS_ALIVE after completion
	result, err = minionClass.PublicFunctions["IS_ALIVE"].NativeImpl(nil, minion, []environment.Value{})
	if err != nil {
		t.Fatalf("IS_ALIVE failed: %v", err)
	}
	if result != environment.NO {
		t.Error("IS_ALIVE should be NO after completion")
	}
}

func TestMinionErrorHandling(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)
	RegisterPROCESSInEnv(env)

	// Test empty command line
	emptyBukkit, err := env.NewObjectInstance("BUKKIT")
	if err != nil {
		t.Fatalf("Failed to create BUKKIT: %v", err)
	}
	emptyBukkit.NativeData = BukkitSlice{}

	minionClass := env.GetAllClasses()["MINION"]
	minion := &environment.ObjectInstance{
		Class:     minionClass,
		Variables: make(map[string]*environment.Variable),
	}

	_, err = minionClass.PublicFunctions["MINION"].NativeImpl(nil, minion, []environment.Value{emptyBukkit})
	if err == nil {
		t.Error("MINION constructor should fail with empty command line")
	}

	// Test starting before construction
	minion2 := &environment.ObjectInstance{
		Class:     minionClass,
		Variables: make(map[string]*environment.Variable),
	}
	_, err = minionClass.PublicFunctions["START"].NativeImpl(nil, minion2, []environment.Value{})
	if err == nil {
		t.Error("START should fail when MINION not properly constructed")
	}

	// Test wait without start
	bukkit, err := env.NewObjectInstance("BUKKIT")
	if err != nil {
		t.Fatalf("Failed to create BUKKIT: %v", err)
	}
	bukkit.NativeData = BukkitSlice{
		environment.StringValue("echo"),
		environment.StringValue("test"),
	}

	minion3 := &environment.ObjectInstance{
		Class:     minionClass,
		Variables: make(map[string]*environment.Variable),
	}
	_, err = minionClass.PublicFunctions["MINION"].NativeImpl(nil, minion3, []environment.Value{bukkit})
	if err != nil {
		t.Fatalf("MINION constructor failed: %v", err)
	}

	_, err = minionClass.PublicFunctions["WAIT"].NativeImpl(nil, minion3, []environment.Value{})
	if err == nil {
		t.Error("WAIT should fail when process not started")
	}
}
