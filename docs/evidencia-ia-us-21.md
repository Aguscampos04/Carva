# Evidencia de desarrollo asistido — US-21

## Objetivo

Implementar la consulta de horas estimadas, horas reales y desviación por Sprint y por proyecto, conservando la trazabilidad de la Issue #21 y el flujo SDD → BDD → TDD → Go → pruebas.

## Prompt utilizado

“Implementar la US-21 asignada a Loana Vargas; respetar requisitos, arquitectura y decisiones del repositorio; definir fórmulas en SDD; aplicar BDD y TDD; validar; dividir el trabajo en commits locales separados; no hacer merge.”

## Decisión funcional aprobada

Loana aprobó el 2026-10-27 las fórmulas de [SDD-US-21](../sdd/SDD-US-21-metricas-horas.md): estimaciones del Sprint desde su snapshot; horas reales desde registros activos explícitamente asociados; horas estimadas de proyecto contando una vez cada historia actual del Backlog; horas reales de proyecto incluyendo esfuerzo sin Sprint; desviación con signo `reales − estimadas`; datos insuficientes como `null` / «No disponible».

## Resultado producido

- SDD y seis escenarios BDD para sumas, desviación positiva/negativa, datos incompletos, historias trasladadas, esfuerzo sin Sprint, esfuerzo eliminado y errores de consulta.
- Dominio Go para sumas exactas en centésimas y desviación con signo.
- Servicio, repositorio SQLite y endpoints de consulta por proyecto y Sprint.
- Interfaz web que presenta Story Points y las tres métricas de horas, con «No disponible» para datos ausentes.
- Pruebas unitarias, de integración y frontend; sin migraciones ni dependencias nuevas.

## Evidencia TDD

| Etapa | Evidencia |
| --- | --- |
| RED dominio | Commit `92de6aa`; `go test ./internal/domain/metrics -run TestCalculateHours -count=1` falló al no existir `HoursMetrics`, `DecimalHours` ni `CalculateHours`. |
| GREEN dominio | Commit `7c34c4b`; la prueba de cálculo pasó con sumas exactas, desviación positiva/negativa, estimación faltante, suma vacía, desbordamiento y JSON. |
| RED integración | Commit `01f4978`; `go test ./tests/integration -run TestUS21 -count=1` falló porque los endpoints de horas respondían HTTP 404. |
| GREEN backend | Commit `660cf15`; las pruebas de integración de US-21 pasaron con asociación a Sprint, suma de proyecto, entradas sin Sprint y entradas eliminadas. |
| RED frontend | Commit `4ed6ec7`; `node --test tests/web/metrics.test.cjs` falló porque faltaba la visualización de horas. |
| GREEN frontend | Commit `0aa853c`; las pruebas de interfaz pasaron mostrando los valores y restableciendo el estado de «No disponible». |

## Archivos relacionados

- `sdd/SDD-US-21-metricas-horas.md`
- `bdd/BDD-US-21-metricas-horas.feature`
- `internal/domain/metrics/hours.go` y `hours_test.go`
- `internal/application/metrics/service.go`
- `internal/persistence/sqlite/metrics_repository.go`
- `internal/presentation/httpapi/metrics.go`
- `tests/integration/metrics_hours_test.go`
- `web/static/index.html`, `web/static/js/app.js`, `web/static/js/metrics.js`
- `tests/web/metrics.test.cjs`

## Validación registrada

- `gofmt -l .`: PASS, sin archivos listados.
- `go test ./... -count=1`: PASS.
- `go vet ./...`: PASS.
- `go build ./...`: terminó con código 0; Go informó que no pudo escribir un archivo de stat cache en la caché de módulos fuera del worktree.
- `go mod verify`: PASS, todos los módulos verificados.
- `node --test tests/web/*.test.cjs`: PASS, 39 pruebas.

## Revisión humana y publicación

La revisión humana del código todavía está pendiente. Los commits son locales; no se hizo push ni se creó un Pull Request en esta etapa. No se afirma aprobación ni merge del equipo.
