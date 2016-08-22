package parser2

import "fmt"

type Or struct {
	left  Parser
	right Parser
}

func (or *Or) Len() int {
	return 2
}

func (or *Or) String() string {
	return fmt.Sprintf("%s || %s", or.left, or.right)
}

func (or *Or) Qt() QueryType {
	return SHOULD
}

func (or *Or) SetQt(qt QueryType) {}

func (or *Or) Parse() ([]Group, error) {
	leftGroup, err := or.left.Parse()
	if err != nil {
		return nil, err
	}
	rightGroup, err := or.right.Parse()
	if err != nil {
		return nil, err
	}

	ret := make([]Group, len(leftGroup)+len(rightGroup))
	copy(ret, leftGroup)
	copy(ret[len(leftGroup):], rightGroup)

	return ret, nil
}
