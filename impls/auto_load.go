package impls

import "github.com/mz-eco/mz/http"

var (
	Handlers []*http.ApiMeta
)

func AddHandler(meta http.Meta, cb func() http.ApiHandler) {
	Handlers = append(Handlers, &http.ApiMeta{
		Meta: meta,
		New:  cb,
	})
}
