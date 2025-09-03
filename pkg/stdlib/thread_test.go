package stdlib

import (
	"sync"
	"testing"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

func TestRegisterTHREADInEnv(t *testing.T) {
	env := environment.NewEnvironment(nil)
	
	// Test full import
	err := RegisterTHREADInEnv(env)
	if err != nil {
		t.Errorf("Failed to register THREAD module: %v", err)
	}
	
	// Check that YARN class is registered
	yarnClass, err := env.GetClass("YARN")
	if err != nil || yarnClass == nil {
		t.Error("YARN class not found after registration")
	}
	
	// Check that KNOT class is registered
	knotClass, err := env.GetClass("KNOT")
	if err != nil || knotClass == nil {
		t.Error("KNOT class not found after registration")
	}
}

func TestRegisterTHREADInEnvSelective(t *testing.T) {
	env := environment.NewEnvironment(nil)
	
	// Test selective import - only KNOT
	err := RegisterTHREADInEnv(env, "KNOT")
	if err != nil {
		t.Errorf("Failed to register selective THREAD import: %v", err)
	}
	
	// Check that only KNOT class is registered
	knotClass, err := env.GetClass("KNOT")
	if err != nil || knotClass == nil {
		t.Error("KNOT class not found after selective registration")
	}
	
	// Check that YARN class is NOT registered
	yarnClass, err := env.GetClass("YARN")
	if err == nil && yarnClass != nil {
		t.Error("YARN class should not be registered in selective import")
	}
}

func TestRegisterTHREADInEnvInvalidDeclaration(t *testing.T) {
	env := environment.NewEnvironment(nil)
	
	// Test invalid declaration
	err := RegisterTHREADInEnv(env, "INVALID_CLASS")
	if err == nil {
		t.Error("Expected error for invalid declaration, got nil")
	}
}

func TestNewYarnInstance(t *testing.T) {
	yarnObj := NewYarnInstance()
	
	// Check class hierarchy
	if len(yarnObj.Hierarchy) != 1 || yarnObj.Hierarchy[0] != "YARN" {
		t.Errorf("Expected hierarchy [YARN], got %v", yarnObj.Hierarchy)
	}
	
	// Check that NativeData is ThreadData
	if _, ok := yarnObj.NativeData.(*ThreadData); !ok {
		t.Error("Expected NativeData to be *ThreadData")
	}
	
	// Check initial variable values
	if runningVar, exists := yarnObj.Variables["RUNNING"]; exists {
		if runningVar.Value != types.NO {
			t.Error("Expected RUNNING to be initially NO")
		}
	} else {
		t.Error("RUNNING variable not found")
	}
	
	if finishedVar, exists := yarnObj.Variables["FINISHED"]; exists {
		if finishedVar.Value != types.NO {
			t.Error("Expected FINISHED to be initially NO")
		}
	} else {
		t.Error("FINISHED variable not found")
	}
}

func TestNewKnotInstance(t *testing.T) {
	knotObj := NewKnotInstance()
	
	// Check class hierarchy
	if len(knotObj.Hierarchy) != 1 || knotObj.Hierarchy[0] != "KNOT" {
		t.Errorf("Expected hierarchy [KNOT], got %v", knotObj.Hierarchy)
	}
	
	// Check that NativeData is MutexData
	if _, ok := knotObj.NativeData.(*MutexData); !ok {
		t.Error("Expected NativeData to be *MutexData")
	}
	
	// Check initial variable values
	if lockedVar, exists := knotObj.Variables["LOCKED"]; exists {
		if lockedVar.Value != types.NO {
			t.Error("Expected LOCKED to be initially NO")
		}
	} else {
		t.Error("LOCKED variable not found")
	}
}

func TestYarnConstructor(t *testing.T) {
	classes := getThreadClasses()
	yarnClass := classes["YARN"]
	
	constructor := yarnClass.PublicFunctions["YARN"]
	if constructor == nil {
		t.Fatal("YARN constructor not found")
	}
	
	yarnObj := NewYarnInstance()
	result, err := constructor.NativeImpl(nil, yarnObj, []types.Value{})
	
	if err != nil {
		t.Errorf("Constructor failed: %v", err)
	}
	
	if result != types.NOTHIN {
		t.Errorf("Expected constructor to return NOTHIN, got %v", result)
	}
}

func TestYarnSpinAbstract(t *testing.T) {
	classes := getThreadClasses()
	yarnClass := classes["YARN"]
	
	spinMethod := yarnClass.PublicFunctions["SPIN"]
	if spinMethod == nil {
		t.Fatal("SPIN method not found")
	}
	
	yarnObj := NewYarnInstance()
	result, err := spinMethod.NativeImpl(nil, yarnObj, []types.Value{})
	
	if err == nil {
		t.Error("Expected SPIN to throw exception for abstract method")
	}
	
	if result != types.NOTHIN {
		t.Errorf("Expected SPIN to return NOTHIN, got %v", result)
	}
}

func TestKnotConstructor(t *testing.T) {
	classes := getThreadClasses()
	knotClass := classes["KNOT"]
	
	constructor := knotClass.PublicFunctions["KNOT"]
	if constructor == nil {
		t.Fatal("KNOT constructor not found")
	}
	
	knotObj := NewKnotInstance()
	result, err := constructor.NativeImpl(nil, knotObj, []types.Value{})
	
	if err != nil {
		t.Errorf("Constructor failed: %v", err)
	}
	
	if result != types.NOTHIN {
		t.Errorf("Expected constructor to return NOTHIN, got %v", result)
	}
}

func TestKnotTieUntie(t *testing.T) {
	classes := getThreadClasses()
	knotClass := classes["KNOT"]
	
	tieMethod := knotClass.PublicFunctions["TIE"]
	untieMethod := knotClass.PublicFunctions["UNTIE"]
	
	if tieMethod == nil {
		t.Fatal("TIE method not found")
	}
	if untieMethod == nil {
		t.Fatal("UNTIE method not found")
	}
	
	knotObj := NewKnotInstance()
	
	// Test TIE (lock)
	result, err := tieMethod.NativeImpl(nil, knotObj, []types.Value{})
	if err != nil {
		t.Errorf("TIE failed: %v", err)
	}
	if result != types.NOTHIN {
		t.Errorf("Expected TIE to return NOTHIN, got %v", result)
	}
	
	// Check that LOCKED status is updated
	if lockedVar, exists := knotObj.Variables["LOCKED"]; exists {
		if lockedVar.Value != types.YEZ {
			t.Error("Expected LOCKED to be YEZ after TIE")
		}
	}
	
	// Test UNTIE (unlock)
	result, err = untieMethod.NativeImpl(nil, knotObj, []types.Value{})
	if err != nil {
		t.Errorf("UNTIE failed: %v", err)
	}
	if result != types.NOTHIN {
		t.Errorf("Expected UNTIE to return NOTHIN, got %v", result)
	}
	
	// Check that LOCKED status is updated
	if lockedVar, exists := knotObj.Variables["LOCKED"]; exists {
		if lockedVar.Value != types.NO {
			t.Error("Expected LOCKED to be NO after UNTIE")
		}
	}
}

func TestKnotUntieWhenNotLocked(t *testing.T) {
	classes := getThreadClasses()
	knotClass := classes["KNOT"]
	
	untieMethod := knotClass.PublicFunctions["UNTIE"]
	knotObj := NewKnotInstance()
	
	// Try to UNTIE when not locked - should throw exception
	result, err := untieMethod.NativeImpl(nil, knotObj, []types.Value{})
	if err == nil {
		t.Error("Expected UNTIE to throw exception when mutex not locked")
	}
	if result != types.NOTHIN {
		t.Errorf("Expected UNTIE to return NOTHIN, got %v", result)
	}
}

func TestKnotConcurrentAccess(t *testing.T) {
	classes := getThreadClasses()
	knotClass := classes["KNOT"]
	
	tieMethod := knotClass.PublicFunctions["TIE"]
	untieMethod := knotClass.PublicFunctions["UNTIE"]
	
	knotObj := NewKnotInstance()
	
	// Test concurrent access to verify actual mutex behavior
	var counter int
	var wg sync.WaitGroup
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			// Lock
			tieMethod.NativeImpl(nil, knotObj, []types.Value{})
			
			// Critical section
			temp := counter
			time.Sleep(1 * time.Millisecond)
			counter = temp + 1
			
			// Unlock
			untieMethod.NativeImpl(nil, knotObj, []types.Value{})
		}()
	}
	
	wg.Wait()
	
	if counter != 10 {
		t.Errorf("Expected counter to be 10 with proper mutex synchronization, got %d", counter)
	}
}

