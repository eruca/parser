package parser

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	pairs, _ := Tokenizer("(A || B) && C")

	assert.Equal(t, pairs.items[0], &TokenItem{t: _OPEN_PAREN})
	assert.Equal(t, pairs.items[1], &TokenItem{t: _RAW, value: "A"})
	assert.Equal(t, pairs.items[2], &TokenItem{t: _OR})
	assert.Equal(t, pairs.items[3], &TokenItem{t: _RAW, value: "B"})
	assert.Equal(t, pairs.items[4], &TokenItem{t: _CLOSE_PAREN})
	assert.Equal(t, pairs.items[5], &TokenItem{t: _AND})
	assert.Equal(t, pairs.items[6], &TokenItem{t: _RAW, value: "C"})
}

func TestParse(t *testing.T) {
	text := "A && (B || C) D"
	log.Println("text:", text)
	tokenItems, cntOr := Tokenizer(text)
	p := Parse(tokenItems)
	groups := make(Groups, cntOr+1)

	for i := 0; i < len(groups); i++ {
		groups[i] = &Group{}
	}
	Convert(p, &groups)

	for i := 0; i < len(groups); i++ {
		log.Println("group:", i)
		for k, item := range groups[i].items {
			log.Printf("%d : %#v", k, item)
		}
	}
}
