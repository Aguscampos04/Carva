# Carva

## Reportes de proyecto y Sprint (US-24)

Dentro de un proyecto, abrir **Reportes**, seleccionar **Proyecto completo** o un
Sprint y pulsar **Actualizar reporte** o **Descargar PDF**. La descarga se genera
en Go; no requiere imprimir desde el navegador ni instalar software adicional.

- `GET /api/projects/{projectID}/reports/project`: reporte JSON del proyecto.
- `GET /api/projects/{projectID}/reports/sprints/{sprintID}`: reporte JSON del Sprint.
- Las mismas rutas con sufijo `.pdf` descargan el documento PDF correspondiente.

Se incluyen datos generales, estado, miembros, historias, estimaciones, resultados,
esfuerzo activo, métricas y defectos según el alcance. Las métricas sin datos se
muestran como **No disponible**. El historial mantiene sus snapshots, con la
excepción aprobada de horas estimadas vigentes de US-16. El reporte de proyecto
incluye esfuerzo sin Sprint; el de Sprint solo esfuerzo explícitamente asociado.

El PDF puede consultarse en cualquier estado del proyecto; exportarlo no significa
finalizarlo ni aprobar su entrega. Las fuentes se leen al generar cada reporte,
sin promesa de snapshot atómico entre consultas concurrentes. Los caracteres
sin glifo en la fuente PDF aparecen como `[U+XXXX]`; acentos españoles se conservan.

SDD: [US-24](sdd/SDD-US-24-reportes-pdf.md). Evidencia técnica y TDD:
[US-24](docs/evidencia-ia-us-24.md). Licencias PDF/fuentes: [directorio](docs/licencias).

## Ejecutar la aplicación

Requiere Go 1.26 o posterior. Desde la raíz del repositorio:

```powershell
go run ./cmd/carva
```

La aplicación queda disponible en <http://localhost:8080>. La base se crea por
defecto en `.local/carva.db`; para usar otra ubicación, configurá
`CARVA_DB_PATH`.

La vista local de UX muestra **todos los proyectos creados** en Inicio y permite
abrirlos por nombre, crear otro o usar opcionalmente su UUID. `GET /api/projects`
lee la base existente; no depende de un historial del navegador.
Dentro del proyecto, el menú ofrece Resumen, Backlog, Sprints, Equipo y
Datos del proyecto. La URL conserva proyecto y sección al recargar.

