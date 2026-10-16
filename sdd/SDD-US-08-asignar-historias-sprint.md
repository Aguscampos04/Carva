# SDD-US-08 — Asignar historias y preservar lo planificado

**Estado:** alcance de snapshot aprobado por el responsable el 2026-10-16; contrato técnico sujeto a revisión del equipo.
**Trazabilidad:** enunciado §2.3 → EP-03 → US-08 / Issue #8 → este SDD → BDD-US-08 → tests y código.
**Dependencias:** US-05 (Product Backlog) y US-07 (Sprint).

## 1. Objetivo

Asignar historias existentes del Product Backlog a un Sprint y permitir consultar la planificación guardada en ese Sprint. Los datos planificados se conservan aunque el elemento original del Backlog se edite después.

## 2. Inputs

- `projectID`: UUID de un proyecto existente, identificado en la ruta.
- `sprintID`: entero positivo de un Sprint perteneciente al proyecto.
- `backlogItemID`: UUID de una historia del Backlog del mismo proyecto.

## 3. Expected Outputs

- La asignación devuelve HTTP 201 con el identificador del Sprint, el UUID original de la historia y una copia de sus datos de planificación.
- La consulta devuelve HTTP 200 con la lista de copias planificadas; sin asignaciones devuelve una lista vacía.
- Una edición posterior del elemento original no altera la copia consultada desde el Sprint.

## 4. Business Rules

- **BR-01:** solo se asigna una historia existente del Backlog al Sprint del mismo proyecto.
- **BR-02:** al asignar se guarda una copia de `id`, título, descripción, prioridad, estado, Story Points provisionales y definitivos, y criterios de aceptación.
- **BR-03:** la copia refleja los datos del Backlog al momento de la asignación y no se actualiza cuando cambia el elemento original.
- **BR-04:** una misma historia no puede duplicarse dentro del mismo Sprint.
- **BR-05:** asignar una historia no modifica el elemento original del Backlog.

## 5. Constraints

Contrato HTTP propuesto:

- `POST /api/projects/{projectID}/sprints/{sprintID}/backlog-items` con `{ "backlog_item_id": "<UUID>" }`.
- `GET /api/projects/{projectID}/sprints/{sprintID}/backlog-items` devuelve las copias planificadas.
- Los snapshots se persisten en una tabla propia de asignaciones del Sprint; consultar la planificación no depende de volver a leer los campos actuales del Backlog.
- El alcance no incluye quitar historias del Sprint, editar snapshots, cerrar Sprints ni trasladar historias incompletas (US-10).

## 6. Edge Cases

- Un Sprint sin historias devuelve `[]`.
- El Sprint no pertenece al proyecto de la ruta.
- La historia no existe o pertenece a otro proyecto.
- La historia ya está asignada al mismo Sprint.
- Editar en el Backlog título, prioridad, Story Points o criterios después de asignar no cambia el snapshot.

## 7. Error Conditions

- JSON inválido o UUID de historia vacío: HTTP 400.
- Proyecto, Sprint o historia inexistente/no perteneciente al proyecto: HTTP 404.
- Historia duplicada en el mismo Sprint: HTTP 409.
- Error de persistencia inesperado: HTTP 500 sin exponer detalles internos.

## 8. Acceptance Criteria

- **AC-01:** Una historia del proyecto asignada a un Sprint aparece al consultar su planificación.
- **AC-02:** El snapshot conserva los datos de planificación originales cuando se modifica posteriormente el elemento del Backlog.
- **AC-03:** La asignación rechaza historias inexistentes o de otro proyecto, sin modificar el Backlog ni la planificación.
- **AC-04:** Una historia no aparece duplicada en la planificación de un mismo Sprint.

## 9. BDD / Test Traceability

- AC-01 → BDD-US-08-01 → test de asignación y consulta.
- AC-02 → BDD-US-08-02 → test de snapshot inmutable tras edición del Backlog.
- AC-03 → BDD-US-08-03 → tests de historia/Sprint fuera del proyecto.
- AC-04 → BDD-US-08-04 → test de asignación duplicada.
