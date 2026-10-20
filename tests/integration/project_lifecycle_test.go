package integration_test

import (
 "context"
 "database/sql"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "reflect"
 "strings"
 "testing"
 "time"

 "github.com/Aguscampos04/Carva/internal/application/backlog"
 "github.com/Aguscampos04/Carva/internal/application/members"
 "github.com/Aguscampos04/Carva/internal/application/projects"
 "github.com/Aguscampos04/Carva/internal/application/sprints"
 "github.com/Aguscampos04/Carva/internal/persistence/sqlite"
 "github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func lifecycleHandler(db *sql.DB, instant string) http.Handler {
 clock:=func()time.Time{v,err:=time.Parse(time.RFC3339,instant);if err!=nil{panic(err)};return v}
 br:=sqlite.NewBacklogRepository(db)
 return httpapi.NewHandler(projects.NewServiceWithClock(sqlite.NewProjectRepository(db),clock),backlog.NewService(br),sprints.NewService(sqlite.NewSprintRepository(db),br),members.NewService(sqlite.NewMemberRepository(db)))
}
func lifeRequest(t *testing.T,h http.Handler,method,path,body string,status int) map[string]any {
 t.Helper();r:=httptest.NewRecorder();h.ServeHTTP(r,httptest.NewRequest(method,path,strings.NewReader(body)))
 if r.Code!=status{t.Fatalf("%s %s => %d want %d: %s",method,path,r.Code,status,r.Body)}
 var result map[string]any;if err:=json.Unmarshal(r.Body.Bytes(),&result);err!=nil{t.Fatal(err)};return result
}

// US-02 BDD-01/02/03/07/09/10/12/13/14/15/18/19/20.
func TestProjectLifecyclePersistenceAndCompatibility(t *testing.T){
 path,db:=openMigratedProjectDB(t);h:=lifecycleHandler(db,"2026-10-21T02:30:00Z")
 p:=lifeRequest(t,h,"POST","/api/projects",`{"name":"Carva Demo"}`,201);id:=p["id"].(string);url:="/api/projects/"+id
 if p["start_date"]!=nil || p["end_date"]!=nil{t.Fatal(p)}
 lifeRequest(t,h,"PATCH",url,`{"start_date":"2026-10-01"}`,200)
 p=lifeRequest(t,h,"PATCH",url,`{"name":"  Carva II  "}`,200)
 if p["start_date"]!="2026-10-01" || p["status"]!="Planificado" || p["name"]!="  Carva II  "{t.Fatal(p)}
 sprint:=lifeRequest(t,h,"POST",url+"/sprints",`{"goal":"MVP"}`,201)
 itemBody:=`{"title":"PDF","description":"Informe","priority":"Alta","acceptance_criteria":["Descargable"]}`
 item:=lifeRequest(t,h,"POST",url+"/backlog",itemBody,201)
 sprintURL:=url+"/sprints/1"
 lifeRequest(t,h,"POST",sprintURL+"/backlog-items",`{"backlog_item_id":"`+item["id"].(string)+`"}`,201)
 before:=performJSON(t,h,"GET",sprintURL+"/backlog-items","").Body.String()
 final:=lifeRequest(t,h,"POST",url+"/finalize",`{"confirmed":true}`,200)
 if final["status"]!="Finalizado" || final["end_date"]!="2026-10-20"{t.Fatal(final)}
 after:=lifeRequest(t,h,"GET",sprintURL,"",200);if !reflect.DeepEqual(sprint,after){t.Fatalf("changed sprint: %v %v",sprint,after)}
 if got:=performJSON(t,h,"GET",sprintURL+"/backlog-items","").Body.String();got!=before{t.Fatal("changed snapshot")}
 lifeRequest(t,h,"POST",url+"/sprints",`{"goal":"Forbidden"}`,409)
 lifeRequest(t,h,"POST",url+"/members",`{"name":"Ana"}`,201)
 next:=lifeRequest(t,h,"POST",url+"/backlog",itemBody,201)
 lifeRequest(t,h,"PUT",url+"/backlog/"+next["id"].(string),strings.Replace(itemBody,"PDF","PDF II",1),200)
 lifeRequest(t,h,"PATCH",url,`{"start_date":"2026-10-21","name":"Not saved"}`,400)
 p=lifeRequest(t,h,"PATCH",url,`{"start_date":"2026-10-20"}`,200)
 p=lifeRequest(t,h,"PATCH",url,`{"start_date":null}`,200)
 if p["name"]!="  Carva II  " || p["start_date"]!=nil{t.Fatal(p)}
 h=lifecycleHandler(db,"2026-10-22T15:00:00Z");p=lifeRequest(t,h,"POST",url+"/finalize",`{"confirmed":true}`,200)
 if p["end_date"]!="2026-10-20"{t.Fatal(p)}
 if err:=db.Close();err!=nil{t.Fatal(err)}
 reopened,err:=sqlite.OpenPath(context.Background(),path);if err!=nil{t.Fatal(err)};defer reopened.Close()
 got:=lifeRequest(t,lifecycleHandler(reopened,"2026-10-22T15:00:00Z"),"GET",url,"",200)
 if !reflect.DeepEqual(p,got){t.Fatalf("reopen: %v %v",p,got)}
}

