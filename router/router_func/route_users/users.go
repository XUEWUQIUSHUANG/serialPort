package users

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/smtp"
	"strings"
	"math/rand"
	"time"
	"strconv"
	"fmt"
	"log"
	"crypto/tls"

	datastruct "go_code/serialPort/data_struct"
	"go_code/serialPort/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
}

type UsersCtrl struct{}
var cacheMap map[string]int

func (ctrl *UsersCtrl) Register(c *gin.Context) {
	var User datastruct.User
	cache:=c.Query("cache")
	if err := c.ShouldBind(&User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if cacheMap[User.Username] != 0 {
			cacheInt, _ := strconv.Atoi(cache)
			if cacheMap[User.Username] == cacheInt {
				cacheMap[User.Username] = 0
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"data": "验证码错误"})
				return
			}
		}
		hash := md5.New()
		hash.Write([]byte(User.Password + "frost"))
		hashInBytes := hash.Sum(nil)
		hashString := hex.EncodeToString(hashInBytes)
		User.Password = hashString
		User.Cookie = uuid.New().String()
		result := sql.Register(User)
		if result {
			c.SetCookie("frost_cookie", User.Cookie, 1000*60*60*24*7, "/", "localhost", false, false)
			c.JSON(http.StatusOK, gin.H{"data": "注册成功"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"data": "用户名已存在"})
		}
	}
}

func (ctrl *UsersCtrl) Login(c *gin.Context) {
	var UserLogin datastruct.UserLogin
	if err := c.ShouldBind(&UserLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		hash := md5.New()
		hash.Write([]byte(UserLogin.Password + "frost"))
		hashInBytes := hash.Sum(nil)
		hashString := hex.EncodeToString(hashInBytes)
		UserLogin.Password = hashString
		result := sql.Login(UserLogin)
		if result != "error" {
			c.SetCookie("frost_cookie", result, 1000*60*60*24*7, "/", "localhost", false, false)
			c.JSON(http.StatusOK, gin.H{"data": "登录成功"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"data": "登录失败"})
		}
	}
}

func (ctrl *UsersCtrl) Check_Users(c *gin.Context) {
	Cookie, err := c.Cookie("frost_cookie")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		result := sql.Check_Users(Cookie)
		if result {
			c.JSON(http.StatusOK, gin.H{"data": "已登录", "hasuser": result})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"data": "未登录", "hasuser": false})
		}
	}
}

func (ctrl *UsersCtrl) Get_Cache(c *gin.Context) {
	email := c.Query("email")
	fmt.Println(email)
	smtpHost := "smtp.qq.com"
	smtpPort := "465"
	smtpUser := "3052604797@qq.com"
	smtpPass := "adjtfsdeseludfbf"

	source := rand.NewSource(time.Now().UnixNano())
    r := rand.New(source)
	
	from := smtpUser
	to := []string{email}
	subject := "Subject: Cache\n"
	cache:= r.Intn(900000)+100000
	body :=`
            <html>
            <head>
                <title>验证码</title>
            </head>
            <body>
                <p>请将验证码填入进网页中的输入框：</p>
                <b>` + strconv.Itoa(cache) + `</b>
            </body>
            </html>
            `

	cacheMap = make(map[string]int)
	message := []byte(
		"From: \"" + "ZHDP" + "\" <" + from + ">\r\n" +
			"To: " + strings.Join(to, ", ") + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" + body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", smtpHost, smtpPort), tlsConfig)
	if err != nil {
		log.Printf("Error connecting to SMTP server: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SMTP server connection failed"})
		return
	}
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Printf("Error creating SMTP client: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SMTP client"})
		return
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		log.Printf("Error during authentication: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
		return
	}

	err = client.Mail(from)
	if err != nil {
		log.Printf("Error sending MAIL FROM command: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send MAIL FROM"})
		return
	}

	for _, addr := range to {
		err = client.Rcpt(addr)
		if err != nil {
			log.Printf("Error sending RCPT TO command: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send RCPT TO"})
			return
		}
	}

	w, err := client.Data()
	if err != nil {
		log.Printf("Error sending DATA command: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send DATA"})
		return
	}
	defer w.Close()

	_, err = w.Write(message)
	if err != nil {
		log.Printf("Error writing email message: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email message"})
		return
	}

	cacheMap[email] = cache
	c.JSON(http.StatusOK, gin.H{"data": "OK"})

}

func (ctrl *UsersCtrl) Change_Password(c *gin.Context){
	var User datastruct.User
	cache:=c.Query("cache")
	if err := c.ShouldBind(&User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if cacheMap[User.Username] != 0 {
			cacheInt, _ := strconv.Atoi(cache)
			if cacheMap[User.Username] == cacheInt {
				cacheMap[User.Username] = 0
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"data": "验证码错误"})
				return
			}
		}
		hash := md5.New()
		hash.Write([]byte(User.Password + "frost"))
		hashInBytes := hash.Sum(nil)
		hashString := hex.EncodeToString(hashInBytes)
		User.Password = hashString
		User.Cookie = uuid.New().String()
		result := sql.Change_Password(User)
		if result {
			c.JSON(http.StatusOK, gin.H{"data": "密码修改成功"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"data": "密码修改失败"})
		}
	}
}