En Backlog, usá búsqueda/filtros y abrí una historia para consultar Información,
Planning Poker, Esfuerzo o Defectos. Sprints permite crear, asignar historias,
confirmar la verificación de sus criterios, completar y cerrar; el Historial
conserva los resultados de cierre y distingue las horas estimadas vigentes.
El detalle del Sprint y cada historia se pueden plegar; guardar conserva la
posición de lectura. Los criterios quedan visibles dentro de la HU abierta. Cada historia
del Sprint muestra solo su título mientras está plegada. Dentro del Sprint,
**En proceso** y **Finalizadas** separan las historias según su verificación;
completar una HU la mueve a Finalizadas sin cambiar de sección.
**Cargar horas trabajadas** abre un formulario de alta en la HU, con el Sprint
compatible preseleccionado. **Cerrar** permite plegarlo y volver a abrirlo;
si hay cambios, solicita confirmación antes de descartarlos. La consulta,
comparación, edición y eliminación siguen en Backlog → HU → Esfuerzo.
**Miembro responsable** permite elegir un único miembro
opcional del proyecto y guardarlo; se conserva entre Sprints y no determina quién
puede registrar horas. En un Sprint cerrado se consulta el responsable actual.
El Inicio presenta los formularios antes del listado de proyectos.
Resumen alterna Proyecto/Sprint. La interfaz conserva las reglas de las SDD y
no incorpora reportes de US-24. La mejora se sigue en la
[issue #53](https://github.com/Aguscampos04/Carva/issues/53).
Referencia: [SDD de UX](sdd/SDD-UX-local.md).

La API de
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

## Edición y finalización del proyecto (US-02)

El detalle del proyecto permite guardar su nombre y fecha de inicio opcional.
Vaciar el inicio lo elimina. «Finalizar proyecto» solicita confirmación y usa los
datos ya guardados: cancelar no guarda el borrador. El servidor registra el día
actual en `America/Argentina/Buenos_Aires`, con formato `YYYY-MM-DD`.

- `PATCH /api/projects/{id}` — objeto parcial con `name` y/o `start_date`. Omitir conserva; `start_date: null` elimina. No se aceptan campos adicionales ni un objeto vacío.
- `POST /api/projects/{id}/finalize` — `{"confirmed":true}`. Repetir conserva la primera fecha de fin.
- `GET /api/projects/{id}` devuelve también `start_date` y `end_date`, con `null` cuando no hay registro.

Un proyecto Finalizado permite editar nombre/inicio y seguir registrando miembros
y creando/editando Product Backlog. No permite nuevos Sprints (`409`); conserva
los existentes y sus snapshots. No se implementan cierre ni traslado de Sprints.
Inicio posterior al fin se rechaza. La migración aditiva `000008` conserva datos
anteriores sin inventar fechas. Contención de escritura superior a cinco segundos
devuelve `503`; ante pérdida de respuesta, consultar antes de reintentar.

Ver [SDD-US-02](sdd/SDD-US-02-modificar-proyecto.md),
[BDD-US-02](bdd/BDD-US-02-modificar-proyecto.feature) y
[evidencia TDD y validación](docs/evidencia-ia-us-02.md). Pruebas web sin dependencias:

```powershell
node --test tests/web/members.test.cjs tests/web/projects.test.cjs
```

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

## Consultar Sprints anteriores (US-11)

En **Sprints anteriores**, ingresá el UUID de un proyecto existente y pulsá
**Consultar Sprints cerrados**. Elegí un Sprint de la lista para ver su Goal,
planificación y resultados registrados al cierre. La sección también toma el
proyecto creado/consultado en la pantalla actual. No hay que crear otro proyecto
para consultar un UUID conocido.

- `GET /api/projects/{projectID}/sprints/closed`: lista solo cerrados, ordenados por ID; `[]` si no hay ninguno y `404` si el proyecto no existe.
- El detalle y los snapshots se leen con las rutas existentes de Sprint y `/backlog-items`.
- Una historia trasladada y completada después conserva el resultado anterior en
  el Sprint cerrado. La consulta no modifica datos ni reconstruye snapshots desde
  el Backlog vigente. Los Sprints se crean/cierran mediante la API existente.

Ver [SDD-US-11](sdd/SDD-US-11-consultar-sprints-anteriores.md),
[BDD-US-11](bdd/BDD-US-11-consultar-sprints-anteriores.feature) y
[evidencia y limitaciones de validación](docs/evidencia-ia-us-11.md).

```powershell
go test ./internal/application/sprints ./tests/integration -run 'TestListClosed|TestUS11|TestCloseSprint|TestUS10' -count=1
node --test tests/web/*.test.cjs
```

## Horas estimadas y esfuerzo real (US-16)

La implementación se entrega en [PR #45](https://github.com/Aguscampos04/Carva/pull/45)
para aceptación final. Reglas y contratos: [SDD](sdd/SDD-US-16-registrar-horas-esfuerzo.md),
[46 escenarios BDD](bdd/BDD-US-16-registrar-horas-esfuerzo.feature) y
[evidencia TDD, matriz y verificaciones](docs/evidencia-ia-us-16.md).

Para probar sin tocar la base habitual, en PowerShell:

```powershell
$env:CARVA_DB_PATH = '.local/prueba-us16.db'
go run ./cmd/carva
```

Abrí [Carva local](http://localhost:8080). Creá un proyecto y registrá al menos
un miembro. Creá una User Story con todos sus campos y **Horas estimadas**
(por ejemplo 2.50). En su tarjeta, **Registrar y consultar esfuerzo** abre
el detalle con total y formulario. Elegí el miembro por nombre y UUID, fecha,
actividad y horas; Sprint es opcional y solo ofrece asignaciones existentes.
Registrar 3.25 muestra real 3.25 y conserva estimación 2.50.

Desde el detalle podés editar una entrada o eliminarla con confirmación.
Cancelar no guarda ni elimina. Al cambiar historia se conserva Sprint compatible;
si no lo es, elegí otro válido o **Sin Sprint** explícitamente. Las historias
históricas sin estimación muestran **Pendiente de regularización**: pueden
recibir esfuerzo, pero editar sus datos exige completar estimación válida.

Se permite esfuerzo en Completadas, Sprints cerrados y proyectos Finalizados
con fecha de fin conocida, respetando inicio/fin y el día actual de Buenos Aires.
Un Finalizado histórico con fin desconocido permite consulta/baja y bloquea
alta/edición de esfuerzo. Cambiar inicio considera únicamente entradas activas.
El máximo diario es 24 horas por miembro/proyecto entre todas sus historias;
horas y Story Points son independientes. La migración se aplica al arrancar y
conserva datos históricos sin inventar estimaciones ni fechas.

API nueva: GET de detalle `/api/projects/{projectID}/backlog/{itemID}`;
POST/GET `.../backlog/{itemID}/effort`; GET `.../backlog/{itemID}/effort-context`;
GET/PUT `.../effort/{entryID}` y POST `.../effort/{entryID}/delete` con
`{"confirmed":true}`. POST/PUT de historias requieren ahora `estimated_hours`.
El contrato completo y cuerpos de esfuerzo están en la SDD.

Verificaciones reproducibles en PowerShell:

```powershell
gofmt -l .
go test ./... -count=1
go test ./... -count=1 '-coverpkg=./internal/...' '-coverprofile=.local/us16-coverage.out'
go tool cover '-func=.local/us16-coverage.out'
go vet ./...
go build ./...
node --test (Get-ChildItem tests/web/*.test.cjs).FullName
```

CI ejecuta las mismas verificaciones en Linux y conserva el perfil de cobertura
como artefacto. No se hace merge ni se cierra la Issue #16 automáticamente.

## Integrantes y Roles Scrum

| Nombre y apellido |       GitHub     | Rol Scrum |
| :--- | :--- | :---|       
| Pablo Pratts      |     -            | **Product Architect** |
| Agustin Campos    | [@Aguscampos04]  | **Agile Enabler** |
| Loana Vargas      | [@Loana-Vargas]  | **Product Builder** |
| Delfina Carribero | [@DelfiCarribero]| **Product Builder** |
| Leandro Canton    | [@leacanton04]   | **Product Builder** |
| Bruno Carrieri    | [@Brunocarri]    | **Product Builder** |
