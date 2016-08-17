package parser

import (
	"errors"
	"log"
	"regexp"
)

type Type int8

const (
	_OPEN_PAREN  Type = iota // (
	_CLOSE_PAREN             // )
	_OPEN_BRACK              // [
	_CLOSE_BRACK             // ]
	_OR                      // ||
	_AND                     // &&
	_TEXT                    // "
	_NUMBER                  // 0,1,2
	_EMPTYSPACE
	_RAW
)

var (
	re_empty_space = regexp.MustCompile(`\s`)
	re_number      = regexp.MustCompile(`\d`)
)

var (
	ErrInval = errors.New("invalid")
)

type Token struct {
	query   []rune
	scope   int
	current int
}

func (t *Token) next() rune {
	t.current++
	return t.query[t.current]
}

func (t *Token) hasNext() bool {
	if tmp := t.current + 1; tmp < len(t.query) {
		return true
	}
	return false
}

func (t *Token) peek(n int) (rune, error) {
	pos := t.current + n
	if pos < 0 || pos >= len(t.query) {
		return 0, ErrInval
	}

	return t.query[pos], nil
}

type Pair struct {
	t     Type
	scope int
	value string
}

func Tokenizer(query string) *Pairs {
	tokens := &Token{
		query:   []rune(query),
		current: -1,
	}

	pairs := []*Pair{}
	forSort := []*Pair{}

	var char rune
	for tokens.hasNext() {
		char = tokens.next()
		log.Println(string(char))

		value := []rune{}
		if re_empty_space.MatchString(string(char)) {
			value = append(value, char)

			var r rune
			for tokens.hasNext() {
				r, _ = tokens.peek(1)

				if re_empty_space.MatchString(string(r)) {
					value = append(value, r)
					tokens.next()
				} else {
					break
				}
			}
		}
		if len(value) > 0 {
			pairs = append(pairs, &Pair{t: _EMPTYSPACE, scope: tokens.scope, value: string(value)})
			continue
		}

		if re_number.MatchString(string(char)) {
			var r rune
			for tokens.hasNext() {
				r, _ = tokens.peek(1)

				if re_number.MatchString(string(r)) {
					value = append(value, r)

					tokens.next()
				} else {
					break
				}
			}
		}
		if len(value) > 0 {
			pairs = append(pairs, &Pair{t: _NUMBER, scope: tokens.scope, value: string(value)})
			continue
		}

		switch char {
		case '(':
			pairs = append(pairs, &Pair{t: _OPEN_PAREN, scope: tokens.scope})
			tokens.scope++
		case ')':
			tokens.scope--
			pairs = append(pairs, &Pair{t: _CLOSE_PAREN, scope: tokens.scope})
		case '[':
			pairs = append(pairs, &Pair{t: _OPEN_BRACK, scope: tokens.scope})
		case ']':
			pairs = append(pairs, &Pair{t: _CLOSE_BRACK, scope: tokens.scope})
		case '|':
			if r, err := tokens.peek(1); err == nil {
				if r == '|' {
					pair := &Pair{t: _OR, scope: tokens.scope}
					pairs = append(pairs, pair)
					forSort = append(forSort, pair)

					tokens.next()
				} else {
					pairs = append(pairs, &Pair{t: _RAW, scope: tokens.scope, value: string(char)})
				}
			}

		case '&':
			if r, err := tokens.peek(1); err == nil {
				if r == '&' {
					pair := &Pair{t: _AND, scope: tokens.scope}
					pairs = append(pairs, pair)
					forSort = append(forSort, pair)

					tokens.next()
				} else {
					pairs = append(pairs, &Pair{t: _RAW, scope: tokens.scope, value: string(char)})
				}
			}

		case '"':
			var r rune
			for tokens.hasNext() {
				r, _ = tokens.peek(1)

				if r == '"' {
					if prevRune, err := tokens.peek(-1); err == nil && prevRune == '\\' {
						value = append(value, '"')
						tokens.next()
					} else {
						break
					}
				} else {
					value = append(value, r)
					tokens.next()
				}
			}

			pairs = append(pairs, &Pair{t: _TEXT, scope: tokens.scope, value: string(value)})

		default:
			pairs = append(pairs, &Pair{t: _RAW, scope: tokens.scope, value: string(char)})
		}
	}

	return &Pairs{pairs: pairs, forSort: forSort}
}
