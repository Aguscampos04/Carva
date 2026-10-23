# US-11 — Evidencia de implementación, IA y validación

## Contexto y autorización

Solicitud: implementar Issue #11 de extremo a extremo con SDD, BDD, TDD,
Go/SQLite, API y web; publicar rama y PR sin merge. Fuente: instrucción del
responsable en este chat del 2026-10-23. Se autorizan commits y push únicamente
de US-11. El responsable confirmó «ya esta confirmado» al solicitarse la
validación formal de US-10; no se inventaron reviews ni se modificó la Issue #10.

Comprobado antes de escribir código: #11 asignada a Aguscampos04 y Milestone
Sprint 2 — Flujo principal; PR #39 integrado en `fc572ef`; main local coincide
con main remoto; checkout limpio y sin PR equivalente abierto. Skills consultadas:
project-manager, SDD, BDD, TDD, Go, testing-quality y GitHub en `.agents/skills/`.

## Resultado y decisiones

Solo se agrega el listado de Sprints cerrados y su acceso web. Se reutilizan las
entidades Sprint/PlannedBacklogItem y las consultas de snapshots existentes.
Sin dependencias ni migraciones nuevas y sin modificar el cierre/traslado.
Los contratos de detalle y errores existentes se conservan. Las decisiones
técnicas de lectura están descritas en SDD-US-11; revisión final pendiente por PR.

## Trazabilidad y TDD real

Enunciado §2.3 → EP-03 → #11 → SDD-US-11 → AC-01/02 → BDD-US-11 → tests → código.

| Etapa | Evidencia real |
| --- | --- |
| Especificación | `8783f8d`: SDD y seis escenarios BDD antes del código. |
| RED backend | `020e4fc`: `go test ./internal/application/sprints ./tests/integration -run 'TestListClosed\|TestUS11' -count=1` falló: método `ListClosed` ausente y GET `/sprints/closed` devolvía 400. |
| GREEN backend | `11943aa`: mismo comando pasó en aplicación e integración. |
| RED web | `24682e5`: cuatro tests Node fallaron porque no existía `sprint-history.js`. |
| GREEN web | `bea1913`: interfaz y montaje; 15/15 tests web pasaron. El harness de montaje incorpora el controlador nuevo igual que su stub de miembros. |
| REFACTOR | `52d21e5`: guardia de selección extraída y fixture de aislamiento simplificado; regresión web 15/15 y suite focalizada US-10/11 pasaron. |

`TestUS11HistorySurvivesTransferAndLaterChanges` verifica AC-01/02 y BDD-01/02/03:
cierre con historias completa/incompleta, traslado, edición de título/descripción/
prioridad/puntos/criterios, completitud posterior, reapertura de SQLite y dos
lecturas del historial idénticas al registro anterior. Se verifica por separado
que la historia del destino sí quedó completada.

`TestUS11ClosedSprintListing` verifica vacío, orden, filtro de cerrados,
aislamiento de proyectos, recursos inexistentes y almacenamiento cerrado.
`TestListClosedValidatesProjectAndPropagatesErrors` cubre validación y propagación
de errores en aplicación. Los tests web cubren BDD-01/04/05/06: datos originales,
texto seguro, solo GET, errores/vacíos/reintentos y descarte de respuestas tardías.

## Resultados ejecutados

Entorno: Go 1.27.1, Windows amd64. No se modificaron las versiones del módulo.

- Suite focalizada US-10/11: PASS, incluyendo tests existentes de cierre/traslado.
- `node --test tests/web/*.test.cjs`: PASS, 15 tests.
- `go test ./...`: dos ejecuciones fallaron por un pánico preexistente en
  `TestProjectFinalizationWinsBeforeSprintCreation`, línea 92, `close of closed channel`.
  El primer fallo ocurrió sobre main limpio antes de implementar US-11.
- `go test ./... -count=1` fuera del entorno restringido: PASS en todos los paquetes.
  La repetición aislada del test anterior volvió a fallar: la regresión completa
  no se considera estable pese a la ejecución exitosa. No se modificó US-02.
- `go vet ./...`: PASS.
- `go build ./...`: PASS.
- `go mod verify`: PASS, `all modules verified`.
- `git diff --check` y `git diff main --check`: PASS.
- `gofmt -l .`: ejecutado; reporta numerosos archivos preexistentes ajenos a US-11.
  `gofmt -l` sobre todos los Go agregados/modificados para US-11 no devuelve archivos.
  No se realizó reformateo masivo ajeno.
- Cobertura focalizada (US-10/11, tres paquetes instrumentados): total 49.4%;
  `sprints.Service.ListClosed`: 100%; `sqlite.SprintRepository.ListClosed`: 76.9%.
  El porcentaje global incluye funciones anteriores fuera del alcance; no indica
  cobertura completa del producto. Perfil local ignorado: `.local/us11.cover`.

Durante validación, algunos comandos restringidos informaron acceso denegado a
la caché Go; vet/build/módulos se comprobaron después con acceso autorizado.
Los primeros intentos de exportar el perfil temporal no produjeron archivo; solo
el perfil `.local/us11.cover` es la fuente de los porcentajes anteriores.

## Prueba funcional en navegador

Servidor real iniciado con `CARVA_DB_PATH` apuntando a una base temporal nueva,
separada de `.local/carva.db`. Se crearon un proyecto, dos Sprints y una historia;
se asignó al primero, se cerró/trasladó y se completó en el segundo por API.
En el navegador se ingresó el UUID en la sección Sprints anteriores, se consultó
y seleccionó el Sprint 1. Se observaron el Goal «Conservar el historial original»,
«Historia original», criterio original, 3 puntos provisionales y resultado
«Incompleta»/estado «Pendiente», pese a estar completada en el Sprint posterior.
Esta prueba fue ejecutada por el agente; no se presenta como revisión humana.

## Estado de entrega

AC-01/02 comprobados técnicamente. SDD/BDD/tests/UI implementados y trazables.
Revisión humana de US-11 y merge manual pendientes. La inestabilidad preexistente
del test de concurrencia y el formato global impiden afirmar una DoD global sin
reservas. El PR debe explicitarlo y referenciar #11 sin `Closes #11`; no se cierra
ni se modifica la Issue y no se activa auto-merge.

Para verificar manualmente: `go run ./cmd/carva`, abrir `http://localhost:8080`,
ingresar el UUID de un proyecto con Sprints cerrados y elegir uno. Comparar sus
snapshots antes/después de editar o completar en un Sprint posterior. Probar un
proyecto sin cerrados y uno inexistente. Los endpoints para preparar datos son
los ya disponibles de proyectos, Backlog, asignación, completion y close.
