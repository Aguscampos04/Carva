# Evidencia de IA, TDD y validación — US-02

**Fecha:** 2026-10-20. **Estado:** implementada y verificada técnicamente; revisión y aprobación final del equipo pendientes.
**Issue:** [#2 — Modificar proyecto y registrar fechas](https://github.com/Aguscampos04/Carva/issues/2).
**Rama:** `codex/us-02-sdd-bdd`, base `37c3752`, coincidente con `origin/main` después de fetch al iniciar.
**Cadena:** [enunciado §2.1](enunciado.md) → EP-01 → US-02 / #2 → [D-01…D-10](contexto-proyecto.md#modificación-y-finalización-del-proyecto--us-02--issue-2) → [SDD](../sdd/SDD-US-02-modificar-proyecto.md) → AC → [BDD](../bdd/BDD-US-02-modificar-proyecto.feature) → pruebas → Go/web → commits / PR.

## Autorización y alcance

El responsable aprobó SDD/BDD y PT-01 a PT-05 y autorizó implementación completa, commits, push y PR hacia main, sin merge ni cierre de Issues. Se ajustó y versionó documentación antes del código: BDD-07 ya no incluye En progreso sin Sprints; BDD-18 verifica identificadores, objetivos y snapshots actuales, sin estados o reglas futuras de cierre. No se implementa US-11 ni se inventan aprobaciones finales o ceremonias Scrum.

Se aplicaron AGENTS.md y skills project-manager, SDD, BDD, GitHub, Go, TDD y testing-quality. Issue #2 fue consultada y estaba abierta. Se preservó trabajo local y no se integraron cambios remotos sobre la copia existente.

## Implementación

- Dominio: día calendario validado y opcional, edición parcial, rango inicio ≤ fin, finalización manual e inmutabilidad del primer fin; elegibilidad para crear Sprint.
- Aplicación: puertos transaccionales, reloj inyectable y zona IANA embebida mediante biblioteca estándar. El reloj se lee después de adquirir escritura, solo para una transición efectiva.
- SQLite: migración aditiva `000008_projects_dates.sql`, columnas NULL sin inferencias históricas. Edición/finalización/creación de Sprint comparten coordinación `BEGIN IMMEDIATE`, lectura vigente, transformación de dominio y commit. Espera acotada a cinco segundos; rollback antes de reintentar BUSY/LOCKED y descarte de conexión si falla la limpieza. Errores de commit inciertos no se repiten automáticamente.
- REST: PATCH parcial y POST finalize con confirmación verdadera; JSON estricto exclusivamente en las nuevas escrituras, sin endurecer contratos anteriores. Errores controlados 400/404/409/500/503.
- Web: formulario de nombre/inicio y confirmación nativa del navegador; finalizar usa datos guardados, cancelar no escribe. Controles editables después del fin, fechas como días sin conversión de zona en JavaScript, descarte de respuestas obsoletas y doble envío bloqueado.
- Finalizar solo cambia el proyecto: conserva Sprints y snapshots y bloquea su nueva creación, incluso ante concurrencia. Miembros y Backlog continúan operativos. No se modifican reglas de cierre o asignación de Sprints.

## Evidencia RED → GREEN → REFACTOR

Cada RED indicado fue ejecutado antes del cambio que lo resolvió. Las pruebas suplementarias de concurrencia, migración y aplicación ampliaron la validación; no se presentan como ciclos RED si pasaron inicialmente.

| Etapa | Commit | Evidencia real |
| --- | --- | --- |
| Documentación previa | `b750d83` | SDD/BDD aprobados y ajustes BDD-07/18 registrados antes de programar. |
| RED inicial | `c46e06f` | `go test ./internal/domain/project` falló por ausencia de Date/ParseDate/Edit/Finalize; integración falló por ausencia de NewServiceWithClock; cuatro pruebas Node fallaron por falta de projects.js. |
| GREEN inicial | `94f09e6` | Suite Go completa PASS y diez pruebas web PASS. Incluye migración, JSON estricto, transacciones, concurrencia y conservación de relaciones. |
| RED contención al commit | `e9cd71b` | `TestProjectFinalizationRetriesAfterReaderBlocksCommit` falló con SQLITE_BUSY: un lector permite BEGIN IMMEDIATE pero impide COMMIT. |
| GREEN contención | `465f61e` | La prueba anterior, concurrencia de escritores y rollback PASS tras reintentar la unidad completa solo después de rollback confirmado. |
| REFACTOR | `dbf1575` | Lectura de proyecto unificada para pool y conexión transaccional; se preservan las fechas en todas las lecturas. Suite Go completa PASS; ampliación de verificaciones de reloj y compatibilidad. |
| RED respuesta obsoleta UI | `dbfb4bb` | Prueba del cableado real app.js/projects.js falló: mostraba Planificado tras una finalización exitosa porque se intercalaba una consulta anterior. |
| GREEN respuesta obsoleta UI | `dee728e` | Invalidar consultas al iniciar escritura y deshabilitar consulta durante guardado: once pruebas web PASS. |

No se reescribió historial para simular TDD. Se corrigieron dos problemas de herramientas sin contarlos como RED del producto: un import sobrante en tests nuevos y argumentos de cobertura separados por PowerShell. Una fixture de migración de miembros borraba todo el ledger posterior a v5; se acotó a v6/v7 para no reaplicar la nueva migración de fechas sobre columnas existentes. Se normalizaron finales de línea CRLF que gofmt reportaba, sin cambios funcionales ajenos a US-02.

## Trazabilidad BDD → pruebas

El `.feature` es especificación; no se agregó framework Gherkin. Se automatiza el comportamiento mediante Go testing/httptest, SQLite temporal real y Node test/vm con DOM mínimo y transporte controlado. No equivale a una suite E2E de navegador.

| BDD-US-02 | Prueba / comportamiento |
| --- | --- |
| 01 | `TestProjectLifecyclePersistenceAndCompatibility`: nombre/inicio/UUID conservados al reabrir DB y proyecto vecino sin cambios. |
| 02 | `TestProjectLifecyclePersistenceAndCompatibility`, `TestCalendarDatesAndProjectEdits`: nombre exacto y conservación de inicio omitido. |
| 03 | Las mismas pruebas: reemplazo y eliminación explícita del inicio, sin cambiar nombre/fin/estado. |
| 04 | `TestProjectLifecycleRejectsInvalidRequests`, `TestLifecycleApplicationDoesNotPersistInvalidChanges`: nombre nulo/blanco y rechazo atómico. |
| 05 | `TestCalendarDatesAndProjectEdits`, `TestProjectLifecycleRejectsInvalidRequests`: formatos, calendario, rango de años y vacío inválidos. |
| 06 | `TestCalendarDatesAndProjectEdits`: 2024-02-29 válido. |
| 07 | `TestFinalizeIsManualAndRepeatable`, `TestProjectLifecycleRejectsInvalidRequests`, `TestProjectLifecyclePersistenceAndCompatibility`: ambos estados, sin inicio/Sprints cuando Planificado; persistencia del fin tras reinicio. |
| 08 | Node `BDD-08`: cancelación conserva borrador y no envía solicitud; prueba Node de finalización envía solo confirmed. Diálogo nativo real no automatizado. |
| 09 | `TestProjectLifecyclePersistenceAndCompatibility`, `TestProjectLifecycleRejectsInvalidRequests` y Node: edición posterior a Finalizado y eliminación de inicio sin reapertura. |
| 10 | `TestCalendarDatesAndProjectEdits`, integración lifecycle: igualdad aceptada e inicio posterior rechazado sin escritura parcial. |
| 11 | `TestProjectDateDoesNotDriveState`: registrar inicio, avanzar reloj y consultar conserva Planificado sin fin; el reloj no se consulta al editar/leer. |
| 12 | `TestCreateProjectAndFirstSprintThroughAPI`, `TestSecondSprintKeepsProjectInProgress`, integración lifecycle: alta Planificado y primer Sprint En progreso, sin inferir inicio. |
| 13 | `TestLifecycleClockBoundaryAndRepeatedFinalization`: 02:30Z/03:00Z, frontera del día en Buenos Aires; integración con instante UTC controlado. |
| 14 | `TestProjectLifecycleRejectsInvalidRequests`, `TestFinalizeIsManualAndRepeatable`: inicio futuro editable, finalización rechazada, eliminación y reintento válido. |
| 15 | `TestProjectLifecyclePersistenceAndCompatibility`, `TestLifecycleClockBoundaryAndRepeatedFinalization`: fin original al reintentar otro día y reloj no vuelto a consultar; PATCH idéntico repetido conserva toda la representación. |
| 16 | `TestProjectLifecycleRejectsInvalidRequests`: id/status/end_date rechazados sin alteración. |
| 17 | La misma prueba: confirmed ausente/null/false/string/extras/duplicado inválidos; false tampoco finaliza de nuevo. |
| 18 | Integración lifecycle: dos Sprints conservan sus representaciones y snapshot consultable byte a byte; sin modelar cierre. |
| 19 | Integración lifecycle: Sprint nuevo en Finalizado devuelve 409. |
| 20 | Integración lifecycle: alta/listado de miembro y alta/edición/listado de Backlog en Finalizado, con reglas vigentes. Node comprueba que el editor no oculta Backlog. |
| 21 | `TestProjectLifecycleConcurrentWriters`: dos finalizaciones simultáneas y un solo fin persistido, entre pools independientes. |
| 22 | `TestProjectFinalizationWinsBeforeSprintCreation`, `TestProjectSprintCreationWinsBeforeFinalization` y escritores concurrentes: se fuerzan ambos órdenes; finalizar primero rechaza Sprint; crearlo primero permite finalizar conservándolo. |
| 23 | `TestProjectLifecycleRejectsInvalidRequests`: proyecto inexistente devuelve 404 sin crearlo. |
| 24 | `TestProjectLifecycleRollbackAndContention`: trigger ABORT al editar/finalizar y al insertar Sprint después de transición; rollback, sin SQL expuesto. DB cerrada no devuelve éxito. |
| 25 | La misma prueba: espera de cinco segundos → 503, contexto corto → error tipado, luego liberación y operación válida. Prueba de lector bloqueando COMMIT cubre el reintento. |
| 26 | `TestProjectLifecycleRejectsInvalidRequests`: objeto completo, claves exactas/duplicadas, extras, tipos, valores posteriores y UTF-8 inválido. |
| 27 | `TestProjectConcurrentFutureStartAndFinalization`: se fuerza cada orden bajo bloqueo; la segunda operación falla, nunca inicio > fin. |
| 28 | `TestProjectLifecycleConcurrentWriters`: nombre/inicio/finalizaciones simultáneos conservan campos distintos y estado/fin. |

Adicionalmente `TestProjectDatesMigrationPreservesExistingData` construye el esquema real v7 a partir de migraciones sin editar, conserva proyectos/miembro/Sprint al migrar dos veces y no inventa fin para un Finalizado histórico. El smoke comprueba persistencia tras reinicio del proceso real.

## Validación ejecutada

- Go `go1.27.1 windows/amd64` (mínimo del módulo 1.26.0), Node `v20.17.0`. Sin dependencias nuevas.
- `go test ./...`: PASS, incluida regresión de US-01/03/05/06/07/08. Se usan archivos SQLite temporales, no la DB de desarrollo.
- `go vet ./...`: PASS. `gofmt -l .`: sin resultados tras normalizar formato/finales de línea.
- `go build -o .local/carva-us-02.exe ./cmd/carva`: PASS.
- `node --test tests/web/members.test.cjs tests/web/projects.test.cjs`: 11 tests PASS. `node --check` de app.js/projects.js: PASS.
- `go test '-coverpkg=./internal/...' '-coverprofile=.local/us-02-coverage.out' ./...`: PASS. `go tool cover '-func=.local/us-02-coverage.out'`: cobertura combinada interna **85,9%**. ParseDate, Validate, Edit, Finalize y CanCreateSprint de dominio **100%**; Update de aplicación, Mutate de persistencia y rutas/traducción de errores de mutación **100%**; decoder del objeto **86,5%**. Incluye integración, no solo tests por paquete. No se interpreta cobertura como prueba de ausencia de defectos.
- `go test -race ./...`: **no pudo compilar** runtime/cgo: `cc1.exe: sorry, unimplemented: 64-bit mode not compiled in`. GCC disponible `mingw32`, sin compilador de 64 bits detectado. No se afirma PASS de race; queda pendiente en un entorno compatible. Las pruebas funcionales de concurrencia SQLite sí pasan.
- Smoke con binario real en puerto libre 8080 y DB aislada `.local/us-02-smoke-30eea22b236e427cb7bcaa5d19c8fc6c.db`: PASS en HTML/JS embebidos, alta, edición, creación de Sprint, finalización, rechazo 409 de Sprint nuevo, conservación de Sprint y consulta tras detener/reiniciar. Se detuvieron solo los procesos de este smoke; no se modificó la DB de desarrollo.

## Revisión y límites

AC-01 a AC-15: comportamientos automatizados según matriz, sin evidencia de revisión humana final. La presentación nativa de confirmación y recorrido visual completo no se probaron en un navegador real; se verificó el controlador, su cableado y la entrega del HTML/JS en el smoke. No se instala un framework E2E para esta historia.

Revisión técnica: reglas en Go y separación de capas, SQL parametrizado, migración anterior intacta, conservación de datos y errores sin detalles internos. La política de campos concurrentes sigue la SDD aprobada: campos distintos se conservan; para el mismo campo prevalece la última escritura serializada válida, sin ETag ni historial nuevo.

US-11 conserva la decisión pendiente sobre cierre con historias incompletas cuando no se pueden crear más Sprints. No bloquea finalizar y no se resuelve aquí. Sin merge, cierre de Issues ni aprobación final atribuida al equipo. Publicación y PR se informan únicamente después de confirmarse.
