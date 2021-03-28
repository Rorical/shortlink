package hcaptcha

import (
	"fmt"
	"shortlink/internal/config"

	"github.com/dghubble/sling"

	"github.com/gin-gonic/gin"
)

const VERIFY_URL = "https://hcaptcha.com/siteverify"

type HCaptcha struct {
	sling     *sling.Sling
	secretKey string
}

func NewhCaptcha(conf *config.Setting) *HCaptcha {
	s := sling.New().Base(VERIFY_URL)
	return &HCaptcha{
		sling:     s,
		secretKey: conf.HCaptcha.SecretKey,
	}
}

type VerifyParams struct {
	Secret   string `url:"secret,omitempty"`
	Response string `url:"response,omitempty"`
	RemoteIP string `url:"remoteip,omitempty"`
	SiteKey  string `url:"sitekey,omitempty"`
}

type VerifyResponse struct {
	Success bool `url:"success,omitempty"`
	//Challenge_TS string `url:"challenge_ts,omitempty"`
	//HostName     string `url:"hostname,omitempty"`
}

func (self *HCaptcha) Verify(token string, ip string, siteKey string) (*VerifyResponse, error) {
	verifyParams := &VerifyParams{
		Secret:   self.secretKey,
		Response: token,
		RemoteIP: ip,
		SiteKey:  siteKey,
	}
	verifyResponse := &VerifyResponse{}
	_, err := self.sling.New().BodyForm(verifyParams).ReceiveSuccess(verifyResponse)
	fmt.Println(verifyResponse)
	return verifyResponse, err
}

func (self *HCaptcha) GinVerify(c *gin.Context) bool {
	siteKey := c.DefaultPostForm("site-key", "")
	token := c.DefaultPostForm("h-captcha-response", "")
	ip := c.ClientIP()
	if token == "" || siteKey == "" {
		return false
	}

	res, err := self.Verify(token, ip, siteKey)
	if err != nil {
		return false
	}
	return res.Success
}
