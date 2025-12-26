package pool

import (
	"sync"

	. "github.com/godot-go/godot-go/pkg/builtin"
)

var (
	variantPool = sync.Pool{New: func() any {
		v := NewVariantNil()
		return &v
	}}
	stringPool     = sync.Pool{New: func() any { return &String{} }}
	stringNamePool = sync.Pool{New: func() any { return &StringName{} }}
	arrayPool      = sync.Pool{New: func() any {
		arr := NewArray()
		return &arr
	}}
	dictPool = sync.Pool{New: func() any {
		dict := NewDictionary()
		return &dict
	}}
)

// AcquireVariant returns a reusable Variant initialized as nil.
func AcquireVariant() *Variant {
	return variantPool.Get().(*Variant)
}

// ReleaseVariant resets the Variant and returns it to the pool.
func ReleaseVariant(v *Variant) {
	if v == nil {
		return
	}
	v.Clear()
	variantPool.Put(v)
}

// AcquireString returns a reusable String initialized as empty.
func AcquireString() *String {
	s := stringPool.Get().(*String)
	*s = NewString()
	return s
}

// ReleaseString destroys the String contents and returns it to the pool.
func ReleaseString(s *String) {
	if s == nil {
		return
	}
	s.Destroy()
	*s = String{}
	stringPool.Put(s)
}

// AcquireStringName returns a reusable StringName initialized as empty.
func AcquireStringName() *StringName {
	sn := stringNamePool.Get().(*StringName)
	*sn = NewStringName()
	return sn
}

// ReleaseStringName destroys the StringName contents and returns it to the pool.
func ReleaseStringName(sn *StringName) {
	if sn == nil {
		return
	}
	sn.Destroy()
	*sn = StringName{}
	stringNamePool.Put(sn)
}

// AcquireArray returns a reusable Array (cleared between uses).
func AcquireArray() *Array {
	return arrayPool.Get().(*Array)
}

// ReleaseArray clears the Array contents and returns it to the pool.
func ReleaseArray(arr *Array) {
	if arr == nil {
		return
	}
	arr.Clear()
	arrayPool.Put(arr)
}

// AcquireDictionary returns a reusable Dictionary (cleared between uses).
func AcquireDictionary() *Dictionary {
	return dictPool.Get().(*Dictionary)
}

// ReleaseDictionary clears the Dictionary contents and returns it to the pool.
func ReleaseDictionary(dict *Dictionary) {
	if dict == nil {
		return
	}
	dict.Clear()
	dictPool.Put(dict)
}
