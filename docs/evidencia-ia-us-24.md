# Evidencia de implementación — US-24

**Fecha:** 2026-10-30. **Issue:** [#24](https://github.com/Aguscampos04/Carva/issues/24).
**Rama:** `codex/us-24-reportes-pdf`. **Base:** `main` en `83eb916`, incluye UX #53/PR #54.
**Responsable GitHub:** @Loana-Vargas; la asignación no se modificó.

## Solicitud y revisión humana

El Agile Enabler pidió «hace la issue 24». Se verificó que no hubiera código,
especificación o PR equivalente y se preparó una propuesta SDD/BDD. Aprobó
P-01…P-04 y PT-01…PT-04 mediante «sisi», incluida la biblioteca PDF.
La IA escribió especificación, tests, implementación y documentos; ejecutó las
verificaciones descritas. Esta aprobación de diseño no equivale a revisión
humana final de código, aceptación funcional o DoD del equipo.

Dependencias #1/#5/#7/#8/#9/#16/#18/#19/#20/#21/#22: cerradas, ocho ítems DoD
marcados por Issue e implementaciones integradas. No se modificaron Issues,
milestones, responsables, prioridades ni Story Points durante US-24.

## Funcionalidad y decisiones implementadas

- Reporte de proyecto completo y reporte de Sprint individual en JSON y PDF.
- Sección Reportes en la navegación del proyecto, selector de alcance,
  reconsulta y descarga real de archivo. La UI usa secciones textuales generadas
  desde Go y no implementa cálculos de negocio.
- Datos generales y estado; miembros; Backlog actual para proyecto; snapshots,
  Goal y resultados de cada Sprint; esfuerzo activo; métricas y defectos.
- Uso de métricas existentes: no hay fórmulas alternativas. Datos insuficientes
  como No disponible; ceros válidos se conservan.
- Sprint: esfuerzo asociado explícitamente; defectos cuya detección o resolución
  corresponde al Sprint, sin duplicar el mismo registro. Su lista no sustituye
  las reglas de conteo de US-22.
- Historial desde snapshots, con la excepción de horas estimadas vigentes de
  US-16 identificada en el reporte. Trasladar/completar después no cambia el
  resultado registrado en el origen.
- PDF paginado con acentos, textos extensos y datos exactos de horas. Fuente
  embebida Go Regular; caracteres no cubiertos se representan como `[U+XXXX]`,
  en lugar de desaparecer. JSON/UI conservan el texto original.
- Sin tablas, migraciones, servidor PDF externo, navegador headless ni CGO.

Dependencias aprobadas y fijadas: `github.com/signintech/gopdf v0.38.3` y
`golang.org/x/image v0.47.0` para la fuente. Se agregaron sus transitivas
gofpdi y pkg/errors; x/image requiere actualizar x/sys de v0.48.0 a v0.49.0.
Go mínimo 1.26.0 y SQLite modernc v1.60.1 se mantienen. Licencias PDF/fuente
conservadas en `docs/licencias`.

## Trazabilidad y TDD real

Enunciado §2.9/entrega → EP-09 → #24 →
[SDD](../sdd/SDD-US-24-reportes-pdf.md) → AC →
[BDD](../bdd/BDD-US-24-reportes-pdf.feature) → Go/Node/SQLite → interfaz y PDF.

| Etapa | Evidencia ejecutada |
| --- | --- |
| Especificación | `fa6de45`: SDD y BDD aprobados antes del código. |
| RED aplicación | `e86d7d2`: `go test ./internal/application/reports -count=1` falló porque el paquete aún no tenía implementación. |
| GREEN aplicación | `c9ecae0`: el mismo paquete y SQLite compilaron/pasaron; scopes, errores y snapshots comprobados. |
| RED API | `25eefbd`: `go test ./tests/integration -run TestUS24 -count=1` falló por WithReports inexistente. |
| RED PDF | `5e80985`: paquete reportpdf sin código de producción; prueba de secciones/archivo/paginación aún no podía ejecutarse. |
| GREEN API/PDF | `247b5b8`: tests focalizados de aplicación, PDF y HTTP/SQLite PASS. |
| RED web | `f90a934`: cuatro tests fallaron por reports.js inexistente. |
| GREEN web | `438cef4`: suite completa web de ese momento PASS, 66 tests. |
| Regresión de historial/esfuerzo | `f20e0eb`: verificación de campos originales, completitud posterior, registros activos, eliminados y sin Sprint. |
| RED caracteres | `ca007d6`: safeText ausente. Un primer intento usando sfnt también informó una dependencia faltante; se evitó introducirla usando el parser ya incluido en gopdf. |
| GREEN caracteres | `cc0747b`: acentos preservados, controles normalizados y escape explícito de glifos ausentes; PDF/API PASS. |
| Refactor | `c368a08`: reutilización de lectores de esfuerzo/defectos por reporte y verificación de navegación; focalizados Go y web PASS. |

La prueba de regresión de snapshots tuvo un primer fixture inválido que intentó
insertar `planned_status` como columna SQLite. Ese campo es una proyección Go
del estado planificado. Se corrigió únicamente el fixture para usar el esquema
real y se ejecutó otra vez: PASS. No se presenta ese fallo como defecto del producto.

## Criterios de aceptación

| Criterio | Evidencia | Resultado técnico |
| --- | --- | --- |
| AC-01: reporte de proyecto y Sprint con datos de su alcance | `TestUS24ReportsMatchSourcesAndPreserveScope`, `TestSprintScopePreservesSnapshotAndFiltersRecords`, navegación web | PASS |
| AC-02: PDF final completo | `TestPDFPreservesSectionsUnavailableAndAccentedText`, endpoints project.pdf/Sprint.pdf; extracción pypdf y revisión visual de cinco páginas | PASS |
| AC-03: concordancia y No disponible | comparación con métricas fuente; `TestUS24HistoricalPlanningAndEffortAreNotRebuiltOrDoubleCounted`; datos null frente a cero | PASS |
| Errores/ámbito | `TestUS24MissingInvalidAndUnavailableSources`: inválidos, ajenos, inexistentes, base cerrada; errores sin falso PDF | PASS |
| Texto y Unicode | `TestPDFTextPreservesAccentsAndEscapesUnsupportedGlyphs`; PDF real con áéíóú/ñ | PASS |
| UI | `tests/web/reports.test.cjs`: consulta segura, descarga, reintento, error PDF, envíos repetidos, navegación y respuestas tardías | PASS |

BDD-01/02/04/05/07 cubiertos por aplicación/integración; BDD-03/06 por PDF y
prueba real; BDD-08 por Node. No se agregó un runner Gherkin: se automatizaron
los comportamientos relacionados mediante Go y Node.

## Validación técnica final

| Comando ejecutado | Resultado |
| --- | --- |
| `go test ./... -count=1` | PASS; todos los paquetes, incluida integración |
| `node --test tests/web/*.test.cjs` | PASS; 67 tests |
| `gofmt -l .` | PASS; no reporta archivos |
| `go vet ./...` | PASS, salida 0 |
| `go build ./...` | PASS, salida 0 |
| `go mod verify` | PASS, all modules verified |
| `git diff --check` | PASS |
| `go test -race ./...` | NO DISPONIBLE; CGO deshabilitado, salida 2. No se afirma PASS. |

Logs locales ignorados por Git: `.local/us24-go-final.log` y
`.local/us24-web-final.log`. Cobertura focalizada con cuatro paquetes
instrumentados y selección `TestReport|TestSprintScope|TestPDF|TestUS24`:

- Project 87,1%; Sprint 83,3%; ForSprint 100%.
- ReportSource.Load 78,4%; WithReports 95,5%.
- Render 75%; Sections y safeText 100%.
- Total combinado del perfil: 36,8%. Incluye código HTTP/SQLite anterior que
  estos tests focalizados no ejecutan: **no es cobertura global del producto**.
- Perfil: `.local/us24.cover`, generado realmente y analizado con go tool cover.

## Prueba real de aplicación, navegador y PDF

Se inició Carva con CARVA_DB_PATH temporal independiente de la base habitual.
Por API se crearon proyecto, dos Sprints, historia con texto extenso y acentos,
miembro, esfuerzo de 3,75 h sobre 5,75 estimadas y un defecto detectado/resuelto
en Sprints diferentes. Se cerró el primer Sprint y se completó la historia en
el segundo.

- Reporte de proyecto: PDF de cinco páginas, métricas 5,75/3,75/−2, datos generales,
  estado, miembros, planificación/resultados, esfuerzo y defecto.
- Reporte de Sprint 1: PDF de dos páginas; extracción confirma `resultado: Incompleta`
  aunque la historia actual ya esté completada en Sprint 2.
- pypdf verificó contenido y acentos. Poppler renderizó las cinco páginas del
  proyecto; se inspeccionaron visualmente sin solapamientos ni texto recortado.
- En navegador se abrió la ruta Reportes, se eligió Sprint 1 y se pulsó Descargar
  PDF. La UI confirmó la descarga y el archivo real quedó en
  `C:/Users/agust/Downloads/carva-sprint-1.pdf`; su texto y número de páginas
  volvieron a comprobarse con pypdf.
- Artefactos locales ignorados: `.local/us24-project.pdf`, `.local/us24-sprint.pdf`,
  `.local/us24-reportes.png` y páginas PNG de comprobación.

Estas pruebas las ejecutó la IA. No se presentan como revisión visual humana.

## Límites y entrega

Las lecturas reutilizadas no garantizan una instantánea atómica de todas las
fuentes ante escrituras concurrentes. El PDF no certifica aceptación ni cierra
el proyecto. Los caracteres sin glifo usan representación Unicode explícita.
No se ejecutó detector de carreras por la limitación indicada.

Implementación local y criterios técnicamente comprobados. Pendientes: revisión
humana, confirmación final de DoD y autorización para push/PR. No se realizó
merge, cierre ni actualización de la Issue. El merge queda siempre a cargo del
responsable. Al aprobar la entrega, el PR deberá incluir `Closes #24`.

Verificación manual: iniciar Carva con base de prueba, abrir proyecto → Reportes,
elegir alcance, comparar valores con sus fuentes y descargar ambos PDF. Verificar
datos insuficientes, historia trasladada, lista vacía y error/reintento.
