# Borradores de Issues para las User Stories US-01 a US-24

**Estado:** borrador local para revisión del equipo. Estos textos no son Issues creadas en GitHub ni representan trabajo implementado.

**Fuentes:** `docs/enunciado.md`, `docs/contexto-proyecto.md`, las nueve Épicas y las 24 User Stories aprobadas como versión inicial del Product Backlog, y las decisiones funcionales aprobadas que allí se documentan.

**Prioridad:** MVP identifica el conjunto candidato del Sprint 1 (US-01, US-03, US-04, US-05, US-06, US-07, US-08 y US-09); importante y secundaria ordenan el trabajo posterior. Todas las funciones exigidas por el enunciado siguen siendo obligatorias para la entrega final. No se asignan Story Points a estas Issues antes de que el equipo las estime.

## Etiquetas sugeridas

Usar una etiqueta de cada grupo, si el equipo aprueba crearlas: `tipo:user-story`; `epic:EP-01` a `epic:EP-09`; `prioridad:mvp`, `prioridad:importante` o `prioridad:secundaria`. Son nombres propuestos, no etiquetas ya existentes en GitHub.

## Definition of Done común

Para cada borrador, la **DoD aplicable** es esta lista más la comprobación particular indicada en su sección. Se aplicará cuando la Issue se ejecute, no a estos borradores:

- [ ] Los Acceptance Criteria de la Issue se verificaron y las decisiones funcionales que afectan su comportamiento están documentadas.
- [ ] La especificación SDD correspondiente está versionada antes de implementar; cubre entradas, salidas, reglas, restricciones, errores y casos límite pertinentes.
- [ ] Los Acceptance Criteria relevantes tienen escenarios BDD Given–When–Then; se automatizaron cuando es viable.
- [ ] El trabajo de reglas, cálculos y validaciones siguió RED → GREEN → REFACTOR con evidencia real; existen tests automatizados pertinentes, incluidos casos de error y límite.
- [ ] El núcleo y las reglas de negocio correspondientes están en Go; se revisó la integración de la funcionalidad y la documentación.
- [ ] Cuando haya código Go, `gofmt -l .`, `go test ./...` y `go vet ./...` se ejecutaron con resultados revisados; se reportó la cobertura cuando corresponda.
- [ ] Queda trazabilidad entre requisito, Épica, User Story, SDD, Acceptance Criteria, BDD, TDD, código Go y tests; el equipo revisó y validó el resultado.
- [ ] Cuando se publique el trabajo, commit, push y eventual Pull Request/Review se registran como evidencia real. No se cierra la Issue solo porque el código compile.

En cada sección, **SDD, BDD, TDD, Go y tests** son artefactos previstos, todavía no creados. Las dependencias se expresan mediante IDs de User Story porque aún no hay números de Issue.

## EP-01 — Gestión de proyectos y equipo

### US-01 — Crear proyecto

**ID / Épica:** US-01 / EP-01. **Prioridad:** MVP. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-01`, `prioridad:mvp`.

**Como / Quiero / Para:** Como integrante del equipo, quiero crear un proyecto para organizar su trabajo.

**Contexto y alcance:** Alta de proyecto, según el §2.1 del enunciado. El estado inicial es **Planificado**; el cambio al crear el primer Sprint pertenece a US-07.

**Acceptance Criteria:**
- Dado un alta de proyecto válida, cuando se registra, entonces el proyecto puede consultarse con los datos ingresados.
- Dado un proyecto recién creado y sin Sprint, cuando se consulta su estado, entonces figura **Planificado**.

**Dependencias:** ninguna.

**Definition of Done aplicable:** DoD común; demostrar el alta y el estado inicial **Planificado** con pruebas y trazabilidad a §2.1.

**Trazabilidad prevista:** §2.1 → EP-01 → US-01 → SDD de US-01 → BDD de US-01 → TDD de US-01 → Go → tests automatizados.

### US-02 — Modificar proyecto y registrar fechas

**ID / Épica:** US-02 / EP-01. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-01`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero modificar un proyecto y registrar sus fechas de inicio y finalización para mantener sus datos.

