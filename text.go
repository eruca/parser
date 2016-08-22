package parser2

type Text struct {
	qt   QueryType
	text string
}

func (t *Text) Len() int {
	return 1
}

func (t *Text) String() string {
	return t.text
}

func (t *Text) Qt() QueryType {
	return t.qt
}

func (t *Text) SetQt(qt QueryType) {
	t.qt = qt
}

func (t *Text) Parse() ([]Group, error) {
	return []Group{Group{Items: []QueryItem{QueryItem{QT: t.qt, Value: t.text, Offset: true}}}}, nil
}
