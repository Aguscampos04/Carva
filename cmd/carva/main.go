package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	db, err := sqlite.Open(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()
	if err := sqlite.Migrate(ctx, db); err != nil {
		return err
	}

	projectRepository := sqlite.NewProjectRepository(db)
	sprintRepository := sqlite.NewSprintRepository(db)
	handler := httpapi.NewHandler(projects.NewService(projectRepository), sprints.NewService(sprintRepository))
	server := &http.Server{Addr: ":8080", Handler: handler}
	log.Printf("Carva API listening on %s", server.Addr)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
