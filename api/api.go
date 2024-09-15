package api

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/m-row/finder-example/controller"
	"github.com/m-row/finder-example/controllers"
	"github.com/m-row/finder-example/database"
	"github.com/m-row/finder-example/models"
)

var (
	CommitCount    = "0"
	CommitDescribe = "dev"
	Version        = "1." + CommitCount + "." + CommitDescribe
)

type Application struct {
	Controllers *controllers.Controllers
	DB          *sqlx.DB
	Models      *models.Models
}

func NewAPI() *Application {
	db, err := database.OpenSQLX()
	if err != nil {
		log.Fatalf("couldn't open db: %s", err.Error())
	}
	log.Println("database connection pool established")
	tc := make(map[string][]string)

	if err := database.PGInfo(db, tc); err != nil {
		log.Fatalf("couldn't get pgInfo: %s", err.Error())
	}
	log.Println("database table and column info saved")

	m := models.Setup(db, tc)

	deps := &controller.Dependencies{
		Models: m,
	}
	ctrls := controllers.Setup(deps)

	newApi := &Application{
		Controllers: ctrls,
		DB:          db,
		Models:      m,
	}
	return newApi
}
