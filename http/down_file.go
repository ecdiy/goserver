package http

import (
	"net/http"
	"github.com/cihub/seelog"
	"io/ioutil"
	"io"
	"encoding/json"
	"os"
	"bytes"
	"archive/zip"
	"encoding/hex"
	"crypto/md5"
	"utils"
	"strings"
)

//func (th *Http) Get(Url string) (string, int, error) {
//	return th.http("GET", Url)
//}
//
//func (th *Http) Post(Url string) (string, int, error) {
//	return th.http("POST", Url)
//}

func (th *Http) DownFile(url, dir, defaultExt string, force bool) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(url))
	cipherStr := md5Ctx.Sum(nil)
	md5str := hex.EncodeToString(cipherStr)
	dInd := strings.LastIndex(url, ".")
	ext := defaultExt
	if dInd > 0 {
		ext = url[dInd:]
	}
	lf, uri := utils.FmtImgDir(dir, md5str)
	localFile := lf+ext
	if !force {
		_, e := os.Stat(localFile)
		if os.IsExist(e) {
			return uri+ext
		}
	}
	bs, _, _ := th.GetBody("GET", url)
	ioutil.WriteFile(localFile, bs, 0644)
	return uri+ext
}

func (th *Http) GetResponse(method, Url string) (*http.Response, error) {
	req, _ := http.NewRequest(method, Url, th.Param)
	if th.UA == "" {
		req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 UBrowser/6.1.2716.5 Safari/537.36"}
	} else {
		req.Header["User-Agent"] = []string{th.UA}
	}
	req.Header["Connection"] = []string{"keep-alive"}
	req.Header["Accept"] = []string{"application/json, text/javascript, */*; q=0.01"}
	if method == "POST" {
		req.Header["X-Requested-With"] = []string{"User-Agent:Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 UBrowser/6.1.2716.5 Safari/537.36"}
	}
	if th.Origin != "" {
		req.Header["Origin"] = []string{th.Origin}
	}
	if th.Referer != "" {
		req.Header["Referer"] = []string{th.Referer}
	}
	if th.ContentType != "" {
		req.Header["content-type"] = []string{th.ContentType}
	} else {
		req.Header["content-type"] = []string{"application/x-www-form-urlencoded; charset=UTF-8"}
	}
	if th.Cookie != "" {
		req.Header["cookie"] = []string{th.Cookie}
	}
	if th.Head != nil {
		for k, v := range th.Head {
			req.Header[k] = []string{v}
		}
	}
	c := http.Client{}
	return c.Do(req)

}

func (th *Http) GetUnzip(url, dir string) error {
	os.Mkdir(dir, 0777)
	resp, err := th.GetResponse("GET", url)
	if err != nil {
		seelog.Info("Http Error:", url, err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			os.Mkdir(dir+"/"+file.Name, 0777)
		} else {
			w, err := os.Create(dir + "/" + file.Name)
			if err != nil {
				continue
			}
			rc, err := file.Open()
			if err != nil {
				w.Close()
				continue
			}
			io.Copy(w, rc)
			w.Close()
			rc.Close()
		}
	}
	return nil
}