func TestUpdateYarnStatus(t *testing.T) {
	yarnObj := NewYarnInstance()
	threadData := &ThreadData{
		goroutineRunning: true,
		finished:         false,
	}
	
	updateYarnStatus(yarnObj, threadData)
	
	if runningVar, exists := yarnObj.Variables["RUNNING"]; exists {
		if runningVar.Value != types.YEZ {
			t.Error("Expected RUNNING to be YEZ after update")
		}
	}
	
	if finishedVar, exists := yarnObj.Variables["FINISHED"]; exists {
		if finishedVar.Value != types.NO {
			t.Error("Expected FINISHED to be NO after update")
		}
	}
}

func TestUpdateKnotStatus(t *testing.T) {
	knotObj := NewKnotInstance()
	mutexData := &MutexData{
		locked: true,
	}
	
	updateKnotStatus(knotObj, mutexData)
	
	if lockedVar, exists := knotObj.Variables["LOCKED"]; exists {
		if lockedVar.Value != types.YEZ {
			t.Error("Expected LOCKED to be YEZ after update")
		}
	}
}

// TestYarnStartWithoutInterpreter tests START method behavior without proper interpreter context
func TestYarnStartWithoutInterpreter(t *testing.T) {
	classes := getThreadClasses()
	yarnClass := classes["YARN"]
	
	startMethod := yarnClass.PublicFunctions["START"]
	yarnObj := NewYarnInstance()
	
	// Call START without proper interpreter context
	result, err := startMethod.NativeImpl(nil, yarnObj, []types.Value{})
	
	if err == nil {
		t.Error("Expected START to fail without proper interpreter context")
	}
	
	if result != types.NOTHIN {
		t.Errorf("Expected START to return NOTHIN, got %v", result)
	}
}

// TestYarnJoin tests the JOIN method
func TestYarnJoin(t *testing.T) {
	classes := getThreadClasses()
	yarnClass := classes["YARN"]
	
	joinMethod := yarnClass.PublicFunctions["JOIN"]
	yarnObj := NewYarnInstance()
	
	// Setup thread data
	threadData := &ThreadData{
		finished: true,
		result:   types.NOTHIN, // Set result to NOTHIN
	}
	yarnObj.NativeData = threadData
	
	// Call JOIN - should return immediately since thread is marked as finished
	result, err := joinMethod.NativeImpl(nil, yarnObj, []types.Value{})
	
	if err != nil {
		t.Errorf("JOIN failed: %v", err)
	}
	
	// Should return the result from the thread
	if result != types.NOTHIN {
		t.Errorf("Expected JOIN to return NOTHIN, got %v", result)
	}
}