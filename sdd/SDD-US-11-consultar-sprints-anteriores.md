# SDD-US-11 — Consultar Sprints anteriores

**Issue:** #11 · EP-03 · enunciado §2.3 · Sprint 2 — Flujo principal.
**Estado:** implementación autorizada por el responsable; revisión final por PR pendiente.
**Dependencia:** US-10, PR #39 integrado en `fc572ef`; el Agile Enabler confirmó la validación del equipo en este chat antes de implementar US-11.

## 1. Objetivo

**Integración aprobada US-16:** mostrar horas estimadas **vigentes**, indicando que pueden actualizarse, desde la columna persistida de cada snapshot. D-21 permite actualizar exclusivamente ese dato incluso en cerrado; no reconstruir otros campos desde Backlog ni cambiar resultados, puntos o pertenencias. Véase [SDD-US-16](SDD-US-16-registrar-horas-esfuerzo.md).

Consultar los Sprints cerrados de un proyecto y sus objetivos, planificación y resultados persistidos. La consulta no cambia el historial cuando una historia se traslada, se edita en el Backlog o se completa en otro Sprint.

## 2. Inputs

- `projectID`: identificador del proyecto existente.
- `sprintID`: entero positivo, para los endpoints de detalle existentes.
- La interfaz admite el identificador de un proyecto existente, además del proyecto recién creado/consultado en pantalla.

## 3. Expected Outputs / contratos de consulta

- `GET /api/projects/{projectID}/sprints/closed`: `200`, arreglo de Sprints con `id`, `projectID`, `goal`, `closed: true`, ordenados por ID ascendente (orden existente de creación). Proyecto sin cerrados: `[]`.
- Se reutiliza `GET /api/projects/{projectID}/sprints/{sprintID}` para el detalle, sin cambiar su contrato ni impedir la consulta preexistente de Sprints abiertos.
- Se reutiliza `GET /api/projects/{projectID}/sprints/{sprintID}/backlog-items`: snapshots con título, descripción, prioridad, puntos y criterios originales; `planned_status`, `status` y `acceptance_criteria_verified` persistidos por US-08/09/10. Sin historias: `[]`.
- La UI lista los Sprints cerrados, permite seleccionarlos y muestra su Goal, identificadores y snapshots con resultado registrado. Incluye carga, vacío y error; descarta respuestas de un proyecto/selección anterior.

## 4. Business Rules

- BR-01: listar únicamente Sprints con `closed = 1` pertenecientes al proyecto solicitado; no filtrar datos de otro proyecto.
- BR-02: leer los snapshots de `sprint_backlog_items`, nunca reconstruirlos desde el Backlog actual.
- BR-03: `acceptance_criteria_verified` representa la finalización registrada en ese Sprint, no la finalización posterior en otro. No se agregan estados nuevos al dominio.
- BR-04: consultar no realiza escrituras, cierres ni traslados. Los cambios posteriores del Backlog no modifican el historial.

## 5. Constraints / alcance

Sin migraciones, dependencias externas ni cambios a US-10. Se agrega solo el listado de cerrados en las capas existentes y una sección web de consulta. No se incorpora gestión web de altas/cierre de Sprint ni métricas. Los contratos técnicos de lectura siguen las convenciones existentes; no crean nuevas reglas funcionales. La selección de datos es responsabilidad de Go/SQLite, no del frontend.

## 6. Edge Cases

Proyecto sin Sprints, con abiertos únicamente o con cerrados vacíos; varios proyectos; ID con caracteres que deben codificarse en URL; historia trasladada y completada luego; datos de texto HTML deben renderizarse como texto; respuestas tardías al cambiar selección/proyecto.

## 7. Error Conditions

- Proyecto inexistente en listado: `404`, formato de error existente.
- Sprint inexistente o de otro proyecto en detalle: `404`; ID no positivo/no entero: `400` (contratos existentes).
- Falla de almacenamiento: `500`, sin detalles SQL.
- La UI conserva un mensaje de error visible y permite reintentar; borra detalles anteriores al cambiar de proyecto o Sprint.

## 8. Acceptance Criteria / trazabilidad

| Criterio | BDD | Tests |
| --- | --- | --- |
| AC-01: consultar cerrados con Goal, plan y resultados | BDD-01, 04, 05 | `TestUS11ClosedSprintListing`, `TestUS11HistorySurvivesTransferAndLaterChanges`, tests Node del historial |
| AC-02: traslado no altera historial | BDD-02, 03 | `TestUS11HistorySurvivesTransferAndLaterChanges` |
| BR-01, errores y aislamiento | BDD-04, 05 | tests de aplicación, SQLite/HTTP y web |

Pruebas de integración con SQLite temporal y reapertura. TDD y validaciones reales se registran en `docs/evidencia-ia-us-11.md`. La validación humana final de US-11 corresponde al equipo.
