package utils

import "github.com/valyala/fasthttp"

func PostDiscord(tkn, path string, jsonBody []byte, setHeaders func(*fasthttp.RequestHeader)) (body string, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURIBytes([]byte("https://discord.com/api/v8/" + path))
	req.SetBody(jsonBody)

	req.Header.Set("authorization", tkn)
	req.Header.SetMethodBytes([]byte("POST"))
	req.Header.SetContentType("application/json")
	if setHeaders != nil {
		setHeaders(&req.Header)
	}

	if fasthttp.Do(req, resp) == nil {
		body = string(resp.Body())
	}

	fasthttp.ReleaseResponse(resp)
	return
}

func GetDiscord(tkn, path string, setHeaders func(*fasthttp.RequestHeader)) (body string, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURIBytes([]byte("https://discord.com/api/v8/" + path))

	req.Header.Set("authorization", tkn)
	req.Header.SetMethodBytes([]byte("GET"))
	if setHeaders != nil {
		setHeaders(&req.Header)
	}

	if fasthttp.Do(req, resp) == nil {
		body = string(resp.Body())
	}

	fasthttp.ReleaseResponse(resp)
	return
}
