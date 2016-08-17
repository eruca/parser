package parser

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	pairs := Tokenizer("(A || B) && C")

	assert.Equal(t, pairs.items[0], &TokenItem{t: _OPEN_PAREN})
	assert.Equal(t, pairs.items[1], &TokenItem{t: _RAW, value: "A"})
	assert.Equal(t, pairs.items[2], &TokenItem{t: _OR})
	assert.Equal(t, pairs.items[3], &TokenItem{t: _RAW, value: "B"})
	assert.Equal(t, pairs.items[4], &TokenItem{t: _CLOSE_PAREN})
	assert.Equal(t, pairs.items[5], &TokenItem{t: _AND})
	assert.Equal(t, pairs.items[6], &TokenItem{t: _RAW, value: "C"})
}

func TestParse(t *testing.T) {
	// tokenItems := Tokenizer("A && B")

	// p := Parse(tokenItems)
	// assert.Equal(t, p.String(), "A && B")

	tokenItems := Tokenizer("(A (B || C) D) E")
	p := Parse(tokenItems)

	if ps, ok := p.(Parsers); ok {
		for k, p1 := range ps {
			log.Println(k, p1.Len())
			if k == 0 {
				log.Println(p1.(Parsers)[1].String())
			}
		}
	} else {
		log.Panic("not right")
	}
}
