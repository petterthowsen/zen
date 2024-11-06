package interpreter

import (
	"testing"
)

func TestPrint(t *testing.T) {
	_, err := InterpretString("print(\"hello, zen!\")")
	if err != nil {
		t.Errorf("Failed to interpret code: %v", err)
	}
}
