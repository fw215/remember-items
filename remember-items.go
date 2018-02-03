package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	v2 "google.golang.org/api/oauth2/v2"
)

// version chache用version
var version string = "0.0.3"

// AppConf app用config
type AppConf struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
	AuthCode     string `json:"auth_code"`
	CookieSecret string `json:"cookie_secret"`
}

// CallbackRequest コールバックリクエスト
type CallbackRequest struct {
	Code  string `form:"code"`
	State string `form:"state"`
}

var appConf AppConf
var googleConf oauth2.Config

func main() {
	GinRun()
}

// createMyRender テンプレートファイル
func createMyRender() multitemplate.Render {
	r := multitemplate.New()
	r.AddFromFiles("Login", "./templates/login.html")
	r.AddFromFiles("Index", "./templates/index.html")
	r.AddFromFiles("Items", "./templates/items.html")
	r.AddFromFiles("Error", "./templates/error.html")
	return r
}

// GinRun gin実行
func GinRun() {
	SetAppConfig()

	router := gin.Default()
	router.Static("/css", "./css")
	router.Static("/js", "./js")
	router.Static("/img", "./img")
	router.StaticFile("/favicon.ico", "./favicon.ico")

	router.HTMLRender = createMyRender()
	store := sessions.NewCookieStore([]byte(appConf.CookieSecret))
	router.Use(sessions.Sessions("RememberItems", store))

	router.GET("/", Index)
	router.GET("/login", Login)
	router.GET("/items/:CategoryID", Items)
	v1 := router.Group("/v1")
	{
		v1.GET("/google/login", v1Login)
		v1.GET("/google/oauthcallback", v1GoogleCallback)
		v1.GET("/categories", v1Categories)
		v1.GET("/category/:CategoryID", v1CategoryGET)
		v1.POST("/category", v1CategoryPOST)
		v1.GET("/items/:CategoryID", v1Items)
		v1.GET("/item/:ItemID", v1ItemGET)
		v1.DELETE("/item/:ItemID", v1ItemDELETE)
		v1.POST("/item", v1ItemPOST)
	}

	router.NoRoute(NoRoute)
	router.Run(":8080")
}

// Login ログイン
func Login(c *gin.Context) {
	ClearSession(c)

	c.HTML(200, "Login", gin.H{
		"title":   "ログイン｜アイテム管理",
		"version": version,
	})
}

// Index トップページ
func Index(c *gin.Context) {
	if err := LoginCheck(c); err != nil {
		ClearSession(c)
		c.HTML(200, "Error", gin.H{
			"title":       "エラーが発生しました｜アイテム管理",
			"error":       err,
			"description": "5秒後にリダイレクトします...",
			"version":     version,
		})
		return
	}
	c.HTML(200, "Index", gin.H{
		"title":   "アイテム管理",
		"version": version,
	})
}

// Items トップページ
func Items(c *gin.Context) {
	if err := LoginCheck(c); err != nil {
		ClearSession(c)
		c.HTML(200, "Error", gin.H{
			"title":       "エラーが発生しました｜アイテム管理",
			"error":       err,
			"description": "5秒後にリダイレクトします...",
			"version":     version,
		})
		return
	}

	CategoryID := c.Param("CategoryID")
	c.HTML(200, "Items", gin.H{
		"title":      "アイテム管理",
		"CategoryID": CategoryID,
		"version":    version,
	})
}

// v1Login ログイン
func v1Login(c *gin.Context) {
	ClearSession(c)
	url, err := GetGoogleAuthURL()
	if err != nil {
		c.HTML(200, "Error", gin.H{
			"title":       "エラーが発生しました｜アイテム管理",
			"error":       "システムエラーが発生中です",
			"description": "5秒後にリダイレクトします...",
			"version":     version,
		})
		return
	}

	c.Redirect(302, url+"&access_type=offline")
}

