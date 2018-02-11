package ui

// Element is the general interface for all elements to be used
type Element interface {
	Update(float32)
	Draw()
}
