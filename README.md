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
## Integrantes y Roles Scrum

| Nombre y apellido |       GitHub     | Rol Scrum |
| :--- | :--- | :---|       
| Pablo Pratts      |     -            | **Product Architect** |
| Agustin Campos    | [@Aguscampos04]  | **Agile Enabler** |
| Loana Vargas      | [@Loana-Vargas]  | **Product Builder** |
| Delfina Carribero | [@DelfiCarribero]| **Product Builder** |
| Leandro Canton    | [@leacanton04]   | **Product Builder** |
| Bruno Carrieri    | [@Brunocarri]    | **Product Builder** |
