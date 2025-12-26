package util

import (
	"reflect"
	"strings"
	"testing"
)

func TestBoolToUint8(t *testing.T) {
	if got := BoolToUint8(true); got != 1 {
		t.Fatalf("BoolToUint8(true) = %d, want 1", got)
	}
	if got := BoolToUint8(false); got != 0 {
		t.Fatalf("BoolToUint8(false) = %d, want 0", got)
	}
}

func TestIff(t *testing.T) {
	if got := Iff(true, "yes", "no"); got != "yes" {
		t.Fatalf("Iff(true) = %q, want %q", got, "yes")
	}
	if got := Iff(false, "yes", "no"); got != "no" {
		t.Fatalf("Iff(false) = %q, want %q", got, "no")
	}
}

func TestReflectValueSliceToString(t *testing.T) {
	values := []reflect.Value{
		reflect.ValueOf("hi"),
		reflect.ValueOf(int64(42)),
	}
	out := ReflectValueSliceToString(values)

	expected := "(string/string(hi),int64/int64(" + values[1].String() + "))"
	if out != expected {
		t.Fatalf("ReflectValueSliceToString() = %q, want %q", out, expected)
	}
	if !strings.HasPrefix(out, "(") || !strings.HasSuffix(out, ")") {
		t.Fatalf("ReflectValueSliceToString() = %q, want wrapping parentheses", out)
	}
}

func TestSyncMap(t *testing.T) {
	m := NewSyncMap[string, int]()
	if m.HasKey("missing") {
		t.Fatalf("HasKey returned true for missing key")
	}
	if _, ok := m.Get("missing"); ok {
		t.Fatalf("Get returned ok for missing key")
	}

	m.Set("a", 1)
	if got, ok := m.Get("a"); !ok || got != 1 {
		t.Fatalf("Get(%q) = (%d, %v), want (1, true)", "a", got, ok)
	}
	if !m.HasKey("a") {
		t.Fatalf("HasKey returned false for existing key")
	}

	keys := m.Keys()
	if len(keys) != 1 || keys[0] != "a" {
		t.Fatalf("Keys() = %v, want [a]", keys)
	}
	values := m.Values()
	if len(values) != 1 || values[0] != 1 {
		t.Fatalf("Values() = %v, want [1]", values)
	}

	m.Delete("a")
	if _, ok := m.Get("a"); ok {
		t.Fatalf("Get returned ok after Delete")
	}

	m.Set("a", 1)
	m.Set("b", 2)
	m.Clear()
	if len(m.Keys()) != 0 {
		t.Fatalf("Keys() not empty after Clear")
	}
	if len(m.Values()) != 0 {
		t.Fatalf("Values() not empty after Clear")
	}
}
