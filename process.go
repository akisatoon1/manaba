package manaba

import (
	"fmt"
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