// v1GoogleCallback Google認証
func v1GoogleCallback(c *gin.Context) {
	InitDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		c.HTML(200, "Error", gin.H{
			"title":       "エラーが発生しました｜アイテム管理",
			"error":       "データベース接続エラーが発生しました",
			"description": "5秒後にリダイレクトします...",
			"version":     version,
		})
		return
	}

	code := c.Query("code")
	state := c.Query("state")

	token, err := GetGoogleCallback(code, state)
	if err != nil {
		c.HTML(200, "Error", gin.H{
			"title":       "エラーが発生しました｜アイテム管理",
			"error":       "認証に失敗しました",
			"description": "5秒後にリダイレクトします...",
			"version":     version,
		})
		return
	}

	callbackID, Email, err := GetGoogleInformaion(token)
	if err != nil {
		c.HTML(200, "Error", gin.H{
			"title":       "エラーが発生しました｜アイテム管理",
			"error":       "認証に失敗しました",
			"description": "5秒後にリダイレクトします...",
			"version":     version,
		})
		return
	}

	userID = ""
	if err := db.QueryRow("SELECT `user_id` FROM `users` WHERE `google_id` = ? LIMIT 1", callbackID).Scan(&userID); err != nil {
		if err != sql.ErrNoRows {
			c.HTML(200, "Error", gin.H{
				"title":       "エラーが発生しました｜アイテム管理",
				"error":       "データベースエラーが発生しました",
				"description": "5秒後にリダイレクトします...",
				"version":     version,
			})
			return
		}
	}

	transaction, err := db.Begin()
	if err != nil {
		c.HTML(200, "Error", gin.H{
			"title":       "エラーが発生しました｜アイテム管理",
			"error":       "データベースエラーが発生しました",
			"description": "5秒後にリダイレクトします...",
			"version":     version,
		})
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	Expiry := token.Expiry.Format("2006-01-02 15:04:05")
	if userID == "" {
		insertSQL := "INSERT INTO `users`(`google_id`, `google_email`, `google_access_token`, `google_expiry`, `google_refresh_token`, `created`, `modified`) VALUES (?,?,?,?,?,?,?)"
		_, err = transaction.Exec(insertSQL, callbackID, Email, token.AccessToken, Expiry, token.RefreshToken, now, now)
	} else {
		if token.RefreshToken == "" {
			updateSQL := "UPDATE `users` SET `google_access_token` = ?, `google_expiry` = ?, `modified` = ? WHERE `user_id` = ?"
			_, err = transaction.Exec(updateSQL, token.AccessToken, Expiry, now, userID)
		} else {
			updateSQL := "UPDATE `users` SET `google_access_token` = ?, `google_expiry` = ?, `google_refresh_token` = ?, `modified` = ? WHERE `user_id` = ?"
			_, err = transaction.Exec(updateSQL, token.AccessToken, Expiry, token.RefreshToken, now, userID)
		}
	}
	if err != nil {
		transaction.Rollback()
		c.HTML(200, "Error", gin.H{
			"title":       "エラーが発生しました｜アイテム管理",
			"error":       "データベースエラーが発生しました",
			"description": "5秒後にリダイレクトします...",
			"version":     version,
		})
		return
	}
	transaction.Commit()

	session := sessions.Default(c)
	session.Set("accessToken", token.AccessToken)
	session.Save()

	c.Redirect(302, "/")
}

// v1Categories category一覧
func v1Categories(c *gin.Context) {
	if err := LoginCheck(c); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "ログインしてください",
		})
		return
	}

	InitDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "データベース接続エラーが発生しました",
		})
		return
	}

	rows, err := db.Query("SELECT `category_id`, `category_name`, `modified` FROM `categories` WHERE `user_id` = ?", userID)
	if err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "データベースエラーが発生しました",
		})
		return
	}
	defer rows.Close()

	var list []gin.H
	var categoryID, categoryName, modified string
	for i := 0; rows.Next(); i++ {
		if err := rows.Scan(&categoryID, &categoryName, &modified); err != nil {
			c.JSON(200, gin.H{
				"code":  500,
				"error": "データベースエラーが発生しました",
			})
			return
		}
		data := gin.H{
			"category_id":   categoryID,
			"category_name": categoryName,
			"modified":      modified,
		}
		list = append(list, data)
	}

	c.JSON(200, gin.H{
		"code":       200,
		"categories": list,
	})
	return
}

// v1Category category詳細
func v1CategoryGET(c *gin.Context) {
	if err := LoginCheck(c); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "ログインしてください",
		})
		return
	}

	CategoryID := c.Param("CategoryID")

	InitDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "データベース接続エラーが発生しました",
		})
		return
	}

	var CategoryName string
	if err := db.QueryRow("SELECT `category_name` FROM `categories` WHERE `category_id` = ? AND `user_id` = ? LIMIT 1", CategoryID, userID).Scan(&CategoryName); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(200, gin.H{
				"code":   500,
				"errors": "カテゴリが見つかりませんでした",
			})
		} else {
			c.JSON(200, gin.H{
				"code":   500,
				"errors": "データベースエラーが発生しました",
			})
		}
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"category": gin.H{
			"category_id":   CategoryID,
			"category_name": CategoryName,
		},
	})
}

