package parser

type QueryType int

const (
	MUSTNOT QueryType = iota - 1
	SHOULD
	MUST
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
