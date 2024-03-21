package auth

import (
	"crypto/ecdsa"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateTokens(t *testing.T) {
	atKey, _ := GenerateECDSAPrivateKey()
	rtKey, _ := GenerateECDSAPrivateKey()
	middleware := NewAuthMiddleware(atKey, rtKey, nil)

	tokens, err := middleware.CreateTokens(123)
	assert.NoError(t, err)
	assert.NotNil(t, tokens.Access)
	assert.NotNil(t, tokens.Refresh)
}

func TestExtractToken(t *testing.T) {
	atKey, _ := GenerateECDSAPrivateKey()
	rtKey, _ := GenerateECDSAPrivateKey()
	middleware := NewAuthMiddleware(atKey, rtKey, nil)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer fake_token")

	token := middleware.ExtractToken(req)
	assert.Equal(t, "Bearer fake_token", token)
}

func TestValidateToken(t *testing.T) {
	atKey, _ := GenerateECDSAPrivateKey()
	rtKey, _ := GenerateECDSAPrivateKey()
	middleware := NewAuthMiddleware(atKey, rtKey, nil)

	accessClaims, _ := GenerateClaims(123)
	accessToken, _ := GenerateToken(atKey, accessClaims)

	claims, err := middleware.Validate(accessToken)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
}

func TestAuthorize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer fake_token")
		c.Next()
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authorized"})
	})

	atKey, _ := GenerateECDSAPrivateKey()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	accessToken, _ := GenerateToken(atKey, &AccessClaims{BaseClaims: NewClaims(123, AccessTokenTTL)})
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func GenerateToken(key *ecdsa.PrivateKey, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(key)
}
