package util

import (
	"github.com/aoaostar/v8cdn_panel/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	Username   string `json:"username"`
	UserKey    string `json:"user_key"`
	UserApiKey string `json:"user_api_key"`
	jwt.StandardClaims
}

// 指定加密密钥
var jwtSecret = []byte(config.Conf.JwtSecret)

// GenerateToken 根据用户的用户名和密码产生token
func GenerateToken(username, userKey string, userApiKey string) (string, error) {
	//设置token有效时间
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour * 10000)

	claims := Claims{
		Username:   username,
		UserKey:    userKey,
		UserApiKey: userApiKey,
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
func ParseToken(token string) (*Claims, error) {

	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err

}
