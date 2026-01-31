package config

var AliPayConfig = struct {
	AppID        string
	PrivateKey   string
	AliPublicKey string
	NotifyURL    string
	ReturnURL    string
}{
	AppID:        "11111111111111111111111",
	PrivateKey:   "********************************",
	AliPublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBC",
	NotifyURL:    "http://b23996e2.natappfree.cc/pay/alipay/notify",
	ReturnURL:    "http://b23996e2.natappfree.cc/pay/alipay/return",
}
