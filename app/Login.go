package app

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"vote_cj/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const ACCESS_TOKEN_URL = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"

func Login(c *gin.Context) {
	code := c.PostForm("code")
	if code == "" {
		c.JSON(450, gin.H{"msg": "登录失败"})
		c.Abort()
		return
	}

	at_url := fmt.Sprintf(ACCESS_TOKEN_URL, config.Wx.Appid, config.Wx.Appsecret, code)

	res := make(map[string]interface{})
	err := json.Unmarshal(httpGet(at_url), &res)

	if err != nil {
		log.Println(err)
		c.JSON(451, gin.H{"msg": "登录失败"})
		c.Abort()
		return
	}
	if _, ok := res["openid"]; !ok {
		c.JSON(452, gin.H{"msg": "登录失败"})
		c.Abort()
		return
	}
	token, err := CreateToken(res["openid"].(string))
	if err != nil {
		c.JSON(453, gin.H{"msg": "登录失败"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

type Claims struct {
	OpenId    string
	CheckCode string
	jwt.StandardClaims
}

//颁发token
func CreateToken(openid string) (string, error) {
	checkcode := CreateCheckCode(openid)
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		OpenId:    openid,
		CheckCode: checkcode,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "www.scdjw.com.cn", // 签名颁发者
			Subject:   "警察故事投票",           //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.System.Secret))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return tokenString, nil
}

func CreateCheckCode(openid string) string {
	h := md5.New()
	h.Write([]byte(openid))
	return hex.EncodeToString(h.Sum([]byte(config.System.Secret)))
}

//发起get请求
func httpGet(url string) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return []byte("")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("")
	}
	return body
}
