package model

import (
	"io"
	"net/url"
)

type ReqData struct {
	Method string
	Url    string
	Header map[string]string
	Params url.Values
	Body   io.Reader
}
