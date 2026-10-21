# US-09 / Issue #9
Feature: Registrar historias completadas

  Scenario: Rechazar finalización sin confirmación de verificación
    Given que una historia del Product Backlog está asignada a un Sprint
    When se intenta completarla sin confirmar que todos sus Acceptance Criteria fueron verificados
    Then la operación es rechazada
    And la historia conserva su estado en el Sprint y en el Backlog

  Scenario: Rechazar un cuerpo sin el campo de verificación
    Given que una historia del Product Backlog está asignada a un Sprint
    When se solicita completarla sin enviar acceptance_criteria_verified
    Then la solicitud es rechazada como inválida
    And la historia conserva su estado en el Sprint y en el Backlog

  Scenario: Completar una historia verificada
    Given que una historia del Product Backlog está asignada a un Sprint
    And sus datos de planificación están guardados en el Sprint
    When un integrante confirma que todos sus Acceptance Criteria fueron verificados y registra su finalización
    Then la historia figura como Completada en el Backlog
    And figura como Completada en el estado actual del Sprint
    And los campos originales de planificación del Sprint se conservan

  Scenario: Completar una historia que no está asignada al Sprint
    Given que una historia pertenece al proyecto pero no está asignada al Sprint
    When se intenta registrarla como completada desde ese Sprint
    Then la operación es rechazada
    And el estado del Backlog no cambia
