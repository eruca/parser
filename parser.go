package parser

type QueryType int

const (
	MUST QueryType = iota
	SHOULD
	MUSTNOT
)

type QueryItem struct {
	QT     QueryType
	Text   string
	Offset []int
}

type Group struct {
	items []*QueryItem
}

type Groups []*Group