**Contexto y alcance:** Modificación de proyecto y registro de fechas exigidos por §2.1. Al terminar el proyecto, su estado pasa a **Finalizado**, según la decisión del equipo. La forma concreta de declarar que el proyecto terminó debe precisarse en el SDD sin presumir que basta una fecha planificada.

**Acceptance Criteria:**
- Dado un proyecto existente, cuando se modifican sus datos o fechas, entonces los valores registrados se muestran al consultarlo de nuevo.
- Dado que el proyecto ha terminado conforme a la regla definida en el SDD, cuando se consulta su estado, entonces figura **Finalizado**.

**Dependencias:** US-01.

**Definition of Done aplicable:** DoD común; probar modificación, persistencia de fechas y transición a **Finalizado** conforme al SDD aprobado.

**Trazabilidad prevista:** §2.1 → EP-01 → US-02 → SDD de US-02 → BDD de US-02 → TDD de US-02 → Go → tests automatizados.

### US-03 — Registrar e identificar miembros del proyecto

**ID / Épica:** US-03 / EP-01. **Prioridad:** MVP. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-01`, `prioridad:mvp`.

**Como / Quiero / Para:** Como integrante del equipo, quiero registrar e identificar miembros del proyecto para atribuirles votos de Planning Poker y registros de esfuerzo.

**Contexto y alcance:** Registro de miembros exigido por §2.1 e identificación necesaria para las funciones aprobadas. El MVP no incluye un sistema complejo de permisos; esta Issue no presupone una tecnología de autenticación.

**Acceptance Criteria:**
- Dado un proyecto, cuando se registra un miembro, entonces puede identificarse entre los miembros de ese proyecto.
- Los votos y registros de esfuerzo, cuando existan, pueden atribuirse al miembro identificado sin confundirlos con los de otro miembro.

**Dependencias:** US-01.

**Definition of Done aplicable:** DoD común; verificar que la identidad registrada permite atribución inequívoca en US-12 y US-16 cuando estén disponibles.

**Trazabilidad prevista:** §2.1 y decisión de MVP → EP-01 → US-03 → SDD de US-03 → BDD de US-03 → TDD de US-03 → Go → tests automatizados.

### US-04 — Consultar estado del proyecto

**ID / Épica:** US-04 / EP-01. **Prioridad:** MVP. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-01`, `prioridad:mvp`.

**Como / Quiero / Para:** Como integrante del equipo, quiero consultar si un proyecto está Planificado, En progreso o Finalizado para conocer su situación.

**Contexto y alcance:** Consulta del estado exigida por §2.1. Comparte el mismo estado con US-23; no implementa un segundo cálculo. El proyecto empieza **Planificado**, pasa a **En progreso** al crear el primer Sprint (US-07) y a **Finalizado** al terminarse (US-02).

**Acceptance Criteria:**
- Al consultar un proyecto se muestra exactamente uno de los tres estados aprobados: **Planificado**, **En progreso** o **Finalizado**.
- Al crear el primer Sprint de un proyecto no finalizado, su estado consultado pasa a **En progreso**.
- Al terminar el proyecto conforme al SDD de US-02, el estado consultado pasa a **Finalizado**.

**Dependencias:** US-01 para consultar un proyecto; US-07 y US-02 aportan las transiciones posteriores.

**Definition of Done aplicable:** DoD común; probar los tres estados y la concordancia de este dato con US-23 cuando se implemente el dashboard.

**Trazabilidad prevista:** §2.1 → EP-01 → US-04 → SDD de US-04 → BDD de US-04 → TDD de US-04 → Go → tests automatizados.

## EP-02 — Product Backlog

### US-05 — Registrar y mantener elementos del Product Backlog

