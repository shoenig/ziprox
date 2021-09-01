package ziprox

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const each = `
"10000","10001",1.0000
"10000","10005",5.0000
"10000","10010",10.0000
"10000","10020",20.0000
"10000","10050",50.0000
"10000","10100",100.0000
"10000","10200",200.0000
"10000","10500",500.0000
"10000","10999",999.9000
`

const data = `
"10000","10001",1.0000
"10000","10001",2.0000
"10000","10006",6.0000
"10000","10007",7.0000
`

func TestMap(t *testing.T) {
	r := strings.NewReader(data)
	m, err := New(r)
	require.NoError(t, err)
	_ = m
}

func TestMap_insert_each(t *testing.T) {
	r := strings.NewReader(each)
	m, err := New(r)
	require.NoError(t, err)

	const origin = 10000
	require.Equal(t, []Zip{10001}, m.groups[origin].under5)
	require.Equal(t, []Zip{10005}, m.groups[origin].under10)
	require.Equal(t, []Zip{10010}, m.groups[origin].under20)
	require.Equal(t, []Zip{10020}, m.groups[origin].under50)
	require.Equal(t, []Zip{10050}, m.groups[origin].under100)
	require.Equal(t, []Zip{10100}, m.groups[origin].under200)
	require.Equal(t, []Zip{10200}, m.groups[origin].under500)
}

func TestMap_Within_absent(t *testing.T) {
	r := strings.NewReader(each)
	m, err := New(r)
	require.NoError(t, err)

	result := m.Within(99999, 100)
	require.Empty(t, result)
}

func TestMap_Within_each(t *testing.T) {
	r := strings.NewReader(each)
	m, err := New(r)
	require.NoError(t, err)

	try := func(dist int, exp []Zip) {
		result := m.Within(10000, dist)
		require.Equal(t, exp, result)
	}

	exp := []Zip{10001}
	try(1, exp)

	exp = append(exp, 10005)
	try(6, exp)

	exp = append(exp, 10010)
	try(11, exp)

	exp = append(exp, 10020)
	try(21, exp)

	exp = append(exp, 10050)
	try(51, exp)

	exp = append(exp, 10100)
	try(101, exp)

	exp = append(exp, 10200)
	try(201, exp)
}

func Test_Parse(t *testing.T) {
	try := func(in string, expZip Zip, expErr error) {
		z, err := Parse(in)
		require.Equal(t, expZip, z)
		require.Equal(t, expErr, err)
	}

	try("00000", 0, nil)
	try("78701", 78701, nil)
	try("00631", 631, nil)
	try("abcde", 0, errors.New(`"abcde" is not a zip code`))
	try("1234", 0, errors.New(`"1234" is not a zip code`))
}

func TestZip_String(t *testing.T) {
	try := func(in Zip, exp string) {
		result := in.String()
		require.Equal(t, exp, result)
	}

	try(Zip(12345), "12345")
	try(Zip(78701), "78701")
	try(Zip(321), "00321")
}
