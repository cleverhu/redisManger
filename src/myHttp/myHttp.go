package myHttp

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
)

func Request(method string, url string, data []byte, cookie *string, header string, proxy string, timeout int) (response []byte, resStr string, responseHeader http.Header, err error) {
	if cookie == nil {
		cookie = new(string)
	}
	if strings.ToLower(method) != "get" && strings.ToLower(method) != "post"&& strings.ToLower(method) != "options" {
		return nil, "", nil, errors.New("请求方式错误")
	}
	transport := http.Transport{
		Proxy:           nil,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy != "" {
		if !strings.HasPrefix(proxy, "http") {
			proxy = "http://" + proxy
			parse, err := url2.Parse(proxy)
			if err == nil {
				transport = http.Transport{Proxy: http.ProxyURL(parse)}
			}
		}
	}
	//fmt.Println(url)
	request, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewBuffer(data))
	//fmt.Println(string(data))
	if err != nil {
		return nil, "", nil, err
	}
	if strings.ToUpper(method) == "POST" && header == "" {
		header = "Content-Type:application/x-www-form-urlencoded"
	}
	if header != "" {
		headers := strings.Split(header, "\n")
		for _, h := range headers {
			//fmt.Println(len(h))
			if len(h) == 0 {
				continue
			}
			if fmt.Sprintf("%c", h[0]) == ":" {
				index := strings.Index(h[1:], ":")
				if index != -1 && strings.Index(h, "gzip") == -1 {
					request.Header.Add(h[:index], h[index+1:])
					//fmt.Println(h[:index],h[index+1:])
				}
			} else {
				index := strings.Index(h, ":")
				if index != -1 && strings.Index(h[1:], "gzip") == -1 {
					request.Header.Add(h[:index], h[index+1:])
					if strings.ToLower(h[:index]) == "cookie" {
						if strings.Index(*cookie, strings.TrimSpace(h[index+1:])) == -1 {
							*cookie += h[index+1:]
						}
						//cookie += h[index+1:]
						if strings.HasSuffix(*cookie, ";") == false {
							*cookie += ";"
						}
					}
					//fmt.Println(h[:index], h[index+1:])
				}
			}

		}
	}else{

	}
	client := http.Client{
		Transport: &transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar:     nil,
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Do(request)

	if err != nil {
		return nil, "", nil, err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	values := resp.Header.Values("set-cookie")
	responseHeader = resp.Header

	for i := 0; i < len(values); i++ {
		//fmt.Println(values[i])
		index := strings.Index(values[i], ";")
		*cookie += values[i][:index+1]
	}

	return all, string(all), responseHeader, nil
}
