package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"server/config"
	redisInit "server/setup/redis"
	"strconv"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

const (
	// TokenTypeAccess 用于普通接口鉴权，有效期通常较短。
	TokenTypeAccess = "access"
	// TokenTypeRefresh 只用于刷新 token，有效期通常比 access token 长。
	TokenTypeRefresh = "refresh"
)

var (
	ErrEmptyJwtSecret      = errors.New("jwt secret is empty")
	ErrInvalidTokenType    = errors.New("invalid token type")
	ErrRedisNotInitialized = errors.New("redis client is not initialized")
	ErrTokenRevoked        = errors.New("token revoked or not found")
	ErrTokenExpired        = errors.New("token expired")
)

// CustomClaims 是项目自定义 JWT 载荷。
// 这里只放必要身份信息，不放密码、权限明细等敏感或过大的数据。
type CustomClaims struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
	JTI       string `json:"jti"`
	gojwt.RegisteredClaims
}

// TokenPair 是登录成功后返回给前端的一组 token。
// access_token 用于访问接口，refresh_token 用于 access_token 过期后的续签。
// 两个 JTI 字段只在服务端内部使用（在线用户会话记录），不下发给前端。
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	AccessJTI    string `json:"-"`
	RefreshJTI   string `json:"-"`
}

// GenerateToken 签发一组新的 access token 和 refresh token，并把 token 状态写入 Redis。
// Redis 状态用于支持主动失效、单一登录、多端登录等能力。
func GenerateToken(userID uint, username string) (*TokenPair, error) {
	jwtConfig := config.JwtConfig()
	if jwtConfig.Secret == "" {
		return nil, ErrEmptyJwtSecret
	}

	// access token 用于接口鉴权，refresh token 只用于换取新 token。
	accessExpires := time.Duration(jwtConfig.AccessExpiresMinute) * time.Minute
	refreshExpires := time.Duration(jwtConfig.RefreshExpiresHour) * time.Hour

	accessToken, accessJTI, err := signToken(userID, username, TokenTypeAccess, accessExpires, jwtConfig)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshJTI, err := signToken(userID, username, TokenTypeRefresh, refreshExpires, jwtConfig)
	if err != nil {
		return nil, err
	}

	if err := storeTokenState(userID, TokenTypeAccess, accessJTI, accessExpires, jwtConfig.LoginMode); err != nil {
		return nil, err
	}
	if err := storeTokenState(userID, TokenTypeRefresh, refreshJTI, refreshExpires, jwtConfig.LoginMode); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		// expires_in 只表示 access token 的剩余有效秒数，方便前端做刷新判断。
		ExpiresIn:  int64(accessExpires.Seconds()),
		AccessJTI:  accessJTI,
		RefreshJTI: refreshJTI,
	}, nil
}

