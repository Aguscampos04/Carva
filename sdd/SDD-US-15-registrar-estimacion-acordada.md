# SDD-US-15 — Registrar estimación acordada

**Estado:** implementación preparada para revisión; los criterios de la Issue #15 y las decisiones existentes de Planning Poker son la fuente del comportamiento.
**Trazabilidad:** enunciado §2.4 → EP-04 → US-15 / [Issue #15](https://github.com/Aguscampos04/Carva/issues/15) → esta SDD → [BDD-US-15](../bdd/BDD-US-15-registrar-estimacion-acordada.feature) → tests → Go.
**Dependencias:** US-13 y US-14, integradas en `main` mediante PR #42 y PR #43.

## Objective

**Integración US-16:** registrar acuerdo sigue modificando solo Story Points; no altera estimación de horas ni esfuerzo. Las horas se pueden editar aun con definitivo y la edición concurrente del Backlog conserva acuerdo y compleción vigentes. Véase [SDD-US-16](SDD-US-16-registrar-horas-esfuerzo.md).

Registrar los Story Points acordados para una historia a partir de una ronda de Planning Poker finalizada y unánime, conservando el valor actual si la ronda no alcanzó consenso.

## Inputs

- `projectID`: UUID del proyecto en la ruta.
- `storyID`: identificador de la User Story del Backlog en la ruta.
- `roundID`: identificador de una ronda de Planning Poker finalizada.
- `member_id`: identificador del miembro seleccionado de los participantes de esa ronda, enviado como query parameter.

## Expected Outputs

- El registro exitoso devuelve HTTP 200 y `{"story_points": N}`; el Backlog devuelve ese valor como Story Points actual y conserva el valor provisional existente.
- Si los votos no son unánimes, la operación devuelve HTTP 409 y no modifica el Story Point actual.
- Si la ronda no está finalizada, la operación devuelve HTTP 409 y no modifica el Backlog.
- Ronda, votos y actualización del Backlog se validan y persisten en una operación SQLite atómica.

## Business Rules

- **BR-01:** solo una ronda finalizada puede usarse para registrar un acuerdo.
- **BR-02:** existe consenso si todos los participantes registrados en la ronda tienen un voto y todos los votos tienen el mismo valor.
- **BR-03:** el valor acordado es el valor común de los votos; solo al registrarlo se actualiza `story_points` de la historia.
- **BR-04:** una ronda con diferencias no cambia `story_points`; se conserva exactamente su valor previo, incluso cuando todavía no existe.
- **BR-05:** registrar un acuerdo conserva el `provisional_story_points` original.
- **BR-06:** solo un participante de la ronda puede solicitar el registro. `member_id` identifica la selección funcional del MVP; no implica autenticación ni autorización fuerte.

## Constraints

- Mantener Go, REST, SQLite y las capas existentes; reutilizar `backlog_items.story_points`, sin una migración de esquema nueva.
- Contrato propuesto: `POST /api/projects/{projectID}/backlog/{storyID}/estimation/rounds/{roundID}/agreement?member_id={memberID}`.
- El handler valida formato y transporte; dominio/aplicación validan estado de la ronda, participación y consenso.
- La persistencia debe comprobar el resultado de la ronda y actualizar la historia dentro de una transacción inmediata para impedir que la actualización se base en votos o estado inconsistentes.
- No se cambia el estado de la ronda ni se modifica su historial de votos.

## Edge Cases

- Los votos `0, 0` son un consenso válido; Story Points aceptan enteros no negativos según US-12.
- Los votos `0, 1` no son consenso y conservan el valor actual.
- Una ronda abierta se rechaza, incluso si los votos cargados hasta ese momento coinciden.
- Un `member_id` ajeno a la ronda se rechaza sin cambios.
- Una historia o ronda inexistente dentro del proyecto no modifica ningún registro.

## Error Conditions

- `member_id` ausente o vacío: HTTP 400.
- Participante ajeno a la ronda: HTTP 400, consistente con US-12/US-13.
- Proyecto, historia o ronda inexistente/fuera del proyecto: HTTP 404.
- Ronda abierta/no finalizada o votos con diferencias: HTTP 409.
- Error de persistencia inesperado: HTTP 500 sin exponer detalles internos.

## Acceptance Criteria

- **AC-01:** si todos los votos de una ronda finalizada son iguales, registrar el acuerdo establece ese valor como Story Points actual de la historia.
- **AC-02:** si los votos no son todos iguales, la estimación no se registra y el Story Point actual permanece sin cambios.
- **AC-03:** revelar votos o iniciar otra ronda no actualiza Story Points por sí solo; solo el registro explícito de un acuerdo unánime lo hace.

## BDD / Test Traceability

- AC-01 → BDD-US-15-01/02 → tests de dominio e integración para unanimidad y persistencia.
- AC-02 → BDD-US-15-03 → tests de preservación del Story Point actual, tanto presente como ausente.
- AC-03 → BDD-US-15-04/05/06 → tests de ausencia de mutación al revelar/iniciar y rechazo de una ronda no finalizada.
