package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/howtri/goRate/database"
	"github.com/howtri/goRate/skills"
)

// AddTodoHandler adds a new todo to the todo list
func AddSkillHandler(c *gin.Context) {
	skillItem, statusCode, err := convertHTTPBodyToSkill(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(statusCode, gin.H{"id": skills.AddSkill(skillItem)})
}

func GetSkillHandler(c *gin.Context) {
	skillID := c.Param("id")
	c.JSON(http.StatusAccepted, gin.H{"skill": database.GetSkill(skillID)})
}

func SearchSkillsHandler(c *gin.Context) {
	skillItem, statusCode, err := convertHTTPBodyToSkill(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"skill": database.SearchSkills(skillItem)})
}

func RankSkillHandler(c *gin.Context) {
	skillItem, statusCode, err := convertHTTPBodyToSkill(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	log.Printf("Calling rank")
	database.RankSkill(skillItem)
	log.Printf("Finished")
	c.JSON(http.StatusOK, "")
}

func convertHTTPBodyToSkill(httpBody io.ReadCloser) (database.Skill, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return database.Skill{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToSkill(body)
}

func convertJSONBodyToSkill(jsonBody []byte) (database.Skill, int, error) {
	var skillItem database.Skill
	err := json.Unmarshal(jsonBody, &skillItem)
	if err != nil {
		return database.Skill{}, http.StatusBadRequest, err
	}
	return skillItem, http.StatusOK, nil
}

func convertHTTPBodyToRanking(httpBody io.ReadCloser) (skills.Ranking, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return skills.Ranking{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToRanking(body)
}

func convertJSONBodyToRanking(jsonBody []byte) (skills.Ranking, int, error) {
	var rankingItem skills.Ranking
	err := json.Unmarshal(jsonBody, &rankingItem)
	if err != nil {
		return skills.Ranking{}, http.StatusBadRequest, err
	}
	return rankingItem, http.StatusOK, nil
}
