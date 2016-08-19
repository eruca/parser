package parser

import (
	"errors"
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

func (t Type) String() string {
	switch t {
	case _OPEN_PAREN:
		return "("
	case _CLOSE_PAREN:
		return ")"
	case _OPEN_BRACK:
		return "["
	case _CLOSE_BRACK:
		return "]"
	}

	return ""
}

var (
	re_empty_space = regexp.MustCompile(`\s`)
	re_number      = regexp.MustCompile(`\d`)
)

var (
	ErrInval = errors.New("invalid")
)

type Token struct {
	query   []rune
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

type TokenItem struct {
	t     Type
	value string
	index int
}

type TokenItems struct {
	items   []*TokenItem
	current int
}

func NewTokenItems(items []*TokenItem) *TokenItems {
	return &TokenItems{items: items, current: -1}
}

func (tis *TokenItems) hasNext() bool {
	if tmp := tis.current + 1; tmp < len(tis.items) {
		return true
	}
	return false
}

func (tis *TokenItems) next() *TokenItem {
	tis.current++
	return tis.items[tis.current]
}

func (tis *TokenItems) peek(n int) (*TokenItem, error) {
	if tmp := tis.current + n; tmp >= 0 && tmp < len(tis.items) {
		return tis.items[tmp], nil
	}

	return nil, ErrInval
}

func (tis TokenItems) TopAndOr() (item *TokenItem, pos int) {
	pos = -1
	paren := 0

	for tis.hasNext() {
		item = tis.next()

		switch item.t {
		case _OPEN_PAREN:
			paren++
		case _CLOSE_PAREN:
			paren--
		case _AND:
			if paren == 0 {
				pos = tis.current
			}
		case _OR:
			if paren == 0 {
				pos = tis.current
			}
		}
	}

	if pos == -1 {
		return nil, -1
	}

	return tis.items[pos], pos
}

func Tokenizer(query string) (*TokenItems, int) {
	tokens := &Token{
		query:   []rune(query),
		current: -1,
	}

	items := []*TokenItem{}
	cntOr := 0

	var char rune
	for tokens.hasNext() {
		char = tokens.next()

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
			// items = append(items, &TokenItem{t: _EMPTYSPACE, value: string(value)})
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
			items = append(items, &TokenItem{t: _NUMBER, value: string(value)})
			continue
		}

		switch char {
		case '(':
			items = append(items, &TokenItem{t: _OPEN_PAREN})
		case ')':
			items = append(items, &TokenItem{t: _CLOSE_PAREN})
		case '[':
			items = append(items, &TokenItem{t: _OPEN_BRACK})
		case ']':
			items = append(items, &TokenItem{t: _CLOSE_BRACK})
		case '|':
			if r, err := tokens.peek(1); err == nil {
				if r == '|' {
					items = append(items, &TokenItem{t: _OR, index: cntOr})
					cntOr++

					tokens.next()
				} else {
					items = append(items, &TokenItem{t: _RAW, value: string(char)})
				}
			}

		case '&':
			if r, err := tokens.peek(1); err == nil {
				if r == '&' {
					items = append(items, &TokenItem{t: _AND})
					tokens.next()
				} else {
					items = append(items, &TokenItem{t: _RAW, value: string(char)})
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

			items = append(items, &TokenItem{t: _TEXT, value: string(value)})

		default:
			items = append(items, &TokenItem{t: _RAW, value: string(char)})
		}
	}

	return NewTokenItems(items), cntOr
}
