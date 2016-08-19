package parser

import (
	"fmt"
	"log"
	"strings"
)

type Parser interface {
	Len() int
	String() string
	Parse(groups Groups)
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
		log.Println("into _And")
		return &AND{
			left:  Parse(NewTokenItems(tokenItems.items[:pos])),
			right: Parse(NewTokenItems(tokenItems.items[pos+1:])),
		}
	} else {
		log.Println("into _OR")
		return &OR{
			left:  Parse(NewTokenItems(tokenItems.items[:pos])),
			right: Parse(NewTokenItems(tokenItems.items[pos+1:])),
			index: item.index,
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

func (ps Parsers) Parse(groups Groups) {
	for _, p := range ps {
		p.Parse(groups)
	}
}

// 代表 左右两个都必须有
type AND struct {
	left  Parser
	right Parser
}

func (and *AND) Len() int {
	return 2
}

func (and *AND) String() string {
	return fmt.Sprintf("%s && %s", and.left, and.right)
}

func (and *AND) Parse(groups Groups) {
	and.left.Parse(groups)
	and.right.Parse(groups)
}

// 代表 或者
type OR struct {
	left  Parser
	right Parser
	index int
}

func (or *OR) Len() int {
	return 1
}

func (or *OR) String() string {
	return fmt.Sprintf("%s || %s", or.left, or.right)
}

func (or *OR) Parse(groups Groups) {
	index := or.index + 1

	leftGroup := groups[:index]
	rightGroup := groups[index:]
	or.left.Parse(leftGroup)
	or.right.Parse(rightGroup)
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

func (t *Text) Parse(groups Groups) {
	for _, group := range groups {
		group.items = append(group.items, &QueryItem{Text: t.text, Offset: true})
	}
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

func (r *Raw) Parse(groups Groups) {
	for _, group := range groups {
		group.items = append(group.items, &QueryItem{Text: r.text})
	}
}
