# Evidencia de implementación — US-19

Fecha de validación: 2026-10-27. Issue [#19](https://github.com/Aguscampos04/Carva/issues/19).
Rama `codex/us-19-actualizar-defecto`, base `ed27692` de main (US-18 / PR #47 integrado).

## Alcance y revisión humana

Solicitud: implementar únicamente US-19. @Aguscampos04 aprobó la especificación
con «apruebo, acordate de hacer la US-19 nomas». La IA preparó especificación,
escenarios, pruebas, implementación y esta evidencia. La aprobación de diseño
no sustituye la revisión humana del código ni la validación final de la DoD.
No se publicaron cambios, no se creó PR ni se modificaron Issues.

Dependencias #18 y #7 verificadas cerradas, con DoD marcada e integración en main.
No se agregaron migraciones, dependencias ni funcionalidades de otras historias.

## Resultado y trazabilidad

[SDD](../sdd/SDD-US-19-actualizar-defecto.md) →
[BDD](../bdd/BDD-US-19-actualizar-defecto.feature) → dominio/aplicación/SQLite/API/UI.

PATCH parcial de estado y Sprints de detección/resolución. Los campos omitidos
se conservan y `null` quita una asociación. Identidad, historia, descripción y
severidad permanecen intactas. La escritura se realiza en una transacción
inmediata y preserva cambios concurrentes de campos diferentes.

| Criterio | Evidencia | Resultado técnico |
| --- | --- | --- |
| AC-01 asociación posterior de detección | `TestUS19UpdatesPersistAndPreserveOmittedFields` | PASS |
| AC-02 resolución en otro Sprint | mismo test, detección 1/resolución 2, reapertura SQLite | PASS |
| AC-03 consulta refleja cambios | mismo test, GET posterior; tests web de edición | PASS |

Errores y ámbito: `TestUS19RejectsInvalidAndForeignUpdates` cubre JSON inválido,
IDs ajenos/inexistentes, rollback, errores genéricos y conservación del registro.
`TestUS19ConcurrentUpdatesPreserveDistinctFields` cubre dos conexiones y ocho
rondas de cambios paralelos. `TestUS19UpdateHonorsContentionDeadline` verifica
el error de contención y deadline. Las reglas P-01…P-04 se ejercitan con estados,
desasociación, Sprint cerrado, historia completada y proyecto finalizado.
Los escenarios Gherkin se verifican mediante tests Go/Node, sin runner Gherkin.

## TDD realmente ejecutado

| Incremento | RED real | GREEN |
| --- | --- | --- |
| Dominio | `f17ba53`: tests fallaron por Patch/Update/ErrEmptyPatch inexistentes | `42d4285`: paquete dominio PASS |
| Aplicación/SQLite/HTTP | `4033ee1`: tests fallaron por Update inexistente | `afb0fb7`: tests focalizados PASS |
| Interfaz | `856f893`: cuatro tests fallaron por ausencia del botón Editar | `ee08076`: seis tests web US-18/19 PASS |
| Refactor | extracción de validación de Sprints compartida | `9378a65`: mismos tests Go y web PASS |

Comandos focalizados:

```text
go test ./internal/domain/defect ./internal/application/defects ./tests/integration -run 'TestPatch|TestUpdate|TestUS18|TestUS19' -count=1
node --test tests/web/defect*.test.cjs
```

El primer intento Node estuvo limitado por `spawn EPERM` del sandbox; se ejecutó
nuevamente con permisos adecuados y se obtuvo RED funcional antes del código.
Se agregó además una regresión para descartar una respuesta PATCH tardía tras
cambiar de contexto; pasa junto al resto de la suite.

## Validación final

| Verificación ejecutada | Resultado |
| --- | --- |
| `gofmt -l .` | PASS, sin archivos |
| `go test ./... -count=1` | PASS, incluye integración e historial de Sprints |
| `node --test tests/web/*.test.cjs` | PASS, 36 tests |
| `go vet ./...` | PASS, salida 0 |
| `go build ./...` | PASS, salida 0 |
| `go mod verify` | PASS, all modules verified |
| `git diff --check` | PASS |
| `go test -race ./internal/domain/defect ./internal/application/defects` | NO DISPONIBLE: requiere CGO, deshabilitado en este entorno |

Cobertura focalizada con `-coverpkg` sobre dominio, aplicación, persistencia y
HTTP, perfil local ignorado `.local/us19.cover`: Update dominio/aplicación 100%,
Update SQLite 84,2%, decoder PATCH 95%, helper de Sprints 88,9%. El total 25%
incluye funciones ajenas a US-19 que no ejecuta la selección focalizada;
no representa cobertura global del producto ni ausencia de defectos.

Smoke real: aplicación iniciada con SQLite temporal fuera de la base del usuario;
POST proyecto/historia/defecto, PATCH Abierto→Resuelto, GET confirma estado y UUID.
El JS de defectos se sirvió con HTTP 200. El primer intento omitió estimated_hours
del fixture y fue rechazado por el alta de Backlog; corregido el fixture, el flujo
completo pasó con errores PowerShell terminantes. Servidor detenido al terminar.
No se realizó validación visual humana en navegador; interfaz cubierta por Node.

## Verificación manual pendiente del responsable

1. Iniciar `go run ./cmd/carva` con una base de prueba.
2. Crear proyecto, historia y dos Sprints; registrar un defecto sin Sprint.
3. En defectos de esa historia, elegir **Editar defecto**.
4. Asociar detección al primer Sprint y resolución al segundo; guardar.
5. Volver a consultar y comprobar ambos IDs y el estado.
6. Cancelar otra edición y comprobar que no cambia el registro.
7. Probar un Sprint ajeno: debe verse el error y conservarse el borrador.

Pendientes: revisión humana, confirmación DoD, push manual y PR autorizado.
Cuando se confirme la DoD, el PR debe incluir `Closes #19` para el cierre al merge.
