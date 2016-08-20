package parser

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	pairs, err := Tokenizer("(A || B) && C \\\"ABC")
	assert.NoError(t, err)

	assert.Equal(t, pairs.items[0], &TokenItem{t: _OPEN_PAREN})
	assert.Equal(t, pairs.items[1], &TokenItem{t: _RAW, value: "A"})
	assert.Equal(t, pairs.items[2], &TokenItem{t: _EMPTYSPACE, value: " "})
	assert.Equal(t, pairs.items[3], &TokenItem{t: _OR})
	assert.Equal(t, pairs.items[4], &TokenItem{t: _EMPTYSPACE, value: " "})
	assert.Equal(t, pairs.items[5], &TokenItem{t: _RAW, value: "B"})
	assert.Equal(t, pairs.items[6], &TokenItem{t: _CLOSE_PAREN})
	assert.Equal(t, pairs.items[7], &TokenItem{t: _EMPTYSPACE, value: " "})
	assert.Equal(t, pairs.items[8], &TokenItem{t: _AND})
	assert.Equal(t, pairs.items[9], &TokenItem{t: _EMPTYSPACE, value: " "})
	assert.Equal(t, pairs.items[10], &TokenItem{t: _RAW, value: "C"})
	assert.Equal(t, pairs.items[11], &TokenItem{t: _EMPTYSPACE, value: " "})
	assert.Equal(t, pairs.items[12], &TokenItem{t: _RAW, value: "\\\"ABC"})
}

func TestParse(t *testing.T) {
	text := "(A || B || C) (D && E) F"
	log.Println("text:", text)
	tokenItems, err := Tokenizer(text)
	assert.NoError(t, err)

	p, err := Parse(tokenItems)
	assert.NoError(t, err)

	// assert.Equal(t, recur_count_or(p), 6)

	groups := make(Groups, recur_count_or(p))

	for i := 0; i < len(groups); i++ {
		groups[i] = &Group{}
	}
	groups = p.Parse()

	log.Println("len(groups):", len(groups))

	for i := 0; i < len(groups); i++ {
		log.Println("group:", i)
		for k, item := range groups[i].items {
			log.Printf("%d : %#v", k, item)
		}
	}
}
