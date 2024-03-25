package handlers

import (
	"net/http"
	"time"

	"github.com/dgb35/telemogus_backend/db"
	"github.com/dgb35/telemogus_backend/models"
	"github.com/dgb35/telemogus_backend/utils"
	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	user := models.User{Username: input.Username, PasswordHash: string(hashedPassword)}
	db.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	db.DB.Where("username = ?", credentials.Username).First(&user)

	if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect username or password"})
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"userId":   user.ID,
		"exp":      expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(utils.JWTKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}