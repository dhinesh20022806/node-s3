package filenamegenerator

import (
	"fmt"
	"sync/atomic"
)

type FileNameGenerator struct {
	counter atomic.Uint64
}

func NewFileNameGenerator(count uint64) *FileNameGenerator {
	counter := &FileNameGenerator{}
	counter.counter.Store(count)
	return counter
}

func (g *FileNameGenerator) Next() string {
	n := g.counter.Add(1) - 1

	letterIndex := n % 26
	group := n / 26

	letter := byte('a' + letterIndex)

	if group == 0 {
		return string(letter)
	}

	return fmt.Sprintf("%s-%d", string(letter), group)
}
