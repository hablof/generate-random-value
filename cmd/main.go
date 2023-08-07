package main

import (
	"github.com/gofiber/fiber/v2/log"

	"github.com/hablof/generate-random-value/internal/controller"
	"github.com/hablof/generate-random-value/internal/db"
	"github.com/hablof/generate-random-value/internal/repository"
	"github.com/hablof/generate-random-value/internal/service"
)

func main() {
	log.SetLevel(log.LevelDebug)

	db, err := db.NewDB()
	if err != nil {
		log.Error(err)
		return
	}

	r := repository.NewRepository(db)
	g := service.Generator{}
	app := controller.NewServer(r, &g)

	log.Fatal(app.Listen(":3000"))
}
