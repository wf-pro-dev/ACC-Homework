package client

import (
	"net/http"
	"net/url"
)

type CookieJar struct {
	cookies []*http.Cookie
}

func (j *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.cookies = cookies
}

func (j *CookieJar) Cookies(u *url.URL) []*http.Cookie {
	return j.cookies
}

func NewClient() (*http.Client, error) {
	cookies, err := LoadCookies()
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Jar: &CookieJar{cookies: cookies},
	}, nil
}
