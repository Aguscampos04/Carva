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

Estas decisiones no definen un contrato REST general ni un framework web. Para el alcance mínimo de US-01/US-07, las SDD correspondientes proponen rutas, payloads y tablas; el equipo debe revisar esos contratos antes de extenderlos.

## 3. Estructura de carpetas aprobada para Sprint 1

La siguiente estructura fue aprobada por el equipo y creada como base del repositorio. La implementación inicial de proyectos y Sprints se especifica en `sdd/SDD-US-01-soporte-proyecto.md` y `sdd/SDD-US-07-crear-sprint.md`.

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

El módulo Go es `github.com/Aguscampos04/Carva`. `cmd/carva/main.go` compone la API inicial y la persistencia SQLite para proyectos y Sprints. Las carpetas `src/`, `sdd/`, `bdd/` y `tests/` que ya existían se conservan; no se utiliza `src/` para el código Go.

Los paquetes iniciales son `presentation/httpapi`, `application/projects`, `application/members`, `application/backlog`, `application/sprints`, `domain/project`, `domain/member`, `domain/backlog`, `domain/sprint` y `persistence/sqlite`. Los paquetes de miembros y Backlog aún no implementan comportamiento; proyectos y Sprints cubren el alcance mínimo descrito en sus SDD. Los directorios aún vacíos se conservan con `.gitkeep`.

El paquete `persistence/sqlite` obtiene su ruta desde `CARVA_DB_PATH` y usa `.local/carva.db` si la variable no está definida o está vacía. `Open` valida la conexión y crea el directorio padre del archivo si hace falta. Las migraciones se almacenan como SQL versionado (`NNNNNN_descripcion.sql`) en `internal/persistence/sqlite/migrations/`, se embeben y se aplican en orden; `carva_schema_migrations` registra las versiones aplicadas. `000001_baseline.sql` conserva la línea de base de infraestructura y `000002_projects_sprints.sql` crea las tablas mínimas de proyectos y Sprints. `.local/` y los archivos SQLite locales están excluidos de Git.

`modernc.org/sqlite v1.60.1` declara compatibilidad desde Go 1.26.0. Su módulo incorpora dependencias transitivas, entre ellas `modernc.org/libc`, `modernc.org/mathutil`, `modernc.org/fileutil` y `golang.org/x/sys`; sus versiones quedan fijadas en `go.mod`/`go.sum`.

## 4. Responsabilidades por capa

### Presentación

- Recibir solicitudes HTTP y exponer la API REST.
- Traducir solicitudes y respuestas HTTP a los casos de uso de aplicación.
- Gestionar el transporte, los formatos de entrada y salida y los códigos HTTP cuando se definan sus contratos.
- Servir o entregar los recursos HTML, CSS y JavaScript según la estrategia que apruebe el equipo.
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
- Gestionar apertura de la base, consultas, escrituras, transacciones y migraciones, una vez definidas sus estrategias.
- No decidir reglas funcionales ni exponer detalles SQL al dominio.

### Composición de dependencias

`cmd/carva` funcionaría como punto de composición: carga configuración, inicializa SQLite, construye adaptadores y casos de uso, configura la API HTTP e inicia el servidor. Esta función describe la propuesta de estructura y no implica que esa implementación ya exista.

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

1. Contratos REST adicionales a los endpoints mínimos propuestos en SDD-US-01 y SDD-US-07: versionado, formatos, códigos HTTP, errores y validaciones.
2. Atributos, identificadores, invariantes y relaciones del modelo de dominio más allá de los mínimos de proyectos y Sprints descritos en esas SDD.
3. Esquema de negocio SQLite más allá de las tablas mínimas de proyectos y Sprints; restricciones, índices y política general de concurrencia/transacciones.
4. Cómo se sirven los recursos estáticos y cómo configuran su URL de API los distintos entornos.
5. Configuración general, logging, apagado ordenado y despliegue.
6. Si se requiere autenticación en una etapa posterior. Para la etapa inicial está aprobada la ausencia de autenticación compleja; el modo de identificar miembros dentro de cada proyecto aún debe especificarse.

Estas decisiones no deben convertirse en reglas funcionales por inferencia. Deben resolverse mediante SDD y aprobación del equipo cuando afecten el comportamiento o la arquitectura.

## 8. Trazabilidad

La arquitectura soporta, pero no reemplaza, la cadena exigida por el proyecto:

```text
Requisito → Epic → User Story / Issue → SDD → Acceptance Criteria → BDD → Test → Go → Commit / PR
```

Cada implementación deberá derivarse de sus artefactos aprobados y conservar evidencia verificable del proceso TDD y de las pruebas ejecutadas.
