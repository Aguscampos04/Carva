# Carva

## Ejecutar la aplicación

Requiere Go 1.26 o posterior. Desde la raíz del repositorio:

```powershell
go run ./cmd/carva
```

La aplicación queda disponible en <http://localhost:8080>. La base se crea por
defecto en `.local/carva.db`; para usar otra ubicación, configurá
`CARVA_DB_PATH`.

La pantalla inicial permite crear un proyecto y volver a consultarlo. La API de
US-01 expone `POST /api/projects` con `{"name":"..."}` y
`GET /api/projects/{id}`. El servidor genera el UUID y establece el estado
`Planificado`.
La API también permite crear un Sprint para un proyecto existente y consultar
su objetivo:

- `POST /api/projects/{projectID}/sprints` — body: `{ "goal": "Entregar el MVP" }`.
- `GET /api/projects/{projectID}/sprints/{sprintID}` — consulta el Sprint dentro de su proyecto.

`projectID` es el UUID devuelto al crear el proyecto; `sprintID` es un entero
positivo. Crear el primer Sprint cambia el estado del proyecto de `Planificado`
a `En progreso`.

Para US-08, las historias del Backlog pueden asignarse y consultarse en la
planificación de un Sprint:

- `POST /api/projects/{projectID}/sprints/{sprintID}/backlog-items` — body: `{ "backlog_item_id": "<UUID>" }`.
- `GET /api/projects/{projectID}/sprints/{sprintID}/backlog-items` — devuelve los snapshots planificados.

Cada asignación conserva los datos de la historia al momento de planificarla,
aunque después cambien en el Backlog.

## Miembros del proyecto (US-03)

Al crear o consultar un proyecto, su sección **Miembros del proyecto** permite
registrar un miembro ingresando únicamente su nombre y consultar los registrados.
Se recortan los espacios exteriores. Se permiten homónimos: cada alta recibe un
UUID independiente, visible como dato secundario debajo del nombre. Los miembros
de otros proyectos no se incluyen. No hay edición, eliminación ni autenticación.

- `POST /api/projects/{projectID}/members` — cuerpo: `{"name":"Ana Pérez"}`; devuelve `201` con `id`, `project_id` y `name`.
- `GET /api/projects/{projectID}/members` — devuelve `200` con una lista, o `[]` para un proyecto sin miembros.

El alta exige un único objeto JSON completo: rechaza tipos incorrectos, claves
duplicadas y contenido posterior. Las propiedades adicionales no controlan UUID
ni proyecto. Datos inválidos: `400`; proyecto inexistente: `404`; errores internos:
`500` con mensaje seguro. Si falla una solicitud sin confirmación, consultá la lista
antes de repetirla: un nuevo envío puede registrar otro homónimo.

Ver [SDD-US-03](sdd/SDD-US-03-registrar-miembros.md) y
[evidencia TDD / matriz BDD–tests](docs/evidencia-ia-us-03.md). Para verificar:

```powershell
gofmt -l .
go test ./...
go vet ./...
node --test tests/web/members.test.cjs
```

Las pruebas web usan el runner incluido en Node.js (verificadas con Node 20),
sin instalar paquetes. Los selectores para votos y esfuerzo se incorporarán en
US-12 y US-16; US-03 proporciona sus identidades estables.

## Integrantes y Roles Scrum

| Nombre y apellido |       GitHub     | Rol Scrum |
| :--- | :--- | :---|       
| Pablo Pratts      |     -            | **Product Architect** |
| Agustin Campos    | [@Aguscampos04]  | **Agile Enabler** |
| Loana Vargas      | [@Loana-Vargas]  | **Product Builder** |
| Delfina Carribero | [@DelfiCarribero]| **Product Builder** |
| Leandro Canton    | [@leacanton04]   | **Product Builder** |
| Bruno Carrieri    | [@Brunocarri]    | **Product Builder** |
