package handlers

import (
	"context"
	"net/http"
	"time"

	"auth/config"
	"auth/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-contrib/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ShowSignup(c *gin.Context) {

	c.HTML(200, "signup.html", nil)

}

func ShowLogin(c *gin.Context) {

	c.HTML(200, "login.html", nil)

}

func Signup(c *gin.Context) {

	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	if name == "" || email == "" || password == "" {
		c.String(http.StatusBadRequest, "All fields are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User

	err := config.UserCollection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&existingUser)

	if err != nil && err != mongo.ErrNoDocuments {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}

	if err == nil {
		c.String(http.StatusBadRequest, "Email already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.String(http.StatusInternalServerError, "Could not hash password")
		return
	}

	user := models.User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	_, err = config.UserCollection.InsertOne(ctx, user)

	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func Login(c *gin.Context) {

	email := c.PostForm("email")
	password := c.PostForm("password")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	err := config.UserCollection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&user)

	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid email or password")
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid email or password")
		return
	}

	session := sessions.Default(c)

	session.Set("user_id", user.ID.Hex())
	session.Set("user_name", user.Name)

	err = session.Save()

	if err != nil {
		c.String(http.StatusInternalServerError, "Could not create session")
		return
	}

	c.Redirect(http.StatusSeeOther, "/dashboard")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}

func Dashboard(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("user_name")
	c.HTML(200, "dashboard.html", gin.H{
		"userName": userName,
	})
}
