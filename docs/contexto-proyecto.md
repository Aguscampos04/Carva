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

La estructura de carpetas incluida en ese documento es una propuesta para revisión; no constituye una decisión aprobada sobre la organización definitiva del repositorio. Los detalles de contratos, modelo de dominio, dependencias y configuración deben precisarse antes de implementarlos.

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
- Se utilizará Scrum.
- Se utilizarán SDD, BDD y TDD.
- Se utilizará Git.
- Se utilizará IA como soporte.
- El proyecto se desarrollará incrementalmente mediante Sprints.

### Product Backlog inicial aprobado

El responsable del proyecto confirma la aprobación de las **9 Épicas (EP-01 a EP-09)** y las **24 User Stories (US-01 a US-24)** como versión inicial del Product Backlog. Esta aprobación establece la base inicial del alcance planificado; el Product Backlog podrá evolucionar mediante decisiones posteriores aprobadas por el equipo, manteniendo la trazabilidad con los requisitos del enunciado.

### Decisiones funcionales aprobadas

Estas decisiones definen el comportamiento funcional del producto. Las fórmulas concretas de las métricas continúan pendientes y deberán documentarse en sus SDD.

1. **Estado del proyecto:** al crearlo queda **Planificado**; al crear su primer Sprint pasa a **En progreso**; pasa a **Finalizado** cuando se cumple la condición de finalización del proyecto. La condición exacta que determina que el proyecto está terminado sigue pendiente de definición.
2. **Estado de una User Story:** los estados son **Pendiente**, **En progreso** y **Completada**.
3. **Finalización de una User Story:** se considera completada cuando se verifican sus Acceptance Criteria.
4. **Cierre de un Sprint:** las User Stories incompletas se trasladan automáticamente al siguiente Sprint. Si existen historias incompletas, no se puede cerrar el Sprint hasta que haya un siguiente Sprint disponible. Las historias completadas no se trasladan. El Sprint cerrado conserva su historial y el conjunto originalmente planificado.
5. **Story Points iniciales:** una User Story puede crearse con Story Points provisionales.
6. **Actualización de Story Points:** el Story Point actual solo se actualiza cuando se registra una estimación acordada mediante Planning Poker.
7. **Planning Poker:** la votación es individual y oculta; las estimaciones se revelan simultáneamente; se discuten las diferencias; pueden realizarse nuevas rondas; y se registra el valor acordado.
8. **Diferencia de estimaciones:** existe cuando los valores votados no son todos iguales.
9. **Estimación acordada:** hay consenso cuando todos los participantes de la ronda votan el mismo valor. Ese valor se registra como Story Point definitivo. Sin consenso, la estimación no se finaliza ni se actualiza el Story Point actual.
10. **Horas:** las horas estimadas se registran para la User Story durante la planificación; las horas reales provienen del registro de esfuerzo.
11. **Fórmulas de métricas:** las fórmulas de las métricas obligatorias se definirán explícitamente en sus SDD antes de implementarlas.
12. **Alcance temporal de las métricas:** se calculan por Sprint y pueden agregarse a nivel proyecto. El conjunto planificado de un Sprint debe preservarse aunque el Backlog cambie posteriormente.
13. **Datos faltantes en métricas:** cuando falten datos para una métrica, se muestra **«No disponible»**, nunca cero en reemplazo de esos datos.
14. **Defectos:** sus estados son **Abierto**, **En progreso**, **Resuelto** y **Cerrado**. Sus severidades son **Baja**, **Media**, **Alta** y **Crítica**.
15. **Defectos y Sprints:** los defectos pueden existir sin Sprint, asociarse posteriormente a uno y resolverse en un Sprint diferente al de detección.
16. **Reportes:** habrá un reporte de proyecto completo y un reporte de Sprint individual. Incluirán User Stories, Story Points, horas estimadas y reales, métricas y defectos.
17. **Reporte final:** debe poder generarse en PDF e incluir información general del proyecto, Sprints, User Stories, Story Points, horas, métricas, defectos y estado general.
18. **MVP e identificación:** se identificará a los miembros para Planning Poker y el registro de esfuerzo. No se implementará un sistema complejo de permisos en el MVP.

### Decisiones pendientes

Todavía deben definirse y documentarse, entre otras:

- Organización definitiva de carpetas y paquetes; la propuesta de `docs/arquitectura.md` requiere revisión.
- Contratos concretos de la API REST: rutas, formatos, versionado, códigos HTTP y formato de errores.
- Modelo de dominio.
- Driver y versión de SQLite, estrategia de migraciones y configuración de concurrencia.
- Estrategia para servir los archivos estáticos del frontend y configurar su acceso a la API.
- Herramientas concretas para pruebas de integración y configuración de los entornos de prueba.
- Mecanismo concreto para identificar miembros registrados por proyecto. Para la etapa inicial está aprobada la ausencia de autenticación compleja.
- Si se requiere autenticación en una etapa posterior.
- Estrategia de despliegue, si corresponde.
- Diseño visual de la interfaz.
- Condición exacta que determina cuándo el proyecto está terminado y pasa a **Finalizado**.
- Fórmulas de las métricas obligatorias, que deberán especificarse en sus SDD antes de implementarlas.
- Tratamiento de los defectos sin Sprint en los SDD correspondientes.

Estas decisiones deberán tomarse antes de implementarlas y quedar documentadas cuando corresponda.

## 13. Regla fundamental

El enunciado oficial del TPI es la fuente de verdad para los requisitos obligatorios.

Cuando exista una diferencia entre una propuesta del agente y el enunciado, deberá prevalecer el enunciado.

El agente no deberá inventar requisitos obligatorios que no estén respaldados por el enunciado o por una decisión explícita del equipo.
