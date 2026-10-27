# SDD-US-21 — Consultar métricas de horas

## 1. Objetivo

Consultar horas estimadas, horas reales y su desviación por Sprint y por proyecto, reutilizando las estimaciones de US-16 y los registros de esfuerzo activos.

## 2. Inputs

- Identificador positivo de Sprint para consultar sus métricas.
- Identificador de proyecto para consultar las métricas agregadas.
- Las consultas no reciben ni modifican horas; usan datos persistidos del proyecto.

## 3. Expected Outputs

- Horas estimadas, horas reales y desviación en centésimas de hora, serializadas como números decimales con dos posiciones.
- Las métricas no disponibles se representan como `null`; la interfaz las presenta como **«No disponible»**.
- La consulta de proyecto incluye sus métricas y las métricas por Sprint para permitir la selección de un Sprint en la interfaz.
- Las consultas son de solo lectura.

## 4. Business Rules

- **BR-01:** Las horas estimadas de un Sprint son la suma de `estimated_hours_hundredths` de las historias en el snapshot de ese Sprint. Los snapshots reflejan la estimación vigente conforme a D-21.
- **BR-02:** Las horas reales de un Sprint son la suma de registros de esfuerzo activos asociados explícitamente a ese Sprint. Registros sin Sprint no forman parte de una métrica individual de Sprint.
- **BR-03:** Las horas estimadas del proyecto son la suma de las estimaciones vigentes de las User Stories distintas del Product Backlog actual. Una historia que aparece en varios snapshots se cuenta una sola vez.
- **BR-04:** Las horas reales del proyecto son la suma de todos los registros de esfuerzo activos del proyecto, incluidos los registros sin Sprint.
- **BR-05:** En ambos alcances, la desviación es `horas reales − horas estimadas`. Un valor positivo indica exceso sobre la estimación; uno negativo indica horas restantes respecto de la estimación.
- **BR-06:** Si alguna historia incluida en una suma de horas estimadas carece de estimación, las horas estimadas de ese alcance y su desviación son no disponibles. Las horas reales siguen disponibles independientemente.
- **BR-07:** Una suma vacía tiene valor cero: un Sprint sin historias tiene 0 horas estimadas y uno sin registros activos asociados tiene 0 horas reales. Un proyecto sin historias tiene 0 horas estimadas; sin registros activos, 0 horas reales.
- **BR-08:** Las entradas con `deleted_at` no nulo se excluyen. No se convierten Story Points a horas.
- **BR-09:** Un Sprint o proyecto inexistente responde recurso no encontrado. Un identificador de Sprint que no sea un entero positivo se rechaza como solicitud inválida.
- **BR-10:** Las métricas usan horas exactas en centésimas, evitando aritmética de punto flotante durante los cálculos.

## 5. Constraints

- Solo se consultan registros pertenecientes al proyecto solicitado.
- La asociación de una entrada a Sprint se determina por su `sprint_id` explícito, incluso si una historia se trasladó luego a otro Sprint.
- No se agregan tablas, migraciones, nuevos frameworks ni cambios en las reglas de registro de esfuerzo.
- No se exponen datos personales ni el detalle de las entradas en estas métricas agregadas.

## 6. Edge Cases

- Sprint sin historias y sin esfuerzo: estimadas 0, reales 0, desviación 0.
- Historias con estimación nula: estimadas y desviación no disponibles; reales disponibles.
- Registros activos sin Sprint: cuentan para el proyecto, no para métricas de Sprint.
- Entradas eliminadas lógicamente: excluidas de ambos niveles.
- Un valor real puede superar o quedar por debajo de la estimación.
- Cambios posteriores de estimación se reflejan de acuerdo con D-21, también en snapshots cerrados; el esfuerzo histórico y su asociación no se modifica.

## 7. Error Conditions

- Proyecto inexistente: HTTP 404.
- Sprint inexistente o perteneciente a otro proyecto: HTTP 404.
- `sprintID` vacío, no numérico, cero o negativo: HTTP 400.
- Error interno de lectura: HTTP 500 con mensaje genérico, sin detalles SQL.

## 8. Acceptance Criteria

- **AC-01:** Para un Sprint se muestran horas estimadas, reales y desviación usando las fuentes de US-16 y las fórmulas aprobadas.
- **AC-02:** Para un proyecto se muestran esas métricas agregadas, contando cada historia actual una sola vez en estimadas e incluyendo esfuerzo sin Sprint en horas reales.
- **AC-03:** Si falta una estimación necesaria, estimadas y desviación muestran «No disponible»; no se muestra cero ficticio.
- **AC-04:** Las consultas excluyen esfuerzo eliminado lógicamente y no modifican datos.
- **AC-05:** Un exceso real produce desviación positiva; un esfuerzo inferior produce desviación negativa.
- **AC-06:** Los errores de proyecto, Sprint e identificador inválidos usan los códigos HTTP definidos.

## 9. BDD / Test Traceability

- AC-01/05 → BDD-US-21-01/02 → pruebas de sumas exactas por Sprint y desviación con signo.
- AC-02/03 → BDD-US-21-03/04 → pruebas de agregado de proyecto, historias sin estimación y esfuerzo sin Sprint.
- AC-04 → BDD-US-21-05 → prueba de exclusión de borrados y lectura sin escrituras.
- AC-06 → BDD-US-21-06 → pruebas de identificadores y recursos inexistentes.

## 10. Decision Record

La propuesta de fórmulas BR-01…BR-10 fue aprobada por Loana Vargas el 2026-10-27 en la conversación de desarrollo de US-21.
