package util

import (
	"fmt"
	"reflect"
)

// Destroyable is implemented by types that require manual cleanup.
type Destroyable interface {
	Destroy()
}

// DestroySlice cleans up each element in the slice.
// It supports both value and pointer receivers.
func DestroySlice[T any](items []T) {
	for i := range items {
		if d, ok := any(items[i]).(Destroyable); ok {
			if !isNilDestroyable(d) {
				d.Destroy()
			}
			continue
		}
		if d, ok := any(&items[i]).(Destroyable); ok {
			if !isNilDestroyable(d) {
				d.Destroy()
			}
			continue
		}
		panic(fmt.Sprintf("DestroySlice: element type %T does not implement Destroy()", items[i]))
	}
}

func isNilDestroyable(d Destroyable) bool {
	value := reflect.ValueOf(d)
	switch value.Kind() {
	case reflect.Pointer, reflect.Func, reflect.Map, reflect.Slice, reflect.Interface, reflect.Chan:
		return value.IsNil()
	default:
		return false
	}
}
