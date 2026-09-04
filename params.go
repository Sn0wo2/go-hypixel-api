package hypixel

import (
	"fmt"
	"net/url"
)

type Params map[string]any

func (p Params) Get[T any](k string) T {
	val, ok := p[k]
	if !ok {
		var zero T
		return zero
	}
	typed, _ := val.(T)
	return typed
}

func (p Params) Set[T any](k string, v T) {
	p[k] = v
}

func (p Params) Del(k string) {
	delete(p, k)
}

func (p Params) Has(k string) bool {
	_, ok := p[k]
	return ok
}

func (p Params) buildURL(full string) string {
	if len(p) == 0 {
		return full
	}
	u, err := url.Parse(full)
	if err != nil {
		return full
	}

	q := u.Query()
	for k, v := range p {
		if v == nil {
			continue
		}
		val := fmt.Sprint(v)
		if val == "" {
			continue
		}
		q.Set(k, val)
	}

	u.RawQuery = q.Encode()

	return u.String()
}
