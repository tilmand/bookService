package auth

import (
	"bookService/model"
	"bookService/store"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	AccessTokenTTL  = time.Hour * 1
	RefreshTokenTTL = time.Hour * 24 * 7
	StringsNumber   = 2
)

type BaseClaims struct {
	jwt.StandardClaims
	ID   uint64 `json:"id"`
	Role uint64 `json:"role"`
}

type AccessClaims struct {
	BaseClaims
	AccessUUID string `json:"access_uuid"`
}

type RefreshClaims struct {
	BaseClaims
	RefreshUUID string `json:"refresh_uuid"`
	UserID      uint64 `json:"user_id"`
}

type Tokens struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}

type Middleware struct {
	atKey *ecdsa.PrivateKey
	rtKey *ecdsa.PrivateKey
	mongo *store.MongoStore
}

type AuthMiddleware interface {
	Authorize(c *gin.Context)
	CreateTokens(id uint64) (*Tokens, error)
	Refresh(tokens Tokens) (*Tokens, error)
	ExtractToken(r *http.Request) string
	Validate(raw string) (*AccessClaims, error)
	GetUserID(token string) (uint64, error)
}

func NewAuthMiddleware(atKey, rtKey *ecdsa.PrivateKey, mongo *store.MongoStore) *Middleware {
	var middleware = &Middleware{
		atKey: atKey,
		rtKey: rtKey,
		mongo: mongo,
	}

	return middleware
}

func (m *Middleware) Authorize(c *gin.Context) {
	tokenString := m.ExtractToken(c.Request)
	claims, err := m.Validate(tokenString)
	if err != nil {
		log.Println("Authorize Validate err: ", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrUnauthorized)

		return
	}

	if claims == nil {
		log.Println("Authorize err: empty claims")
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrUnauthorized)

		return
	}

	_, err = m.mongo.UsersRepository.Find(claims.BaseClaims.ID)
	if err != nil {
		log.Println("Authorize Find err: ", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrUnauthorized)

		return
	}
}

func (m *Middleware) CreateTokens(id uint64) (*Tokens, error) {
	accessClaims, refreshClaims := GenerateClaims(id)

	at := jwt.NewWithClaims(jwt.SigningMethodES256, accessClaims)
	accessToken, err := at.SignedString(m.atKey)
	if err != nil {
		log.Println("CreateTokens SignedString atKey err: ", err)

		return nil, err
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims)
	refreshToken, err := rt.SignedString(m.rtKey)
	if err != nil {
		log.Println("CreateTokens SignedString rtKey err: ", err)

		return nil, err
	}

	return &Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (m *Middleware) Refresh(refreshToken string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				log.Printf("Refresh unexpected signing method: %v", token.Header["alg"])

				return nil, model.ErrUnauthorized
			}

			return &m.rtKey.PublicKey, nil
		})
	if err != nil {
		log.Println("Refresh ParseWithClaims err: ", err)

		return "", model.ErrUnauthorized
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok {
		log.Printf("Refresh invalid token claims: %v", token.Claims)

		return "", model.ErrUnauthorized
	}

	if !token.Valid {
		log.Println("Refresh not valid token err")

		return "", model.ErrUnauthorized
	}

	_, err = m.mongo.UsersRepository.Find(claims.BaseClaims.ID)
	if err != nil {
		log.Println("Refresh Find", err)

		return "", model.ErrInternalServerError

	}

	return m.GenerateAccessToken(claims.ID)
}

func (m *Middleware) GenerateAccessToken(id uint64) (string, error) {
	accessClaims, _ := GenerateClaims(id)

	at := jwt.NewWithClaims(jwt.SigningMethodES256, accessClaims)
	accessToken, err := at.SignedString(m.atKey)
	if err != nil {
		log.Println("GenerateAccessToken SignedString err: ", err)

		return "", err
	}

	return accessToken, nil
}

func GenerateClaims(id uint64) (*AccessClaims, *RefreshClaims) {
	access := AccessClaims{
		BaseClaims: NewClaims(id, AccessTokenTTL),
	}

	refresh := RefreshClaims{
		BaseClaims: NewClaims(id, RefreshTokenTTL),
		UserID:     id,
	}

	access.AccessUUID = refresh.Id
	refresh.RefreshUUID = access.Id

	return &access, &refresh
}

func NewClaims(id uint64, ttl time.Duration) BaseClaims {
	return BaseClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			Id:        uuid.NewV4().String(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "bookService",
		},
		ID: id,
	}
}

func (m *Middleware) RefreshTokenKey() *ecdsa.PrivateKey {
	return m.rtKey
}

func (m *Middleware) ExtractToken(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func (m *Middleware) Validate(raw string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(raw, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			log.Printf("Validate unexpected signing method: %v", token.Header["alg"])

			return nil, model.ErrUnauthorized
		}

		return &m.atKey.PublicKey, nil
	})
	if err != nil {
		log.Println("Validate ParseWithClaims err: ", err)

		return nil, model.ErrUnauthorized
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok {
		log.Printf("Validate invalid token claims: %v", token.Claims)

		return nil, model.ErrUnauthorized
	}

	if !token.Valid {
		log.Println("Validate not valid token err")

		return nil, model.ErrUnauthorized
	}

	return claims, nil
}

func GenerateECDSAPrivateKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("generateECDSAPrivateKey GenerateKey err: %v", err)

		return nil, err
	}

	return privateKey, nil
}
