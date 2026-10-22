# SDD-US-10 — Cerrar Sprint y trasladar historias incompletas

**Issue:** [#10](https://github.com/Aguscampos04/Carva/issues/10)
**Estado:** criterios de la Issue y decisión funcional de traslado aprobados; contrato de cierre acordado con la responsable el 2026-10-22.
**Trazabilidad:** enunciado §2.3 → EP-03 → US-10 / Issue #10 → decisiones funcionales aprobadas en `docs/contexto-proyecto.md` → esta SDD → [BDD-US-10](../bdd/BDD-US-10-cerrar-sprint.feature) → tests de dominio e integración.
**Dependencias:** US-07 (Sprint) y US-08 (asignación e historial), ambas integradas en `main`; US-09 provee el registro de finalización verificada.

## 1. Objetivo

Cerrar un Sprint y, en la misma operación, trasladar sus historias no completadas al siguiente Sprint disponible del proyecto. El cierre debe conservar intactos el plan y los resultados históricos del Sprint de origen.

## 2. Inputs

| Entrada | Tipo / origen | Regla |
| --- | --- | --- |
| `projectID` | UUID de la ruta | Identifica el proyecto propietario del Sprint. |
| `sprintID` | Entero positivo de la ruta | Identifica el Sprint que se desea cerrar. |
| Cuerpo | Ninguno | `POST /api/projects/{projectID}/sprints/{sprintID}/close` no recibe datos de negocio. |

La condición de historia incompleta usa el registro de US-09: `acceptance_criteria_verified = false`. Una verificación confirmada significa que la historia está completada.

## 3. Expected Outputs

- Cierre exitoso: `200 OK` con el Sprint y `closed: true`.
- Cierre con historias incompletas pero sin otro Sprint abierto posterior: `409 Conflict`; no cierra ni traslada parcialmente.
- Sprint o proyecto no existente/no relacionado: `404 Not Found`.
- Cierre repetido, asignación a un Sprint cerrado o finalización de una historia desde un Sprint cerrado: `409 Conflict`.
- Error inesperado de persistencia: `500 Internal Server Error`; se revierten conjuntamente el traslado y el cierre.
- El conjunto planificado y el resultado de cada historia del Sprint cerrado permanecen consultables sin cambios.

## 4. Business Rules

- **BR-01:** un Sprint se considera cerrado cuando su indicador persistido `closed` es `true`.
- **BR-02:** se consideran incompletas las historias planificadas cuyo registro de US-09 no confirma la verificación de todos sus Acceptance Criteria.
- **BR-03:** si hay historias incompletas, se elige como destino el Sprint abierto del mismo proyecto con menor ID mayor al del Sprint que se cierra.
- **BR-04:** si hay historias incompletas y no existe destino disponible, se rechaza el cierre sin cambios.
- **BR-05:** al cerrar con destino, se copian al destino las historias incompletas con los datos del snapshot planificado de origen; la marca de finalización del destino queda sin verificar. Una historia que ya esté asignada al destino no se duplica.
- **BR-06:** las historias completadas no se trasladan.
- **BR-07:** no se modifican las historias ni los resultados del Sprint de origen. Su plan y su estado de ejecución continúan disponibles como historial.
- **BR-08:** un Sprint sin historias incompletas puede cerrarse aunque no exista un Sprint posterior.
- **BR-09:** traslado y cierre constituyen una sola transacción: ambos se confirman o ninguno se aplica.
- **BR-10:** un Sprint cerrado no admite nuevos elementos planificados ni registro de finalización de historias. Las consultas históricas continúan permitidas.
- **BR-11:** un Sprint cerrado no puede cerrarse nuevamente.

## 5. Constraints

- API REST: `POST /api/projects/{projectID}/sprints/{sprintID}/close` devuelve el Sprint actualizado con su indicador `closed`.
- La selección de destino respeta el orden de ID existente; los IDs de Sprint son asignados por SQLite y la consulta se limita al proyecto.
- Se mantiene el snapshot de US-08; el traslado reutiliza sus campos planificados sin reconsultar ni mutar el Product Backlog.
- La lógica de reglas y orquestación permanece en Go; SQLite persiste el cierre y la transferencia; el frontend actual no ofrece operaciones de Sprint y esta historia sigue el límite REST de US-07/US-08.
- No se incluyen reapertura, borrado de Sprints ni edición de snapshots históricos.

## 6. Edge Cases

- Sprint sin historias: cerrar sin requerir destino.
- Sprint con todas las historias completadas: cerrar sin requerir destino ni trasladarlas.
- Sprint con combinación de historias incompletas y completadas: trasladar solo las incompletas.
- Historia incompleta ya presente en el Sprint destino: conservar una sola asignación en destino.
- El siguiente Sprint cronológico ya está cerrado: seleccionar el Sprint abierto posterior de menor ID; si no existe, rechazar el cierre.
- Sprint ya cerrado: rechazar la repetición y conservar los datos.
- Falla durante la transferencia o al marcar el cierre: rollback completo.

## 7. Error Conditions

- Sprint inexistente o perteneciente a otro proyecto: `404`.
- Historia incompleta sin Sprint posterior abierto: `409`.
- Cierre repetido o mutación de un Sprint cerrado: `409`.
- Falla de almacenamiento: `500`, sin traslado ni cierre parcial.

## 8. Acceptance Criteria

- **AC-01:** Si hay historias incompletas y no existe un siguiente Sprint disponible, el cierre se rechaza y no se modifica el Sprint.
- **AC-02:** Con un destino disponible, se trasladan automáticamente las historias incompletas; las completadas no se trasladan y no se duplican asignaciones existentes.
- **AC-03:** El Sprint cerrado conserva su plan original y el resultado de cada historia después del traslado.
- **AC-04:** Un Sprint sin historias incompletas puede cerrarse sin Sprint posterior.
- **AC-05:** El cierre persiste; las consultas históricas siguen disponibles y el Sprint cerrado rechaza nuevas asignaciones, finalizaciones y cierres repetidos.

## Trazabilidad

| Criterio | Escenario BDD | Verificación |
| --- | --- | --- |
| AC-01 | BDD-01 | Integración: bloqueo sin destino y sin cambio parcial. |
| AC-02 | BDD-02 | Integración: transferencia selectiva, destino abierto y sin duplicados. |
| AC-03 | BDD-03 | Integración: snapshots y resultados de origen sin cambios. |
| AC-04 | BDD-04 | Integración: Sprint vacío/completo se cierra sin destino. |
| AC-05 | BDD-05/06 | Integración: persistencia, consulta y rechazo de mutaciones repetidas. |
