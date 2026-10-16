# US-09 — Registrar historias completadas
# Fuente: Issue #9, SDD-US-09 y decisión de confirmación única de Loana Vargas (2026-10-16).

Feature: Registrar una User Story como completada
  Como integrante del equipo
  Quiero registrar la finalización de una historia después de verificar sus criterios
  Para reflejar el trabajo efectivamente terminado

  # AC-01 → BDD-01 → tests de dominio, aplicación e integración HTTP.
  Scenario: Rechazar la finalización si los criterios no fueron verificados
    Given que existe una historia en estado "En progreso"
    And no se verificaron todos sus Acceptance Criteria
    When se intenta registrar la historia como completada
    Then el sistema rechaza la operación
    And la historia conserva el estado "En progreso"

  # AC-02 → BDD-02 → tests de aplicación e integración HTTP/SQLite.
  Scenario: Completar una historia con confirmación de todos sus criterios
    Given que existe una historia en el Product Backlog
    And se verificaron todos sus Acceptance Criteria
    When se registra su finalización
    Then la historia queda en estado "Completada"
    And el Product Backlog devuelve la historia con ese estado

  # BR-04 → BDD-03 → test de integración HTTP de US-05/US-09.
  Scenario: La edición genérica no permite completar una historia
    Given que existe una historia que no está completada
    When se edita la historia enviando el estado "Completada" por el endpoint genérico
    Then la historia conserva su estado anterior

  # BR-05 / dependencia US-08; validación de Sprint pendiente de integración.
  Scenario: Consultar como completada una historia planificada en un Sprint
    Given que una historia está asignada a un Sprint
    And se verificaron todos sus Acceptance Criteria
    When se consulta la historia desde el Sprint después de registrarla como completada
    Then la historia figura como "Completada"
