# Ingeniería y Calidad de Software

## Trabajo Práctico Integrador 2026

### 1. Objetivo

Desarrollar una aplicación destinada a estimar, realizar el seguimiento y medir proyectos de software.

La solución deberá desarrollarse utilizando:

* Go (Golang)
* Scrum
* SDD (Specification Driven Development)
* BDD (Behavior Driven Development)
* TDD (Test Driven Development)
* Git
* Herramientas de Inteligencia Artificial como soporte al desarrollo

La solución podrá implementarse como aplicación web, de escritorio o de consola.

El núcleo de la aplicación y las reglas de negocio deberán estar implementados en Go.

---

## 2. Requerimientos funcionales

### 2.1 Gestión de proyectos

La aplicación deberá permitir:

* Crear proyectos.
* Modificar proyectos.
* Registrar miembros del equipo.
* Registrar fecha de inicio.
* Registrar fecha de finalización.
* Consultar el estado del proyecto.

### 2.2 Product Backlog

Cada elemento del Product Backlog deberá contar con:

* Identificador.
* Título.
* Descripción.
* Prioridad.
* Estado.
* Story Points.
* Criterios de aceptación.

### 2.3 Gestión de Sprints

La aplicación deberá permitir:

* Crear Sprints.
* Definir el Sprint Goal.
* Asignar historias a un Sprint.
* Registrar historias completadas.
* Cerrar un Sprint.
* Consultar Sprints anteriores.

### 2.4 Estimación

La aplicación deberá permitir trabajar con Story Points y Planning Poker.

Deberá contemplar:

* Estimaciones individuales.
* Ocultar las estimaciones hasta finalizar la votación.
* Mostrar las estimaciones una vez finalizada la votación.
* Detectar diferencias entre estimaciones.
* Realizar nuevas rondas de votación.
* Registrar la estimación acordada.

### 2.5 Esfuerzo

Se deberá poder registrar:

* Miembro del equipo.
* Fecha.
* Actividad.
* Horas trabajadas.

La aplicación deberá permitir comparar el esfuerzo estimado con el esfuerzo real.

### 2.6 Defectos

Los defectos deberán contar con:

* Descripción.
* Severidad.
* Estado.
* Historia relacionada.
* Sprint de detección.
* Sprint de resolución.

### 2.7 Métricas

La aplicación deberá contemplar, como mínimo:

* Story Points planificados.
* Story Points completados.
* Velocidad del equipo.
* Horas estimadas.
* Horas reales.
* Desviación entre esfuerzo estimado y real.
* Porcentaje de historias completadas.
* Defectos detectados.
* Defectos resueltos.

### 2.8 Dashboard

La aplicación deberá contar con un dashboard que permita consultar el estado del proyecto y visualizar gráficos.

### 2.9 Reportes

Se deberán poder generar reportes de proyecto o de Sprint que incluyan:

* Historias planificadas.
* Historias completadas.
* Estimaciones.
* Esfuerzo.
* Métricas.
* Defectos.

---

## 3. Scrum

Se deberán contemplar los siguientes roles:

### Product Architect

Corresponde a los profesores.

### Agile Enabler

Será uno de los integrantes del equipo.

### Product Builders

Serán los integrantes del equipo.

El trabajo deberá organizarse utilizando:

* Product Backlog.
* Sprints.
* Sprint Planning.
* Daily.
* Sprint Review.
* Sprint Retrospective.

---

## 4. SDD — Specification Driven Development

Las funcionalidades principales deberán especificarse antes de su implementación.

Cada especificación deberá contemplar:

* Objetivo.
* Entradas.
* Salidas esperadas.
* Reglas de negocio.
* Restricciones.
* Casos borde.
* Condiciones de error.
* Criterios de aceptación.

Las especificaciones deberán estar versionadas junto con el proyecto.

---

## 5. BDD — Behavior Driven Development

Las funcionalidades deberán describirse utilizando escenarios en formato:

**Given – When – Then**

Cuando corresponda, se deberán contemplar:

* Escenarios normales.
* Escenarios alternativos.
* Escenarios de borde.
* Escenarios de error.

Los escenarios deberán automatizarse cuando sea posible.

---

## 6. TDD — Test Driven Development

El desarrollo deberá seguir el ciclo:

**RED → GREEN → REFACTOR**

Como mínimo deberán existir pruebas unitarias para:

* Métricas.
* Cálculos de estimación.
* Reglas de negocio.
* Validaciones.

Se deberá contar con evidencia del proceso TDD mediante el historial del repositorio.

---

## 7. Uso de Inteligencia Artificial

Se permite utilizar Inteligencia Artificial como soporte para:

* Análisis de requerimientos.
* Especificaciones.
* Generación de código.
* Generación de pruebas.
* Refactorización.
* Revisión de código.
* Documentación.

Los resultados generados por IA deberán ser:

* Comprendidos por el equipo.
* Revisados.
* Validados.

El equipo será responsable de todo el código incorporado al proyecto.

---

## 8. Plan de trabajo

### Sprint 0 — Preparación

Se deberá realizar:

* Conformación de equipos.
* Definición de roles.
* Preparación del entorno.
* Creación del repositorio GitHub/GitLab.
* Configuración del tablero de GitHub Projects.
* Reunión inicial con el profesor.
* Creación del Product Backlog inicial.
* Creación de las primeras User Stories junto con el cliente.

### Sprint 1 — MVP

Desarrollo de una versión funcional básica.

Deberán realizarse las ceremonias:

* Planning.
* Daily.
* Review.
* Retrospective.

### Sprint 2 — Interfaz

Desarrollo de una interfaz usable.

### Sprint 3 — Funcionalidad y Calidad

Implementación de las funcionalidades principales y mejora de la robustez del sistema.

También se deberá implementar la visualización.

### Sprint 4 — Cierre y entrega final

Se deberá:

* Realizar el pulido final.
* Completar la documentación.
* Preparar la entrega.
* Exportar el informe del proyecto a PDF.
* Realizar la Sprint Review final y presentación formal.

---

## 9. Trazabilidad

Se deberá mantener trazabilidad entre:

**User Story → SDD Specification → Acceptance Criteria → BDD Scenarios → Tests → Go Code**

Al menos una funcionalidad deberá demostrar el recorrido completo durante la presentación final.

---

## 10. Entregables

El proyecto deberá contar con:

* Repositorio Git con historial de contribuciones.
* Tablero Scrum.
* Product Backlog.
* Sprint Backlogs.
* SDD.
* BDD.
* Pruebas automatizadas.
* Código fuente en Go.
* Evidencia de TDD.
* Software funcional.
* Métricas.
* Reporte de cobertura de tests.
* Actas de retrospectivas.
* Documentación técnica.
* Breve manual de usuario.
* Presentación y demostración final.

---

## 11. Evaluación

La evaluación contempla:

| Criterio                         | Peso |
| -------------------------------- | ---: |
| Producto funcional               |  25% |
| SDD / BDD / TDD                  |  25% |
| Calidad del software             |  20% |
| Gestión del proyecto             |  20% |
| Trabajo en equipo y presentación |  10% |

---

## 12. Fuente

Este documento corresponde al enunciado oficial del Trabajo Práctico Integrador 2026 de Ingeniería y Calidad de Software.

La versión original del enunciado es la fuente de verdad ante cualquier diferencia o interpretación.
