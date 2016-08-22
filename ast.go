package parser2

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

var (
	ErrQueryTypeConflict = errors.New("查询类型相冲突")
)

type Parser interface {
	Len() int
	String() string
	Qt() QueryType
	SetQt(qt QueryType)
	Parse() ([]Group, error)
}

// Parsers 集合
type Parsers struct {
	QT    QueryType
	Items []Parser
}

func (ps *Parsers) Len() int {
	return len(ps.Items)
}

func (ps *Parsers) String() string {
	result := make([]string, len(ps.Items))

	for i := 0; i < len(ps.Items); i++ {
		result[i] = ps.Items[i].String()
	}

	return strings.Join(result, " ")
}

func (ps *Parsers) Qt() QueryType {
	return ps.QT
}

func (ps *Parsers) SetQt(qt QueryType) {
	ps.QT = qt
}

func (ps *Parsers) Parse() ([]Group, error) {
	groups := make([][]Group, 0, len(ps.Items))
	for _, p := range ps.Items {
		if p.Len() > 0 {
			qt, err := calcQueryType(ps.QT, SHOULD, p.Qt())
			if err != nil {
				return nil, err
			}
			p.SetQt(qt)
			group, err := p.Parse()
			if err != nil {
				return nil, err
			}
			groups = append(groups, group)
		}
	}
	if len(groups) == 1 {
		return groups[0], nil
	}

	cnt := 1
	for _, group := range groups {
		cnt *= len(group)
	}
	ret := make([]Group, cnt)

	index := 0
	for _, group := range groups {
		for i := 0; i < cnt; i++ {
			index = i % len(group)
			ret[i].Items = append(ret[i].Items, group[index].Items...)
		}
	}

	return ret, nil
}

func Parse(tokenItems *TokenItems) (Parser, error) {
	if tokenItems == nil {
		return nil, nil
	}

	item, pos := tokenItems.TopAndOr()
	log.Println("pos:", pos)
	if pos == -1 {
		return simple(tokenItems)
	}

	if item.t == _AND {
		log.Println("_AND")
		if tokenItems.baseQT == MUSTNOT {
			return nil, ErrQueryTypeConflict
		}

		left, err := Parse(NewTokenItems(tokenItems.items[:pos], MUST))
		if err != nil {
			return nil, err
		}

		right, err := Parse(NewTokenItems(tokenItems.items[pos+1:], MUST))
		if err != nil {
			return nil, err
		}
		return &And{left: left, right: right}, nil
	} else {
		// 因为是OR,所以就以父的QueryType决定
		log.Println("_OR")
		left, err := Parse(NewTokenItems(tokenItems.items[:pos], tokenItems.baseQT))
		if err != nil {
			return nil, err
		}
		log.Println("Or left:", left.String())

		right, err := Parse(NewTokenItems(tokenItems.items[pos+1:], tokenItems.baseQT))
		if err != nil {
			return nil, err
		}

		log.Println("Or right:", right.String())

		return &Or{left: left, right: right /*, index: item.index*/}, nil
	}

	return nil, nil
}

func simple(ts *TokenItems) (Parser, error) {
	ret := Parsers{}

	start, parens := 0, 0
	for ts.hasNext() {
		item := ts.next()

		switch item.t {
		case _OPEN_PAREN:
			start = ts.current + 1
			var (
				parser Parser
				err    error
			)

		LOOP_OPEN_PAREN:
			for ts.hasNext() {
				next := ts.next()

				switch next.t {
				case _OPEN_PAREN:
					parens++
				case _CLOSE_PAREN:
					if parens == 0 {
						parser, err = Parse(NewTokenItems(ts.items[start:ts.current], ts.baseQT))
						if err != nil {
							return nil, err
						}
						log.Println("start:", start, "end:", ts.current, "baseQT", ts.baseQT,
							"return qt", parser.Qt())

						break LOOP_OPEN_PAREN
					} else {
						parens--
					}
				}
			}
			err = prevAttribute(&ret, ts.baseQT, parser)
			if err != nil {
				return nil, err
			}
			log.Println("after attr", parser.Qt())

		// { [
		case _OPEN_BRACK, _OPEN_BRACE:
			value := make([]string, 1, 5)

			if item.t == _OPEN_BRACE {
				value[0] = "{"
			} else {
				value[0] = "["
			}

		LOOP_OPEN_BRACK:
			for ts.hasNext() {
				next := ts.next()

				switch next.t {
				case _EMPTYSPACE:
					continue

				case _CLOSE_BRACK:
					value = append(value, "]")
					break LOOP_OPEN_BRACK

				case _CLOSE_BRACE:
					value = append(value, "}")
					break LOOP_OPEN_BRACK

				default:
					value = append(value, next.value)
				}
			}

			// log.Println("ret.Len():", ret.Len(), "value:", strings.Join(value, ""))
			if ret.Len() > 0 {
				last := ret.Items[ret.Len()-1]
				if attr, ok := last.(*Attribute); ok {
					qt, err := calcQueryType(ts.baseQT, attr.qt, SHOULD)
					if err != nil {
						return nil, err
					}
					last.(*Attribute).right = &Range{qt: qt, text: strings.Join(value, "")}
					continue
				}
			}

			return nil, fmt.Errorf("以[]{}定义一个范围查询: price:[2~3],but %s", strings.Join(value, ""))

		case _COLON:
			if len(ret.Items) == 0 {
				return nil, errors.New("price:[1~2] 不能以:开头,只能表示属性或以转义符'\\'开始")
			}

			var (
				qt  QueryType
				err error
			)

			last := ret.Items[len(ret.Items)-1]
			switch last.(type) {
			case *Raw, *Text:
				log.Println("last qt", last.Qt())
				qt, err = calcQueryType(ts.baseQT, last.Qt(), SHOULD)
				if err != nil {
					return nil, err
				}

			default:
				return nil, errors.New(":前面不能以字符以外的其他")
			}

			ret.Items[len(ret.Items)-1] = &Attribute{left: last, qt: qt}

		case _CLOSE_PAREN:
			panic("never happen")

		case _RAW:
			log.Println("into _RAW --", item.value)
			err := prevAttribute(&ret, ts.baseQT, &Raw{qt: SHOULD, text: item.value})
			if err != nil {
				return nil, err
			}

		case _TEXT:
			log.Println("into _TEXT --", item.value)
			err := prevAttribute(&ret, ts.baseQT, &Text{qt: SHOULD, text: item.value})
			if err != nil {
				return nil, err
			}

		case _PLUS, _SUB:
			var qt QueryType

			if item.t == _PLUS {
				qt = MUST
				if ts.baseQT == MUSTNOT {
					return nil, ErrQueryTypeConflict
				}
			} else {
				qt = MUSTNOT
				if ts.baseQT == MUST {
					return nil, ErrQueryTypeConflict
				}
			}

		LOOP_PLUS_SUB:
			for ts.hasNext() {
				next, _ := ts.peek(1)

				switch next.t {
				case _RAW:
					ret.Items = append(ret.Items, &Raw{qt: qt, text: next.value})
				case _TEXT:
					ret.Items = append(ret.Items, &Text{qt: qt, text: next.value})

				case _EMPTYSPACE, _COLON:
					break LOOP_PLUS_SUB

				default:
					return nil, errors.New("+(-)后面只能是字符或带引号字符, +(-)A +(-)\"AB\" +(-)price:[1~2]")
				}
				ts.next()
			}
		case _EMPTYSPACE:
			ret.Items = append(ret.Items, Sep(0))
		}
	}

	return &ret, nil
}
