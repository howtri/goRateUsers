package skills

import (
	"github.com/howtri/goRate/database"
	"github.com/rs/xid"
)

type Ranking struct {
	ID      string `json:"id"`
	Ranking int    `json:"ranking"`
}

func AddSkill(s database.Skill) string {
	s.ID = xid.New().String()
	database.AddSkill(s)
	return s.ID
}
