package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-contrib/multitemplate"
	_ "github.com/go-sql-driver/mysql"

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

// createMyRender テンプレートファイル
func createMyRender() multitemplate.Render {
	r := multitemplate.New()
	r.AddFromFiles("Login", "./templates/login.html")
	r.AddFromFiles("Index", "./templates/index.html")
	r.AddFromFiles("Items", "./templates/items.html")
	return r
}

// GinRun gin実行
func GinRun() {
	router := gin.Default()
	router.Static("/css", "./css")
	router.Static("/js", "./js")
	router.Static("/img", "./img")
	router.StaticFile("/favicon.ico", "./favicon.ico")

	router.HTMLRender = createMyRender()

	router.GET("/", Index)
	router.GET("/login", Login)
	router.GET("/items", Items)
	v1 := router.Group("/v1")
	{
		v1.GET("/google/login", v1Login)
		v1.GET("/google/oauthcallback", v1GoogleOAuth)
	}

	router.NoRoute(NoRoute)
	router.Run(":8080")
}

// Login ログイン
func Login(c *gin.Context) {
	c.HTML(200, "Login", gin.H{
		"title": "ログイン｜持ち物管理",
	})
}

// Index トップページ
func Index(c *gin.Context) {
	c.HTML(200, "Index", gin.H{
		"title": "持ち物管理",
	})
}

// Items トップページ
func Items(c *gin.Context) {
	c.HTML(200, "Items", gin.H{
		"title": "持ち物管理",
	})
}

// v1Login ログイン
func v1Login(c *gin.Context) {
	url, err := GetGoogleAuthURL()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "システムエラーが発生中です",
		})
		return
	}

	c.Redirect(302, url)
}

// v1GoogleOAuth Google認証
func v1GoogleOAuth(c *gin.Context) {
	InitDB()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "データベース接続エラーが発生しました",
		})
		return
	}

	code := c.Query("code")
	state := c.Query("state")

	token, err := GetGoogleCallback(code, state)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "認証に失敗しました",
		})
		return
	}

	ID, Email, err := GetGoogleInformaion(token)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "認証に失敗しました",
		})
		return
	}

	transaction, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "データベース接続エラーが発生しました",
		})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	insertSQL := "INSERT INTO `users`(`google_id`, `google_email`, `google_access_token`, `google_expiry`, `created`, `modified`) VALUES (?,?,?,?,?,?)"
	_, err = transaction.Exec(insertSQL, ID, Email, token.AccessToken, token.Expiry, now, now)
	if err != nil {
		transaction.Rollback()
		c.JSON(200, gin.H{
			"code":    500,
			"message": "データベース接続エラーが発生しました",
		})
		return
	}
	transaction.Commit()

	c.Redirect(302, "/")
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

// DbConfig データベース接続用struct
type DbConfig struct {
	Dsn      string `json:"dsn"`
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Database string `json:"database"`
	Charset  string `json:"charset"`
}

var dbConf DbConfig
var db *sql.DB

// InitDB データベース接続
func InitDB() error {
	jsonString, err := ioutil.ReadFile("dbConf.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonString, &dbConf)
	if err != nil {
		return err
	}

	connect := fmt.Sprintf(dbConf.Dsn, dbConf.Username, dbConf.Password, dbConf.Server, dbConf.Database, dbConf.Charset)
	fmt.Println(dbConf)
	db, err = sql.Open("mysql", connect)

	if err != nil {
		return err
	}

	return nil
}
