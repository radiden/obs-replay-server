package services

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/andreykaipov/goobs"
	"github.com/charmbracelet/log"
	"github.com/radiden/obs-replay-server/models"
)

type Database struct {
	Client  *sql.DB
	Queries *models.Queries
}

type Services struct {
	OBS *goobs.Client
	DB  *Database
}

func InitServices(initialSchema string) *Services {
	db, err := initDB(initialSchema)
	if err != nil {
		log.Fatal("couldn't initialize db", "error", err)
	}

	obs, err := initOBS()
	if err != nil {
		log.Fatal("couldn't connect to obs", "error", err)
	}
	return &Services{
		DB:  db,
		OBS: obs,
	}
}

func initDB(initialSchema string) (*Database, error) {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "file:db.sqlite")
	if err != nil {
		return nil, err
	}

	db.ExecContext(ctx, initialSchema)

	return &Database{
		Client:  db,
		Queries: models.New(db),
	}, nil
}

func initOBS() (*goobs.Client, error) {
	obs, err := goobs.New("localhost:4455", goobs.WithPassword("ThpNdGXBwneYn9X7"))
	return obs, err
}
