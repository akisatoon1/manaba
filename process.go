package manaba

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func setFormData(doc *goquery.Document) (url.Values, error) {
	inputs := doc.Find("input")

	data := url.Values{}
	inputs.Each(func(i int, sel *goquery.Selection) {
		name, nameExist := sel.Attr("name")
		if !nameExist {
			return
		}
		val, valueExist := sel.Attr("value")
		if valueExist {
			data.Set(name, val)
		} else {
			data.Set(name, "")
		}
	})
	return data, nil
}

func getMetaUrl(doc *goquery.Document) (string, error) {
	meta := doc.Find("meta[http-equiv=refresh]").First()
	content, isExist := meta.Attr("content")
	if !isExist {
		return "", fmt.Errorf("goquery.Selection.Attr: content attribute doesn't exist")
	}

	subs := ";URL="
	i := strings.Index(content, subs)
	if i == -1 {
		return "", nil
	}
	return content[i+len(subs):], nil
}

func setCommonPart(jar *cookiejar.Jar, url string, mw *multipart.Writer) error {
	part, _ := mw.CreateFormField("manaba-form")
	io.WriteString(part, "1")

	part, _ = mw.CreateFormField("SessionValue")
	io.WriteString(part, "@1")

	// get "SessionValue1" field value and write it
	res, err := get(jar, url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil && err != io.EOF {
		return e("goquery.NewDocumentFromReader", err)
	}
	val, isExist := doc.Find("input[name=SessionValue1]").First().Attr("value")
	if !isExist {
		return fmt.Errorf("'value' attribute doesn't exist")
	}
	part, _ = mw.CreateFormField("SessionValue1")
	io.WriteString(part, val)

	return nil
}

func getFileIDs(jar *cookiejar.Jar, url string) ([]string, error) {
	res, err := get(jar, url)
	if err != nil {
		return nil, e("get", err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil && err != io.EOF {
		return nil, e("goquery.NewDocumentFromReader", err)
	}

	var IDs []string
	doc.Find("input.inline").Each(func(_ int, s *goquery.Selection) {
		val, isExist := s.Attr("name")
		if isExist {
			IDs = append(IDs, val)
		}
	})
	return IDs, nil
}