// v1Category category登録更新
func v1CategoryPOST(c *gin.Context) {
	if err := LoginCheck(c); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "ログインしてください",
		})
		return
	}

	CategoryID := c.PostForm("category_id")
	CategoryName := c.PostForm("category_name")

	var resError []string
	if CategoryName == "" {
		resError = append(resError, "カテゴリ名を入力してください")
	}
	if len(resError) > 0 {
		c.JSON(200, gin.H{
			"code":   300,
			"errors": resError,
		})
		return
	}

	InitDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "データベース接続エラーが発生しました",
		})
		return
	}

	var updateID string
	if err := db.QueryRow("SELECT `category_id` FROM `categories` WHERE `category_id` = ? AND `user_id` = ? LIMIT 1", CategoryID, userID).Scan(&updateID); err != nil {
		if err != sql.ErrNoRows {
			c.JSON(200, gin.H{
				"code":   500,
				"errors": "データベースエラーが発生しました",
			})
			return
		}
	}

	transaction, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{
			"code":   500,
			"errors": "データベースエラーが発生しました",
		})
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")

	if updateID == "" {
		insertSQL := "INSERT INTO `categories`(`user_id`, `category_name`, `created`, `modified`) VALUES (?,?,?,?)"
		_, err = transaction.Exec(insertSQL, userID, CategoryName, now, now)
	} else {
		updateSQL := "UPDATE `categories` SET `category_name` = ?, `modified` = ? WHERE `category_id` = ? AND `user_id` = ? "
		_, err = transaction.Exec(updateSQL, CategoryName, now, updateID, userID)
	}
	if err != nil {
		transaction.Rollback()
		c.JSON(200, gin.H{
			"code":   500,
			"errors": "データベースエラーが発生しました",
		})
		return
	}
	transaction.Commit()

	c.JSON(200, gin.H{
		"code": 200,
	})
}

// v1Items category詳細
func v1Items(c *gin.Context) {
	if err := LoginCheck(c); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "ログインしてください",
		})
		return
	}

	InitDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "データベース接続エラーが発生しました",
		})
		return
	}

	CategoryID := c.Param("CategoryID")
	rows, err := db.Query("SELECT `item_id`, `item_name`, `item_image`, `modified` FROM `items` WHERE `category_id` = ? AND `user_id` = ?", CategoryID, userID)
	if err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "データベースエラーが発生しました",
		})
		return
	}
	defer rows.Close()

	var list []gin.H
	var itemID, itemName, itemImage, modified string
	for i := 0; rows.Next(); i++ {
		if err := rows.Scan(&itemID, &itemName, &itemImage, &modified); err != nil {
			c.JSON(200, gin.H{
				"code":  500,
				"error": "データベースエラーが発生しました",
			})
			return
		}
		data := gin.H{
			"item_id":    itemID,
			"item_name":  itemName,
			"item_image": itemImage,
			"modified":   modified,
		}
		list = append(list, data)
	}

	c.JSON(200, gin.H{
		"code":  200,
		"items": list,
	})
}

// v1ItemGET items詳細
func v1ItemGET(c *gin.Context) {
	if err := LoginCheck(c); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "ログインしてください",
		})
		return
	}

	ItemID := c.Param("ItemID")

	InitDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "データベース接続エラーが発生しました",
		})
		return
	}

	var ItemName, ItemImage string
	if err := db.QueryRow("SELECT `item_name`, `item_image` FROM `items` WHERE `item_id` = ? AND `user_id` = ? LIMIT 1", ItemID, userID).Scan(&ItemName, &ItemImage); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(200, gin.H{
				"code":   500,
				"errors": "アイテムが見つかりませんでした",
			})
		} else {
			c.JSON(200, gin.H{
				"code":   500,
				"errors": "データベースエラーが発生しました",
			})
		}
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"item": gin.H{
			"item_id":    ItemID,
			"item_name":  ItemName,
			"item_image": ItemImage,
		},
	})
}

