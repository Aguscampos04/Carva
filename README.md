# Carva
## Integrantes y Roles Scrum

| Nombre y apellido |       GitHub     | Rol Scrum |
| :--- | :--- | :---|       
| Pablo Pratts      |     -            | **Product Architect** |
| Agustin Campos    | [@Aguscampos04]  | **Agile Enabler** |
| Loana Vargas      | [@Loana-Vargas]  | **Product Builder** |
| Delfina Carribero | [@DelfiCarribero]| **Product Builder** |
| Leandro Canton    | [@leacanton04]   | **Product Builder** |
| Bruno Carrieri    | [@Brunocarri]    | **Product Builder** |

## API inicial de proyectos y Sprints

Ejecutar `go run ./cmd/carva` inicia la API local en `http://localhost:8080` y crea/actualiza la base SQLite configurada por `CARVA_DB_PATH` (por defecto `.local/carva.db`).

- `POST /api/projects` — body: `{ "name": "Carva Demo" }`.
- `GET /api/projects/{projectID}` — consulta nombre y estado del proyecto.
- `POST /api/projects/{projectID}/sprints` — body: `{ "goal": "Entregar el MVP" }`.
- `GET /api/projects/{projectID}/sprints/{sprintID}` — consulta el Sprint dentro de su proyecto.

Las respuestas usan JSON. La creación responde `201`; las consultas responden `200`; entradas inválidas responden `400`; recursos inexistentes responden `404`.
