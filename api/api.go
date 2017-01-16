// Copyright © 2015 Alexandr Medvedev <alexandr.mdr@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	apiRequestGetToken string = "https://pddimp.yandex.ru/api2/admin/get_token_result"
	apiRequestPrefix   string = "https://pddimp.yandex.ru/api2/admin/dns/"
	ErrorAnswer        string = "error"

	Type_SRV   string = "SRV"
	Type_TXT   string = "TXT"
	Type_NS    string = "NS"
	Type_MX    string = "MX"
	Type_SOA   string = "SOA"
	Type_A     string = "A"
	Type_AAAA  string = "AAAA"
	Type_CNAME string = "CNAME"
)

type Record struct {
	RecordId   int         `json:"record_id"`
	RecordType string      `json:"type"`
	Domain     string      `json:"domain"`
	Content    string      `json:"content"`
	TTL        int         `json:"ttl"`
	MinTTL     int         `json:"minttl"`
	FQDN       string      `json:"fqdn"`
	Priority   interface{} `json:"priority, string"` // Required only for SRV or MX records
	Subdomain  string      `json:"subdomain"`
	Weight     int         `json:"weight"`     // Required only for SRV records
	Port       int         `json:"port"`       // Required only for SRV records
	Target     string      `json:"target"`     // Required only for SRV records
	AdminMail  string      `json:"admin_mail"` // Required only for SOA records
	Refresh    int         `json:"refresh"`    // Required only for SOA records
	Retry      int         `json:"retry"`      // Required only for SOA records
	Expire     int         `json:"expire"`     // Required only for SOA records
	NegCache   int         `json:"neg_cache"`  // Required only for SOA records
	Operation  string      `json:"operation"`
}

type Response struct {
	Records  []Record `json:"records"`
	Record   Record   `json:"record"`
	RecordId int      `json:"record_id"`
	Domain   string   `json:"domain"`
	Success  string   `json:"success"`
	Error    string   `json:"error"`
	Json     string
}

type ApiError struct {
	msg string
}

func (e ApiError) Error() string {
	return e.msg
}

func GetTokenLink() string {
	return apiRequestGetToken
}

func doRequest(method string, command string, getParams string, token string) (Response, error) {
	var response Response
	client := &http.Client{}
	urlStr := apiRequestPrefix + command + "?" + getParams
	fmt.Printf("Request URL: %s\n\n", urlStr)
	res, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		return response, err
	}

	res.Header.Set("PddToken", token)
	resp, err := client.Do(res)
	defer resp.Body.Close()

	if err != nil {
		return response, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	response.Json = string(body)
	return response, err
}

func recordToQueryString(r Record) string {
	var query []string
	if 0 != r.RecordId {
		query = append(query, "record_id="+strconv.Itoa(r.RecordId))
	}
	if "" != r.RecordType {
		query = append(query, "type="+r.RecordType)
	}
	if "" != r.Content {
		query = append(query, "content="+r.Content)
	}
	if 0 != r.TTL {
		query = append(query, "ttl="+strconv.Itoa(r.TTL))
	}
	if "" != r.AdminMail {
		query = append(query, "admin_mail="+r.AdminMail)
	}
	if "" != r.Priority.(string) {
		query = append(query, "priority="+r.Priority.(string))
	}
	if 0 != r.Weight {
		query = append(query, "weight="+strconv.Itoa(r.Weight))
	}
	if 0 != r.Port {
		query = append(query, "port="+strconv.Itoa(r.Port))
	}
	if "" != r.Target {
		query = append(query, "target="+r.Target)
	}
	if "" != r.Subdomain {
		query = append(query, "subdomain="+r.Subdomain)
	}
	if 0 != r.Refresh {
		query = append(query, "refresh="+strconv.Itoa(r.Refresh))
	}
	if 0 != r.Retry {
		query = append(query, "retry="+strconv.Itoa(r.Retry))
	}
	if 0 != r.Expire {
		query = append(query, "expire="+strconv.Itoa(r.Expire))
	}
	if 0 != r.NegCache {
		query = append(query, "neg_cache="+strconv.Itoa(r.NegCache))
	}
	return strings.Join(query, "&")
}

func copyRecordParams(dst, src *Record) {
	dst.RecordId = src.RecordId
	dst.RecordType = src.RecordType
	dst.Domain = src.Domain
	dst.Subdomain = src.Subdomain
	dst.FQDN = src.FQDN
	dst.Content = src.Content
	dst.TTL = src.TTL
	dst.Priority = src.Priority
}

func GetList(domain, token string) (Response, error) {
	res, err := doRequest("GET", "list", "domain="+domain, token)
	if err != nil {
		return res, err
	}
	if res.Success == ErrorAnswer {
		return res, ApiError{res.Error}
	}
	return res, err
}

func AddRecord(r *Record, domain, token string) (Response, error) {
	query := "domain=" + domain + "&" + recordToQueryString(*r)
	res, err := doRequest("POST", "add", query, token)
	if err != nil {
		return res, err
	}
	if res.Success == ErrorAnswer {
		return res, ApiError{res.Error}
	}
	copyRecordParams(r, &res.Record)

	return res, nil
}

func EditRecord(r *Record, domain, token string) (Response, error) {
	query := "domain=" + domain + "&" + recordToQueryString(*r)
	res, err := doRequest("POST", "edit", query, token)
	if err != nil {
		return res, err
	}
	if res.Success == ErrorAnswer {
		return res, ApiError{res.Error}
	}
	copyRecordParams(r, &res.Record)

	return res, nil
}

func DeleteRecord(r *Record, domain, token string) (Response, error) {
	return DeleteRecordById(r.RecordId, domain, token)
}
func DeleteRecordById(id int, domain, token string) (Response, error) {
	query := "domain=" + domain + "&record_id=" + strconv.Itoa(id)
	res, err := doRequest("POST", "del", query, token)
	if res.Success == ErrorAnswer {
		return res, ApiError{res.Error}
	}

	return res, err
}
