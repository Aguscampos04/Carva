# Contexto del Proyecto

## 1. Identificación

**Proyecto:** Aplicación web para gestión, estimación, seguimiento y medición de proyectos de software.

**Materia:** Ingeniería y Calidad de Software.

**Año:** 2026.

## 2. Objetivo del proyecto

Desarrollar una aplicación web que permita gestionar proyectos de software y sus elementos asociados, incluyendo Product Backlog, Sprints, estimaciones, esfuerzo, defectos y métricas.

La aplicación deberá cumplir con los requerimientos funcionales y metodológicos establecidos en el enunciado oficial.

## 3. Tipo de solución

Se decidió desarrollar una **aplicación web**.

Esta decisión ya está tomada y no debe ser modificada por el agente salvo que el equipo lo solicite explícitamente.

## 4. Tecnologías y prácticas obligatorias

El proyecto deberá utilizar:

* Go (Golang) para el núcleo de la aplicación y las reglas de negocio.
* Scrum como marco de trabajo.
* SDD para especificar funcionalidades antes de implementarlas.
* BDD para definir el comportamiento mediante escenarios.
* TDD para desarrollar y validar las reglas y funcionalidades mediante pruebas.
* Git para el control de versiones.
* Herramientas de Inteligencia Artificial como soporte al desarrollo.

## 5. Principios de trabajo

El proyecto se desarrollará de forma incremental.

No se deberá implementar toda la aplicación de una sola vez.

El trabajo se organizará mediante:

1. Sprint 0 — Preparación.
2. Sprint 1 — MVP.
3. Sprint 2 — Interfaz.
4. Sprint 3 — Funcionalidad y Calidad.
5. Sprint 4 — Cierre y entrega final.

El avance entre etapas deberá realizarse de forma controlada y con revisión del equipo.

## 6. Arquitectura

La arquitectura concreta de la aplicación todavía no está cerrada.

El agente deberá proponer una arquitectura adecuada para una aplicación web cuyo núcleo y reglas de negocio estén implementados en Go.

Antes de comenzar la implementación deberá existir una decisión arquitectónica documentada.

La arquitectura deberá favorecer:

* Separación de responsabilidades.
* Testabilidad.
* Mantenibilidad.
* Trazabilidad.
* Evolución incremental.
* Claridad de las reglas de negocio.

## 7. Funcionalidades principales

La solución deberá contemplar, como mínimo:

* Gestión de proyectos.
* Gestión de miembros.
* Product Backlog.
* Gestión de Sprints.
* Estimación mediante Story Points.
* Planning Poker.
* Registro y comparación de esfuerzo.
* Gestión de defectos.
* Métricas.
* Dashboard.
* Reportes.

El detalle de cada funcionalidad deberá derivarse del enunciado oficial y posteriormente formalizarse mediante User Stories, SDD, BDD y TDD.

## 8. Gestión del proyecto

El proyecto utilizará un repositorio Git y un tablero Scrum.

Las unidades principales de trabajo serán:

* Epic.
* User Story.
* Issue.
* Sprint.
* Task, cuando sea necesario.

Cada elemento deberá mantener relaciones con los artefactos correspondientes.

## 9. Trazabilidad

Se deberá mantener la siguiente cadena de trazabilidad:

**Requerimiento → Epic → User Story → Issue → SDD → Criterios de aceptación → BDD → Test → Código → Commit/PR**

No será necesario que todas las decisiones se resuelvan de antemano.

La trazabilidad se construirá progresivamente durante el desarrollo.

Como mínimo, una funcionalidad deberá demostrar el recorrido completo durante la presentación final.

## 10. Uso del agente de IA

El agente será utilizado como asistente técnico y de gestión del proyecto.

Podrá colaborar en:

* Análisis de requerimientos.
* Identificación y refinamiento de User Stories.
* Creación y actualización de Issues.
* Elaboración de especificaciones SDD.
* Elaboración de escenarios BDD.
* Diseño y seguimiento de TDD.
* Generación y revisión de código Go.
* Generación y revisión de tests.
* Análisis de cobertura.
* Revisión de calidad.
* Control de trazabilidad.
* Documentación.
* Organización Scrum.

El agente no reemplaza la responsabilidad del equipo.

Todo resultado generado por IA deberá ser comprendido, revisado y validado antes de incorporarse al proyecto.

## 11. Regla de aprobación

El agente podrá proponer soluciones, estructuras, historias, especificaciones, código o modificaciones.

Sin embargo, no deberá considerar una propuesta como decisión definitiva hasta que el equipo la apruebe cuando dicha decisión afecte significativamente:

* Arquitectura.
* Requerimientos.
* Reglas de negocio.
* Alcance.
* Tecnologías.
* Organización del proyecto.

## 12. Estado actual

### Decisiones tomadas

* El proyecto será una aplicación web.
* Se utilizará Go para el núcleo y las reglas de negocio.
* Se utilizará Scrum.
* Se utilizarán SDD, BDD y TDD.
* Se utilizará Git.
* Se utilizará IA como soporte.
* El proyecto se desarrollará incrementalmente mediante Sprints.

### Decisiones pendientes

Todavía deben definirse y documentarse, entre otras:

* Arquitectura concreta.
* Tecnología del frontend.
* Tecnología de persistencia/base de datos.
* Organización definitiva de carpetas.
* API y contratos.
* Modelo de dominio.
* Estrategia de autenticación, si corresponde.
* Estrategia de despliegue, si corresponde.
* Diseño visual de la interfaz.

Estas decisiones deberán tomarse antes de implementarlas y quedar documentadas cuando corresponda.

## 13. Regla fundamental

El enunciado oficial del TPI es la fuente de verdad para los requisitos obligatorios.

Cuando exista una diferencia entre una propuesta del agente y el enunciado, deberá prevalecer el enunciado.

El agente no deberá inventar requisitos obligatorios que no estén respaldados por el enunciado o por una decisión explícita del equipo.
