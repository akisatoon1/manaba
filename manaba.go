package manaba

import (
	"fmt"
	"net/http/cookiejar"
)

func e(funcName string, e error) error {
	return fmt.Errorf("%v: %v", funcName, e.Error())
}

func Login(jar *cookiejar.Jar, username string, password string) error {
	// 1th flow
	res1, err := get(jar, HOME_URL)
	if err != nil {
		return e("get", err)
	}
	defer res1.Body.Close()

	doc1, err := checkResponse(res1, res1.Body, FIRST_POST_URL, SSO_TITLE)
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
	res2, err := post(jar, FIRST_POST_URL, data1)
	if err != nil {
		return e("post", err)
	}
	defer res2.Body.Close()

	doc2, err := checkResponse(res2, res2.Body, FIRST_POST_URL, META_TITLE)
	if err != nil {
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

	doc3, err := checkResponse(res3, res3.Body, META_URL, SAML_TITLE)
	if err != nil {
		return e("checkResponse 3", err)
	}

	data2, err := setFormData(doc3)
	if err != nil {
		return e("setFormData", err)
	}

	// 4th flow
	res4, err := post(jar, SECOND_POST_URL, data2)
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
