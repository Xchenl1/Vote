package tools

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	AccessTokenDuration  = 2 * time.Hour
	RefreshTokenDuration = 30 * 24 * time.Hour
	TokenIssuer          = "xinxuecheng-vote"
)

var Token VoteJwt

func NewToken(s string) {
	b := []byte("图书管理系统")
	fmt.Println(b)
	if s != "" {
		b = []byte(s)
	}

	Token = VoteJwt{Secret: b}
}

type VoteJwt struct {
	Secret []byte
}

// Claim 自定义的数据结构，这里使用了结构体的组合
type Claim struct {
	jwt.RegisteredClaims
	ID   int64  `json:"user_id"`
	Name string `json:"username"`
}

func (j *VoteJwt) getTime(t time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(t))
}

// 得到密钥
func (j *VoteJwt) keyFunc(token *jwt.Token) (interface{}, error) {
	return j.Secret, nil
}

// GetToken 颁发token access token 和 refresh token
func (j *VoteJwt) GetToken(id int64, name string) (aToken, rToken string, err error) {

	//这一行代码定义了一个 JWT 注册声明（RegisteredClaims）类型的变量 rc，
	//并初始化了其中的 ExpiresAt 和 Issuer 字段。ExpiresAt 表示访问令牌的过期时间，
	//使用 j.getTime(AccessTokenDuration) 函数计算得出，Issuer 表示颁发该令牌的发行者，这里为 TokenIssuer。
	rc := jwt.RegisteredClaims{
		ExpiresAt: j.getTime(AccessTokenDuration),
		Issuer:    TokenIssuer,
	}
	//这一行代码定义了一个自定义声明（Claim）类型的变量 claim，
	//并初始化了其中的 ID、Name 和 RegisteredClaims 字段。ID 和 Name 字段分别表示用户的唯一标识和用户名，
	//RegisteredClaims 字段则表示 JWT 的注册声明，这里使用之前定义的 rc。
	claim := Claim{
		ID:               id,
		Name:             name,
		RegisteredClaims: rc,
	}
	//这一行代码生成访问令牌，使用 jwt.NewWithClaims() 函数创建一个 JWT 令牌，
	//并指定签名方法为 HMAC-SHA256。claim 参数表示自定义声明，j.Secret 表示 JWT 的密钥，
	//SignedString() 方法将 JWT 令牌签名并生成字符串形式的令牌。
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(j.Secret)

	// refresh token 不需要保存任何用户信息
	rc.ExpiresAt = j.getTime(RefreshTokenDuration)
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, rc).SignedString(j.Secret)
	return
}

// VerifyToken 验证Token
func (j *VoteJwt) VerifyToken(tokenID string) (*Claim, error) {
	claim := &Claim{}
	token, err := jwt.ParseWithClaims(tokenID, claim, j.keyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("access token 验证失败")
	}

	return claim, nil
}

// RefreshToken 通过 refresh token 刷新 access token
func (j *VoteJwt) RefreshToken(a, r string) (aToken, rToken string, err error) {
	// r 无效直接返回
	if _, err = jwt.Parse(r, j.keyFunc); err != nil {
		return
	}
	// 从旧access token 中解析出claims数据
	claim := &Claim{}
	_, err = jwt.ParseWithClaims(a, claim, j.keyFunc)
	// 判断错误是不是因为access token 正常过期导致的
	if errors.Is(err, jwt.ErrTokenExpired) {
		return j.GetToken(claim.ID, claim.Name)
	}
	return
}
