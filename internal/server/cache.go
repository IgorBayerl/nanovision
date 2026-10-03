package server

import (
	"bytes"
	"compress/gzip"
	"container/list"
	"sync"
)

// viewCache keeps built views gzip-compressed. Stored runs never change, so an
// entry never goes stale; only memory pressure evicts it. Concurrent requests
// for a view that is being built wait for that build.
type viewCache struct {
	mu       sync.Mutex
	maxBytes int64
	bytes    int64
	order    *list.List
	items    map[string]*list.Element
	pending  map[string]*pendingView
}

type cachedView struct {
	key string
	gz  []byte
}

type pendingView struct {
	done chan struct{}
	gz   []byte
	err  error
}

func newViewCache(maxBytes int64) *viewCache {
	return &viewCache{
		maxBytes: maxBytes,
		order:    list.New(),
		items:    make(map[string]*list.Element),
		pending:  make(map[string]*pendingView),
	}
}

// get returns the gzip bytes of a view, building it with build on a miss.
func (c *viewCache) get(key string, build func() ([]byte, error)) ([]byte, error) {
	c.mu.Lock()
	if el, ok := c.items[key]; ok {
		c.order.MoveToFront(el)
		c.mu.Unlock()
		return el.Value.(*cachedView).gz, nil
	}
	if p, ok := c.pending[key]; ok {
		c.mu.Unlock()
		<-p.done
		return p.gz, p.err
	}
	p := &pendingView{done: make(chan struct{})}
	c.pending[key] = p
	c.mu.Unlock()

	raw, err := build()
	if err == nil {
		p.gz, err = gzipBytes(raw)
	}
	p.err = err
	close(p.done)

	c.mu.Lock()
	delete(c.pending, key)
	if err == nil {
		c.add(key, p.gz)
	}
	c.mu.Unlock()
	return p.gz, p.err
}

func (c *viewCache) add(key string, gz []byte) {
	if int64(len(gz)) > c.maxBytes {
		return
	}
	c.items[key] = c.order.PushFront(&cachedView{key: key, gz: gz})
	c.bytes += int64(len(gz))
	for c.bytes > c.maxBytes {
		last := c.order.Back()
		v := last.Value.(*cachedView)
		c.order.Remove(last)
		delete(c.items, v.key)
		c.bytes -= int64(len(v.gz))
	}
}

func gzipBytes(raw []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(len(raw)/6 + 64)
	zw, err := gzip.NewWriterLevel(&buf, gzip.DefaultCompression)
	if err != nil {
		return nil, err
	}
	if _, err := zw.Write(raw); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
