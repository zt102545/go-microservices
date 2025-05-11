package utils

import (
	"crypto/md5"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func Md5(input string) string {
	data := []byte(input)
	hash := md5.Sum(data)
	output := hex.EncodeToString(hash[:])
	return output
}

// CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
type CustomClaims struct {
	// 可根据需要自行添加字段
	UserID               int64 `json:"user_id"`
	jwt.RegisteredClaims       // 内嵌标准的声明
}

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init1() {
	var err error
	var bytes []byte
	bytes, err = os.ReadFile("/root/uccs/realworld/private.pem")
	if err != nil {
		panic(err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		panic(err)
	}

	bytes, err = os.ReadFile("/root/uccs/realworld/public.pem")
	if err != nil {
		panic(err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		panic(err)
	}
}

// GenToken 生成JWT
func GenToken(userId int64, secretKey string) (string, error) {
	// 创建一个我们自己声明的数据
	claims := CustomClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 定义过期时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串
	return token.SignedString([]byte(secretKey))
}

// ParseToken 解析JWT
func ParseToken(tokenString string, secretKey string) (*CustomClaims, error) {
	// 解析token
	var mc = new(CustomClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
