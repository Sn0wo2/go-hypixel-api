package hypixel

import (
	"fmt"
	"net/url"
)

// Params http query
// use fmt.Sprint parse to string
type Params map[any]any

func (p *Params) String(full string) string {
	if len(*p) == 0 {
		return full
	}
	u, err := url.Parse(full)
	if err != nil {
		return full
	}

	q := u.Query()
	for k, v := range *p {
		q.Set(fmt.Sprint(k), fmt.Sprint(v))
	}
	u.RawQuery = q.Encode()

	return u.String()
}

func (p *Params) Get(k any) any {
	return (*p)[k]
}

func (p *Params) Set(k, v any) {
	(*p)[k] = v
}

func (p *Params) Del(k any) {
	delete(*p, k)
}

func (p *Params) Has(k any) bool {
	_, ok := (*p)[k]
	return ok
}
