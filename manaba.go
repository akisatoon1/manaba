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

	err = setCommonPart(jar, url, mw)
	if err != nil {
		return err
	}

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

func SubmitReports(jar *cookiejar.Jar, url string) error {
	body := &bytes.Buffer{} // request body
	mw := multipart.NewWriter(body)

	//
	// create body for multipart/form-data
	//
	part, _ := mw.CreateFormField("action_ReportStudent_commitdone")
	io.WriteString(part, "提出")

	err := setCommonPart(jar, url, mw)
	if err != nil {
		return err
	}

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
