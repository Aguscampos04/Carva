# Arquitectura de Carva

## 1. Estado de las decisiones

El equipo aprobó las decisiones arquitectónicas de la sección 2 y la estructura de carpetas y paquetes de la sección 3 para Sprint 1. Esta estructura inicial ya está creada en el repositorio. Los detalles técnicos que el equipo no aprobó expresamente permanecen pendientes y no se consideran definidos por la creación del esqueleto.

## 2. Decisiones aprobadas

| Área | Decisión aprobada |
| --- | --- |
| Backend | Go para el backend, el núcleo y las reglas de negocio. |
| Frontend | HTML, CSS y JavaScript. |
| Comunicación | API REST entre frontend y backend. |
| Persistencia | SQLite. |
| Arquitectura | Separación en presentación, aplicación, dominio y persistencia. |
| Identificación | Los miembros se registran asociados a proyectos; inicialmente no se implementa autenticación compleja. |
| Pruebas | Testing de Go y pruebas de integración. |
| Go mínimo | `go 1.26.0` en `go.mod`. |
| Driver SQLite | `modernc.org/sqlite v1.60.1`, sin CGO. |
| Ruta de base local | `CARVA_DB_PATH`; por defecto `.local/carva.db`. |
| Migraciones | Archivos SQL versionados e incluidos mediante `embed`. |
| Bases de pruebas | Archivos temporales independientes por prueba de integración. |

Estas decisiones no definen rutas HTTP, esquemas JSON, tablas de negocio, columnas, reglas de validación no aprobadas ni un framework web.

## 3. Estructura de carpetas aprobada para Sprint 1

La siguiente estructura fue aprobada por el equipo y creada como base del repositorio. US-01 implementa el alta y consulta de proyectos; US-07 agrega creación y consulta de Sprints vinculados a esos proyectos.

```text
.
├── cmd/
│   └── carva/
│       └── main.go
├── internal/
│   ├── presentation/
│   │   └── httpapi/
│   │       └── doc.go
│   ├── application/
│   │   ├── projects/
│   │   ├── members/
│   │   ├── backlog/
│   │   └── sprints/
│   ├── domain/
│   │   ├── project/
│   │   ├── member/
│   │   ├── backlog/
│   │   └── sprint/
│   └── persistence/
│       └── sqlite/
│           ├── db.go
│           ├── migrate.go
│           └── migrations/
│               ├── embed.go
│               └── README.md
├── web/
│   └── static/
│       ├── css/
│       └── js/
├── bdd/
├── sdd/
├── tests/
│   └── integration/
│       └── sqlite_test.go
├── docs/
└── go.mod
```

El módulo Go es `github.com/Aguscampos04/Carva`. US-01 incorpora el primer caso de uso de proyectos, persistencia, API y UI en las ubicaciones aprobadas. US-07 reutiliza el proyecto identificado por UUID y agrega los Sprints; US-08 agrega la asignación de historias con snapshots de planificación. Las carpetas `src/`, `sdd/`, `bdd/` y `tests/` que ya existían se conservan; no se utiliza `src/` para el código Go.

Los paquetes iniciales son `presentation/httpapi`, `application/projects`, `application/members`, `application/backlog`, `application/sprints`, `domain/project`, `domain/member`, `domain/backlog`, `domain/sprint` y `persistence/sqlite`. US-01 implementa proyectos; US-03 implementa registro y listado de miembros por proyecto; US-05/US-06 implementan Product Backlog; US-07 implementa Sprints; US-08 conserva los elementos asignados como snapshots. Los paquetes de las historias restantes continúan como esqueletos.

US-12 agrega `application/planningpoker` y `domain/planningpoker` para las rondas y votos. SQLite mantiene su adaptador dentro de `persistence/sqlite`; la migración `000011_planning_poker.sql` persiste la ronda inicial, participantes y votos. La UI usa JavaScript existente y no define reglas de estimación.

US-15 reutiliza ese caso de uso de Planning Poker para registrar un acuerdo unánime de una ronda finalizada. SQLite valida el resultado y actualiza `backlog_items.story_points` en una transacción inmediata; no requiere una migración adicional. La interfaz ofrece la acción cuando la ronda finalizada no tiene diferencias y recarga el Backlog tras registrarla.

El paquete `persistence/sqlite` obtiene su ruta desde `CARVA_DB_PATH` y usa `.local/carva.db` si la variable no está definida o está vacía. `Open` valida la conexión y crea el directorio padre del archivo si hace falta. Las migraciones se almacenan como SQL versionado (`NNNNNN_descripcion.sql`) en `internal/persistence/sqlite/migrations/`, se embeben y se aplican en orden; `carva_schema_migrations` registra las versiones aplicadas. Las migraciones `000001_baseline.sql`, `000002_projects.sql`, `000003_backlog_items.sql`, `000004_sprints.sql` y `000005_sprint_backlog_items.sql` crean las estructuras para US-01, US-05/US-06, US-07 y US-08. Los Sprints se vinculan por UUID al proyecto y sus historias planificadas mantienen copias de los datos vigentes al asignarlas. La tabla de registro es metadato de infraestructura. `.local/` y los archivos SQLite locales están excluidos de Git.