**ID / Épica:** US-05 / EP-02. **Prioridad:** MVP. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-02`, `prioridad:mvp`.

**Como / Quiero / Para:** Como integrante del equipo, quiero registrar y actualizar elementos del Product Backlog para describir y organizar el trabajo.

**Contexto y alcance:** Cada elemento cuenta con identificador, título, descripción, prioridad, estado, Story Points y Acceptance Criteria (§2.2). Sus estados son **Pendiente**, **En progreso** y **Completada**. Puede crearse con Story Points provisionales; luego el Story Point actual cambia solo mediante US-15. El paso a **Completada** requiere US-09.

**Acceptance Criteria:**
- Al registrar un elemento, quedan disponibles todos los campos exigidos por §2.2 y se admite un valor de Story Points provisional.
- Los datos editables del Backlog pueden actualizarse y consultarse sin que ello registre una estimación acordada.
- Un elemento no pasa a **Completada** únicamente por editar su estado sin verificar sus Acceptance Criteria.

**Dependencias:** US-01.

**Definition of Done aplicable:** DoD común; probar campos exigidos, valor provisional y protección de las reglas de finalización y estimación.

**Trazabilidad prevista:** §2.2 → EP-02 → US-05 → SDD de US-05 → BDD de US-05 → TDD de US-05 → Go → tests automatizados.

### US-06 — Consultar Product Backlog

**ID / Épica:** US-06 / EP-02. **Prioridad:** MVP. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-02`, `prioridad:mvp`.

**Como / Quiero / Para:** Como integrante del equipo, quiero consultar los elementos del Product Backlog para conocer su prioridad, estado y datos antes de planificar.

**Contexto y alcance:** Lectura de los elementos creados en US-05; modificar prioridad y otros campos corresponde a US-05. La consulta permite seleccionar historias para US-08 sin duplicar datos.

**Acceptance Criteria:**
- Los elementos del Backlog se pueden consultar con sus campos exigidos por §2.2.
- La prioridad, estado y Story Points mostrados coinciden con los valores vigentes del elemento.

**Dependencias:** US-05.

**Definition of Done aplicable:** DoD común; verificar lectura de los campos exigidos y consistencia con US-05.

**Trazabilidad prevista:** §2.2 → EP-02 → US-06 → SDD de US-06 → BDD de US-06 → TDD de US-06 → Go → tests automatizados.

## EP-03 — Gestión de Sprints

### US-07 — Crear Sprint y definir Sprint Goal

**ID / Épica:** US-07 / EP-03. **Prioridad:** MVP. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-03`, `prioridad:mvp`.

**Como / Quiero / Para:** Como integrante del equipo, quiero crear un Sprint y definir su Sprint Goal para establecer el objetivo del trabajo.

**Contexto y alcance:** Alta y objetivo de Sprint (§2.3). Al crear el primer Sprint de un proyecto, el proyecto pasa de **Planificado** a **En progreso**.

**Acceptance Criteria:**
- Un Sprint creado queda asociado al proyecto y permite consultar su Sprint Goal.
- Al crear el primer Sprint de un proyecto **Planificado**, el estado del proyecto pasa a **En progreso**.

**Dependencias:** US-01.

**Definition of Done aplicable:** DoD común; probar alta, objetivo y transición de estado del proyecto.

**Trazabilidad prevista:** §2.3 y decisión de estado del proyecto → EP-03 → US-07 → SDD de US-07 → BDD de US-07 → TDD de US-07 → Go → tests automatizados.

### US-08 — Asignar historias y preservar lo planificado

**ID / Épica:** US-08 / EP-03. **Prioridad:** MVP. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-03`, `prioridad:mvp`.

**Como / Quiero / Para:** Como integrante del equipo, quiero asignar historias del Backlog a un Sprint para definir su trabajo planificado.

**Contexto y alcance:** Asignación de historias (§2.3) y preservación del conjunto planificado aunque cambie el Backlog posteriormente. US-10 reutiliza ese historial al trasladar historias incompletas.

**Acceptance Criteria:**
- Las historias asignadas aparecen en el Sprint correspondiente.
- El conjunto planificado de ese Sprint puede consultarse sin alterarse por modificaciones posteriores del Backlog.

**Dependencias:** US-05, US-07.

**Definition of Done aplicable:** DoD común; probar asignación y preservación histórica del conjunto planificado.

