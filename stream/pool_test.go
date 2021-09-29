package stream

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCtxPool(t *testing.T) {

	pool := newParseCtxPool()
	assert.NotNil(t, pool)

	r := &bytes.Buffer{}
	w := &bytes.Buffer{}

	// initial get
	ctx := pool.Get(r, w)
	assert.NotNil(t, pool)
	assert.Equal(t, ctx.r, r)
	assert.Equal(t, ctx.w, w)

	// put
	pool.Put(ctx)

	// get from pool
	ctx = pool.Get(r, w)
	assert.NotNil(t, pool)
}
