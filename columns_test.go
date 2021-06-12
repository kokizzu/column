// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// BenchmarkProperty/set-8         	344183034	         3.498 ns/op	       0 B/op	       0 allocs/op
// BenchmarkProperty/get-8         	1000000000	         1.123 ns/op	       0 B/op	       0 allocs/op
// BenchmarkProperty/replace-8     	291245523	         4.157 ns/op	       0 B/op	       0 allocs/op
func BenchmarkProperty(b *testing.B) {
	b.Run("update", func(b *testing.B) {
		p := newColumnAny()
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			p.Update(5, "hello")
		}
	})

	b.Run("fetch", func(b *testing.B) {
		p := newColumnAny()
		p.Update(5, "hello")
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			p.Value(5)
		}
	})

	b.Run("replace", func(b *testing.B) {
		p := newColumnAny()
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			p.Update(5, "hello")
			p.Delete(5)
		}
	})
}

func TestProperty(t *testing.T) {
	p := newColumnAny().(*columnAny)

	{ // Set the value at index
		p.Update(9, 99.5)
		assert.Equal(t, 10, len(p.data))
	}

	{ // Get the value
		v, ok := p.Value(9)
		assert.Equal(t, 99.5, v)
		assert.True(t, ok)
	}

	{ // Remove the value
		p.Delete(9)
		v, ok := p.Value(9)
		assert.Equal(t, nil, v)
		assert.False(t, ok)
	}

	{ // Set a couple of values, should only take 2 slots
		p.Update(5, "hi")
		p.Update(1000, "roman")
		assert.Equal(t, 1001, len(p.data))

		v1, ok := p.Value(5)
		assert.True(t, ok)
		assert.Equal(t, "hi", v1)

		v2, ok := p.Value(1000)
		assert.True(t, ok)
		assert.Equal(t, "roman", v2)
	}

}

func TestPropertyOrder(t *testing.T) {

	// TODO: not sure if it's all correct, what happens if
	// we have 2 properties?

	p := newColumnAny()
	for i := uint32(100); i < 200; i++ {
		p.Update(i, i)
	}

	for i := uint32(100); i < 200; i++ {
		x, ok := p.Value(i)
		assert.True(t, ok)
		assert.Equal(t, i, x)
	}

	for i := uint32(150); i < 180; i++ {
		p.Delete(i)
		p.Update(i, i)
	}

	for i := uint32(100); i < 200; i++ {
		x, ok := p.Value(i)
		assert.True(t, ok)
		assert.Equal(t, i, x)
	}
}