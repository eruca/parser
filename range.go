package parser

type Range struct {
	qt   QueryType
	text string
}

func (r *Range) Len() int {
	return 1
}

func (r *Range) String() string {
	return r.text
}

func (r *Range) Qt() QueryType {
	return r.qt
}

func (r *Range) SetQt(qt QueryType) {
	r.qt = qt
}

func (r *Range) Parse() ([]Group, error) {
	return []Group{Group{Items: []QueryItem{{QT: r.qt, Value: r.text, isRange: true}}}}, nil
}