**Trazabilidad prevista:** §2.3 y decisión sobre métricas → EP-03 → US-08 → SDD de US-08 → BDD de US-08 → TDD de US-08 → Go → tests automatizados.

### US-09 — Registrar historias completadas

**ID / Épica:** US-09 / EP-03. **Prioridad:** MVP. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-03`, `prioridad:mvp`.

**Como / Quiero / Para:** Como integrante del equipo, quiero registrar una historia como Completada después de verificar sus Acceptance Criteria para reflejar el trabajo efectivamente terminado.

**Contexto y alcance:** Registro de historias completadas (§2.3) conforme al estado aprobado. No basta con que termine el Sprint o exista código.

**Acceptance Criteria:**
- Una historia cuyos Acceptance Criteria no se verificaron no puede pasar a **Completada**.
- Tras verificar esos criterios y registrar su finalización, la historia figura como **Completada** en el Sprint y en su Backlog.

**Dependencias:** US-05, US-08.

**Definition of Done aplicable:** DoD común; probar rechazo de finalización sin verificación y consistencia del estado entre Sprint y Backlog.

**Trazabilidad prevista:** §2.3 y decisión sobre completitud → EP-03 → US-09 → SDD de US-09 → BDD de US-09 → TDD de US-09 → Go → tests automatizados.

### US-10 — Cerrar Sprint y trasladar historias incompletas

**ID / Épica:** US-10 / EP-03. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-03`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero cerrar un Sprint y trasladar automáticamente sus historias incompletas al siguiente Sprint para continuar el trabajo sin perder el resultado anterior.

**Contexto y alcance:** Cierre (§2.3), traslado automático y conservación del historial. Si hay historias incompletas, el siguiente Sprint debe existir antes de cerrar el actual.

**Acceptance Criteria:**
- Dado un Sprint con historias incompletas y sin siguiente Sprint disponible, cuando se intenta cerrarlo, entonces el cierre no se realiza.
- Dado un siguiente Sprint disponible, cuando se cierra el Sprint actual, entonces las historias incompletas pasan automáticamente al siguiente y las completadas no se trasladan.
- El Sprint cerrado conserva su conjunto planificado y el resultado de cada historia, incluso después del traslado.
- Un Sprint sin historias incompletas puede cerrarse sin necesitar un siguiente Sprint.

**Dependencias:** US-07; US-08 cuando el Sprint tiene historias asignadas.

**Definition of Done aplicable:** DoD común; probar ambos casos de cierre, bloqueo sin siguiente Sprint y preservación del historial.

**Trazabilidad prevista:** §2.3 y decisión de traslado → EP-03 → US-10 → SDD de US-10 → BDD de US-10 → TDD de US-10 → Go → tests automatizados.

### US-11 — Consultar Sprints anteriores

**ID / Épica:** US-11 / EP-03. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-03`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero consultar Sprints cerrados para revisar sus objetivos y resultados.

**Contexto y alcance:** Consulta histórica exigida por §2.3, con objetivo, historias planificadas y finalización registrada. No modifica ni vuelve a calcular el plan histórico.

**Acceptance Criteria:**
- Los Sprints cerrados pueden consultarse con su Sprint Goal, historias planificadas y resultado registrado.
- El traslado de una historia a otro Sprint no cambia el registro histórico del Sprint cerrado.

**Dependencias:** US-10.

**Definition of Done aplicable:** DoD común; probar la consulta de un Sprint anterior tras trasladar una historia incompleta.

**Trazabilidad prevista:** §2.3 → EP-03 → US-11 → SDD de US-11 → BDD de US-11 → TDD de US-11 → Go → tests automatizados.

## EP-04 — Estimación y Planning Poker

### US-12 — Registrar votos individuales ocultos

**ID / Épica:** US-12 / EP-04. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-04`, `prioridad:importante`.

**Como / Quiero / Para:** Como miembro identificado, quiero registrar mi estimación individual sin ver los votos ajenos durante la votación para estimar de forma independiente.

**Contexto y alcance:** Votos por historia y ronda (§2.4), ocultos hasta finalizar la votación. Los participantes de cada ronda deben poder distinguirse para determinar consenso.

