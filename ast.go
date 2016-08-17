package parser

import (
	"fmt"
	"log"
	"strings"
)

type Parser interface {
	Len() int
	String() string
	Parse() []*QueryItem
}

type Pairs struct {
	pairs []*Pair
	// forSort []*Pair
	current int
}

// func (p *Pairs) Len() int {
// 	return len(p.forSort)
// }

// func (p *Pairs) Less(i, j int) bool {
// 	p.forSort[i].scope < p.forSort[j].scope
// }

// func (p *Pairs) Swap(i, j int) {
// 	p.forSort[i], p.forSort[j] = p.forSort[j], p.forSort[i]
// }

func (p *Pairs) hasNext() bool {
	if tmp := p.current + 1; tmp < len(p.pairs) {
		return true
	}

	return false
}

func (p *Pairs) next() *Pair {
	p.current++
	return p.pairs[p.current]
}

func (p *Pairs) peek(n int) *Pair {
	if tmp := p.current + n; tmp >= 0 && tmp < len(p.pairs) {
		return p.pairs[tmp]
	}

	return nil
}

func Parse(ps []*Pair) Parser {
	pairs := &Pairs{pairs: ps, current: -1}

	result := Parsers{}
	for pairs.hasNext() {
		pair := pairs.next()
		log.Println(pair.value)

		switch pair.t {
		case _OR:
			var left, right *Pair

			peek := -1
			for left = pairs.peek(peek); left != nil; left = pairs.peek(peek) {
				switch left.t {
				case _RAW:
					result = append(result, &Raw{text: left.value})
				}
				peek--
			}

			peek = 1
			for right = pairs.peek(peek); right != nil; right = pairs.peek(peek) {
				switch right.t {
				case _OPEN_PAREN:
				case _CLOSE_PAREN:
				case _RAW:
					result = append(result, &Raw{text: right.value})
				}
				peek++
			}

		case _AND:
		}

	}

	return result
}

type Parsers []Parser

func (p Parsers) Len() int {
	return len(p)
}

func (p Parsers) String() string {
	result := make([]string, len(p))

	for i := 0; i < len(p); i++ {
		result[i] = p[i].String()
	}

	return strings.Join(result, " ")
}

func (p Parsers) Parse() []*QueryItem {
	res := []*QueryItem{}

	for _, parser := range p {
		res = append(res, parser.Parse()...)
	}

	return res
}

// 代表 左右两个都必须有
type AND struct {
	left  Parser
	right Parser
}

func (and *AND) String() string {
	return fmt.Sprintf("%s && %s", and.left, and.right)
}

// 代表 或者
type OR struct {
	left  Parser
	right Parser
}

func (or *OR) Len() int {
	return 1
}

func (or *OR) String() string {
	return fmt.Sprintf("%s || %s", or.left, or.right)
}

func (or *OR) Parse() []*QueryItem {
	result := make([]*QueryItem, 0, 2)
	result = append(result, or.left.Parse()...)
	result = append(result, or.right.Parse()...)

	return result
}

// Text 代表以""包饶的
type Text struct {
	text string
}

func (t *Text) Len() int {
	return 1
}

func (t *Text) String() string {
	return t.text
}

func (t *Text) Parse() []*QueryItem {
	return []*QueryItem{&QueryItem{Text: t.text, Offset: []int{}}}
}

// Raw 代表 毫无修饰的词项
type Raw struct {
	text string
}

func (r *Raw) Len() int {
	return 1
}

func (r *Raw) String() string {
	return r.text
}

func (r *Raw) Parse() []*QueryItem {
	return []*QueryItem{&QueryItem{Text: r.text}}
}
