# SDD-US-16 — Registrar horas estimadas y esfuerzo real

**Estado:** implementación autorizada por el responsable; revisión final pendiente; 31 decisiones funcionales D-01…D-31 aprobadas por el responsable, incluida la corrección de D-26 y las resoluciones D-30/31 incorporadas en la actualización del PR #45. El responsable aprobó formalmente las 31 decisiones y PT-01…PT-06 y autorizó implementación completa en la misma rama del PR #45, con SDD → BDD → TDD y evidencia real. La aceptación final e integración siguen reservadas al responsable.
**User Story:** Como integrante del equipo, quiero registrar horas estimadas para una User Story durante la planificación y horas trabajadas en el registro de esfuerzo para disponer de ambos datos.
**Trazabilidad:** [enunciado](../docs/enunciado.md) §2.5 y soporte de datos para §2.7 → EP-05 → [US-16 / Issue #16](https://github.com/Aguscampos04/Carva/issues/16) → [decisiones](../docs/contexto-proyecto.md#horas-estimadas-y-esfuerzo-real--us-16--issue-16) → esta SDD → AC → [BDD-US-16](../bdd/BDD-US-16-registrar-horas-esfuerzo.feature) → futuros tests TDD → futuro código Go.
**Base revisada:** `main` actualizado desde `origin/main`, commit `a128429e5e72209e0455644d2f6e3e09c83f0544`; rama `codex/docs-us16-horas-esfuerzo`. Issue #16 abierta, EP-05, prioridad importante; su cuerpo conserva SP y Sprint pendientes. No se infieren ceremonias ni aprobaciones finales.
**Fuentes:** AGENTS.md; skills efectivas en `.agents/skills/` (project-manager, sdd, bdd, github); enunciado, contexto, reglas y arquitectura; SDD/BDD US-02/03/05/08/10/11/15 y código de Backlog, miembros, proyectos, Sprints, Planning Poker, SQLite y UI. US-06/07/09 aportan listado, asignaciones y resultados usados aquí. La arquitectura es Go → REST → SQLite, con presentación, aplicación, dominio y persistencia separados.

## Objective

Mantener una estimación vigente de horas por historia y administrar registros independientes de esfuerzo real desde su detalle en el Product Backlog. Consultar la estimación, el total activo y su detalle, sin derivar horas desde Story Points ni alterar automáticamente estimación, estados o resultados históricos al cargar esfuerzo.

Alcance: alta/edición de estimaciones mediante US-05, consulta y regularización histórica, alta/consulta/edición/baja lógica de esfuerzo, identidad de US-03, validaciones de fechas de US-02, propagación de estimación a snapshots y coordinación transaccional. PT-01…PT-06 concretan las decisiones y fueron aprobadas formalmente para implementar.

Exclusiones: métricas de desviación, velocidad, dashboards y reportes de otras US; conversión horas/SP; reapertura de proyectos/Sprints; edición de resultados históricos; gestión nueva de miembros; autenticación/permisos; restauración o interfaz de eliminados; bitácora de todas las versiones de cada edición. La etapa de implementación está autorizada; se conserva el alcance funcional aprobado.

## Inputs

| Entrada | Tipo / obligatoriedad propuesta | Regla |
| --- | --- | --- |
| `projectID` | String de ruta, UUID emitido por US-01 | Proyecto existente; determina el ámbito y no se cambia en el cuerpo. ID que no resuelve: 404, sin nuevo validador global de UUID. |
| `itemID` | String de ruta, UUID de US-05 | Historia existente del proyecto para detalle/alta. |
| `entryID` | String de ruta, UUID generado por servidor | Identidad del esfuerzo para consulta/edición/baja. |
| `estimated_hours` | Número decimal JSON, obligatorio en POST/PUT de US-05 | Mayor que cero, hasta dos decimales. Se agrega a los campos editables actuales; no depende de tener Sprint, esfuerzo o SP definitivo. No se permite borrar una estimación vigente. |
| `member_id` | String UUID, obligatorio en esfuerzo | Selección de `GET /api/projects/{projectID}/members`; validar existencia y pertenencia, nunca resolver por nombre. |
| `work_date` | String `YYYY-MM-DD`, obligatorio en esfuerzo | Día de calendario válido, actual o pasado; respetar inicio y fin conocido. La UI completa hoy; omisión en API se rechaza. |
| `activity` | String, obligatorio | Texto libre no blanco según `strings.TrimSpace`; se propone conservar el texto original, sin máximo funcional ni interpretación HTML. |
| `hours` | Número decimal JSON, obligatorio en esfuerzo | Mayor que cero y hasta dos decimales; máximo acumulado miembro/día/proyecto 24. |
| `backlog_item_id` | String UUID obligatorio en PUT de esfuerzo | Historia destino del mismo proyecto. En POST se toma de la ruta, sin dos fuentes de historia. |
| `sprint_id` | Entero positivo o `null` | Opcional en alta (omitir = sin Sprint). En PUT, omitir conserva la asociación persistida si es válida para la historia destino, incluso si es nula. Número o `null` expresan una elección explícita. Un Sprint no nulo requiere asignación historia/Sprint del proyecto, incluso cerrado; si el anterior es incompatible, no se permite omitirlo para resolver el conflicto. |
| `confirmed` | Booleano requerido en baja lógica | Solo `true`; cancelar el diálogo no envía escritura. |
| Hoy | Reloj inyectable del servidor | Día en `America/Argentina/Buenos_Aires`, ya autorizado e implementado por US-02; no depende del navegador ni zona del host. |

Los campos editables previos de US-05 mantienen sus validaciones. `estimated_hours` ausente, nulo o inválido invalida el alta y toda edición genérica, incluso si solo cambia un título histórico. Completar una historia por US-09 o registrar acuerdo por US-15 son acciones específicas, no una edición genérica de US-05: no se impone estimación retroactiva en esos comandos. Registrar esfuerzo tampoco la exige (D-28).

## Expected Outputs

- Historia consultable con su UUID, datos actuales, `estimated_hours` (número o `null` únicamente en legado) y `actual_hours` (suma de registros activos, cero si no hay). La UI muestra «Pendiente de regularización» para estimación nula; nunca sustituye ese faltante por cero. El cero del esfuerzo corresponde a ausencia real de registros activos, no a una métrica sin datos.
- Alta de esfuerzo: UUID nuevo, proyecto e historia obligatorios, miembro, fecha, actividad, horas y Sprint opcional persistidos. Edición conserva UUID y proyecto; puede cambiar historia dentro del proyecto y los otros campos editables.
- Baja confirmada: registro almacenado con marca de eliminación, oculto de consultas normales y excluido de sumas. No hay eliminación física ni restauración implícita.
- Cambiar estimación actualiza todos sus snapshots, abiertos y cerrados, en la misma operación. El resto del plan, resultados y metadatos de Sprint permanece igual.
- Errores controlados sin éxito aparente ni escrituras parciales. Lecturas y reinicio con la misma base conservan identidad, relaciones, importes exactos y marcas de eliminación.

## Business Rules

IDs D con alcance exclusivo a US-16. Cada BR corresponde a la D del mismo número; el texto íntegro aprobado se conserva en contexto.

| Regla | Comportamiento verificable | AC |
| --- | --- | --- |
| BR-01 / D-01 | Cualquier historia del Product Backlog admite estimación sin exigir Sprint. | AC-01 |
| BR-02 / D-02 | Estimación modificable en cualquier momento, aun con esfuerzo; se conserva la vigente. | AC-01 |
| BR-03 / D-03 | Todo esfuerzo pertenece obligatoriamente a una historia del proyecto. | AC-04 |
| BR-04 / D-04 | Esfuerzo permitido con o sin asignación de la historia a Sprint. | AC-04 |
| BR-05 / D-05 | Horas reales > 0, máximo dos decimales. | AC-05 |
| BR-06 / D-06 | Estimaciones admiten hasta dos decimales. | AC-01, AC-02 |
| BR-07 / D-07 | Nueva historia requiere estimación > 0 desde el alta. | AC-02 |
| BR-08 / D-08 | Historia histórica sin estimación sigue consultable; editarla exige regularizar. | AC-03 |
| BR-09 / D-09 | Fecha de trabajo ≤ hoy del servidor. | AC-07 |
| BR-10 / D-10 | Actividad obligatoria, no vacía ni solo espacios. | AC-05 |
| BR-11 / D-11 | Editar miembro, fecha, actividad, horas e historia conserva ID. | AC-10 |
| BR-12 / D-12 | Eliminar requiere confirmación; cancelar no modifica. | AC-11 |
| BR-13 / D-13 | Baja lógica conserva fila; consultas y sumas usan solo activos. | AC-11, AC-12 |
| BR-14 / D-14 | Historia Completada admite alta, edición y baja sin cambiar estado. | AC-13 |
| BR-15 / D-15 | Proyecto Finalizado admite las tres operaciones sin cambiar estado/fechas. | AC-08, AC-13 |
| BR-16 / D-16 | Varias entradas mismo miembro/día/historia son independientes y acumulables. | AC-09 |
| BR-17 / D-17 | En Finalizado, fecha de trabajo ≤ fecha real de fin. | AC-08 |
| BR-18 / D-18 | Fecha de trabajo ≥ inicio si existe; rechazar cambio posterior de inicio incompatible, sin mover registros. | AC-07 |
| BR-19 / D-19 | Cada historia muestra estimación, total activo y acceso al detalle. | AC-12 |
| BR-20 / D-20 | Administrar esfuerzo desde detalle del Product Backlog. | AC-12 |
| BR-21 / D-21 | Cambiar estimación propaga a todos los Sprints asignados, incluidos cerrados; ningún otro dato histórico cambia. | AC-16 |
| BR-22 / D-22 | Cada entrada puede tener un Sprint de la historia o ninguno. | AC-14 |
| BR-23 / D-23 | CRUD lógico de esfuerzo en Sprint cerrado no reabre ni cambia su plan. | AC-15 |
| BR-24 / D-24 | Suma activa miembro/fecha/proyecto ≤ 24 incluyendo todas las historias y Sprints. | AC-09 |
| BR-25 / D-25 | Real puede superar estimado sin modificarlo automáticamente. | AC-09 |
| BR-26 / D-26 | Al cambiar historia conservar el Sprint anterior si sigue siendo válido; si es incompatible exigir otro Sprint válido o quitar explícitamente la asociación. Nunca modificarla automáticamente. | AC-10, AC-14 |
| BR-27 / D-27 | No reasignar entre proyectos; corregir con baja en origen y nueva alta en destino. | AC-06, AC-10 |
| BR-28 / D-28 | Admitir esfuerzo en historia histórica sin estimación; conservar indicador pendiente. | AC-03 |
| BR-29 / D-29 | Formulario de alta propone fecha actual editable bajo las restricciones. | AC-07, AC-12 |
| BR-30 / D-30 | En proyecto Finalizado con `end_date = NULL`, bloquear toda alta/edición de esfuerzo; permitir consultas y baja lógica sin inventar fin ni cambiar estado. | AC-19 |
| BR-31 / D-31 | Validar D-18 únicamente contra registros activos; eliminados no bloquean inicio y conservan datos históricos; validación y cambio son atómicos. | AC-07, AC-18 |

Identificación adicional aprobada: reutilizar US-03, UUID seleccionado de la lista del proyecto; homónimos son identidades distintas. Estimación, esfuerzo y Story Points son independientes. D-30 precisa la excepción a D-15 para Finalizado sin fin; D-31 precisa los registros considerados en D-18. No se agregan otras prohibiciones por estado ni reglas de fecha de Sprint: hoy los Sprints no tienen fechas de inicio/fin.

## Casos de uso

1. **Estimar/regularizar:** abrir alta o edición de historia, completar campos US-05 y estimación válida; guardar todo o nada; propagar solo horas estimadas a asignaciones existentes. Cancelar conserva datos. Se permite modificar estimación aunque haya SP definitivo, esfuerzo, historias completadas o Sprints/proyecto cerrados.
2. **Registrar esfuerzo:** desde detalle de historia cargar miembros del proyecto, proponer hoy y Sprint sin asociación; elegir miembro, fecha, actividad y horas; opcionalmente Sprint de la historia; validar ámbito, fechas y acumulado vigente y guardar; refrescar detalle/total tras confirmación del servidor.
3. **Consultar:** lista de Backlog presenta estimación/total y abre detalle con entradas activas; sin entradas mostrar total 0 y lista vacía. No confundir error de consulta con vacío. Lectura no regulariza, reasigna ni modifica snapshots.
4. **Editar esfuerzo:** abrir registro activo, mostrar datos guardados; permitir cambiar miembro/día/actividad/horas/historia. Al cambiar historia conservar la asociación anterior si sigue válida (incluida ausencia de Sprint); solo si es incompatible exigir elegir otro Sprint válido o quitarla explícitamente. No elegir sustituto ni borrar asociación automáticamente. Confirmar con validación integral contra datos vigentes y mismo UUID. Refrescar totales de origen y destino. D-30 bloquea toda edición de esfuerzo en Finalizado con fin NULL, aunque no cambie la fecha; no bloquea la edición de estimaciones de horas de la historia.
5. **Eliminar lógicamente:** identificar el registro y pedir confirmación; cancelar/cerrar no envía escritura ni guarda borradores. Confirmar marca `deleted_at` sin cambiar sus datos de negocio; refrescar total. Repetición confirmada no duplica efecto ni cambia la primera marca; editar un eliminado no lo reactiva.

## Constraints

### PT-01 — Precisión de horas y fechas (aprobada para implementar)

Representar horas en centésimas enteras en dominio y SQLite (`int64`/`INTEGER`), no minutos ni `REAL`/`float64`. 0.01 horas equivale a una centésima de hora, no a un minuto. Decodificar el número decimal JSON exactamente, rechazar fracciones con precisión mayor a dos sin redondear ni truncar. Propuesta de sintaxis de escritura: números sin exponente, enteros o fracción de uno/dos dígitos (1, 1.5, 1.50); `1.500`, `1e1`, strings y separador coma no son el formato contratado. La UI puede aceptar entrada decimal local y serializar el número canónico, sin conversiones binarias que alteren centésimas. Esta sintaxis es una precisión técnica pendiente, no una decisión funcional adicional.

Estimada mínima 0.01; sin máximo funcional aprobado (24 solo aplica al real diario). Rechazar desbordamientos del almacenamiento antes de escribir, identificándolos como límite técnico de representación para revisar, sin wraparound. Real individual hasta 24 por consecuencia de D-24; 24 exactas se admite si no existe otro activo de ese miembro/día/proyecto. Sumas exactas y comprobadas contra overflow; salida como número decimal, no prometer ceros finales. No almacenar total como campo editable ni derivarlo de snapshots duplicados.

Reutilizar `project.Date`/`ParseDate`: `YYYY-MM-DD`, fechas existentes, años 0001…9999; no aceptar timestamps, año cero, fechas imposibles ni espacios. Comparar días, sin componentes horarios. Capturar hoy del reloj del servidor en la zona aprobada al validar la operación dentro de la unidad de trabajo; usar el mismo día para todas sus comprobaciones. El valor propuesto en formulario no evita revalidación en servidor si pasa la medianoche.

### PT-02 — Contratos REST aprobados para implementar

Rutas sobre `/api/projects`, nombres `snake_case` como US-05 y error `{"error":"mensaje"}`. No renombrar `projectID` de las respuestas existentes de Sprint. Todos los objetos de esfuerzo devuelven `id`, `project_id`, `backlog_item_id`, `member_id`, `work_date`, `activity`, `hours`, `sprint_id` (número o `null`); `deleted_at` es técnico, no se publica en consultas normales.

| Operación | Ruta y cuerpo | Éxito |
| --- | --- | --- |
| Crear historia, existente a ampliar | `POST /api/projects/{projectID}/backlog`: campos actuales + `estimated_hours` obligatorio | 201, historia con estimación y `actual_hours: 0` |
| Editar historia, existente a ampliar | `PUT /api/projects/{projectID}/backlog/{itemID}`: conjunto completo editable actual + `estimated_hours` obligatorio | 200, historia vigente, total activo |
| Listar Backlog, existente a ampliar | `GET /api/projects/{projectID}/backlog` | 200, arreglo actual enriquecido con estimación/total |
| Detalle de historia, nuevo | `GET /api/projects/{projectID}/backlog/{itemID}` | 200, historia con estimación/total |
| Alta esfuerzo, nuevo | `POST /api/projects/{projectID}/backlog/{itemID}/effort`: `member_id`, `work_date`, `activity`, `hours`, opcional `sprint_id` | 201, entrada creada |
| Listar esfuerzo de historia, nuevo | `GET /api/projects/{projectID}/backlog/{itemID}/effort` | 200, arreglo de activos; `[]` si no hay |
| Consultar entrada activa, nuevo | `GET /api/projects/{projectID}/effort/{entryID}` | 200, entrada activa |
| Editar esfuerzo, nuevo | `PUT /api/projects/{projectID}/effort/{entryID}`: `member_id`, `work_date`, `activity`, `hours`, `backlog_item_id`; `sprint_id` opcional con semántica de conservación | 200, entrada actualizada con mismo ID |
| Baja lógica, nuevo | `POST /api/projects/{projectID}/effort/{entryID}/delete` con `{"confirmed":true}` | 204 sin cuerpo |
| Contexto de formulario, nuevo | `GET /api/projects/{projectID}/backlog/{itemID}/effort-context` | 200, `{"today":"2026-10-26","sprints":[{"id":1,"closed":true}]}` con asignaciones de esa historia |

El último endpoint evita calcular hoy con el navegador y evita asumir un listado general de Sprints o pertenencias que no existe. Devuelve solo metadatos de las asignaciones existentes (abiertas/cerradas); no crea asignaciones ni una pantalla de gestión de Sprints. Sin asignaciones: `sprints: []`. Reutilizar `GET /members` para el selector. Al elegir otra historia, obtener su contexto y comprobar si la asociación persistida sigue válida; únicamente si es incompatible mantener la resolución pendiente hasta selección explícita. El proyecto consultado por la API existente aporta `status` y `end_date` para presentar D-30; no inventar fechas en el contexto.

Ejemplo válido de alta: `{"member_id":"b1d168da-325a-40d1-90b7-b5a1ba006389","work_date":"2026-10-26","activity":"Revisar criterios","hours":1.25,"sprint_id":null}`. La ruta aporta historia/proyecto. La respuesta añade UUID y ámbito. PUT requiere los campos editables salvo `sprint_id`: omitirlo conserva el valor persistido bajo validación contra la historia destino; `null` quita explícitamente y un número selecciona ese Sprint. Si se cambia historia y el Sprint conservado no es válido, rechazar con 409 sin reasignar ni quitar silenciosamente. La omisión sobre asociación nula la conserva sin exigir nueva selección. No se cambia `project_id` en PUT; mandar ese campo se rechaza, no se ignora.

Propuesta de decodificación para rutas nuevas: un único objeto JSON UTF-8, claves exactas y no duplicadas, EOF verificado; rechazar campos desconocidos/protegidos, raíz nula, tipos inválidos y campos obligatorios nulos. En POST de esfuerzo `sprint_id` ausente/nulo significa sin Sprint. No extender silenciosamente esta política estricta al decoder compartido de US-05, que hoy es permisivo con extras: ampliar únicamente su campo de estimación y validación, sin cambiar otras reglas de transporte. En US-05 ausente/null sigue siendo inválido para estimación.

Baja por acción POST sigue la convención `/finalize` y `/close`, permite confirmación explícita sin depender de cuerpo DELETE y evita exponer eliminación física. Confirmación API expresa intención, no demuestra que la UI mostró diálogo. Repetir baja confirmada de una fila ya eliminada del mismo proyecto devuelve 204 conservando la marca inicial; inexistente o de otro proyecto: 404. GET/PUT de eliminados: 404. No se agrega `include_deleted`, búsqueda global, paginación ni restauración. Orden propuesto de entradas: fecha ascendente e ID ascendente para desempatar; no confundirlo con orden de carga.

### PT-03 — Modelo de dominio y persistencia (aprobada para implementar)

| Entidad / estructura | Cambio futuro y responsabilidad |
| --- | --- |
| `backlog.Input` / `backlog.Item` (`internal/domain/backlog/item.go`) | Entrada estimación requerida en nuevas escrituras; entidad estimación opcional solo para hidratar legado. Tipo exacto de horas reutilizable por esfuerzo. Mantener estado y SP bajo sus comandos actuales. |
| Nueva `EffortEntry`, dominio propuesto `internal/domain/effort` | UUID, proyecto inmutable, historia, miembro, día, actividad, centésimas, Sprint opcional y marca técnica de baja. Invariantes y edición sin dependencia de HTTP/SQLite. |
| Nuevo caso de uso `internal/application/effort` | Coordinar alta/consulta/edición/baja, reloj y puertos transaccionales; verificar relaciones, fechas y cupo sobre datos vigentes. No importar SQLite. |
| `sprint.PlannedBacklogItem` / `Snapshot` | Agregar estimación opcional y distinguir dato vigente de los demás datos planificados inmutables; serializar `estimated_hours`. |
| Proyectos y miembros | Reutilizar entidades y UUID existentes. Incorporar consulta de límites temporales y validación transaccional de cambio de inicio; no inventar miembro global ni fechas de Sprint. |
| Lecturas de Backlog | Proyección `actual_hours` calculada sobre activos; no input editable ni fuente de verdad del agregado. |

Esquema conceptual, **sin SQL ni migración en esta etapa**:

- `backlog_items.estimated_hours_hundredths INTEGER NULL`, CHECK nulo o positivo. NULL reservado a datos históricos; no DEFAULT cero ni estimación inferida. La escritura de aplicación distingue reconstrucción histórica de alta/edición para exigir D-07/08.
- `sprint_backlog_items.estimated_hours_hundredths INTEGER NULL`, mismo CHECK. No cambia IDs, claves únicas, `status`, `completion_verified` ni SP. Lecturas usan la columna persistida, sin reconstruir el resto del snapshot desde Backlog.
- Nueva tabla propuesta `effort_entries`: `id TEXT PRIMARY KEY NOT NULL`; `project_id`, `backlog_item_id`, `member_id`, `work_date`, `activity` NOT NULL; `hours_hundredths INTEGER NOT NULL` positivo; `sprint_id INTEGER NULL`; `deleted_at TEXT NULL` (instante UTC generado al eliminar, no día de trabajo). No hace falta historial completo de ediciones para cumplir lo aprobado. Datos de negocio y primera marca de baja permanecen almacenados.
- FK a proyectos, miembros, historias y, si se informa, asignación `(sprint_id, backlog_item_id)` que ya es única en `sprint_backlog_items`. Validar además que Sprint/miembro/historia comparten `project_id`; las FK simples por ID no garantizan ese aislamiento. Se propone reforzarlo con índices UNIQUE `(project_id,id)` en padres y FK compuestas de historia/miembro/Sprint. Con Sprint nulo, las FK obligatorias de historia/miembro siguen vigentes. No borrado en cascada de esfuerzo ni snapshots.
- Índices propuestos por `(project_id,backlog_item_id,deleted_at)` y `(project_id,member_id,work_date,deleted_at)` para detalle/suma; no UNIQUE miembro/día/historia (D-16). CHECK y FK respaldan, pero no sustituyen validación Go ni cupo transaccional. `foreign_keys` está habilitado por conexión en el DSN actual.

### PT-04 — Integridad, atomicidad y concurrencia (aprobada para implementar)

1. Reutilizar el patrón `withImmediate` de `internal/persistence/sqlite/immediate.go`: conexión reservada, `BEGIN IMMEDIATE`, lecturas de estado vigente, validación Go a través de puerto, escrituras específicas y COMMIT/ROLLBACK. Espera acotada de 5 segundos y cancelación; BUSY/LOCKED solo se reintenta tras rollback confirmado, sin repetir a ciegas resultados de commit inciertos. No confiar exclusivamente en un mutex de proceso ni validar fuera de la transacción.
2. Alta/edición de esfuerzo: releer proyecto y bloquear por D-30 si está Finalizado con fin NULL, antes de escribir cualquier campo. Comprobar historia/miembro y asociación opcional, fechas y suma activa miembro/día/proyecto dentro de la misma unidad. En edición resolver `sprint_id` omitido desde el registro vigente y validar su compatibilidad con historia destino; si falla exigir elección explícita, nunca cambiar a otro Sprint o NULL por cuenta del servidor. Calcular `suma destino sin entryID + nuevas horas`; excluir el ID una sola vez si permanece en el grupo. Si cambia miembro o fecha, liberar su contribución de origen y verificar destino antes de escribir; cambiar solo historia/Sprint no crea otro cupo. Baja libera cupo activo y no exige fin conocido ni revalidar la fecha histórica. Nunca eliminar primero ni dejar movimientos parciales al fallar validación.
3. D-18/D-31 protegen también el PATCH de US-02: validar nuevo inicio contra el mínimo de fechas de **registros activos (`deleted_at IS NULL`)** del proyecto dentro de la misma transacción que edita proyecto. Si no hay activos no existe límite derivado de esfuerzo; siguen vigentes las reglas inicio/fin de US-02. Los eliminados no bloquean el cambio aunque sus fechas queden anteriores al nuevo inicio y conservan íntegramente sus datos y marca de baja. Inicio nulo elimina el límite; omisión conserva. Al rechazar PATCH con nombre e inicio, tampoco guardar el nombre. El filtro activo y la atomicidad son decisiones aprobadas por D-31; el mecanismo SQLite es propuesta técnica.
4. Serializar cambio de inicio con altas/ediciones/bajas de esfuerzo: quien confirma segundo revalida los activos vigentes; nunca admitir esfuerzo activo anterior al inicio persistido. Si una baja confirma primero, su fila deja de bloquear el nuevo inicio; si el cambio de inicio se valida antes, puede rechazarse y la baja posterior no lo reintenta automáticamente. Serializar finalización con esfuerzo: antes de finalizar validar fecha máxima activa contra el fin que se registraría y conservar reglas US-02. Con reloj estable y D-09 esto ya se cumple; no alterar fin automáticamente ni denegar esfuerzo solo por Finalizado. Si el Finalizado conserva fin NULL, D-30 bloquea altas y toda edición, permite consulta/baja y no habilita completar fin mediante estos comandos ni reabrirlo.
5. Modificar estimación en US-05 y propagarla a **todas** las filas de asignación de esa historia/proyecto en una sola transacción, aun cerradas. Nuevas asignaciones y traslado de US-10 deben leer/copiar la estimación vigente bajo la misma coordinación. Si destino ya existe, no duplicarlo y conservar su estimación sincronizada. La asignación actual recibe snapshot armado antes de su transacción diferida; no reutilizar esa estimación obsoleta: leer solo ese dato vigente dentro de la unidad. No volver a copiar título/puntos/criterios para aplicar la excepción.
6. El actual `backlog.Service.Update` lee antes de escribir y `BacklogRepository.Update` reescribe también `status` y `story_points`: necesita puerto de mutación transaccional y UPDATE de campos editables para no sobrescribir un acuerdo US-15 o finalización US-09 concurrente. Revalidar bloqueo de cambio provisional si se registra definitivo concurrentemente. US-15 sigue actualizando solo `story_points`, sin tocar horas.
7. Dos ediciones del mismo registro activo: releer en transacción y validar el reemplazo de campos requeridos y la conservación/elección de Sprint según PT-02; prevalece última escritura serializada válida. No se proponen ETag ni versionado. Baja contra edición: si baja confirma primero, PUT falla 404; si edición primero, baja conserva esa versión y la marca eliminada. Ninguna resurrección. Dos altas independientes pueden crear dos UUID; no deduplicar por mismos datos. Si se pierde respuesta de POST, consultar detalle antes de repetir.
8. Totales sin multiplicar por asignaciones: agregar desde `effort_entries` por historia/proyecto con filtro activo, nunca sumar un JOIN que repita entradas por cada Sprint. Leer detalle de historia con estimación/total en una vista consistente (misma sentencia o transacción de lectura). Si UI combina listado de entradas y detalle en solicitudes separadas durante una escritura, volver a cargar tras éxito; no prometer snapshot distribuido entre peticiones.

### PT-05 — Migración histórica (aprobada para implementar)

En futura implementación crear la siguiente migración libre según el estado de `main` de ese momento (actualmente la última es `000011_planning_poker.sql`); no modificar migraciones aplicadas ni fijar una versión que pueda colisionar. Agregar columnas anulables y tabla/índices/FK, manteniendo el mecanismo versionado/embebido y su aplicación transaccional existente.

Historias y snapshots previos sin horas quedan NULL, visibles y pendientes; esfuerzo inicial vacío y total activo 0. No estimar por SP, copiar ceros ni fechas actuales, recrear IDs, resetear resultados o perder vínculos. Cuando una historia se regularice, propagar a todas sus asignaciones incluidas las antiguas. Si en un entorno posterior ya existiera estimación válida, conservarla y alinear únicamente ese campo de sus snapshots; no inferir ni modificar otros datos.

Conservar también proyectos históricos Finalizados con `end_date = NULL`; D-30 se aplica al operar, sin una migración que complete el fin o cambie estado. Las filas de esfuerzo eliminadas conservan fechas originales aunque el inicio cambie posteriormente según D-31; no aplicarles un backfill ni exigirles cumplir el nuevo límite de activos.

La verificación futura de migración debe partir de una base temporal de esquema anterior con proyecto, miembros homónimos, historias y snapshots abiertos/cerrados/trasladados; aplicar y reabrir; comparar datos/IDs/resultados; consultar legado, registrar esfuerzo sin regularizar y luego regularizar. Reaplicar no debe duplicar ni reinicializar; fallo revierte la migración. No inspeccionar ni alterar la base local de desarrollo para generar evidencia.

### PT-06 — Interfaz (aprobada para implementar)

Ampliar alta/edición actual de Backlog en `web/static/index.html` y `web/static/js/app.js`: horas estimadas obligatorias, mínimo 0.01, paso 0.01. No bloquear estimación por SP definitivo ni esfuerzo; conservar validaciones US-05. En cada tarjeta mostrar estimación o «Pendiente de regularización», total real y acción de detalle; desde ese detalle gestionar entradas activas, sin formulario general sin historia.

Selector de miembro con nombre y UUID completo como dato secundario, valor UUID; sin miembros, indicar que debe registrarse uno mediante US-03. Fecha propuesta proviene del contexto del servidor y sigue editable. Selector opcional de Sprint solo con pertenencias existentes, incluidos cerrados identificados como tales. Si se cambia la historia en edición, conservar y mostrar el Sprint anterior cuando sigue válido; conservar ausencia de Sprint si ya era nula. Si es incompatible, mostrarlo como pendiente de resolver sin borrar el dato ni elegir un reemplazo; exigir otro Sprint válido o «Sin Sprint» explícito antes de guardar. El servidor revalida aun si la UI lo considera compatible.

En Finalizado con fin NULL, mostrar aviso de fecha de fin desconocida y deshabilitar alta/edición de esfuerzo por D-30; mantener consultas, totales y baja lógica confirmada disponibles. No bloquear toda la pantalla ni la edición de estimación de la historia. No ofrecer completar fin, inventarlo o cambiar estado desde esfuerzo. Esta regla se comprueba también en API, incluidas ediciones de actividad o miembro sin cambio de fecha.

Confirmación de baja identifica registro, fecha, miembro y horas; cancelar/cerrar mantiene todo, restaura foco y no guarda borradores. API valida `confirmed:true` independientemente del diálogo. Tras alta/edición/baja confirmada, recargar los totales afectados; conservar borrador si falla y mostrar error. Evitar doble envío, descartar respuestas de proyecto/historia anteriores y renderizar actividad/nombres como texto. En un Sprint histórico mostrar horas estimadas como **vigentes**, aclarando que pueden actualizarse; el resto del plan y resultados es histórico. No agregar métricas ni cambiar la UI de historial a una pantalla de gestión de esfuerzo.

## Compatibilidad e inventario de cambios futuros

| Historia / contrato actual comprobado | Impacto estrictamente necesario |
| --- | --- |
| US-05/06: `Input`, `Item`, POST/PUT/GET Backlog, `backlog.Service`, `Repository`, SQL y formularios no tienen horas | Agregar campo requerido en escrituras y opcional legado en lectura; total calculado; detalle nuevo; puerto de mutación atómica. Clientes y fixtures que crean/editan sin estimación dejarán de ser válidos: actualizar en implementación, no mediante valores inventados. |
| US-02: PATCH inicio opcional/nulo, `projects.Repository.Mutate` solo recibe proyecto y valida inicio/fin | Ampliar puerto/contexto transaccional para consultar esfuerzo activo y aplicar D-18/D-31, preservando semántica de omisión/null, nombre, estados y fin inmutable. Reloj/zona reutilizados. D-30 conserva Finalizado con fin NULL y bloquea únicamente alta/edición de esfuerzo, sin modificar sus fechas. |
| US-03: GET `/members`, UUID/proyecto/nombre; homónimos e identidades independientes | Reutilizar selector/lista; nuevo puerto de búsqueda miembro por ID/proyecto para validar en transacción. No cambiar alta/listado ni agregar edición/baja de miembros. |
| US-08: `PlannedBacklogItem`, `Snapshot`, tabla de asignaciones y POST/GET `/sprints/{sprintID}/backlog-items` | Agregar columna/JSON de estimación y única excepción D-21; coordinar asignación para evitar valor obsoleto. Sin endpoint genérico para editar snapshots ni cambios a los demás datos. |
| US-09: `CompleteBacklogItem` cambia Backlog y `completion_verified` atómicamente | Conservar transición y verificación; no exigir regularizar en este comando ni mutar esfuerzo. Edición US-05 concurrente no puede revertir su resultado. |
| US-10: `SprintRepository.Close` copia snapshots con `INSERT OR IGNORE`; conserva origen y no duplica destino | Incluir estimación en copia/coordinación; no copiar ni mover esfuerzo con traslado. Registro mantiene Sprint elegido (origen cerrado, destino válido o ninguno); cada entrada suma una sola vez por historia. Cierre no cambia sus horas. |
| US-11: GET `/sprints/closed`, detalle y planificación usan snapshots, UI `sprint-history.js` | Exponer estimación sincronizada, distinguirla del resto histórico; no reconstruir plan desde Backlog. D-21 es excepción explícita a su conservación general; no modifica resultado ni pertenencias. |
| US-15: `RegisterAgreement` actualiza únicamente `backlog_items.story_points` bajo `withImmediate` | Mantener independencia; estimación de horas editable aun con definitivo. Reemplazo genérico obsoleto no debe borrar acuerdo. No convertir puntos en horas ni actualizar puntos de snapshots por D-21. |
| Composición `cmd/carva`, interfaces `httpapi` y adaptadores SQLite | Incorporar servicio/puertos de esfuerzo, rutas/DTO/errores, reloj y consultas; no framework ni dependencia nueva necesaria. |

Los SDD/BDD previos permanecen sin modificaciones en esta entrega. Su futura sincronización está identificada: US-05 debe incluir horas requeridas en datos y ejemplos válidos (también fixtures de US-02 que crean historias); US-08 BR-03 y US-11 reglas de historial deben explicitar solo la excepción D-21. Los escenarios actuales que cambian título/criterios/puntos continúan válidos al aportar estimación donde corresponda. US-10 continúa prohibiendo nueva planificación/compleción en cerrado; administrar esfuerzo por US-16 no es esa operación. Ninguna funcionalidad pendiente se da por existente.

## Edge Cases

- Estimación 0.01, 1, 1.50 y mayor que 24; rechazo de ausencia, cero, negativos, nulo o precisión excedida sin redondeo. Esfuerzo 0.01 y cupo exacto 24, rechazo a 24.01 acumuladas.
- Legado consultable/recibe esfuerzo sin estimación si no está bloqueado por D-30; editar título sin regularizar se rechaza entero, regularización propaga y no inventa estimación original.
- Historia sin Sprint, en varios Sprints o trasladada; elegir cualquiera de sus asignaciones históricas válidas o ninguno, nunca distribuir/deduplicar horas por snapshot. D-22 no exige que día de trabajo caiga en fechas de Sprint inexistentes.
- Real > estimado, múltiples registros y eliminación de uno; total cambia sin tocar estimación.
- Fecha igual a inicio, hoy o fin permitido; futura, anterior a inicio, posterior a fin rechazadas. Proyecto sin inicio solo aplica límite superior; inicio futuro deja sin fecha válida para esfuerzo hasta corregirlo, sin borrar sus datos.
- Finalizado sin fin: bloquear alta y toda edición de esfuerzo, aun si conserva fecha; permitir consultar/eliminar lógicamente y conservar fin NULL/estado. Inicio contra activos únicamente: eliminados pueden conservar fecha anterior al nuevo inicio sin impedir el cambio; nunca corregir su historia.
- Cambio de historia con Sprint anterior válido, incompatible o nulo: conservar el válido/nulo si se omite el campo; incompatible exige resolución explícita, incluso si está cerrado, sin borrado o reemplazo automático.
- Miembros homónimos seleccionados por UUID; ausencia de miembros; historias/miembros/Sprint/registro inexistentes o de otro proyecto; registro eliminado no se puede consultar/editar normalmente.
- Contención, fallos antes de commit, respuesta perdida y operaciones concurrentes según PT-04. Sin promesa de entrega exactamente una vez.

## Error Conditions

Códigos/mensajes siguientes son contrato técnico propuesto. Estructura JSON de error existente; UI presenta texto claro en español sin filtrar SQL, DSN o rutas.

| Condición | HTTP propuesto | Efecto |
| --- | --- | --- |
| JSON/tipos/formato decimal/fecha inválidos, campos requeridos ausentes/nulos, actividad blanca, horas ≤ 0, precisión excedida; `sprint_id` no entero positivo ni null; confirmación distinta de true | 400 | Sin cambios; mensaje específico de campo. Omitir Sprint en PUT tiene la semántica de conservación descrita. |
| Proyecto/historia/miembro/Sprint/entrada inexistente o fuera de ámbito; GET/PUT de eliminado | 404 | Sin crear recursos ni alterar otro proyecto. |
| Sprint existe en proyecto pero historia no está asignada a él | 409 | Conservar asociación y datos anteriores. |
| Cambio de historia con Sprint anterior incompatible y `sprint_id` omitido | 409 | Exigir otro Sprint válido o null explícito; conservar todo el registro. |
| Alta o cualquier edición de esfuerzo en Finalizado con fin NULL (D-30) | 409 | Informar fecha de fin desconocida; no inventarla ni cambiar estado, entrada o totales. Consultas y baja lógica permanecen permitidas. |
| Fecha futura, antes de inicio o después de fin; cambio de inicio incompatible; cupo diario excedido | 409 | Sin cambios parciales, fechas nunca corregidas automáticamente. |
| Cambiar proyecto por cuerpo PUT | 400 | Campo protegido; explicar corrección por baja y nueva alta. |
| Contención de escritura agotada | 503 | Sin éxito aparente; reintento con nueva validación. |
| Fallo inesperado, overflow al calcular agregado o commit incierto | 500 | No mostrar total incorrecto/éxito; consultar antes de repetir escritura si resultado incierto. |

Precedencia para alta/edición de esfuerzo: validar transporte/formatos, resolver proyecto e historia de ruta o entrada activa, comprobar D-30, miembro e historia destino y Sprint en ese ámbito, compatibilidad de asociación conservada/elegida, límites de fechas y cupo; escribir. En baja, validar confirmación y pertenencia del registro, sin exigir fin conocido ni validar de nuevo su fecha. Consultas tampoco exigen fin conocido. `estimated_hours` inválida en US-05 se integra en su validación de negocio después de resolver el proyecto, sin cambiar la precedencia actual de otros campos. Un error en cualquier etapa conserva íntegramente el estado anterior. No usar 409 por historia completada, Sprint cerrado o proyecto finalizado por sí solos: la excepción precisa es Finalizado con fin NULL en alta/edición.

## Acceptance Criteria

AC locales a esta SDD. Los tres AC originales de #16 quedan cubiertos por AC-01/12 (registrar/consultar estimadas), AC-04/06 (campos y atribución) y AC-09/17 (independencia). Los demás precisan D y robustez técnica propuesta.

| AC | Criterio verificable | D / fuente | BDD-US-16 |
| --- | --- | --- | --- |
| AC-01 | Crear/editar estimación positiva con hasta dos decimales, sin Sprint y aun con esfuerzo; conservar vigente. | D-01/02/06 | 01, 02, 03 |
| AC-02 | Rechazar nueva historia sin estimación válida y cualquier precisión excedida sin redondeo ni escritura parcial. | D-06/07 | 04 |
| AC-03 | Legado consultable y admite esfuerzo con indicador salvo bloqueo D-30; edición de historia exige regularizar y conserva identidad. | D-08/28/30 | 05, 06 |
| AC-04 | Alta real conserva miembro, fecha, actividad, horas, historia obligatoria e identidad, sin exigir Sprint. | D-03/04 | 07, 08 |
| AC-05 | Rechazar horas/actividad/fecha mal formadas o requeridos ausentes sin escritura. | D-05/10, PT-01/02 | 08, 09, 10, 40 |
| AC-06 | Elegir UUID desde lista, distinguir homónimos y validar ámbito de todos los recursos. | US-03, D-27 | 11, 12, 27 |
| AC-07 | Proponer hoy de zona aprobada, permitir pasado e igual inicio/hoy y rechazar fuera de límites; validar cambio de inicio atómicamente solo contra activos, conservando eliminados sin bloquearlo. | D-09/18/29/31 | 13, 14, 15, 16, 33, 40, 45 |
| AC-08 | CRUD lógico en Finalizado con fin conocido conserva estado/fechas y valida trabajo ≤ fin incluso cargado después; fin desconocido se rige por AC-19. | D-15/17/30 | 17, 18, 40, 41 |
| AC-09 | Acumular múltiples registros activos, máximo 24 por miembro/día/proyecto; superar estimación no la cambia. | D-16/24/25 | 19, 20, 21, 22, 23, 32, 42 |
| AC-10 | Editar campos permitidos conserva ID/proyecto, recalcula origen/destino y valida cupo; conserva Sprint compatible y exige decisión explícita si es incompatible. | D-11/26/27 | 22, 23, 26, 27, 40, 42, 46 |
| AC-11 | Confirmar baja lógica, cancelar sin escritura; eliminados almacenados, invisibles y fuera de sumas; no resucitan. | D-12/13 | 24, 25, 31, 35 |
| AC-12 | Cada historia presenta estimación/total/detalle; administración se realiza desde ese detalle y errores no se confunden con vacío. | D-19/20/29 | 05, 07, 13, 24, 30, 39 |
| AC-13 | Alta/edición/baja de esfuerzo en Completada y Finalizado con fin conocido conserva estados, fechas y resultados; aplicar AC-19 si falta fin. | D-14/15/30 | 17 |
| AC-14 | Sprint opcional de la historia/proyecto; cambio de historia conserva asociación válida o nula y exige elegir otro válido o quitar explícitamente si es incompatible. | D-22/26 | 07, 12, 26, 28, 46 |
| AC-15 | Esfuerzo en cerrados y trasladadas no reabre, mueve ni duplica registros/plan/resultados. | D-23, US-10/11 | 28, 29, 34 |
| AC-16 | Cambiar estimación propaga a todos los snapshots, aun cerrados/trasladados; otros datos inmutables y atomicidad. | D-21, US-08/10/11 | 06, 29, 34, 36 |
| AC-17 | Horas y Story Points independientes; edición concurrente no revierte acuerdo ni compleción. | #16, US-15, PT-04 | 02, 19, 37 |
| AC-18 | Migración, reapertura, rollback y concurrencia conservan invariantes y cambio de inicio atómico contra activos, sin éxito falso ni duplicados derivados de traslados. | D-13/18/21/24/31, PT-03/04/05 | 30, 31, 32, 33, 34, 35, 36, 37, 38, 41, 42, 45 |
| AC-19 | En Finalizado con fin NULL rechazar toda alta/edición de esfuerzo y permitir consultas/baja lógica, conservando datos, estado y fin desconocido. | D-30 | 43, 44 |

## BDD / Test Traceability

Cada escenario tiene `@BDD-US-16-NN`, `@AC-NN` y D pertinentes; `@PT-*` distingue precisiones técnicas propuestas. Para cada BDD-US-16-NN se prevé la familia **TEST-US-16-NN** (identificador de planificación, no test existente). En la fase autorizada se registrarán nombres reales de tests, evidencia RED → GREEN → REFACTOR y código correspondiente; hoy la cadena se detiene en BDD.

| Familias futuras | Nivel y comprobación prevista |
| --- | --- |
| TEST-US-16-01…10, 13…16, 18…23 | Unitarios Go de horas exactas, fechas, validación, cupo y regularización; aplicación con reloj controlado y puertos dobles. |
| TEST-US-16-07…12, 17, 22…29, 40, 43, 44, 46 | Integración SQLite/HTTP de contratos, UUID, ámbito, conservación/elección de Sprint, reemplazos, confirmación y baja lógica; D-30 bloquea toda edición sin fin, pero no consulta/baja. Comprobar filas almacenadas sin exponer endpoint de eliminados. |
| TEST-US-16-30…38, 41, 42, 45 | Migración y reapertura sobre archivo temporal; fallos inyectados antes de commit, dos conexiones/repositorios independientes, coordinación con US-02/08/09/10/15 y baja contra cambio de inicio solo sobre activos; comprobar valores y filas finales, no solo status HTTP. |
| TEST-US-16-05/07/11/13/24/25/26/39/43/44/46 | UI con dobles HTTP y reloj conocido: selectores, estado pendiente, fecha predeterminada, conservación de Sprint compatible/resolución incompatible, bloqueo D-30, diálogo, totales, borradores y respuestas tardías. |

Implementación y automatización autorizadas. La evidencia real de ciclos RED → GREEN → REFACTOR, verificaciones y correspondencia BDD/tests se registra en docs/evidencia-ia-us-16.md durante el desarrollo. No se requiere framework BDD adicional.

## Riesgos, incompatibilidades y revisión humana

**Cuestiones funcionales resueltas, fuera de los pendientes:** R-01 queda resuelto por D-30 (alta/edición bloqueadas sin fin; consulta/baja permitidas); R-02 por D-31 (solo activos para cambio de inicio atómico). D-26 corregida conserva Sprint compatible y exige resolución explícita solo ante incompatibilidad. Su comportamiento ya está incorporado en BR, AC, BDD y propuestas técnicas; no se solicita una nueva decisión sobre esos puntos.

- **R-03 — Compatibilidad de contratos:** exigir estimación rompe los clientes/fixtures antiguos de alta/PUT sin ese campo. Es el impacto aprobado D-07/08; no simular compatibilidad llenándolo con valores ficticios. Desplegar contrato y formularios juntos en implementación; datos históricos sobreviven.
- **R-04 — Inmutabilidad documental previa:** D-21 establece una excepción aprobada que debe explicitarse en US-08/11 al implementar. Solo horas estimadas se propagan; no Story Points ni resultados. Una comparación histórica futura deberá reconocer que esa estimación es vigente, sin decidir ahora fórmulas de métricas.
- **R-05 — Carreras existentes:** asignación con snapshot previo y UPDATE genérico del Backlog necesitan coordinación descrita en PT-04. Limitar la futura modificación a proteger estimación/estado/acuerdo; no sustituir la arquitectura ni agregar versionado innecesario.
- **R-06 — Revisión técnica pendiente:** aprobar sintaxis decimal/representación y límites técnicos de PT-01, rutas/JSON/errores/contexto y acción de baja PT-02, FK/esquema/migración PT-03/05, mecanismo de concurrencia PT-04 e interfaz PT-06. Son propuestas técnicas para cumplir las 31 decisiones funcionales ya aprobadas; el alcance de D-30/31 y la corrección de D-26 no siguen pendientes.

La aprobación de esta documentación habilitaría únicamente la etapa que el responsable autorice expresamente después. Este PR no hace merge, no cierra #16; aceptación final y revisión de la Definition of Done pendientes del responsable.
