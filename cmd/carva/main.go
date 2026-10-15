package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Aguscampos04/Carva/internal/application/backlog"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
	"github.com/Aguscampos04/Carva/web/static"
)

func main() {
	ctx := context.Background()
	db, err := sqlite.Open(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := sqlite.Migrate(ctx, db); err != nil {
		log.Fatal(err)
	}

	service := projects.NewService(sqlite.NewProjectRepository(db))
	backlogService := backlog.NewService(sqlite.NewBacklogRepository(db))
	mux := http.NewServeMux()
	mux.Handle("/api/", httpapi.NewHandler(service, backlogService))
	mux.Handle("/", http.FileServer(http.FS(static.Files)))

	log.Print("Carva listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