**Acceptance Criteria:**
- Cada voto registrado se atribuye al miembro, historia y ronda correspondientes.
- Antes de finalizar una ronda, sus votos individuales no son visibles para los otros participantes.

**Dependencias:** US-03, US-05.

**Definition of Done aplicable:** DoD común; probar atribución individual y ocultamiento antes de finalizar la ronda.

**Trazabilidad prevista:** §2.4 → EP-04 → US-12 → SDD de US-12 → BDD de US-12 → TDD de US-12 → Go → tests automatizados.

### US-13 — Revelar votos y detectar diferencias

**ID / Épica:** US-13 / EP-04. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-04`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero revelar simultáneamente los votos al finalizar la ronda y conocer sus diferencias para discutir la estimación.

**Contexto y alcance:** Revelación conjunta y detección de diferencias (§2.4). Existe diferencia cuando los valores votados no son todos iguales.

**Acceptance Criteria:**
- Antes de finalizar la ronda, los votos individuales de otros participantes continúan ocultos; al finalizarla se muestran simultáneamente.
- Si todos los valores votados son iguales, no se indica diferencia; si al menos uno difiere, se indica que existe diferencia.

**Dependencias:** US-12.

**Definition of Done aplicable:** DoD común; probar ocultamiento, revelación simultánea e igualdad/desigualdad de votos.

**Trazabilidad prevista:** §2.4 → EP-04 → US-13 → SDD de US-13 → BDD de US-13 → TDD de US-13 → Go → tests automatizados.

### US-14 — Discutir diferencias y realizar nuevas rondas

**ID / Épica:** US-14 / EP-04. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-04`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero revisar las diferencias reveladas y realizar otra ronda cuando sea necesario para reconsiderar la estimación.

**Contexto y alcance:** Los votos revelados quedan disponibles para discusión del equipo; la función requerida es iniciar nuevas rondas, no incorporar mensajería. Cada nueva ronda vuelve a tener votos individuales ocultos hasta su revelación.

**Acceptance Criteria:**
- Después de revelar una ronda, sus valores y la existencia de diferencias pueden consultarse para discutirlos.
- Puede iniciarse otra ronda para la misma historia; los votos de esa nueva ronda permanecen ocultos hasta finalizarla.

**Dependencias:** US-13.

**Definition of Done aplicable:** DoD común; probar una segunda ronda y que sus votos no se revelan prematuramente.

**Trazabilidad prevista:** §2.4 → EP-04 → US-14 → SDD de US-14 → BDD de US-14 → TDD de US-14 → Go → tests automatizados.

### US-15 — Registrar estimación acordada

**ID / Épica:** US-15 / EP-04. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-04`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero registrar el valor consensuado para establecer el Story Point definitivo de la historia.

**Contexto y alcance:** Hay consenso cuando todos los participantes de la ronda votan el mismo valor. Solo el registro del acuerdo finaliza la estimación y actualiza el Story Point actual. Una ronda sin consenso puede dar lugar a US-14, pero no obliga a iniciarla automáticamente.

**Acceptance Criteria:**
- Si todos los participantes de la ronda votaron el mismo valor, ese valor puede registrarse como estimación acordada y pasa a ser el Story Point definitivo y actual.
- Si los votos de los participantes no son todos iguales, la estimación no se finaliza y el Story Point actual permanece sin cambios.
- Revelar votos o iniciar otra ronda, por sí solos, no modifican el Story Point actual.

**Dependencias:** US-13; US-14 solo si se hacen nuevas rondas.

**Definition of Done aplicable:** DoD común; probar consenso unánime, ausencia de consenso y protección del Story Point actual.

**Trazabilidad prevista:** §2.4 → EP-04 → US-15 → SDD de US-15 → BDD de US-15 → TDD de US-15 → Go → tests automatizados.

## EP-05 — Seguimiento de esfuerzo

### US-16 — Registrar horas estimadas y esfuerzo real

**ID / Épica:** US-16 / EP-05. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-05`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero registrar horas estimadas para una User Story durante la planificación y horas trabajadas en el registro de esfuerzo para disponer de ambos datos.

