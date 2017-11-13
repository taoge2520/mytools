//get httpcode
package models

import (
	"net/http"
	"net/url"
)

//给定网址，获取该网址的http状态码
func GetSiteCode(domain string) (resCode int, err error) {
	u, _ := url.Parse("http://www." + domain)
	q := u.Query()
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return
	}
	resCode = res.StatusCode
	res.Body.Close()
	if resCode != 200 {
		resCode, err = ReTry(domain)
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	return
}
func ReTry(domain string) (resCode int, err error) {
	req, err := http.NewRequest("GET", "https://www."+domain, nil)
	// ...
	req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36`)
	cli := http.Client{}
	resp, err := cli.Do(req)

	resCode = resp.StatusCode
	resp.Body.Close()
	if err != nil {

		return
	}

	return
}
