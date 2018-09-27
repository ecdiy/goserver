package http

import (
	"strings"
	"io"
	"github.com/cihub/seelog"
	"io/ioutil"
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
