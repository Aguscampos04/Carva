# SDD-US-05 — Registrar y mantener elementos del Product Backlog

**Estado:** Especificación aprobada por DelfiCarribero el 2026-10-15.  
**Trazabilidad:** enunciado §2.2 → EP-02 → US-05 → Issue #5 → este SDD → BDD-US-05 → tests y código.

## 1. Objetivo

Permitir registrar una User Story en el Backlog de un proyecto y actualizar sus datos editables, sin alterar su estado ni registrar una estimación acordada.

## 2. Datos y validaciones aprobados

| Campo | Regla |
| --- | --- |
| `id` | UUID generado por el servidor; no se recibe del cliente. |
| `project_id` | Proyecto existente al que pertenece el elemento. |
| `title` | Requerido; se rechaza si queda vacío al recortar espacios exteriores para validar. Se conserva el texto original y no hay máximo. |
| `description` | Requerida; misma regla de texto no vacío que el título, sin máximo. |
| `priority` | Requerida; valores admitidos: `Alta`, `Media`, `Baja`, en ese orden. |
| `status` | Generado por el servidor como `Pendiente` al crear. Estados del dominio: `Pendiente`, `En progreso`, `Completada`. US-05 no modifica el estado. |
| `provisional_story_points` | Opcional al crear; si se informa, entero no negativo. Se puede editar mediante US-05 hasta que US-15 registre la estimación definitiva. |
| `story_points` | Estimación definitiva, separada del valor provisional. US-05 no la establece ni modifica; corresponde a US-15. |
| `acceptance_criteria` | Lista requerida con al menos un texto no vacío. Los espacios exteriores se recortan para validar; se conserva el texto ingresado. |

El identificador sigue la convención aprobada en US-01. Los textos requeridos se validan en Go, no solamente en la interfaz.

## 3. Entradas y salidas

El alta recibe título, descripción, prioridad, criterios de aceptación y, opcionalmente, Story Points provisionales. El servidor agrega el UUID, el proyecto y el estado inicial. La respuesta contiene todos los campos persistidos.

La actualización recibe el conjunto completo de campos editables —título, descripción, prioridad, valor provisional y criterios— para el elemento identificado dentro del proyecto. No recibe `id`, `project_id`, `status` ni `story_points` definitivo. Devuelve el elemento vigente.

## 4. Reglas de negocio

- **BR-01:** cada elemento pertenece a un proyecto existente.
- **BR-02:** el estado inicial es `Pendiente`.
- **BR-03:** los únicos estados del Backlog son `Pendiente`, `En progreso` y `Completada`.
- **BR-04:** US-05 no cambia el estado. La transición a `Completada` corresponde a US-09, tras verificar criterios de aceptación.
- **BR-05:** US-05 nunca registra ni cambia Story Points definitivos; eso corresponde a US-15.
- **BR-06:** el valor provisional puede modificarse hasta que US-15 registre el definitivo y se mantiene separado de este.
- **BR-07:** al actualizar un elemento, los campos no editables y su pertenencia al proyecto se conservan.
- **BR-08:** los cambios válidos se guardan sin escrituras parciales; errores no producen éxito aparente.

## 5. Contrato HTTP propuesto para revisión

La implementación reutilizaría el estilo REST JSON de US-01:

- `POST /api/projects/{projectID}/backlog`: crea un elemento; devuelve `201` y el objeto creado.
- `PUT /api/projects/{projectID}/backlog/{itemID}`: reemplaza los campos editables; devuelve `200` y el objeto vigente.
- Cuerpo JSON: `title`, `description`, `priority`, `provisional_story_points` (opcional) y `acceptance_criteria` (lista de textos).
- `status`, `id`, `project_id` y `story_points` definitivo no forman parte del contrato de escritura.
- `400` para JSON o datos inválidos; `404` para proyecto o elemento no encontrado; `500` para fallos de almacenamiento.

El contrato fue aprobado por DelfiCarribero el 2026-10-15. US-06 especifica la consulta de lista.

## 6. Casos límite y errores

- No existe el proyecto: `404`; no se crea el elemento.
- No existe el elemento en ese proyecto: `404`; no se modifica otro elemento.
- Título, descripción o criterio requerido vacío/blanco, prioridad desconocida o Story Points provisional negativo/no entero: `400`; no hay cambios parciales.
- Intento de enviar un estado `Completada` mediante escritura genérica: no completa el elemento; el estado se mantiene bajo el flujo de US-09.
- Falla de SQLite: `500`, sin filtrar detalles SQL ni devolver éxito.
- Identificadores diferentes representan elementos independientes.

## 7. Criterios de aceptación

- **AC-01:** Registrar elemento con los campos requeridos y admitir Story Points provisionales opcionales; asignar ID y estado inicial, y persistirlo.
- **AC-02:** Actualizar los campos editables y consultar luego sus valores sin cambiar una estimación definitiva.
- **AC-03:** Una edición genérica no lleva el elemento a `Completada` sin la verificación de US-09.

## 8. Fuera de alcance

Consulta/listado (US-06), asignación a Sprint (US-08), finalización (US-09), Planning Poker o consenso (US-15), autenticación y gestión de proyectos.

## 9. Decisiones

- Reglas funcionales de campos y validaciones aprobadas por DelfiCarribero el 2026-10-15.
- Contrato HTTP aprobado por DelfiCarribero el 2026-10-15.
