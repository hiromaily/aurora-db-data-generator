package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

type Iterator interface {
	Next() (string, bool)
}

//
// Collection Iterator
//

type collectionIterator struct { //nolint: unused
	current int
	data    []string
}

func newCollectionIterator(data []string) Iterator { //nolint: unused
	return &collectionIterator{
		current: 0,
		data:    data,
	}
}

func (c *collectionIterator) Next() (string, bool) { //nolint: unused
	if c.current >= len(c.data) {
		return "", false
	}
	defer func() { c.current++ }()
	return c.data[c.current], true
}

//
// Randomly select one from the list Iterator
//

type randomListIterator struct { //nolint: unused
	data []string
}

func newRandomListIterator(data []string) (Iterator, error) { //nolint: unused
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	return &randomListIterator{
		data: data,
	}, nil
}

func (c *randomListIterator) Next() (string, bool) { //nolint: unused
	if len(c.data) == 0 {
		return "", false
	}

	idx := rand.IntN(len(c.data))
	return c.data[idx], true
}

//
// Custom string Iterator
//

type customStringIterator struct {
	format  string
	current int
}

func newCustomStringIterator(format string) Iterator {
	// format must be like "something%06d"
	return &customStringIterator{
		format:  format,
		current: 0,
	}
}

func (c *customStringIterator) Next() (string, bool) {
	c.current++
	return fmt.Sprintf(c.format, c.current), true
}
