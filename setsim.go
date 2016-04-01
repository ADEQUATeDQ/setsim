// Copyright (c) 2016 Johann Höchtl
// See LICENSE for license

/*
Package setsim provides fuctionality to calculate the distance between two ordered
lists. The items within the list are interpreted as bytes.
*/
package setsim

import (
	"bytes"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// Distance is a measure for the structural difference of two byte lists, which
// can be interpreted as strings. It first calculates the difference between two
// set items using Levenshtein distance. This distance measures serves as a weight to multiply with
// Spearman Foot Distance.
//
// See also the documentation on SpearmanFootDistance concerning the normalization of strings.
//
// For more information see http://theory.stanford.edu/~sergei/slides/www10-metrics.pdf
func Distance(a, b [][]byte) int {
	maxlen := max(len(a), len(b))
	minlen := min(len(a), len(b))

	diffset := make([]int, maxlen)
	for i := 0; i < maxlen; i++ {
		if i < minlen {
			diffset[i] = levenshtein.DistanceForStrings(bytes.Runes(a[i]), bytes.Runes(b[i]), levenshtein.DefaultOptions)
		} else {
			if len(a) > i {
				diffset[i] = len(a[i])
			} else {
				diffset[i] = len(b[i])
			}
		}
	}

	sfd := SpearmanFootDistance(a, b)

	var adiff int
	for _, diff := range diffset {
		adiff += sfd * diff
	}
	return adiff
}

// StringDistance is a helper function to Distance which accepts an array of strings. See the documentation on
// SpearmanFootDistance concerning the normalization of strings.
func StringDistance(a, b []string) int {
	bytesa := make([][]byte, len(a))
	bytesb := make([][]byte, len(b))

	for i := 0; i < len(a); i++ {
		bytesa[i] = []byte(a[i])
	}
	for i := 0; i < len(b); i++ {
		bytesb[i] = []byte(b[i])
	}

	return Distance(bytesa, bytesb)
}

// SpearmanFootDistance calculates the Spearman Footrule Distance of two byte arrays
// For more information see https://mikespivey.wordpress.com/2014/01/20/the-maximum-value-of-spearmans-footrule-distance/
// When comparing two strings, you might consider normalizing these strings eg. http://unicode.org/faq/normalization.html

// This function extends the classical definition of Spearman Footrule Distance in the sense that
// if one element of set a or set b is not contained within the other set, max(len a , len b ) is added to the Foot Distance.
func SpearmanFootDistance(a, b [][]byte) int {
	var odistance int
	for i1, aval := range a {
		for i2, bval := range b {
			if bytes.Equal(aval, bval) {
				diff := i1 - i2
				if diff < 0 {
					diff = -diff
				}
				odistance += diff
				goto next
			}
		}
		odistance += max(len(a), len(b))
	next:
	}
	return odistance
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
