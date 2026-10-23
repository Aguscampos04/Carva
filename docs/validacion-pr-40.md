# PR #40 — Corrección y validación final de US-11

## Alcance y autorización

El responsable solicitó investigar el test de concurrencia preexistente de US-02,
resolver el formato global y actualizar el mismo PR #40 sin reimplementar US-11,
sin nuevas historias ni merge. Rama: `codex/us-11-consultar-sprints-anteriores`.
Checkout inicial limpio; HEAD local/remoto `8f20873`; main `fc572ef`, PR #39 incluido.
El PR estaba abierto y ya no era borrador al iniciar esta etapa. No se cambió ese
estado ni se aprobó en nombre del responsable. No se modificaron dependencias,
migraciones ni reglas de negocio. La aclaración posterior del responsable sobre
Issues futuras se respetó: ninguno de estos problemas depende de una historia futura.

## Investigación y reproducción sobre main

Se extrajo `git archive main` a una carpeta temporal fuera del checkout y se ejecutó:

```powershell
go test ./tests/integration -run '^TestProjectFinalizationWinsBeforeSprintCreation$' -count=50
```

La ejecución abortó con `panic: close of closed channel` en la línea 92, antes de
completar las 50 repeticiones. Confirma que el problema existe sin US-11. No se
cambió main ni se eliminó trabajo local.

`ProjectRepository.Mutate` documenta que su transformación puede repetirse después
de un BUSY/LOCKED con rollback confirmado. `Finalize` invoca el reloj dentro de
esa transformación. El reloj del test cerraba `entered` en cada invocación: la
segunda ejecución cerraba un canal ya cerrado. El stack y una prueba nueva que
fuerza un COMMIT fallido verificaron esa causa. Es un efecto lateral del fixture,
no un defecto demostrado en las reglas de negocio ni una race de acceso a memoria.

La prueba nueva mantiene una transacción lectora SQLite, espera la entrada del
escritor, permite su primer intento y espera la segunda invocación del reloj.
Solo entonces libera el lector y permite confirmar. Antes de corregir el reloj
reprodujo el pánico de forma controlada (RED, commit `6b2e4e6`). No usa Sleep ni recover.

## Correcciones de sincronización

1. `gatedProjectClock` señala la primera entrada con `sync.Once` y conserva la
   misma barrera y fecha para cualquier reintento. Commit `1b54144`.
2. La primera repetición después de esa corrección eliminó el pánico, pero obtuvo
   ocho respuestas 201 en el caso que esperaba 409. La prueba suponía que ganar
   el primer BEGIN IMMEDIATE equivalía a ganar el COMMIT. Un rollback permite
   legítimamente que la creación gane antes de la finalización; eso no viola el
   dominio. No se amplió la aserción para admitir 201.
3. Se estableció con un canal que la finalización debe haber confirmado antes de
   que la solicitud de creación acceda al repositorio. Se conserva el escenario
   BDD-22 de orden finalización→creación y la aserción original 409. Commit `28aacd6`.
4. La nueva prueba `TestProjectFinalizationClockSurvivesCommitRetry` sigue forzando
   contención y replay reales, y verifica estado Finalizado y fecha final persistidos.
   `TestProjectLifecycleConcurrentWriters` y el orden inverso existente siguen
   verificando concurrencia real y sus resultados permitidos. No se retiraron tests.

No se tocó código de producción para resolver estos fallos. No se desactivaron
aserciones ni se ocultaron errores. Las señales son eventos ordenados por canales;
el timeout solo detecta bloqueos, no organiza el orden de ejecución.

## Investigación del formato y corrección independiente

`gofmt -l .` informó inicialmente 56 archivos (listados al final). Todos son
fuentes Go versionadas; ninguno tiene marca de código generado. Ninguno de esos
56 fue modificado por la implementación original de US-11 (el test se modificó
recién durante esta etapa). Los cinco Go originales afectados por US-11 estaban
formateados y con LF.

Se compararon los bytes de cada blob `HEAD:<ruta>` con `gofmt` mediante Node y
subprocesos binarios: **cero diferencias de formato en los blobs versionados**.
Para los 56 archivos del checkout, eliminar CRLF→LF produjo exactamente la salida
de gofmt. `core.autocrlf=true` y `git ls-files --eol` confirmaron `i/lf w/crlf`.

Por tanto no había 56 archivos con indentación incorrecta: el problema era la
conversión automática del checkout. Se agregó `.gitattributes` con `*.go text eol=lf`
y se ejecutó gofmt únicamente sobre las fuentes versionadas reportadas, verificando
que no fueran generadas. El test ya había sido normalizado en su commit, por lo
que esa pasada normalizó las otras 55 copias. La renormalización del índice no
produjo cambios de contenido en ellas. Commit independiente `d25e666`: solo la
regla de atributos, sin refactorizaciones ni modificaciones funcionales.

