# Contexto del Proyecto

## 1. Identificación

**Proyecto:** Aplicación web para gestión, estimación, seguimiento y medición de proyectos de software.

**Materia:** Ingeniería y Calidad de Software.

**Año:** 2026.

## 2. Objetivo del proyecto

Desarrollar una aplicación web que permita gestionar proyectos de software y sus elementos asociados, incluyendo Product Backlog, Sprints, estimaciones, esfuerzo, defectos y métricas.

La aplicación deberá cumplir con los requerimientos funcionales y metodológicos establecidos en el enunciado oficial.

## 3. Tipo de solución

Se decidió desarrollar una **aplicación web**.

Esta decisión ya está tomada y no debe ser modificada por el agente salvo que el equipo lo solicite explícitamente.

## 4. Tecnologías y prácticas obligatorias

El proyecto deberá utilizar:

- Go (Golang) para el núcleo de la aplicación y las reglas de negocio.
- HTML, CSS y JavaScript para el frontend.
- Una API REST para la comunicación entre frontend y backend.
- SQLite como base de datos.
- Separación en las capas de presentación, aplicación, dominio y persistencia.
- Scrum como marco de trabajo.
- SDD para especificar funcionalidades antes de implementarlas.
- BDD para definir el comportamiento mediante escenarios.
- TDD para desarrollar y validar las reglas y funcionalidades mediante pruebas.
- Testing de Go y pruebas de integración.
- Git para el control de versiones.
- Herramientas de Inteligencia Artificial como soporte al desarrollo.

## 5. Principios de trabajo

El proyecto se desarrollará de forma incremental.

No se deberá implementar toda la aplicación de una sola vez.

El trabajo se organizará mediante:

1. Sprint 0 — Preparación.
2. Sprint 1 — MVP.
3. Sprint 2 — Interfaz.
4. Sprint 3 — Funcionalidad y Calidad.
5. Sprint 4 — Cierre y entrega final.

El avance entre etapas deberá realizarse de forma controlada y con revisión del equipo.

## 6. Arquitectura

El equipo aprobó oficialmente la arquitectura descrita en [`docs/arquitectura.md`](arquitectura.md):

- Backend y núcleo de negocio en Go.
- Frontend con HTML, CSS y JavaScript.
- Comunicación entre frontend y backend mediante una API REST.
- Persistencia en SQLite.
- Separación de responsabilidades en presentación, aplicación, dominio y persistencia.
- Registro de miembros asociados a cada proyecto, inicialmente sin autenticación compleja.
- Pruebas de Go y pruebas de integración.

La estructura de carpetas y paquetes para Sprint 1 descrita en ese documento fue aprobada por el equipo y creada en el repositorio. El módulo Go aprobado es `github.com/Aguscampos04/Carva`. Los contratos de API, el modelo y esquema de negocio, y las decisiones de despliegue aún deben precisarse antes de implementarlos.

La arquitectura deberá favorecer:

- Separación de responsabilidades.
- Testabilidad.
- Mantenibilidad.
- Trazabilidad.
- Evolución incremental.
- Claridad de las reglas de negocio.

## 7. Funcionalidades principales

La solución deberá contemplar, como mínimo:

- Gestión de proyectos.
- Gestión de miembros.
- Product Backlog.
- Gestión de Sprints.
- Estimación mediante Story Points.
- Planning Poker.
- Registro y comparación de esfuerzo.
- Gestión de defectos.
- Métricas.
- Dashboard.
- Reportes.

El detalle de cada funcionalidad deberá derivarse del enunciado oficial y posteriormente formalizarse mediante User Stories, SDD, BDD y TDD.

## 8. Gestión del proyecto

El proyecto utilizará un repositorio Git y un tablero Scrum.

Las unidades principales de trabajo serán:

- Epic.
- User Story.
- Issue.
- Sprint.
- Task, cuando sea necesario.

Cada elemento deberá mantener relaciones con los artefactos correspondientes.

