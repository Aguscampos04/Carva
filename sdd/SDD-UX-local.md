# UX — Espacio de trabajo por proyecto

Issue: [#53](https://github.com/Aguscampos04/Carva/issues/53).
Estado: implementación para revisión; integración y aceptación final pendientes.

## Objetivo
Organizar Carva por proyecto en Resumen, Backlog, Sprints, Equipo y Datos del proyecto. Mantener las funcionalidades existentes y añadir el listado de proyectos y el responsable opcional por HU solicitados por el usuario.

## Entradas y salidas
La URL hash identifica proyecto, sección e historia opcional. Recarga y Atrás/Adelante recuperan el contexto. Inicio muestra primero Abrir por identificador y Crear proyecto; debajo, Proyectos creados con acceso directo. Backlog permite búsqueda y filtros, muestra completadas al final y abre Información, Poker, Esfuerzo y Defectos bajo demanda. Resumen distingue Proyecto/Sprint y datos insuficientes de cero.

Cada Sprint puede abrirse y plegarse. Tiene En proceso y Finalizadas según acceptance_criteria_verified, sin nuevos estados de negocio. Cada HU plegada muestra solo el título; abierta muestra sus detalles y los criterios como lista visible sin un segundo desplegable. Completar exige verificación explícita y mueve la HU a Finalizadas conservando la pestaña actual. Guardar conserva despliegues y posición de lectura dentro de los límites del contenido.

Miembro responsable presenta etiqueta, selector y Guardar responsable en una fila compacta. En Sprint abierto permite un miembro opcional del proyecto. Se conserva entre Sprints. En historial se consulta el responsable actual, no una atribución histórica. No determina quién registra horas.

Cargar horas trabajadas abre dentro de la HU un formulario de alta con miembro, fecha, actividad y horas. HU y Sprint compatible se preseleccionan desde el contexto. Guardar confirma y vuelve al botón de carga. Cerrar permite reabrir; solicita confirmación si hay cambios y conserva el borrador si se cancela el descarte. Durante el guardado se bloquea el cierre y el doble envío. La consulta, comparación, edición y eliminación completas siguen en Backlog → HU → Esfuerzo.

## Contratos nuevos
- GET /api/projects: 200 con lista, vacía como []; orden por nombre COLLATE NOCASE e identificador. Devuelve 500 seguro ante error de almacenamiento. No requiere migración.
- GET /api/projects/{projectID}/backlog/{itemID}/assignee: 200 con {member_id: UUID|null}.
- PUT en la misma ruta exige únicamente member_id; null desasigna. UUID inválido o miembro ajeno: 400. Historia inexistente o ajena: 404. Error interno: 500 seguro.
- Migración aditiva 000014_story_assignees.sql: una asignación por HU, referencias a backlog_items y members. La pertenencia al proyecto se valida dentro de una transacción. No se asignan responsables por defecto.

## Reglas y restricciones
Mantener Go, REST, SQLite, HTML, CSS y JavaScript, sin frameworks nuevos. Reutilizar las reglas de US-01 a US-23: estimaciones, snapshots, completitud, traslado y cierre se resuelven en Go. No modificar fórmulas ni inventar datos históricos. Completar no genera horas ni exige cargarlas. US-24 y autenticación quedan fuera de alcance.

## Errores y casos límite
Las respuestas tardías y widgets descartados no pueden escribir sobre otro contexto. Los borradores requieren confirmación antes de descartarse; las escrituras bloquean navegación y envíos repetidos. Si una escritura se confirma y la actualización posterior falla, informar el guardado y evitar repetición automática. Distinguir cero, falta de datos, vacío y error. Desambiguar nombres iguales mediante identificadores. En una lista vacía o sin miembros, explicar cómo continuar. Una URL inválida muestra Inicio con explicación.

## Criterios y trazabilidad
Los identificadores UX-01 a UX-11 corresponden a los Acceptance Criteria de #53:

| Criterio | Validación y artefactos |
| --- | --- |
| UX-01 Inicio y proyectos | tests/integration/project_list_test.go; tests/web/workspace.test.cjs |
| UX-02 Navegación y respuestas tardías | tests/web/workspace.test.cjs; tests/web/projects.test.cjs |
| UX-03 Backlog y filtros | tests/web/workspace.test.cjs y regresión de módulos web |
| UX-04 Plegado y criterios visibles | tests/web/sprints-workspace.test.cjs; revisión visual |
| UX-05 Completitud y pestañas | tests/web/sprints-workspace.test.cjs; reglas existentes US-09 |
| UX-06 Contexto al guardar | tests/web/sprints-workspace.test.cjs; revisión en navegador |
| UX-07 Responsable | tests/integration/assignees_test.go; tests/web/sprints-workspace.test.cjs |
| UX-08 Horas contextuales y cierre | tests/web/effort.test.cjs; tests/web/sprints-workspace.test.cjs |
| UX-09 Borradores, errores y doble envío | tests/web/workspace.test.cjs; tests/web/effort.test.cjs |
| UX-10 Adaptación, teclado y foco | revisión visual y funcional; no equivale a auditoría completa de accesibilidad |
| UX-11 Métricas y regresión | tests/web/dashboard.test.cjs y tests Go existentes |

Enunciado §2.1–§2.8 y decisiones explícitas del usuario → US relacionadas en #53 → este SDD → bdd/BDD-UX-local.feature → tests → web/static y servicios Go. Los nombres locales de SDD/BDD se conservan para no romper referencias. Los tests existen; su agrupación en commits no representa una secuencia TDD histórica. La revisión del equipo sigue pendiente.
