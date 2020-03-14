package app

// Error customize error interface
type Error interface {
	error
	Status() int
}
// StatusError include error code
type StatusError struct {
	Code int
	Err error
}

// Error implement StatusError interface
func (s StatusError) Error() string {
	return s.Err.Error()
}

// Status implement StatusError interface
func (s StatusError) Status() int {
	return s.Code
}