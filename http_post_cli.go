package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    //"net/url"
)

func main() {
	str :=`{ "id": 1330395827, "service": "account.verifySession", "data": { "sid": "80a5fe53d3540300005a17e308a4b1fb" }, "game": { "gameId": 12345 }, "encrypt": "md5", "sign": "6e9c3c1e7d99293dfc0c81442f9a9984" }`

    url := "http://passport_i.25pp.com:8080/account?tunnel-command=2852126760"
    client := &http.Client{
        CheckRedirect: nil,
    }
    postBytesReader := bytes.NewReader([]byte(str))
    reqest, _ := http.NewRequest("POST", url, postBytesReader)
    reqest.Header.Set("User-Agent", " Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
    reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    reqest.Header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
    reqest.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
    reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
    reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
    reqest.Header.Set("Cache-Control", "max-age=0")
    reqest.Header.Set("Connection", "keep-alive")
    reqest.Header.Set("Referer", url)
	fmt.Println(reqest.PostForm)
    response, err := client.Do(reqest)

    if err != nil {
        fmt.Println("err response", err)
        return
    }
    body, e := ioutil.ReadAll(response.Body)
    if e != nil {
        fmt.Println("err response")
        return
    }
    fmt.Println(string(body))
    response.Body.Close()
}
