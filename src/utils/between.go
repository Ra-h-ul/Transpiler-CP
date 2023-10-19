package utils

import "strings"

// Between returns the string Between two given strings, or the original string
// firstA specifies if the first or last instance of a should be used
// firstB specifies if the first or last instance of b should be used
func Between(s, a, b string, lastA, lastB bool) string {
	var aPos int
	if lastA {
		aPos = strings.LastIndex(s, a)
	} else {
		aPos = strings.Index(s, a)
	}
	if aPos == -1 {
		return s
	}
	var bPos int
	if lastB {
		bPos = strings.LastIndex(s, b)
	} else {
		bPos = strings.Index(s, b)
	}
	if bPos == -1 {
		return s
	}
	if bPos < aPos {
		return s[bPos+len(b) : aPos]
	}
	return s[aPos+len(a) : bPos]
}

// leftBetween searches from the left and returns the first string that is
// Between a and b.
func LeftBetween(s, a, b string) string {
	return Between(s, a, b, false, false)
}

// like leftBetween, but use the rightmost instance of b
func LeftBetweenRightmost(s, a, b string) string {
	return Between(s, a, b, false, true)
}

// greedyBetween searches from the left for a then
// searches as far as possible for b, before returning
// the string that is Between a and b.
func GreedyBetween(s, a, b string) string {
	return Between(s, a, b, false, true)
}
