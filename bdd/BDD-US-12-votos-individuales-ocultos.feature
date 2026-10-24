Feature: Registrar votos individuales ocultos (US-12)
  Como miembro identificado
  Quiero registrar mi estimación individual sin ver los votos ajenos durante la votación
  Para estimar de forma independiente

  # Los participantes son los miembros registrados cuando se abre la ronda.

  @BDD-US-12-01 @AC-01
  Scenario: Atribuir un voto al miembro, historia y ronda correctos
    Given que existe una User Story con una ronda abierta
    And existe un miembro del mismo proyecto
    When el miembro registra un valor de Story Points para esa ronda
    Then el voto queda asociado a ese miembro, esa User Story y esa ronda

  @BDD-US-12-02 @AC-02
  Scenario: Mantener ocultos los votos de otros miembros durante una ronda abierta
    Given que existe una ronda abierta con votos registrados por distintos miembros
    When un miembro consulta el estado de votación antes de finalizar la ronda
    Then no puede ver los valores individuales de los otros miembros

  @BDD-US-12-03 @AC-01
  Scenario: Rechazar un miembro de otro proyecto
    Given que existe una User Story en un proyecto
    And existe un miembro asociado a otro proyecto
    When se intenta registrar el voto de ese miembro para la User Story
    Then el voto se rechaza
    And no se registra ningún voto

  @BDD-US-12-04 @AC-01
  Scenario: Permitir corregir el propio voto antes de finalizar la ronda
    Given que un miembro ya registró su voto en una ronda abierta
    When ese mismo miembro registra otro valor para la misma ronda
    Then el nuevo valor reemplaza al anterior
    And el voto sigue oculto para los otros participantes

  @BDD-US-12-05 @AC-01
  Scenario: Incluir a los miembros registrados al abrir la ronda
    Given que un proyecto tiene miembros registrados
    When se abre la ronda inicial de estimación para una User Story
    Then todos los miembros registrados quedan como participantes de esa ronda