// v1ItemPOST item登録
func v1ItemPOST(c *gin.Context) {
	if err := LoginCheck(c); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "ログインしてください",
		})
		return
	}

	CategoryID := c.PostForm("category_id")
	ItemID := c.PostForm("item_id")
	ItemName := c.PostForm("item_name")
	ItemImage := c.PostForm("item_image")

	var resError []string
	if ItemName == "" {
		resError = append(resError, "アイテム名を入力してください")
	}
	if len(resError) > 0 {
		c.JSON(200, gin.H{
			"code":   300,
			"errors": resError,
		})
		return
	}

	InitDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "データベース接続エラーが発生しました",
		})
		return
	}

	var updateID string
	if err := db.QueryRow("SELECT `item_id` FROM `items` WHERE `category_id` = ? AND `item_id` = ? AND `user_id` = ? LIMIT 1", CategoryID, ItemID, userID).Scan(&updateID); err != nil {
		if err != sql.ErrNoRows {
			c.JSON(200, gin.H{
				"code":   500,
				"errors": "データベースエラーが発生しました",
			})
			return
		}
	}

	transaction, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{
			"code":   500,
			"errors": "データベースエラーが発生しました",
		})
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")

	if updateID == "" {
		insertSQL := "INSERT INTO `items`(`user_id`, `category_id`, `item_name`, `item_image`, `created`, `modified`) VALUES (?,?,?,?,?,?)"
		_, err = transaction.Exec(insertSQL, userID, CategoryID, ItemName, ItemImage, now, now)
	} else {
		updateSQL := "UPDATE `items` SET `item_name` = ?, `item_image` = ?, `modified` = ? WHERE `item_id` = ? AND `user_id` = ? AND `category_id` = ? "
		_, err = transaction.Exec(updateSQL, ItemName, ItemImage, now, updateID, userID, CategoryID)
	}
	if err != nil {
		transaction.Rollback()
		c.JSON(200, gin.H{
			"code":   500,
			"errors": "データベースエラーが発生しました",
		})
		return
	}
	transaction.Commit()

	c.JSON(200, gin.H{
		"code": 200,
	})
}

// v1ItemDELETE item登録
func v1ItemDELETE(c *gin.Context) {
	if err := LoginCheck(c); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "ログインしてください",
		})
		return
	}

	ItemID := c.Param("ItemID")

	InitDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": "データベース接続エラーが発生しました",
		})
		return
	}

	var deleteID string
	if err := db.QueryRow("SELECT `item_id` FROM `items` WHERE `item_id` = ? AND `user_id` = ? LIMIT 1", ItemID, userID).Scan(&deleteID); err != nil {
		if err != sql.ErrNoRows {
			c.JSON(200, gin.H{
				"code":   500,
				"errors": "データベースエラーが発生しました",
			})
			return
		}
	}

	transaction, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{
			"code":   500,
			"errors": "データベースエラーが発生しました",
		})
		return
	}

	if deleteID == "" {
		c.JSON(200, gin.H{
			"code":   500,
			"errors": "データベースエラーが発生しました",
		})
		return
	}

	deleteSQL := "DELETE FROM `items` WHERE `item_id` = ? AND `user_id` = ?"
	_, err = transaction.Exec(deleteSQL, deleteID, userID)
	if err != nil {
		transaction.Rollback()
		c.JSON(200, gin.H{
			"code":   500,
			"errors": "データベースエラーが発生しました",
		})
		return
	}
	transaction.Commit()

	c.JSON(200, gin.H{
		"code": 200,
	})
}

// NoRoute (404)Not Foundページ
func NoRoute(c *gin.Context) {
	c.HTML(404, "Error", gin.H{
		"title":       "ページが見つかりません",
		"error":       "ページが見つかりません",
		"description": "5秒後にリダイレクトします...",
		"version":     version,
	})
}

