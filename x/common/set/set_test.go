package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var elementSlice = []string{"fire", "earth", "water", "air"}

func TestAdd(t *testing.T) {
	elements := New(elementSlice...)

	assert.False(t, elements.Has("lava"))
	assert.False(t, elements.Has("mud"))

	elements.Add("lava")
	elements.Add("mud")
	assert.True(t, elements.Has("lava"))
	assert.True(t, elements.Has("mud"))

	assert.Equal(t, 6, elements.Len())

	// Add blank string
	elements.Add("")
	assert.True(t, elements.Has(""))
	assert.Equal(t, 7, elements.Len())
}

func TestRemove(t *testing.T) {
	elements := New(elementSlice...)
	elem := "water"
	assert.True(t, elements.Has(elem))

	elements.Remove(elem)
	assert.False(t, elements.Has(elem))
}

func TestHas(t *testing.T) {
	elements := New(elementSlice...)

	assert.True(t, elements.Has("fire"))
	assert.True(t, elements.Has("water"))
	assert.True(t, elements.Has("air"))
	assert.True(t, elements.Has("earth"))
	assert.False(t, elements.Has(""))
	assert.False(t, elements.Has("foo"))
	assert.False(t, elements.Has("bar"))
}

func TestLen(t *testing.T) {
	elements := New(elementSlice...)
	assert.Equal(t, elements.Len(), 4)

	elements.Remove("fire")
	elements.Remove("water")
	assert.Equal(t, elements.Len(), 2)
}

func TestList(t *testing.T) {
	elements := New(elementSlice...)
	assert.Contains(t, elements.List(), "fire")
	assert.Contains(t, elements.List(), "water")
	assert.Contains(t, elements.List(), "air")
	assert.Contains(t, elements.List(), "earth")
}

func TestToMap(t *testing.T) {
	elements := New(elementSlice...)
	elements.Add("lava")

	m := elements.ToMap()
	for _, elem := range elementSlice {
		assert.Contains(t, m, elem)
	}
	assert.Contains(t, m, "lava")
	assert.NotContains(t, m, "mud")
}

func TestIterate(t *testing.T) {
	elements := New(elementSlice...)
	elements.Add("lava")
	elements.Add("mud")

	elements.Iterate(func(elem string) bool {
		assert.True(t, elements.Has(elem))
		return false
	})
}

func TestIterateAll(t *testing.T) {
	elements := New(elementSlice...)
	elements.Add("lava")
	elements.Add("mud")

	elements.IterateAll(func(elem string) {
		assert.True(t, elements.Has(elem))
	})
}

func TestUnion(t *testing.T) {
	elements := New(elementSlice...)
	elements.Add("lava")

	other := New("lava", "mud")

	union := elements.Union(other)
	assert.True(t, union.Has("fire"))
	assert.True(t, union.Has("earth"))
	assert.True(t, union.Has("water"))
	assert.True(t, union.Has("air"))
	assert.True(t, union.Has("lava"))
	assert.True(t, union.Has("mud"))
	assert.Equal(t, union.Len(), 6)
}
