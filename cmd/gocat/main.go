// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Originally authored by Russ Cox (https://code.google.com/p/rsc/source/browse/cmd/bundle/main.go)
//
// Modified by Joseph Naegele, 2015

// Gocat combines multiple Go source files into a single source file,
// optionally adding a prefix to all top-level names.
//
// Usage:
//	gocat [-p pkgname] [-x prefix] file_0 [file_n...] >file.go
//
// Example
//
//	gocat -p acl github.com/naegelejd/go-acl/acl* > acl.go
//
// Bugs
//
// Gocat has many limitations, most of them not fundamental.
//
// It does not work with cgo.
//
// It does not work with renamed imports.
//
// It does not correctly translate struct literals when prefixing is enabled
// and a field key in the literal key is the same as a top-level name.
//
// It does not work with embedded struct fields.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/itsmontoya/gocat"
	"io"
)

var (
	pkgname    = flag.String("p", "", "package name to use in output file")
	prefix     = flag.String("x", "", "prefix to add to all top-level names")
	nocomments = flag.Bool("c", false, "ignore comments")
	notest     = flag.Bool("n", false, "ignore test files")
	kill       = flag.Bool("k", false, "delete concatenated files from disk")
	out        = flag.String("o", "", "output file (if empty, stdout will be used")
)

func die(v ...interface{}) {
	fmt.Fprint(os.Stderr, v)
	os.Exit(1)
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}

	var outW io.Writer
	if *out == "" {
		outW = os.Stdout
	} else {
		f, err := os.Create(*out)
		if err != nil {
			die(err)
		}

		defer f.Close()
		outW = f
	}

	if err := gocat.Cat(outW, *pkgname, *prefix, flag.Args(), *nocomments, *notest, *kill); err != nil {
		die(err)
	}
}
