package ptr

// Int64 turns an int64 into a *int64
func Int64(i int64) *int64 {
	p := i
	return &p
}

// Bool turns a bool into a *bool
func Bool(i bool) *bool {
	p := i
	return &p
}

// String turns a string into *string
func String(s string) *string {
	p := s
	return &p
}
