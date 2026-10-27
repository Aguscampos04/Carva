# SDD-US-17 — Comparar esfuerzo estimado y real

**Issue:** [#17 — US-17](https://github.com/Aguscampos04/Carva/issues/17)  
**Épica:** EP-05 — Seguimiento de esfuerzo  
**Requisito:** `docs/enunciado.md` §2.5 y §2.7  
**Dependencia:** US-16 — horas estimadas y registros de esfuerzo real

## 1. Objetivo

Permitir que un integrante compare, para una User Story del Product Backlog, las horas estimadas con las horas reales registradas para esa misma historia. La comparación hace visibles ambos datos con sus unidades y fuentes diferenciadas.

## 2. Inputs

| Entrada | Origen | Restricción |
| --- | --- | --- |
| Identificador de proyecto e historia | Historia seleccionada en el Product Backlog | Deben identificar recursos del mismo proyecto; un recurso inexistente conserva el error 404 existente. |
| Horas estimadas | Campo de planificación de US-16 | Se muestran con hasta dos decimales. Puede no existir en historias históricas pendientes de regularización. |
| Registros de esfuerzo | Registros US-16 vinculados a la historia | El total usa horas de registros activos; cada valor admite centésimas. |

No se reciben ni se convierten Story Points como entrada de horas.

## 3. Expected Outputs

El detalle de la historia presenta dos valores distinguibles:

- **Horas estimadas:** valor de planificación de la historia, o «Pendiente de regularización» si no existe.
- **Horas reales:** suma de las horas de los registros de esfuerzo activos relacionados con esa historia.

El detalle del esfuerzo sigue permitiendo consultar los registros individuales mediante la funcionalidad US-16.

## 4. Business Rules

- **BR-01:** La estimación se lee del campo de horas estimadas de la User Story.
- **BR-02:** Las horas reales se obtienen sumando los registros activos asociados a esa historia, sin incluir registros dados de baja.
- **BR-03:** Las horas estimadas y reales se muestran por separado y no se sobrescriben entre sí.
- **BR-04:** Story Points no se interpretan ni convierten a horas.
- **BR-05:** Esta historia no calcula una desviación ni agrega valores por Sprint/proyecto. Las fórmulas de esas métricas corresponden a US-21 y deben definirse en su SDD.

## 5. Constraints

- Reutilizar las fuentes de datos de US-16; no duplicar ni volver a registrar esfuerzo.
- Conservar precisión de centésimas implementada por US-16.
- No modificar estimaciones, Story Points, registros, estados ni asociaciones al consultar la comparación.
- No presentar un total real de una historia distinta.

## 6. Edge Cases

- Sin horas estimadas: indicar «Pendiente de regularización», manteniendo visible el total real.
- Sin registros activos: mostrar cero horas reales según la representación existente de US-16.
- Varios registros activos para la historia: mostrar su suma.
- Registros eliminados lógicamente: excluirlos de la suma.
- Story Points presentes: conservarlos como dato independiente y no usarlos en el total horario.

## 7. Error Conditions

- Proyecto o User Story inexistente o ajeno al proyecto: conservar el 404 del detalle del Product Backlog.
- Fallo de consulta: conservar el manejo de error existente de la interfaz/API; no mostrar un valor de comparación inventado.

## 8. Acceptance Criteria

- **AC-01:** Para una historia con estimación y registros de esfuerzo, se muestran por separado las horas estimadas y las reales correspondientes.
- **AC-02:** La comparación utiliza los registros de esfuerzo relacionados con esa historia, sin tratar los Story Points como horas.

Trazabilidad: §2.5/§2.7 → EP-05 → US-17 / #17 → esta SDD → `bdd/BDD-US-17-comparar-esfuerzo.feature` → pruebas web e integración → UI y consulta del Product Backlog.