## 9. Trazabilidad

Se deberá mantener la siguiente cadena de trazabilidad:

**Requerimiento → Epic → User Story → Issue → SDD → Criterios de aceptación → BDD → Test → Código → Commit/PR**

No será necesario que todas las decisiones se resuelvan de antemano.

La trazabilidad se construirá progresivamente durante el desarrollo.

Como mínimo, una funcionalidad deberá demostrar el recorrido completo durante la presentación final.

## 10. Uso del agente de IA

El agente será utilizado como asistente técnico y de gestión del proyecto.

Podrá colaborar en:

- Análisis de requerimientos.
- Identificación y refinamiento de User Stories.
- Creación y actualización de Issues.
- Elaboración de especificaciones SDD.
- Elaboración de escenarios BDD.
- Diseño y seguimiento de TDD.
- Generación y revisión de código Go.
- Generación y revisión de tests.
- Análisis de cobertura.
- Revisión de calidad.
- Control de trazabilidad.
- Documentación.
- Organización Scrum.

El agente no reemplaza la responsabilidad del equipo.

Todo resultado generado por IA deberá ser comprendido, revisado y validado antes de incorporarse al proyecto.

## 11. Regla de aprobación

El agente podrá proponer soluciones, estructuras, historias, especificaciones, código o modificaciones.

Sin embargo, no deberá considerar una propuesta como decisión definitiva hasta que el equipo la apruebe cuando dicha decisión afecte significativamente:

- Arquitectura.
- Requerimientos.
- Reglas de negocio.
- Alcance.
- Tecnologías.
- Organización del proyecto.

## 12. Estado actual

### Decisiones tomadas

- El proyecto será una aplicación web.
- El backend, núcleo y reglas de negocio se implementarán en Go.
- El frontend utilizará HTML, CSS y JavaScript.
- El frontend y el backend se comunicarán mediante una API REST.
- La base de datos será SQLite.
- La aplicación se organizará en capas de presentación, aplicación, dominio y persistencia.
- Los miembros se registrarán asociados a cada proyecto, inicialmente sin autenticación compleja.
- Se realizarán pruebas de Go y pruebas de integración.
- Se aprobó y creó la estructura inicial de carpetas y paquetes para Sprint 1, con el módulo `github.com/Aguscampos04/Carva`.
- Se aprobó Go 1.26.0 como versión mínima del módulo y `modernc.org/sqlite v1.60.1` como driver SQLite sin CGO.
- La ruta de la base se configura con `CARVA_DB_PATH`; el valor predeterminado local es `.local/carva.db`.
- Las migraciones serán SQL versionado e incorporado mediante `embed`; las pruebas de integración usarán bases temporales independientes.
- Git ignora `.local/` y los archivos de bases SQLite locales.
- Se utilizará Scrum.
- Se utilizarán SDD, BDD y TDD.
- Se utilizará Git.
- Se utilizará IA como soporte.
- El proyecto se desarrollará incrementalmente mediante Sprints.

### Product Backlog inicial aprobado

El responsable del proyecto confirma la aprobación de las **9 Épicas (EP-01 a EP-09)** y las **24 User Stories (US-01 a US-24)** como versión inicial del Product Backlog. Esta aprobación establece la base inicial del alcance planificado; el Product Backlog podrá evolucionar mediante decisiones posteriores aprobadas por el equipo, manteniendo la trazabilidad con los requisitos del enunciado.

### Decisiones funcionales aprobadas

Estas decisiones definen el comportamiento funcional del producto. Las fórmulas concretas de las métricas continúan pendientes y deberán documentarse en sus SDD.

