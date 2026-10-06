# GitHub Skill

## Propósito

Gestionar la interacción del agente con el repositorio GitHub del proyecto.

La Skill debe permitir que el trabajo realizado por el agente quede correctamente versionado, trazable y publicado cuando la integración con GitHub esté disponible y autorizada.

Flujo principal:

**Issue → Branch → Desarrollo → Tests → Commit → Push → Pull Request → Review → Merge**

---

## Regla fundamental

Nunca afirmar que una operación fue realizada en GitHub si no fue realmente ejecutada y confirmada.

Por ejemplo:

* no afirmar que una Issue fue creada si no fue creada;
* no afirmar que un commit fue enviado si el push falló;
* no afirmar que una Pull Request existe si no fue creada;
* no afirmar que un cambio está en GitHub si solamente existe localmente.

---

## Fuentes de verdad

Antes de modificar el repositorio:

1. consultar el estado actual;
2. revisar la estructura existente;
3. revisar Issues relacionadas;
4. revisar ramas relevantes;
5. revisar cambios pendientes;
6. consultar `docs/enunciado.md`;
7. consultar `docs/contexto-proyecto.md`;
8. consultar `docs/reglas-agente.md`.

No asumir que el repositorio está limpio.

---

## Seguridad

Nunca subir:

* contraseñas;
* tokens;
* claves privadas;
* credenciales;
* secretos;
* archivos `.env` con información sensible.

Antes de hacer push revisar cambios sensibles cuando corresponda.

---

## Estado inicial

Antes de comenzar una tarea:

1. verificar rama actual;
2. verificar estado de Git;
3. actualizar información del repositorio;
4. verificar si existen cambios locales;
5. revisar Issues relacionadas;
6. identificar dependencias.

No sobrescribir cambios existentes sin autorización.

---

## Issues

Antes de crear una Issue:

1. buscar Issues similares;
2. verificar si la User Story ya tiene Issue;
3. evitar duplicados;
4. determinar título y descripción;
5. incluir criterios de aceptación;
6. relacionarla con la Épica/User Story cuando sea posible.

Una Issue debería permitir identificar claramente:

* qué debe hacerse;
* por qué;
* cómo saber cuándo está terminado.

---

## User Stories y Issues

Mantener la relación:

```text id="svk7yr"
Épica
  ↓
User Story
  ↓
Issue
```

La Issue debe poder rastrearse hasta el requisito de origen.

---

## Branches

Para cambios de funcionalidad, utilizar una rama apropiada.

Formato recomendado:

```text id="q4j4k8"
feature/<nombre-corto>
```

Para correcciones:

```text id="6c1g1v"
fix/<nombre-corto>
```

Para documentación:

```text id="v4m2nj"
docs/<nombre-corto>
```

No crear ramas innecesarias para cambios triviales cuando el flujo del proyecto no las requiera.

---

## Commits

Los commits deben ser:

* pequeños cuando sea razonable;
* coherentes;
* descriptivos;
* relacionados con una tarea;
* fáciles de revisar.

Ejemplos:

```text id="w5t0g9"
feat: agregar creación de proyectos
test: agregar validación de nombre de proyecto
fix: corregir cálculo de velocidad
docs: actualizar SDD de gestión de proyectos
refactor: simplificar validación de sprint
```

Evitar mensajes genéricos como:

```text id="5p6j7q"
update
changes
final
cosas
```

---

## Relación Issue → Commit

Cuando sea posible, relacionar commits con la Issue correspondiente.

Utilizar referencias de Issue según las capacidades disponibles de GitHub.

El objetivo es mantener:

```text id="v8x7c3"
Issue
  ↓
Commit
  ↓
Código
```

---

## Pull Requests

Cuando el flujo del proyecto utilice Pull Requests:

1. crear la rama;
2. implementar;
3. ejecutar tests;
4. revisar cambios;
5. crear commit;
6. hacer push;
7. crear Pull Request;
8. describir los cambios;
9. relacionar la Issue;
10. revisar;
11. fusionar cuando corresponda.

No crear Pull Requests innecesarias si el flujo acordado permite realizar directamente el cambio.

---

## Pull Request

La descripción debe incluir, cuando corresponda:

* objetivo;
* Issue relacionada;
* cambios realizados;
* tests ejecutados;
* resultado;
* consideraciones relevantes.

Ejemplo:

```text id="1p9q0r"
## Objetivo

Implementar la creación de proyectos.

## Issue

Relacionado con #XX.

## Cambios

- agregado dominio Project;
- agregada validación;
- agregado servicio de creación;
- agregados tests.

## Tests

go test ./...

Resultado: PASS
```

---

## Push

Antes de realizar push:

1. verificar los archivos modificados;
2. revisar el diff;
3. ejecutar tests;
4. verificar que no existan secretos;
5. confirmar el commit;
6. realizar push.

Después del push:

1. comprobar que el cambio fue publicado;
2. verificar la rama remota;
3. informar el resultado.

---

## Cambios locales

Si existen cambios locales que no pertenecen a la tarea actual:

* no sobrescribirlos;
* no eliminarlos;
* no incluirlos accidentalmente en el commit.

Separar los cambios cuando sea posible.

Si no es posible hacerlo de forma segura, informar la situación.

---

## GitHub Projects

Si el proyecto utiliza GitHub Projects y la integración disponible permite gestionarlo:

* mantener Issues organizadas;
* actualizar estados;
* asociar trabajo al Sprint;
* evitar duplicaciones;
* respetar la estructura definida por el equipo.

Si la herramienta disponible no permite modificar Projects, no simular esa operación.

---

## Releases

Crear releases únicamente cuando el proyecto lo requiera.

No crear versiones ficticias.

Los releases deben corresponder a versiones reales del software.

---

## CI

Cuando exista integración continua:

* revisar ejecuciones;
* identificar fallos;
* corregir problemas relevantes;
* no ignorar fallos de CI.

Una Pull Request no debería considerarse lista si los checks obligatorios fallan.

---

## Recuperación ante errores

Si una operación Git falla:

1. identificar el error;
2. determinar si es local, remoto o de permisos;
3. no repetir ciegamente la operación;
4. corregir la causa;
5. volver a intentar cuando sea seguro;
6. informar el resultado.

---

## Trazabilidad completa

Mantener la cadena:

```text id="g2q3pl"
Requisito
   ↓
Épica
   ↓
User Story
   ↓
Issue
   ↓
Branch
   ↓
SDD
   ↓
BDD
   ↓
Tests
   ↓
Commit
   ↓
Pull Request
   ↓
Merge
```

No todos los elementos necesitan existir para cada cambio, pero las relaciones relevantes deben mantenerse.

---

## Finalización

Una tarea puede considerarse publicada cuando:

* el código está implementado;
* los tests relevantes pasan;
* la documentación está actualizada;
* los cambios están committeados;
* el push fue exitoso;
* la Issue está actualizada;
* la trazabilidad está mantenida.

---

## Regla principal

GitHub no es solamente un lugar para almacenar código.

Debe utilizarse como evidencia del proceso de desarrollo, manteniendo:

**trabajo real + historial real + trazabilidad real**.

Nunca fabricar actividad para aparentar que se realizó un proceso que no ocurrió.