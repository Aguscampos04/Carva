# TDD Skill

## Propósito

Aplicar Test-Driven Development al desarrollo del proyecto.

El objetivo es que las funcionalidades relevantes sean desarrolladas mediante el ciclo:

**RED → GREEN → REFACTOR**

Los tests deben servir como evidencia del comportamiento esperado y contribuir a la trazabilidad del proyecto.

---

## Fuentes de verdad

Antes de crear o modificar tests, consultar:

1. `docs/enunciado.md`
2. `docs/contexto-proyecto.md`
3. `docs/reglas-agente.md`
4. User Story;
5. SDD;
6. criterios de aceptación;
7. escenarios BDD.

No inventar comportamientos que no estén definidos.

---

## Ciclo TDD

Para cada funcionalidad relevante:

### 1. RED

Primero definir y escribir el test correspondiente al comportamiento esperado.

Ejecutar el test.

El test debe fallar cuando todavía no existe la implementación necesaria.

La falla debe ser entendible y estar relacionada con la funcionalidad que se está desarrollando.

---

### 2. GREEN

Implementar la cantidad mínima de código necesaria para que el test pase.

No implementar funcionalidades adicionales que todavía no hayan sido requeridas.

Ejecutar nuevamente los tests.

El resultado esperado es:

```text
PASS
```

---

### 3. REFACTOR

Una vez que los tests estén en verde:

* mejorar estructura;
* eliminar duplicación;
* mejorar nombres;
* simplificar código;
* mejorar mantenibilidad.

Después del refactor, ejecutar nuevamente los tests.

Los tests deben continuar pasando.

---

## Qué debe probarse

Como mínimo, priorizar tests automatizados para:

* métricas;
* cálculos de estimación;
* reglas de negocio;
* validaciones;
* casos límite;
* condiciones de error;
* lógica crítica del dominio.

También crear tests para otros comportamientos cuando aporten valor.

---

## Tests unitarios

El código Go debe contar con tests unitarios para la lógica de negocio.

Preferir tests:

* deterministas;
* independientes;
* rápidos;
* fáciles de mantener;
* enfocados en una responsabilidad.

Evitar tests que dependan innecesariamente de:

* red;
* servicios externos;
* bases de datos reales;
* estado global;
* tiempo real.

Cuando una dependencia externa sea necesaria, utilizar una estrategia apropiada de aislamiento.

---

## Table-Driven Tests

En Go, cuando existan múltiples casos similares, preferir tests table-driven.

Ejemplo conceptual:

```go
tests := []struct {
    name    string
    input   ...
    want    ...
}{
    {
        name:  "caso válido",
        input: ...,
        want:  ...,
    },
    {
        name:  "caso inválido",
        input: ...,
        want:  ...,
    },
}
```

No utilizar esta técnica cuando haga que el test sea menos claro.

---

## Casos a cubrir

Cuando corresponda, contemplar:

### Caso válido

Los datos cumplen todas las reglas.

### Caso inválido

Los datos violan una regla.

### Caso límite

Los datos están en el límite permitido.

### Caso de error

La operación debe fallar de forma controlada.

---

## Relación con BDD

Cada escenario BDD relevante debe poder relacionarse con uno o más tests.

Ejemplo:

```text
BDD-01 → TEST-01
BDD-02 → TEST-02
BDD-03 → TEST-03
```

No es necesario que exista una correspondencia estrictamente uno a uno.

Un test puede cubrir varios criterios relacionados cuando tenga sentido.

---

## Relación con SDD

Las reglas de negocio deben poder rastrearse hasta los tests.

Ejemplo:

```text
BR-01 → TEST-01
BR-02 → TEST-02
AC-01 → BDD-01 → TEST-01
```

La trazabilidad debe permitir demostrar que las reglas importantes fueron verificadas.

---

## Ubicación

Los tests deben mantenerse junto al código Go correspondiente cuando sea apropiado, utilizando las convenciones estándar de Go:

```text
*_test.go
```

La carpeta:

```text
tests/
```

puede utilizarse para pruebas o evidencias adicionales que no correspondan directamente a tests unitarios del código Go.

No duplicar innecesariamente tests entre ambas ubicaciones.

---

## Ejecución

Antes de considerar una funcionalidad terminada, ejecutar:

```bash
go test ./...
```

Si corresponde, también utilizar:

```bash
go test -cover ./...
```

La cobertura es una métrica de apoyo y no reemplaza la calidad de los tests.

---

## Fallos

Si un test falla:

1. identificar el motivo;
2. determinar si falla el código o el propio test;
3. corregir el problema;
4. volver a ejecutar los tests;
5. verificar que no existan regresiones.

No eliminar o debilitar un test simplemente para conseguir que pase.

---

## Evidencia TDD

El proceso debe poder demostrarse mediante:

* tests;
* commits;
* historial Git;
* mensajes de commit;
* documentación cuando corresponda.

Cuando sea útil, utilizar commits que permitan identificar el ciclo TDD.

Ejemplo:

```text
test: agregar prueba para validar proyecto sin nombre
feat: implementar validación de proyecto
refactor: simplificar validación de proyecto
```

No modificar el historial únicamente para simular un proceso que no ocurrió.

---

## Calidad de los tests

Un buen test debe:

* comprobar comportamiento real;
* detectar regresiones;
* ser comprensible;
* fallar ante una implementación incorrecta;
* evitar depender de detalles internos innecesarios.

No escribir tests que simplemente repitan la implementación.

---

## Regla principal

El objetivo de TDD no es tener muchos tests.

El objetivo es utilizar tests para definir y proteger el comportamiento esperado antes y durante la implementación.

El agente debe respetar:

**RED → GREEN → REFACTOR**

y mantener los tests en verde al finalizar cada cambio.