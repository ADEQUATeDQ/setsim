// Copyright (c) 2016 Johann Höchtl
// See LICENSE for license

/*
Package setsim TODO: Add package documentation
*/
package setsim

import (
	"bytes"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// Distance is a measure for the structural difference of two byte sets
// TODO: expand documentation
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
	for diff := range diffset {
		adiff += sfd * diff
	}
	return adiff
}

// SpearmanFootDistance calculates the Spearman Footrule Distance of two byte arrays
// When comparing two strings, you might consider normalizing these strings eg. http://unicode.org/faq/normalization.html
// For more information see https://mikespivey.wordpress.com/2014/01/20/the-maximum-value-of-spearmans-footrule-distance/
// This function extends the classical definition of Spearman Footrule Distance in the sense that
// if one element of set a is not contained within set b, this max(len a , len b ) to the Foot Distance
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
