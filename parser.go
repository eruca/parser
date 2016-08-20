package parser

type QueryType int

const (
	MUST QueryType = iota
	SHOULD
	MUSTNOT
)

type QueryItem struct {
	QT        QueryType
	Attribute string // 属性值
	Text      string
	Offset    bool
}

type Group struct {
	items []*QueryItem
}

type Groups []*Group
