package main

import "github.com/the42/setsim"

func main() {
	var doc1 = [][]byte{
		[]byte("Zeit"), []byte("Datum"), []byte("Temperatur"), /* []byte("Temperatur2",*/
	}
	var doc2 = [][]byte{
		[]byte("Datum"), []byte("Zeit"), []byte("Temperatur"),
	}

	println(setsim.Distance(doc1, doc2))
}
