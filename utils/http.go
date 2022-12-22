package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/izacgaldino23/binance-consult-trade-api/config"
)

type Param struct {
	Key   string
	Value string
}

type ParamList []Param

type URL string

func (p *ParamList) ToQuery() (query string) {
	params := make([]string, 0)

	for i := range *p {
		params = append(params, fmt.Sprintf("%v=%v", (*p)[i].Key, (*p)[i].Value))
	}

	if len(params) > 0 {
		query = "?" + strings.Join(params, "&")
	}

	return
}

func (p *ParamList) AddParam(key, value string) *ParamList {
	*p = append(*p, Param{key, value})

	return p
}

func (u *URL) parseParams(params *ParamList) (newURL string) {
	newURL = string(*u)

	for i := range *params {
		var keyInURL = fmt.Sprintf("{%v}", (*params)[i].Key)

		newURL = strings.Replace(newURL, keyInURL, (*params)[i].Value, 1)
	}

	return
}

func Get(url URL, params *ParamList, query *ParamList) (body []byte, statusCode int, err error) {
	finalURL := config.Environment.BinanceEndpoint + url.parseParams(params) + query.ToQuery()

	fmt.Println(finalURL)

	res, err := http.Get(finalURL)
	if err != nil {
		return nil, 0, err
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	statusCode = res.StatusCode

	return
}
