package core

type responseStorage struct {
	header map[string]string
	body   string
}

func (rs *responseStorage) setResponseBody(body string) {
	rs.body = body
}
