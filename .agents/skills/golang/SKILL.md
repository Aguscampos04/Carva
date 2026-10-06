---
name: golang
description: Define y valida las prácticas de desarrollo en Go para la lógica de negocio, arquitectura, calidad, seguridad y testing.
---

# Golang Skill

## Propósito

Guiar el desarrollo del núcleo de la aplicación utilizando Go (Golang), respetando los requisitos del Trabajo Práctico Integrador.

El código Go debe contener la lógica de negocio y las reglas principales del sistema.

---

## Fuente de verdad

Antes de implementar:

1. consultar `docs/enunciado.md`;
2. consultar `docs/contexto-proyecto.md`;
3. consultar `docs/reglas-agente.md`;
4. revisar la User Story;
5. revisar la SDD;
6. revisar los escenarios BDD;
7. revisar los tests definidos mediante TDD.

No implementar comportamiento que contradiga estos artefactos.

---

## Responsabilidad de Go

Go debe contener el núcleo y las reglas de negocio de la aplicación.

Esto incluye, según corresponda:

* entidades del dominio;
* reglas de negocio;
* validaciones;
* cálculos;
* métricas;
* estimaciones;
* gestión de proyectos;
* gestión de backlog;
* gestión de sprints;
* defectos;
* lógica relacionada con reportes;
* otras funcionalidades requeridas por el dominio.

La interfaz web no debe contener reglas de negocio críticas.

---

## Separación de responsabilidades

Mantener una separación clara entre:

```text
Presentación
    ↓
Aplicación
    ↓
Dominio
    ↓
Persistencia / Infraestructura
```

La arquitectura concreta será definida por el proyecto.

No asumir frameworks o librerías sin necesidad.

---

## Dominio

Las entidades y reglas principales deben mantenerse independientes de detalles externos siempre que sea razonable.

Evitar que la lógica de negocio dependa directamente de:

* HTML;
* JavaScript;
* HTTP;
* base de datos;
* servicios externos.

La lógica central debe poder probarse mediante tests unitarios.

---

## Estructura

Utilizar una estructura de proyecto coherente con el tamaño de la aplicación.

Una posible organización es:

```text
src/
├── cmd/
├── internal/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   └── ...
└── ...
```

Esta estructura es una propuesta, no una obligación.

Si el proyecto adopta otra arquitectura, respetar la decisión aprobada.

---

## Idiomatic Go

Escribir código siguiendo las convenciones habituales de Go.

Priorizar:

* simplicidad;
* claridad;
* composición;
* interfaces pequeñas;
* errores explícitos;
* nombres descriptivos;
* funciones con responsabilidades claras.

Evitar complejidad innecesaria.

---

## Errores

Los errores deben manejarse explícitamente.

Preferir:

```go
result, err := operation()
if err != nil {
    return err
}
```

sobre ignorar errores.

No utilizar `panic` para errores esperables de negocio.

Utilizar errores apropiados y mensajes útiles.

---

## Validaciones

Las validaciones relevantes deben realizarse en el núcleo de negocio.

Ejemplos:

* campos obligatorios;
* valores permitidos;
* estados válidos;
* relaciones entre entidades;
* límites;
* reglas de negocio.

No confiar únicamente en validaciones de la interfaz.

---

## Funciones

Mantener funciones razonablemente pequeñas y enfocadas.

Una función debería tener una responsabilidad clara.

Si una función comienza a acumular múltiples responsabilidades, evaluar su división.

No dividir funciones artificialmente solo para reducir su tamaño.

---

## Interfaces

Utilizar interfaces cuando aporten desacoplamiento real.

No crear interfaces innecesarias únicamente por seguir un patrón.

Las interfaces deben ser pequeñas y representar comportamientos relevantes.

---

## Dependencias

Evitar introducir dependencias externas sin una razón concreta.

Antes de agregar una librería:

1. determinar qué problema resuelve;
2. verificar si puede resolverse con la biblioteca estándar;
3. evaluar mantenimiento y complejidad;
4. verificar compatibilidad con el proyecto;
5. documentar la decisión cuando sea relevante.

---

## Tests

Toda lógica crítica debe poder probarse.

Crear tests para:

* reglas de negocio;
* validaciones;
* cálculos;
* métricas;
* estimaciones;
* casos límite;
* errores.

Seguir las reglas de `skills/tdd/SKILL.md`.

Ejecutar:

```bash
go test ./...
```

antes de considerar terminada una funcionalidad.

---

## Concurrencia

No utilizar concurrencia simplemente porque Go la permite.

Utilizar goroutines, channels u otros mecanismos únicamente cuando exista una necesidad real.

La seguridad y claridad tienen prioridad sobre la complejidad.

---

## Seguridad

No almacenar secretos en el código.

No incluir:

* contraseñas;
* tokens;
* claves privadas;
* credenciales;
* datos sensibles.

Utilizar configuración apropiada para los secretos.

Validar entradas externas antes de utilizarlas.

---

## Configuración

Separar configuración del código cuando corresponda.

No hardcodear valores que deban cambiar entre entornos.

Los valores sensibles nunca deben versionarse.

---

## Observabilidad

Cuando sea necesario, utilizar:

* logs;
* mensajes de error;
* métricas;
* trazas apropiadas.

No llenar el código de logs innecesarios.

No registrar información sensible.

---

## Formato y herramientas

Mantener el código formateado mediante:

```bash
gofmt
```

Cuando estén disponibles, utilizar herramientas estándar como:

```bash
go vet ./...
```

y otras herramientas de análisis apropiadas.

---

## Validación previa al commit

Antes de realizar un commit de una funcionalidad:

```bash
gofmt -w .
go test ./...
go vet ./...
```

Si alguna herramienta falla, investigar y corregir el problema antes de continuar, salvo que exista una razón documentada.

---

## Trazabilidad

El código implementado debe poder relacionarse con:

```text
User Story
    ↓
SDD
    ↓
BDD
    ↓
TDD
    ↓
Código Go
```

Cuando sea útil, utilizar referencias a Issues en commits o documentación.

---

## Cambios

Antes de modificar código existente:

1. comprender el comportamiento actual;
2. revisar tests;
3. revisar SDD y BDD relacionados;
4. identificar posibles regresiones;
5. realizar el cambio;
6. ejecutar los tests.

No modificar código crítico sin verificar sus dependencias.

---

## Regla principal

El código Go debe priorizar:

**correctitud → claridad → mantenibilidad → simplicidad**

antes que la cantidad de funcionalidades implementadas.

La lógica de negocio no debe quedar escondida dentro de la interfaz web.