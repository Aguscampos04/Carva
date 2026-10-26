Feature: Revelar votos y detectar diferencias (US-13)
  Como integrante de un equipo
  Quiero revelar los votos al finalizar una ronda
  Para conocer las diferencias y discutir la estimación

  @BDD-US-13-01 @AC-01
  Scenario: Mantener ocultos los votos ajenos mientras haya votos pendientes
    Given que una ronda abierta tiene participantes y votos registrados
    And al menos un participante todavía no votó
    When un participante consulta sus votos visibles
    Then solo ve su propio voto
    And la ronda continúa abierta

  @BDD-US-13-02 @AC-01 @AC-03
  Scenario: Revelar todos los votos cuando todos votaron y se finaliza la ronda
    Given que todos los participantes de una ronda abierta registraron un voto
    When un participante finaliza la ronda
    Then la ronda queda finalizada
    And cada participante puede consultar simultáneamente todos los votos

  @BDD-US-13-03 @AC-02
  Scenario: Informar que no hay diferencias cuando todos votaron el mismo valor
    Given que todos los participantes votaron 5 Story Points
    When se consulta el resultado de la ronda finalizada
    Then se muestran todos los votos
    And el resultado indica que no hay diferencias

  @BDD-US-13-04 @AC-02
  Scenario: Informar diferencias cuando los votos no son iguales
    Given que los participantes votaron 3 y 8 Story Points
    When se consulta el resultado de la ronda finalizada
    Then se muestran ambos votos
    And el resultado indica que hay diferencias

  @BDD-US-13-05 @AC-03
  Scenario: Impedir la finalización si falta algún voto
    Given que una ronda abierta tiene al menos un participante sin voto
    When un participante intenta finalizar la ronda
    Then la operación informa un conflicto
    And la ronda permanece abierta
    And los votos ajenos siguen ocultos

  @BDD-US-13-06 @AC-01
  Scenario: Impedir consultar el resultado antes de revelar la ronda
    Given que existe una ronda abierta
    When un participante consulta el resultado
    Then la operación informa un conflicto
    And no se revelan votos ajenos
