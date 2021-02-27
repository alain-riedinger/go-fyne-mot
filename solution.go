package main

// Solution contains the information needed for the recursive solving
type Solution struct {
	BestLen int
	Best    []string
	Current string
}

// NewSolution initializes a Solution structure
func NewSolution() *Solution {
	s := new(Solution)
	s.BestLen = 0
	return s
}
