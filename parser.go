package parser

import "log"

type QueryType int

const (
	MUST QueryType = iota
	SHOULD
	MUSTNOT
)

type QueryItem struct {
	QT     QueryType
	Text   string
	Offset bool
}

type Group struct {
	items []*QueryItem
}

type Groups []*Group

func Convert(parser Parser, groups *Groups) {
	if ps, ok := parser.(Parsers); ok {
		log.Println("Is Parsers")
		for _, p := range ps {
			Convert(p, groups)
		}
	}

	switch parser.(type) {
	case *Raw:
		log.Println("len(groups)", len(*groups))
		for _, group := range *groups {
			item := &QueryItem{Text: parser.(*Raw).text}
			group.items = append(group.items, item)
		}
	case *Text:
		for _, group := range *groups {
			group.items = append(group.items, &QueryItem{Text: parser.(*Text).text})
		}
	case *OR:
		or := parser.(*OR)
		index := or.index + 1

		leftgroup := (*groups)[:index]
		rightGroup := (*groups)[index:]

		Convert(or.left, &leftgroup)
		Convert(or.right, &rightGroup)
	case *AND:
		and := parser.(*AND)
		Convert(and.left, groups)
		Convert(and.right, groups)
	}

}
