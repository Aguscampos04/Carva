# Requirements Skill

## Propósito

Analizar los requisitos del proyecto y transformarlos en elementos de trabajo gestionables:

**Requisito → Épica → User Story → Issue → Criterios de aceptación**

La Skill debe priorizar el cumplimiento del enunciado oficial y mantener la trazabilidad de cada requisito.

---

## Fuentes de verdad

Consultar siempre:

1. `docs/enunciado.md`
2. `docs/contexto-proyecto.md`
3. `docs/reglas-agente.md`
4. decisiones aprobadas por el equipo.

No inventar requisitos.

Si una necesidad no aparece en el enunciado ni fue aprobada por el equipo, tratarla como propuesta y no como requisito obligatorio.

---

## Análisis inicial

Al analizar el proyecto:

1. identificar los requisitos funcionales;
2. identificar los requisitos de proceso;
3. identificar restricciones tecnológicas;
4. identificar entregables;
5. identificar requisitos de calidad;
6. identificar requisitos de Scrum;
7. identificar requisitos SDD, BDD y TDD;
8. identificar requisitos de trazabilidad;
9. identificar requisitos relacionados con IA.

Agrupar requisitos relacionados.

---

## Épicas

Las Épicas representan grandes áreas funcionales o de trabajo.

Para este proyecto, considerar inicialmente las áreas indicadas por el enunciado, por ejemplo:

* Gestión de proyectos;
* Product Backlog;
* Sprints;
* Estimación;
* Planning Poker;
* Seguimiento de esfuerzo;
* Defectos;
* Métricas;
* Dashboard;
* Reportes.

Estas Épicas son una propuesta inicial.

Antes de crear nuevas Épicas que cambien significativamente el alcance, verificar el enunciado y las decisiones aprobadas.

No crear Épicas duplicadas.

---

## User Stories

Cada User Story debe representar una funcionalidad o necesidad concreta.

Utilizar preferentemente:

> Como [rol], quiero [acción], para [beneficio].

Una User Story debe:

* ser clara;
* ser específica;
* ser verificable;
* tener criterios de aceptación;
* tener un tamaño razonable;
* poder implementarse y probarse;
* pertenecer a una Épica.

Evitar User Stories demasiado grandes.

Si una User Story contiene varias funcionalidades independientes, proponer su división.

---

## Roles

Utilizar roles relevantes para el sistema.

No inventar roles de usuario si no están definidos.

Cuando una funcionalidad no dependa de un usuario específico, utilizar un actor apropiado como:

* sistema;
* administrador;
* responsable del proyecto;
* miembro del equipo.

La elección del rol debe estar justificada por el comportamiento de la funcionalidad.

---

## Criterios de aceptación

Cada User Story debe tener criterios de aceptación claros.

Deben permitir determinar objetivamente cuándo una historia está terminada.

Ejemplo:

```text
Dado que existe un proyecto creado
Cuando el usuario agrega una User Story válida
Entonces la User Story queda registrada en el Product Backlog
```

Los criterios deben cubrir:

* comportamiento esperado;
* validaciones relevantes;
* errores relevantes;
* casos límite cuando correspondan.

---

## Issues

Cada User Story implementable debe poder convertirse en una o más Issues.

Una Issue debe contener como mínimo:

* título claro;
* descripción;
* relación con la User Story;
* criterios de aceptación;
* prioridad cuando corresponda;
* información suficiente para comenzar el trabajo.

No crear Issues duplicadas.

No crear una Issue simplemente por cada requisito textual si eso produce tareas artificiales.

---

## Dependencias

Detectar dependencias entre funcionalidades.

Ejemplos:

* no se puede crear un Sprint sin un proyecto;
* no se puede agregar una User Story a un Product Backlog inexistente;
* no se puede calcular una métrica si no existen los datos necesarios.

Registrar dependencias cuando afecten el orden de implementación.

---

## Priorización

Priorizar considerando:

1. requisitos obligatorios del enunciado;
2. funcionalidades necesarias para el MVP;
3. dependencias;
4. valor funcional;
5. complejidad;
6. riesgos técnicos.

No priorizar únicamente por facilidad de implementación.

---

## MVP

Identificar las funcionalidades mínimas necesarias para demostrar un producto funcional.

El MVP debe permitir construir progresivamente las funcionalidades requeridas por el enunciado.

No agregar funcionalidades no obligatorias antes de completar las esenciales.

---

## Detección de ambigüedades

Cuando un requisito sea ambiguo:

1. identificar la ambigüedad;
2. explicar qué interpretación sería posible;
3. proponer una alternativa;
4. solicitar decisión si afecta el alcance o comportamiento.

No ocultar ambigüedades mediante suposiciones.

---

## Trazabilidad

Cada User Story debe poder relacionarse con:

* requisito de origen;
* Épica;
* Issue;
* SDD;
* criterios de aceptación;
* BDD;
* tests;
* código.

La trazabilidad debe mantenerse durante todo el proyecto.

---

## Cambios de requisitos

Cuando se modifique un requisito:

1. identificar las User Stories afectadas;
2. identificar las Issues afectadas;
3. revisar SDD;
4. revisar BDD;
5. revisar tests;
6. revisar implementación;
7. actualizar documentación.

No modificar solamente el código si el cambio afecta otros artefactos.

---

## Integración con GitHub

Cuando GitHub esté disponible:

* consultar Issues existentes antes de crear nuevas;
* evitar duplicados;
* crear Épicas/User Stories/Issues según las capacidades disponibles;
* actualizar Issues cuando cambie el estado;
* mantener referencias entre Issues y artefactos del proyecto;
* informar qué elementos fueron creados o modificados.

Nunca afirmar que una Issue fue creada si la operación no fue realmente ejecutada.

---

## Resultado esperado

Cuando se solicite analizar requisitos, producir una estructura clara similar a:

```text
Épica
├── User Story
│   ├── Criterios de aceptación
│   └── Issue
│
├── User Story
│   ├── Criterios de aceptación
│   └── Issue
│
└── User Story
    ├── Criterios de aceptación
    └── Issue
```

El resultado debe ser suficiente para que las siguientes Skills puedan generar SDD, BDD, TDD y código.

---

## Regla principal

El objetivo no es crear la mayor cantidad posible de Issues.

El objetivo es crear una **descomposición útil, completa, trazable y ejecutable del proyecto**.