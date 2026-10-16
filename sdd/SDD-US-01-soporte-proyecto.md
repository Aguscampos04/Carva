# SDD-US-01 — Soporte mínimo de proyectos

**Estado:** alcance aprobado por el responsable para habilitar US-07; contrato técnico sujeto a revisión del equipo.
**Trazabilidad:** §2.1 → EP-01 → US-01 / Issue #1.

## Objective

Crear y consultar un proyecto con los datos mínimos necesarios para asociarle un Sprint. Esta implementación acotada de US-01 no cubre miembros, fechas ni edición del proyecto.

## Inputs

- `name`: string obligatorio para crear un proyecto. Se eliminan espacios iniciales y finales; vacío se rechaza.
- `projectID`: entero positivo para consultar un proyecto.

## Expected Outputs

- Al crear un proyecto, se devuelve su identificador, nombre y estado inicial `Planificado`.
- Al consultar un proyecto existente, se devuelven esos mismos datos.

## Business Rules

- BR-01: Todo proyecto creado comienza con estado `Planificado`.
- BR-02: El nombre no puede quedar vacío después de quitar espacios exteriores.
- BR-03: Un identificador que no corresponda a un proyecto produce un resultado no encontrado.

## Constraints

- La persistencia es SQLite y la lógica central está en Go.
- Contrato REST propuesto: `POST /api/projects` con `{ "name": string }` y `GET /api/projects/{projectID}`.
- La respuesta de creación es HTTP 201; la consulta existente es HTTP 200.
- Se persisten únicamente ID, nombre y estado, como soporte mínimo de US-07.
- El servidor local escucha en `:8080`.

## Edge Cases

- Nombres con espacios exteriores se normalizan antes de guardar.
- Un nombre compuesto solo por espacios se rechaza.
- Un proyecto creado sin Sprint conserva estado `Planificado`.

## Error Conditions

- JSON inválido, nombre ausente/vacío o ID inválido: HTTP 400.
- Proyecto inexistente: HTTP 404.
- Error de persistencia inesperado: HTTP 500 sin exponer detalles internos.

## Acceptance Criteria

- AC-01: Un proyecto creado con nombre válido puede consultarse con los datos ingresados y estado `Planificado`.
- AC-02: Un nombre vacío o compuesto solo por espacios es rechazado.

## BDD / Test Traceability

- AC-01 → BDD-US-01-01 → `TestCreateProjectAndFirstSprintThroughAPI`.
- AC-02 → BDD-US-01-02 → `TestCreateProjectRejectsBlankName`.
