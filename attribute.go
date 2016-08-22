package parser2

import (
	"fmt"
	"log"
)

type Attribute struct {
	qt    QueryType
	left  Parser
	right Parser
}

func (a *Attribute) Len() int {
	return a.right.Len()
}

func (a *Attribute) String() string {
	return fmt.Sprintf("%s : %s", a.left, a.right)
}

func (a *Attribute) Qt() QueryType {
	return a.qt
}

func (a *Attribute) SetQt(qt QueryType) {
	log.Println("Attr", a.qt, qt)
	a.qt = qt
}

func (a *Attribute) Parse() ([]Group, error) {
	groups, err := a.right.Parse()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(groups); i++ {
		for j := 0; j < len(groups[i].Items); j++ {
			groups[i].Items[j].Attribute = a.left.String()
		}
	}

	return groups, nil
}
