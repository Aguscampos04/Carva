# SDD-US-12 — Registrar votos individuales ocultos

**Estado:** especificación aprobada para implementar por Loana Vargas el 2026-10-24, con la composición de participantes y el reemplazo de votos confirmados durante el análisis; revisión final del equipo pendiente.
**Trazabilidad:** enunciado §2.4 → EP-04 → US-12 / [Issue #12](https://github.com/Aguscampos04/Carva/issues/12) → esta SDD → [BDD-US-12](../bdd/BDD-US-12-votos-individuales-ocultos.feature).
**Dependencias:** US-03 y US-05, ambas cerradas al consultar GitHub el 2026-10-24.
**Fuentes:** Issue #12 y su AC; `docs/enunciado.md` §2.4; decisiones Planning Poker de `docs/contexto-proyecto.md`; D-05/D-07 y SDD-US-03; arquitectura y código vigente en `origin/main`.

## Objetivo

Registrar el voto individual de un miembro identificado para una User Story y una ronda de Planning Poker, manteniendo ocultos los votos de los demás participantes mientras la ronda no haya sido finalizada.

## Entradas

| Campo | Origen | Validación propuesta |
| --- | --- | --- |
| `projectID` | Ruta del proyecto | Debe identificar el proyecto existente. |
| `backlogItemID` | Ruta de la User Story | Debe pertenecer al proyecto indicado. |
| `roundID` | Recurso de ronda | Debe pertenecer a esa User Story y estar abierta. |
| `memberID` | Miembro seleccionado de la lista del proyecto (D-05) | Debe pertenecer al proyecto indicado. |
| `storyPoints` | Valor del voto | Entero no negativo, consistente con el tipo actual de Story Points. No se establece una escala (por ejemplo, Fibonacci) porque no está aprobada. |

## Salidas esperadas

- Un voto aceptado queda asociado a un miembro, una User Story y una ronda.
- La respuesta de registro no incluye votos de otros miembros.
- Antes de finalizar la ronda, una consulta no revela votos ajenos.
- La operación no altera los Story Points provisionales ni definitivos del Backlog.

## Reglas de negocio

| ID | Regla | Origen / estado |
| --- | --- | --- |
| BR-01 | Cada voto debe conservar la identidad del miembro, la User Story y la ronda a la que pertenece. | AC-01 de Issue #12. |
| BR-02 | Los votos ajenos permanecen ocultos antes de finalizar la ronda. | AC-02 de Issue #12; Planning Poker aprobado. |
| BR-03 | La identidad se selecciona entre los miembros del proyecto y se valida su pertenencia. | D-05 de US-03, aprobado. |
| BR-04 | El registro de votos no actualiza Story Points del Backlog; solo US-15 registra el valor acordado. | Contexto aprobado. |
| BR-05 | Se aceptan enteros no negativos, sin imponer una secuencia de estimación no aprobada. | Compatible con el tipo actual de Story Points. |
| BR-06 | Todos los miembros registrados del proyecto participan en la ronda; la lista se conserva como participantes de esa ronda. | Decisión de Loana Vargas del 2026-10-24. |
| BR-07 | Un miembro puede reemplazar su propio voto mientras la ronda esté abierta. | Decisión de Loana Vargas del 2026-10-24. |

## Restricciones

- Mantener Go, API REST, SQLite, capas aprobadas y migraciones SQL embebidas.
- No agregar autenticación o permisos; D-07 los excluye del MVP. El ocultamiento es comportamiento del producto para miembros identificados, no una garantía frente a un cliente que manipule identificadores.
- La revelación simultánea y detección de diferencias pertenecen a US-13; nuevas rondas posteriores pertenecen a US-14.
- La creación de la ronda inicial es parte de esta historia; crear rondas posteriores pertenece a US-14.
- Sin cambios de Story Points ni de la estimación definitiva.

### Contrato REST

| Operación | Solicitud | Respuesta correcta |
| --- | --- | --- |
| Abrir ronda inicial | `POST /api/projects/{projectID}/backlog/{storyID}/estimation/rounds`, sin cuerpo | `201` y ronda abierta con sus participantes. Solo crea la primera ronda; las siguientes son US-14. |
| Consultar rondas | `GET /api/projects/{projectID}/backlog/{storyID}/estimation/rounds` | `200` y arreglo, vacío antes de iniciar. No contiene votos. |
| Registrar o corregir voto | `PUT /api/projects/{projectID}/backlog/{storyID}/estimation/rounds/{roundID}/votes/{memberID}`, cuerpo `{"story_points": 5}` | `200` y `{"saved":true}`. No contiene votos. |
| Consultar votos visibles | `GET /api/projects/{projectID}/backlog/{storyID}/estimation/rounds/{roundID}/votes?member_id={memberID}` | `200` y solo el voto propio mientras esté abierta. Una ronda finalizada permite revelar el conjunto a sus participantes cuando US-13 implemente esa transición. |

Los errores de transporte o valor inválido devuelven `400`; proyecto, historia o ronda inexistente devuelven `404`; ronda inicial duplicada o cerrada devuelve `409`; proyecto sin miembros devuelve `422`; fallos inesperados devuelven `500`. Se usa el formato existente `{"error":"mensaje"}`.

## Casos límite y errores

- Proyecto, historia o miembro inexistente: rechazo controlado; no se registra voto.
- Miembro existente de otro proyecto: rechazo; no se registra voto.
- Historia que no pertenece al proyecto de la ruta: rechazo; no se registra voto.
- Identificador de ronda inexistente, de otra historia o cerrada: rechazo; no se registra voto.
- Valor ausente, no entero o negativo: rechazo; no se registra voto.
- Dos registros concurrentes para el mismo miembro y ronda: debe persistir un solo voto vigente para ese miembro; el último reemplazo confirmado es el valor visible mientras la ronda siga abierta.
- Error de almacenamiento: error controlado y sin éxito aparente.
- Lectura de votos antes de finalizar: ningún voto ajeno se incluye en la respuesta.

## Criterios de aceptación

Se conservan los dos AC de la Issue #12:

- **AC-01:** cada voto queda atribuido al miembro, historia y ronda correspondientes.
- **AC-02:** antes de finalizar una ronda, sus votos individuales no son visibles para los otros participantes.

La ronda inicial se abre de forma explícita y captura los miembros del proyecto en ese momento. La UI identifica al miembro seleccionándolo de esa lista; no se agrega autenticación ni permisos. El reemplazo usa la misma identidad de miembro y solo se permite mientras la ronda esté abierta. La revelación y detección de diferencias quedan para US-13; las rondas posteriores, para US-14.
