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
	Parse() Groups
}

func recur_split(parser Parser) Parsers {
	res := Parsers{}

	if ps, ok := parser.(Parsers); ok {
		for _, p := range ps {
			ps_child := recur_split(p)
			res = append(res, ps_child...)
		}
	} else {
		res = append(res, parser)
	}

	return res
}

func recur_count_or(parser Parser) int {
	sum := 1
	if ps, ok := parser.(Parsers); ok {
		for _, p := range ps {
			if p.Len() > 0 {
				num := recur_count_or(p)
				if num > 0 {
					sum *= num
				}
			}
		}

		if sum == 1 {
			return 0
		}
		return sum
	}

	switch parser.(type) {
	case *OR:
		or := parser.(*OR)
		left := recur_count_or(or.left)
		right := recur_count_or(or.right)

		if left+right == 0 {
			return 2
		}

		return left + right + 1
	case *AND:
		and := parser.(*AND)
		left := recur_count_or(and.left)
		right := recur_count_or(and.right)

		if left == 0 {
			return right
		}

		if right == 0 {
			return left
		}

		return left * right
	}
	return 0
}

func Parse(tokenItems *TokenItems) (Parser, error) {
	if tokenItems == nil {
		return nil, nil
	}

	item, pos := tokenItems.TopAndOr()

	// log.Println("pos", pos)
	if pos == -1 {
		return Simple(tokenItems)
	}

	if item.t == _AND {
		// log.Println("into _And")
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
		// log.Println("into _OR")
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
						// log.Println("start:", start, "end:", ts.current)
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
			// log.Println("into raw", item.value)
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
			// log.Println("into text", item.value)
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

// func (ps Parsers) Parse(groups Groups) {
// 	for _, p := range ps {
// 		p.Parse(groups)
// 	}
// }

func (ps Parsers) Parse() Groups {
	groups := make([]Groups, 0, len(ps))
	for _, p := range ps {
		if p.Len() > 0 {
			groups = append(groups, p.Parse())
		}
	}
	if len(groups) == 1 {
		return groups[0]
	}

	cnt := 1
	for _, gs := range groups {
		log.Println("gs:", len(gs))
		cnt *= len(gs)
	}

	log.Println("len(cnt):", cnt)
	ret := make(Groups, cnt)

	for i := 0; i < cnt; i++ {
		ret[i] = &Group{}
	}

	for _, gs := range groups {
		for i := 0; i < cnt; i++ {
			index := i % len(gs)
			ret[i].items = append(ret[i].items, gs[index].items...)
		}
	}
	return ret
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

// func (a *Attribute) Parse(groups Groups) {
// 	ps := recur_split(a.right)
// 	items := make([]*QueryItem, 0, len(ps))

// 	for i := 0; i < len(ps); i++ {
// 		if ps[i].Len() > 0 {
// 			items = append(items, &QueryItem{
// 				Attribute: a.left.String(),
// 				Text:      ps[i].String(),
// 			})
// 		}
// 	}

// 	for _, group := range groups {
// 		for _, item := range items {
// 			group.items = append(group.items, item)
// 		}
// 	}
// }
func (a *Attribute) Parse() Groups {
	groups := a.right.Parse()

	for _, group := range groups {
		for _, item := range group.items {
			item.Attribute = a.left.String()
		}
	}

	return groups
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

// func (and *AND) Parse(groups Groups) {
// 	and.left.Parse(groups)
// 	and.right.Parse(groups)
// }
func (and *AND) Parse() Groups {
	ps := Parsers{and.left, and.right}

	return ps.Parse()
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

// func (or *OR) Parse(groups Groups) {
// 	index := or.index + 1
// 	log.Println("index:", index)
// 	leftGroup := groups[:index]
// 	rightGroup := groups[index:]
// 	or.left.Parse(leftGroup)
// 	or.right.Parse(rightGroup)
// }
func (or *OR) Parse() Groups {
	leftGroup := or.left.Parse()
	rightGroup := or.right.Parse()

	log.Println("len", len(leftGroup), len(rightGroup))

	ret := make(Groups, len(leftGroup)+len(rightGroup))

	copy(ret, leftGroup)
	copy(ret[len(leftGroup):], rightGroup)
	return ret
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

func (t *Text) Parse() Groups {
	// func (t *Text) Parse(groups Groups) {
	// 	for _, group := range groups {
	// 		group.items = append(group.items, &QueryItem{Text: t.text, Offset: true})
	// 	}
	return Groups{&Group{items: []*QueryItem{&QueryItem{Text: t.text, Offset: true}}}}
}

type Sep int8

func (s Sep) Len() int {
	return 0
}

func (s Sep) String() string {
	return " "
}

func (s Sep) Parse() Groups { return nil }

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

func (r *Raw) Parse() Groups {
	// for _, group := range groups {
	// 	group.items = append(group.items, &QueryItem{Text: r.text})
	// }
	return Groups{&Group{items: []*QueryItem{&QueryItem{Text: r.text}}}}
}
