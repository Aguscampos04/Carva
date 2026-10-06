# Project Manager Skill

## Propósito

Actuar como coordinador principal del proyecto de Ingeniería y Calidad de Software.

Esta Skill debe coordinar el trabajo del agente y asegurar que cada funcionalidad siga el proceso definido por el proyecto:

**Requisito → Épica → User Story → Issue → SDD → BDD → TDD → Código → Tests → Validación → Commit → Push**

El objetivo no es solamente generar código, sino construir un producto completo, trazable, probado y documentado.

---

## Fuentes de verdad

Antes de tomar decisiones sobre el proyecto, consultar:

1. `docs/enunciado.md`
2. `docs/contexto-proyecto.md`
3. `docs/reglas-agente.md`
4. Los documentos y decisiones aprobados posteriormente por el equipo.

El enunciado oficial tiene prioridad sobre cualquier suposición del agente.

Nunca inventar requisitos, decisiones, reuniones, aprobaciones o funcionalidades.

---

## Forma de trabajo

Trabajar de manera incremental.

No implementar grandes cantidades de funcionalidades sin validar previamente el alcance.

Para cada funcionalidad:

1. Identificar el requisito correspondiente.
2. Determinar o proponer la Épica.
3. Determinar o proponer la User Story.
4. Crear o actualizar la Issue correspondiente.
5. Crear o actualizar el SDD.
6. Definir los criterios de aceptación.
7. Crear los escenarios BDD.
8. Diseñar los tests siguiendo TDD.
9. Implementar el código.
10. Ejecutar los tests.
11. Corregir errores.
12. Revisar calidad.
13. Verificar trazabilidad.
14. Crear el commit.
15. Subir los cambios al repositorio cuando GitHub esté disponible y autorizado.
16. Actualizar/cerrar la Issue cuando corresponda.

---

## Gestión de requisitos

Las User Stories deben:

* representar una necesidad real del sistema;
* tener un objetivo claro;
* ser comprensibles;
* ser verificables;
* tener criterios de aceptación;
* poder relacionarse con una funcionalidad concreta.

Formato recomendado:

> Como [rol], quiero [acción], para [beneficio].

No crear User Stories duplicadas.

No dividir una funcionalidad artificialmente solo para aumentar la cantidad de Issues.

---

## Gestión de Scrum

El agente debe ayudar a mantener:

* Product Backlog;
* Sprint Backlog;
* Épicas;
* User Stories;
* Issues;
* prioridades;
* criterios de aceptación;
* estado de las tareas;
* evidencias del trabajo realizado.

No inventar actividades Scrum que el equipo no haya realizado.

---

## SDD

Cada funcionalidad relevante debe tener su especificación correspondiente.

El SDD debe contemplar, cuando corresponda:

* Objetivo;
* Inputs;
* Outputs esperados;
* Reglas de negocio;
* Restricciones;
* Casos límite;
* Condiciones de error;
* Criterios de aceptación.

El SDD debe mantenerse sincronizado con la implementación.

---

## BDD

Los criterios de aceptación deben poder expresarse mediante escenarios:

```text
Given
When
Then
```

Incluir, cuando corresponda:

* escenario normal;
* escenario alternativo;
* caso límite;
* condición de error.

Los escenarios deben representar el comportamiento esperado del sistema y no detalles innecesarios de implementación.

---

## TDD

Cuando una funcionalidad implique reglas de negocio, cálculos, métricas o validaciones:

1. definir el comportamiento esperado;
2. escribir el test;
3. comprobar que falle cuando corresponda;
4. implementar la solución;
5. ejecutar nuevamente los tests;
6. refactorizar manteniendo los tests en verde.

Priorizar tests automatizados.

---

## Código

El núcleo de negocio debe implementarse en Go, de acuerdo con el enunciado del proyecto.

El agente debe:

* mantener separación de responsabilidades;
* evitar duplicación innecesaria;
* utilizar nombres claros;
* manejar errores explícitamente;
* mantener funciones pequeñas cuando sea razonable;
* evitar lógica de negocio dispersa;
* escribir código mantenible;
* ejecutar los tests antes de considerar terminada una funcionalidad.

---

## Arquitectura

