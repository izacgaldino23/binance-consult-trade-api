package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/izacgaldino23/binance-consult-trade-api/config"
)

type Request struct {
	URL    URL
	Params *PathParamList
	Query  *QueryParamList
	Body   map[string]interface{}
}

type PathParam struct {
	Key   string
	Value string
}

type QueryParam struct {
	Key   string
	Value []string
}

type PathParamList []PathParam

type QueryParamList []QueryParam

type URL string

func (p *QueryParamList) ToQuery() (query string) {
	params := make([]string, 0)

	for _, param := range *p {
		for _, value := range param.Value {
			params = append(params, fmt.Sprintf("%v=%v", param.Key, value))
		}
	}

	if len(params) > 0 {
		query = "?" + strings.Join(params, "&")
	}

	return
}

func (p *QueryParamList) AddParam(key, value string) *QueryParamList {
	*p = append(*p, QueryParam{key, []string{value}})

	return p
}

func (u *URL) bindParams(params *PathParamList) (newURL string) {
	newURL = string(*u)

	if params != nil {
		for i := range *params {
			var keyInURL = fmt.Sprintf("{%v}", (*params)[i].Key)

			newURL = strings.Replace(newURL, keyInURL, (*params)[i].Value, 1)
		}
	}

	return
}

func (r *Request) generateFinalURL() string {
	finalUrl := config.Environment.BinanceEndpoint + r.URL.bindParams(r.Params)

	if r.Query != nil {
		finalUrl += r.Query.ToQuery()
	}

	return finalUrl
}

func (r *Request) getParsedBody() (io.Reader, error) {
	if r.Body == nil {
		return nil, nil
	}

	encoded, err := json.Marshal(r.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(encoded), nil
}

func (r *Request) ParseParams(c *fiber.Ctx) *Request {
	params := c.AllParams()

	if len(params) > 0 && r.Params == nil {
		r.Params = &PathParamList{}
	}

	for i := range params {
		(*r.Params) = append((*r.Params), PathParam{i, params[i]})
	}

	return r
}

func Get(req Request) (body []byte, statusCode int, err error) {
	finalURL := req.generateFinalURL()

	log.Println(finalURL)

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

func Post(req *Request) (body []byte, statusCode int, err error) {
	finalURL := req.generateFinalURL()

	parsedBody, err := req.getParsedBody()
	if err != nil {
		return nil, 0, err
	}

	newReq, err := http.NewRequest("POST", finalURL, parsedBody)
	if err != nil {
		return nil, 0, err
	}

	newReq.Header.Set("X-MBX-APIKEY", config.Environment.APIKey)

	client := http.Client{}

	res, err := client.Do(newReq)
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
