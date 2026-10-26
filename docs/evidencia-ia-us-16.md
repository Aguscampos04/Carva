# Evidencia de implementación — US-16

**Issue:** [#16](https://github.com/Aguscampos04/Carva/issues/16). **PR único:** [#45](https://github.com/Aguscampos04/Carva/pull/45), rama `codex/docs-us16-horas-esfuerzo`. Sin merge ni cierre. Aceptación final por el responsable pendiente.

El responsable aprobó formalmente las 31 decisiones funcionales (incluida D-26 corregida) y PT-01…PT-06 y autorizó implementación, automatización, commits y push. La aprobación se registra en SDD, BDD y contexto. No se atribuyen reuniones ni aceptación final.

## Resultado implementado

- Go conserva horas como centésimas enteras, sin float ni conversión a Story Points. Nuevas historias y ediciones requieren estimación; legado NULL sigue consultable y admite esfuerzo salvo D-30.
- Servicios de esfuerzo, SQLite y REST implementan alta, consulta, edición con UUID estable, asociación opcional y baja lógica confirmada. Se validan relaciones dentro del proyecto, fechas, texto, precisión y máximo activo de 24 horas por miembro/día/proyecto.
- D-30 bloquea toda alta/edición en Finalizado sin fin, pero permite consulta/baja. D-31 controla inicio únicamente contra activos, sin modificar eliminados. D-26 conserva Sprint compatible o nulo; incompatible requiere elección explícita.
- Edición de Backlog y propagación de estimación son atómicas; los demás datos históricos permanecen. La asignación relee horas vigentes bajo bloqueo, y cierre/traslado conserva esfuerzo en su asociación elegida. Acuerdos y compleción concurrentes se conservan.
- Migración aditiva `000012_effort.sql`: columnas anulables históricas, tabla de esfuerzo, FK compuestas e índices. No se alteran migraciones aplicadas ni se inventan fechas/estimaciones.
- Interfaz en detalle del Backlog: estimación/total, selección nombre+UUID, fecha del servidor, edición, resolución de Sprint y confirmación de baja. Preserva centésimas también al leer y enviar JSON en JavaScript. Respuestas tardías se descartan.

## Evidencia TDD real

| Incremento | RED observado | Implementación / refactor verificables |
| --- | --- | --- |
| `83aa25d` | `go test ./internal/domain/hours ./internal/domain/effort`: FAIL por Value, Input y reglas aún inexistentes | `fc1f838`: tipos exactos y reglas; mismo comando PASS. |
| `9c73126` | `go test ./tests/integration -run TestUS16 -count=1`: FAIL por paquete application/effort aún inexistente | `11e4b4d`: capas, SQL y API; tests US-16 y suite Go completa PASS tras adaptar fixtures requeridos. |
| `73d459e` | `node --test tests/web/effort.test.cjs`: 5 FAIL por effort.js inexistente | Widget implementado; pruebas de detalle, UUID, D-26, D-30 y baja pasan. |
| `834086c` | Caso adicional de parseo exacto: FAIL `parseHoursJSON is not a function` | `7fe7a1e` agrega interfaz/parser; la regresión de app detecta referencia `selectedProject` ausente y se corrige en `16548ce`. No se presenta aquel commit intermedio como suite completa en verde. |
| `78718a5` | `TestUS16SnapshotOwnsEstimateValue`: FAIL por alias del puntero de horas | `16548ce`: copia independiente, parser simplificado y correcciones de concurrencia; tests unitarios PASS. |

Durante regresión completa se observó `panic: close of closed channel` en un test previo de US-02: su callback no toleraba reintentos, pese a que el puerto ya los admite. Se corrigió el doble de sincronización con `sync.Once`, conservando sus aserciones y comportamiento esperado.

La repetición con cobertura también detectó un 500 intermitente al editar estimación contra compleción: la comprobación de proyecto se hacía antes de la transacción. Se movió al puerto `BacklogRepository.Mutate` dentro de `withImmediate`, sin sobrescribir estado/puntos. El grupo `TestUS16ConcurrentPlanning` pasó después 20 repeticiones; tras añadir destino ya asignado y pools independientes, pasó otras 10. No se ocultan estos fallos intermedios ni se reescribe el historial. CI inicial también detectó en `TestProjectConcurrentFutureStartAndFinalization` la suposición incorrecta de que obtener primero el bloqueo garantiza el primer COMMIT. Se ajustó el test para validar exactamente uno de los dos resultados serializados válidos (incluidos códigos/errores de la operación rechazada y fechas persistidas), contemplando rollback/reintento; no se cambió la regla de negocio.

## Verificaciones locales efectivamente realizadas

Entorno: Windows/PowerShell, Go 1.27.1, Node 20.17.0; módulo conserva mínimo Go 1.26.0. Bases de integración temporales por test, con pools SQLite independientes para concurrencia.

| Verificación | Resultado |
| --- | --- |
| `gofmt -l .` | PASS, sin archivos listados. |
| `go test ./... -count=1` | PASS en la regresión completa. |
| `go test ./... -count=1 '-coverpkg=./internal/...' '-coverprofile=.local/us16-coverage.out'` | PASS; cobertura agregada del núcleo `internal/...`: **84,6 %**. |
| `go vet ./...` | PASS. |
| `go build ./...` | PASS. |
| `node --test (Get-ChildItem tests/web/*.test.cjs).FullName` | PASS: **28 tests**, sin fallos/cancelados/omitidos. |
| `git diff --check` | PASS, sin errores de whitespace. |

Los argumentos de cobertura van entre comillas en PowerShell. Un intento inicial sin ellas no generó el perfil esperado; se repitió con argumentos explícitos y la cifra se obtuvo de `go tool cover '-func=.local/us16-coverage.out'`, no de estimaciones ni del porcentaje aislado de un paquete. El perfil local está ignorado por Git; CI conserva un artefacto propio.

CI se configura en `.github/workflows/ci.yml` para push a `codex/**` y PR a main: formato, tests/cobertura, vet, build y frontend. La ejecución inicial [40417571666](https://github.com/Aguscampos04/Carva/actions/runs/40417571666) completó correctamente. La ejecución push [40417554737](https://github.com/Aguscampos04/Carva/actions/runs/40417554737) detectó la suposición de orden del test previo descrita arriba. Tras corregir los tests de US-02, ambos casos pasaron 30 repeticiones locales. El estado del commit final y su URL se verifican y publican en la descripción del PR y en [sus checks](https://github.com/Aguscampos04/Carva/pull/45/checks); no se confunde configuración con CI exitosa.

## Matriz de los 46 comportamientos BDD y pruebas reales

El archivo `.feature` permanece como especificación Given–When–Then. **No se ejecuta un runner Gherkin**: los comportamientos se automatizan con testing de Go y el runner Node existentes, sin agregar dependencias. Los esquemas de ejemplos se cubren mediante tablas de validación y límites en esas pruebas; esta matriz relaciona comportamiento y regresión, no afirma que cada línea del feature haya sido interpretada por una herramienta.

Abreviaturas: `I` = [effort_test.go](../tests/integration/effort_test.go); `R` = [effort_regression_test.go](../tests/integration/effort_regression_test.go); `W` = [effort.test.cjs](../tests/web/effort.test.cjs); `H` = [hours/value_test.go](../internal/domain/hours/value_test.go); `D` = [effort/entry_test.go](../internal/domain/effort/entry_test.go); `A` = [application/effort/service_test.go](../internal/application/effort/service_test.go).

| BDD-US-16 | Prueba / comportamiento automatizado |
| --- | --- |
| 01 | R `TestUS16EstimatesAndRegularization`; H `TestUS16ExactHours`: positivos, 0.01, decimales, sin puntos/Sprint. |
| 02 | R `TestUS16EstimatesAndRegularization`: esfuerzo y SP conservados al cambiar horas. |
| 03 | R `TestUS16EstimatesAndRegularization`: edición sin Sprint. |
| 04 | R `TestUS16EstimatesAndRegularization`; H: ausencia, NULL, cero, negativo y precisión. |
| 05 | I `TestUS16EffortLifecycle`; W detalle pendiente de regularización. |
| 06 | I `TestUS16SnapshotsSprintAndStart`; R migración y planificación: regularizar propaga sin reescribir plan. |
| 07 | I ciclo de vida y W detalle/envío. |
| 08 | I `TestUS16InvalidEffortAndDates`; D campos requeridos. |
| 09 | I fechas inválidas; D fechas y `project.Date` reutilizado. |
| 10 | I JSON inválido; R `TestUS16ScopedResources`; W serialización. |
| 11 | I ciclo de vida (miembro n); W selección de UUID entre homónimos. |
| 12 | R `TestUS16ScopedResources`: miembro, historia, proyecto y Sprint. |
| 13 | A `TestUS16ServerDayAndApplicationRules` y W: frontera 02:30Z/03:00Z y fecha predeterminada. |
| 14 | D `TestUS16EffortRules`; I límites temporales. |
| 15 | I snapshots/inicio: PATCH rechazado íntegro e inicio borrable. |
| 16 | I snapshots/inicio y R concurrencia delete-start: eliminados dejan de bloquear, datos preservados. |
| 17 | R `TestUS16FinalizedEffort`: Completada, fin conocido y edición/baja sin cambios de estado. |
| 18 | R finalizados y D fechas: límite del fin. |
| 19 | I ciclo de vida y R estimaciones: múltiples entradas y real superior al estimado. |
| 20 | I `TestUS16ConcurrentDailyLimit`; R reemplazos: límite 24 exacto y exceso. |
| 21 | R `TestUS16ReplacementAndSprintValidation`: cupos por miembro/día/proyecto. |
| 22 | R reemplazos: excluir el registro que se está editando. |
| 23 | I ciclo de vida y R reemplazos: cambiar miembro/fecha/historia y rechazar cupo lleno. |
| 24 | I baja repetida y almacenamiento histórico; W diálogo confirmable. |
| 25 | W confirmación falsa: sin solicitud ni guardado de borrador; API exige true. |
| 26 | I snapshots y R reemplazos; W conservación compatible/resolución incompatible. |
| 27 | I aislamiento y R recursos: ni reasignación ni baja desde otro proyecto. |
| 28 | I snapshots/Sprint cerrado, con esfuerzo sin reapertura. |
| 29 | R `TestUS16ConcurrentPlanning/transfer`: esfuerzo conserva Sprint origen y total único. |
| 30 | R `TestUS16HistoricalMigrationAndReopen`: esquema previo, snapshots abiertos/cerrados, migración repetida, reapertura y baja persistida. |
| 31 | R `TestUS16Rollback`: triggers de fallo y rollback de alta, edición y baja. |
| 32 | I `TestUS16ConcurrentDailyLimit`: conexiones independientes. |
| 33 | R `TestUS16ConcurrentLifecycle/start`: inicio contra alta. |
| 34 | R planificación `assign`, `transfer`, `transfer-existing`: horas sincronizadas. |
| 35 | R concurrencia `delete-edit`: sin resurrección. |
| 36 | R rollback de propagación a snapshot cerrado. |
| 37 | R planificación `agreement` y `complete`: usa operaciones reales de US-15/09. |
| 38 | R `TestUS16BoundedContention`: espera acotada y HTTP 503. |
| 39 | W respuestas tardías; regresión de wiring en projects.test.cjs. |
| 40 | R finalizados: edición revalida fecha sin sobrescritura parcial. |
| 41 | R concurrencia `finalize`: ambas operaciones válidas conservadas. |
| 42 | R concurrencia `edit-budget`: solo una edición consume cupo destino. |
| 43 | R finalizados, I validación de estados y W bloqueo UI por fin desconocido. |
| 44 | R finalizados: consulta y baja con fin NULL; W mantiene consulta/baja disponibles. |
| 45 | R concurrencia `delete-start`: cambio de inicio atómico sobre activos vigentes. |
| 46 | I ciclo de vida sin Sprint al cambiar historia; W selector y contrato de conservación. |

## Prueba funcional en navegador y límites

Se inició `go run ./cmd/carva` con `CARVA_DB_PATH=.local/us16-smoke.db`, sin modificar la base de desarrollo. Se creó un proyecto y miembro de prueba, una historia con estimación 2.50 sin SP, y esfuerzo real 3.25 sin Sprint. La pantalla mostró ambos datos independientes y el miembro por UUID. Se editó el esfuerzo a 1.50 y la consulta actualizó el total. Tras la baja, la pantalla mostró total 0.00. El control de confirmación/cancelación y el bloqueo D-30 tienen además pruebas Node/HTTP reproducibles.

La captura local `.local/us16-smoke.png` corresponde al estado con estimación 2.50 y real 3.25. Es evidencia visual local, ignorada por Git. No se afirma un recorrido manual en navegador de todas las combinaciones ni un análisis de accesibilidad completo; las carreras, legado y fallos se verifican automáticamente con datos controlados.

No se agregan métricas, dashboards, reportes, autenticación, reapertura ni interfaz de eliminados. Los clientes anteriores que omiten estimación en POST/PUT deben actualizarse según el contrato aprobado. No hay decisiones funcionales/técnicas pendientes; revisión humana final e integración siguen pendientes.
