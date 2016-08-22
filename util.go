package parser2

import "log"

// Seperator
type Sep int8

func (s Sep) Len() int {
	return 0
}

func (s Sep) String() string {
	return " "
}

func (s Sep) Qt() QueryType {
	return SHOULD
}

func (s Sep) SetQt(QueryType) {}

func (s Sep) Parse() ([]Group, error) { return nil, nil }

// Slash
type Slash int8

func (s Slash) Len() int {
	return 0
}

func (s Slash) String() string {
	return "\\"
}

func (s Slash) Qt() QueryType {
	return SHOULD
}

func (s Slash) SetQt(QueryType) {}

func (s Slash) Parse() ([]Group, error) { return nil, nil }

// utils
func calcQueryType(baseQt, attrQt, selfQt QueryType) (QueryType, error) {
	switch baseQt {
	case MUST:
		if attrQt >= 0 && selfQt >= 0 {
			return MUST, nil
		}
	case MUSTNOT:
		if attrQt <= 0 && selfQt <= 0 {
			return MUSTNOT, nil
		}
	case SHOULD:
		if attrQt == MUST && selfQt >= 0 {
			return MUST, nil
		} else if attrQt == MUSTNOT && selfQt <= 0 {
			return MUSTNOT, nil
		} else if attrQt == SHOULD {
			return selfQt, nil
		}
	}

	return 0, ErrQueryTypeConflict
}

func prevAttribute(ps *Parsers, baseQt QueryType, p Parser) error {
	if ps.Len() > 0 {
		last := ps.Items[ps.Len()-1]
		if attr, ok := last.(*Attribute); ok {
			qt, err := calcQueryType(baseQt, attr.qt, p.Qt())
			if err != nil {
				return err
			}

			log.Println("prev Attribute:", baseQt, attr.qt, p.Qt(), qt)
			p.SetQt(qt)
			attr.right = p

			log.Println("p.Qt()", p.Qt())
			return nil
		}
	}
	qt, err := calcQueryType(baseQt, SHOULD, p.Qt())
	if err != nil {
		return err
	}
	p.SetQt(qt)
	ps.Items = append(ps.Items, p)

	return nil
}
