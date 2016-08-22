package parser2

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	tokenItems, err := Tokenizer("A B")
	assert.NoError(t, err)
	assert.Equal(t, tokenItems.baseQT, SHOULD)
	assert.Equal(t, len(tokenItems.items), 3)
	assert.Equal(t, tokenItems.items[0], &TokenItem{t: _RAW, value: "A"})
	assert.Equal(t, tokenItems.items[1], &TokenItem{t: _EMPTYSPACE})
	assert.Equal(t, tokenItems.items[2], &TokenItem{t: _RAW, value: "B"})

	tokenItems, err = Tokenizer("(A || B) +price:[10~20} -title:\"golang\"")
	assert.NoError(t, err)
	assert.Equal(t, tokenItems.baseQT, SHOULD)
	assert.Equal(t, len(tokenItems.items), 19)
	assert.Equal(t, tokenItems.items[0], &TokenItem{t: _OPEN_PAREN})
	assert.Equal(t, tokenItems.items[1], &TokenItem{t: _RAW, value: "A"})
	assert.Equal(t, tokenItems.items[2], &TokenItem{t: _EMPTYSPACE})
	assert.Equal(t, tokenItems.items[3], &TokenItem{t: _OR})
	assert.Equal(t, tokenItems.items[4], &TokenItem{t: _EMPTYSPACE})
	assert.Equal(t, tokenItems.items[5], &TokenItem{t: _RAW, value: "B"})
	assert.Equal(t, tokenItems.items[6], &TokenItem{t: _CLOSE_PAREN})
	assert.Equal(t, tokenItems.items[7], &TokenItem{t: _EMPTYSPACE})
	assert.Equal(t, tokenItems.items[8], &TokenItem{t: _PLUS})
	assert.Equal(t, tokenItems.items[9], &TokenItem{t: _RAW, value: "price"})
	assert.Equal(t, tokenItems.items[10], &TokenItem{t: _COLON})
	assert.Equal(t, tokenItems.items[11], &TokenItem{t: _OPEN_BRACK})
	assert.Equal(t, tokenItems.items[12], &TokenItem{t: _RAW, value: "10~20"})
	assert.Equal(t, tokenItems.items[13], &TokenItem{t: _CLOSE_BRACE})
	assert.Equal(t, tokenItems.items[14], &TokenItem{t: _EMPTYSPACE})
	assert.Equal(t, tokenItems.items[15], &TokenItem{t: _SUB})
	assert.Equal(t, tokenItems.items[16], &TokenItem{t: _RAW, value: "title"})
	assert.Equal(t, tokenItems.items[17], &TokenItem{t: _COLON})
	assert.Equal(t, tokenItems.items[18], &TokenItem{t: _TEXT, value: "golang"})
}

func TestParse(t *testing.T) {
	text := "(A || B) +price:[10~20} -title:\"golang\""
	log.Println("text:", text)

	tokenItems, err := Tokenizer(text)
	assert.NoError(t, err)

	p, err := Parse(tokenItems)
	assert.NoError(t, err)

	p.Parse()
}

func TestParse2(t *testing.T) {
	text := "(A || B) -C +title:(D E) F "
	log.Println(text)

	tis, err := Tokenizer(text)
	assert.NoError(t, err)

	p, err := Parse(tis)
	assert.NoError(t, err)
	group, err := p.Parse()
	assert.NoError(t, err)
	printGroup(group)
}

func printGroup(groups []Group) {
	for i := 0; i < len(groups); i++ {
		log.Println("group:", i)
		for k, item := range groups[i].Items {
			log.Printf("%d : %#v", k, item)
		}
	}
}
