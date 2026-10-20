# SDD-US-02 — Modificar proyecto y registrar fechas

**Estado:** SDD/BDD y PT-01 a PT-05 aprobados para implementar por el responsable en la solicitud del 2026-10-20, con ajustes previos de BDD-07/18. Revisión final de la implementación por el equipo pendiente.
**Aprobación de esta etapa:** Las secciones que conservan la denominación «propuesta técnica en revisión» documentan su origen; PT-01 a PT-05 quedan autorizadas por esta solicitud. La evidencia posterior se registra en `docs/evidencia-ia-us-02.md`, sin atribuir pruebas a la etapa de especificación. BDD-07 no incluye En progreso sin Sprints y BDD-18 comprueba conservación de datos actuales sin depender de US-11.
**Fecha:** 2026-10-20.
**User Story:** Como integrante del equipo, quiero modificar un proyecto y registrar sus fechas de inicio y finalización para mantener sus datos.
**Trazabilidad:** [enunciado §2.1](../docs/enunciado.md) → EP-01 → [US-02 / Issue #2](https://github.com/Aguscampos04/Carva/issues/2) → [decisiones D-01…D-10](../docs/contexto-proyecto.md#modificación-y-finalización-del-proyecto--us-02--issue-2) → esta SDD → AC → [BDD-US-02](../bdd/BDD-US-02-modificar-proyecto.feature).
**Fuentes:** AGENTS.md; skills project-manager, sdd, bdd, github y golang en `.agents/skills/` (ubicación efectiva, no `skills/`); enunciado, contexto, reglas, arquitectura; Issue #2 abierta, sin comentarios al consultar; SDD/BDD US-01/07, SDD US-03/05 y código actual de proyectos, Sprints, miembros, Backlog e interfaz.
**Base revisada:** `main` y `origin/main` en `37c3752218c68e82a48691d8de5e4b180909f5b4` después de `git fetch origin main`; sin cambios locales iniciales. Rama documental: `codex/us-02-sdd-bdd`. No se integró ni sobrescribió trabajo de terceros.

## Objective

Permitir editar nombre y fecha de inicio de un proyecto existente y finalizarlo manualmente, registrando su fecha real de finalización. La consulta debe reflejar datos persistidos y el mismo estado compartido con US-01/07 y las consultas futuras US-04/23.

Los dos AC originales de #2 se conservan en AC-01 (modificación y consulta posterior) y AC-04 (consulta del estado Finalizado). Los demás precisan las decisiones aprobadas y el contrato propuesto. No se incluyen eliminación, reapertura, fechas de Sprints, permisos ni implementación de otras historias.

## Inputs

| Entrada | Origen / tipo propuesto | Regla |
| --- | --- | --- |
| Identificador | Ruta; string UUID emitido por US-01 | Proyecto existente. No se cambia el ID. Como GET actual, identificador que no resuelve devuelve no encontrado, sin imponer un nuevo validador global. |
| `name` | Edición; string opcional en PATCH | Si está presente debe ser no blanco según `strings.TrimSpace`. Conservar el texto original, como hace `project.New` de US-01: no recortar al guardar, sin unicidad ni máximo funcional nuevos. `null` no es nombre válido. |
| `start_date` | Edición; string de fecha o `null`, opcional | Fecha manual, opcional y modificable (D-02). Omisión conserva; `null` elimina la fecha; string válido establece/reemplaza. Cadena vacía no equivale a ausencia. |
| `confirmed` | Finalización; booleano requerido | Solo `true` solicita la finalización confirmada. `false` no produce cambios; el flujo UI de cancelación no envía la solicitud. |
| Fecha actual | Reloj del servidor, no del cliente | Insumo del caso de uso de finalización; precisión PT-01. |

El cliente no escribe `id`, `status` ni `end_date`. Edición y finalización son acciones separadas: finalizar no guarda silenciosamente un formulario de edición pendiente.

## Expected Outputs

- Edición válida: mismo UUID, nombre/fecha de inicio actualizados y estado/fecha de finalización conservados; consulta posterior y reinicio sobre la misma base devuelven esos valores.
- Finalización confirmada válida: estado `Finalizado` y fecha de finalización calculada por el servidor, guardados juntos. No se requiere fecha de inicio ni Sprint.
- Cancelación: ningún cambio, sin guardar ediciones pendientes ni tocar Sprints, miembros o Backlog.
- Error: respuesta controlada, sin éxito aparente ni cambios parciales. Una pérdida de respuesta después de commit deja resultado incierto para el cliente: consultar de nuevo antes de decidir reintentar.
- Consulta REST propuesta, conservando los nombres existentes `id`, `name`, `status` y añadiendo fechas siempre presentes (ausencia como `null`):

```json
{
  "id": "0828dc84-050f-421b-b842-420f63851a2b",
  "name": "Carva Demo",
  "status": "Finalizado",
  "start_date": "2026-10-01",
  "end_date": "2026-10-20"
}
```

## Business Rules

Los identificadores D de esta sección corresponden exclusivamente a US-02.

| ID | Regla verificable | Fuente |
| --- | --- | --- |
| BR-01 | Solo una acción manual confirmada finaliza el proyecto. Alcanzar/pasar una fecha o consultar no cambia su estado. | D-01/07 |
| BR-02 | Fecha de inicio manual, opcional y modificable; no se infiere de alta, primer Sprint ni finalización. | D-02 |
| BR-03 | Finalizar registra automáticamente la fecha actual y el estado Finalizado en una única operación. | D-03 |
| BR-04 | Después de finalizar solo nombre e inicio siguen editables en este recurso. No se permite reabrir ni corregir fin mediante edición. | D-04 |
| BR-05 | Se finaliza desde Planificado o En progreso, con cero o más Sprints. | D-05 |
| BR-06 | Cuando existen ambas fechas, inicio ≤ fin. Se permite igualdad y fin sin inicio. Se valida también al editar un finalizado. | D-06 |
| BR-07 | Cancelar la confirmación no produce cambios. | D-07 |
| BR-08 | Se rechaza crear un Sprint nuevo si el proyecto ya está Finalizado; sus Sprints siguen consultables. | D-08 |
| BR-09 | Finalizar no cierra, elimina, mueve ni altera Sprints, objetivos o snapshots existentes, aunque haya un Sprint abierto. | D-09 |
| BR-10 | Se conservan alta/listado de miembros y creación/edición/consulta de Product Backlog en Finalizado, con sus validaciones vigentes. | D-10 |
| BR-11 | Alta de proyecto → Planificado. Primer Sprint en Planificado → En progreso. Otro Sprint en En progreso conserva estado. Editar nombre/inicio no cambia estado. | US-01, US-07 |
| BR-12 | El nombre conserva la regla de US-01: debe contener algún carácter no blanco; identidad y pertenencias permanecen estables. | US-01, #2 |

No se introduce una prohibición de fecha de inicio futura para proyectos sin fecha de fin: D-06 compara contra el fin cuando exista. Si al finalizar el inicio es posterior al día actual, se rechaza la finalización por incompatibilidad con D-03/06; el usuario puede corregir/eliminar el inicio antes de volver a confirmar. La fecha de inicio por sí sola no inicia el proyecto.

## Constraints

### PT-01 — Fechas y reloj (propuesta técnica en revisión)

- Ambas fechas son días de calendario gregoriano, sin hora ni offset; string canónico `YYYY-MM-DD`, años `0001` a `9999`, mes/día reales y validación de años bisiestos. No aceptar espacios, formatos locales, timestamps ni normalizar fechas imposibles.
- Ausencia: JSON `null`, SQL `NULL`, valor opcional de fecha en Go. No usar cadena vacía, fecha cero ni el instante Unix como sustitutos.
- Proponer la zona IANA **`America/Argentina/Buenos_Aires`** para interpretar «fecha actual», común a todos los proyectos en este alcance. Requiere aprobación técnica; no está incluida en D-03. No depende de zona del navegador ni zona local del host ni de `CURRENT_DATE` de SQLite (UTC).
- Aplicación obtiene un instante de un reloj inyectable del servidor, lo convierte a esa zona y extrae el día. Se toma después de adquirir el acceso de escritura de la operación exitosa y antes de validar/mutar; esa es la fecha del acto de finalizar, aunque el commit termine después de medianoche. Un intento abortado sin commit puede tomar un nuevo instante al reintentarse.
- El dominio recibe el día ya determinado: sin acceso a reloj del sistema o configuración de infraestructura. Composición asegura disponibilidad de la zona, incluso en despliegues sin base de zonas (por ejemplo con `time/tzdata` estándar); un fallo de configuración no habilita fallback silencioso.
- Ejemplo determinista: instante `2026-10-21T02:30:00Z` → fin `2026-10-20` en la zona propuesta. `2026-10-21T03:00:00Z` → `2026-10-21`. La UI muestra días sin convertirlos a instantes JavaScript que puedan desplazar la fecha.

### PT-02 — REST y semántica de edición (propuesta técnica en revisión)

Se conservan prefijo `/api/projects`, JSON y errores `{"error":"mensaje"}` con `Content-Type: application/json; charset=utf-8`. Nuevas fechas en `snake_case`, como Backlog/miembros; no renombrar `projectID` del Sprint existente.

| Operación | Contrato propuesto | Éxito |
| --- | --- | --- |
| Consultar | `GET /api/projects/{id}` existente | `200`, objeto ampliado con las dos fechas. |
| Editar parcialmente | `PATCH /api/projects/{id}`, por ejemplo `{"name":"Carva II","start_date":"2026-10-01"}` | `200`, objeto vigente completo. |
| Quitar inicio | Mismo PATCH, `{"start_date":null}` | `200`, inicio `null`, resto conservado. |
| Finalizar | `POST /api/projects/{id}/finalize`, `{"confirmed":true}` | `200`, objeto vigente completo; no crea otro proyecto ni devuelve Location. |
| Crear proyecto | `POST /api/projects`, sigue recibiendo solo nombre | `201` y Location vigentes; se propone añadir fechas `null` en la respuesta. No ampliar entradas de US-01. |
| Crear Sprint | `POST /api/projects/{projectID}/sprints`, cuerpo existente `{"goal":"..."}` | `201` para Planificado/En progreso; nuevo `409` para Finalizado. |
| Consultar Sprint | `GET /api/projects/{projectID}/sprints/{sprintID}` existente | `200` con el mismo contrato y pertenencia vigente, también en Finalizado. |

Se elige PATCH para expresar omisión/conservación y `null`/eliminación del inicio; no es un reemplazo PUT como la edición completa de Backlog. Es un objeto JSON parcial propio, no JSON Patch ni una adopción global de Merge Patch. Al menos un campo editable debe estar presente; `{}` se rechaza. Todas las propiedades proporcionadas se validan antes de guardar; un nombre válido con fecha inválida no cambia nada.

Proponer validación estricta solo para las nuevas escrituras: exactamente un objeto JSON y fin de entrada verificado con el mismo decoder; rechazar raíz `null`/array/escalar, UTF-8 inválido, claves duplicadas, desconocidas o con capitalización distinta; tipos inválidos y varios valores JSON. Campos protegidos se rechazan, aunque repitan su valor vigente; no se ignoran ni se aplican parcialmente. No endurecer los handlers existentes de otras historias como parte de US-02.

Finalizar no acepta nombre, inicio, estado ni fecha de fin. `confirmed` ausente, `null`, no booleano o `false` responde `400`, sin cambios. Esto expresa intención del cliente; no prueba que se mostró un diálogo, por lo que D-07 también requiere validación de UI. Métodos no expuestos conservan la respuesta del router (`405` para método no permitido).

### PT-03 — Repetición y concurrencia (propuesta técnica en revisión)

- Repetir un PATCH con los mismos valores devuelve `200` sin diferencias de datos. No hay alta, nueva identidad ni eventos adicionales. En ediciones concurrentes se aplican solo los campos presentes sobre el registro vigente: campos distintos se conservan; para un mismo campo prevalece la última escritura serializada válida. No se propone historial, ETag ni versionado en esta historia.
- Repetir finalización confirmada sobre un proyecto ya Finalizado devuelve `200` con su representación vigente y **conserva el primer fin**, aunque sea otro día. No vuelve a evaluar inicio contra el nuevo día ni produce efectos sobre Sprints. Ediciones válidas posteriores se reflejan en esa representación; no se promete respuesta byte a byte idéntica. Cuerpo inválido sigue siendo `400` incluso si ya finalizó.
- Dos finalizaciones simultáneas producen una sola transición efectiva y un único día de fin; la segunda observa el resultado de la primera. El servidor no confía en el estado mostrado al abrir la confirmación.
- Las tres escrituras (PATCH, finalizar, crear Sprint) deben releer el proyecto **dentro de la misma transacción que escribe**. Proponer `BEGIN IMMEDIATE` en una conexión SQLite reservada, con COMMIT/ROLLBACK y liberación garantizados. No combinar un BEGIN manual con un `sql.Tx` paralelo ni leer por otra conexión fuera de esa transacción.
- Tras obtener el bloqueo, la aplicación valida el estado/fechas vigentes con las reglas del dominio mediante un puerto transaccional; persistencia implementa el bloqueo y SQL. `ProjectRepository.Save` actual es INSERT y debe conservarse para alta; agregar operaciones de actualización, no sustituirlo por un upsert ni guardar una copia completa obsoleta.
- US-07 actualmente lee un proyecto en aplicación y vuelve a leer solo `status` en `SprintRepository.Create`; no rechaza Finalizado. Esa transacción deberá incorporar comprobación de elegibilidad y transición basada en el estado vigente antes de insertar. Un control exclusivo en UI o previo a la transacción no satisface D-08.
- Si finalizar confirma primero, crear Sprint se rechaza sin inserción. Si Sprint confirma primero, se crea/inicia el proyecto y luego finalizar puede finalizarlo, conservando ese Sprint. Nunca volver de Finalizado a En progreso ni crear un Sprint después de observar Finalizado dentro de la transacción.
- Edición y finalización también se serializan: si se guarda primero un inicio posterior al día actual, finalizar se rechaza; si finaliza primero, esa edición se rechaza al comparar con el fin persistido. Nunca confirmar ambos cambios dejando inicio > fin. Una edición de nombre no sobrescribe la transición ni la fecha de fin.
- Proponer espera de bloqueo acotada de hasta 5 segundos por operación, respetando cancelación del contexto; agotada la espera, `503` sin cambios de esa operación. Si el driver exige reintento por BUSY/LOCKED, repetir la transacción completa solo tras rollback confirmado y dentro del mismo límite. No reintentar a ciegas después de un commit de resultado incierto. Verificar esta coordinación con conexiones independientes en la etapa futura de integración.

### PT-04 — Dominio, capas y SQLite (propuesta técnica en revisión)

| Capa / ubicación | Cambio previsto, sin implementar |
| --- | --- |
| `internal/domain/project` | Añadir `StartDate` y `EndDate` opcionales con tipo fecha de calendario validado; operaciones de edición, finalización e invariantes; elegibilidad para crear Sprint; mantener estados y `StartFirstSprint`. Errores de fecha, rango y proyecto finalizado tipados, sin HTTP/SQL. |
| `internal/application/projects` | Coordinar Update/Finalize mediante puertos, reloj y transacción; aplicar PATCH al proyecto vigente; no sobrescribir campos omitidos. |
| `internal/application/sprints` | Evaluar elegibilidad y transición sobre el proyecto vigente en la unidad transaccional de creación; conservar objetivo/identificadores y atomicidad actuales. |
| `internal/persistence/sqlite` | Implementar nuevos puertos, lectura de fechas con NULL y actualizaciones específicas; transacción común de lectura-validación-escritura. No dependencia de SQLite desde aplicación/dominio. |
| `internal/presentation/httpapi`, `web/static` | Nuevas rutas/DTO, traducción de errores y formulario/confirmación; validación central sigue en Go. |
| `cmd/carva` | Componer reloj/zona y puertos con repositorios/servicios actuales; sin nueva tecnología o framework. |

Migración **aditiva propuesta**, no archivo ejecutable en esta etapa: agregar a `projects` columnas `start_date TEXT NULL` y `end_date TEXT NULL`, sin DEFAULT de fecha actual. Nueva versión siguiente a la última disponible al implementar (hoy `000007`; candidato `000008_projects_dates.sql`, sin reservarlo frente a trabajo ajeno). Mantener migraciones aplicadas intactas y registro versionado/embebido existente. No recrear IDs, tablas de miembros, Sprints, Backlog ni snapshots.

Los registros previos conservan UUID, nombre, estado y relaciones y reciben fechas NULL: no inferirlas de Sprints, archivos ni fecha de migración. El CHECK de estados actual ya admite Finalizado. Las invariantes de calendario, rango e inmutabilidad se validan en Go dentro de la unidad transaccional; no depender de funciones de fecha SQLite que normalicen entradas. No se propone reconstruir la tabla ni añadir triggers en esta etapa.

La ruta actual no permite finalizar, por lo que no hay evidencia de finales históricos con fecha faltante. Si se encuentran registros Finalizado previos por cargas externas, conservar fin NULL como dato desconocido (y mostrarlo como «Sin registrar»), nunca inventar una fecha o permitir que un reintento la reescriba; relevar esos datos antes de implementar. Para **nuevas** finalizaciones, fin no nulo es obligatorio. Esta compatibilidad no autoriza nuevas finalizaciones sin fecha.

### PT-05 — Interfaz (propuesta técnica en revisión)

- Ampliar el detalle actual con fechas (ausentes: «Sin registrar»), formulario de nombre/inicio y acción «Guardar cambios». Inicio como control de día, opcional, borrable; al guardar vacío enviar `null`. No alterar el formulario de alta de US-01.
- «Finalizar proyecto» solo disponible en Planificado/En progreso. Abrir confirmación identificando proyecto e indicando que fin/estado quedan fijos y que no se cierran Sprints. Confirmar envía `confirmed:true`; cancelar o cerrar el diálogo no envía escritura ni guarda el formulario. Restaurar foco y permitir teclado.
- Si hay edición sin guardar, finalizar opera sobre los datos **persistidos**, no sobre el borrador; indicarlo en el diálogo para evitar un guardado implícito. No combinar dos acciones silenciosamente.
- Tras éxito, mostrar representación del servidor. En Finalizado mantener edición de nombre/inicio; estado y fin de solo lectura; mantener secciones de miembros y Backlog operativas. Nunca un bloqueo genérico de toda la pantalla.
- No existe formulario de Sprints en el frontend actual: no crear una pantalla de Sprints dentro de US-02. El bloqueo se aplica en API/Go. Las futuras interfaces de US-07 deberán reflejarlo conservando consulta de existentes.
- Evitar doble envío mientras hay una solicitud; aun así el backend cubre repetición. Mostrar errores sin afirmar éxito ni borrar el borrador de edición; consultar de nuevo tras pérdida de respuesta. Ignorar respuestas de un proyecto anterior si cambia la selección, siguiendo el patrón de miembros. Renderizar nombres como texto.

## Edge Cases

- Inicio ausente, agregado, corregido o eliminado; editar solo nombre no borra el inicio. Se permite inicio igual a fin y año bisiesto válido.
- Alta y primer Sprint no generan fecha de inicio automáticamente. Llegar al día de inicio no inicia ni finaliza el proyecto.
- Finalizar desde ambos estados, con o sin inicio, con cero Sprints o con Sprints existentes. Un Sprint abierto no bloquea finalizar.
- Inicio futuro sin fin: editable; finalizar con ese inicio mayor al día actual falla sin cambios. Después de finalizar, inicio mayor que fin también falla.
- Doble clic, reintento al día siguiente y dos finalizaciones simultáneas: fecha original estable.
- Cambio de día UTC distinto al día de la zona propuesta; reloj del cliente adelantado no afecta el fin.
- Escrituras concurrentes según PT-03; operaciones sobre otro proyecto no lo alteran. Fallar cualquiera de las escrituras de finalización o creación de Sprint revierte toda su unidad.

## Error Conditions

Mensajes y códigos son parte del contrato técnico propuesto, no nuevas reglas funcionales.

| Condición | HTTP propuesto / mensaje orientativo | Resultado |
| --- | --- | --- |
| JSON inválido, campos desconocidos/protegidos, PATCH vacío o tipos inválidos | `400` / `invalid JSON request` | Ningún campo se modifica. |
| Nombre nulo, vacío o blanco | `400` / `project name must not be blank` | Sin edición parcial. |
| Fecha inválida o inicio > fin ya registrado | `400` / `invalid project date` | Sin edición parcial. |
| Confirmación ausente, falsa o inválida | `400` / `project finalization must be confirmed` | Sin finalización. |
| ID que no resuelve a proyecto | `404` / `project not found` | No crear proyecto ni modificar otro. |
| Inicio > fecha actual al finalizar | `409` / `start date is after finalization date` | Conserva estado previo y fin nulo. |
| Crear Sprint en Finalizado | `409` / `project is finalized` | No crear Sprint ni cambiar estado. |
| Contención de escritura que supera el límite | `503` / `project is busy; retry later` | Rollback; sin éxito aparente. |
| Fallo inesperado de lectura/escritura/commit | `500` / `could not update project` o `could not finalize project` | Sin filtrar SQL/rutas. Si commit se confirmó pero se perdió la respuesta, consultar el recurso. |

Precedencia propuesta: validar estructura/tipos/fechas/nombre del cuerpo antes de buscar proyecto; luego resolverlo en la unidad de trabajo; finalmente validar invariantes contra estado/fechas actuales y persistir. Cuerpo válido e ID desconocido → 404; cuerpo inválido incluso con ID desconocido → 400. En finalizar repetido se reconoce primero Finalizado después de validar cuerpo, conservando el fin. Se mantienen los errores existentes de US-07 para goal vacío o proyecto/Sprint no encontrado.

## Acceptance Criteria

Todos se encuentran en revisión como parte de esta especificación. Los que detallan PT no implican aprobación técnica por el equipo.

| ID | Criterio verificable | Fuente | Escenarios BDD-US-02 |
| --- | --- | --- | --- |
| AC-01 | Editar nombre/inicio de un proyecto conserva ID y permite consultar los mismos valores después de reiniciar con la misma base, sin alterar otro proyecto. | #2 AC original 1; D-02 | 01, 02 |
| AC-02 | Omitir inicio conserva, `null` lo elimina, fecha válida lo reemplaza; la edición no cambia estado ni fin. | D-02/04; PT-02 | 02, 03, 09 |
| AC-03 | Rechazar nombre blanco/nulo, fecha mal formada/imposible o inicio > fin sin cambios parciales; admitir igualdad y bisiesto válido. | US-01; D-06; PT-01/02 | 04, 05, 06, 10 |
| AC-04 | Confirmar finalización desde Planificado o En progreso, aun sin Sprints o inicio, persiste estado Finalizado y fin actual automático consultables tras reinicio. | #2 AC original 2; D-01/03/05/06 | 07 |
| AC-05 | Pedir confirmación antes de finalizar; cancelar no envía escritura ni guarda borrador. API sin confirmación verdadera no finaliza. | D-07; PT-02/05 | 08, 17 |
| AC-06 | Un Finalizado admite editar nombre/inicio y borrar inicio, pero no estado/fin ni reapertura. | D-02/04/06; PT-02 | 09, 10, 16 |
| AC-07 | Fechas y paso del tiempo no finalizan/inician automáticamente; se conserva alta Planificado y transición de US-07 sin inferir inicio. | D-01/02; US-01/07 | 11, 12 |
| AC-08 | Fin usa día del servidor en zona propuesta, independientemente de cliente/host; inicio > ese día impide finalizar. | D-03/06; PT-01 | 13, 14 |
| AC-09 | Repetir finalización conserva primer fin y una sola transición; repetir edición conserva datos. | D-04; PT-03 | 15, 21 |
| AC-10 | Finalizar preserva todos los Sprints y snapshots; sus consultas siguen disponibles aunque alguno esté abierto. | D-08/09 | 18 |
| AC-11 | Crear Sprint en Finalizado falla sin inserción ni regresión de estado, incluso concurrentemente con finalizar. | D-08; US-07; PT-03 | 19, 22 |
| AC-12 | Miembros y Product Backlog siguen operativos en Finalizado con sus reglas actuales y sin reabrirlo. | D-10; US-03/05/06 | 20 |
| AC-13 | Ausencia de proyecto, fallo de persistencia y contención producen error controlado sin cambios parciales de la operación. | #2; PT-03 | 23, 24, 25 |
| AC-14 | Contrato estricto de nuevas escrituras rechaza cuerpos/campos inválidos sin alterar el recurso; campos protegidos no son editables. | D-03/04; PT-02 | 16, 17, 26 |
| AC-15 | Edición y finalización concurrentes respetan el rango de fechas y no pierden estado/fin; cambios parciales de campos distintos se conservan. | D-04/06; PT-03 | 27, 28 |

## Integración, revisión y continuidad

- **US-01:** ampliar lectura/representación de Project con fechas opcionales. Alta sigue solo nombre, UUID estable y Planificado. Preservar validación/conservación de nombre actual y contratos anteriores; fechas nulas no se rellenan con el alta.
- **US-07:** incorporar D-08 a servicio y transacción real, preservando primer Sprint → En progreso y atomicidad. La SDD/BDD actuales no describen Finalizado; esta SDD aporta el incremento a revisar, sin reescribirlos ni cambiar su implementación en esta etapa.
- **US-03/05/06:** sus servicios comprueban existencia y no bloquean por estado; D-10 confirma esa continuidad. Conservar validaciones, UUID, estimaciones y consultas existentes.
- **US-08:** finalizar conserva snapshots/relaciones. No decidir aquí nuevas restricciones de asignación a un Sprint existente.
- **US-04/23/24:** consumir el mismo estado y las fechas del proyecto; no calcular finalización por fechas o métricas. No implementar dashboard/reportes aquí.
- **US-11 (cierre de Sprint):** existe una ambigüedad funcional transversal real: la regla vigente exige un siguiente Sprint para cerrar uno con historias incompletas, mientras D-08 impide crearlo después de finalizar y D-09 conserva el abierto. El equipo deberá decidir en US-11 si ese Sprint queda abierto o qué tratamiento reciben las incompletas al cerrarlo en un proyecto Finalizado. Esta SDD permite finalizar y conserva todo; no añade cierre, traslado, excepción a D-08 ni una condición previa para finalizar. No bloquea la especificación de US-02.

No se detectan otras ambigüedades funcionales que impidan definir US-02 con D-01…D-10. La aprobación de PT-01…PT-05 (zona, fechas, PATCH/null, confirmación API, repetición, política concurrente y migración) es revisión técnica pendiente, no reapertura de esas diez decisiones.

**Validación de esta etapa:** revisión documental y trazabilidad AC/BDD; no escenarios ejecutados, tests implementados, migraciones aplicadas ni evidencia TDD. Después de aprobar SDD/BDD, planificar TDD con unitarias de dominio/aplicación, integración SQLite/API (incluidas conexiones concurrentes) y confirmación/cancelación en UI. Los tests existentes son referencias de regresión, no evidencia de US-02. No se hace merge ni se cierra #2.
