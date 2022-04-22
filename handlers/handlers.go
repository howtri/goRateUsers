package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/howtri/goRateUsers/database"
)

func AddUserHandler(c *gin.Context) {
	userItem, statusCode, err := convertHTTPBodyToSkill(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(userItem.PassHash), bcrypt.MinCost)
	userItem.PassHash = string(hash)
	database.AddUser(userItem)
	c.JSON(statusCode, "")
}

func VerifyUserHandler(c *gin.Context) {
	userItem, statusCode, err := convertHTTPBodyToSkill(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	dbUser := database.GetUser(userItem.Username)
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PassHash), []byte(userItem.PassHash))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden, "")
	}
	c.JSON(statusCode, "")
}

func convertHTTPBodyToSkill(httpBody io.ReadCloser) (database.User, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return database.User{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToSkill(body)
}

func convertJSONBodyToSkill(jsonBody []byte) (database.User, int, error) {
	var skillItem database.User
	err := json.Unmarshal(jsonBody, &skillItem)
	if err != nil {
		return database.User{}, http.StatusBadRequest, err
	}
	return skillItem, http.StatusOK, nil
}
