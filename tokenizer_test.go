package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	pairs := Tokenizer("(A || B) && C")

	assert.Equal(t, pairs[0], &Pair{t: _OPEN_PAREN})
	assert.Equal(t, pairs[1], &Pair{t: _RAW, scope: 1, value: "A"})
	assert.Equal(t, pairs[2], &Pair{t: _EMPTYSPACE, scope: 1, value: " "})
	assert.Equal(t, pairs[3], &Pair{t: _OR, scope: 1})
	assert.Equal(t, pairs[4], &Pair{t: _EMPTYSPACE, scope: 1, value: " "})
	assert.Equal(t, pairs[5], &Pair{t: _RAW, scope: 1, value: "B"})
	assert.Equal(t, pairs[6], &Pair{t: _CLOSE_PAREN})
	assert.Equal(t, pairs[8], &Pair{t: _AND})
	assert.Equal(t, pairs[10], &Pair{t: _RAW, value: "C"})
}
