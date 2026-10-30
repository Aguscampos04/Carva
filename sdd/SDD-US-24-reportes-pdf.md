# SDD-US-24 — Reportes de proyecto, Sprint y PDF final

**Estado:** P-01…P-04 y PT-01…PT-04 aprobadas por el Agile Enabler el 2026-10-30 en este chat («sisi»); implementación autorizada. DoD final pendiente.
**Fecha:** 2026-10-30. **Issue:** #24, EP-09, Sprint 4.
**Responsable GitHub:** @Loana-Vargas; sin cambios de asignación.
**Solicitud:** Agile Enabler pidió desarrollar US-24 en este chat.
**Fuentes:** enunciado §2.9 y plan de entrega; contexto decisiones 16/17;
SDD de Sprints, esfuerzo y métricas; Issue #24 y código integrado.

## 1. Objetivo

Consultar reportes de proyecto completo y de Sprint individual y descargar un
PDF final del proyecto con información general, estado, Sprints, historias,
estimaciones, esfuerzo, métricas y defectos. Los valores provienen de los casos
de uso existentes; no crear métricas alternativas ni modificar datos.

Dependencias directas #1/#5/#7/#8/#9/#16/#18/#19/#20/#21/#22 verificadas cerradas
con ocho puntos DoD marcados por Issue e implementaciones integradas en main.
Base `83eb916` incluye además UX de #53/PR #54. No hay reporte implementado,
SDD de US-24 ni PR equivalente detectados al iniciar.

## 2. Inputs

- Proyecto identificado por UUID.
- Sprint opcional, entero positivo perteneciente a ese proyecto.
- Datos persistidos de proyecto, miembros, Backlog, Sprints/snapshots, esfuerzo
  activo y defectos; métricas calculadas por US-20/21/22.
- Ningún parámetro permite sustituir los datos o introducir resultados inventados.

## 3. Expected Outputs — propuesta PT-01

| Operación GET | Salida |
| --- | --- |
| `/api/projects/{projectID}/reports/project` | JSON con proyecto, miembros, historias actuales, Sprints con planificación/resultados, esfuerzo activo, métricas y defectos. |
| `/api/projects/{projectID}/reports/sprints/{sprintID}` | JSON del Sprint, proyecto identificable, snapshots, resultados, esfuerzo asociado, métricas y defectos pertinentes. |
| `/api/projects/{projectID}/reports/project.pdf` | PDF descargable con el contenido del reporte de proyecto; satisface el PDF final. |
| `/api/projects/{projectID}/reports/sprints/{sprintID}.pdf` | PDF del mismo reporte de Sprint, reutilizando el generador. |

JSON/PDF usan el mismo modelo de reporte en Go. JSON incluye además `sections`,
proyección textual utilizada por la UI y el generador. PDF: `application/pdf`,
`Content-Disposition: attachment` con nombre seguro por alcance y número de Sprint.
Documento paginado, encabezados, acentos legibles, texto largo sin recortes,
listas vacías explícitas y unidades claras; sin depender de imprimir el navegador.

## 4. Business Rules — propuestas P-01…P-04

- P-01: proyecto incluye cada historia actual del Backlog una vez y el desglose
  de cada Sprint con sus propios snapshots. Una historia trasladada puede
  aparecer en varios Sprints, sin duplicar el total de esfuerzo del proyecto.
- P-02: reporte del Sprint lee snapshots, Goal y resultados registrados; no
  reemplaza esos datos por el Backlog actual. Respeta la excepción aprobada
  D-21 de US-16 para horas estimadas vigentes; las etiqueta como tales.
- P-03: esfuerzo del proyecto incluye registros activos con/sin Sprint;
  esfuerzo del Sprint incluye solo registros activos asociados explícitamente
  a ese Sprint. Métricas mantienen las reglas de sus SDD; no sumarlas desde
  listas en JavaScript. Defectos del Sprint: unión sin duplicados de los que
  indiquen ese Sprint en detección o resolución; ambas asociaciones se muestran
  separadas. Esta lista no sustituye las reglas de conteo de US-22.