Resultado: `gofmt -l .` no devuelve archivos, incluidos los no versionados;
`git ls-files --eol '*.go'` no muestra fuentes con `w/crlf`. La regla evita que un
checkout futuro vuelva a provocar el mismo diagnóstico en Windows.

## Validaciones posteriores a las correcciones

| Comando | Resultado |
| --- | --- |
| `go test ./tests/integration -run '^TestProjectFinalization(WinsBeforeSprintCreation\|ClockSurvivesCommitRetry)$' -count=50` | PASS: 50 ejecuciones por caso, 100 en total. |
| `gofmt -l .` | PASS: sin salida. |
| `go test ./...` | PASS, todos los paquetes. |
| `go test ./... -count=1` | PASS, todos los paquetes sin caché de resultados. |
| `go vet ./...` | PASS. |
| `go build ./...` | PASS. |
| `go mod verify` | PASS: all modules verified. |
| `git diff --check` y `git diff main --check` | PASS. |
| `go test ./internal/application/sprints ./tests/integration -run 'TestListClosed\|TestUS11\|TestCloseSprint\|TestUS10' -count=1` | PASS. |
| `node --test tests/web/*.test.cjs` | PASS: 15/15. |
| `go test -race ./...` | No ejecutable: CGO deshabilitado por defecto. |
| `CGO_ENABLED=1`, luego `go test -race ./...` | No ejecutable: compilador C gcc ausente de PATH. No se afirma PASS. |

CGO se habilitó solo en el proceso de comprobación y se restauró después; no se
instaló un toolchain ni se cambiaron Go, SQLite o dependencias del proyecto.
Las repeticiones aportan evidencia de estabilidad, no una demostración matemática
de ausencia de errores concurrentes. La prueba con detector de carreras debe
completarse en un entorno con compilador C compatible si se requiere ese control.

## Conservación funcional y entrega

Los tests de US-11 vuelven a verificar Goal y planificación, traslado, edición y
completitud posterior, reapertura de SQLite y snapshots originales intactos.
Los tests de cierre/traslado US-10 y la regresión web pasan. No se cambió el código
de US-11 en esta etapa y no se detectaron regresiones en los controles ejecutados.

El PR queda listo para revisión humana y merge manual, con la limitación de
`-race` indicada. No hay cierre manual de #11, auto-merge, aprobación automática
ni otro PR. La validación humana final continúa siendo responsabilidad del equipo.

## Inventario exacto reportado inicialmente por gofmt

```text
cmd/carva/main.go
internal/application/backlog/doc.go
internal/application/backlog/service.go
internal/application/backlog/service_test.go
internal/application/members/doc.go
internal/application/members/service.go
internal/application/members/service_test.go
internal/application/projects/doc.go
internal/application/projects/lifecycle_test.go
internal/application/projects/service.go
internal/application/projects/service_test.go
internal/application/sprints/doc.go
internal/domain/backlog/doc.go
internal/domain/backlog/item.go
internal/domain/backlog/item_test.go
internal/domain/member/doc.go
internal/domain/member/member.go
internal/domain/member/member_test.go
internal/domain/project/doc.go
internal/domain/project/lifecycle.go
internal/domain/project/lifecycle_test.go
internal/domain/project/project.go
internal/domain/project/project_test.go
internal/domain/sprint/doc.go
internal/domain/sprint/sprint.go
internal/domain/sprint/sprint_test.go
internal/persistence/sqlite/backlog_repository.go
internal/persistence/sqlite/db.go
internal/persistence/sqlite/doc.go
internal/persistence/sqlite/immediate.go
internal/persistence/sqlite/member_repository.go
internal/persistence/sqlite/migrate.go
internal/persistence/sqlite/migrations/embed.go
internal/persistence/sqlite/project_repository.go
internal/persistence/sqlite/sprints.go
internal/presentation/httpapi/backlog.go
internal/presentation/httpapi/doc.go
internal/presentation/httpapi/members.go
internal/presentation/httpapi/project_mutations.go
internal/presentation/httpapi/projects.go
tests/integration/backlog_http_test.go
tests/integration/backlog_test.go
tests/integration/members_test.go
tests/integration/project_http_test.go
tests/integration/project_lifecycle_test.go
tests/integration/project_lifecycle_transactions_test.go
tests/integration/project_reader_contention_test.go
tests/integration/project_sprint_api_test.go
tests/integration/project_status_test.go
tests/integration/project_test_helpers_test.go
tests/integration/projects_test.go
tests/integration/sprint_backlog_api_test.go
tests/integration/sprint_closure_http_test.go
tests/integration/sprint_closure_repository_test.go
tests/integration/sqlite_test.go
web/static/assets.go
```
