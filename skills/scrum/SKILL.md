# Scrum Skill

## Propósito

Ayudar a gestionar el proyecto utilizando Scrum de acuerdo con el Trabajo Práctico Integrador.

La Skill debe mantener el trabajo organizado en:

* Product Backlog;
* Sprints;
* Sprint Backlog;
* User Stories;
* Issues;
* planificación;
* seguimiento;
* Review;
* Retrospective.

---

## Fuente de verdad

Consultar:

1. `docs/enunciado.md`;
2. `docs/contexto-proyecto.md`;
3. `docs/reglas-agente.md`;
4. decisiones reales del equipo;
5. Issues y estado real del repositorio cuando GitHub esté disponible.

No inventar reuniones ni actividades Scrum que no hayan ocurrido.

---

## Roles

Respetar los roles definidos por el enunciado:

### Product Architect

Corresponde a los profesores.

### Agile Enabler

Corresponde a un integrante del equipo.

### Product Builders

Corresponde a los integrantes del equipo.

No reasignar estos roles arbitrariamente.

---

## Product Backlog

El Product Backlog debe contener el trabajo necesario para cumplir el proyecto.

Cada elemento debería tener:

* título;
* descripción;
* User Story cuando corresponda;
* prioridad;
* criterios de aceptación;
* relación con Épica;
* estado.

No agregar funcionalidades sin justificación.

---

## Sprint

Cada Sprint debe tener:

* objetivo;
* conjunto de trabajo seleccionado;
* User Stories;
* Issues;
* criterios de finalización;
* evidencia del trabajo realizado.

No considerar que una tarea está terminada solamente porque se creó el código.

---

## Sprint Goal

Cada Sprint debería tener un objetivo claro.

El objetivo debe describir el resultado que se espera alcanzar.

Ejemplo:

```text id="1kzqfd"
Sprint Goal:
Disponer de un MVP capaz de crear proyectos,
gestionar el Product Backlog y administrar Sprints.
```

El objetivo debe ser coherente con el alcance del Sprint.

---

## Planificación

Antes de comenzar un Sprint:

1. revisar Product Backlog;
2. revisar prioridades;
3. verificar dependencias;
4. seleccionar User Stories;
5. definir Sprint Goal;
6. identificar riesgos;
7. verificar capacidad del equipo.

No asignar trabajo ficticio.

---

## Seguimiento

Durante el Sprint, mantener actualizado:

* estado de Issues;
* avance de User Stories;
* bloqueos;
* defectos;
* cambios de alcance.

Cuando exista GitHub, utilizar el estado real de Issues y Pull Requests.

---

## Daily

La Skill puede ayudar a estructurar la información para una Daily:

```text id="f35n0s"
¿Qué hice?
¿Qué voy a hacer?
¿Qué bloqueos tengo?
```

No inventar respuestas.

Utilizar solamente información real proporcionada por el equipo o disponible en el repositorio.

---

## Sprint Review

Registrar:

* qué funcionalidades fueron completadas;
* qué funcionalidades no fueron completadas;
* qué comportamiento puede demostrarse;
* qué criterios de aceptación fueron cumplidos;
* qué problemas quedaron pendientes.

No declarar una funcionalidad terminada si sus criterios obligatorios no fueron cumplidos.

---

## Retrospective

Ayudar a identificar:

### Qué salió bien

Procesos o prácticas que funcionaron.

### Qué podría mejorar

Problemas o dificultades reales.

### Acciones de mejora

Acciones concretas para el siguiente Sprint.

No inventar problemas ni conclusiones.

---

## Gestión de cambios

Si durante un Sprint aparece una nueva necesidad:

1. determinar si proviene del enunciado;
2. determinar si fue solicitada por los profesores;
3. determinar si es una mejora;
4. evaluar impacto;
5. actualizar el backlog;
6. decidir cuándo incorporarla.

No introducir cambios silenciosamente.

---

## Priorización

Priorizar considerando:

1. cumplimiento del enunciado;
2. dependencias;
3. MVP;
4. riesgos;
5. valor funcional;
6. esfuerzo.

Las funcionalidades obligatorias tienen prioridad sobre mejoras opcionales.

---

## Issues

Cuando GitHub esté disponible:

* consultar Issues existentes;
* evitar duplicados;
* asociar Issues con User Stories;
* mantener estados actualizados;
* registrar bloqueos;
* vincular Pull Requests cuando corresponda.

No cerrar una Issue si la funcionalidad todavía no cumple sus criterios de aceptación.

---

## Definition of Done

Utilizar una Definition of Done consistente.

Una User Story puede considerarse terminada cuando:

```text id="w7f1tv"
[ ] cumple los criterios de aceptación
[ ] SDD actualizada
[ ] BDD actualizado
[ ] tests implementados
[ ] tests ejecutados correctamente
[ ] código revisado
[ ] documentación actualizada
[ ] trazabilidad verificada
[ ] cambios versionados
```

Agregar otros criterios cuando el proyecto los defina.

---

## Evidencia Scrum

El proyecto debe conservar evidencia real del proceso.

Puede incluir:

* Issues;
* Sprint Backlogs;
* Product Backlog;
* GitHub Project;
* Pull Requests;
* commits;
* Reviews;
* Retrospectives;
* documentación.

No fabricar evidencia de actividades que no ocurrieron.

---

## Relación con desarrollo

Scrum debe integrarse con el flujo técnico:

```text id="7f4vks"
Product Backlog
      ↓
Sprint
      ↓
User Story
      ↓
Issue
      ↓
SDD
      ↓
BDD
      ↓
TDD
      ↓
Código
      ↓
Tests
      ↓
Review
      ↓
Done
```

---

## Regla principal

Scrum debe ayudar a organizar y demostrar el trabajo real del equipo.

No se debe utilizar para generar documentación ficticia.

La prioridad es mantener un proceso simple, visible, trazable y coherente con el desarrollo real.
