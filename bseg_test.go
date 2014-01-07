package bseg

import (
	"log"
	"testing"
)

func TestBSeg(t *testing.T) {
	s := NewBSeg()

	tokens := []string{"hello", "world"}
	segments := []uint8{SEG}

	s.ProcessText(tokens, segments)

	log.Print(s.dict)
	log.Print(s.unigram)
}
