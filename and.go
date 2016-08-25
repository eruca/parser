package parser

import "fmt"

type And struct {
	left  Parser
	right Parser
}

func (and *And) Len() int {
	return 1
}

func (and *And) String() string {
	return fmt.Sprintf("%s && %s", and.left, and.right)
}

func (and *And) Qt() QueryType {
	return MUST
}

func (and *And) SetQt(QueryType) {}

func (and *And) Parse() ([]Group, error) {
	ps := Parsers{QT: MUST, Items: []Parser{and.left, and.right}}

	return ps.Parse()
}
