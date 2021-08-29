package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"go.etcd.io/bbolt"
	"gophers.dev/pkgs/ignore"
	"gophers.dev/pkgs/regexplus"
)

func usage(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	_, _ = fmt.Fprintf(os.Stderr, "usage: csv2db <input:csv> <output:db>\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 3 {
		usage(nil)
	}

	inFile := os.Args[1]
	outFile := os.Args[2]

	fmt.Println("output:", outFile)
	fmt.Println("input:", inFile)

	bb, err := bbolt.Open(outFile, 0660, nil)
	if err != nil {
		usage(err)
	}

	f, err := os.Open(inFile)
	if err != nil {
		usage(err)
	}

	if err := convert(f, bb); err != nil {
		usage(err)
	}
}

var (
	lineRe = regexp.MustCompile(`^"(?P<source>[\d]{5})","(?P<dest>[\d]{5})",(?P<dist>[\d.]+)$`)
)

func convert(in io.Reader, bb *bbolt.DB) error {
	scanner := bufio.NewScanner(in)
	count := 0
	scanner.Scan() // skip header line

	tx, err := bb.Begin(true)
	if err != nil {
		return err
	}

	for scanner.Scan() {
		line := scanner.Text()
		count++

		m := regexplus.FindNamedSubmatches(lineRe, line)
		if len(m) != 3 {
			fmt.Println("skip line:", line)
			continue
		}
		fmt.Println("m:", m)

		source, err := strconv.Atoi(m["source"])
		if err != nil {
			return err
		}

		dest, err := strconv.Atoi(m["dest"])
		if err != nil {
			return err
		}

		dist, err := strconv.ParseFloat(m["dist"], 64)
		if err != nil {
			return err
		}

		fmt.Printf("%d -> %d = %f\n", source, dest, dist)

		b, err := tx.CreateBucketIfNotExists(encode32(source))
		if err != nil {
			return err
		}

		k := uint64(dist * 1000000000000000)
		v := encode32(dest)
		fmt.Println("k:", k, "v:", v)

		_ = b

		if count > 10 {
			break
		}

	}
	return scanner.Err()
}

func getBucket(bb *bbolt.DB, zip int) *bbolt.Bucket {
	if err := bb.Update(func(tx *bbolt.Tx) error {

	}); err != nil {
		ignore.Panic(err)
	}
}

func encode32(i int) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	return b
}
