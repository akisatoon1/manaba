package manaba

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func makeClient(jar *cookiejar.Jar) *http.Client {
	return &http.Client{
		Jar: jar,
	}
}

func statusCodeErr(c int) error {
	return fmt.Errorf("status code is not 200 but %v", c)
}

func get(jar *cookiejar.Jar, url string) (*http.Response, error) {
	client := makeClient(jar)
	res, err := client.Get(url)
	if err != nil {
		return nil, e("http.Client.Get", err)
	}
	if c := res.StatusCode; c != 200 {
		return nil, statusCodeErr(c)
	}
	return res, nil
}

func postForm(jar *cookiejar.Jar, url string, data url.Values) (*http.Response, error) {
	client := makeClient(jar)
	res, err := client.PostForm(url, data)
	if err != nil && err != io.EOF {
		return nil, e("http.client.PostForm", err)
	}
	if c := res.StatusCode; c != 200 {
		return nil, statusCodeErr(c)
	}
	return res, nil
}

func postMultipart(jar *cookiejar.Jar, url string, contentType string, body io.Reader) error {
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", contentType)
	client := makeClient(jar)
	r, err := client.Do(req)
	if err != nil && err != io.EOF {
		return e("client.Do", err)
	}
	defer r.Body.Close()
	if c := r.StatusCode; c != 200 {
		return statusCodeErr(c)
	}
	return nil
}
