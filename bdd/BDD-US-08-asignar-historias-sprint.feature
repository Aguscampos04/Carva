# US-08 / Issue #8 — Asignar historias y preservar lo planificado
# Fuente: SDD-US-08 y criterios de aceptación de Issue #8.

Feature: Asignar historias del Product Backlog a un Sprint

  Scenario: Consultar una historia planificada en su Sprint
    Given que existe un Sprint y una historia del Backlog del mismo proyecto
    When se asigna la historia al Sprint
    Then la historia aparece en la planificación consultada del Sprint
    And se conserva su identificador y sus datos de planificación

  Scenario: Preservar el snapshot ante una edición posterior del Backlog
    Given que una historia fue asignada al Sprint con título "Plan inicial"
    When el título y los criterios originales se actualizan en el Backlog
    Then la planificación del Sprint conserva el título "Plan inicial" y sus criterios originales
    And el Backlog muestra los datos actualizados

  Scenario: Rechazar una historia que no pertenece al proyecto del Sprint
    Given que existe un Sprint en un proyecto y una historia en otro proyecto
    When se intenta asignar esa historia al Sprint
    Then la operación indica que el recurso no fue encontrado
    And la planificación del Sprint no cambia

  Scenario: Rechazar una asignación duplicada en el mismo Sprint
    Given que una historia ya está asignada a un Sprint
    When se intenta asignar nuevamente la misma historia al mismo Sprint
    Then la operación indica un conflicto
    And la historia aparece una sola vez en la planificación