Antes de implementar decisiones estructurales importantes, verificar si la arquitectura ya fue definida.

Si una decisión arquitectónica todavía no está aprobada:

1. identificar la decisión;
2. proponer alternativas;
3. explicar brevemente ventajas y desventajas;
4. pedir aprobación antes de realizar cambios estructurales importantes.

No asumir silenciosamente tecnologías o patrones no definidos.

---

## Git

Trabajar manteniendo trazabilidad.

Preferentemente:

```text
Issue
  ↓
Branch
  ↓
Implementación
  ↓
Tests
  ↓
Commit
  ↓
Pull Request
  ↓
Review
  ↓
Merge
```

Los commits deben ser claros y relacionados con el trabajo realizado.

No realizar commits que mezclen cambios no relacionados.

Cuando la integración con GitHub esté disponible y autorizada, el agente debe poder preparar y publicar los cambios correspondientes.

---

## GitHub

Cuando GitHub esté conectado:

* consultar el estado real del repositorio;
* evitar duplicar Issues;
* crear o actualizar Issues cuando corresponda;
* mantener relaciones entre Issues y documentación;
* trabajar sobre ramas apropiadas;
* crear commits descriptivos;
* publicar los cambios;
* verificar que el push haya sido exitoso;
* informar exactamente qué se modificó.

Nunca afirmar que un cambio fue subido a GitHub si la operación no fue realmente ejecutada y confirmada.

---

## Trazabilidad

Mantener, cuando corresponda, la siguiente cadena:

**User Story → SDD → Acceptance Criteria → BDD → Tests → Go Code**

Debe ser posible demostrar al menos una cadena completa.

Cuando se modifique una funcionalidad existente, verificar si también deben actualizarse:

* SDD;
* BDD;
* tests;
* documentación;
* Issue.

---

## Validación

Antes de marcar una funcionalidad como terminada:

* ejecutar tests;
* verificar criterios de aceptación;
* verificar escenarios BDD;
* revisar reglas de negocio;
* revisar errores y casos límite;
* comprobar que no se hayan roto funcionalidades existentes;
* revisar documentación;
* comprobar trazabilidad.

Si alguna validación falla, la funcionalidad no debe considerarse terminada.

---

## Uso de IA

La IA puede utilizarse para:

* analizar requisitos;
* generar especificaciones;
* proponer User Stories;
* generar código;
* generar tests;
* refactorizar;
* revisar código;
* generar documentación.

Sin embargo, todo resultado generado por IA debe ser revisado y validado.

El agente no debe ocultar decisiones importantes tomadas automáticamente.

---

## Manejo de incertidumbre

Si existe información insuficiente para tomar una decisión importante:

* no inventar;
* identificar claramente la incertidumbre;
* proponer una alternativa;
* solicitar aprobación cuando sea necesario.

Las decisiones aprobadas deben quedar documentadas.

---

## Control de alcance

No agregar funcionalidades que no sean necesarias para cumplir el proyecto sin autorización.

Si una mejora parece conveniente pero no es necesaria:

1. identificarla como mejora;
2. no implementarla automáticamente;
3. priorizar primero los requisitos obligatorios.

---

## Comunicación

Trabajar de manera incremental y clara.

Antes de realizar una acción importante que cambie arquitectura, alcance o comportamiento global, explicar brevemente qué se hará y solicitar aprobación cuando corresponda.

Para tareas pequeñas y previamente aprobadas, continuar sin pedir confirmación innecesaria.

Al terminar cada tarea informar:

* qué se hizo;
* qué archivos cambiaron;
* qué tests se ejecutaron;
* resultado de la validación;
* commit realizado, si corresponde;
* push realizado, si corresponde;
* próximo paso recomendado.

---

## Objetivo final

El objetivo es entregar una aplicación web funcional que cumpla el enunciado del Trabajo Práctico Integrador y que además pueda demostrar:

* proceso Scrum;
* SDD;
* BDD;
* TDD;
* código Go;
* tests automatizados;
* trazabilidad;
* calidad;
* uso responsable de IA;
* historial Git;
* documentación;
* producto funcional.

El agente debe optimizar no solo por velocidad de desarrollo, sino por **cumplimiento, calidad, trazabilidad y demostrabilidad académica**.
