# Evidencia de desarrollo asistido por IA — US-20

## Objetivo e instrucción

Implementar la Issue [#20 — Consultar métricas de Story Points y avance](https://github.com/Aguscampos04/Carva/issues/20) como parte del trabajo indicado por Loana Vargas para el Sprint 4, y preparar un Pull Request para revisión y merge de un compañero.

## Decisiones aprobadas

Loana aprobó el 2026-10-27:

- Story Points planificados por Sprint: suma de los Story Points definitivos de los snapshots originales.
- Story Points completados: suma de los puntos de las historias completadas.
- Porcentaje por Sprint: historias completadas / historias planificadas × 100; sin historias, «No disponible».
- Velocidad de un Sprint cerrado: puntos completados. Velocidad del proyecto: promedio de puntos completados de Sprints cerrados.
- Agregado del porcentaje del proyecto: proporción ponderada entre todos los snapshots.
- Agregado de puntos del proyecto: suma entre todos los snapshots; una historia trasladada cuenta en cada Sprint que la planificó.
- Las métricas de puntos que requieran una estimación faltante se presentan como «No disponible»; en UI se representa API `null` con ese texto.

Las demás reglas y criterios se derivan del enunciado, contexto aprobado e Issue #20. Las decisiones pendientes de aprobación no se inventan.

## Cambios

- SDD y BDD específicos de US-20.
- Dominio Go para las métricas, aplicación, repositorio SQLite de solo lectura y endpoints REST por proyecto/Sprint.
- Interfaz web con resumen por proyecto, selector de Sprint y presentación de «No disponible».
- Sin migraciones ni dependencias nuevas.

## Evidencia TDD

- **RED (dominio):** se ejecutó `go test ./internal/domain/metrics` después de escribir las pruebas y antes de implementar el paquete. Falló al compilar porque faltaban `StorySnapshot`, `StoryPointMetrics` y las funciones de cálculo.
- **GREEN (dominio):** tras implementar el cálculo, las pruebas del paquete verificaron valores de cero, puntos faltantes, porcentaje ponderado y velocidad por Sprints cerrados; finalizaron en PASS.
- Las pruebas de integración y de UI se añadieron para validar el contrato REST/SQLite y el comportamiento visible. No se documenta una etapa RED separada para esas pruebas.

## Revisión humana

La aprobación de las fórmulas por Loana quedó registrada durante esta conversación. La revisión de código del equipo queda pendiente; no se afirma que ya haya ocurrido.

## Validación y publicación

- `go test ./...`: PASS; todos los paquetes, incluida la integración.
- `go vet ./...`: PASS, sin diagnósticos.
- `go build ./...`: PASS (exit code 0); Go mostró una advertencia no fatal al escribir una caché de estadísticas en el directorio de módulos externo.
- `go mod verify`: PASS; todos los módulos verificados.
- `node --test tests/web/*.test.cjs`: PASS, 39/39.
- `gofmt -l .`: PASS; no reportó archivos.
- `git diff --check`: PASS.

El commit, el push y el Pull Request se registrarán después de ejecutarlos y verificar su estado real. No se hará merge ni se cerrará la Issue desde esta tarea.
