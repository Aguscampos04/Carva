# SDD-US-20 — Consultar métricas de Story Points y avance

**Issue:** #20
**Estado:** fórmulas aprobadas por Loana Vargas el 2026-10-27; contrato técnico implementado en revisión.
**Trazabilidad:** enunciado §2.7 → EP-07 → US-20 / Issue #20 → este SDD → BDD-US-20 → tests y código.
**Dependencias:** US-05, US-08 y US-09 (cerradas).

## 1. Objetivo

Consultar, por Sprint y por proyecto, Story Points planificados, Story Points completados, velocidad del equipo y porcentaje de historias completadas, usando los snapshots de planificación y las finalizaciones registradas.

## 2. Inputs

- `projectID`: UUID obligatorio del proyecto en la ruta.
- `sprintID`: entero positivo obligatorio para la consulta de un Sprint.

## 3. Expected Outputs

- La API devuelve HTTP 200 con `planned_story_points`, `completed_story_points`, `team_velocity` y `completion_percentage`.
- Las métricas numéricas sin datos suficientes se representan como `null`; la interfaz las muestra como **«No disponible»**.
- La consulta del proyecto también devuelve las métricas por cada Sprint para seleccionar y consultar un Sprint desde la interfaz.
- Las lecturas no modifican Sprints, snapshots ni historias.

## 4. Business Rules

- **BR-01:** se utilizan exclusivamente los snapshots inmutables de `sprint_backlog_items` y su indicador de finalización. Las ediciones posteriores del Backlog no cambian la medición histórica.
- **BR-02:** Story Points planificados por Sprint son la suma de los Story Points definitivos de todas las historias planificadas en ese Sprint. Si una historia planificada no tiene Story Points definitivos, la métrica es no disponible.
- **BR-03:** Story Points completados por Sprint son la suma de los Story Points definitivos de las historias completadas en ese Sprint. Si alguna historia completada no tiene Story Points definitivos, la métrica es no disponible.
- **BR-04:** velocidad de un Sprint cerrado equivale a sus Story Points completados. Para un Sprint abierto, la velocidad es no disponible.
- **BR-05:** velocidad del proyecto es el promedio aritmético de Story Points completados en sus Sprints cerrados. Si no hay Sprints cerrados o algún Sprint cerrado tiene Story Points completados no disponibles, la velocidad del proyecto es no disponible.
- **BR-06:** porcentaje de historias completadas por Sprint es `historias completadas / historias planificadas × 100`. Un Sprint sin historias tiene porcentaje no disponible.
- **BR-07:** el porcentaje de proyecto es ponderado: `historias completadas en todos los snapshots / historias planificadas en todos los snapshots × 100`. Si no hay historias planificadas, es no disponible.
- **BR-08:** Story Points planificados y completados a nivel proyecto son la suma de las métricas de todos los snapshots de Sprint. Una historia trasladada se cuenta en cada Sprint cuyo snapshot la incluye. Si alguna historia relevante no tiene los Story Points definitivos necesarios, la suma afectada es no disponible.
- **BR-09:** un Sprint vacío tiene sumas de puntos iguales a cero; no tiene porcentaje de historias disponible. Para velocidad, un Sprint vacío cerrado aporta cero al promedio.
- **BR-10:** una Story Point de valor cero es un dato válido y no equivale a un dato ausente.

## 5. Constraints

- `GET /api/projects/{projectID}/metrics/story-points` devuelve `{ "metrics": {...}, "sprints": [{ "id", "goal", "closed", "metrics": {...} }] }`.
- `GET /api/projects/{projectID}/sprints/{sprintID}/metrics/story-points` devuelve las cuatro métricas del Sprint.
- Story Points se devuelven como enteros; porcentajes y velocidad del proyecto se devuelven como números decimales. La interfaz presenta los porcentajes con dos decimales.
- No se agregan migraciones ni dependencias.
- La interfaz permite consultar métricas del proyecto y elegir un Sprint existente.

## 6. Edge Cases

- Proyecto existente sin Sprints: todas las métricas del proyecto son no disponibles; la lista de Sprints es `[]`.
- Sprint vacío: Story Points planificados y completados son `0`; porcentaje no disponible; velocidad `0` solo si está cerrado.
- Historia planificada sin Story Points definitivos: la suma planificada correspondiente no está disponible.
- Historia completada sin Story Points definitivos: puntos completados y velocidad correspondientes no están disponibles.
- Historia incompleta sin Story Points no impide calcular puntos completados, pero sí puntos planificados.
- Cambios posteriores en Backlog o en el Story Point actual no alteran los snapshots.
- Proyecto o Sprint inexistente o ajeno al proyecto: HTTP 404.

## 7. Error Conditions

- `sprintID` no entero positivo: HTTP 400.
- Proyecto o Sprint inexistente/no perteneciente al proyecto: HTTP 404.
- Error inesperado de persistencia: HTTP 500 sin detalles internos.

## 8. Acceptance Criteria

- **AC-01:** Para un Sprint se presentan las cuatro métricas según las fórmulas aprobadas.
- **AC-02:** La agregación de proyecto conserva los datos del snapshot de cada Sprint aunque cambie después el Backlog.
- **AC-03:** Los datos insuficientes se muestran como «No disponible», nunca como cero ficticio.
- **AC-04:** Las consultas no modifican los datos fuente.

## 9. BDD / Test Traceability

- AC-01 → BDD-US-20-01/02/03 → pruebas de cálculo por Sprint, proyecto y datos incompletos.
- AC-02 → BDD-US-20-04 → prueba de invariancia ante edición posterior del Backlog.
- AC-03 → BDD-US-20-02/03 → pruebas con Story Points faltantes y sin historias.
- AC-04 → BDD-US-20-05 → prueba de lectura sin escrituras.
