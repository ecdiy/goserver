package pay

import (
	"net/url"
	"time"
	"encoding/json"
)

type BizContent struct {
	Subject        string  `json:"subject"`
	OutTradeNo     string  `json:"out_trade_no"`
	TotalAmount    float32 `json:"total_amount"`
	ProductCode    string  `json:"product_code"`
	TimeoutExpress string  `json:"timeout_express"`
	Body           string  `json:"body"`
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

func AliAppSign(appId, privateKey, returnUrl, notifyUrl, bizContent string) (string, url.Values) {
	var data = url.Values{}
	data.Add("app_id", appId)
	data.Add("method", "alipay.trade.app.pay")
	if returnUrl != "" {
		data.Add("return_url", returnUrl)
	}
	if notifyUrl != "" {
		data.Add("notify_url", notifyUrl)
	}
	data.Add("format", "json")
	data.Add("charset", "UTF-8")
	data.Add("sign_type", "RSA2")
	data.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	data.Add("version", "1.0")
	data.Add("biz_content", bizContent)
	return aliSign(data, privateKey), data
}

func AliAppPayInfo(p url.Values, sign string) string {
	var data = url.Values{}
	data.Add("app_id", p.Get("app_id"))
	data.Add("return_url", p.Get("return_url"))
	data.Add("notify_url", p.Get("notify_url"))
	data.Add("biz_content", p.Get("biz_content"))
	data.Add("method", "alipay.trade.app.pay")
	data.Add("format", "json")
	data.Add("charset", "UTF-8")
	data.Add("sign_type", "RSA2")
	data.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	data.Add("version", "1.0")
	data.Add("sign", sign)
	return data.Encode()
}
