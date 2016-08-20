package parser

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Parser interface {
	Len() int
	String() string
	Parse(groups Groups)
}

func Parse(tokenItems *TokenItems) (Parser, error) {
	if tokenItems == nil {
		return nil, nil
	}

	item, pos := tokenItems.TopAndOr()

	log.Println("pos", pos)
	if pos == -1 {
		return Simple(tokenItems)
	}

	if item.t == _AND {
		log.Println("into _And")
		left, err := Parse(NewTokenItems(tokenItems.items[:pos]))
		if err != nil {
			return nil, err
		}
		right, err := Parse(NewTokenItems(tokenItems.items[pos+1:]))
		if err != nil {
			return nil, err
		}
		return &AND{left: left, right: right}, nil
	} else {
		log.Println("into _OR")
		left, err := Parse(NewTokenItems(tokenItems.items[:pos]))
		if err != nil {
			return nil, err
		}
		right, err := Parse(NewTokenItems(tokenItems.items[pos+1:]))
		if err != nil {
			return nil, err
		}
		return &OR{left: left, right: right, index: item.index}, nil
	}

	return nil, nil
}

// 没有or and
func Simple(ts *TokenItems) (Parser, error) {
	ret := Parsers{}

	start, parens := 0, 0
	for ts.hasNext() {
		item := ts.next()

		switch item.t {
		case _OPEN_PAREN:
			start = ts.current + 1
			var parser Parser
			var err error

		INNER:
			for ts.hasNext() {
				next := ts.next()

				switch next.t {
				case _OPEN_PAREN:
					parens++
				case _CLOSE_PAREN:
					if parens == 0 {
						log.Println("start:", start, "end:", ts.current)
						parser, err = Parse(NewTokenItems(ts.items[start:ts.current]))
						if err != nil {
							return nil, err
						}

						break INNER
					} else {
						parens--
					}
				}
			}

			if len(ret) > 0 {
				last := ret[len(ret)-1]
				if attr, ok := last.(*Attribute); ok {
					attr.right = parser
					break
				}
			}
			ret = append(ret, parser)

		case _COLON:
			if len(ret) == 0 {
				return nil, errors.New("错误语法，不能以:开头")
			}
			ret[len(ret)-1] = &Attribute{
				left: ret[len(ret)-1],
			}

		case _CLOSE_PAREN:
			panic("never happen")
		case _RAW:
			log.Println("into raw", item.value)
			// ret = append(ret, &Raw{text: item.value})
			if len(ret) > 0 {
				last := ret[len(ret)-1]
				if attr, ok := last.(*Attribute); ok {
					attr.right = &Raw{text: item.value}
					break
				}
			}
			ret = append(ret, &Raw{text: item.value})

		case _TEXT:
			log.Println("into text", item.value)
			if len(ret) > 0 {
				last := ret[len(ret)-1]
				if attr, ok := last.(*Attribute); ok {
					attr.right = &Text{text: item.value}
					break
				}
			}
			ret = append(ret, &Text{text: item.value})

		case _EMPTYSPACE:
			ret = append(ret, Sep(0))
		default:
		}
	}

	return ret, nil
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

// 属性
type Attribute struct {
	left  Parser
	right Parser
}

func (a *Attribute) Len() int {
	return a.right.Len()
}

func (a *Attribute) String() string {
	return fmt.Sprintf("%s : %s", a.left, a.right)
}

func (a *Attribute) Parse(groups Groups) {
	var items []*QueryItem
	if ps, ok := a.right.(Parsers); ok {
		items = make([]*QueryItem, 0, len(ps))

		for i := 0; i < len(ps); i++ {
			if ps[i].Len() > 0 {
				items = append(items, &QueryItem{
					Attribute: a.left.String(),
					Text:      ps[i].String(),
				})
			}
		}
	} else {
		items = []*QueryItem{&QueryItem{
			Attribute: a.left.String(),
			Text:      a.right.String(),
		}}
	}

	for _, group := range groups {
		for _, item := range items {
			group.items = append(group.items, item)
		}
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

type Sep int8

func (s Sep) Len() int {
	return 0
}

func (s Sep) String() string {
	return " "
}

func (s Sep) Parse(Groups) {}

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
