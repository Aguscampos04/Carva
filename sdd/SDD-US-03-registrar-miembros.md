# SDD-US-03 — Registrar e identificar miembros del proyecto

**Estado:** especificación autorizada para implementar por el responsable el 2026-10-16, con los ajustes de revisión de esta etapa; implementación realizada y verificación técnica registrada en [evidencia](../docs/evidencia-ia-us-03.md). Revisión y aprobación final de la historia por el equipo pendientes.
**Fecha:** 2026-10-16.
**Trazabilidad:** [enunciado](../docs/enunciado.md) §2.1 → EP-01 → US-03 / [Issue #3](https://github.com/Aguscampos04/Carva/issues/3) → esta SDD → AC → [BDD-US-03](../bdd/BDD-US-03-registrar-miembros.feature).
**Dependencia:** US-01, proyecto existente con UUID. Consumidores futuros: US-12 (§2.4, votos individuales) y US-16 (§2.5, esfuerzo).
**Fuentes consultadas:** Issue #3 y sus comentarios (sin comentarios al consultar), AGENTS.md, enunciado, contexto, reglas del agente, arquitectura, SDD de US-01/05/06/07/08 y convenciones del código existente. Decisiones funcionales: [registro D-01 a D-07](../docs/contexto-proyecto.md#miembros-del-proyecto--us-03--issue-3).

## Objective

Como integrante del equipo, quiero registrar e identificar miembros del proyecto para atribuirles votos de Planning Poker y registros de esfuerzo.

US-03 entrega el alta y la consulta de la lista de miembros en una sección del proyecto. Cada miembro conserva una identidad estable independiente de su nombre. La selección en votos y esfuerzo corresponde a las futuras US-12 y US-16, usando el UUID del miembro de ese proyecto.

## Inputs

| Entrada | Tipo y origen | Regla |
| --- | --- | --- |
| `projectID` | String en la ruta; UUID emitido por US-01 | Proyecto existente; no se ingresa en el formulario de miembro. |
| `name` | String JSON, único dato del formulario | Obligatorio. Recortar espacios exteriores; rechazar si queda vacío; conservar el interior, mayúsculas, tildes y demás caracteres. |

El cliente no elige el UUID ni la pertenencia mediante el cuerpo. Propuesta de transporte: ignorar propiedades adicionales, siguiendo las altas existentes; solo `name` se consume, y la ruta determina el proyecto. No se añaden email, rol, credenciales, fecha ni identidad global de persona.

Propuesta de precisión técnica de D-06: usar `strings.TrimSpace` en Go, incluyendo tabulaciones, saltos de línea y espacios Unicode exteriores. No colapsar espacios internos ni normalizar Unicode. No se propone un máximo funcional de longitud, restricción de caracteres ni cupo de miembros porque no están exigidos ni aprobados.

## Expected Outputs

- Alta válida: miembro persistido con `id` UUID generado por el servidor, `project_id` y `name` normalizado; respuesta HTTP `201`.
- Lista: HTTP `200` con arreglo JSON de todos los miembros del proyecto, cada uno con esos tres campos. Sin miembros: `[]`, nunca `null`.
- Consultas posteriores, incluso después de reiniciar la aplicación sobre la misma base, conservan UUID, nombre y proyecto.
- La sección muestra formulario de nombre y lista; tras un alta confirmada actualiza la lista. Un error se muestra como tal, sin agregar un miembro aparente ni presentar una lista fallida como vacía.
- UI para distinguir homónimos: mostrar nombre como dato principal y UUID completo como dato secundario legible en cada entrada. Los futuros selectores mostrarán ambos y enviarán el UUID; no se decide por el nombre ni por la posición en la lista. No se crea un selector de identidad global o sesión en US-03.

## Business Rules

| ID | Regla | Origen |
| --- | --- | --- |
| BR-01 | Todo miembro pertenece a un proyecto existente; registrar un miembro no crea un proyecto. | §2.1, D-04 |
| BR-02 | Nombre requerido, recortado en los extremos y no vacío después del recorte; interior conservado. | D-01, D-06 |
| BR-03 | UUID automático, estable y distinto para registros independientes; no se deriva del nombre. | D-02, D-03 |
| BR-04 | Dos altas con igual nombre en un mismo proyecto crean dos miembros; no se deduplican ni sobrescriben. | D-03 |
| BR-05 | La misma persona en dos proyectos tiene registros y UUID independientes; no existe catálogo global de personas. | D-04 |
| BR-06 | La lista contiene exclusivamente los miembros del proyecto consultado. | D-04, D-05 |
| BR-07 | La identidad usada en funciones futuras es el UUID seleccionado de la lista del proyecto; esas funciones deberán verificar su pertenencia. | D-05; AC de Issue #3 |
| BR-08 | La sección permite alta y listado, sin edición, eliminación, autenticación ni permisos. | D-07 |

Validación central en Go. El alta no modifica datos ni estado del proyecto ni crea votos o registros de esfuerzo. No se introduce una restricción por estado del proyecto, dado que no fue solicitada.

## Constraints

### Contratos REST implementados

Se reutilizan `/api/projects`, JSON, nombres de campos `snake_case` para la relación como en Backlog, y el formato de error existente `{"error":"mensaje"}`. No se cambia el contrato de otros recursos.

| Operación | Solicitud | Respuesta correcta |
| --- | --- | --- |
| Alta | `POST /api/projects/{projectID}/members`, cuerpo `{"name":"  Ana  Pérez  "}` | `201` y objeto miembro con `name` igual a `Ana  Pérez`. |
| Listado | `GET /api/projects/{projectID}/members`, sin cuerpo requerido | `200` y arreglo de objetos miembro, posiblemente vacío. |

Ejemplo de objeto devuelto (UUID ilustrativos):

```json
{
  "id": "b1d168da-325a-40d1-90b7-b5a1ba006389",
  "project_id": "0828dc84-050f-421b-b842-420f63851a2b",
  "name": "Ana  Pérez"
}
```

- Respuestas JSON: `Content-Type: application/json; charset=utf-8`.
- POST recibe exactamente un objeto JSON. Cuerpo vacío, JSON malformado, varios valores JSON, raíz `null`/arreglo/escalar y `name` ausente, `null` o no string se rechazan con `400`. Una cadena vacía o solo espacios también devuelve `400`.
- Ajuste de revisión: validar el objeto completo y el fin de entrada con un mismo decoder; rechazar claves superiores duplicadas y bytes UTF-8 inválidos. El campo aceptado es exactamente `name`, sensible a mayúsculas. Las propiedades adicionales siguen ignorándose, tal como se especificó; no permiten alterar ID ni proyecto. Estos ajustes se limitan al alta de miembros.
- GET verifica la existencia del proyecto antes de devolver la lista. Una ruta con identificador que no resuelve a un proyecto devuelve `404`, incluso si no tiene forma de UUID; se mantiene el criterio de búsqueda existente de US-01, sin añadir validación global de UUID.
- Propiedades extra como `id`, `project_id` o `email` no se almacenan ni controlan el alta. Un nombre válido con extras crea un miembro con UUID nuevo en el proyecto de la ruta.
- Sin orden garantizado, paginación, filtros ni búsqueda en esta historia, siguiendo el listado de US-06.
- No se propone GET individual ni `Location` a un recurso individual inexistente. La identificación se obtiene del objeto creado y del listado. PUT, PATCH y DELETE no se exponen para miembros.
- Cada POST es una nueva alta: no se agrega idempotencia ni deduplicación. Si se pierde la respuesta después de persistir, debe consultarse la lista antes de decidir repetir; repetir puede crear un homónimo con otro UUID.

### Persistencia implementada

Migraciones SQL versionadas y embebidas: `000006_members.sql` crea la tabla y `000007_members_name_storage.sql` ajusta su CHECK para medir bytes (`length(CAST(name AS BLOB)) > 0`). La segunda es aditiva: conserva los registros y evita rechazar nombres que comienzan con NUL, ya que SQLite `length(TEXT)` se detiene en ese carácter. No se modificaron migraciones aplicadas. Se verificó la actualización desde el esquema anterior, la conservación de UUID/datos y la aplicación idempotente de las migraciones.

| Tabla `members`: columna | Definición propuesta |
| --- | --- |
| `id` | `TEXT PRIMARY KEY NOT NULL`: UUID generado en aplicación; reutilizar `github.com/google/uuid`, ya presente. |
| `project_id` | `TEXT NOT NULL REFERENCES projects(id)`: asociación obligatoria. |
| `name` | `TEXT NOT NULL CHECK (length(CAST(name AS BLOB)) > 0)`: texto normalizado y validado por Go, sin restricciones adicionales de caracteres. |

Índice no único sobre `project_id` para el listado. Sin unicidad sobre nombre ni sobre `(project_id, name)`. No agregar tabla de personas ni relaciones entre miembros de diferentes proyectos. El CHECK no reemplaza la validación Unicode en Go; SQLite `trim` no debe definir D-06.

Integridad: verificar proyecto existente y guardar sin alta parcial. `OpenPath` configura `_pragma=foreign_keys(1)` en el DSN del driver actual, ejecutado para cada conexión del pool. Se verificó con tres conexiones simultáneas, rechazo de inserciones huérfanas, una base en memoria y rutas con espacios y `#`. Se ejecutó regresión de los repositorios existentes; no se cambió el driver ni se agregaron dependencias.

Una inserción es atómica; si falla, no queda un miembro incompleto. La lista se filtra siempre por `project_id`. Una colisión extraordinaria de UUID no sobrescribe al registro existente: la PK rechaza la inserción y se informa error interno controlado. Fallos de generación de UUID o persistencia no producen respuesta exitosa. No se introducen cascadas de eliminación ni políticas de borrado.

### Responsabilidades por capa

| Capa / ubicación existente | Responsabilidad prevista |
| --- | --- |
| Presentación / `internal/presentation/httpapi`, `web/static` | Traducir contratos, mostrar sección con formulario y lista para el proyecto consultado; renderizar nombres como texto; evitar mezclar respuestas de distintos proyectos al cambiar la selección. |
| Aplicación / `internal/application/members` | Coordinar registrar y listar, generar UUID, comprobar proyecto mediante puertos y persistir. |
| Dominio / `internal/domain/member` | Miembro e invariantes de nombre e identidad, sin dependencia de HTTP o SQLite. |
| Persistencia / `internal/persistence/sqlite` | Implementar puertos, almacenamiento e integridad referencial, listado aislado y migración embebida. |
| Composición / `cmd/carva` | Conectar repositorio, servicio y rutas al servidor y frontend actuales. |

Se conservan Go, SQLite, frontend HTML/CSS/JS embebido y arquitectura aprobada. No se incorpora framework BDD ni nueva tecnología.

## Edge Cases

- Primer miembro y lista inicialmente vacía; listar no crea ni modifica registros.
- Nombre con un solo carácter significativo: válido; no se exige nombre y apellido.
- Nombres con tildes, apóstrofes, guiones, espacios múltiples internos y cambios de mayúsculas: conservados; no se comparan para identidad.
- Nombres que quedan iguales al recortar exteriores: siguen siendo registros independientes.
- Alta de homónimos consecutiva o concurrente: UUID distintos y ninguna sobrescritura si ambas altas se confirman.
- Dos proyectos representan a la misma persona mediante registros independientes y listas separadas.
- Reinicio y recarga: identidad y pertenencia conservadas.
- Fallo al listar: error visible, sin mostrar `[]` como si la consulta hubiera tenido éxito.
- Respuesta de un proyecto anterior recibida tras cambiar de proyecto: no se muestra en el proyecto actual.
- Texto que parece HTML se conserva como nombre y se presenta como texto, sin interpretarlo como contenido ejecutable; no se agrega una restricción funcional de caracteres.
- Error previo a confirmar la escritura: no se informa éxito; una pérdida de conexión posterior a confirmar puede dejar el alta realizada, por lo que no se promete entrega exactamente una vez.

## Error Conditions

Los mensajes siguientes son propuestas de contrato, no decisiones funcionales aprobadas. La UI puede presentarlos en español. No se expone SQL, rutas de archivos ni detalles internos.

| Condición | HTTP | Cuerpo propuesto | Efecto |
| --- | --- | --- | --- |
| Cuerpo/JSON inválido o `name` numérico, booleano, arreglo u objeto | 400 | `{"error":"invalid JSON request"}` | Sin alta. |
| Nombre ausente, `null`, vacío o blanco | 400 | `{"error":"member name must not be blank"}` | Sin alta. |
| Proyecto no encontrado, en alta o listado | 404 | `{"error":"project not found"}` | Sin crear proyecto ni miembro. |
| Fallo interno, generación de UUID o escritura | 500 | `{"error":"could not register member"}` | Sin éxito aparente ni sobrescritura. |
| Fallo de lectura del proyecto o lista | 500 | `{"error":"could not retrieve members"}` | Sin lista parcial ni vacía de éxito. |

Propuesta de precedencia: en POST, validar cuerpo y nombre antes de buscar proyecto; para cuerpo válido, resolver proyecto y luego escribir. En GET, resolver proyecto y luego listar. Un proyecto inexistente con nombre inválido devuelve `400`; uno inexistente con nombre válido devuelve `404`.

## Acceptance Criteria

Los dos criterios originales de la Issue se conservan; los restantes detallan decisiones aprobadas y la robustez del contrato propuesto.

| ID | Criterio verificable | Fuente |
| --- | --- | --- |
| AC-01 | Dado un proyecto, al registrar un miembro puede identificarse entre los miembros de ese proyecto por su UUID estable y nombre. | Issue #3, primer AC; D-01/02 |
| AC-02 | La identidad permite atribuir votos y esfuerzo al miembro correcto, aunque haya homónimos; las funciones futuras usarán UUID y verificarán proyecto. US-03 entrega esos datos; la atribución efectiva debe validarse en US-12 y US-16 cuando existan. | Issue #3, segundo AC; D-03/04/05 |
| AC-03 | Se rechaza un nombre ausente, vacío o blanco sin crear miembro; se recortan extremos y se conserva el resto. | D-01/06 |
| AC-04 | Homónimos en un proyecto tienen UUID diferentes, sin deduplicación; altas en proyectos distintos son independientes. | D-03/04 |
| AC-05 | La lista devuelve solo miembros del proyecto y `[]` si existe sin miembros; proyecto inexistente se informa sin crear registros. | D-04/05, BR-01/06 |
| AC-06 | La sección permite alta con solo nombre y muestra la lista actualizada, distinguiendo homónimos; no ofrece edición, eliminación, autenticación ni permisos. | D-01/07; propuesta de UI |
| AC-07 | Alta y listado manejan errores de transporte y almacenamiento de forma controlada; no hay altas parciales ni listas parciales presentadas como éxito. | Robustez del contrato propuesto |
| AC-08 | El servidor controla UUID y pertenencia; campos extra no los alteran; UUID, nombre y proyecto sobreviven recarga y reinicio. | D-02/04; persistencia aprobada |

## BDD / Test Traceability

Cada escenario lleva etiqueta `@BDD-US-03-NN` y referencias `@AC-NN` en el archivo BDD. Su automatización concreta y las limitaciones se detallan en la [matriz BDD–tests](../docs/evidencia-ia-us-03.md#trazabilidad-bdd-y-pruebas).

| AC | Escenarios BDD-US-03 | Comportamiento validado / pendiente |
| --- | --- | --- |
| AC-01 | 01, 07, 09 | Alta/listado y estabilidad de identidad. |
| AC-02 | 04, 05 | Identidades distinguibles en US-03; atribución efectiva pendiente en US-12/US-16. |
| AC-03 | 02, 03, 17 | Validación y normalización del nombre. |
| AC-04 | 04, 05, 18 | Homónimos y registros independientes, incluso concurrencia. |
| AC-05 | 05, 06, 10, 11 | Lista vacía, pertenencia y proyecto inexistente. |
| AC-06 | 04, 07, 15, 16 | Formulario, listado, texto seguro y contexto de proyecto. |
| AC-07 | 08, 10, 11, 12, 13, 14, 19 | Contrato HTTP, fallos y ausencia de cambios parciales. |
| AC-08 | 01, 09, 14 | Autoridad del servidor y persistencia. |

Cadena actual: §2.1 → EP-01 → US-03 / #3 → SDD → AC → BDD → tests Go/Node → dominio/aplicación/SQLite/HTTP/web → commits RED/GREEN/REFACTOR. Véase la evidencia para resultados reales. AC-02 se valida aquí como disponibilidad de identidades distintas; atribuir votos/esfuerzo efectivamente sigue pendiente en US-12/US-16.

## Revisión y pendientes

No se detectaron ambigüedades funcionales bloqueantes con D-01 a D-07. Rutas y JSON, errores/precedencia, precisión Unicode del recorte, esquema/integridad SQLite y UUID secundario fueron autorizados para implementar en esta etapa y están implementados. Falta revisión y aprobación final de la historia por el equipo.

No se requieren decisiones nuevas de autenticación, datos personales ni de identidad global para US-03. El detalle del selector y la validación de atribución en votos/esfuerzo se tratará en US-12/US-16, sin bloquear alta y listado. Cualquier límite de longitud/cupo o regla por estado del proyecto que se solicite posteriormente será una decisión funcional nueva.

El responsable autorizó implementar esta SDD y sus BDD con los ajustes de revisión indicados: JSON estricto, claves foráneas efectivas, escenarios 07/17 sin múltiples acciones y UUID completo secundario. También autorizó commits TDD, push y un único PR hacia main relacionado con #3, sin merge ni cierre. Las menciones a propuestas técnicas en el cuerpo describen el diseño ahora autorizado. La revisión y aprobación final de la historia por el equipo continúan pendientes.
