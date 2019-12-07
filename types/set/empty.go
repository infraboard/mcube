package set

// Empty is public since it is used by some internal API objects for conversions between external
// string arrays and internal sets, and conversion logic requires public types today.
type Empty struct{}
