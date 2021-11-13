package util

import (
	"errors"
	"github.com/aoaostar/v8cdn_panel/pkg"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type User struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	UserKey    string `json:"user_key"`
	UserApiKey string `json:"user_api_key"`
	//认证类型，partner or user_api_key
	AuthType string `json:"auth_type"`
	jwt.StandardClaims
}

// 指定加密密钥
var jwtSecret = []byte(pkg.Conf.JwtSecret)

// GenerateToken 根据用户的用户名和密码产生token
func GenerateToken(user User) (string, error) {
	//设置token有效时间
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour * 7)

	claims := User{
		ID:      user.ID,
		Email:      user.Email,
		UserKey:    user.UserKey,
		UserApiKey: user.UserApiKey,
		AuthType:   user.AuthType,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 指定token发行人
			Issuer: "v8cdn",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseToken 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func ParseToken(token string) (*User, error) {

	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &User{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*User); ok && tokenClaims.Valid {
			if claims.Email == "" || claims.AuthType == "" || claims.UserApiKey == "" {
				return nil, errors.New("token无效")
			}
			return claims, nil
		}
	}
	return nil, err

}
