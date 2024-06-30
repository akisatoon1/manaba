package manaba

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func Login(jar *cookiejar.Jar, username string, password string) error {
	// 1th flow
	res1, err := get(jar, gHomeUrl)
	if err != nil {
		return e("get", err)
	}
	defer res1.Body.Close()

	doc1, err := checkResponse(res1, res1.Body, gFirstPostUrl, gSsoTitle)
	if err != nil {
		return e("checkResponse 1", err)
	}

	data1, err := setFormData(doc1)
	if err != nil {
		return e("setFormData", err)
	}
	data1.Set("username", username)
	data1.Set("password", password)

	// 2th flow
	res2, err := post(jar, gFirstPostUrl, data1)
	if err != nil {
		return e("post", err)
	}
	defer res2.Body.Close()

	doc2, err := checkResponse(res2, res2.Body, gFirstPostUrl, gMetaTitle)
	if err != nil {
		if isIncorrectUsernameOrPassword(doc2) {
			return fmt.Errorf("incorrect username or password")
		}
		return e("checkResponse 2", err)
	}

	url, err := getMetaUrl(doc2)
	if err != nil {
		return e("getMetaUrl", err)
	}

	// 3th flow
	res3, err := get(jar, url)
	if err != nil {
		return e("get", err)
	}
	defer res3.Body.Close()

	doc3, err := checkResponse(res3, res3.Body, gMetaUrl, gSamlTitle)
	if err != nil {
		return e("checkResponse 3", err)
	}

	data2, err := setFormData(doc3)
	if err != nil {
		return e("setFormData", err)
	}

	// 4th flow
	res4, err := post(jar, gSecondPostUrl, data2)
	if err != nil {
		return e("post", err)
	}
	defer res4.Body.Close()

	err = checkResult(res4)
	if err != nil {
		return e("checkResult", err)
	}

	return nil
}

func UploadFile(jar *cookiejar.Jar, url string, filePath string) error {
	// check url
	isMatch, _ := regexp.MatchString("^https://room.chuo-u.ac.jp/ct/course_[0-9]+_report_[0-9]+$", url)
	if !isMatch {
		return fmt.Errorf("invalid url for uploading file")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{} // request body
	mw := multipart.NewWriter(body)

	//
	// create body for multipart/form-data
	//
	_, fileName := filepath.Split(filePath)
	part, err := mw.CreateFormFile("RptSubmitFile", fileName)
	if err != nil {
		return err
	}
	io.Copy(part, file)

	part, _ = mw.CreateFormField("action_ReportStudent_submitdone")
	io.WriteString(part, "アップロード")

	part, _ = mw.CreateFormField("manaba-form")
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

	mw.Close()

	// POST to url
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", mw.FormDataContentType())
	client := makeClient(jar)
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if c := r.StatusCode; c != 200 {
		return fmt.Errorf("status code is not 200 but %v", c)
	}

	return nil
}
