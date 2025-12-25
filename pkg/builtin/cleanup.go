package builtin

// WithString creates a temporary String and destroys it after fn returns.
func WithString(content string, fn func(String)) {
	s := NewStringWithUtf8Chars(content)
	defer s.Destroy()
	fn(s)
}

// WithStringName creates a temporary StringName and destroys it after fn returns.
func WithStringName(content string, fn func(StringName)) {
	sn := NewStringNameWithUtf8Chars(content)
	defer sn.Destroy()
	fn(sn)
}
