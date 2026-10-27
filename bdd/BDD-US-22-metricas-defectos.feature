# US-22 / Issue #22 — Consultar métricas de defectos
# Fuente: criterios de aceptación de Issue #22 y SDD-US-22.

Feature: Consultar defectos detectados y resueltos por Sprint y proyecto

  @AC-01 @AC-02
  Scenario: Atribuir detección y resolución a Sprints distintos
    Given que un defecto del proyecto fue detectado en el Sprint 1
    And el defecto está resuelto y su Sprint de resolución es el Sprint 2
    When se consultan las métricas del proyecto
    Then el Sprint 1 muestra 1 defecto detectado y 0 resueltos
    And el Sprint 2 muestra 0 defectos detectados y 1 resuelto
    And el agregado del proyecto muestra 1 defecto detectado y 1 resuelto

  @AC-01 @AC-05
  Scenario: Incluir defectos sin Sprint en el agregado del proyecto
    Given que existe un defecto detectado sin Sprint asociado
    And existe un defecto cerrado con Sprint de resolución registrado
    When se consultan las métricas del proyecto
    Then el agregado muestra 2 defectos detectados y 1 resuelto
    And la métrica de defectos detectados por Sprint muestra «No disponible»

  @AC-03
  Scenario: No atribuir resoluciones sin Sprint de resolución
    Given que existe un defecto en estado Resuelto sin Sprint de resolución
    When se consultan las métricas del proyecto
    Then el agregado muestra 1 defecto resuelto
    And la métrica de defectos resueltos por Sprint muestra «No disponible»

  @AC-04
  Scenario: Mostrar cero cuando la atribución está completa y no hay coincidencias
    Given que todos los defectos del proyecto tienen Sprint de detección
    And todos los defectos Resueltos o Cerrados tienen Sprint de resolución
    And ningún defecto fue detectado ni resuelto en el Sprint 3
    When se consultan las métricas del Sprint 3
    Then la métrica de defectos detectados muestra 0
    And la métrica de defectos resueltos muestra 0

  @AC-01 @AC-05
  Scenario: Consultar un proyecto sin defectos ni Sprints
    Given que existe un proyecto sin defectos ni Sprints
    When se consultan las métricas del proyecto
    Then el agregado muestra 0 defectos detectados y 0 resueltos
    And la lista de Sprints está vacía

  @AC-05
  Scenario: Rechazar recursos inexistentes
    Given que existe un proyecto
    When se consulta un proyecto inexistente o un Sprint ajeno al proyecto
    Then la API responde recurso no encontrado

  @AC-05
  Scenario: Rechazar un identificador de Sprint inválido
    Given que existe un proyecto
    When se consulta con un ID de Sprint no entero positivo
    Then la API responde solicitud inválida