// ParseToken 只校验 JWT 本身：签名算法、签名密钥、过期时间等。
// 它不会检查 Redis 中 token 是否仍然有效；需要完整鉴权时请使用 ValidateToken。
func ParseToken(tokenString string) (*CustomClaims, error) {
	jwtConfig := config.JwtConfig()
	if jwtConfig.Secret == "" {
		return nil, ErrEmptyJwtSecret
	}

	// 只接受 HMAC 签名算法，避免被伪造 alg 绕过校验。
	token, err := gojwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *gojwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtConfig.Secret), nil
	})
	if err != nil {
		if errors.Is(err, gojwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateToken 是接口鉴权应使用的方法。
// 它会先解析 JWT，再校验 token_type，最后确认 Redis 中仍存在对应 jti。
func ValidateToken(tokenString string, tokenType string) (*CustomClaims, error) {
	if tokenType != TokenTypeAccess && tokenType != TokenTypeRefresh {
		return nil, ErrInvalidTokenType
	}

	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != tokenType {
		return nil, ErrInvalidTokenType
	}

	if err := validateTokenState(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// RefreshToken 使用 refresh token 换取一组新的 token。
// 注意：传入 access token 会因为 token_type 不匹配而失败。
func RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := ValidateToken(refreshToken, TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	return GenerateToken(claims.UserID, claims.Username)
}

// RevokeToken 撤销一个已签发 token。
// 撤销后 Redis 中对应 jti 会被删除，后续 ValidateToken 会失败。
func RevokeToken(tokenString string) error {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return err
	}

	client := redisInit.Redis.Get()
	if client == nil {
		return ErrRedisNotInitialized
	}

	ctx := context.Background()
	if err := client.Del(ctx, tokenKey(claims.TokenType, claims.JTI)).Err(); err != nil {
		return err
	}

	if config.JwtConfig().LoginMode == config.JwtLoginModeSingle {
		if err := client.Del(ctx, userTokenKey(claims.UserID, claims.TokenType)).Err(); err != nil {
			return err
		}
	}

	return nil
}

// signToken 负责真正生成 JWT 字符串。
// jti 是每个 token 的唯一编号，会同时写入 JWT 和 Redis。
func signToken(userID uint, username string, tokenType string, expires time.Duration, jwtConfig config.Jwt) (string, string, error) {
	jti, err := newJTI()
	if err != nil {
		return "", "", err
	}

	now := time.Now()
	claims := CustomClaims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		JTI:       jti,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(now.Add(expires)),
			IssuedAt:  gojwt.NewNumericDate(now),
			NotBefore: gojwt.NewNumericDate(now),
			Issuer:    jwtConfig.Issuer,
			Subject:   strconv.FormatUint(uint64(userID), 10),
			ID:        jti,
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", "", err
	}

	return signedToken, jti, nil
}

// storeTokenState 把 token 的 jti 写入 Redis。
// 多端登录：只记录 jwt:{type}:{jti}，每个 token 独立有效。
// 单一登录：额外写 jwt:user:{userID}:{type}，新 token 会覆盖旧 token 的当前 jti。
func storeTokenState(userID uint, tokenType string, jti string, expires time.Duration, loginMode string) error {
	client := redisInit.Redis.Get()
	if client == nil {
		return ErrRedisNotInitialized
	}

	ctx := context.Background()
	if err := client.Set(ctx, tokenKey(tokenType, jti), strconv.FormatUint(uint64(userID), 10), expires).Err(); err != nil {
		return err
	}

	if loginMode == config.JwtLoginModeSingle {
		if err := client.Set(ctx, userTokenKey(userID, tokenType), jti, expires).Err(); err != nil {
			return err
		}
	}

	return nil
}

// validateTokenState 校验 Redis 中的 token 状态。
// 如果 Redis 中不存在该 jti，说明 token 已过期、被撤销或从未签发过。
func validateTokenState(claims *CustomClaims) error {
	client := redisInit.Redis.Get()
	if client == nil {
		return ErrRedisNotInitialized
	}

	ctx := context.Background()
	exists, err := client.Exists(ctx, tokenKey(claims.TokenType, claims.JTI)).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return ErrTokenRevoked
	}

	if config.JwtConfig().LoginMode == config.JwtLoginModeSingle {
		// 单一登录模式下，用户当前有效 jti 必须和 token 中的 jti 一致。
		currentJTI, err := client.Get(ctx, userTokenKey(claims.UserID, claims.TokenType)).Result()
		if err != nil {
			return err
		}
		if currentJTI != claims.JTI {
			return ErrTokenRevoked
		}
	}

	return nil
}

// TokenStateKey 返回 token 状态在 Redis 中的 key。
// 导出给在线用户管理使用（踢下线需要按 jti 删除对应 key）。
func TokenStateKey(tokenType string, jti string) string {
	return fmt.Sprintf("jwt:%s:%s", tokenType, jti)
}

// tokenKey 是 TokenStateKey 的包内简写。
func tokenKey(tokenType string, jti string) string {
	return TokenStateKey(tokenType, jti)
}

// userTokenKey 是单一登录模式下记录用户当前有效 token 的 Redis key。
func userTokenKey(userID uint, tokenType string) string {
	return fmt.Sprintf("jwt:user:%d:%s", userID, tokenType)
}

// newJTI 生成随机 token ID。
// 这里使用 crypto/rand，避免可预测的 token 编号。
func newJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
