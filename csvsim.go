package main

import (
	"bytes"
	"fmt"

	"github.com/mfonda/simhash"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

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

func main() {

	var doc1 = [][]byte{
		[]byte("Zeit1"), []byte("Datum"), []byte("Temperatur"), /* []byte("Temperatur2",*/
	}

	var doc2 = [][]byte{
		[]byte("Datum"), []byte("Zeit"), []byte("Temperatur"),
	}

	hashes1 := make([]uint64, len(doc1))
	for i, d := range doc1 {
		hashes1[i] = simhash.Simhash(simhash.NewWordFeatureSet(d))
		fmt.Printf("Simhash of %s: %x\n", d, hashes1[i])
	}
	hashes2 := make([]uint64, len(doc2))
	for i, d := range doc2 {
		hashes2[i] = simhash.Simhash(simhash.NewWordFeatureSet(d))
		fmt.Printf("Simhash of %s: %x\n", d, hashes2[i])
	}

	maxlen := max(len(hashes1), len(hashes2))
	minlen := min(len(hashes1), len(hashes2))

	diffhash := make([]int, maxlen)
	var i int
	for i = 0; i < maxlen; i++ {
		if i < minlen && i < maxlen {
			diffhash[i] = levenshtein.DistanceForStrings(bytes.Runes(doc1[i]), bytes.Runes(doc2[i]), levenshtein.DefaultOptions)
		} else {
			if len(hashes1) > i {
				diffhash[i] = len(doc1[i])
			} else {
				diffhash[i] = len(doc2[i])
			}

		}
	}

	sfd := SpearmanFootDistance(doc1, doc2)
	fmt.Println(sfd)
	var adiff int
	for _, val := range diffhash {
		adiff += sfd * val
		fmt.Println(val)
	}
	println(adiff)

	//fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[1], simhash.Compare(hashes[0], hashes[1]))
	//fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[2], simhash.Compare(hashes[0], hashes[2]))
}
