# SDD-US-22 — Consultar métricas de defectos

- **Issue:** [#22 — US-22](https://github.com/Aguscampos04/Carva/issues/22)
- **Épica:** EP-07 — Métricas
- **Requisito:** `docs/enunciado.md` §2.7
- **Dependencias:** US-18 y US-19 (cerradas e integradas).
- **Trazabilidad:** §2.7 → EP-07 → US-22 → #22 → esta SDD → BDD → tests → Go/API/UI.
- **Decisiones funcionales:** aprobadas por DelfiCarribero el 2026-10-27 en la conversación de implementación de US-22.

## 1. Objetivo

Consultar la cantidad de defectos detectados y resueltos por Sprint y como agregado del proyecto. La detección y la resolución se atribuyen de forma independiente, pues pueden ocurrir en Sprints diferentes.

## 2. Inputs

- `projectID`: UUID del proyecto indicado en la ruta.
- `sprintID`: entero positivo del proyecto para la consulta individual por Sprint.
- Las consultas no reciben parámetros de conteo ni modifican datos.

## 3. Expected Outputs

- `GET /api/projects/{projectID}/metrics/defects` devuelve HTTP 200 con `metrics` agregadas y `sprints`, cada uno con su ID, objetivo, estado de cierre y métricas.
- `GET /api/projects/{projectID}/sprints/{sprintID}/metrics/defects` devuelve HTTP 200 con las dos métricas del Sprint.
- `detected_defects` y `resolved_defects` son enteros no negativos cuando están disponibles y `null` cuando falta la atribución necesaria. La interfaz presenta `null` como **«No disponible»**.
- Los defectos sin Sprint siguen formando parte de los totales del proyecto. No se inventa una atribución a un Sprint.
- Las lecturas no modifican defectos, Sprints ni historias.

## 4. Business Rules

- **BR-01:** El total de defectos detectados del proyecto es la cantidad de registros de defectos del proyecto. La existencia de un registro significa que el defecto fue detectado, aunque no tenga Sprint de detección.
- **BR-02:** El total de defectos resueltos del proyecto es la cantidad de defectos cuyo estado actual es `Resuelto` o `Cerrado`.
- **BR-03:** Para un Sprint, defectos detectados son los registros cuyo `detection_sprint_id` coincide con el Sprint consultado. No se deduce el Sprint desde la historia ni desde la creación del defecto.
- **BR-04:** Para un Sprint, defectos resueltos son los registros cuyo estado actual es `Resuelto` o `Cerrado` y cuyo `resolution_sprint_id` coincide con el Sprint consultado. El Sprint de resolución es independiente del de detección.
- **BR-05:** Si al menos un defecto del proyecto no tiene `detection_sprint_id`, la métrica de defectos detectados por Sprint es no disponible: el total no puede atribuirse con certeza entre Sprints.
- **BR-06:** Si al menos un defecto en estado `Resuelto` o `Cerrado` no tiene `resolution_sprint_id`, la métrica de defectos resueltos por Sprint es no disponible.
- **BR-07:** Cuando los datos necesarios para una métrica por Sprint están completos y no hay coincidencias, su valor es cero. Cero no representa datos faltantes.
- **BR-08:** Las métricas por Sprint se calculan con los defectos del proyecto solicitado. Los agregados del proyecto no se obtienen sumando Sprints; así incluyen defectos sin Sprint y no dependen de la cantidad de Sprints.

## 5. Constraints

- Las métricas son de solo lectura y no agregan tablas ni migraciones.
- Las rutas y el envoltorio de respuesta siguen el patrón de US-20 y US-21.
- Los datos se limitan al proyecto indicado.
- Los contadores se serializan como enteros JSON; la falta de datos necesarios se serializa como `null`.
- El frontend presenta el agregado de proyecto y el desglose del Sprint seleccionado en la sección existente de métricas.

## 6. Edge Cases

- Proyecto sin defectos: ambos totales del proyecto son cero; si tiene Sprints, los conteos por Sprint son cero.
- Defecto sin Sprint de detección: cuenta como detectado en el proyecto y vuelve no disponibles los conteos de detección por Sprint.
- Defecto resuelto/cerrado sin Sprint de resolución: cuenta en el total resuelto del proyecto y vuelve no disponibles los conteos de resolución por Sprint.
- Defecto detectado en un Sprint y resuelto en otro: suma una detección al primer Sprint y una resolución al segundo.
- Proyecto sin Sprints: devuelve los totales y una lista de Sprints vacía.
- Sprint sin coincidencias: devuelve cero si la atribución de esa métrica está completa para el proyecto.

## 7. Error Conditions

- Proyecto inexistente: HTTP 404.
- Sprint inexistente o ajeno al proyecto: HTTP 404.
- `sprintID` vacío, no numérico, cero o negativo: HTTP 400.
- Error inesperado de lectura: HTTP 500 con mensaje genérico, sin detalles de SQLite.

## 8. Acceptance Criteria

- **AC-01:** Se muestran las métricas de detección y resolución por Sprint y agregadas al proyecto según BR-01…BR-08.
- **AC-02:** La detección y la resolución se atribuyen mediante sus respectivos campos, aun si los Sprints son diferentes.
- **AC-03:** Si falta la atribución requerida para alguna métrica por Sprint, se muestra «No disponible» y no un cero ficticio.
- **AC-04:** Con datos completos, la ausencia de coincidencias se muestra como cero.
- **AC-05:** Las consultas son de solo lectura y respetan el límite de proyecto.

## 9. BDD / Test Traceability

- AC-01/02 → BDD-US-22-01 → defecto detectado y resuelto en Sprints distintos.
- AC-01/05 → BDD-US-22-02 → agregado de proyecto, inclusión de defectos sin Sprint y consulta sin escrituras.
- AC-03 → BDD-US-22-03 → detección o resolución sin Sprint necesario.
- AC-04 → BDD-US-22-04 → cero válido con datos completos.
- AC-05 → BDD-US-22-05 → recursos inexistentes e ID de Sprint inválido.
