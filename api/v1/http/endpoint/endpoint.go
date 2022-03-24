package endpoint

type Endpoint[Params, Body, Entity any] struct {
	Method string
	Path   string
}

func NewEndpoint[Params, Body, Entity any](method, path string) *Endpoint[Params, Body, Entity] {
	return &Endpoint[Params, Body, Entity]{
		method,
		path,
	}
}