1. **Estado del proyecto:** al crearlo queda **Planificado**; al crear su primer Sprint pasa a **En progreso**; pasa a **Finalizado** mediante la acción manual confirmada «Finalizar proyecto», nunca automáticamente por una fecha. Véanse las decisiones de US-02 / Issue #2.
2. **Estado de una User Story:** los estados son **Pendiente**, **En progreso** y **Completada**.
3. **Finalización de una User Story:** se considera completada cuando se verifican sus Acceptance Criteria.
4. **Cierre de un Sprint:** las User Stories incompletas se trasladan automáticamente al siguiente Sprint. Si existen historias incompletas, no se puede cerrar el Sprint hasta que haya un siguiente Sprint disponible. Las historias completadas no se trasladan. El Sprint cerrado conserva su historial y el conjunto originalmente planificado.
5. **Story Points iniciales:** una User Story puede crearse con Story Points provisionales.
6. **Actualización de Story Points:** el Story Point actual solo se actualiza cuando se registra una estimación acordada mediante Planning Poker.
7. **Planning Poker:** la votación es individual y oculta; las estimaciones se revelan simultáneamente; se discuten las diferencias; pueden realizarse nuevas rondas; y se registra el valor acordado.
8. **Diferencia de estimaciones:** existe cuando los valores votados no son todos iguales.
9. **Estimación acordada:** hay consenso cuando todos los participantes de la ronda votan el mismo valor. Ese valor se registra como Story Point definitivo. Una ronda finalizada puede revelar diferencias; sin consenso no se registra un Story Point definitivo y el equipo puede discutir y realizar otra ronda.
10. **Horas:** las horas estimadas se registran para la User Story durante la planificación; las horas reales provienen del registro de esfuerzo.
11. **Fórmulas de métricas:** las fórmulas de las métricas obligatorias se definirán explícitamente en sus SDD antes de implementarlas.
12. **Alcance temporal de las métricas:** se calculan por Sprint y pueden agregarse a nivel proyecto. El conjunto planificado de un Sprint debe preservarse aunque el Backlog cambie posteriormente.
13. **Datos faltantes en métricas:** cuando falten datos para una métrica, se muestra **«No disponible»**, nunca cero en reemplazo de esos datos.
14. **Defectos:** sus estados son **Abierto**, **En progreso**, **Resuelto** y **Cerrado**. Sus severidades son **Baja**, **Media**, **Alta** y **Crítica**.
15. **Defectos y Sprints:** los defectos pueden existir sin Sprint, asociarse posteriormente a uno y resolverse en un Sprint diferente al de detección.
16. **Reportes:** habrá un reporte de proyecto completo y un reporte de Sprint individual. Incluirán User Stories, Story Points, horas estimadas y reales, métricas y defectos.
17. **Reporte final:** debe poder generarse en PDF e incluir información general del proyecto, Sprints, User Stories, Story Points, horas, métricas, defectos y estado general.
18. **MVP e identificación:** se identificará a los miembros para Planning Poker y el registro de esfuerzo. No se implementará un sistema complejo de permisos en el MVP.
19. **Finalización de Planning Poker:** una ronda solo puede finalizarse cuando todos sus participantes registrados ya votaron.
20. **Participantes de Planning Poker:** todos los miembros registrados del proyecto participan en cada ronda; la ronda conserva la lista de participantes que tenía al abrirse.
21. **Corrección de voto:** un miembro puede reemplazar su propio voto mientras la ronda permanezca abierta; una ronda finalizada no admite cambios.
22. **Nuevas rondas de Planning Poker:** la siguiente ronda se abre solo después de finalizar la anterior, recibe el número consecutivo y conserva los votos/resultados revelados de las rondas previas.

