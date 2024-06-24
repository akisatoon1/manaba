package manaba

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func checkResponse(res *http.Response, body io.ReadCloser, stdUrl string, stdTitle string) (*goquery.Document, error) {
	Url := res.Request.URL
	if u := getUrl(Url); u != stdUrl {
		return nil, fmt.Errorf("unexpected url: '%v'", u)
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil && err != io.EOF {
		return nil, e("NewDocumentFromReader", err)
	}

	title, err := getTitle(doc)
	if err != nil {
		return nil, e("getTitle", err)
	}

	if title != stdTitle {
		return nil, fmt.Errorf("unexpected title: '%v'", title)
	}

	return doc, nil
}

func checkResult(res *http.Response) error {
	Url := res.Request.URL
	if u := getUrl(Url); u != gHomeUrl {
		return fmt.Errorf("fail to get resource. unexpected url: '%v'", u)
	}
	return nil
}

func getUrl(u *url.URL) string {
	return fmt.Sprintf("%v://%v%v", u.Scheme, u.Host, u.Path)
}

func getTitle(doc *goquery.Document) (string, error) {
	title := doc.Find("title")
	if l := title.Length(); l != 1 {
		return "", fmt.Errorf("the number of title is not 1 but %v", l)
	}
	return title.Text(), nil
}
