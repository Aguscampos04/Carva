# Evidencia de desarrollo asistido por IA — US-14

## Objetivo

Implementar la consulta de resultados históricos de Planning Poker y la creación de rondas sucesivas para la Issue [#14](https://github.com/Aguscampos04/Carva/issues/14), dependiente de US-13.

## Instrucción y decisiones aplicadas

La instrucción recibida fue continuar con la Issue 14. La Issue define que una nueva ronda se inicia después de revelar la anterior y que los votos/resultados previos deben seguir disponibles. Se aplicó la decisión funcional del contexto existente: participan los miembros registrados del proyecto y cada ronda conserva su propia lista de participantes. No se agregaron mensajería ni cambios del Story Point acordado (US-15).

## Resultado y cambios

- Se versionaron SDD-US-14 y BDD-US-14 con criterios para iniciar rondas sucesivas, preservar resultados y mantener votos ocultos.
- `POST /rounds` crea la primera ronda o la consecutiva cuando la última está finalizada; la transacción SQLite rechaza intentos duplicados o mientras haya una ronda abierta.
- La interfaz muestra el resultado de cada ronda finalizada, permite iniciar otra ronda y conserva visible el historial durante la votación nueva.
- Se actualizó `docs/contexto-proyecto.md` para corregir la numeración duplicada de las decisiones Planning Poker y aclarar que revelar diferencias no equivale a registrar una estimación acordada.

## Evidencia TDD

- **RED integración:** antes de implementar, el test de Issue 14 recibió HTTP 409 al solicitar la ronda 2 después de finalizar la primera; esperaba HTTP 201.
- **RED web:** antes de implementar, el botón para iniciar otra ronda seguía oculto después de revelar la ronda 1.
- **GREEN:** los tests comprueban la numeración consecutiva, conservación y consulta del resultado anterior, ausencia inicial de votos, ocultamiento del voto propio frente a los otros y rechazo de una tercera ronda mientras la segunda está abierta.

## Validación

- `go test ./...`: PASS en la repetición aislada. Una ejecución concurrente con `go vet` tuvo un panic intermitente en `TestProjectSprintCreationWinsBeforeFinalization` (`close of closed channel`); la repetición completa pasó.
- Tests de integración Planning Poker: PASS.
- `go vet ./...`: PASS.
- `gofmt -l .`: sin archivos pendientes.
- `node --test tests/web/*.test.cjs`: PASS, 18 de 18.
- `git diff --check`: PASS.

## Revisión humana y publicación

La revisión humana queda pendiente. US-13 sigue en el PR #42 abierto, por lo que US-14 se prepara como rama dependiente de `codex/issue-13`; no se realizó merge ni se cerró la Issue #14.
