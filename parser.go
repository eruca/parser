package parser

import (
	"errors"
	"log"
	"strings"
)

type QueryType int

var (
	ErrNotRange         = errors.New("the QueryItem is not range")
	ErrRangeInvalFormat = errors.New("range invalid format, format: [{ A ~ B }], '[':>=,'{':>,']':<=,'}':<")
)

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

func (item *QueryItem) Parse() (start string, startEqual bool, end string, endEqual bool) {
	if !item.isRange {
		log.Panicln(ErrNotRange)
	}

	if item.Value[0] == '[' {
		startEqual = true
	} else if item.Value[0] == '{' {
		startEqual = false
	} else {
		log.Panicln(ErrRangeInvalFormat)
	}

	if item.Value[len(item.Value)-1] == ']' {
		endEqual = true
	} else if item.Value[len(item.Value)-1] == '}' {
		endEqual = false
	} else {
		log.Panicln(ErrRangeInvalFormat)
	}

	result := strings.Split(item.Value[1:len(item.Value)-1], "~")
	if len(result) != 2 {
		log.Panicln(ErrRangeInvalFormat)
	}

	start = result[0]
	end = result[1]
	return
}

type Group struct {
	QT    QueryType
	Items []QueryItem
}
