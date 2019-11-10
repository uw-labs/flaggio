package flaggio

// Variant represents a value that can be returned by the evaluation
// process of a Flag.
type Variant struct {
	ID          string
	Description *string
	Value       interface{}
}
