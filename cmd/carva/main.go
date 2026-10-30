package main

import (
	"context"
	appeffort "github.com/Aguscampos04/Carva/internal/application/effort"
	"log"
	"net/http"
	"time"

	"github.com/Aguscampos04/Carva/internal/application/backlog"
	"github.com/Aguscampos04/Carva/internal/application/defects"
	"github.com/Aguscampos04/Carva/internal/application/members"
	appmetrics "github.com/Aguscampos04/Carva/internal/application/metrics"
	"github.com/Aguscampos04/Carva/internal/application/planningpoker"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/reports"
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
	api := httpapi.NewHandlerWithPlanningPoker(projectService, backlogService, sprintService, planningPokerService, memberService)
	api = httpapi.WithEffort(api, appeffort.NewService(sqlite.NewEffortRepository(db), time.Now))
	api = httpapi.WithDefects(api, defects.NewService(sqlite.NewDefectRepository(db)))
	api = httpapi.WithMetrics(api, appmetrics.NewService(sqlite.NewMetricsRepository(db)))
	api = httpapi.WithReports(api, reports.NewService(sqlite.NewReportSource(db), appmetrics.NewService(sqlite.NewMetricsRepository(db))))
	api = httpapi.WithAssignees(api, backlog.NewAssigneeService(sqlite.NewAssigneeRepository(db)))
	mux.Handle("/api/", api)
	mux.Handle("/", http.FileServer(http.FS(static.Files)))

	log.Print("Carva listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
