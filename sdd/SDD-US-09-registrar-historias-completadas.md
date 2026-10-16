# SDD-US-09 — Registrar historias completadas

**Estado:** Implementación de US-09; la mecánica de verificación fue confirmada por Loana Vargas el 2026-10-16. La integración visible con el Sprint depende de US-08.
**Trazabilidad:** enunciado §2.3 → EP-03 → US-09 → Issue #9 → este SDD → BDD-US-09 → tests y código.

## 1. Objetivo

Registrar una User Story como `Completada` solamente después de que una persona confirme explícitamente que verificó todos sus Acceptance Criteria. El estado consultado en el Backlog debe reflejar esa finalización.

## 2. Inputs

| Campo | Ubicación | Obligatorio | Validación |
| --- | --- | --- | --- |
| `projectID` | Ruta | Sí | Debe identificar un proyecto existente. |
| `itemID` | Ruta | Sí | Debe identificar una historia del proyecto. |
| `acceptance_criteria_verified` | JSON | Sí | Debe ser `true`; `false` o su valor omitido no autoriza la finalización. |

La confirmación única representa que se verificó el conjunto completo de Acceptance Criteria de la historia. No se registran resultados individuales por criterio.

## 3. Expected Outputs

- Finalización válida: HTTP `200` con la historia persistida y estado `Completada`.
- Consulta posterior del Backlog: la historia conserva el estado `Completada`.
- Confirmación falsa u omitida: HTTP `400`; el estado no cambia.
- Proyecto o historia inexistente: HTTP `404`.
- Fallo de persistencia: HTTP `500` sin detalles internos.

## 4. Business Rules

- **BR-01:** una historia solo puede registrarse como `Completada` si se confirma explícitamente que se verificaron todos sus Acceptance Criteria.
- **BR-02:** la confirmación no sustituye ni modifica los Acceptance Criteria almacenados.
- **BR-03:** la finalización cambia el estado de la historia persistida; no cambia sus otros datos.
- **BR-04:** la edición genérica de US-05 continúa sin poder establecer el estado `Completada`.
- **BR-05:** el estado del Backlog es el dato canónico de la historia. La consulta de historias planificadas en un Sprint deberá reflejar el mismo estado al integrarse US-08.

## 5. Constraints

- Endpoint propuesto e implementado: `POST /api/projects/{projectID}/backlog/{itemID}/complete`.
- Cuerpo: `{"acceptance_criteria_verified": true}`.
- La persistencia reutiliza el campo `status` ya definido en la tabla `backlog_items`; esta historia no agrega columnas ni cambia el esquema.
- No se agregan criterios individuales de verificación, autor, fecha de finalización ni reglas de transición adicionales.
- US-08 (asignación y consulta del conjunto planificado del Sprint) sigue abierta. La coherencia del estado en el Sprint se completa mediante la integración de US-08 con el estado canónico del Backlog; esta rama no modifica dicha funcionalidad.

## 6. Edge Cases

- Confirmación omitida se trata como no verificada y no cambia el estado.
- Una confirmación falsa no cambia el estado, incluso si la historia ya estaba `En progreso`.
- El identificador de historia perteneciente a otro proyecto no permite finalizarla.
- La edición genérica que incluye `status: "Completada"` no cambia el estado.
- La historia finalizada sigue apareciendo en el listado del Backlog con su estado persistido.

## 7. Error Conditions

- Cuerpo JSON inválido o confirmación que no es `true`: HTTP `400`.
- Proyecto o historia inexistente/en otro proyecto: HTTP `404`.
- Error inesperado de almacenamiento: HTTP `500`; no se devuelve éxito.

## 8. Acceptance Criteria

- **AC-01:** Una historia cuyos Acceptance Criteria no se verificaron no puede pasar a `Completada`.
- **AC-02:** Tras verificar los criterios y registrar su finalización, la historia figura como `Completada` en el Backlog y su estado puede consultarse desde el Sprint.

**Límite de validación de AC-02 en esta entrega:** el estado persistido y consultable en el Backlog se verifica en US-09. La vista de Sprint requiere que US-08, actualmente abierta, lea el estado canónico de la misma historia; no se implementa ni se afirma aquí.
