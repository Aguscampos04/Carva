package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Aguscampos04/Carva/internal/application/backlog"
	"github.com/Aguscampos04/Carva/internal/application/members"
	"github.com/Aguscampos04/Carva/internal/application/planningpoker"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
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

	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	backlogRepository := sqlite.NewBacklogRepository(db)
	backlogService := backlog.NewService(backlogRepository)
	sprintService := sprints.NewService(sqlite.NewSprintRepository(db), backlogRepository)
	mux := http.NewServeMux()
	memberService := members.NewService(sqlite.NewMemberRepository(db))
	planningPokerService := planningpoker.NewService(sqlite.NewPlanningPokerRepository(db), backlogRepository, memberService)
	mux.Handle("/api/", httpapi.NewHandlerWithPlanningPoker(projectService, backlogService, sprintService, planningPokerService, memberService))
	mux.Handle("/", http.FileServer(http.FS(static.Files)))

	log.Print("Carva listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
