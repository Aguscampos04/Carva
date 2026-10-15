# SDD-US-06 — Consultar Product Backlog

**Estado:** Especificación aprobada por DelfiCarribero el 2026-10-15.  
**Trazabilidad:** enunciado §2.2 → EP-02 → US-06 → Issue #6 → este SDD → BDD-US-06 → tests y código.

## 1. Objetivo

Permitir consultar los elementos registrados en el Backlog de un proyecto, mostrando los datos vigentes que US-05 mantiene sin duplicar ni alterar registros.

## 2. Dependencia y datos consultados

Depende de US-05 y del proyecto creado mediante US-01. Cada elemento contiene: identificador, proyecto, título, descripción, prioridad, estado, Story Points provisionales y definitivos separados, y criterios de aceptación.

La respuesta expone por separado los campos `provisional_story_points` y `story_points` definitivo. En la interfaz, se muestra el valor definitivo cuando existe; de lo contrario, se muestra el provisional identificado como tal. Así no se confunde una estimación provisional con un consenso de US-15.

## 3. Salidas esperadas

- Para un proyecto existente, se devuelve una lista JSON, posiblemente vacía, con todos sus elementos y los campos exigidos por §2.2.
- Los valores consultados coinciden con los vigentes en SQLite; la lectura no modifica registros.
- No se incluyen elementos de otros proyectos.

## 4. Contrato HTTP propuesto para revisión

- `GET /api/projects/{projectID}/backlog`: devuelve `200` y un arreglo JSON de elementos.
- Si el proyecto existe y no tiene elementos, devuelve `200` con `[]`.
- Proyecto inexistente: `404`; fallo de almacenamiento: `500` con mensaje seguro.
- Cada elemento contiene `id`, `project_id`, `title`, `description`, `priority`, `status`, `provisional_story_points`, `story_points` y `acceptance_criteria`.
- La lectura no ofrece parámetros de filtro ni orden en esta historia. La lista no garantiza un orden específico.

El contrato fue aprobado por DelfiCarribero el 2026-10-15. La consulta de un elemento individual no está requerida por la Issue y queda fuera de alcance.

## 5. Criterios de aceptación

- **AC-01:** Se pueden consultar los elementos del Backlog con todos los campos exigidos por §2.2.
- **AC-02:** Prioridad, estado y valor de Story Points presentado corresponden al dato vigente del elemento; se distingue un valor provisional de uno definitivo.

## 6. Casos límite y errores

- Proyecto existente sin elementos: lista vacía, sin error.
- Proyecto inexistente: `404`, sin crear registros.
- Error de SQLite: `500`; no se devuelve información parcial como si fuera una lista completa.
- Elementos de otro proyecto no aparecen en la respuesta.

## 7. Restricciones y fuera de alcance

- Es una operación de solo lectura y no altera prioridad, estado, estimaciones ni criterios.
- La consulta utiliza Go, API REST y SQLite, manteniendo las capas aprobadas.
- No incluye búsqueda, filtros, paginación, orden configurable, edición (US-05), asignación a Sprint (US-08) ni estimación (US-15).

## 8. Decisiones

- La Issue #6 depende de US-05 y sus criterios de aceptación son la fuente del comportamiento.
- Contrato HTTP aprobado por DelfiCarribero el 2026-10-15.
