# SDD-US-04 — Consultar estado del proyecto

**Estado:** Especificación derivada de los criterios de aceptación vigentes de Issue #4. La lógica que satisface los estados fue incorporada por US-01, US-07 y US-02; esta historia documenta y verifica su consulta.
**Trazabilidad:** requisito §2.1 → EP-01 → US-04 / Issue #4 → US-01, US-07 y US-02 → esta SDD → AC → [BDD-US-04](../bdd/BDD-US-04-consultar-estado-proyecto.feature) → `TestUS04ProjectStatusCanBeConsultedThroughoutLifecycle`.

## 1. Objetivo

Permitir consultar el estado vigente y persistido de un proyecto para conocer si está **Planificado**, **En progreso** o **Finalizado**.

## 2. Entradas

| Entrada | Tipo / origen | Regla |
| --- | --- | --- |
| Identificador del proyecto | UUID en `GET /api/projects/{id}` | Debe identificar un proyecto existente. |

La consulta no recibe ni modifica campos del proyecto. El estado no se calcula a partir de fechas ni se elige en la solicitud.

## 3. Salidas esperadas

- Proyecto existente: `200 OK` con el recurso persistido, incluido `status` con uno de los tres estados admitidos.
- Proyecto inexistente: `404 Not Found`.
- Error de persistencia: `500 Internal Server Error`, sin respuesta de éxito.
- La consulta no altera el proyecto.

## 4. Reglas de negocio

- **BR-01:** el estado consultado es el estado persistido del proyecto.
- **BR-02:** los únicos estados del flujo son `Planificado`, `En progreso` y `Finalizado`.
- **BR-03:** US-01 crea el proyecto como `Planificado`.
- **BR-04:** la creación del primer Sprint cambia `Planificado` a `En progreso`; los Sprints siguientes conservan `En progreso` (US-07).
- **BR-05:** la finalización manual confirmada de US-02 cambia el estado a `Finalizado`.
- **BR-06:** consultar no produce transiciones de estado.

## 5. Restricciones

- Se reutiliza `GET /api/projects/{id}` y el recurso compartido de proyectos; no se agrega una ruta ni un campo de estado alternativo.
- Go mantiene las reglas de dominio, la API REST expone JSON y SQLite conserva el estado, según la arquitectura aprobada.
- La consulta no finaliza proyectos automáticamente ni deriva estado de fechas.
- Los cambios de estado pertenecen a US-01, US-07 y US-02; US-04 los presenta al consultarlos.

## 6. Casos límite

- Un proyecto recién creado se consulta como `Planificado`.
- El primer Sprint se consulta como `En progreso`.
- Un proyecto finalizado se consulta como `Finalizado`.
- Consultar reiteradamente devuelve el mismo estado mientras ninguna operación autorizada lo cambie.

## 7. Condiciones de error

- ID sin proyecto correspondiente: `404`.
- Error al leer SQLite: `500`.
- Ninguna condición de error debe representarse como una consulta exitosa.

## 8. Criterios de aceptación

- **AC-01:** consultar un proyecto devuelve exactamente un estado válido del flujo.
- **AC-02:** un proyecto sin finalizar, antes de iniciar su primer Sprint, se consulta como `Planificado`.
- **AC-03:** después de crear el primer Sprint, se consulta como `En progreso`.
- **AC-04:** después de la finalización manual confirmada de US-02, se consulta como `Finalizado`.
- **AC-05:** consultar un proyecto inexistente devuelve `404`; un error de persistencia devuelve `500`.

## Trazabilidad de aceptación

| Criterio | Escenario BDD | Prueba automatizada |
| --- | --- | --- |
| AC-01, AC-02 | BDD-01 | `TestUS04ProjectStatusCanBeConsultedThroughoutLifecycle` — alta |
| AC-01, AC-03 | BDD-02 | `TestUS04ProjectStatusCanBeConsultedThroughoutLifecycle` — primer Sprint |
| AC-01, AC-04 | BDD-03 | `TestUS04ProjectStatusCanBeConsultedThroughoutLifecycle` — finalización |
| AC-05 | BDD-04 | Cobertura existente en `TestProjectHTTPValidationAndMissingProject` y `TestProjectHTTPDoesNotReportPersistenceFailureAsSuccess`; esta última comprueba fallo de escritura y se complementa con la cobertura de lectura existente del servicio. |
