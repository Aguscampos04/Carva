# SDD-US-13 — Revelar votos y detectar diferencias

**Estado:** la regla de finalización con todos los votos fue aprobada por el responsable el 2026-10-26; contrato técnico sujeto a revisión del equipo.
**Trazabilidad:** enunciado §2.4 → EP-04 → US-13 / [Issue #13](https://github.com/Aguscampos04/Carva/issues/13) → esta SDD → [BDD-US-13](../bdd/BDD-US-13-revelar-votos-detectar-diferencias.feature) → tests → Go.
**Dependencia:** US-12 (registro de votos individuales ocultos), integrada en `main` mediante PR #41.

## Objective

Finalizar una ronda de Planning Poker para revelar a sus participantes, al mismo tiempo, todos los votos registrados y señalar si sus Story Points difieren.

## Inputs

- `projectID`: UUID del proyecto en la ruta.
- `storyID`: identificador de la historia del Backlog en la ruta.
- `roundID`: identificador de la ronda en la ruta.
- `member_id`: identificador del participante que solicita finalizar o consultar el resultado, enviado como query parameter.

## Expected Outputs

- La finalización exitosa devuelve HTTP 200 con la ronda finalizada, todos sus votos y `has_differences`.
- La consulta posterior del resultado devuelve los mismos votos y la misma indicación de diferencias.
- Antes de finalizar, los participantes siguen viendo únicamente su propio voto mediante el endpoint existente de votos.
- Si todos los valores son iguales, `has_differences` es `false`; si alguno difiere, es `true`.

## Business Rules

- **BR-01:** la ronda solo puede finalizarse cuando todos sus participantes registrados en la ronda tienen un voto.
- **BR-02:** la transición de abierta a finalizada es atómica respecto del registro de votos.
- **BR-03:** la respuesta de finalización y la consulta de resultados muestran el conjunto completo de votos a los participantes de la ronda; ningún participante ve un subconjunto revelado antes de la transición.
- **BR-04:** existe diferencia si y solo si al menos dos votos tienen valores distintos.
- **BR-05:** finalizar la ronda no modifica los Story Points del Backlog; registrar el valor acordado corresponde a US-15.
- **BR-06:** una ronda finalizada no admite nuevos votos ni correcciones.
- **BR-07:** únicamente un miembro participante puede finalizar la ronda o consultar su resultado. No se agrega autenticación ni permisos al MVP.

## Constraints

- Mantener Go, REST, SQLite y las capas existentes; la columna `status` de `estimation_rounds` ya contempla `open` y `finalized`.
- Contrato propuesto:
  - `POST /api/projects/{projectID}/backlog/{storyID}/estimation/rounds/{roundID}/finalize?member_id={memberID}` finaliza y devuelve el resultado.
  - `GET /api/projects/{projectID}/backlog/{storyID}/estimation/rounds/{roundID}/result?member_id={memberID}` consulta el resultado después de finalizar.
- El resultado contiene `round`, `votes` y `has_differences`; cada voto conserva `member_id` y `story_points`.
- La decisión sobre todos los participantes debe votar fue confirmada por el responsable el 2026-10-26. La ronda se finaliza explícitamente cuando el equipo decide revelar.
- Crear rondas posteriores, registrar un nuevo valor acordado y actualizar Story Points quedan fuera de US-13 (US-14 y US-15).

## Edge Cases

- Con algún voto pendiente, la ronda permanece abierta y no se revela ningún voto ajeno.
- Una sola ronda con todos los votos iguales produce `has_differences: false`.
- Una ronda con al menos dos valores diferentes produce `has_differences: true`.
- Un participante inexistente o ajeno a la ronda no obtiene resultados.
- Una solicitud repetida de finalización sobre una ronda ya finalizada se rechaza como conflicto; consultar el resultado sigue permitido.

## Error Conditions

- `member_id` ausente o vacío: HTTP 400.
- Proyecto, historia o ronda inexistente o fuera del alcance del proyecto: HTTP 404.
- Participante que no pertenece a la ronda: HTTP 400, siguiendo el contrato existente de US-12.
- Votos pendientes, ronda ya finalizada al volver a finalizarla, o consulta antes de la revelación: HTTP 409.
- Error de persistencia inesperado: HTTP 500 sin exponer detalles internos.

## Acceptance Criteria

- **AC-01:** antes de finalizar una ronda, los votos de otros participantes siguen ocultos; al finalizarla, todos los votos se revelan juntos.
- **AC-02:** valores iguales producen `has_differences: false`; si algún valor difiere, producen `has_differences: true`.
- **AC-03:** la ronda no se puede finalizar mientras falte el voto de algún participante.

## BDD / Test Traceability

- AC-01 → BDD-US-13-01/02 → tests de visibilidad antes y después de finalizar.
- AC-02 → BDD-US-13-03/04 → tests de unanimidad y diferencia.
- AC-03 → BDD-US-13-05 → test de intento de finalización con votos pendientes.