- P-04: permitir consultar/exportar proyectos Planificados, En progreso o
  Finalizados y Sprints abiertos/cerrados. «PDF final» describe el entregable;
  no certifica aceptación ni finaliza el proyecto. No exigir un cierre nuevo.
- Los datos numéricos ausentes se muestran como **No disponible**, nunca cero
  ficticio. Un cero válido conserva su significado. Sin asociación: Sin Sprint.
- Conservar centésimas de horas y distinguir puntos provisionales/definitivos.
- Consultar/exportar no escribe, completa historias, cierra ni traslada Sprints.

## 5. Constraints — propuestas PT-02…PT-04

- PT-02: caso de uso `application/reports` que reúne datos existentes y devuelve
  un modelo compartido. Persistencia de solo lectura y reutilización de cálculos
  Go existentes. Sin nuevas tablas ni modificaciones de migraciones previas.
- PT-03 aprobada: `github.com/signintech/gopdf v0.38.3` para PDF en Go,
  junto a `golang.org/x/image v0.47.0/font/gofont/goregular` para la fuente
  TTF redistribuible embebida. Licencias MIT/Go en `docs/licencias`.
  Evitar ejecutables externos, navegador headless o servicios de terceros.
  Se conservan Go mínimo 1.26.0 y el driver SQLite aprobado.
  Referencia primaria: https://github.com/signintech/gopdf.
- La fuente admite acentos españoles; caracteres sin glifo se representan como
  `[U+XXXX]` para no perderlos silenciosamente. JSON y UI conservan el texto original.
- PT-04: integrar sección Reportes en la navegación del proyecto de UX #53;
  elegir proyecto completo o Sprint, consultar y descargar PDF. Carga/error/vacío,
  reintento y descarte de respuestas de un contexto anterior. No rediseñar toda
  la aplicación ni implementar reglas de negocio en JavaScript.
- No prometer snapshot atómico entre todas las fuentes sin soporte transaccional
  verificado. El reporte muestra los datos consultados al generarlo; no es un
  cierre contable ni un historial nuevo. Esta limitación se documentará.
- No modificar Issues, milestones, responsables, prioridades ni estados.

## 6. Edge Cases

Proyecto sin historias/Sprints/defectos; Sprint vacío; datos insuficientes;
historias trasladadas; proyecto con esfuerzo sin Sprint; esfuerzo eliminado;
defecto detectado y resuelto en distintos Sprints; textos largos, saltos de línea,
acentos y caracteres especiales; nombres idénticos con UUID distinto; muchas
páginas; consulta de proyecto/Sprint ajeno; cambio de ruta durante descarga.

## 7. Error Conditions

400: ID de Sprint no positivo/no entero. 404: proyecto o Sprint inexistente/ajeno.
500: lectura o generación fallida, mensaje genérico sin SQL ni rutas internas.
No responder éxito ni un PDF incompleto después de un error; preparar el
documento antes de enviar los headers de éxito.

## 8. Acceptance Criteria / trazabilidad

| AC oficial | Escenarios / tests previstos |
| --- | --- |
| AC-01 reportes de proyecto y Sprint con datos de su alcance | BDD-01/02/04; tests de aplicación, SQLite/HTTP y web. |
| AC-02 PDF final con todas las secciones requeridas | BDD-03/06; validación de archivo, extracción de texto y revisión visual multipágina. |
| AC-03 concordancia con fuentes y No disponible | BDD-02/04/05; comparación con casos de uso existentes y null frente a cero. |

SDD → BDD → RED real → implementación mínima → GREEN → refactor con regresión.
Evidencia de pruebas y TDD ejecutados en `docs/evidencia-ia-us-24.md`.
Publicación/PR requieren autorización específica; merge siempre manual por el responsable.

## 9. Estado de implementación

P-01…P-04 y PT-01…PT-04 aprobadas explícitamente antes de implementar, mediante
«sisi» en este chat el 2026-10-30. Implementación local de los cuatro GET, modelos,
generación PDF e integración web terminada; validación técnica documentada.
Revisión humana, confirmación final de DoD y publicación pendientes. No se
modificaron responsable, milestone, estado ni otros campos de la Issue #24.
