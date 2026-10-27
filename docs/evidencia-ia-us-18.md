# Evidencia de implementación — US-18

- **Issue:** [#18 — Registrar defecto](https://github.com/Aguscampos04/Carva/issues/18)
- **Rama:** `codex/issue-18`
- **Base:** `main` con US-16 integrada. US-17 no se incluye porque aún no está mergeada y no es dependencia de esta historia.
## Resultado

- Registro de defectos desde cada User Story del Product Backlog y consulta de la lista asociada.
- Go valida descripción, los cuatro estados y las cuatro severidades aprobadas; genera el UUID.
- Los Sprints de detección y resolución quedan opcionales (`null` hasta que se informen). Cuando se proveen, deben pertenecer al mismo proyecto. No se deduce ni se asigna un Sprint ficticio.
- SQLite agrega `000013_defects.sql` con estados/severidades restringidos y FK compuestas para proyecto/historia/Sprints.
- La UI permite crear y consultar defectos por historia, conserva los detalles como texto y no incluye reglas de actualización de US-19.

No se añadieron transiciones de estado, edición, eliminación ni métricas; son parte de otras historias.

## Trazabilidad

§2.6 → EP-06 → US-18 / #18 → [SDD](../sdd/SDD-US-18-registrar-defecto.md) → [BDD](../bdd/BDD-US-18-registrar-defecto.feature) → tests de dominio, integración y web → `internal/domain/defect`, SQLite, API y UI.

AC-01/AC-02 están cubiertos con alta sin Sprint y con Sprint de detección; el test de dominio recorre las 16 combinaciones permitidas de estados/severidades y rechaza valores no aprobados. La prueba de integración también comprueba el aislamiento entre historias y proyectos.

## Evidencia TDD en commits

- `b1f2583` agrega tests de dominio y falla por los tipos/reglas aún inexistentes; `c192ddb` implementa el dominio y los tests pasan.
- `346640d` agrega la prueba API/integración y falla porque no existe el servicio de defectos; `8758a06` completa la API y la composición, y la prueba pasa.
- `bfa0a26` agrega pruebas de interfaz y falla porque `defects.js` no existe; `5215d65` implementa el formulario/listado y ambas pasan.
- `7c9f8a7` actualiza la regresión de migraciones al total nuevo (13 migraciones, 10 tablas de negocio).
- `ede39f2` actualiza el fixture de `app.js` para cargar el nuevo módulo, igual que hace la página real.

## Validación realizada

| Comando/verificación | Resultado |
| --- | --- |
| `gofmt -l .` | PASS, sin archivos. |
| `go test ./... -count=1` | PASS, suite completa, incluida integración. |
| `go vet ./...` | PASS. |
| Suite web Node | PASS: 30 tests, 0 fallos. |
| `git diff --check origin/main...HEAD` | PASS. |

La Issue #18 permanece abierta para revisión del equipo; no se ha hecho merge ni cierre.
