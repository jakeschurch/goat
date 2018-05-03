// Copyright (c) 2017 Jake Schurch
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package instruments

import (
	"strconv"
)

// ----------------------------------------------------------------------------

// Amount is the representation of a Volume and an Amount.
type Amount int

func (amt Amount) ToPercent() string {
	var str, end []byte
	str = []byte(strconv.Itoa(int(amt)))
	str, end = str[0:len(str)-2], str[len(str)-2:]

	return toString(str) + "." + string(end) + "%"
}

// Volume represents a quantity or volume.
type Volume int

// Price is an integer representation of a decimal price amount.
type Price int

// NewVolume instantiates a volume struct from a float.
func NewVolume(f float64) Volume {
	return Volume(f)
}

// String returns a string representation of volume.
func (v Volume) String() string {
	var amt []byte = []byte(strconv.Itoa(int(v)))
	// amt, end = amt[0:len(amt)-2], amt[len(amt)-2:]

	return toString(amt) + ".00"
}

// ----------------------------------------------------------------------------

// NewPrice instantiates a price struct from a float.
func NewPrice(f float64) Price {
	return Price(f * 100)
}

// String returns a string representation of a price value.
func (p Price) String() string {
	var amt, end []byte
	amt = []byte(strconv.Itoa(int(p)))
	amt, end = amt[0:len(amt)-2], amt[len(amt)-2:]

	return "$" + toString(amt) + "." + string(end)
}

// ----------------------------------------------------------------------------

// Divide returns the product of two price values.
func Divide(top, bottom Price) Amount {
	return Amount((top*200 + bottom) / (bottom * 2))
}

// toString takes a byte slice and converts it to a string representation.
func toString(amt []byte) string {
	var result []byte
	nextComma := 3
	for i := 0; len(amt)-1 >= 0; i++ {
		if i == nextComma {
			result = append([]byte{','}, result...)
			nextComma += 3
		}
		result = append(amt[len(amt)-1:], result...)
		amt = amt[0 : len(amt)-1]
	}
	return string(result)
}
