# BDD Skill

## Propósito

Crear y mantener escenarios de comportamiento del sistema utilizando Behavior-Driven Development (BDD).

La Skill debe transformar los criterios de aceptación definidos en la SDD en escenarios verificables.

Flujo:

**User Story → SDD → Acceptance Criteria → BDD → Tests → Código**

---

## Fuentes de verdad

Antes de crear o modificar escenarios BDD, consultar:

1. `docs/enunciado.md`
2. `docs/contexto-proyecto.md`
3. `docs/reglas-agente.md`
4. User Story correspondiente;
5. SDD correspondiente;
6. decisiones aprobadas por el equipo.

La SDD y los criterios de aceptación son la principal fuente para definir el comportamiento esperado.

---

## Estructura

Los escenarios deben utilizar la estructura:

```text
Feature: [nombre de la funcionalidad]

Scenario: [nombre del escenario]
  Given [contexto inicial]
  When [acción]
  Then [resultado esperado]
```

Cuando sea necesario pueden utilizarse múltiples condiciones:

```text
Given ...
And ...
When ...
And ...
Then ...
And ...
```

---

## Tipos de escenarios

Cuando corresponda, contemplar:

1. escenario normal;
2. escenario alternativo;
3. caso límite;
4. condición de error.

No es obligatorio crear todos los tipos para todas las funcionalidades.

Solo deben incluirse los escenarios relevantes al comportamiento real.

---

## Relación con criterios de aceptación

Cada criterio de aceptación debe poder relacionarse con uno o más escenarios.

Ejemplo:

```text
AC-01 → BDD-01
AC-02 → BDD-02
AC-03 → BDD-03
```

No crear escenarios que no tengan una justificación en los requisitos, SDD o decisiones aprobadas.

---

## Escenario normal

Representar el flujo esperado cuando todos los datos y condiciones son válidos.

Ejemplo:

```text
Scenario: Crear un proyecto correctamente
  Given que el usuario está creando un proyecto
  And proporciona un nombre válido
  When confirma la creación
  Then el sistema crea el proyecto
  And el proyecto queda disponible para su gestión
```

---

## Escenario alternativo

Representar una forma válida diferente de completar la operación.

Ejemplo:

```text
Scenario: Crear un proyecto con datos opcionales
  Given que el usuario proporciona un nombre válido
  And no proporciona un campo opcional
  When confirma la creación
  Then el sistema crea el proyecto utilizando los valores permitidos
```

Solo utilizar escenarios alternativos cuando exista realmente un comportamiento alternativo definido.

---

## Casos límite

Representar valores o condiciones en los límites del comportamiento permitido.

Ejemplos:

* cantidad mínima;
* cantidad máxima;
* lista vacía;
* primer o último elemento;
* valores límite permitidos.

Ejemplo:

```text
Scenario: Crear un proyecto con el límite permitido de caracteres
  Given que el usuario proporciona un nombre con la longitud máxima permitida
  When confirma la creación
  Then el sistema acepta el nombre
```

---

## Condiciones de error

Representar situaciones donde la operación debe ser rechazada o producir un error controlado.

Ejemplo:

```text
Scenario: Intentar crear un proyecto sin nombre
  Given que el usuario no proporciona un nombre
  When intenta crear el proyecto
  Then el sistema rechaza la operación
  And informa que el nombre es obligatorio
```

---

## Lenguaje

Los escenarios deben:

* ser claros;
* describir comportamiento observable;
* evitar detalles internos de implementación;
* evitar mencionar funciones, clases o variables;
* poder ser comprendidos por personas no técnicas.

No escribir:

```text
When se ejecuta la función CreateProject()
```

Preferir:

```text
When el usuario confirma la creación del proyecto
```

---

## Determinismo

Los escenarios deben poder ejecutarse de forma reproducible.

Evitar depender de:

* fechas aleatorias;
* datos aleatorios;
* estados externos no controlados;
* servicios externos innecesarios;
* información que pueda cambiar sin control.

Cuando sea necesario utilizar datos controlados.

---

## Datos de prueba

Cuando un escenario requiera datos concretos, definirlos claramente.

Ejemplo:

```text
Given que existe un proyecto llamado "Proyecto Demo"
And tiene 3 User Stories
```

No utilizar datos ambiguos como:

```text
Given que existe un proyecto cualquiera
```

---

## Organización

Los escenarios BDD deben almacenarse dentro de:

```text
bdd/
```

Utilizar nombres identificables.

Formato recomendado:

```text
BDD-[ID]-[nombre-corto].feature
```

Ejemplo:

```text
BDD-001-crear-proyecto.feature
```

---

## Trazabilidad

Mantener la relación:

```text
User Story
    ↓
SDD
    ↓
Acceptance Criteria
    ↓
BDD Scenario
    ↓
Test
    ↓
Código
```

Cuando sea útil, incluir identificadores explícitos.

Ejemplo:

```text
AC-01 → BDD-01 → TEST-01
```

---

## Automatización

Cuando sea técnicamente viable, los escenarios BDD deben poder utilizarse como base para pruebas automatizadas.

No es obligatorio forzar una herramienta BDD específica si el proyecto no la requiere.

La automatización debe mantener el comportamiento definido por los escenarios.

---

## Consistencia

Si cambia una regla de negocio o criterio de aceptación:

1. revisar los escenarios afectados;
2. actualizar los escenarios;
3. revisar los tests;
4. revisar la implementación;
5. verificar trazabilidad.

No dejar escenarios que describan un comportamiento que el sistema ya no debe tener.

---

## Validación

Antes de considerar una especificación BDD terminada:

* cada escenario tiene Given/When/Then;
* describe comportamiento observable;
* deriva de un criterio de aceptación;
* contempla errores relevantes;
* contempla casos límite relevantes;
* no contiene detalles innecesarios de implementación;
* puede relacionarse con tests;
* no contradice la SDD.

---

## Regla principal

BDD debe describir **qué comportamiento espera el usuario o el sistema**, no cómo está implementado internamente.

Un escenario BDD debe poder leerse como una especificación ejecutable del comportamiento esperado.