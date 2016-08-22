package parser2

type QueryType int

const (
	MUSTNOT QueryType = iota - 1
	SHOULD
	MUST
)

type QueryItem struct {
	QT        QueryType
	Attribute string
	Value     string
	isRange   bool
	Offset    bool
}

func (item *QueryItem) IsRange() bool {
	return item.isRange
}

func (item *QueryItem) Start() (value string, isEqual bool) {
	// todo
	return "", false
}

func (item *QueryItem) End() (value string, isEqual bool) {
	// todo
	return
}

type Group struct {
	QT    QueryType
	Items []QueryItem
}
