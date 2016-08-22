package parser2

type Raw struct {
	qt   QueryType
	text string
}

func (r *Raw) Len() int {
	return 1
}

func (r *Raw) String() string {
	return r.text
}

func (r *Raw) Qt() QueryType {
	return r.qt
}

func (r *Raw) SetQt(qt QueryType) {
	r.qt = qt
}

func (r *Raw) Parse() ([]Group, error) {
	return []Group{Group{Items: []QueryItem{{QT: r.qt, Value: r.text}}}}, nil
}
