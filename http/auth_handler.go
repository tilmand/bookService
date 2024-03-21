package http

import (
	"bookService/model"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

const salt = "book-service-v1"

type AuthHandlerInterface interface {
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)
	Refresh(c *gin.Context)
	Recover(c *gin.Context)
	SetNewPassword(c *gin.Context)
}

type AuthHandler struct {
	api *api
}

func NewAuthHandler(a *api) *AuthHandler {
	return &AuthHandler{
		api: a,
	}
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	creds := &model.User{}
	err := c.ShouldBindJSON(&creds)
	if err != nil {
		log.Println("SignIn ShouldBindJSON err: ", err)
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	login := strings.TrimSpace(creds.Login)
	pass := strings.TrimSpace(creds.Password)
	if login == "" || pass == "" {
		log.Println("SignIn Empty login or pass")
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	user, err := h.api.mongo.UsersRepository.GetByLogin(creds.Login)
	if err != nil {
		log.Println("SignIn GetByLogin err: ", err)
		if err.Error() == "not found" {
			c.JSON(http.StatusUnauthorized, model.ErrUnauthorized)
		} else {
			c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)
		}

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		log.Println("SignIn CompareHashAndPassword err: ", err)
		c.JSON(http.StatusUnauthorized, model.ErrUnauthorized)

		return
	}

	tokens, err := h.api.auth.CreateTokens(user.ID)
	if err != nil {
		log.Println("SignIn CreateTokens err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}

	answer := map[string]interface{}{
		"accessToken":  tokens.Access,
		"refreshToken": tokens.Refresh,
	}

	c.JSON(http.StatusOK, answer)
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	user := &model.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		log.Println("SignUp ShouldBindJSON err: ", err)
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	if user.Login == "" || user.Password == "" {
		log.Println("SignIn Empty login or pass")
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("SignUp GenerateFromPassword err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}
	user.Password = string(hashedPassword)

	if err := h.api.mongo.UsersRepository.Insert(*user); err != nil {
		log.Println("SignUp Insert err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user created successfully"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken := c.PostForm("refreshToken")
	if refreshToken == "" {
		log.Println("Refresh empty refreshToken")
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	newAccessToken, err := h.api.auth.Refresh(refreshToken)
	if err != nil {
		log.Println("Refresh Refresh err: ", err)
		c.JSON(http.StatusUnauthorized, model.ErrUnauthorized)

		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": newAccessToken})
}

func (h *AuthHandler) Recover(c *gin.Context) {
	var emailRequest struct {
		Email string `json:"email"`
	}
	err := c.ShouldBindJSON(&emailRequest)
	if err != nil {
		log.Println("Recover ShouldBindJSON err: ", err)
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	user, err := h.api.mongo.UsersRepository.GetByLogin(emailRequest.Email)
	if err != nil {
		log.Println("Recover GetByLogin err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})

		return
	}

	recoveryToken, err := generateRecoveryToken()
	err = h.api.mongo.UsersRepository.SaveRecoveryToken(user.ID, recoveryToken)
	if err != nil {
		log.Println("Recover SaveRecoveryToken err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}

	go func() {
		err := sendRecoveryEmail(emailRequest.Email, recoveryToken)
		if err != nil {
			log.Println("Recover sendRecoveryEmail err: ", err)

			return
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "recovery email sent"})
}

func (h *AuthHandler) SetNewPassword(c *gin.Context) {
	recoveryToken := c.Param("token")

	userID, err := h.api.mongo.UsersRepository.VerifyRecoveryToken(recoveryToken)
	if err != nil {
		log.Println("SetNewPassword VerifyRecoveryToken err:", err)
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	var newPasswordRequest struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&newPasswordRequest); err != nil {
		log.Println("SetNewPassword ShouldBindJSON err: ", err)
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPasswordRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("SetNewPassword GenerateFromPassword err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}

	err = h.api.mongo.UsersRepository.SetPassword(userID, string(hashedPassword))
	if err != nil {
		log.Println("SetNewPassword SetPassword err: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrInternalServerError)

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

func IsPasswordMatch(password, hashedPassword string) bool {
	userPasswordHash := H3hash(password + salt)

	return userPasswordHash == hashedPassword
}

func H3hash(s string) string {
	h3 := sha3.New512()
	if _, err := io.WriteString(h3, s); err != nil {
		log.Printf("H3hash WriteString err: %v", err)
	}

	return fmt.Sprintf("%x", h3.Sum(nil))
}

func generateRecoveryToken() (string, error) {
	tokenLength := 32

	tokenBytes := make([]byte, tokenLength)

	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	tokenHex := hex.EncodeToString(tokenBytes)

	return tokenHex, nil
}

func sendRecoveryEmail(email, recoveryToken string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := "testSendEmail456@gmail.com"
	smtpPassword := "ulnsooneferjalvw"

	subject := "Password Recovery"
	body := "Dear User,\n\nPlease click on the following link to reset your password:\n\n"
	body += "http://localhost:8080/api/v1/recoverPassword?token=" + recoveryToken + "\n\n"
	body += "Best regards,\nBook Service Team"

	msg := "From: Book Service <" + smtpUsername + ">\n"
	msg += "To: " + email + "\n"
	msg += "Subject: " + subject + "\n\n"
	msg += body

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err := smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, smtpUsername, []string{email}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
