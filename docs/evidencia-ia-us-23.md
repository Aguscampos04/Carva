# Evidencia de implementación — US-23

**Fecha:** 2026-10-27. **Issue:** [#23](https://github.com/Aguscampos04/Carva/issues/23).
**Rama:** `codex/us-23-dashboard`. **Base:** `main` en `3c5a6b1`.

## Solicitud, decisiones y responsabilidad

El usuario pidió actualizar el repositorio e iniciar US-23 después del avance
de sus compañeros. Se actualizó main mediante fast-forward, preservando la rama
US-19. Se verificaron las dependencias #4/#20/#21/#22, cerradas con DoD marcada
e integradas mediante PR #38/#49/#50/#51. Se consultaron sus contratos y pruebas.

Se prepararon [SDD](../sdd/SDD-US-23-dashboard.md) y
[BDD](../bdd/BDD-US-23-dashboard.feature) como propuesta. @Aguscampos04 aprobó
P-01/P-02 y PT-01…PT-03 con «apruebo» antes de comenzar la implementación.
La IA generó especificación, pruebas, interfaz y documentación y ejecutó las
verificaciones aquí descritas. La revisión humana del resultado y confirmación
final de DoD quedan pendientes; la aprobación de diseño no las sustituye.

## Resultado

Dashboard integrado en la página del proyecto, con nombre, UUID y el estado
recibido de la misma selección que US-04. Cinco grupos de gráficos de barras
SVG: puntos, porcentaje de avance, horas, desviación y defectos. Valores y
etiquetas accesibles también como texto. El selector consulta un Sprint y muestra
su objetivo; el botón Actualizar métricas reintenta y conserva el Sprint elegido.

No disponible no tiene barra ni se transforma en cero. El cero válido se presenta
como cero con barra de ancho cero. La desviación negativa conserva el signo y se
dibuja a la izquierda del origen. Los campos de horas preservan el valor textual
de la respuesta al formatearse; la conversión numérica se usa para geometría.

Se reutilizan las APIs y cálculos Go de US-20/21/22 sin sumar métricas en JavaScript.
No hubo cambios Go, tablas, migraciones, dependencias, prioridades, estados de
Issues ni funcionalidades de US-24. Los datos de selecciones previas se limpian
al iniciar una carga y las respuestas tardías no reemplazan el contexto vigente.

## Evidencia TDD ejecutada

| Etapa | Evidencia |
| --- | --- |
| Especificación | `0a74747`, SDD y BDD aprobados antes del código |
| RED gráficos | `43454a5`: `node --test tests/web/dashboard.test.cjs`, tres fallos por dashboard inexistente y datos del Sprint anterior sin limpiar |
| GREEN gráficos | `e091a82`: siete tests de dashboard/métricas PASS |
| RED integración | `3f37ce1`: test de app falla por configuración de dashboard inexistente |
| GREEN integración | `4b561b0`: suite web completa PASS, 43 tests |
| RED reintento | `7414512`: falla porque Actualizar métricas no vuelve a cargar el Sprint seleccionado |
| GREEN y refactor | `481aa39`: extracción de loadSprint, limpieza compartida y protección del contexto; 13 tests focalizados PASS |

Durante el refactor, una eliminación de limpieza en el lugar incorrecto hizo
fallar el test que comprueba horas/defectos del Sprint anterior. Se restauró la
limpieza en `clearSprint`; las aserciones se conservaron y volvieron a pasar.
No se presenta una ejecución fallida como exitosa.

## Trazabilidad de criterios y escenarios

| Criterio / escenario | Prueba | Resultado técnico |
| --- | --- | --- |
| AC-01 / US23-01 | dashboard.test + integración app en projects.test | PASS: estado oficial, incluidas actualización y finalización |
| AC-01 / US23-02 | dashboard.test | PASS: mismos valores, cinco grupos y Sprint con objetivo |
| AC-02 / US23-03 | dashboard.test | PASS: null sin barra; cero válido con ancho cero |
| AC-01 / US23-04 | dashboard.test | PASS: desviación negativa, posición anterior al origen |
| US23-05 | dashboard.test | PASS: limpieza completa, cambio de contexto y respuestas tardías |
| US23-06 | dashboard.test | PASS: sin Sprints y sin proyecto seleccionado |
| US23-07 | dashboard.test | PASS: error visible, reintento y conservación de selección |

Los escenarios Gherkin se cubren mediante Go/Node; no se añadió runner Gherkin.
Las reglas de negocio permanecen en Go. La suite Go existente verifica las fuentes,
persistencia SQLite, snapshots e historial; no se crearon cálculos duplicados.

## Validaciones ejecutadas

| Comando | Resultado |
| --- | --- |
| `go test ./... -count=1` | PASS, todos los paquetes e integración |
| `node --test tests/web/*.test.cjs` | PASS, 44 tests |
| `gofmt -l .` | PASS, no reportó archivos |
| `go vet ./...` | PASS, salida 0 |
| `go build ./...` | PASS, salida 0 |
| `go mod verify` | PASS, all modules verified |
| `git diff --check` | PASS |

No se calculó cobertura Go nueva: US-23 reutiliza Go sin modificarlo.
No se ejecutó detector de carreras en esta tarea; no se declara aprobado.

## Prueba en navegador con aplicación real

Se inició Carva con SQLite temporal independiente mediante CARVA_DB_PATH.
Se creó un proyecto desde la UI: estado Planificado, puntos/avance no disponibles,
horas y defectos cero según las fuentes. Se prepararon datos por REST en esa base:
inicio 2026-10-27, primer Sprint, historia estimada en 5,75 h, miembro y esfuerzo
de 3,75 h y un defecto resuelto vinculado al Sprint.

Un primer intento usó una ruta incorrecta del fixture (`/backlog` en vez de
`/backlog-items` del Sprint) y devolvió 404. Se corrigió la ruta del fixture;
no fue un fallo de producto ni requirió modificar código.

La API de horas confirmó 5,75/3,75/−2. En navegador, Consultar de nuevo mostró
En progreso, esos mismos valores y un defecto detectado/resuelto. Seleccionar
Sprint 1 mostró el objetivo y cinco grupos de gráficos. Se inspeccionaron
visualmente las barras y la desviación negativa. Captura local ignorada por Git:
`.local/dashboard-us23.png`. No se usó la base habitual del usuario.

## Revisión manual recomendada

1. Iniciar Carva con una base de prueba y crear/consultar un proyecto.
2. Verificar que estado del dashboard y estado del proyecto coincidan.
3. Comparar los valores de gráficos con las consultas numéricas de métricas.
4. Seleccionar un Sprint y comprobar su objetivo y métricas.
5. Verificar null frente a cero y desviaciones positivas/negativas.
6. Usar Actualizar métricas y comprobar que conserva el Sprint seleccionado.

Pendientes: revisión humana, confirmación DoD, autorización/publicación y PR.
El PR deberá usar `Closes #23` cuando la DoD esté confirmada; merge exclusivamente
manual por el responsable. No se cerró ni actualizó la Issue desde esta tarea.
