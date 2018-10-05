package pay

import (
	"net/url"
	"time"
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/axgle/mahonia"
	"encoding/json"
)

type BizContent struct {
	Subject     string  `json:"subject"`
	OutTradeNo  string  `json:"out_trade_no"`
	TotalAmount float32 `json:"total_amount"`

	ProductCode    string `json:"product_code"`
	TimeoutExpress string `json:"timeout_express"`
	Body           string `json:"body"`
	//"timeout_express":"30m","product_code":"QUICK_MSECURITY_PAY","total_amount":"0.01","subject":"1","body":"我是测试数据","out_trade_no":"IQJZSRC1YMQB5HU"
}

func (bc *BizContent) Marshal() string {
	if len(bc.ProductCode) < 2 {
		bc.ProductCode = "QUICK_MSECURITY_PAY"
	}
	if len(bc.TimeoutExpress) < 1 {
		bc.TimeoutExpress = "30m"
	}
	if len(bc.Body) < 1 {
		bc.Body = bc.Subject
	}
	bs, _ := json.Marshal(bc)
	return string(bs)
}

func Sign(appId, privateKey, returnUrl, notifyUrl, bizContent string) (string, url.Values) {
	var data = url.Values{}
	data.Add("app_id", appId)
	data.Add("method", "alipay.trade.app.pay")
	if returnUrl != "" {
		data.Add("return_url", returnUrl)
	}
	if notifyUrl != "" {
		data.Add("notify_url", returnUrl)
	}
	data.Add("format", "json")
	data.Add("charset", "UTF-8")
	data.Add("sign_type", "RSA2")
	data.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	data.Add("version", "1.0")
	data.Add("biz_content", bizContent)
	//data.Add("AliSign", AliSign(data, privateKey))
	return AliSign(data, privateKey), data
}

func AliPayInfo(p url.Values, sign string) string {
	var data = url.Values{}
	data.Add("app_id", p.Get("app_id"))
	data.Add("return_url", p.Get("return_url"))
	data.Add("notify_url", p.Get("notify_url"))
	data.Add("biz_content",  p.Get("biz_content"))

	data.Add("method", "alipay.trade.app.pay")
	data.Add("format", "json")
	data.Add("charset", "UTF-8")
	data.Add("sign_type", "RSA2")
	data.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	data.Add("version", "1.0")
	data.Add("sign", sign)
	return data.Encode()
}

func Pay(appId, returnUrl, notifyUrl, bizContent, sign string) (s string, err error) {
	var data = url.Values{}
	data.Add("app_id", appId)
	data.Add("method", "alipay.trade.app.pay")
	if returnUrl != "" {
		data.Add("return_url", returnUrl)
	}
	if notifyUrl != "" {
		data.Add("notify_url", returnUrl)
	}
	data.Add("format", "json")
	data.Add("charset", "UTF-8")
	data.Add("sign_type", "RSA2")
	data.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	data.Add("version", "1.0")
	data.Add("biz_content", bizContent)
	data.Add("sign", sign)

	resp, err := http.PostForm( //"https://openapi.alipay.com/gateway.do"
		"https://openapi.alipaydev.com/gateway.do", data)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	d := mahonia.NewDecoder("utf-8")
	return d.ConvertString(string(body)), err
}