`modernc.org/sqlite v1.60.1` declara compatibilidad desde Go 1.26.0. Sus dependencias transitivas quedan fijadas en `go.mod`/`go.sum`.

## 4. Responsabilidades por capa

### Presentación

- Recibir solicitudes HTTP y exponer la API REST.
- Traducir solicitudes y respuestas HTTP a los casos de uso de aplicación.
- Gestionar el transporte y los contratos HTTP documentados por cada SDD.
- Servir los recursos HTML, CSS y JavaScript embebidos desde el mismo proceso HTTP.
- No concentrar reglas centrales de negocio en handlers ni en JavaScript.

### Aplicación

- Coordinar casos de uso, por ejemplo crear proyectos, registrar miembros, administrar el Backlog y planificar Sprints.
- Orquestar operaciones y transacciones necesarias para cada caso de uso.
- Definir puertos o interfaces que necesiten los casos de uso, como repositorios.
- Traducir resultados del dominio a resultados consumibles por la presentación.

### Dominio

- Contener entidades, tipos y reglas de negocio independientes de HTTP, HTML, JavaScript y SQLite.
- Proteger invariantes que se definan en las SDD y decisiones funcionales aprobadas.
- Ser testeable sin iniciar servidor ni conectarse a una base de datos.

### Persistencia

- Implementar los puertos de almacenamiento usando SQLite.
- Traducir entre los datos persistidos y los tipos del dominio o de aplicación.
- Gestionar apertura de la base, consultas, escrituras, transacciones y migraciones.
- No decidir reglas funcionales ni exponer detalles SQL al dominio.

### Composición de dependencias

`cmd/carva` es el punto de composición: inicializa SQLite, aplica migraciones, construye los repositorios y casos de uso de proyectos y Sprints, configura la API HTTP y sirve el frontend.

## 5. Dependencias permitidas

La dirección prevista es:

```text
Frontend (HTML/CSS/JavaScript) → API REST / Presentación → Aplicación → Dominio
                                                        ↑
                                             Persistencia SQLite

cmd/carva compone Presentación, Aplicación y Persistencia.
```

- La presentación puede depender de la aplicación.
- La aplicación puede depender del dominio y de los puertos que define.
- Persistencia puede depender de esos puertos y de los tipos de dominio necesarios para implementarlos.
- El dominio no depende de las otras capas ni de detalles tecnológicos.
- La aplicación no importa directamente handlers HTTP ni implementaciones SQLite.
- El punto de composición puede importar las implementaciones concretas para conectarlas.
- Las dependencias no deben formar ciclos. El frontend consume el contrato HTTP y no accede directamente a SQLite.

## 6. Estrategia de pruebas propuesta

- **Pruebas unitarias de dominio:** reglas e invariantes definidas en SDD, sin infraestructura externa.
- **Pruebas unitarias de aplicación:** casos de uso con adaptadores de prueba para los puertos.
- **Pruebas de integración:** verificar conjuntamente los límites acordados, por ejemplo persistencia SQLite y API HTTP, usando bases de archivos temporales independientes dentro de `t.TempDir()`. Las pruebas aplican las migraciones y descartan sus archivos al finalizar; no utilizan la ruta configurada para desarrollo.
- Los escenarios BDD deberán derivarse de criterios de aceptación. Las pruebas deben mantener trazabilidad con User Story, SDD y comportamiento.
- El driver elegido es `modernc.org/sqlite`; no requiere CGO. No se agrega un framework BDD ni una herramienta externa de migraciones.

## 7. Decisiones técnicas pendientes

Antes de implementar los aspectos correspondientes, el equipo debe revisar o decidir:

