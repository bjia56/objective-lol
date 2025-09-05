package environment

import (
	"reflect"
	"testing"
)

func TestC3LinearizationSimple(t *testing.T) {
	env := NewEnvironment(nil)

	// Create a simple inheritance hierarchy: C -> A, B -> A
	classA := NewClass("A", "main.olol", []string{})
	classB := NewClass("B", "main.olol", []string{"main.olol.A"})
	classC := NewClass("C", "main.olol", []string{"main.olol.A"})
	classD := NewClass("D", "main.olol", []string{"main.olol.B", "main.olol.C"})

	env.DefineClass(classA)
	env.DefineClass(classB)
	env.DefineClass(classC)
	env.DefineClass(classD)

	// Test MRO computation
	mro, err := env.computeC3Linearization(classD)
	if err != nil {
		t.Fatalf("Failed to compute C3 linearization: %v", err)
	}

	expected := []string{"main.olol.D", "main.olol.B", "main.olol.C", "main.olol.A"}
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
	classA := NewClass("A", "main.olol", []string{})
	classB := NewClass("B", "main.olol", []string{"main.olol.A"})
	classC := NewClass("C", "main.olol", []string{"main.olol.A"})
	classD := NewClass("D", "main.olol", []string{"main.olol.B", "main.olol.C"})

	env.DefineClass(classA)
	env.DefineClass(classB)
	env.DefineClass(classC)
	env.DefineClass(classD)

	mro, err := env.computeC3Linearization(classD)
	if err != nil {
		t.Fatalf("Failed to compute C3 linearization: %v", err)
	}

	expected := []string{"main.olol.D", "main.olol.B", "main.olol.C", "main.olol.A"}
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
	classO := NewClass("O", "main.olol", []string{})
	classA := NewClass("A", "main.olol", []string{"main.olol.O"})
	classB := NewClass("B", "main.olol", []string{"main.olol.O"})
	classC := NewClass("C", "main.olol", []string{"main.olol.O"})
	classD := NewClass("D", "main.olol", []string{"main.olol.A"})
	classE := NewClass("E", "main.olol", []string{"main.olol.C"})
	classF := NewClass("F", "main.olol", []string{"main.olol.D", "main.olol.E"})

	env.DefineClass(classO)
	env.DefineClass(classA)
	env.DefineClass(classB)
	env.DefineClass(classC)
	env.DefineClass(classD)
	env.DefineClass(classE)
	env.DefineClass(classF)

	mro, err := env.computeC3Linearization(classF)
	if err != nil {
		t.Fatalf("Failed to compute C3 linearization: %v", err)
	}

	expected := []string{"main.olol.F", "main.olol.D", "main.olol.A", "main.olol.E", "main.olol.C", "main.olol.O"}
	if !reflect.DeepEqual(mro, expected) {
		t.Errorf("Expected MRO %v, got %v", expected, mro)
	}
}

func TestC3LinearizationConflict(t *testing.T) {
	env := NewEnvironment(nil)

	// Create problematic inheritance that should fail C3 linearization
	// This creates an inconsistent hierarchy
	classA := NewClass("A", "main.olol", []string{})
	classB := NewClass("B", "main.olol", []string{"main.olol.A"})
	classC := NewClass("C", "main.olol", []string{"main.olol.A"})
	classX := NewClass("X", "main.olol", []string{"main.olol.B", "main.olol.C"})
	classY := NewClass("Y", "main.olol", []string{"main.olol.C", "main.olol.B"}) // Different order
	classZ := NewClass("Z", "main.olol", []string{"main.olol.X", "main.olol.Y"})

	env.DefineClass(classA)
	env.DefineClass(classB)
	env.DefineClass(classC)
	env.DefineClass(classX)
	env.DefineClass(classY)
	env.DefineClass(classZ)

	_, err := env.computeC3Linearization(classZ)
	if err == nil {
		t.Error("Expected C3 linearization to fail for inconsistent hierarchy, but it succeeded")
	}
}

func TestC3LinearizationNoParents(t *testing.T) {
	env := NewEnvironment(nil)

	// Test class with no parents
	classA := NewClass("A", "main.olol", []string{})
	env.DefineClass(classA)

	mro, err := env.computeC3Linearization(classA)
	if err != nil {
		t.Fatalf("Failed to compute C3 linearization: %v", err)
	}

	expected := []string{"main.olol.A"}
	if !reflect.DeepEqual(mro, expected) {
		t.Errorf("Expected MRO %v, got %v", expected, mro)
	}
}

func TestC3LinearizationSingleParent(t *testing.T) {
	env := NewEnvironment(nil)

	// Test simple single inheritance
	classA := NewClass("A", "main.olol", []string{})
	classB := NewClass("B", "main.olol", []string{"main.olol.A"})

	env.DefineClass(classA)
	env.DefineClass(classB)

	mro, err := env.computeC3Linearization(classB)
	if err != nil {
		t.Fatalf("Failed to compute C3 linearization: %v", err)
	}

	expected := []string{"main.olol.B", "main.olol.A"}
	if !reflect.DeepEqual(mro, expected) {
		t.Errorf("Expected MRO %v, got %v", expected, mro)
	}
}

func TestMROCaching(t *testing.T) {
	env := NewEnvironment(nil)

	// Create simple hierarchy
	classA := NewClass("A", "main.olol", []string{})
	classB := NewClass("B", "main.olol", []string{"main.olol.A"})

	env.DefineClass(classA)
	env.DefineClass(classB)

	// First computation should cache MRO
	mro1, err := env.computeOrGetMRO(classB)
	if err != nil {
		t.Fatalf("Failed to compute MRO: %v", err)
	}

	// Second computation should use cached result
	mro2, err := env.computeOrGetMRO(classB)
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
	classA := NewClass("A", "main.olol", []string{})
	classA.PublicVariables["A_VAR"] = &Variable{
		Name: "A_VAR", Type: "STRIN", Value: StringValue("from A"), IsPublic: true,
	}

	classB := NewClass("B", "main.olol", []string{})
	classB.PublicVariables["B_VAR"] = &Variable{
		Name: "B_VAR", Type: "STRIN", Value: StringValue("from B"), IsPublic: true,
	}

	classC := NewClass("C", "main.olol", []string{"A", "B"})
	classC.PublicVariables["C_VAR"] = &Variable{
		Name: "C_VAR", Type: "STRIN", Value: StringValue("from C"), IsPublic: true,
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
	expectedMRO := []string{"main.olol.C", "main.olol.A", "main.olol.B"}
	if !reflect.DeepEqual(instance.MRO, expectedMRO) {
		t.Errorf("Expected instance MRO %v, got %v", expectedMRO, instance.MRO)
	}

	// Verify all variables are initialized
	if _, exists := instance.Variables["A_VAR"]; !exists {
		t.Error("A_VAR not initialized in instance")
	}
	if _, exists := instance.Variables["B_VAR"]; !exists {
		t.Error("B_VAR not initialized in instance")
	}
	if _, exists := instance.Variables["C_VAR"]; !exists {
		t.Error("C_VAR not initialized in instance")
	}
}
