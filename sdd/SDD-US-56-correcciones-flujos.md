# SDD — US-56: Identificación, esfuerzo al completar y Planning Poker

## Objetivo

Corregir tres problemas reportados por el equipo: distinguir historias homónimas en Backlog y Sprints, exigir horas trabajadas antes de completar una historia y permitir votar en Planning Poker desde una computadora compartida sin perder la identidad seleccionada.

Issue: [#56](https://github.com/Aguscampos04/Carva/issues/56). Trazabilidad: enunciado §2.2–§2.5; US-05, US-06, US-08, US-09, US-12–US-15 y US-16–US-17; corrección transversal #56.

## Inputs

- Historia de Backlog ya persistida, cuyo UUID existente es su identidad estable.
- Sprint abierto al completar, aceptación explícita de sus criterios y, solo si la historia no tiene esfuerzo activo, una entrada de esfuerzo con miembro del proyecto, fecha, actividad y horas.
- Votante seleccionado entre los participantes congelados de la ronda; voto entero no negativo.

## Expected Outputs

- UUID completo visible junto al título en las filas del Backlog y en los resúmenes de historias planificadas en Sprints.
- La finalización devuelve la historia completada solo cuando se verificaron sus criterios y existe esfuerzo activo. Cuando no hay esfuerzo activo, el formulario incorporado al flujo registra la nueva entrada y completa la historia en una sola operación.
- Planning Poker conserva al miembro seleccionado al recargar la vista y devuelve votos ajenos únicamente después de finalizar la ronda. Si faltan votos, informa los participantes pendientes sin incluir estimaciones.

## Business Rules

1. El identificador visible es el UUID que la historia ya tiene. No se crea otro ID, no se usa el número de Issue y la identidad permanece al editar la historia o trasladar su snapshot a otro Sprint.
2. Una entrada de esfuerzo activa asociada a la historia satisface la condición de horas aunque pertenezca a un Sprint anterior o no tenga Sprint. Una entrada eliminada lógicamente no cuenta. Story Points y horas estimadas no cuentan como esfuerzo real.
3. Si no existe esfuerzo activo, la petición de finalización debe incluir una entrada válida: miembro del proyecto, fecha admitida por el proyecto, actividad no vacía y horas positivas con precisión y límites definidos por US-16. La entrada se asocia a la historia y al Sprint abierto desde el que se finaliza.
4. La inserción de la entrada y la finalización son atómicas. Un error no completa la historia ni persiste una entrada parcial. Una repetición con los mismos datos de la entrada previamente confirmada completa la historia sin insertar otra fila ni modificar la existente. Una petición obsoleta que intenta agregar horas cuando ya existe una entrada distinta falla sin alterar registros.
5. Se conserva la verificación explícita de todos los criterios de aceptación y el Sprint debe seguir abierto.
6. Cambiar el votante solo cambia el contexto de edición y consulta de su propio voto; nunca revela estimaciones de otros participantes mientras la ronda está abierta.
7. Una ronda solo se revela cuando cada participante congelado tiene un voto válido. Ante una ronda incompleta, el error puede incluir únicamente las identidades faltantes; nunca incluye sus Story Points.
8. La detección de diferencias, el acuerdo unánime, la estimación definitiva y las rondas posteriores conservan las reglas de US-12–US-15.

## Constraints

- Go mantiene reglas de negocio, validación de horas y coordinación transaccional; HTML/CSS/JavaScript presenta los identificadores y formularios; la API REST y SQLite mantienen sus capas existentes.
- Se reutilizan el UUID y la tabla de esfuerzo actuales. No se requiere una nueva ruta ni una migración de esquema.
- El endpoint existente `POST /api/projects/{projectID}/sprints/{sprintID}/backlog-items/{backlogItemID}/completion` conserva `acceptance_criteria_verified` y admite opcionalmente `effort: {member_id, work_date, activity, hours}`. La asociación de Sprint se deriva del `{sprintID}` de la ruta. Sin esfuerzo previo ni objeto `effort`, responde conflicto con código `effort_required`; ante concurrencia con un esfuerzo distinto responde `effort_already_registered`.
- El endpoint existente de finalización de Planning Poker mantiene su ruta. Si faltan votos, el error 409 agrega `code: "votes_pending"` y `missing_participants` con IDs, nunca votos.

## Edge Cases

- Dos historias con el mismo título deben distinguirse por sus UUID completos en las dos listas.
- Cambiar título, mover una historia incompleta al próximo Sprint o cargar un snapshot histórico no genera una identidad nueva.
- Esfuerzo activo previo en cualquier Sprint cuenta; esfuerzo eliminado no cuenta; falta de puntos definitivos o de horas estimadas no sustituye el esfuerzo real.
- Entrada nueva inválida, miembro ajeno al proyecto, Sprint no asociado, Sprint cerrado, fecha inválida o límite diario excedido deben conservar la historia incompleta y el borrador en pantalla.
- Repetir exactamente una petición cuyo commit tuvo respuesta incierta no duplica el esfuerzo. Una entrada guardada diferente no se reemplaza.
- Cambiar repetidamente entre miembros conserva cada voto propio y mantiene ocultos los votos ajenos. Una ronda realmente incompleta identifica solo a quienes faltan.

## Error Conditions

- 400 para aceptación ausente, JSON inválido o entrada con campos/formato inválidos.
- 409 para criterios no verificados, Sprint cerrado, esfuerzo requerido, una entrada distinta ya registrada o restricciones de fechas/límite diario de US-16.
- Errores de almacenamiento no cambian parcialmente ni la historia ni el esfuerzo.
- 409 `votes_pending` incluye las identidades participantes que aún no votaron y ninguna estimación.

## Acceptance Criteria

- **AC-56-01:** Backlog y Sprint muestran el UUID completo junto al título, distinguen títulos iguales y conservan el ID tras edición y traslado.
- **AC-56-02:** No se completa una historia sin criterios verificados y sin al menos una entrada de esfuerzo activa.
- **AC-56-03:** Si no hay esfuerzo previo, se puede registrar una entrada válida dentro del flujo de finalización; el esfuerzo y la finalización se guardan juntos o ninguno se guarda.
- **AC-56-04:** Si falla la validación o persistencia, la historia sigue incompleta, los registros anteriores no cambian y el formulario conserva sus valores.
- **AC-56-05:** Reintentar una finalización confirmada no duplica la entrada ni reemplaza otra; una petición obsoleta tampoco modifica registros.
- **AC-56-06:** En una computadora, cada miembro puede seleccionar su identidad, registrar o corregir su propio voto y conservar los votos de los demás sin verlos antes de revelar.
- **AC-56-07:** La revelación se permite con todos los votos; si faltan, se identifican los participantes faltantes sin mostrar estimaciones.
- **AC-56-08:** Diferencias, acuerdo, historia de rondas y estimación acordada continúan funcionando.

## Reproducción y evidencia

- La Issue no incluye versión desplegada, historia/ronda concreta, capturas ni mensajes de error; no se atribuye una sesión manual al equipo.
- La revisión del código actual reproduce de forma determinista la causa candidata de Planning Poker: el evento `change` vuelve a ejecutar `load()`, que elimina y reconstruye el selector sin conservar el miembro previo, por lo que el navegador selecciona el primero.
- La regresión automatizada documentará los pasos reproducibles de cada hallazgo. La nueva prueba manual del equipo de los tres flujos queda pendiente de realizarse y registrarse por sus integrantes.
