# SDD-US-07 — Crear Sprint y definir Sprint Goal

**Estado:** alcance aprobado por el responsable; contrato técnico sujeto a revisión del equipo.
**Trazabilidad:** §2.3 y decisión aprobada de estado del proyecto → EP-03 → US-07 / Issue #7.
**Dependencia:** US-01 (soporte mínimo de proyectos descrito en SDD-US-01).

## Objective

Crear un Sprint asociado a un proyecto y permitir consultar su Sprint Goal. Al crear el primer Sprint, un proyecto `Planificado` pasa a `En progreso`.

## Inputs

- `projectID`: UUID no vacío de un proyecto existente, identificado en la ruta.
- `goal`: string obligatorio, enviado en el cuerpo JSON. Se eliminan espacios iniciales y finales; vacío se rechaza.
- `sprintID`: entero positivo para consultar un Sprint dentro del proyecto indicado.

## Expected Outputs

- La creación devuelve HTTP 201 y los campos `id`, `projectID` y `goal`.
- La consulta devuelve HTTP 200 con esos campos.
- Al crear el primer Sprint, el estado del proyecto se devuelve como `En progreso` en consultas posteriores.

## Business Rules

- BR-01: Cada Sprint pertenece a un proyecto existente.
- BR-02: Un Sprint debe tener un Sprint Goal no vacío.
- BR-03: La creación del Sprint y la transición del proyecto se guardan de forma atómica.
- BR-04: Si el proyecto está `Planificado`, su creación del primer Sprint lo cambia a `En progreso`.
- BR-05: Si el proyecto ya está `En progreso`, crear otro Sprint conserva ese estado.
- BR-06: Un Sprint solo puede consultarse dentro del proyecto al que pertenece.

## Constraints

- Contrato REST propuesto: `POST /api/projects/{projectID}/sprints` con `{ "goal": string }` y `GET /api/projects/{projectID}/sprints/{sprintID}`. `projectID` reutiliza el UUID emitido por US-01; el identificador interno del Sprint es entero.
- El proyecto y los Sprints se persisten en SQLite; el alta de Sprint y la transición de estado forman una transacción.
- La API local se sirve en `:8080`.
- No se incluyen cierre de Sprint, asignación de historias, fechas, nombre/numero de Sprint ni frontend.

## Edge Cases

- Un segundo Sprint de un proyecto `En progreso` no vuelve a cambiar el estado.
- Un Sprint de otro proyecto no se devuelve al consultar con un `projectID` diferente.
- Un `goal` con espacios exteriores se normaliza; uno vacío después de normalizar se rechaza.

## Error Conditions

- JSON inválido, Sprint ID no positivo o Sprint Goal vacío: HTTP 400.
- Proyecto inexistente o Sprint inexistente/no perteneciente al proyecto: HTTP 404.
- Error de persistencia inesperado: HTTP 500 sin exponer detalles internos.

## Acceptance Criteria

- AC-01: Un Sprint creado queda asociado al proyecto y permite consultar su Sprint Goal.
- AC-02: Al crear el primer Sprint de un proyecto `Planificado`, el proyecto pasa a `En progreso`.
- AC-03: La creación de Sprint rechaza proyecto inexistente y Sprint Goal vacío sin persistir cambios parciales.

## BDD / Test Traceability

- AC-01 → BDD-US-07-01 → `TestCreateProjectAndFirstSprintThroughAPI`.
- AC-02 → BDD-US-07-02 → `TestCreateProjectAndFirstSprintThroughAPI`.
- AC-03 → BDD-US-07-03/04 → `TestInvalidSprintCreationDoesNotChangeProject` and `TestCreateSprintRejectsUnknownProject`.
