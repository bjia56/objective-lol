package environment

import (
	"testing"
	"reflect"
	"github.com/bjia56/objective-lol/pkg/types"
)

func TestC3LinearizationSimple(t *testing.T) {
	env := NewEnvironment(nil)
	
	// Create a simple inheritance hierarchy: C -> A, B -> A
	classA := NewClass("A", "", []string{})
	classB := NewClass("B", "", []string{"A"})
	classC := NewClass("C", "", []string{"A"})
	classD := NewClass("D", "", []string{"B", "C"})
	
	env.DefineClass(classA)
	env.DefineClass(classB)
	env.DefineClass(classC)
	env.DefineClass(classD)
	
	// Test MRO computation
	mro, err := env.computeC3Linearization("D")
	if err != nil {
		t.Fatalf("Failed to compute C3 linearization: %v", err)
	}
	
	expected := []string{"D", "B", "C", "A"}
	if !reflect.DeepEqual(mro, expected) {
		t.Errorf("Expected MRO %v, got %v", expected, mro)
	}
}

func TestC3LinearizationDiamond(t *testing.T) {
	env := NewEnvironment(nil)
	
	// Create diamond inheritance: D -> B, C -> A
	//        A
	//       / \
	//      B   C
	//       \ /
	//        D
	classA := NewClass("A", "", []string{})
	classB := NewClass("B", "", []string{"A"})
	classC := NewClass("C", "", []string{"A"})
	classD := NewClass("D", "", []string{"B", "C"})
	
	env.DefineClass(classA)
	env.DefineClass(classB)
	env.DefineClass(classC)
	env.DefineClass(classD)
	
	mro, err := env.computeC3Linearization("D")
	if err != nil {
		t.Fatalf("Failed to compute C3 linearization: %v", err)
	}
	
	expected := []string{"D", "B", "C", "A"}
	if !reflect.DeepEqual(mro, expected) {
		t.Errorf("Expected MRO %v, got %v", expected, mro)
	}
}

func TestC3LinearizationComplex(t *testing.T) {
	env := NewEnvironment(nil)
	
	// Create complex inheritance hierarchy from Python's C3 examples
	//    O
	//   /|\
	//  A B C
	//  |   |
	//  D   E
	//   \ /
	//    F
	classO := NewClass("O", "", []string{})
	classA := NewClass("A", "", []string{"O"})
	classB := NewClass("B", "", []string{"O"})
	classC := NewClass("C", "", []string{"O"})
	classD := NewClass("D", "", []string{"A"})
	classE := NewClass("E", "", []string{"C"})
	classF := NewClass("F", "", []string{"D", "E"})
	
	env.DefineClass(classO)
	env.DefineClass(classA)
	env.DefineClass(classB)
	env.DefineClass(classC)
	env.DefineClass(classD)
	env.DefineClass(classE)
	env.DefineClass(classF)
	
	mro, err := env.computeC3Linearization("F")
	if err != nil {
		t.Fatalf("Failed to compute C3 linearization: %v", err)
	}
	
	expected := []string{"F", "D", "A", "E", "C", "O"}
	if !reflect.DeepEqual(mro, expected) {
		t.Errorf("Expected MRO %v, got %v", expected, mro)
	}
}

func TestC3LinearizationConflict(t *testing.T) {
	env := NewEnvironment(nil)
	
	// Create problematic inheritance that should fail C3 linearization
	// This creates an inconsistent hierarchy
	classA := NewClass("A", "", []string{})
	classB := NewClass("B", "", []string{"A"})
	classC := NewClass("C", "", []string{"A"})
	classX := NewClass("X", "", []string{"B", "C"})
	classY := NewClass("Y", "", []string{"C", "B"}) // Different order
	classZ := NewClass("Z", "", []string{"X", "Y"})
	
	env.DefineClass(classA)
	env.DefineClass(classB)
	env.DefineClass(classC)
	env.DefineClass(classX)
	env.DefineClass(classY)
	env.DefineClass(classZ)
	
	_, err := env.computeC3Linearization("Z")
	if err == nil {
		t.Error("Expected C3 linearization to fail for inconsistent hierarchy, but it succeeded")
	}
}

