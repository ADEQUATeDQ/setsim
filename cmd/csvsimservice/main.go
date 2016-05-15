// Copyright (c) Johann Höchtl 2016
//
// See LICENSE for License

// Beginnings of a test to compare the similarity between two CSV-headers
package main

import "github.com/the42/setsim"

func main() {
	var doc1 = []string{
		"Date", "Time", "ZIP", "Address", "Reason",
	}
	var doc2 = []string{
		"Time", "Date", "ZIP", "Address", "Reason",
	}
	var doc3 = []string{
		"Date1", "Time", "ZIP", "Address", "Reason",
	}

	println(setsim.StringDistance(doc1, doc2))
	println(setsim.StringDistance(doc1, doc3))
}