**Contexto y alcance:** Las horas estimadas pertenecen a la historia; las horas reales provienen de registros con miembro, fecha, actividad y horas trabajadas (§2.5). Los registros reales deben poder relacionarse con la historia para permitir la comparación aprobada. No se equiparan Story Points con horas.

**Acceptance Criteria:**
- Durante la planificación se pueden registrar y consultar horas estimadas para una User Story.
- Un registro de esfuerzo conserva miembro, fecha, actividad y horas trabajadas, y permite identificar la historia a la que corresponde el esfuerzo que se comparará.
- Registrar horas reales no sustituye ni altera automáticamente las horas estimadas.

**Dependencias:** US-03, US-05.

**Definition of Done aplicable:** DoD común; probar separación de ambos datos y los cuatro campos obligatorios de cada registro real.

**Trazabilidad prevista:** §2.5 y §2.7 → EP-05 → US-16 → SDD de US-16 → BDD de US-16 → TDD de US-16 → Go → tests automatizados.

### US-17 — Comparar esfuerzo estimado y real

**ID / Épica:** US-17 / EP-05. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-05`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero comparar las horas estimadas de una historia con sus horas reales registradas para identificar diferencias de esfuerzo.

**Contexto y alcance:** Comparación funcional por historia, reutilizando los datos de US-16. US-21 presenta las métricas de horas y su desviación por Sprint/proyecto según fórmulas que se definirán en SDD.

**Acceptance Criteria:**
- Para una historia con estimación y registros de esfuerzo, se muestran por separado las horas estimadas y las reales correspondientes.
- La comparación utiliza los registros de esfuerzo relacionados con esa historia, sin tratar los Story Points como horas.

**Dependencias:** US-16.

**Definition of Done aplicable:** DoD común; verificar que la comparación utiliza las mismas fuentes de datos que US-16 y US-21.

**Trazabilidad prevista:** §2.5 → EP-05 → US-17 → SDD de US-17 → BDD de US-17 → TDD de US-17 → Go → tests automatizados.

## EP-06 — Gestión de defectos

### US-18 — Registrar defecto

**ID / Épica:** US-18 / EP-06. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-06`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero registrar un defecto relacionado con una historia para seguirlo aunque todavía no tenga Sprint asociado.

**Contexto y alcance:** Se registran descripción, severidad, estado, historia relacionada y los campos de Sprint de detección y resolución (§2.6). Los estados permitidos son **Abierto**, **En progreso**, **Resuelto** y **Cerrado**; las severidades son **Baja**, **Media**, **Alta** y **Crítica**. La decisión aprobada permite que inicialmente no haya Sprint.

**Acceptance Criteria:**
- Puede registrarse un defecto con descripción, estado y severidad aprobados e historia relacionada, aunque no haya Sprint asociado inicialmente.
- Si se conoce el Sprint de detección, queda registrado en el campo correspondiente; la ausencia inicial de Sprint no se reemplaza con uno ficticio.

**Dependencias:** US-05; US-07 no es obligatoria.

**Definition of Done aplicable:** DoD común; probar ambos casos, con y sin Sprint de detección, y los valores admitidos de estado y severidad.

**Trazabilidad prevista:** §2.6 y decisión sobre defectos sin Sprint → EP-06 → US-18 → SDD de US-18 → BDD de US-18 → TDD de US-18 → Go → tests automatizados.

### US-19 — Actualizar defecto y registrar su resolución

