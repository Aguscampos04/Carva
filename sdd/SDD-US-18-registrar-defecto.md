# SDD-US-18 — Registrar defecto

**Issue:** [#18 — US-18](https://github.com/Aguscampos04/Carva/issues/18)  
**Épica:** EP-06 — Gestión de defectos  
**Requisito:** `docs/enunciado.md` §2.6  
**Dependencia:** US-05 — User Stories del Product Backlog  
**Trazabilidad:** §2.6 → EP-06 → US-18 → #18 → esta SDD → BDD → tests → Go/API/SQLite/UI.

## 1. Objetivo

Registrar un defecto del proyecto asociado a una User Story. El defecto puede comenzar sin Sprint; cuando se conoce el Sprint de detección, se guarda en el campo correspondiente, sin inventar un valor.

## 2. Inputs

La API registra un defecto desde el detalle de una historia:

`POST /api/projects/{projectID}/backlog/{backlogItemID}/defects`

| Campo | Obligatorio | Restricción |
| --- | --- | --- |
| `description` | Sí | Texto no vacío ni compuesto solo por espacios. Se conservan espacios internos y se normalizan los externos. |
| `severity` | Sí | Uno de `Baja`, `Media`, `Alta`, `Crítica`. |
| `status` | Sí | Uno de `Abierto`, `En progreso`, `Resuelto`, `Cerrado`. |
| `detection_sprint_id` | No | `null`/ausente o ID positivo de un Sprint del mismo proyecto. |
| `resolution_sprint_id` | No | `null`/ausente o ID positivo de un Sprint del mismo proyecto. Se incluye para persistir el campo exigido por §2.6; su actualización posterior se cubre en US-19. |

El proyecto y la historia relacionada se identifican por la ruta. Los identificadores son provistos por la aplicación, no por el cliente.

## 3. Expected Outputs

- Al crear, HTTP 201 devuelve el defecto con ID, proyecto, historia, descripción, severidad, estado y ambos campos de Sprint nullable.
- `GET /api/projects/{projectID}/backlog/{backlogItemID}/defects` devuelve los defectos de esa historia, incluso cuando no tienen Sprint, en una lista JSON (vacía si no hay defectos).
- La interfaz de la historia ofrece el formulario y consulta la lista.

## 4. Business Rules

- **BR-01:** Todo defecto se asocia a una User Story existente del mismo proyecto.
- **BR-02:** La descripción no puede estar vacía ni contener solo espacios.
- **BR-03:** Solo se aceptan los cuatro estados y severidades aprobados.
- **BR-04:** Los Sprints de detección/resolución son opcionales. Si se especifican, deben existir en el mismo proyecto.
- **BR-05:** No se asigna Sprint por defecto ni se deduce uno del Sprint actual o de la historia.
- **BR-06:** La creación no aplica transiciones de estado ni reglas adicionales sobre la relación entre estado y Sprints; US-19 gestiona la evolución del defecto.

## 5. Constraints

- Go mantiene las validaciones del dominio; la UI no es la única protección.
- SQLite persiste el defecto sin alterar la historia o el Sprint.
- Los estados y severidades se guardan con sus valores aprobados.
- La selección de historia limita el listado al mismo proyecto e historia.
- No incluye edición, eliminación, métricas ni reglas de transición de estados (US-19/US-22).

## 6. Edge Cases

- Ausencia de ambos Sprint IDs: se aceptan y persisten como `null`.
- Sprint de detección conocido y Sprint de resolución ausente: solo se persiste el de detección.
- Historia o Sprint de otro proyecto: se rechaza, no se asocia indirectamente.
- Proyecto/historia sin defectos: se responde una lista vacía.
- Descripción con espacios exteriores: se normalizan antes de persistir; espacios internos se conservan.

## 7. Error Conditions

- JSON inválido, campo requerido ausente, descripción vacía o valor de estado/severidad no aprobado: HTTP 400.
- Proyecto, User Story o Sprint inexistente/no perteneciente al proyecto: HTTP 404.
- Error de persistencia inesperado: HTTP 500 con mensaje genérico.

## 8. Acceptance Criteria

- **AC-01:** Puede registrarse un defecto con descripción, estado y severidad aprobados e historia relacionada, aunque no haya Sprint asociado inicialmente.
- **AC-02:** Si se conoce el Sprint de detección, queda registrado en el campo correspondiente; la ausencia inicial de Sprint no se reemplaza con uno ficticio.
- **AC-03 (DoD):** Se aceptan todos los estados y severidades enumerados en esta SDD; valores no aprobados se rechazan.
- **AC-04 (DoD):** La creación se prueba con y sin Sprint de detección.

Los campos `detection_sprint_id` y `resolution_sprint_id` son nullable. La asignación posterior de Sprints y las actualizaciones del estado pertenecen a US-19.
