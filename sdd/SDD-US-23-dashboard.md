# SDD-US-23 — Consultar dashboard del proyecto

**Estado:** P-01/P-02 y PT-01…PT-03 aprobadas por @Aguscampos04 el 2026-10-27 en este chat («apruebo»). Implementación y validación final pendientes.
**Issue:** [#23](https://github.com/Aguscampos04/Carva/issues/23).
**Responsable:** @Aguscampos04. **Milestone:** Sprint 4 — Métricas y entrega final.
**Trazabilidad:** enunciado §2.8 → EP-08 → US-23 → AC → BDD → tests previstos.

## 1. Objetivo

Consultar el estado oficial del proyecto y gráficos de las métricas de US-20,
US-21 y US-22 sin crear fórmulas alternativas ni modificar los datos.

### Control de dependencias — 2026-10-27

Main actualizado por fast-forward a `3c5a6b1`, sin cambios locales previos.
Issues #4, #20, #21 y #22 cerradas con su DoD marcada; PR #38, #49, #50 y #51
integrados. SDD/BDD, pruebas Go/Node y evidencias disponibles en el repositorio
o descripción de los PR. El PR #51 registra CI verify exitoso.
Esto habilita por dependencias; no aprueba esta especificación nueva.

## 2. Inputs

- UUID del proyecto seleccionado por la aplicación.
- Sprint opcional, seleccionado de la lista del proyecto.
- Estado obtenido por `GET /api/projects/{id}`, misma fuente que US-04.
- Consultas de US-20/21/22: `/api/projects/{id}/metrics/story-points`,
  `/metrics/hours` y `/metrics/defects`, con `metrics` y `sprints` existentes.
- Para el Sprint seleccionado, reutilizar las rutas existentes
  `/api/projects/{id}/sprints/{sprintID}/metrics/{story-points|hours|defects}`.

## 3. Expected Outputs

Propuesta P-01: incorporar un dashboard en la página existente, manteniendo las
consultas numéricas de métricas. Mostrar nombre, UUID y estado oficial; agregado
de proyecto y detalle del Sprint seleccionado, con sus objetivos identificados.

Propuesta P-02: gráficos de barras horizontales separados por unidades:

| Grupo | Series |
| --- | --- |
| Story Points | planificados, completados, velocidad |
| Avance | porcentaje de historias completadas, escala 0–100 |
| Horas | estimadas y reales |
| Desviación | reales menos estimadas, eje con signo |
| Defectos | detectados y resueltos |

Cada serie conserva etiqueta y valor textual; `null`/ausente muestra **No
disponible**, sin barra ni conversión a cero. Un cero válido se etiqueta como 0.
La escala gráfica no recalcula ninguna métrica ni mezcla horas, puntos y porcentajes.
La desviación negativa se representa a la izquierda del origen, sin truncarla.
Los gráficos llevan alternativa textual y no dependen exclusivamente del color.

## 4. Business Rules

- BR-01: estado idéntico al de US-04; no derivarlo de métricas ni de cantidad de Sprints.
- BR-02: valores numéricos idénticos a los contratos de US-20/21/22 y sus SDD.
- BR-03: no sumar Sprints en JavaScript para reconstruir agregados de proyecto.
- BR-04: no tratar errores de red, campos faltantes o `null` como cero.
- BR-05: consultas de solo lectura; no alterar Sprints, snapshots, esfuerzo ni defectos.
- BR-06: al cambiar proyecto/Sprint, quitar resultados anteriores y descartar
  respuestas tardías de selecciones anteriores.
- BR-07: cambio de estado guardado/actualización del proyecto vuelve a consultar
  el dashboard; no agregar actualización automática periódica.

## 5. Constraints

Propuesta PT-01: reutilizar APIs/casos de uso Go existentes. No agregar un endpoint
agregador ni nuevos cálculos en Go: el alcance pendiente es presentación gráfica.
Si las pruebas revelan un defecto de las fuentes, documentarlo antes de ampliar alcance.

Propuesta PT-02: HTML/CSS/JavaScript nativos con gráficos SVG accesibles, sin
librerías externas, cambios de stack, migraciones ni dependencias nuevas.
Reutilizar el módulo de métricas cuando permita compartir lecturas y evitar duplicación.

Propuesta PT-03: carga y error visibles con reintento; ante un error del conjunto
consultado, limpiar sus gráficos y mostrar error sin presentar resultados parciales
como si fueran una carga completa. Proyecto y Sprint tienen contextos diferenciados.
Los endpoints existentes no garantizan snapshot atómico entre consultas; el
dashboard no promete esa consistencia ante escrituras concurrentes.

## 6. Edge Cases

Proyecto sin Sprints; Sprint sin historias; valores cero válidos; estimaciones
faltantes; velocidad sin Sprint cerrado; defecto sin atribución; desviación negativa;
nombres con caracteres HTML; cambio de proyecto/Sprint durante una carga;
desmarcar Sprint; error seguido de reintento; proyecto finalizado consultable.

## 7. Error Conditions

Reutilizar errores 400/404/500 de las APIs existentes. Ningún error se representa
como dato cero. No conservar gráficos del Sprint anterior al cambiar la selección,
incluso si la nueva consulta falla. No mostrar texto externo como HTML ejecutable.

## 8. Acceptance Criteria

- AC-01 oficial: dashboard con estado concordante con US-04 y gráficos basados
  en métricas existentes → BDD-US-23-01/02/04/06.
- AC-02 oficial: No disponible no se representa como cero → BDD-US-23-03.
- Validación de integración: datos de otra selección no aparecen en la actual,
  errores visibles y consulta sin escrituras → BDD-US-23-05/07.

Tests previstos: Node para valores/series, null frente a cero, signo, identidad,
errores, reintentos, selecciones tardías e integración de app; regresión Go completa
para fuentes y snapshots. No se ejecutaron tests nuevos ni hay evidencia TDD todavía.

## 9. Decisión pendiente

Aprobar P-01/P-02 y PT-01…PT-03 antes de implementar. No hay código funcional
nuevo. Rama local `codex/us-23-dashboard`. Publicación/PR y merge no realizados.
