# SDD-US-19 — Actualizar defecto y registrar su resolución

**Estado:** SDD y BDD aprobados para implementar por @Aguscampos04 el 2026-10-27 en este chat: «apruebo, acordate de hacer la US-19 nomas». Aprobación comprende P-01…P-04 y PT-01…PT-04; no implica validación final de la historia.
**Issue:** [#19](https://github.com/Aguscampos04/Carva/issues/19).
**Responsable:** @Aguscampos04. **Milestone:** Sprint 3 — Funcionalidad y calidad.
**Origen:** enunciado §2.6 → EP-06 → US-19 → esta SDD → AC → BDD → TDD registrado → Go/tests.

## 1. Objetivo

Actualizar el estado y los Sprints de detección/resolución de un defecto existente,
permitiendo asociarlo después de registrarlo sin Sprint y resolverlo en uno distinto
del de detección. Las consultas deben reflejar los valores persistidos.

### Dependencias verificadas antes de la preparación

- US-18: Issue #18 cerrada con DoD marcada; PR #47 integrado en `ed27692`.
- US-07, necesaria al vincular Sprints: Issue #7 cerrada con DoD marcada y código integrado.
- Base local y remota al iniciar: `main` en `ed27692`, checkout limpio.
- Dependencia transitiva US-05: ya integrada como soporte del Backlog de US-18.
- No se identificó una implementación existente de actualización: US-18 solo expone Create/ListByStory.
- Habilitada por dependencias; aprobación de especificación y autorización de implementación confirmadas.

### Fuentes y alcance aprobado

Se reutilizan los cuatro estados (`Abierto`, `En progreso`, `Resuelto`, `Cerrado`),
la identidad UUID del defecto, su historia/proyecto y los campos nullable
`detection_sprint_id`/`resolution_sprint_id` de US-18. Cada Sprint informado debe
pertenecer al proyecto. La Issue no prescribe una máquina de transiciones nueva.

Alcance previsto: actualización de esos tres campos y edición web dentro del
listado de defectos existente. No editar descripción/severidad, reasignar historia,
eliminar defectos, introducir fechas de resolución ni implementar métricas US-22.

## 2. Inputs

**PT-01, contrato técnico aprobado:**
`PATCH /api/projects/{projectID}/backlog/{itemID}/defects/{defectID}`.

Proyecto, historia y defecto identificados por la ruta. El defecto debe pertenecer
exactamente al proyecto/historia indicados. Body JSON parcial no vacío:

| Campo | Si se omite | Si se informa |
| --- | --- | --- |
| `status` | Conservar estado vigente | Uno de los cuatro valores exactos; `null`/vacío/tipo incorrecto se rechazan. |
| `detection_sprint_id` | Conservar asociación vigente | Entero positivo del mismo proyecto; `null` quita asociación conforme a P-02 aprobada. |
| `resolution_sprint_id` | Conservar asociación vigente | Entero positivo del mismo proyecto; `null` quita asociación conforme a P-02 aprobada. |

Ejemplo: `{"status":"Resuelto","resolution_sprint_id":2}` conserva la detección.
No permitir claves adicionales, duplicadas, cuerpos no objeto ni contenido JSON
posterior. No coercionar strings/fracciones/booleanos a un ID de Sprint.

## 3. Expected Outputs

- **PT-01:** `200` con el objeto Defect completo existente, mismo UUID, historia,
  proyecto, descripción y severidad; únicamente cambian campos solicitados.
- Reutilizar `GET /api/projects/{projectID}/backlog/{itemID}/defects`: debe devolver
  el defecto actualizado. No hace falta otro endpoint de consulta para esta historia.
- Al reabrir SQLite deben conservarse los cambios.
- Un error no deja actualizaciones parciales ni mensajes falsos de éxito.
- UI: botón Editar por defecto, estado y ambos Sprints precargados, Guardar/Cancelar;
  consulta/listado refleja valores confirmados por el servidor. Carga/error visibles,
  borrador conservado ante error, doble envío bloqueado y respuestas obsoletas descartadas.

## 4. Business Rules

### Reglas derivadas del alcance ya aprobado

- BR-01: conservar ID, proyecto e historia; no crear otro defecto al actualizar.
- BR-02: aceptar solo los cuatro estados aprobados.
- BR-03: detección y resolución son asociaciones distintas; sus Sprints pueden diferir.
- BR-04: cada Sprint informado existe en el mismo proyecto; no asociar Sprints ajenos.
- BR-05: las consultas recuperan los valores almacenados, sin deducir asociaciones
  desde el Sprint actual ni modificar historia, proyecto o Sprint.

### Decisiones funcionales aprobadas

| ID | Propuesta | Motivo / alternativa |
| --- | --- | --- |
| P-01 | Permitir cualquier cambio entre los cuatro estados, sin transiciones automáticas ni obligación de Sprint para Resuelto/Cerrado. Cambiar estado no asigna ni borra Sprints; informar resolución no cambia estado automáticamente. | Mantener la independencia de campos de US-18 y no inventar un flujo no definido. Alternativa: el equipo define restricciones explícitas antes de implementarlas. |
| P-02 | Permitir corregir ambos Sprints, incluyendo quitar una asociación con `null` y volver a Sin Sprint. | Corregir un dato registrado; la Issue exige asociación posterior pero no especifica desasociación. Alternativa: solo permitir asignar/reemplazar y rechazar `null` en PATCH. |
| P-03 | Aceptar Sprints del proyecto abiertos o cerrados, sin exigir que la historia esté asignada a ellos, ni imponer orden de IDs entre detección y resolución. | Reutiliza las validaciones de alta de US-18; el ID de creación no equivale a una fecha de detección/resolución. No altera planificación histórica. |
| P-04 | Permitir estas actualizaciones en proyectos finalizados y en historias completadas. | El alta de US-18 no los bloquea; evita una restricción nueva y permite registrar/corregir resolución. Alternativa: acordar bloqueos explícitos si el equipo los necesita. |

P-01…P-04 y PT-01…PT-04 fueron aprobadas explícitamente en la fecha indicada arriba. Las alternativas describen opciones descartadas por esa aprobación.

## 5. Constraints / diseño técnico aprobado

- PT-02: dominio con cambio parcial (presencia distinta de nulo), función que
  devuelve copia validada y errores de estado/Sprint inválido; identidad inmutable.
- PT-03: caso de uso Update y puerto de repositorio; SQLite lee/valida/escribe en
  una transacción `withImmediate`, siguiendo US-18. No agregar tablas ni modificar
  `000013_defects.sql`; ya almacena todos los datos y tiene FK compuestas.
- Actualizar contra el registro vigente bajo lock: cambios concurrentes de campos
  distintos se conservan; para el mismo campo prevalece el último commit válido.
  Sin versionado optimista nuevo, ni efectos laterales dentro de callbacks repetibles.
- Validación de recursos y escritura en la misma transacción; falla de validación,
  FK o almacenamiento revierte todos los cambios. HTTP no expone SQL.
- PT-04: extender `WithDefects` y `createDefects`, conservando alta/listado existentes,
  sus selectores/identidad de proyecto y estilo. Los IDs opcionales pueden ingresarse
  como hoy; no crear una gestión nueva de Sprints ni agregar frameworks.
- La UI consume PATCH/GET; Go aplica las reglas de negocio. Cancelar no envía PATCH.
- No cambiar asignaciones de GitHub, Milestones, Story Points ni prioridades.

## 6. Edge Cases

Asociación tardía, mismo Sprint en ambos campos, Sprints distintos, edición de un
solo campo sin borrar el resto, repetición del mismo cambio sin duplicar, recursos
ajenos/inexistentes, IDs cero/negativos/fraccionarios, JSON vacío/desconocido/duplicado,
error de persistencia después de validar, solicitudes concurrentes y cambios de
proyecto mientras la UI espera respuesta. Los casos P-01…P-04 están aprobados
conforme al encabezado.

## 7. Error Conditions (PT-01)

| HTTP | Situación |
| --- | --- |
| 400 | JSON inválido, duplicado, desconocido, sin campos editables, estado inválido/nulo o Sprint ID inválido. |
| 404 | Proyecto/historia/defecto inexistente o fuera del ámbito; Sprint inexistente/ajeno. |
| 503 | Contención que excede el límite existente (`projects.ErrBusy`). |
| 500 | Falla inesperada de almacenamiento, con error genérico y rollback. |

Mantener formato `{"error":"..."}`. Mensajes específicos se alinearán con
`writeDefectError`; no filtrar existencia de defectos ajenos ni detalles internos.

## 8. Acceptance Criteria / trazabilidad

| ID | Criterio oficial | BDD previsto |
| --- | --- | --- |
| AC-01 | Un defecto inicialmente sin Sprint puede asociarse posteriormente a uno de detección. | US19-01 |
| AC-02 | Puede registrarse un Sprint de resolución distinto del de detección. | US19-02 |
| AC-03 | Los cambios de estado y asociación se reflejan al consultar. | US19-01/02/03/04 |

Tests implementados: dominio para cambios parciales
y valores; aplicación para propagación de errores; SQLite/HTTP temporal para
persistencia, ámbito, rollback y concurrencia; Node para edición/cancelación,
errores/doble envío/respuestas obsoletas; regresión US-18, Sprints e historial.

Los escenarios etiquetados `@propuesta` conservan el identificador de su origen;
P-01…P-04/PT-01…PT-04 están aprobados. La evidencia real de RED, GREEN y refactor
se registra en `docs/evidencia-ia-us-19.md`.

## 9. Estado de esta etapa

Se autoriza la implementación local de US-19 y sus pruebas tras la aprobación indicada.
La autorización específica de push/PR de US-11 no se extiende a US-19: publicación
pendiente del responsable. DoD y validación final del equipo quedan pendientes.