**ID / Épica:** US-19 / EP-06. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-06`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero actualizar el estado y los Sprints asociados a un defecto para reflejar su evolución y resolución.

**Contexto y alcance:** Un defecto puede asociarse a un Sprint después de creado y resolverse en un Sprint diferente del de detección. Se usan los cuatro estados aprobados; no se inventan transiciones adicionales.

**Acceptance Criteria:**
- Un defecto inicialmente sin Sprint puede asociarse posteriormente a un Sprint de detección cuando corresponda.
- Puede registrarse un Sprint de resolución diferente del Sprint de detección.
- Los cambios de estado y de asociación a Sprints se reflejan al consultar el defecto.

**Dependencias:** US-18; US-07 únicamente si se vincula un Sprint concreto.

**Definition of Done aplicable:** DoD común; probar asociación posterior y resolución en un Sprint distinto.

**Trazabilidad prevista:** §2.6 → EP-06 → US-19 → SDD de US-19 → BDD de US-19 → TDD de US-19 → Go → tests automatizados.

## EP-07 — Métricas

### US-20 — Consultar métricas de Story Points y avance

**ID / Épica:** US-20 / EP-07. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-07`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero consultar Story Points planificados y completados, velocidad y porcentaje de historias completadas para medir el avance.

**Contexto y alcance:** Cuatro métricas obligatorias de §2.7, calculadas por Sprint y agregables al proyecto. El conjunto planificado se preserva aunque cambie el Backlog. Sus fórmulas, períodos exactos y casos límite se definirán en SDD antes de implementar.

**Acceptance Criteria:**
- Para un Sprint se presentan las cuatro métricas según sus SDD aprobados; cuando corresponde, pueden consultarse agregadas para el proyecto.
- Cambios posteriores del Backlog no modifican el conjunto planificado usado para el Sprint.
- Una métrica sin datos suficientes muestra **«No disponible»**, no cero en lugar de esos datos.

**Dependencias:** US-05, US-08, US-09.

**Definition of Done aplicable:** DoD común; SDD con las cuatro fórmulas explícitas y tests para datos completos, cambios de Backlog y datos insuficientes.

**Trazabilidad prevista:** §2.7 → EP-07 → US-20 → SDD de US-20 → BDD de US-20 → TDD de US-20 → Go → tests automatizados.

### US-21 — Consultar métricas de horas

**ID / Épica:** US-21 / EP-07. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-07`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero consultar horas estimadas, horas reales y desviación para medir el esfuerzo.

**Contexto y alcance:** Tres métricas de §2.7, por Sprint y agregables al proyecto. Reutiliza los mismos datos que US-16 y US-17. Las fórmulas y condiciones de agregación se especificarán en SDD, sin asumirlas en esta Issue.

**Acceptance Criteria:**
- Para un Sprint se muestran las tres métricas según sus SDD aprobados y, cuando corresponde, su agregado de proyecto.
- Las horas estimadas provienen de la planificación por historia y las reales del registro de esfuerzo.
- Una métrica sin datos suficientes muestra **«No disponible»**, no cero en lugar de esos datos.

**Dependencias:** US-16, US-17.

**Definition of Done aplicable:** DoD común; SDD con fórmulas de horas y desviación, y tests que reutilicen las fuentes de US-16/US-17.

**Trazabilidad prevista:** §2.7 → EP-07 → US-21 → SDD de US-21 → BDD de US-21 → TDD de US-21 → Go → tests automatizados.

### US-22 — Consultar métricas de defectos

**ID / Épica:** US-22 / EP-07. **Prioridad:** importante. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-07`, `prioridad:importante`.

**Como / Quiero / Para:** Como integrante del equipo, quiero consultar defectos detectados y resueltos para medir su evolución.

**Contexto y alcance:** Dos métricas de §2.7, por Sprint y agregables al proyecto. Su SDD definirá las reglas de conteo y el tratamiento de defectos sin Sprint; esta Issue no presupone una fórmula.

**Acceptance Criteria:**
- Para un Sprint se muestran los defectos detectados y resueltos según el SDD aprobado y, cuando corresponde, su agregado de proyecto.
- El conteo distingue los datos de detección de los de resolución, incluso si están en Sprints diferentes.
- Una métrica sin datos suficientes muestra **«No disponible»**, no cero en lugar de esos datos.

**Dependencias:** US-18, US-19.

**Definition of Done aplicable:** DoD común; SDD con reglas de conteo y tests para Sprints distintos, defectos sin Sprint y datos insuficientes.

