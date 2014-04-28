package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    //"net/url"
)

func main() {
    back := make(map[string]interface{})
    //back["Year"] = 2008
    //back["Mon"] = 7
    //back["Day"] = 1
    //back["N"] = 1
    //back["M"] = 1
    back["name"] = "lul"
    back_json, err := json.Marshal(back)

    url := "http://192.168.1.33:8888/api/CheckUserOnline"
    client := &http.Client{
        CheckRedirect: nil,
    }
    postBytesReader := bytes.NewReader(back_json)
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

    response, err := client.Do(reqest)

    if err != nil {
        fmt.Println("err response", err)
        return
    }
    defer response.Body.Close()
    body, e := ioutil.ReadAll(response.Body)
    if e != nil {
        fmt.Println("err response")
        return
    }
    fmt.Println(string(body))
}
