package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcQueryType(t *testing.T) {
	qt, err := calcQueryType(1, 1, 1)
	assert.Equal(t, qt, QueryType(1))
	assert.NoError(t, err)

	qt, err = calcQueryType(1, 0, 0)
	assert.Equal(t, qt, QueryType(1))
	assert.NoError(t, err)

	qt, err = calcQueryType(0, 0, 0)
	assert.Equal(t, qt, QueryType(0))
	assert.NoError(t, err)

	qt, err = calcQueryType(0, 1, 0)
	assert.Equal(t, qt, QueryType(1))
	assert.NoError(t, err)

	qt, err = calcQueryType(0, 0, 1)
	assert.Equal(t, qt, QueryType(1))
	assert.NoError(t, err)

	qt, err = calcQueryType(1, 0, -1)
	assert.Equal(t, qt, QueryType(0))
	assert.EqualError(t, err, ErrQueryTypeConflict.Error())

	qt, err = calcQueryType(0, 1, -1)
	assert.Equal(t, qt, QueryType(0))
	assert.EqualError(t, err, ErrQueryTypeConflict.Error())

	qt, err = calcQueryType(0, 0, -1)
	assert.Equal(t, qt, QueryType(-1))
	assert.NoError(t, err)
}
