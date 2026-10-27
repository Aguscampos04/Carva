# Evidencia de implementación — US-17

**Issue:** [#17 — Comparar esfuerzo estimado y real](https://github.com/Aguscampos04/Carva/issues/17)  
**Rama:** `codex/issue-17`  
**Base:** `main` con US-16 integrada en `fd8c0a8` (PR #45). La Issue #17 sigue abierta; este trabajo no la cierra ni hace merge.

## Alcance realizado

La comparación usa horas estimadas de la historia y horas reales calculadas desde sus registros activos asociados, como implementó US-16. US-17 hace explícita la presentación: la UI muestra cada valor en una fila independiente dentro de «Comparación de esfuerzo», con la unidad `h`.

No se agregan una resta ni una desviación porcentual: la Issue solo pide mostrar separadas las horas estimadas y reales; la agregación de métricas corresponde a US-21, cuyas fórmulas deben definirse en su SDD. Story Points se mantienen aparte y no se convierten a horas. Sin cambios de persistencia o migraciones.

## Trazabilidad

- Requerimiento: `docs/enunciado.md` §2.5 (esfuerzo) y §2.7 (métricas).
- Épica EP-05 → US-17 → Issue #17.
- SDD: `sdd/SDD-US-17-comparar-esfuerzo.md`.
- BDD: `bdd/BDD-US-17-comparar-esfuerzo.feature`.
- AC-01/AC-02 → escenarios BDD → `TestUS17EffortComparisonUsesOnlyRelatedActiveRecords` y test web `AC-US17-01/02`.
- Código de datos reutilizado: `internal/persistence/sqlite/backlog_repository.go`; presentación: `web/static/js/effort.js`.

## Evidencia TDD en commits

- `6aaef8e` documenta la SDD y BDD antes de implementar la mejora visual.
- `bff89e4` agrega la prueba de integración: 2.50 horas estimadas, 3.75 reales activas, exclusión de una entrada eliminada, aislamiento frente a otra historia y Story Points independientes.
- `d8a090f` agrega la prueba de interfaz. RED observado: falló porque no había términos `<dt>` para «Horas estimadas» y «Horas reales».
- `b501348` implementa la comparación accesible con términos y valores separados. Una primera ejecución encontró que el test no esperaba el sufijo `h`; se ajustó el resultado esperado a la unidad mostrada y pasó.
- `e624124` registra el formato `gofmt` de la prueba de integración.

No se reescribió el historial ni se simularon fallos.

## Verificaciones realizadas

| Verificación | Resultado |
| --- | --- |
| `gofmt -l .` | PASS, sin archivos listados después de formatear la prueba. |
| `go test ./... -count=1` | PASS; `tests/integration` pasó, incluida la prueba US-17. |
| `go vet ./...` | PASS. |
| Suite Node `tests/web/*.test.cjs` | PASS: 29 tests, 0 fallos. |
| `git diff --check` | PASS. |

## Commits de la rama

1. `6aaef8e docs: especificar comparación de esfuerzo US-17`
2. `bff89e4 test: verificar comparación horaria US-17`
3. `d8a090f test: exigir comparación visible de US-17` (RED)
4. `b501348 feat: mostrar comparación de esfuerzo US-17` (GREEN)
5. `e624124 style: formatear prueba de comparación US-17`

El equipo debe revisar y validar el resultado. La evidencia no implica aceptación final, cierre de la Issue ni merge.
