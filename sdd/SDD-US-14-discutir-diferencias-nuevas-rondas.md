# SDD-US-14 — Discutir diferencias y realizar nuevas rondas

**Estado:** derivada de US-14 / Issue #14; contrato técnico sujeto a revisión del equipo.
**Trazabilidad:** enunciado §2.4 → EP-04 → US-14 / [Issue #14](https://github.com/Aguscampos04/Carva/issues/14) → esta SDD → [BDD-US-14](../bdd/BDD-US-14-discutir-diferencias-nuevas-rondas.feature) → tests → Go.
**Dependencia:** US-13, publicada en la rama `codex/issue-13` y en el PR #42 (abierto).

## Objective

Permitir consultar los resultados de las rondas ya reveladas y abrir nuevas rondas de Planning Poker para una misma historia, conservando los resultados anteriores.

## Inputs

- `projectID`: UUID del proyecto en la ruta.
- `storyID`: identificador de la historia del Backlog en la ruta.
- `member_id`: identificador de un participante, requerido al consultar votos/resultados como en US-13.
- `roundID`: identificador de ronda para consultar los votos revelados de esa ronda.

## Expected Outputs

- `GET /api/projects/{projectID}/backlog/{storyID}/estimation/rounds` conserva y lista las rondas en orden ascendente por número.
- `POST /api/projects/{projectID}/backlog/{storyID}/estimation/rounds` crea la ronda inicial cuando no hay rondas, o la siguiente ronda cuando la anterior está finalizada; responde HTTP 201 con la ronda abierta.
- Cada ronda tiene un número consecutivo y su propia lista de votos/participantes.
- La interfaz permite consultar los resultados revelados de rondas anteriores y muestra el resultado de cada una por separado.
- La ronda nueva no muestra votos de participantes distintos del miembro seleccionado hasta finalizarse conforme a US-13.

## Business Rules

- **BR-01:** solo se puede iniciar la ronda siguiente si no existe ninguna ronda previa o si la ronda más reciente está finalizada.
- **BR-02:** no se permite más de una ronda abierta por historia.
- **BR-03:** los números de ronda son consecutivos, comenzando por 1.
- **BR-04:** cada nueva ronda inicia sin votos propios y conserva como participantes a los miembros registrados del proyecto al momento de abrirse, según la decisión de Planning Poker existente.
- **BR-05:** iniciar una ronda no altera ni elimina votos ni resultados de rondas anteriores.
- **BR-06:** los resultados de rondas finalizadas continúan consultables por sus participantes mediante el contrato de resultados de US-13.
- **BR-07:** los votos de una ronda abierta siguen las reglas de ocultamiento de US-12/US-13.

## Constraints

- Mantener Go, REST, SQLite y las capas existentes.
- Reutilizar `POST /rounds` para crear la ronda inicial o la siguiente; no agregar mensajería ni cambiar el resultado acordado del Backlog (US-15).
- Verificar dentro de una transacción que la ronda previa esté finalizada y que el número entrante sea el siguiente, para impedir dos rondas abiertas o numeración duplicada ante solicitudes concurrentes.
- Mantener accesibles los resultados previos desde la interfaz durante y después de rondas nuevas.

## Edge Cases

- Sin rondas previas, se crea la ronda 1.
- Con una ronda anterior abierta, una nueva solicitud de inicio se rechaza y no crea datos.
- Con una ronda anterior finalizada, se crea la siguiente con número consecutivo y sin votos.
- La ronda anterior y sus votos permanecen consultables después de iniciar otra.
- Si no hay miembros registrados, se conserva la respuesta existente HTTP 422.
- Dos solicitudes concurrentes no pueden crear dos rondas siguientes.

## Error Conditions

- Historia o proyecto inexistente: HTTP 404.
- Intento de iniciar otra ronda mientras la más reciente permanece abierta: HTTP 409.
- Historia sin miembros registrados: HTTP 422.
- Error de persistencia inesperado: HTTP 500 sin exponer detalles internos.

## Acceptance Criteria

- **AC-01:** después de revelar una ronda, sus votos y su indicador de diferencias siguen consultables para discutirlos.
- **AC-02:** se puede iniciar una nueva ronda para la misma historia solo después de finalizar la anterior; se asigna un número consecutivo y la anterior se conserva.
- **AC-03:** la nueva ronda comienza abierta y sin votos; los votos permanecen ocultos hasta su finalización, de acuerdo con US-13.
- **AC-04:** no se puede iniciar una segunda ronda mientras haya una ronda abierta para esa historia.

## BDD / Test Traceability

- AC-01 → BDD-US-14-01 → test de consulta/render de los resultados revelados anteriores.
- AC-02/03 → BDD-US-14-02 → test de integración que crea y enumera la segunda ronda y comprueba el ocultamiento de sus votos.
- AC-04 → BDD-US-14-03 → test de conflicto al intentar iniciar otra ronda abierta.
