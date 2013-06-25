package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

const (
	ENDPOINT = "http://mixpanel.com/api"
	VERSION  = "2.0"
)

type ReportClient struct {
	apiKey    string
	apiSecret string
}

func InitClient(apiKey, apiSecret string) *ReportClient {
	return &ReportClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

func (self *ReportClient) Request(path string, params map[string]string) (ret interface{}, err error) {
	args := make(map[string]string, len(params)+3)
	for k, v := range params {
		args[k] = v
	}

	args["api_key"] = self.apiKey
	args["expire"] = fmt.Sprintf("%d", time.Now().Unix()+600)
	args["sig"] = hashArgs(self.apiSecret, args)

	query := url.Values{}
	for k, v := range args {
		query.Add(k, v)
	}
	url := fmt.Sprintf("%s/%s/%s/?%s", ENDPOINT, VERSION, path, query.Encode())
	resp, urlerr := http.Get(url)
	if urlerr != nil {
		err = urlerr
		return
	}

	defer resp.Body.Close()
	body, httperr := ioutil.ReadAll(resp.Body)
	if httperr != nil {
		err = httperr
		return
	}
	jsonerr := json.Unmarshal(body, &ret)
	if jsonerr != nil {
		err = jsonerr
		return
	}
	return ret, nil
}

func hashArgs(secret string, args map[string]string) string {
	var keys sort.StringSlice = make([]string, len(args))
	i := 0
	for k, _ := range args {
		keys[i] = k
		i++
	}
	keys.Sort()
	md5 := md5.New()
	for _, k := range keys {
		v, _ := args[k]
		io.WriteString(md5, k)
		io.WriteString(md5, "=")
		io.WriteString(md5, v)
	}
	io.WriteString(md5, secret)

	sum := md5.Sum(make([]byte, md5.Size()))
	return hex.EncodeToString(sum)[32:]
}
