package bseg

import (
)

type BSeg struct {
	dict map[string]int
	totalCount int
}

func (s *BSeg) IncrDict(word string){
	dict[word]++
	totalCount++
}

func (s *BSeg) DecrDict(word string){
	dict[word]--
	totalCount--
}

func Sample() {

}
