# Evidencia de IA — US-01 Crear proyecto

**Issue:** #1 (referenciada por el Developer; su página no fue accesible desde este entorno).
**Rama local:** `codex/us-01-crear-proyecto`.
**Requisito:** `docs/enunciado.md` §2.1.
**Trazabilidad:** EP-01 → US-01 → SDD-US-01 → AC-01/AC-02 → BDD-US-01 → tests → Go.

## Instrucción recibida

**Prompt:** «Quiero que me ayudes a desarrollar la US-01 — Crear proyecto, del proyecto Carva.» La instrucción completa está en el adjunto `Texto pegado.txt`; solicitó implementar US-01 de extremo a extremo, consultar la Issue #1 y las decisiones vigentes, producir SDD/BDD, seguir RED → GREEN → REFACTOR con commits que reflejen evidencia real, mantener Go/SQLite/API/web en las capas aprobadas y ejecutar `gofmt`, tests, `vet`, `build`, verificación de módulos y cobertura. Indicó expresamente no hacer push, PR, cierre de Issue ni cambios de Sprint.

## Decisiones usadas

- Nombre obligatorio; se rechaza si queda vacío al recortar espacios. Se almacena el texto ingresado, sin límite de longitud ni unicidad.
- UUID generado por el servidor.
- El estado inicial `Planificado` se establece en el dominio Go. El cliente no controla ID ni estado.
- Persistencia SQLite mediante `modernc.org/sqlite v1.60.1` y migración SQL embebida.
- Contrato de US-01: `POST /api/projects`, `GET /api/projects/{id}`; la interfaz web usa la API del mismo proceso.
- La Issue remota no se pudo leer; los criterios utilizados son el texto facilitado por el Developer y el borrador local de Issues.

## Resultado producido

- Especificación y escenarios versionados en `sdd/SDD-US-01-crear-proyecto.md` y `bdd/BDD-US-01-crear-proyecto.feature`.
- Entidad y regla de estado inicial en `internal/domain/project/`.
- Casos de uso para crear y consultar en `internal/application/projects/`.
- Repositorio SQLite y tabla creada por `000002_projects.sql`.
- API REST, servidor de composición y UI básica para alta y consulta.
- Pruebas unitarias e integración con base temporal, API HTTP real en `httptest` y persistencia tras reabrir SQLite.

## Evidencia TDD

| Etapa | Evidencia | Resultado real |
| --- | --- | --- |
| RED dominio | Commit `2210be4` (`test: define US-01 project domain behavior #1`) | `go test ./internal/domain/project` falló porque faltaban `New`, `Planned` y los errores definidos por el dominio. |
| RED aplicación | Commit `311b021` (`test: specify US-01 project use case #1`) | Luego de resolver la caché de módulos, `go test ./internal/application/projects` falló por la ausencia del servicio, entidad y reglas. |
| GREEN dominio/aplicación | Commit `ea296de` (`feat: implement US-01 project domain and use case #1`) | Ambos paquetes pasaron. |
| RED SQLite | Commit `3263061` (`test: cover US-01 SQLite project persistence #1`) | Una vez disponible el driver, la integración falló porque aún no existía `NewProjectRepository`. |
| GREEN SQLite | Commit `36f7c2c` (`feat: persist US-01 projects in SQLite #1`) | Alta, lecturas, estado inicial, dos registros y conservación tras reabrir DB pasaron. |
| RED HTTP | Commit `32a617b` (`test: specify US-01 HTTP project flow #1`) | La integración falló porque aún no existía `httpapi.NewHandler`. |
| GREEN HTTP | Commit `54a3081` (`feat: expose US-01 project creation on web #1`) | El recorrido HTTP contra SQLite pasó, incluido el intento de manipular ID/estado y un fallo de persistencia. |
| REFACTOR | Commit `ee2ee14` (`refactor: share temporary project database setup #1`) | Se unificó la preparación de bases temporales repetida en las pruebas de integración; los paquetes afectados volvieron a pasar. |

Los fallos iniciales de los tres primeros intentos de integración causados por bloqueo de la caché temporal/global de Go no se consideran evidencia RED; solo se registraron fallos después de completar dependencias y compilar las pruebas.

## Validación ejecutada

- Go disponible: `go1.27.1 windows/amd64`; módulo mínimo: Go 1.26.0.
- `gofmt -l` sobre todos los Go modificados: sin resultados.
- `gofmt -l .`: reportó archivos Go preexistentes sin formato fuera de esta contribución; no se reformatearon paquetes ajenos.
- `go test ./...`: PASS.
- `go vet ./...`: PASS.
- `go build ./...`: PASS (repetido con acceso a caché; sin errores).
- `go mod verify`: PASS — `all modules verified`.
- `go test -cover ./...`: PASS con permisos para la carpeta temporal. Dominio de proyectos 100%; aplicación de proyectos 100%; los paquetes HTTP y SQLite aparecen con 0% en el reporte estándar por tener su cobertura ejecutada desde `tests/integration`, y la carpeta de integración no contiene sentencias de producción.
- UI demostrada en navegador: se creó **US-01 Demo Carva**, luego se consultó de nuevo y se observó el mismo UUID y estado **Planificado**. La demostración usó `.local/us-01-demo.db`, separada de `.local/carva.db`.
- No había herramienta de frontend/CI definida en el repositorio; el recorrido web se verificó manualmente en el navegador.

## Revisión humana y publicación

- Revisión humana: **pendiente**, no confirmada.
- Commits locales: registrados en el historial de la rama.
- Push, PR, cierre de Issue, cambio de Sprint e integración a `main`: **no realizados**.
- Próximo paso: revisión del Developer/equipo sobre la especificación, los contratos y el código antes de publicar.