1. Contratos HTTP, atributos de dominio, esquema y reglas de validación de historias futuras; el contrato de US-01 está documentado en `sdd/SDD-US-01-crear-proyecto.md`.
2. Configuración general, logging, apagado ordenado y despliegue.
3. Si se requiere autenticación en una etapa posterior. Para US-03 se aprobaron nombre como único dato de alta, UUID automático, homónimos permitidos, registro independiente por proyecto y selección desde su lista en funcionalidades futuras, sin autenticación ni permisos. Véanse D-01 a D-07 en [contexto del proyecto](contexto-proyecto.md#miembros-del-proyecto--us-03--issue-3). El responsable autorizó implementar la [SDD-US-03](../sdd/SDD-US-03-registrar-miembros.md) con los ajustes de revisión. Los paquetes de miembros implementan alta/listado por proyecto, conectados a REST y a su sección web. La migración `000006_members.sql` agrega `members`; `OpenPath` activa claves foráneas en cada conexión mediante el DSN del driver. La revisión y aprobación final de la historia por el equipo siguen pendientes.

US-03 incorpora además `000007_members_name_storage.sql`: conserva los registros y mide el nombre por bytes para que el CHECK no rechace caracteres admitidos por D-06 (SQLite `length(TEXT)` se detiene en NUL). Las migraciones anteriores se conservan sin modificaciones.

Estas decisiones pendientes no deben convertirse en reglas funcionales por inferencia. Deben resolverse mediante SDD y aprobación del equipo cuando afecten el comportamiento o la arquitectura.

## 8. Trazabilidad

### Incremento US-16 autorizado

El responsable aprobó las 31 decisiones y PT-01…PT-06 de [SDD-US-16](../sdd/SDD-US-16-registrar-horas-esfuerzo.md). Se agregan `domain/hours` (centésimas exactas), `domain/effort` y `application/effort`, con adaptador SQLite, REST y formulario embebido desde el detalle del Product Backlog. `cmd/carva` compone el servicio con reloj y repositorio existentes; no se agregan frameworks ni dependencias.

La migración aditiva `000012_effort.sql` conserva NULL históricos de estimación y agrega `effort_entries` con relaciones de proyecto/miembro/historia/Sprint y marca de baja lógica. Los totales se calculan desde activos, sin multiplicarlos por snapshots. `BacklogRepository.Mutate` lee el registro vigente, actualiza campos editables y propaga exclusivamente horas estimadas a asignaciones abiertas/cerradas en la misma transacción inmediata. Acuerdos y resultados no se sobrescriben.

Esfuerzo y cambios de fechas se coordinan mediante `withImmediate`: fechas, relaciones y cupo se validan sobre datos vigentes. El cambio de inicio considera solo activos (D-31); el Finalizado con fin NULL conserva sus datos y solo permite consulta/baja de esfuerzo (D-30). La API de esfuerzo se agrega con `httpapi.WithEffort` sin sustituir los constructores previos. Evidencia de TDD, compatibilidad, concurrencia y CI: [US-16](evidencia-ia-us-16.md). Revisión final e integración pendientes del responsable.

### Incremento US-18 — Registro de defectos

La [SDD-US-18](../sdd/SDD-US-18-registrar-defecto.md) deriva el registro de
`domain/defect` y `application/defects`. SQLite aplica `000013_defects.sql` con
FK compuestas para mantener el proyecto, la User Story y los Sprints en el mismo
proyecto. Los Sprints de detección y resolución son anulables; sus valores
aprobados se validan en Go y la API `WithDefects` permite registrar y consultar
por historia. La interfaz ofrece la operación desde cada tarjeta del Backlog;
la edición del estado y la asociación posterior de Sprints queda en US-19.
Trazabilidad: Issue #18 → SDD → BDD → tests Go/Node → dominio, SQLite, REST y UI.

### Incremento US-20 — Métricas de Story Points y avance

La [SDD-US-20](../sdd/SDD-US-20-metricas-story-points-avance.md) define las
fórmulas aprobadas para Story Points planificados/completados, velocidad y
porcentaje de historias terminadas. El dominio calcula desde snapshots de
`sprint_backlog_items`; SQLite los lee sin migraciones y `httpapi.WithMetrics`
expone resultados por Sprint y proyecto. El frontend muestra valores faltantes
como «No disponible» y permite seleccionar un Sprint. Los agregados suman
snapshots, por lo que una historia trasladada cuenta en cada Sprint en el que
fue planificada. Trazabilidad: Issue #20 → SDD → BDD → tests Go/Node → dominio,
SQLite, REST y UI. Véase [evidencia de US-20](evidencia-ia-us-20.md).
### Incremento US-02 autorizado el 2026-10-20

La [SDD-US-02](../sdd/SDD-US-02-modificar-proyecto.md) y PT-01 a PT-05 fueron
aprobadas para implementar. `Project` incorpora fechas de calendario opcionales;
Go mantiene las reglas, y aplicación obtiene el día del reloj del servidor en
`America/Argentina/Buenos_Aires`. PATCH y finalización se exponen en la misma API
y el frontend embebido incorpora edición y confirmación manual.

El puerto `projects.Repository.Mutate` ejecuta una transformación de dominio
repetible sobre el registro vigente. La creación de Sprint recibe igualmente una
transición de aplicación para ejecutarla después de leer el proyecto vigente.
SQLite coordina ambas operaciones con `BEGIN IMMEDIATE` en una conexión reservada;
reintenta BUSY/LOCKED solo tras rollback confirmado, con espera total acotada a
cinco segundos. No reintenta errores de commit de resultado incierto. No hay
dependencia de SQLite desde dominio o aplicación.

La migración `000008_projects_dates.sql` añade `start_date` y `end_date` anulables,
sin cambiar las migraciones anteriores ni inferir fechas históricas. Finalizar
solo actualiza el proyecto: los Sprints y snapshots se conservan; su nueva
creación se rechaza en Finalizado. Miembros y Backlog conservan sus operaciones.
La revisión final del equipo permanece pendiente; véase la [evidencia](evidencia-ia-us-02.md).

La arquitectura soporta, pero no reemplaza, la cadena exigida por el proyecto:

```text
Requisito → Epic → User Story / Issue → SDD → Acceptance Criteria → BDD → Test → Go → Commit / PR
```

Cada implementación deberá derivarse de sus artefactos aprobados y conservar evidencia verificable del proceso TDD y de las pruebas ejecutadas.
