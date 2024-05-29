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

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignUp godoc
// @Summary Register a new user
// @Description Create a new user account with a username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body Credentials true "User registration details"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /signup [post]
func SignUp(c *gin.Context) {
	var credentials Credentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(credentials.Password), 14)

	var exists bool
	db.DB.Model(&models.User{}).Select("count(*) > 0").Where("Username = ?", credentials.Username).Find(&exists)

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	user := models.User{Username: credentials.Username, PasswordHash: string(hashedPassword)}
	db.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// Login godoc
// @Summary Log in a user
// @Description Authenticate user with username and password and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body Credentials true "User credentials"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /login [post]
func Login(c *gin.Context) {
	var credentials Credentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	db.DB.Where("username = ?", credentials.Username).First(&user)

	if user.Id == 0 || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect username or password"})
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"userId":   user.Id,
		"exp":      expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(utils.JWTKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	user.LastSeen = time.Now()
	db.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