// GoogleAccessResponse アクセストークン有効チェック
type GoogleAccessResponse struct {
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Scope         string `json:"scope"`
	Exp           string `json:"exp"`
	ExpiresIn     string `json:"expires_in"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	AccessType    string `json:"access_type"`
}

// GoogleRefreshResponse リフレッシュトークン
type GoogleRefreshResponse struct {
	AccessToken string `json:"access_token"`
	ToeknType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
}

var userID string

// LoginCheck ログイン状態
func LoginCheck(c *gin.Context) error {
	session := sessions.Default(c)
	accessToken := session.Get("accessToken")
	if accessToken == nil {
		return errors.New("ログインしてください")
	}

	InitDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		return errors.New("データベース接続エラーが発生しました")
	}

	userID = ""
	var googleRefreshToken string
	if err := db.QueryRow("SELECT `user_id`, `google_refresh_token` FROM `users` WHERE `google_access_token` = ? LIMIT 1", accessToken.(string)).Scan(&userID, &googleRefreshToken); err != nil {
		if err != sql.ErrNoRows {
			return errors.New("データベースエラーが発生しました")
		}
	}
	if userID == "" {
		return errors.New("ログインしてください")
	}

	accessTokenParam := HTTPParam{
		Key:   "access_token",
		Value: accessToken.(string),
	}
	var params []HTTPParam
	params = append(params, accessTokenParam)
	resp, err := RequestGET("https://www.googleapis.com/oauth2/v3/tokeninfo", params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var googleAccessResponse GoogleAccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleAccessResponse); err != nil {
		return err
	}
	if googleAccessResponse.Sub == "" {
		// アクセストークンの有効期限が切れているのでリフレッシュトークンで新しいアクセストークンを取得
		if googleRefreshToken != "" {
			var refreshParams []HTTPParam
			refreshParams = append(refreshParams, HTTPParam{
				Key:   "client_id",
				Value: appConf.ClientID,
			})
			refreshParams = append(refreshParams, HTTPParam{
				Key:   "client_secret",
				Value: appConf.ClientSecret,
			})
			refreshParams = append(refreshParams, HTTPParam{
				Key:   "refresh_token",
				Value: googleRefreshToken,
			})
			refreshParams = append(refreshParams, HTTPParam{
				Key:   "grant_type",
				Value: "refresh_token",
			})
			refreshResp, err := RequestPOST("https://www.googleapis.com/oauth2/v4/token", refreshParams)
			if err != nil {
				return err
			}
			defer refreshResp.Body.Close()

			var googleRefreshResponse GoogleRefreshResponse
			if err := json.NewDecoder(refreshResp.Body).Decode(&googleRefreshResponse); err != nil {
				return err
			}
			if googleRefreshResponse.AccessToken == "" {
				return errors.New("ログインしてください")
			}
			transaction, err := db.Begin()
			if err != nil {
				return errors.New("データベースエラーが発生しました")
			}
			now := time.Now().Format("2006-01-02 15:04:05")

			updateSQL := "UPDATE `users` SET `google_access_token` = ?, `modified` = ? WHERE `user_id` = ?"
			_, err = transaction.Exec(updateSQL, googleRefreshResponse.AccessToken, now, userID)

			if err != nil {
				transaction.Rollback()
				return errors.New("データベースエラーが発生しました")
			}
			transaction.Commit()

			session.Set("accessToken", googleRefreshResponse.AccessToken)
			session.Save()
		} else {
			// リフレッシュトークンがないので再ログイン
			return errors.New("アクセストークンの有効期限が切れています")
		}
	}

	return nil
}

// ClearSession セッションクリア
func ClearSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

// GetGoogleCallback GoogleAuth用callback
func GetGoogleCallback(code, state string) (*oauth2.Token, error) {
	context := context.Background()

	token, err := googleConf.Exchange(context, code)
	if err != nil {
		return nil, err
	}

	if appConf.AuthCode != state {
		return nil, errors.New("不正なアクセスです")
	}

	if token.Valid() == false {
		return nil, errors.New("vaild token")
	}

	return token, nil
}

// GetGoogleInformaion GoogleアカウントのIDとEmailを取得
func GetGoogleInformaion(token *oauth2.Token) (string, string, error) {
	context := context.Background()
	service, _ := v2.New(googleConf.Client(context, token))
	tokenInfo, _ := service.Tokeninfo().AccessToken(token.AccessToken).Context(context).Do()

	ID := tokenInfo.UserId
	Email := tokenInfo.Email

	return ID, Email, nil
}

// HTTPParam リクエスト用パラメータ
type HTTPParam struct {
	Key   string
	Value string
}

// RequestGET GETリクエスト
func RequestGET(target string, params []HTTPParam) (*http.Response, error) {
	values := url.Values{}
	for _, param := range params {
		values.Add(param.Key, param.Value)
	}
	resp, err := http.Get(target + "?" + values.Encode())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// RequestPOST POSTリクエスト
func RequestPOST(target string, params []HTTPParam) (*http.Response, error) {
	values := url.Values{}
	for _, param := range params {
		values.Add(param.Key, param.Value)
	}
	resp, err := http.PostForm(target, values)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SetAppConfig Google用Config
func SetAppConfig() error {
	jsonString, err := ioutil.ReadFile("appConf.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonString, &appConf)
	if err != nil {
		return err
	}

	googleConf = oauth2.Config{
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
	url := googleConf.AuthCodeURL(appConf.AuthCode)
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
	db, err = sql.Open("mysql", connect)

	if err != nil {
		return err
	}

	return nil
}
