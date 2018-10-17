package http

import (
	"strings"
	"io"
	"github.com/cihub/seelog"
	"io/ioutil"
	"encoding/json"
)

type Http struct {
	UA, FromEncoding, ToEncoding, Referer, Origin, ContentType, Cookie string
	Head                                                               map[string]string
	Param                                                              io.Reader
	//CookieArray                                            []map[string]interface{}
}

func (th *Http) ParamString(p string) {
	th.Param = strings.NewReader(p)
}

func (th *Http) SetHead(n, v string) {
	if th.Head == nil {
		th.Head = make(map[string]string)
	}
	th.Head[n] = v
}
func (th *Http) GetBody(method, Url string) ([]byte, int, error) {
	resp, err := th.GetResponse(method, Url)
	if err != nil {
		seelog.Info("Http Error:", Url, err)
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

func (th *Http) Get(url string) (string, error) {
	bd, _, e := th.GetBody("GET", url)
	if e == nil {
		return string(bd), e
	} else {
		return "", e
	}
}

func (th *Http) GetMap(url string) (map[string]interface{}, error) {
	bd, _, ex := th.GetBody("GET", url)
	if ex == nil {
		jsonMap := make(map[string]interface{})
		if e := json.Unmarshal(bd, &jsonMap); e == nil {
			return jsonMap, nil
		} else {
			return nil, e
		}
	} else {
		return nil, ex
	}
}
