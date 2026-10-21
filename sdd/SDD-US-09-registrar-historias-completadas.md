# SDD-US-09 — Registrar historias completadas

**Issue:** #9  
**Trazabilidad:** enunciado §2.3 → EP-03 → US-09 / Issue #9 → este SDD → BDD-US-09 → tests y código.  
**Dependencias:** US-05 (Product Backlog) y US-08 (Sprint Backlog).

## 1. Objetivo

Registrar una historia del Sprint como completada después de que un integrante confirme explícitamente que verificó todos sus Acceptance Criteria. El estado actual debe reflejarse en el Backlog y en la consulta del Sprint sin modificar los datos originales de planificación guardados por US-08.

## 2. Inputs

- `projectID`: UUID del proyecto de la ruta.
- `sprintID`: entero positivo del Sprint de la ruta.
- `backlogItemID`: UUID de la historia asignada al Sprint.
- `acceptance_criteria_verified`: booleano obligatorio; `true` confirma la verificación conjunta de todos los Acceptance Criteria.

## 3. Expected Outputs

- Éxito: HTTP 200 con la historia del Sprint mostrando `status: "Completada"`, su `planned_status` original y `acceptance_criteria_verified: true`.
- La consulta del Product Backlog devuelve la misma historia con `status: "Completada"`.
- El resto de los campos del snapshot de planificación no cambia.
- La actualización de ambos estados es atómica.

## 4. Business Rules

- **BR-01:** solo una confirmación explícita `true` habilita el registro de finalización.
- **BR-02:** una historia completada aparece como `Completada` en el Backlog y en el estado actual del Sprint.
- **BR-03:** la historia debe pertenecer al proyecto y estar asignada al Sprint indicado.
- **BR-04:** los datos de planificación originales de US-08 permanecen inmutables; el estado actual de ejecución se representa por separado.
- **BR-05:** la verificación y los cambios de estado del Backlog y Sprint se registran en una única transacción.

## 5. Constraints

- `POST /api/projects/{projectID}/sprints/{sprintID}/backlog-items/{backlogItemID}/completion`
- Body: `{ "acceptance_criteria_verified": true }`.
- La respuesta de planificación agrega `planned_status` y `acceptance_criteria_verified`; `status` conserva el significado de estado actual que ya consumen los clientes.
- La verificación es una confirmación explícita del integrante; la aplicación no interpreta ni valida semánticamente el cumplimiento de cada criterio.

## 6. Edge Cases

- Body ausente, inválido, campo ausente, valor `false` o campos desconocidos: no completar.
- Proyecto, Sprint, historia o asignación inexistente: no completar.
- Historia ya completada: el resultado sigue siendo `Completada` (operación idempotente).
- Error al actualizar cualquiera de los dos registros: rollback de ambos.

## 7. Error Conditions

- Request JSON inválido o confirmación ausente/falsa: HTTP 400 si el formato/campo no es válido; HTTP 409 si el formato es válido y la verificación no fue confirmada.
- Recurso o asignación inexistente/no perteneciente al proyecto: HTTP 404.
- Error inesperado de persistencia: HTTP 500 sin detalles internos.

## 8. Acceptance Criteria

- **AC-01:** Una historia cuyos Acceptance Criteria no fueron verificados no puede pasar a `Completada`.
- **AC-02:** Tras verificar esos criterios y registrar su finalización, la historia figura como `Completada` en el Sprint y en su Backlog.

## 9. BDD / Test Traceability

- AC-01 → BDD-US-09-01/02 → tests de rechazo sin confirmación.
- AC-02 → BDD-US-09-03 → test de integración de estado en ambos recursos y snapshot conservado.
- BR-05 → BDD-US-09-04 → test de atomicidad ante error de persistencia.