func TestC3LinearizationNoParents(t *testing.T) {
	env := NewEnvironment(nil)
	
	// Test class with no parents
	classA := NewClass("A", "", []string{})
	env.DefineClass(classA)
	
	mro, err := env.computeC3Linearization("A")
	if err != nil {
		t.Fatalf("Failed to compute C3 linearization: %v", err)
	}
	
	expected := []string{"A"}
	if !reflect.DeepEqual(mro, expected) {
		t.Errorf("Expected MRO %v, got %v", expected, mro)
	}
}

func TestC3LinearizationSingleParent(t *testing.T) {
	env := NewEnvironment(nil)
	
	// Test simple single inheritance
	classA := NewClass("A", "", []string{})
	classB := NewClass("B", "", []string{"A"})
	
	env.DefineClass(classA)
	env.DefineClass(classB)
	
	mro, err := env.computeC3Linearization("B")
	if err != nil {
		t.Fatalf("Failed to compute C3 linearization: %v", err)
	}
	
	expected := []string{"B", "A"}
	if !reflect.DeepEqual(mro, expected) {
		t.Errorf("Expected MRO %v, got %v", expected, mro)
	}
}

func TestMROCaching(t *testing.T) {
	env := NewEnvironment(nil)
	
	// Create simple hierarchy
	classA := NewClass("A", "", []string{})
	classB := NewClass("B", "", []string{"A"})
	
	env.DefineClass(classA)
	env.DefineClass(classB)
	
	// First computation should cache MRO
	mro1, err := env.computeOrGetMRO("B")
	if err != nil {
		t.Fatalf("Failed to compute MRO: %v", err)
	}
	
	// Second computation should use cached result
	mro2, err := env.computeOrGetMRO("B")
	if err != nil {
		t.Fatalf("Failed to get cached MRO: %v", err)
	}
	
	if !reflect.DeepEqual(mro1, mro2) {
		t.Errorf("Cached MRO differs from computed MRO: %v vs %v", mro1, mro2)
	}
	
	// Verify MRO is actually cached in the class
	classB, _ = env.GetClass("B")
	if len(classB.MRO) == 0 {
		t.Error("MRO was not cached in class structure")
	}
}

func TestObjectInstanceWithMultipleInheritance(t *testing.T) {
	env := NewEnvironment(nil)
	
	// Create inheritance hierarchy
	classA := NewClass("A", "", []string{})
	classA.PublicVariables["A_VAR"] = &Variable{
		Name: "A_VAR", Type: "STRIN", Value: types.StringValue("from A"), IsPublic: true,
	}
	
	classB := NewClass("B", "", []string{})
	classB.PublicVariables["B_VAR"] = &Variable{
		Name: "B_VAR", Type: "STRIN", Value: types.StringValue("from B"), IsPublic: true,
	}
	
	classC := NewClass("C", "", []string{"A", "B"})
	classC.PublicVariables["C_VAR"] = &Variable{
		Name: "C_VAR", Type: "STRIN", Value: types.StringValue("from C"), IsPublic: true,
	}
	
	env.DefineClass(classA)
	env.DefineClass(classB)
	env.DefineClass(classC)
	
	// Create instance
	instance, err := env.NewObjectInstance("C")
	if err != nil {
		t.Fatalf("Failed to create object instance: %v", err)
	}
	
	// Verify MRO is stored
	expectedMRO := []string{"C", "A", "B"}
	objInstance := instance.(*ObjectInstance)
	if !reflect.DeepEqual(objInstance.MRO, expectedMRO) {
		t.Errorf("Expected instance MRO %v, got %v", expectedMRO, objInstance.MRO)
	}
	
	// Verify all variables are initialized
	if _, exists := objInstance.Variables["A_VAR"]; !exists {
		t.Error("A_VAR not initialized in instance")
	}
	if _, exists := objInstance.Variables["B_VAR"]; !exists {
		t.Error("B_VAR not initialized in instance")
	}
	if _, exists := objInstance.Variables["C_VAR"]; !exists {
		t.Error("C_VAR not initialized in instance")
	}
}