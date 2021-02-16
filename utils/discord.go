package utils

import "github.com/valyala/fasthttp"

func PostDiscord(tkn, path string, jsonBody []byte, setHeaders func(*fasthttp.RequestHeader)) (body string, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURIBytes([]byte("https://discord.com/api/v8/" + path))

	req.Header.Set("authorization", tkn)
	req.Header.SetMethodBytes([]byte("POST"))
	req.Header.SetContentType("application/json")
	req.SetBody(jsonBody)
	if setHeaders != nil {
		setHeaders(&req.Header)
	}

	err = fasthttp.Do(req, resp)
	if err == nil {
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

	err = fasthttp.Do(req, resp)
	if err == nil {
		body = string(resp.Body())
	}

	fasthttp.ReleaseResponse(resp)
	return
}