Las especificaciones de US-12 a US-14 aplican estas decisiones y usan la identidad seleccionada de la lista del proyecto, como define D-05 de US-03. Trazabilidad: §2.4 → EP-04 → [US-14 / Issue #14](https://github.com/Aguscampos04/Carva/issues/14) → [SDD-US-14](../sdd/SDD-US-14-discutir-diferencias-nuevas-rondas.md) → [BDD-US-14](../bdd/BDD-US-14-discutir-diferencias-nuevas-rondas.feature).

### Modificación y finalización del proyecto — US-02 / Issue #2

Decisiones funcionales **aprobadas explícitamente por el responsable en la solicitud del 2026-10-20**. D-01 a D-10 tienen alcance local a US-02; no reemplazan los identificadores homónimos de US-03.

| ID | Decisión aprobada |
| --- | --- |
| D-01 | El proyecto se finaliza manualmente mediante un botón «Finalizar proyecto». Nunca automáticamente por una fecha. |
| D-02 | La fecha de inicio se registra manualmente, es opcional y puede modificarse. |
| D-03 | Al finalizar, el sistema registra automáticamente la fecha actual como fecha de finalización. |
| D-04 | Los proyectos finalizados permiten modificar su nombre y fecha de inicio, pero no su estado ni fecha de finalización. |
| D-05 | Se pueden finalizar proyectos `Planificado` o `En progreso`, incluso sin Sprints. |
| D-06 | La fecha de inicio no puede ser posterior a la fecha de finalización. Se permite finalizar sin fecha de inicio. |
| D-07 | Antes de finalizar, solicitar confirmación. Cancelar no provoca cambios. |
| D-08 | Un proyecto `Finalizado` no permite crear nuevos Sprints, pero sí consultar los existentes. |
| D-09 | Finalizar un proyecto no modifica sus Sprints existentes, aunque haya alguno abierto. |
| D-10 | En proyectos `Finalizados` se pueden seguir registrando miembros y creando/editando historias del Product Backlog. |

Trazabilidad: [enunciado](enunciado.md) §2.1 → EP-01 → [US-02 / Issue #2](https://github.com/Aguscampos04/Carva/issues/2) → D-01…D-10 → [SDD-US-02](../sdd/SDD-US-02-modificar-proyecto.md) → AC → [BDD-US-02](../bdd/BDD-US-02-modificar-proyecto.feature).

Estas decisiones resuelven la condición de finalización antes pendiente. La especificación inicial se preparó en revisión, sin código, migraciones ni tests TDD. El responsable aprobó posteriormente el 2026-10-20 la SDD, el BDD y PT-01 a PT-05 para implementar US-02, con ajustes previos: eliminar de BDD-07 el caso En progreso sin Sprints y hacer BDD-18 independiente de US-11. Formato y zona de fechas, contratos REST, concurrencia, modelo y migración quedan autorizados para esta implementación; no forman parte de D-01 a D-10 ni implican aprobación final de la historia. Se autorizan commits, push y PR hacia main, sin merge ni cierre de Issues.

Implementación realizada y verificación técnica registrada en [evidencia US-02](evidencia-ia-us-02.md), con revisión final del equipo pendiente. US-11 conserva su decisión pendiente de cierre de Sprint con historias incompletas en un proyecto Finalizado; US-02 no agrega reglas sobre ese cierre.

### Miembros del proyecto — US-03 / Issue #3

Decisiones funcionales **aprobadas explícitamente por el responsable en la solicitud del 2026-10-16**. Los identificadores D-01 a D-07 corresponden a US-03:

| ID | Decisión aprobada |
| --- | --- |
| D-01 | El único dato ingresado para registrar un miembro es su nombre. |
| D-02 | Cada miembro recibe automáticamente un UUID. |
| D-03 | Se permiten nombres iguales; la identidad depende del UUID. |
| D-04 | Los miembros se registran independientemente en cada proyecto, aunque representen a la misma persona. |
| D-05 | En las futuras funcionalidades se identificará al miembro seleccionándolo de la lista del proyecto. |
| D-06 | Rechazar nombres vacíos o solo con espacios. Eliminar espacios iniciales/finales y conservar el resto del nombre. |
| D-07 | Incorporar una sección de miembros con formulario de alta y lista de registrados. Sin edición, eliminación, autenticación ni permisos. |

Trazabilidad: enunciado §2.1 → EP-01 → [US-03 / Issue #3](https://github.com/Aguscampos04/Carva/issues/3) → [SDD-US-03](../sdd/SDD-US-03-registrar-miembros.md) → [BDD-US-03](../bdd/BDD-US-03-registrar-miembros.feature). La identificación habilita la futura atribución de votos (§2.4, US-12) y esfuerzo (§2.5, US-16), sin implementar esas historias en US-03.

El responsable autorizó el 2026-10-16 implementar la SDD y los escenarios con los ajustes de revisión: JSON estricto, integridad referencial efectiva en SQLite, claridad de BDD 07/17 y UUID completo presentado como dato secundario. Esta autorización de la etapa no implica aprobación final de la historia por el equipo ni evidencia de tests; la implementación deberá conservar evidencia TDD y entregarse por PR sin merge ni cierre de #3.

### Horas estimadas y esfuerzo real — US-16 / Issue #16

Decisiones funcionales **aprobadas explícitamente por el responsable**: D-01…D-29 en la solicitud del 2026-10-26 y D-30/D-31 junto a la corrección de D-26 en la actualización solicitada del PR #45. Los 31 identificadores tienen alcance local a US-16; no reemplazan las decisiones homónimas de otras historias. El responsable aprobó posteriormente de forma expresa las 31 decisiones y PT-01…PT-06 y autorizó implementar completamente US-16 mediante SDD → BDD → TDD en la misma rama del PR #45, con commits y push; sin merge ni cierre de #16. Aceptación final pendiente.

| ID | Decisión aprobada | AC de SDD-US-16 |
| --- | --- | --- |
| D-01 | Se pueden registrar horas estimadas para cualquier User Story del Product Backlog, aunque no esté asignada a un Sprint. | AC-01 |
| D-02 | Las horas estimadas pueden modificarse en cualquier momento, incluso cuando existen registros de esfuerzo real. Se conserva el valor actualizado como estimación vigente. | AC-01 |
| D-03 | Cada registro debe asociarse obligatoriamente a una User Story del proyecto. No existen registros generales sin historia. | AC-04 |
| D-04 | Se pueden registrar horas reales para historias asignadas o no a un Sprint. | AC-04 |
| D-05 | Las horas trabajadas deben ser mayores que cero y admitir hasta dos decimales. | AC-05 |
| D-06 | Las horas estimadas admiten hasta dos decimales. | AC-01, AC-02 |
| D-07 | Toda User Story nueva requiere horas estimadas mayores que cero desde su creación. | AC-02 |
| D-08 | Las historias existentes sin estimación se conservan como pendientes de regularización. Se pueden consultar, pero al editarlas se exige completar una estimación válida. | AC-03 |
| D-09 | La fecha de trabajo puede ser actual o pasada, nunca futura. | AC-07 |
| D-10 | La actividad es un texto libre obligatorio, no vacío ni compuesto únicamente por espacios. | AC-05 |
| D-11 | Los registros pueden editarse después de guardarse: miembro, fecha, actividad, horas e historia asociada, conservando su identificador. | AC-10 |
| D-12 | Los registros pueden eliminarse únicamente después de una confirmación del usuario. | AC-11 |
| D-13 | La eliminación es lógica (*soft delete*). Los registros eliminados permanecen almacenados para trazabilidad, pero no aparecen en las consultas normales ni participan de las sumatorias. No se requiere interfaz para consultar eliminados. | AC-11, AC-12 |
| D-14 | Se permite registrar, editar y eliminar lógicamente esfuerzo para historias completadas, sin modificar su estado. | AC-13 |
| D-15 | Se permite registrar, editar y eliminar lógicamente esfuerzo en proyectos finalizados, sin cambiar su estado ni sus fechas. | AC-08, AC-13 |
| D-16 | Un miembro puede tener múltiples registros independientes el mismo día para la misma historia. Sus horas se acumulan. | AC-09 |
| D-17 | En proyectos finalizados, la fecha de trabajo no puede ser posterior a la fecha real de finalización, aunque el registro se cargue después. | AC-08 |
| D-18 | Si existe fecha de inicio del proyecto, no se permite esfuerzo anterior a ella. Los cambios posteriores de esa fecha deben rechazarse cuando generen inconsistencias con los registros existentes; nunca modificar automáticamente sus fechas. | AC-07 |
| D-19 | Cada User Story muestra sus horas estimadas, el total de horas reales de sus registros activos y permite consultar el detalle. | AC-12 |
| D-20 | El registro y la administración del esfuerzo real se realizan desde el detalle de cada User Story del Product Backlog. | AC-12 |
| D-21 | Cuando cambian las horas estimadas de una historia, se actualiza también ese valor en todos los Sprints a los que está asignada, incluidos los cerrados. Los demás datos históricos se conservan. | AC-16 |
| D-22 | Cada registro puede asociarse opcionalmente a un Sprint al que pertenezca la historia. También puede permanecer sin Sprint. | AC-14 |
| D-23 | Se permite registrar, editar y eliminar lógicamente esfuerzo asociado a Sprints cerrados, sin reabrirlos ni alterar sus historias planificadas. | AC-15 |
| D-24 | Un miembro puede acumular como máximo 24 horas reales por fecha dentro de su proyecto, considerando todas sus historias y Sprints. Los registros eliminados lógicamente no se contabilizan. | AC-09 |
| D-25 | Las horas reales acumuladas pueden superar las estimadas, sin modificar automáticamente la estimación. | AC-09 |
| D-26 | Al cambiar la User Story de un registro, conservar el Sprint anterior si sigue siendo válido. Si es incompatible, exigir seleccionar otro Sprint válido o quitar explícitamente la asociación. Nunca modificar la asociación automáticamente. | AC-10, AC-14 |
| D-27 | Un registro de esfuerzo solo puede reasignarse a historias del mismo proyecto. Para corregir un proyecto equivocado, debe eliminarse lógicamente el registro y crearse otro en el proyecto correcto. | AC-06, AC-10 |
| D-28 | Se permite registrar esfuerzo en historias históricas sin estimación, mostrando «Pendiente de regularización» hasta completar ese dato. | AC-03 |
| D-29 | El formulario de esfuerzo debe completar por defecto la fecha actual, permitiendo modificarla según las restricciones establecidas. | AC-07, AC-12 |
| D-30 | Proyectos finalizados sin fecha de fin: bloquear altas y ediciones de esfuerzo si `end_date = NULL`. Permitir consultas y eliminación lógica. No inventar fechas ni modificar el estado del proyecto. | AC-19 |
| D-31 | Cambio de fecha de inicio: validar D-18 únicamente contra registros de esfuerzo activos. Los eliminados lógicamente no bloquean el cambio y conservan sus datos históricos. Garantizar atomicidad de la operación. | AC-07, AC-18 |

Integración aprobada con US-03: seleccionar miembros desde la lista del proyecto mediante UUID, nunca nombre libre; validar siempre su pertenencia. Cada esfuerzo conserva como mínimo miembro, fecha, actividad, horas y User Story, además de Sprint opcional y datos técnicos de identidad/eliminación lógica. Horas estimadas y reales son independientes: **no convertir Story Points en horas ni viceversa**. El día actual reutiliza la zona ya aprobada de US-02, `America/Argentina/Buenos_Aires`.

Trazabilidad: [enunciado](enunciado.md) §2.5 y datos para §2.7 → EP-05 → [US-16 / Issue #16](https://github.com/Aguscampos04/Carva/issues/16) → D-01…D-31 → [SDD-US-16](../sdd/SDD-US-16-registrar-horas-esfuerzo.md) → AC-01…AC-19 y su matriz → [BDD-US-16](../bdd/BDD-US-16-registrar-horas-esfuerzo.feature), escenarios BDD-US-16-01…46 → pruebas Go/Node, ciclos TDD y código implementados → [evidencia US-16](evidencia-ia-us-16.md); revisión y aceptación final pendientes. D-30 → AC-19 → BDD-43/44; D-31 → AC-07/18 → BDD-15/16/33/45; D-26 corregida → AC-10/14 → BDD-26/46. La evidencia distingue pruebas automatizadas Go/Node de la especificación feature, sin atribuirle ejecución por un runner Gherkin.

Impactos implementados y cubiertos por regresión (véase evidencia US-16): US-05 requiere estimación en altas/ediciones, con NULL histórico consultable; US-02 debe proteger cambios de inicio frente al esfuerzo; US-03 aporta identidades; US-08 y US-11 incorporan la excepción explícita D-21 únicamente para horas estimadas vigentes, preservando los demás datos históricos; US-10 conserva esfuerzo en su asociación elegida al trasladar historias y sincroniza estimación; US-15 conserva su actualización independiente de Story Points. Estas precisiones no reescriben otras decisiones funcionales.

R-01 y R-02 quedan **resueltos** por D-30 y D-31 y se retiran de los pendientes funcionales. La excepción de D-30 precisa D-15 para Finalizado sin fin; no permite completar esa fecha desde esfuerzo ni bloquea consultas/baja. D-31 precisa D-18: el cambio de inicio se coordina atómicamente contra activos, sin alterar eliminados aunque conserven fechas anteriores al nuevo inicio. D-26 corregida conserva Sprint válido o asociación nula al cambiar historia; solo una incompatibilidad exige resolución explícita, sin sustitución automática.

Las propuestas técnicas PT-01…PT-06 de la SDD (representación exacta, API/errores, esquema, mecanismo transaccional, migración e interfaz) quedan aprobadas formalmente para implementar; ni estas ni las reglas funcionales siguen pendientes de decisión. Se autoriza implementación y automatización con evidencia real, conservando fuera de alcance métricas/reportes. Sin merge, cierre de #16 ni aceptación final implícita.

### Decisiones pendientes

Todavía deben definirse y documentarse, entre otras:

- Contratos REST generales: versionado, formatos, códigos HTTP, errores y rutas fuera del alcance mínimo descrito en `sdd/SDD-US-01-soporte-proyecto.md` y `sdd/SDD-US-07-crear-sprint.md`.
- Modelo de dominio más allá de los campos mínimos de proyectos y Sprints especificados en esas SDD.
- Esquema SQLite y configuración de concurrencia/transacciones más allá de las tablas mínimas descritas en esas SDD.
- Estrategia para servir los archivos estáticos del frontend y configurar su acceso a la API.
- Configuración completa de los entornos de desarrollo y pruebas, más allá de la ruta base aprobada.
- Revisión final de la integración de selección de miembros en votos y esfuerzo. US-16 implementa selección por UUID y validación de proyecto según su SDD autorizada; véase evidencia US-16. El mecanismo funcional está aprobado (D-05), sin autenticación ni permisos (D-07); el alta/listado de US-03 está implementado según su SDD autorizada, pendiente de revisión final del equipo.
- Si se requiere autenticación en una etapa posterior.
- Estrategia de despliegue, si corresponde.
- Diseño visual de la interfaz.
- Revisión final de la implementación de US-02. SDD/BDD y PT-01 a PT-05 están aprobados para implementar; no sigue pendiente la condición funcional de finalización.
- Fórmulas de las métricas obligatorias, que deberán especificarse en sus SDD antes de implementarlas.

Estas decisiones deberán tomarse antes de implementarlas y quedar documentadas cuando corresponda.

## 13. Regla fundamental

El enunciado oficial del TPI es la fuente de verdad para los requisitos obligatorios.

Cuando exista una diferencia entre una propuesta del agente y el enunciado, deberá prevalecer el enunciado.

El agente no deberá inventar requisitos obligatorios que no estén respaldados por el enunciado o por una decisión explícita del equipo.
