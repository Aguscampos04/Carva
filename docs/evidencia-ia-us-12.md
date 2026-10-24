# Evidencia de desarrollo asistido por IA — US-12

## Objetivo

Implementar el registro de votos individuales ocultos de Planning Poker para la Issue [#12](https://github.com/Aguscampos04/Carva/issues/12), incluyendo especificación, escenarios, dominio Go, persistencia SQLite, API REST e interfaz web. Las dependencias US-03 y US-05 estaban cerradas al iniciar el trabajo.

## Instrucción recibida

> Bueno quiero que hagas la 12 termina esa y decime cuando termines mostrame y después hace un pull request, no mergees vos eso lo hace el agile del grupo. Acordate que quería tratar de que se cerrara sola y que crearas el pull request vos; la otra vez no pude. Cuando termines verificamos esa parte.

Durante el análisis, Loana confirmó que todos los miembros registrados del proyecto participan en cada ronda y que cada uno puede reemplazar su voto mientras la ronda siga abierta.

## Resultado y cambios

- SDD-US-12 y BDD-US-12 versionadas con trazabilidad a la Issue.
- Dominio de ronda y voto; snapshot de participantes al abrir la ronda inicial.
- Persistencia de rondas, participantes y votos mediante la migración SQLite `000011_planning_poker.sql`.
- Casos de uso para iniciar la ronda inicial, registrar/corregir votos y consultar solo el voto propio mientras la ronda está abierta.
- Rutas REST e interfaz de Planning Poker integrada en las tarjetas de User Story.
- Actualización de la documentación arquitectónica y de decisiones funcionales.

## Evidencia TDD

- **RED:** `go test ./internal/domain/planningpoker` falló antes de implementar `NewVote`, indicando que la entidad y sus validaciones no existían.
- **GREEN:** el mismo paquete pasó luego de agregar la entidad y validaciones.
- **REFACTOR:** se centralizó la consulta de participantes de una ronda en un helper de dominio; la prueba se ejecutará de nuevo junto con la suite final.

## Validación

- `gofmt -l .`: sin archivos pendientes de formato.
- `go test ./...`: PASS.
- `go vet ./...`: PASS.
- `go build ./...`: PASS.
- `go mod verify`: PASS (`all modules verified`).
- `git diff --check`: PASS.
- `node --test tests/web/*.test.cjs`: PASS, 16 de 16 pruebas.
- Las pruebas de integración usan SQLite temporal independiente; cubren ronda inicial, participantes del proyecto, atribución, ocultamiento, corrección, rechazo de valor negativo y preservación de Story Points definitivos.

## Revisión humana

Pendiente. No se afirma que otra persona haya revisado o aprobado el código.

## GitHub

Commit: pendiente de crear.
Push y Pull Request: pendientes de publicar y verificar. La futura descripción del PR incluirá `Closes #12` para que GitHub cierre la Issue al fusionarlo en la rama base.
