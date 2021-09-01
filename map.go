// Package ziprox is a library for computing proximity between US zip codes.
package ziprox

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Zip represents a zip code.
//
// Encodes a zip code as 32-bit unsigned integer. A zip code in the United States
// is five numeric digits, which may include leading zeroes.
type Zip uint32

// Parse s as a zip code.
//
// If s is not in the form of a valid zip code, an error is returned.
func Parse(s string) (Zip, error) {
	if len(s) != 5 {
		return 0, fmt.Errorf("%q is not a zip code", s)
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%q is not a zip code", s)
	}

	return Zip(i), nil
}

func (z Zip) String() string {
	return fmt.Sprintf("%05d", z)
}

type buckets struct {
	under5   []Zip
	under10  []Zip
	under20  []Zip
	under50  []Zip
	under100 []Zip
	under200 []Zip
	under500 []Zip
}

// Map contains zip code proximity information.
type Map struct {
	groups map[Zip]*buckets
}

func (m Map) insert(origin, dest Zip, dist float64) {
	m.setup(origin)

	switch {
	case dist < 5:
		m.groups[origin].under5 = append(m.groups[origin].under5, dest)
	case dist < 10:
		m.groups[origin].under10 = append(m.groups[origin].under10, dest)
	case dist < 20:
		m.groups[origin].under20 = append(m.groups[origin].under20, dest)
	case dist < 50:
		m.groups[origin].under50 = append(m.groups[origin].under50, dest)
	case dist < 100:
		m.groups[origin].under100 = append(m.groups[origin].under100, dest)
	case dist < 200:
		m.groups[origin].under200 = append(m.groups[origin].under200, dest)
	case dist < 500:
		m.groups[origin].under500 = append(m.groups[origin].under500, dest)
	}
}

func (m Map) setup(origin Zip) {
	if m.groups[origin] == nil {
		m.groups[origin] = &buckets{
			under5:   make([]Zip, 0, 30),
			under10:  make([]Zip, 0, 30),
			under20:  make([]Zip, 0, 30),
			under50:  make([]Zip, 0, 30),
			under100: make([]Zip, 0, 30),
			under200: make([]Zip, 0, 30),
			under500: make([]Zip, 0, 30),
		}
	}
}

// Within returns the set of zip codes within distance of origin.
func (m Map) Within(origin Zip, distance int) []Zip {
	bucket, exists := m.groups[origin]
	if !exists {
		return nil
	}

	var combine []Zip

	combine = append(combine, bucket.under5...)

	if distance > 5 {
		combine = append(combine, bucket.under10...)
	}

	if distance > 10 {
		combine = append(combine, bucket.under20...)
	}

	if distance > 20 {
		combine = append(combine, bucket.under50...)
	}

	if distance > 50 {
		combine = append(combine, bucket.under100...)
	}

	if distance > 100 {
		combine = append(combine, bucket.under200...)
	}

	if distance > 200 {
		combine = append(combine, bucket.under500...)
	}

	return combine
}

func tokensSub(line string) (a, b, d string) {
	return line[1:6], line[9:14], line[16:22]
}

// New extracts in and inflates a new Map.
func New(in io.Reader) (*Map, error) {
	m := &Map{groups: make(map[Zip]*buckets)}

	scanner := bufio.NewScanner(in)
	scanner.Scan() // skip header line

	for scanner.Scan() {
		line := scanner.Text()

		a, b, d := tokensSub(line)

		source, err := strconv.Atoi(a)
		if err != nil {
			return nil, err
		}

		dest, err := strconv.Atoi(b)
		if err != nil {
			return nil, err
		}

		dist, err := strconv.ParseFloat(d, 64)
		if err != nil {
			return nil, err
		}

		m.insert(Zip(source), Zip(dest), dist)
	}
	return m, scanner.Err()
}
