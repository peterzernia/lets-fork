package ptr

// Bool turns a bool into a *bool
func Bool(b bool) *bool {
	p := b
	return &p
}

// Float64 turns a float64 into a *float64
func Float64(f float64) *float64 {
	p := f
	return &p
}

// Int64 turns an int64 into a *int64
func Int64(i int64) *int64 {
	p := i
	return &p
}

// String turns a string into *string
func String(s string) *string {
	p := s
	return &p
}