// US-02 BDD-04/05/16/17/23/26; whole object and atomic rejection.
func TestProjectLifecycleRejectsInvalidRequests(t *testing.T){
 _,db:=openMigratedProjectDB(t);h:=lifecycleHandler(db,"2026-10-20T15:00:00Z")
 p:=lifeRequest(t,h,"POST","/api/projects",`{"name":"Demo"}`,201);url:="/api/projects/"+p["id"].(string)
 for _,body:=range []string{"",`null`,`[]`,`"x"`,`{}`,`{"name":null}`,`{"name":false}`,`{"name":" "}`,`{"Name":"X"}`,`{"name":"X","extra":true}`,`{"name":"A","name":"B"}`,`{"name":"A"} {}`,`{"start_date":true}`,`{"start_date":""}`,`{"start_date":"2026-02-29","name":"X"}`,`{"start_date":"2026-04-31"}`,`{"start_date":"2026-10-20T00:00:00Z"}`,`{"start_date":"2026-1-2"}`,`{"start_date":"0000-01-01"}`,`{"status":"Finalizado"}`,`{"end_date":null}`,`{"id":"else"}`,"{\"name\":\"\xff\"}"}{
  lifeRequest(t,h,"PATCH",url,body,400)
  if got:=lifeRequest(t,h,"GET",url,"",200);!reflect.DeepEqual(p,got){t.Fatalf("partial edit for %q: %v",body,got)}
 }
 for _,body:=range []string{`{}`,`null`,`{"confirmed":false}`,`{"confirmed":null}`,`{"confirmed":"true"}`,`{"Confirmed":true}`,`{"confirmed":true,"confirmed":true}`,`{"confirmed":true,"end_date":"2026-10-20"}`} {lifeRequest(t,h,"POST",url+"/finalize",body,400)}
 lifeRequest(t,h,"PATCH","/api/projects/missing",`{"name":"X"}`,404)
 lifeRequest(t,h,"POST","/api/projects/missing/finalize",`{"confirmed":true}`,404)
 lifeRequest(t,h,"PATCH",url,`{"start_date":"2026-10-21"}`,200)
 lifeRequest(t,h,"POST",url+"/finalize",`{"confirmed":true}`,409)
 lifeRequest(t,h,"PATCH",url,`{"start_date":null}`,200)
 lifeRequest(t,h,"POST",url+"/finalize",`{"confirmed":true}`,200)
 lifeRequest(t,h,"PATCH",url,`{"name":"After","start_date":"2026-10-01"}`,200)
 lifeRequest(t,h,"POST",url+"/finalize",`{"confirmed":false}`,400)
}
