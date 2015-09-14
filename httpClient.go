package main

import (
	"net/http"
	"log"
	"io/ioutil" 
	"net/url"
	"net/http/httputil"
	"fmt"
	"bytes"
	"strings"
	"strconv"
)

type HttpClient struct {
	apps []App
	debug_flag bool
}

type HttpParams map[string]string

func (h *HttpClient)set_debug_mode(new_status bool) {
	h.debug_flag = new_status
}

func (h *HttpClient)is_debug_mode() bool {
	return h.debug_flag
}

func (h *HttpClient)get(_url string) string {
	resp, err := http.Get(_url)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}		

	body_bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(body_bytes)
}

func (h *HttpClient)post(_url string, _params HttpParams) string {
	data := url.Values{}
	if len(_params) > 0 {
		for key, value := range _params {
			data.Add(key, value)
			// fmt.Println(key, value)
		}
	}

	r, err := http.NewRequest("POST", _url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatal(err)
		return ""
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	if h.is_debug_mode() {
		dump, _ := httputil.DumpRequestOut(r, true)
		prefix := "> "
		fmt.Println(prefix + strings.Replace(string(dump), "\n", "\n" + prefix , -1))
		fmt.Println()
	}


	client := &http.Client{}
	resp, _ := client.Do(r)

	if h.is_debug_mode() {
		fmt.Println(resp.Status)
		for k, v := range resp.Header {
			fmt.Printf("- %s: %s", k, v)
		}

		dump, _ := httputil.DumpResponse(resp, true)
		// fmt.Println("<<< " + string(dump[0:100]) + "...")
		prefix := "< "
		// @TODO 테스트로 디버깅 시 메시지 출력 자릿수 제한
		s := string(dump)
		s_len := len(s)
		if s_len > 200 {
			s_len = 200
		}
		fmt.Println(prefix + strings.Replace(s[:s_len], "\n", "\n" + prefix , -1))
		fmt.Println()
	}

	body_bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return string(body_bytes)
}
