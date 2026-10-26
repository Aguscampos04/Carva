Feature: Registrar estimación acordada (US-15)
  Como integrante del equipo
  Quiero registrar el valor unánime de una ronda finalizada
  Para actualizar los Story Points actuales de la historia

  @BDD-US-15-01 @AC-01
  Scenario: Registrar el valor de una ronda unánime
    Given que una historia tiene una ronda finalizada con todos los votos en 5 Story Points
    When un participante registra la estimación acordada
    Then la historia queda con 5 Story Points actuales
    And el valor provisional de la historia se conserva

  @BDD-US-15-02 @AC-01
  Scenario: Registrar cero cuando todos los participantes votaron cero
    Given que una historia tiene una ronda finalizada con todos los votos en 0 Story Points
    When un participante registra la estimación acordada
    Then la historia queda con 0 Story Points actuales

  @BDD-US-15-03 @AC-02
  Scenario: Conservar el Story Point actual si la ronda tiene diferencias
    Given que una historia tiene 3 Story Points actuales
    And una ronda finalizada contiene votos de 3 y 8 Story Points
    When un participante intenta registrar la estimación acordada
    Then la operación informa que no hay consenso
    And la historia conserva sus 3 Story Points actuales

  @BDD-US-15-04 @AC-03
  Scenario: No actualizar Story Points al revelar una ronda
    Given que una historia no tiene Story Points actuales
    And una ronda finalizada contiene votos unánimes
    When se consultan los resultados sin registrar el acuerdo
    Then la historia continúa sin Story Points actuales

  @BDD-US-15-05 @AC-03
  Scenario: No actualizar Story Points al iniciar otra ronda
    Given que una historia tiene 5 Story Points actuales
    And su ronda anterior está finalizada
    When se inicia una nueva ronda de estimación
    Then la historia conserva sus 5 Story Points actuales

  @BDD-US-15-06 @AC-03
  Scenario: Rechazar el registro antes de finalizar la ronda
    Given que una ronda abierta contiene votos unánimes hasta el momento
    When un participante intenta registrar la estimación acordada
    Then la operación informa que la ronda no está finalizada
    And la historia no cambia sus Story Points actuales

  @BDD-US-15-07
  Scenario: Rechazar el registro de un miembro ajeno a la ronda
    Given que una ronda finalizada unánime tiene sus participantes registrados
    When un miembro ajeno a la ronda intenta registrar el acuerdo
    Then la operación rechaza la solicitud
    And la historia no cambia sus Story Points actuales
