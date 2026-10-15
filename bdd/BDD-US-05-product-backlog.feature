# US-05 — Registrar y mantener elementos del Product Backlog
# Fuente: Issue #5, docs/enunciado.md §2.2 y SDD-US-05.

Feature: Registrar y mantener elementos del Product Backlog
  Como integrante del equipo
  Quiero registrar y actualizar User Stories del Product Backlog
  Para describir y organizar el trabajo

  Scenario: Crear una User Story sin Story Points provisionales
    Given que existe un proyecto
    And se envían título, descripción, prioridad válida y al menos un criterio no vacío
    When se registra la User Story en el Backlog del proyecto
    Then se persiste con UUID generado por el servidor
    And su estado es "Pendiente"
    And su Story Points provisional y definitivo no están informados

  Scenario: Crear una User Story con Story Points provisionales
    Given que existe un proyecto
    And se envían los campos requeridos y Story Points provisionales enteros no negativos
    When se registra la User Story
    Then el valor provisional queda persistido y separado del definitivo

  Scenario: Rechazar datos obligatorios inválidos al crear
    Given que existe un proyecto
    When se registra una User Story con título, descripción o criterio vacío, prioridad desconocida o Story Points negativo
    Then se informa un error de validación
    And no se crea el elemento

  Scenario: Actualizar los campos editables sin registrar consenso
    Given que existe una User Story pendiente en el Backlog del proyecto
    When se actualizan sus campos editables con datos válidos
    Then los cambios quedan persistidos
    And no se registra ni modifica una estimación definitiva

  Scenario: Mantener el estado bajo el flujo de US-09
    Given que existe una User Story que no está completada
    When se intenta cambiar su estado a "Completada" mediante la edición genérica
    Then US-05 no cambia su estado

  Scenario: Rechazar un proyecto inexistente al crear
    Given que no existe el proyecto indicado
    When se intenta registrar una User Story para ese proyecto
    Then se informa que el proyecto no fue encontrado
    And no se crea el elemento

  Scenario: Rechazar una User Story inexistente al actualizar
    Given que no existe el elemento indicado dentro del proyecto
    When se intenta actualizarlo
    Then se informa que el elemento no fue encontrado
    And no se crea ni modifica otro elemento

  Scenario: Informar una falla de persistencia
    Given que el almacenamiento no está disponible
    When se intenta registrar o actualizar una User Story
    Then se informa un error controlado
    And no se devuelve una operación exitosa
