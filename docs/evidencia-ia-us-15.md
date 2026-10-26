# Evidencia de desarrollo asistido por IA — US-15

## Objetivo

Implementar la Issue [#15 — Registrar estimación acordada](https://github.com/Aguscampos04/Carva/issues/15): guardar como Story Points actuales el valor unánime de una ronda finalizada de Planning Poker.

## Instrucción y decisiones aplicadas

La instrucción recibida fue implementar la Issue #15 y dejar un Pull Request para que lo revise y fusione un compañero. Se consultaron `docs/enunciado.md`, `docs/contexto-proyecto.md`, `docs/reglas-agente.md`, `docs/arquitectura.md`, las skills del repositorio y las Issues #13–#15.

Se reutilizan las decisiones existentes: todos los miembros registrados participan en la ronda, la ronda conserva su lista de participantes y se identifica al participante seleccionado por su UUID, sin autenticación compleja. El valor se registra solo cuando una ronda finalizada es unánime. Revelar o iniciar otra ronda por sí solo no cambia el Story Point actual. El acuerdo conserva el valor provisional y escribe el Story Point actual en una transacción SQLite, sin sumar una migración ni un historial nuevo que la Issue no pide.

## Resultado y cambios

- `sdd/SDD-US-15-registrar-estimacion-acordada.md` especifica entradas, salidas, reglas, errores, límites y criterios de aceptación.
- `bdd/BDD-US-15-registrar-estimacion-acordada.feature` cubre consenso, cero, diferencias, revelación, ronda abierta y participante ajeno.
- El dominio extrae el valor acordado solo de resultados completos y finalizados; la persistencia actualiza el Backlog de forma atómica.
- La API añade `POST .../rounds/{roundID}/agreement?member_id=...` y devuelve `{"story_points": N}`.
- La UI muestra el botón cuando los votos revelados son unánimes y recarga el Backlog al registrar el resultado.
- No se añade una migración SQLite.

## Evidencia TDD

- **RED dominio:** `go test ./internal/domain/planningpoker` falló porque no existían `ErrNoConsensus` ni `RoundResult.AgreedStoryPoints`.
- **GREEN dominio:** después de implementar la regla, el paquete Planning Poker pasó; se probaron valor cero, unanimidad, diferencia y ronda abierta.
- **RED integración:** `TestPlanningPokerRegistersOnlyFinalizedUnanimousAgreement` recibió HTTP 404 al intentar registrar un acuerdo porque la ruta todavía no existía.
- **GREEN integración:** el test pasó al implementar la ruta y la transacción. Comprueba el guardado de Story Points, la conservación de provisionales, la ausencia de cambios al revelar/iniciar ronda, la conservación del valor previo con diferencias, el rechazo de una ronda abierta y el rechazo de un miembro ajeno.
- **RED interfaz:** el escenario BDD-US-15-01 falló antes de agregar el control, porque no existía el botón de acuerdo.
- **GREEN interfaz:** el test web pasó al mostrar y conectar el botón con la API y la recarga del Backlog.
- **RED historial de interfaz:** BDD-US-15-05 falló cuando se inició una ronda posterior porque la ronda unánime anterior no ofrecía la acción de registro.
- **GREEN historial de interfaz:** el control quedó disponible también en el resultado histórico y el test pasó; así iniciar otra ronda no impide registrar después el acuerdo anterior.

## Validación

| Validación | Resultado |
| --- | --- |
| `go test ./...` | PASS; todos los paquetes, incluida integración. |
| `go vet ./...` | PASS, sin diagnósticos. |
| `go build ./...` | PASS (exit code 0); Go mostró una advertencia no fatal al escribir su caché de estadísticas fuera del repositorio. |
| `go mod verify` | PASS: todos los módulos verificados. |
| `node --test tests/web/*.test.cjs` | PASS: 20/20. |
| `gofmt -l` sobre archivos Go modificados | PASS: sin archivos pendientes. |
| `gofmt -l .` | Reporta 12 archivos preexistentes fuera de esta Issue; ninguno de los archivos Go modificados por US-15 aparece en la lista. |
| `git diff --check` | PASS. |

La prueba focal de integración `go test -v -timeout 20s ./tests/integration -run '^TestPlanningPokerRegistersOnlyFinalizedUnanimousAgreement$' -count=1` pasó, y `node --test tests/web/planning-poker.test.cjs` pasó, 5/5.

## Revisión humana y publicación

La revisión humana del equipo queda pendiente y no se afirma como realizada. El commit, el push y el Pull Request se registrarán cuando se ejecuten; el PR deberá incluir `Closes #15` para que GitHub cierre la Issue al fusionarlo. No se realizará el merge ni se cerrará la Issue desde esta tarea.
