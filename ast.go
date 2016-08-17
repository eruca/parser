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

func Parse(tokenItems *TokenItems) Parser {
	if tokenItems == nil {
		return nil
	}

	item, pos := tokenItems.TopAndOr()

	log.Println("pos", pos)
	if pos == -1 {
		return Simple(tokenItems)
	}

	if item.t == _AND {
		return &AND{
			left:  Parse(NewTokenItems(tokenItems.items[:pos])),
			right: Parse(NewTokenItems(tokenItems.items[pos+1:])),
		}
	} else {
		log.Println("into _OR")
		return &OR{
			left:  Parse(NewTokenItems(tokenItems.items[:pos])),
			right: Parse(NewTokenItems(tokenItems.items[pos+1:])),
		}
	}

	return nil
}

// 没有or and
func Simple(ts *TokenItems) Parser {
	ret := Parsers{}

	start, end, parens := 0, 0, 0
	for ts.hasNext() {
		item := ts.next()

		switch item.t {
		case _OPEN_PAREN:
			start = ts.current + 1

		INNER:
			for ts.hasNext() {
				next := ts.next()

				switch next.t {
				case _OPEN_PAREN:
					parens++
				case _CLOSE_PAREN:
					if parens == 0 {
						end = ts.current

						log.Println("start:", start, "end:", end)
						parser := Parse(NewTokenItems(ts.items[start:end]))
						if parser != nil {
							ret = append(ret, parser)
						}
						break INNER

					} else {
						parens--
					}
				}
			}

		case _CLOSE_PAREN:
			panic("never happen")
		case _RAW:
			log.Println("into raw", item.value)
			ret = append(ret, &Raw{text: item.value})
		case _TEXT:
			ret = append(ret, &Text{text: item.value})
		default:
		}
	}

	return ret
}

// 集合
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

func (and *AND) Len() int {
	return 2
}

func (and *AND) Parse() []*QueryItem {
	result := make([]*QueryItem, 0, 2)
	result = append(result, and.left.Parse()...)
	result = append(result, and.right.Parse()...)

	return result
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
