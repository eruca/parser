package parser

import (
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	pairs := Tokenizer("A || B")
	parser := Parse(pairs)

	// assert.Equal(t, parser.Len(), 1)
	log.Println(parser.String())
}
