package http

import (
	"bookService/auth"
	"bookService/mocks"
	"bookService/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type FakeAPI struct {
	db   *FakeMongoStore
	auth auth.Middleware
}

func NewFakeAPI() *FakeAPI {
	return &FakeAPI{
		db:   &FakeMongoStore{},
		auth: auth.Middleware{},
	}
}

type FakeMongoStore struct{}

func (f *FakeMongoStore) GetByLogin(login string) (*model.User, error) {
	return &model.User{}, nil
}

func (f *FakeMongoStore) Insert(user model.User) error {
	return nil
}

func (f *FakeMongoStore) SaveRecoveryToken(userID uint64, token string) error {
	return nil
}

func (f *FakeMongoStore) VerifyRecoveryToken(token string) (uint64, error) {
	return 123, nil
}

type FakeJWTAuth struct{}

func (f *FakeJWTAuth) CreateTokens(userID uint64) (*auth.Tokens, error) {
	return &auth.Tokens{
		Access:  "fake_access_token",
		Refresh: "fake_refresh_token",
	}, nil
}

func (f *FakeJWTAuth) ParseToken(tokenString string) (uint64, error) {
	return 123, nil
}

func ConvertToAPI(fakeAPI *FakeAPI) *api {
	return &api{}
}

func TestSignInHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthHandler := mocks.NewMockAuthHandlerInterface(ctrl)

	mockAuthHandler.EXPECT().SignIn(gomock.Any()).Return()

	requestData := map[string]string{
		"Login":    "testuser",
		"Password": "testpassword",
	}
	payload, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %s", err)
	}

	req, err := http.NewRequest("POST", "/signin", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := gin.Default()

	router.POST("/signin", func(c *gin.Context) { mockAuthHandler.SignIn(c) })
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSignUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthHandler := mocks.NewMockAuthHandlerInterface(ctrl)

	mockAuthHandler.EXPECT().SignUp(gomock.Any()).Return()

	requestData := map[string]string{
		"Login":    "testuser",
		"Password": "testpassword",
	}
	payload, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %s", err)
	}

	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := gin.Default()

	router.POST("/signup", func(c *gin.Context) { mockAuthHandler.SignUp(c) })
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
func TestRefreshHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthHandler := mocks.NewMockAuthHandlerInterface(ctrl)

	mockAuthHandler.EXPECT().Refresh(gomock.Any()).Return()

	req, err := http.NewRequest("POST", "/refresh", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := gin.Default()

	router.POST("/refresh", func(c *gin.Context) { mockAuthHandler.Refresh(c) })
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRecoverHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthHandler := mocks.NewMockAuthHandlerInterface(ctrl)

	mockAuthHandler.EXPECT().Recover(gomock.Any()).Return()

	requestData := map[string]string{
		"email": "test@example.com",
	}
	payload, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %s", err)
	}

	req, err := http.NewRequest("POST", "/recover", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := gin.Default()

	router.POST("/recover", func(c *gin.Context) { mockAuthHandler.Recover(c) })
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSetNewPasswordHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthHandler := mocks.NewMockAuthHandlerInterface(ctrl)

	mockAuthHandler.EXPECT().SetNewPassword(gomock.Any()).Return()

	requestData := map[string]string{
		"password": "new_password",
	}
	payload, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %s", err)
	}

	req, err := http.NewRequest("POST", "/set-new-password/token123", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := gin.Default()

	router.POST("/set-new-password/:token", func(c *gin.Context) { mockAuthHandler.SetNewPassword(c) })
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func SignInWrapper(handler auth.Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sign In is not implemented"})
	}
}
