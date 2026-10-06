---
name: sdd
description: Define y valida Software Design Documents según los requisitos del TPI y mantiene la trazabilidad con User Stories, criterios de aceptación, BDD, TDD y código.
---

# SDD Skill

## Propósito

Crear y mantener las especificaciones de diseño y comportamiento requeridas por el proyecto.

La SDD debe funcionar como puente entre los requisitos y la implementación.

Flujo:

**User Story → SDD → Acceptance Criteria → BDD → TDD → Código**

---

## Fuente de verdad

Antes de crear o modificar una SDD, consultar:

1. `docs/enunciado.md`
2. `docs/contexto-proyecto.md`
3. `docs/reglas-agente.md`
4. User Story correspondiente;
5. decisiones aprobadas por el equipo.

No inventar reglas de negocio.

---

## Estructura obligatoria

Cada SDD debe contener, cuando corresponda:

1. Objetivo
2. Inputs
3. Expected Outputs
4. Business Rules
5. Constraints
6. Edge Cases
7. Error Conditions
8. Acceptance Criteria

Los nombres de las secciones deben mantenerse consistentes.

---

## Objetivo

Explicar claramente qué problema resuelve la funcionalidad.

Debe ser concreto y estar relacionado con la User Story.

No describir detalles de implementación innecesarios.

---

## Inputs

Identificar todos los datos necesarios para ejecutar la funcionalidad.

Para cada input relevante indicar:

* nombre;
* tipo cuando sea necesario;
* si es obligatorio;
* restricciones relevantes.

Ejemplo:

```text id="4jv1y6"
nombre: título del proyecto
tipo: string
obligatorio: sí
restricción: no puede estar vacío
```

---

## Expected Outputs

Definir qué resultado debe producir la funcionalidad.

Puede incluir:

* datos;
* cambios de estado;
* registros;
* mensajes;
* errores.

El resultado debe poder verificarse mediante tests.

---

## Business Rules

Documentar las reglas de negocio.

Las reglas deben ser:

* claras;
* verificables;
* numerables cuando sea útil;
* independientes de detalles de implementación.

Ejemplo:

```text id="gkgx8h"
BR-01: Un proyecto debe tener un nombre.
BR-02: El nombre no
```