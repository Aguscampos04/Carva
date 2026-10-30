package reports_test

import (
 "context"
 "errors"
 "testing"
 app "github.com/Aguscampos04/Carva/internal/application/reports"
 metricsapp "github.com/Aguscampos04/Carva/internal/application/metrics"
 "github.com/Aguscampos04/Carva/internal/domain/backlog"
 "github.com/Aguscampos04/Carva/internal/domain/defect"
 "github.com/Aguscampos04/Carva/internal/domain/effort"
 "github.com/Aguscampos04/Carva/internal/domain/project"
 "github.com/Aguscampos04/Carva/internal/domain/sprint"
)

type source struct { data app.Data; err error }
func(s source) Load(context.Context,string)(app.Data,error){return s.data,s.err}
type metricSource struct { *metricsapp.Service }

func TestReportPropagatesSourceFailureAndRejectsInvalidSprint(t *testing.T){
 sentinel:=errors.New("read failed")
 service:=app.NewService(source{err:sentinel},metricSource{})
 if _,err:=service.Project(context.Background(),"p"); !errors.Is(err,sentinel){t.Fatal(err)}
 if _,err:=service.Sprint(context.Background(),"p",0);!errors.Is(err,app.ErrInvalidSprintID){t.Fatal(err)}
}

func TestSprintScopePreservesSnapshotAndFiltersRecords(t *testing.T){
 one,two:=int64(1),int64(2)
 original:=app.Report{Scope:"project",Project:project.Project{ID:"p"},
 Stories:[]backlog.Item{{ID:"story",Title:"Changed"}},
 Sprints:[]app.SprintReport{{Sprint:sprint.Sprint{ID:1},Stories:[]sprint.PlannedBacklogItem{{BacklogItemID:"story",Title:"Original",Status:"Pendiente"}}}},
 Effort:[]effort.Entry{{Input:effort.Input{SprintID:&one}},{Input:effort.Input{SprintID:&two}},{Input:effort.Input{}}},
 Defects:[]defect.Defect{{ID:"a",DetectionSprintID:&one,ResolutionSprintID:&two},{ID:"b",DetectionSprintID:&one,ResolutionSprintID:&one},{ID:"c"}},
 }
 report,err:=app.ForSprint(original,1);if err!=nil{t.Fatal(err)}
 if report.Scope!="sprint"||len(report.Stories)!=0||len(report.Effort)!=1||len(report.Defects)!=2{t.Fatalf("scope=%+v",report)}
 if report.Sprints[0].Stories[0].Title!="Original"||report.Sprints[0].Stories[0].Status!="Pendiente"{t.Fatal("snapshot changed")}
 if len(original.Stories)!=1||len(original.Effort)!=3||len(original.Defects)!=3{t.Fatal("source mutated")}
 if _,err:=app.ForSprint(original,999);!errors.Is(err,sprint.ErrNotFound){t.Fatal(err)}
}
