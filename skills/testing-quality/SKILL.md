# Testing & Quality Skill

## Propósito

Garantizar que el software cumpla los requisitos funcionales, las reglas de negocio y los criterios de calidad definidos para el proyecto.

Esta Skill coordina la validación general del sistema y complementa la Skill de TDD.

---

## Fuentes de verdad

Consultar:

1. `docs/enunciado.md`;
2. `docs/contexto-proyecto.md`;
3. `docs/reglas-agente.md`;
4. User Stories;
5. SDD;
6. BDD;
7. tests existentes;
8. decisiones aprobadas.

---

## Objetivos de calidad

Verificar:

* funcionalidad;
* correctitud;
* mantenibilidad;
* testabilidad;
* seguridad básica;
* estabilidad;
* trazabilidad;
* documentación;
* ausencia de regresiones.

No considerar que una funcionalidad está terminada únicamente porque compila.

---

## Validación automática

Antes de considerar terminado un cambio relevante, ejecutar las herramientas disponibles y apropiadas.

Para Go, como mínimo:

```bash
go test ./...
```

Cuando corresponda:

```bash
go test -cover ./...
go vet ./...
gofmt -l .
```

Si alguna herramienta no está disponible, informar la limitación.

---

## Tests

Verificar que existan tests para la lógica relevante.

Priorizar:

* reglas de negocio;
* cálculos;
* métricas;
* estimaciones;
* validaciones;
* casos límite;
* errores;
* funcionalidades críticas.

No crear tests artificiales únicamente para aumentar el porcentaje de cobertura.

---

## Cobertura

La cobertura debe utilizarse como indicador de calidad, no como único objetivo.

Cuando corresponda ejecutar:

```bash
go test -cover ./...
```

Si se obtiene una cobertura baja:

1. identificar las áreas sin cobertura;
2. determinar si son críticas;
3. agregar tests cuando sea necesario;
4. evitar tests sin valor simplemente para aumentar el porcentaje.

---

## Regresión

Después de modificar una funcionalidad:

1. ejecutar los tests relacionados;
2. ejecutar el conjunto completo de tests;
3. verificar funcionalidades dependientes;
4. revisar errores nuevos.

El comando general recomendado es:

```bash
go test ./...
```

---

## Validación de criterios de aceptación

Cada User Story debe validarse contra sus criterios de aceptación.

Para cada criterio determinar:

```text
PASS
FAIL
NOT TESTED
```

Una historia no debe marcarse como terminada si alguno de sus criterios obligatorios falla.

---

## Validación BDD

Comprobar que:

* los escenarios BDD representan el comportamiento esperado;
* no contradicen la SDD;
* los casos importantes están cubiertos;
* los escenarios de error relevantes existen;
* los tests relacionados cubren el comportamiento definido.

---

## Validación de SDD

Verificar que la implementación no contradiga:

* reglas de negocio;
* restricciones;
* outputs esperados;
* condiciones de error;
* casos límite;
* criterios de aceptación.

Si la implementación cambió intencionalmente el comportamiento, actualizar primero la documentación correspondiente.

---

## Revisiones de código

Cuando sea necesario, revisar:

### Correctitud

¿El código hace lo que debería hacer?

### Claridad

¿Otro desarrollador puede entenderlo?

### Mantenibilidad

¿Es razonablemente fácil modificarlo?

### Errores

¿Los errores se manejan correctamente?

### Seguridad

¿Existen entradas no validadas o secretos expuestos?

### Duplicación

¿Existe código repetido innecesariamente?

### Complejidad

¿Existe complejidad innecesaria?

---

## Seguridad básica

Buscar problemas evidentes como:

* credenciales hardcodeadas;
* secretos versionados;
* entradas no validadas;
* exposición innecesaria de información;
* manejo incorrecto de errores;
* configuraciones inseguras.

No afirmar que una aplicación es completamente segura únicamente por pasar estas verificaciones.

---

## Dependencias

Cuando se agregue una dependencia:

* verificar su necesidad;
* revisar su uso;
* evitar dependencias innecesarias;
* documentar decisiones relevantes.

No actualizar dependencias masivamente sin necesidad.

---

## Defectos

Cuando se detecte un defecto:

1. describir el problema;
2. identificar cómo reproducirlo;
3. identificar el comportamiento esperado;
4. identificar el comportamiento actual;
5. determinar el requisito afectado;
6. crear o actualizar la Issue correspondiente;
7. agregar o modificar un test que reproduzca el defecto cuando sea apropiado;
8. corregirlo;
9. ejecutar nuevamente los tests.

Preferentemente seguir:

**Test que reproduce el defecto → Corrección → Test en verde**

---

## Smoke Test

Cuando exista una versión ejecutable de la aplicación, realizar una validación básica:

* la aplicación inicia;
* las funcionalidades principales están disponibles;
* no existen errores bloqueantes;
* las operaciones principales pueden ejecutarse.

No reemplaza los tests automatizados.

---

## Checklist de finalización

Antes de cerrar una funcionalidad:

```text
[ ] User Story identificada
[ ] Criterios de aceptación definidos
[ ] SDD actualizada
[ ] BDD actualizado
[ ] Tests implementados
[ ] Código implementado
[ ] gofmt ejecutado
[ ] go test ./... ejecutado
[ ] go vet ./... ejecutado cuando corresponda
[ ] Cobertura revisada cuando corresponda
[ ] Casos límite revisados
[ ] Errores revisados
[ ] Regresiones verificadas
[ ] Documentación actualizada
[ ] Trazabilidad verificada
```

---

## Resultado de validación

Al finalizar una validación, informar claramente:

```text
Resultado: PASS / FAIL

Tests:
[resultado]

Cobertura:
[resultado si corresponde]

go vet:
[resultado si corresponde]

Formato:
[resultado]

Criterios de aceptación:
[resultado]

Problemas encontrados:
[lista]

Acciones realizadas:
[lista]
```

No ocultar fallos.

---

## Regla principal

La calidad no se mide solamente por si el programa funciona.

Debe verificarse que:

**funciona + cumple requisitos + está probado + es mantenible + es trazable**.