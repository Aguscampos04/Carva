# Arquitectura de Carva

## 1. Estado de las decisiones

El equipo aprobó las decisiones arquitectónicas de la sección 2. Este documento registra esas decisiones y propone una organización de carpetas para revisión. La estructura y los detalles técnicos que el equipo no aprobó expresamente permanecen como propuestas o pendientes.

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

Estas decisiones no definen rutas HTTP, esquemas JSON, tablas, columnas, reglas de validación no aprobadas, ni una biblioteca o framework particular.

## 3. Estructura de carpetas propuesta

La siguiente estructura es una propuesta para revisión. No se crearon estos directorios ni se aprobaron todavía los nombres definitivos de carpetas y paquetes.

```text
.
├── cmd/
│   └── carva/
│       └── main.go
├── internal/
│   ├── presentation/
│   │   └── httpapi/
│   │       ├── handlers/
│   │       ├── requests/
│   │       ├── responses/
│   │       └── router.go
│   ├── application/
│   │   ├── projects/
│   │   ├── members/
│   │   ├── backlog/
│   │   ├── sprints/
│   │   └── ports/
│   ├── domain/
│   │   ├── project/
│   │   ├── member/
│   │   ├── backlog/
│   │   └── sprint/
│   └── persistence/
│       └── sqlite/
│           ├── migrations/
│           └── repositories/
├── web/
│   └── static/
│       ├── index.html
│       ├── css/
│       └── js/
├── bdd/
├── sdd/
├── tests/
│   └── integration/
├── docs/
├── go.mod
└── README.md
```

El repositorio contiene actualmente directorios vacíos `src/`, `sdd/`, `bdd/` y `tests/`. No contiene código de aplicación ni un `go.mod` rastreado. La propuesta ubica `cmd/` e `internal/` en la raíz del módulo Go y no requiere conservar `src/`; el equipo deberá aprobar esa organización antes de crearla o mover contenido.

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
- **Pruebas de integración:** verificar conjuntamente los límites acordados, por ejemplo persistencia SQLite y API HTTP, usando una base controlada para pruebas.
- Los escenarios BDD deberán derivarse de criterios de aceptación. Las pruebas deben mantener trazabilidad con User Story, SDD y comportamiento.
- No se selecciona todavía un framework BDD, driver SQLite ni herramienta adicional de integración.

## 7. Decisiones técnicas pendientes

Antes de implementar los aspectos correspondientes, el equipo debe revisar o decidir:

1. Organización definitiva de carpetas, paquetes y módulo Go (`go.mod`).
2. Rutas, métodos, versionado, formatos JSON, códigos HTTP, errores y validaciones de la API REST.
3. Atributos, identificadores, invariantes y relaciones detalladas del modelo de dominio, manteniendo las reglas funcionales aprobadas.
4. Driver y versión de SQLite, esquema, migraciones, transacciones, restricciones e índices.
5. Cómo se sirven los recursos estáticos y cómo configuran su URL de API los distintos entornos.
6. Configuración, logging, apagado ordenado y despliegue.
7. Librerías y mecanismos concretos para pruebas de integración y aislamiento de dependencias.
8. Si se requiere autenticación en una etapa posterior. Para la etapa inicial está aprobada la ausencia de autenticación compleja; el modo de identificar miembros dentro de cada proyecto aún debe especificarse.

Estas decisiones no deben convertirse en reglas funcionales por inferencia. Deben resolverse mediante SDD y aprobación del equipo cuando afecten el comportamiento o la arquitectura.

## 8. Trazabilidad

La arquitectura soporta, pero no reemplaza, la cadena exigida por el proyecto:

```text
Requisito → Epic → User Story / Issue → SDD → Acceptance Criteria → BDD → Test → Go → Commit / PR
```

Cada implementación deberá derivarse de sus artefactos aprobados y conservar evidencia verificable del proceso TDD y de las pruebas ejecutadas.