**Trazabilidad prevista:** §2.7 → EP-07 → US-22 → SDD de US-22 → BDD de US-22 → TDD de US-22 → Go → tests automatizados.

## EP-08 — Dashboard

### US-23 — Consultar dashboard del proyecto

**ID / Épica:** US-23 / EP-08. **Prioridad:** secundaria. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-08`, `prioridad:secundaria`.

**Como / Quiero / Para:** Como integrante del equipo, quiero ver el estado del proyecto y gráficos de sus métricas para comprender su situación.

**Contexto y alcance:** Dashboard exigido por §2.8. Presenta el estado de US-04 y los resultados de US-20 a US-22, sin crear cálculos alternativos.

**Acceptance Criteria:**
- El dashboard muestra el mismo estado del proyecto que US-04 y gráficos basados en las métricas existentes.
- Cuando una métrica está **«No disponible»**, el dashboard no la representa como si fuera cero.

**Dependencias:** US-04, US-20, US-21, US-22.

**Definition of Done aplicable:** DoD común; verificar concordancia de estado y métricas con sus fuentes, incluidos datos insuficientes.

**Trazabilidad prevista:** §2.8 → EP-08 → US-23 → SDD de US-23 → BDD de US-23 → TDD de US-23 → Go → tests automatizados.

## EP-09 — Reportes

### US-24 — Generar reportes de proyecto, Sprint y PDF final

**ID / Épica:** US-24 / EP-09. **Prioridad:** secundaria. **Etiquetas sugeridas:** `tipo:user-story`, `epic:EP-09`, `prioridad:secundaria`.

**Como / Quiero / Para:** Como integrante del equipo, quiero generar un reporte completo de proyecto, uno de un Sprint individual y un PDF final para consultar y entregar sus resultados.

**Contexto y alcance:** Ambos reportes incluyen User Stories planificadas y completadas, Story Points/estimaciones, horas estimadas y reales, métricas y defectos (§2.9 y decisión aprobada). El PDF final incluye además información general del proyecto, Sprints, User Stories, Story Points, horas estimadas y reales, métricas, defectos y estado general. Reutiliza las métricas de EP-07. Se conserva como una sola User Story conforme a la instrucción del equipo.

**Acceptance Criteria:**
- Se puede generar un reporte completo del proyecto y un reporte de un Sprint individual; ambos contienen los datos exigidos para su alcance.
- Se puede generar un PDF final que incluye información general del proyecto, Sprints, User Stories, Story Points, horas estimadas y reales, métricas, defectos y estado general.
- Los reportes y el PDF reflejan los mismos datos y métricas que las funcionalidades de origen; un dato insuficiente se presenta como **«No disponible»**, no como cero ficticio.

**Dependencias:** US-01, US-05, US-07, US-08, US-09, US-16, US-18, US-19, US-20, US-21 y US-22.

**Definition of Done aplicable:** DoD común; comprobar por separado los dos alcances de reporte, el archivo PDF generado y la concordancia de sus contenidos con las fuentes.

**Trazabilidad prevista:** §2.9 y Sprint 4 → EP-09 → US-24 → SDD de US-24 → BDD de US-24 → TDD de US-24 → Go → tests automatizados.

## Puntos de revisión antes de publicar

- La regla para declarar terminado un proyecto y pasar a **Finalizado** necesita precisión en el SDD de US-02. El cambio a **En progreso** al crear el primer Sprint ya está decidido.
- El mecanismo exacto para finalizar una ronda y determinar quién participa debe definirse en los SDD de Planning Poker; la unanimidad de los participantes de la ronda ya es la regla de consenso aprobada.
- Las fórmulas de las métricas, los casos de datos insuficientes y la atribución de defectos sin Sprint deben quedar explícitos en los SDD de US-20 a US-22 antes de implementar.
- La decisión aprobada de admitir defectos sin Sprint requiere expresar en el SDD cómo se representan temporalmente los campos de Sprint de detección y resolución exigidos por el enunciado.
- US-24 concentra tres salidas y puede ser grande; no se divide en esta etapa. Su tamaño deberá estimarse antes de comprometerla en un Sprint.

