package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	v2 "google.golang.org/api/oauth2/v2"
)

// AppConf GoogleAuth用config
type AppConf struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
	AuthCode     string `json:"auth_code"`
}

// CallbackRequest コールバックリクエスト
type CallbackRequest struct {
	Code  string `form:"code"`
	State string `form:"state"`
}

var appConf AppConf
var config oauth2.Config

func main() {
	GinRun()
}

// GinRun gin実行
func GinRun() {
	router := gin.Default()

	router.StaticFile("/favicon.ico", "./favicon.ico")

	v1 := router.Group("/v1")
	{
		v1.GET("/login", v1Login)
		v1.GET("/google/oauthcallback", v1GoogleOAuth)
	}

	router.NoRoute(NoRoute)
	router.Run(":8080")
}

// v1Login ログイン
func v1Login(c *gin.Context) {
	url, err := GetGoogleAuthURL()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "システムエラーが発生中です",
		})
	} else {
		c.Redirect(302, url)
	}

}

// v1GoogleOAuth Google認証
func v1GoogleOAuth(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	token, err := GetGoogleCallback(code, state)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "認証に失敗しました",
		})
	} else {
		c.JSON(200, gin.H{
			"code":  200,
			"token": token,
		})
		ID, email, err := GetGoogleInformaion(token)
		if err != nil {
			c.JSON(200, gin.H{
				"code":    500,
				"message": "認証に失敗しました",
			})
		} else {
			c.JSON(200, gin.H{
				"code":  200,
				"ID":    ID,
				"Email": email,
			})
		}
	}

}

// NoRoute (404)Not Foundページ
func NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{
		"title": "Not Found",
	})
}

// GetGoogleCallback GoogleAuth用callback
func GetGoogleCallback(code, state string) (*oauth2.Token, error) {
	context := context.Background()

	token, err := config.Exchange(context, code)
	if err != nil {
		return nil, err
	}

	if appConf.AuthCode != state {
		return nil, err
	}

	if token.Valid() == false {
		return nil, errors.New("vaild token")
	}

	return token, nil
}

// GetGoogleInformaion GoogleアカウントのIDとEmailを取得
func GetGoogleInformaion(token *oauth2.Token) (string, string, error) {
	context := context.Background()
	service, _ := v2.New(config.Client(context, token))
	tokenInfo, _ := service.Tokeninfo().AccessToken(token.AccessToken).Context(context).Do()

	ID := tokenInfo.UserId
	Email := tokenInfo.Email

	return ID, Email, nil
}

// SetGoogleConfig Google用Config
func SetGoogleConfig() error {
	jsonString, err := ioutil.ReadFile("appConf.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonString, &appConf)
	if err != nil {
		return err
	}

	config = oauth2.Config{
		ClientID:     appConf.ClientID,
		ClientSecret: appConf.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"openid", "email"},
		RedirectURL:  appConf.RedirectURL,
	}
	return nil
}

// GetGoogleAuthURL GoogleAuthCodeURLを取得
func GetGoogleAuthURL() (string, error) {
	SetGoogleConfig()

	url := config.AuthCodeURL(appConf.AuthCode)
	return url, nil
}
