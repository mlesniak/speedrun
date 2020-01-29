package main

func assert(condition bool) {
	if !condition {
		panic("Assertion failed.")
	}
}
