Setsim is a library to calculate the Similarity of two ordered Lists.

**Installation** `go get github.com/the42/setsim`  
**Build status**  [![Build Status](https://travis-ci.org/the42/setsim.svg?branch=master)](https://travis-ci.org/the42/setsim)  
**Documentation** [![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/the42/setsim)

Copyright (c) 2016 Johann Höchtl. For License see LICENSE

This libray uses a combination of

* [Spearmans Footrule Distance](http://perso.telecom-paristech.fr/~bloch/P6/IREC/Ranking/77_04_spearmans.pdf)
* [Levenshtein Distance](https://en.wikipedia.org/wiki/Levenshtein_distance)

The items of the lists have to be bytes and, for comparability as a string, interpretable as UTF8-Runes.

For further reading see http://theory.stanford.edu/~sergei/slides/www10-metrics.pdf