func (th *Http) Json(method, Url string) (map[string]interface{}, error) {
	resp, err := th.GetResponse(method, Url)
	if err != nil {
		seelog.Info("Http Error:", Url, err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		jsonMap := make(map[string]interface{})
		if e := json.Unmarshal(body, &jsonMap); e == nil {
			return jsonMap, nil
		} else {
			return nil, e
		}
	}
	return nil, err
}

//func (th *Http) http(method, Url string) (string, int, error) {
//	resp, err := th.GetResponse(method, Url)
//	if err != nil {
//		seelog.Info("Http Error:", Url, err)
//		return "", 0, err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		seelog.Error("Http 出错：", err)
//		return "", resp.StatusCode, err
//	} else {
//		b := string(body)
//		if th.FromEncoding != "" && th.ToEncoding != "" {
//			o, errE := iconv.ConvertString(b, th.FromEncoding, th.ToEncoding)
//			if errE != nil {
//				seelog.Error("encode fail:\n", o)
//				return b, resp.StatusCode, nil
//			}
//			return o, resp.StatusCode, nil
//		} else {
//			return b, resp.StatusCode, nil
//		}
//	}
//}

//func (th *Http) getCookie(host string) string {
//	s := ""
//	for _, c := range th.CookieArray {
//		if strings.Index(host, c["domain"].(string)) >= 0 {
//			if len(s) > 0 {
//				s += ";"
//			}
//			s += c["name"].(string) + "=" + c["value"].(string)
//		}
//	}
//	if len(s) == 0 {
//		seelog.Warn("cookie  not match???\n ==> ", host, th.CookieArray)
//	}
//	return s
//}

//-------------

//func HttpGetPost(url, referer string) (*http.Request, error) {
//	return httpReq("POST", url, referer)
//}
//func HttpGetReq(url, referer string) (*http.Request, error) {
//	return httpReq("GET", url, referer)
//}
//func httpReq(method, url, referer string) (*http.Request, error) {
//	req, err := http.NewRequest(method, url, nil)
//	if err != nil {
//		print(err)
//	}
//	req.Header.Set("Pragma", "no-cache")
//	//req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
//	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
//	if referer != "" {
//		req.Header.Set("Referer", referer)
//	}
//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
//	req.Header.Set("Accept", "text/javascript, text/html, application/xml, text/xml, */*")
//	req.Header.Set("Connection", "keep-alive")
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
//	req.Header.Set("Cache-Control", "no-cache")
//	return req, err
//}
//
//func HttpDo(req *http.Request) ([]byte, error) {
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		print(err)
//	}
//	return ioutil.ReadAll(resp.Body)
//}

//func HttpGet(url ... string) (string, error) {
//	resp, err := http.Get(url[0])
//	if resp != nil {
//		defer resp.Body.Close()
//	}
//	if err != nil {
//		seelog.Error("HttpGet出错", url, err)
//		return "", err
//	} else {
//		//fmt.Println("TransferEncoding:", resp.TransferEncoding)
//		input, err := ioutil.ReadAll(resp.Body)
//		if len(url) == 2 {
//			enc := mahonia.NewEncoder(url[1])
//			//converts a  string from UTF-8 to gbk encoding.
//			return enc.ConvertString(string(input)), err
//			//b, e := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(input)), simplifiedchinese.GB18030.NewEncoder()))
//			//return string(b), e
//		} else {
//			if err != nil {
//				seelog.Error("HttpGet出错", url, err)
//				return "", err
//			} else {
//				return string(input), nil
//			}
//		}
//	}
//}

//func HttpGetJsonp(url string) (string, error) {
//	s, e := HttpGet(url)
//	if e != nil {
//		return "", e
//	}
//	en := lang.ConvertUnicode(s)
//	f := strings.Index(en, "(")
//	if f > 0 {
//		en = en[f+1:]
//		f := strings.LastIndex(en, ")")
//		if f > 0 {
//			en = en[0:f ]
//		}
//	}
//	return en, nil
//}

///*
//"http://www.01happy.com/demo/accept.php",
// "name=cjb&xx=000"
//*/
//
///*
//"http://www.01happy.com/demo/accept.php" "a" "123"
//<==> http://www.01happy.com/demo/accept.php?a=123
// */
//func HttpPostForm(args ... string) (string, error) {
//	p := url.Values{}
//	if len(args) > 1 {
//		for i := 1; i < len(args); i += 2 {
//			p[args[i]] = []string{args[i+1]}
//		}
//	}
//	resp, err := http.PostForm(args[0], p)
//	if err != nil {
//		seelog.Error("HttpPostForm 出错：", args[0], err)
//		return "", err
//	} else {
//		defer resp.Body.Close()
//		body, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			seelog.Error("HttpPost 出错：", args[0], err)
//			return "", err
//		} else {
//			return string(body), nil
//		}
//	}
//}
