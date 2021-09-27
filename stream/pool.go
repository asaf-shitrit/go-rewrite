package stream

import (
	"io"
	"sync"
)

var defaultPool = newParseCtxPool()

type parseContextPool struct {
	pool sync.Pool
}

func (pcp *parseContextPool) Get(r io.Reader, w io.Writer) *parseContext {
	pc, ok := pcp.pool.Get().(*parseContext)
	if !ok {
		return newParseCtx(r, w)
	}

	pc.reset(r, w)
	return pc
}

func (pcp *parseContextPool) Put(pc *parseContext) {
	pc.reset(nil, nil)
	pcp.pool.Put(pc)
}

func newParseCtxPool() *parseContextPool {
	return &parseContextPool{
		sync.Pool{New: func() interface{} {
			return newParseCtx(nil, nil)
		}},
	}
}
