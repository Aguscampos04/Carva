Feature: Discutir diferencias y realizar nuevas rondas (US-14)
  Como integrante de un equipo
  Quiero revisar las diferencias reveladas y realizar otra ronda cuando sea necesario
  Para reconsiderar la estimación de una historia

  @BDD-US-14-01 @AC-01
  Scenario: Mantener consultables los resultados revelados de rondas anteriores
    Given que una historia tiene una ronda finalizada con sus votos y diferencias revelados
    When se inicia una ronda posterior para la historia
    Then los votos y las diferencias de la ronda anterior siguen consultables

  @BDD-US-14-02 @AC-02 @AC-03
  Scenario: Iniciar una ronda nueva después de revelar la anterior
    Given que la ronda más reciente de una historia está finalizada
    When un participante inicia otra ronda
    Then se crea una ronda abierta con número consecutivo
    And la ronda nueva no tiene votos registrados
    And los votos de la ronda anterior se conservan
    And cada participante de la ronda nueva solo ve su propio voto mientras siga abierta

  @BDD-US-14-03 @AC-04
  Scenario: Impedir una nueva ronda si la anterior sigue abierta
    Given que la ronda más reciente de una historia permanece abierta
    When se solicita iniciar otra ronda
    Then la operación informa un conflicto
    And la historia conserva una única ronda abierta
