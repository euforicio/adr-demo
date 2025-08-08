package validator

import "testing"

func TestValidateSequentialNumberingIgnoresFileOrder(t *testing.T) {
	v := New(&Config{})
	filenames := []string{"0002-second.md", "0001-first.md", "0003-third.md"}
	result := &ValidationResult{}
	if err := v.validateSequentialNumbering(filenames, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ErrorCount != 0 {
		t.Fatalf("expected no errors, got %d: %#v", result.ErrorCount, result.Issues)
	}
}
