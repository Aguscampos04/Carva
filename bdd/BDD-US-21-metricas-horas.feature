Feature: Consultar métricas de horas
  Como integrante del equipo
  Quiero consultar horas estimadas, reales y su desviación
  Para medir el esfuerzo por Sprint y proyecto

  @BDD-US-21-01 @AC-01
  Scenario: Consultar métricas del Sprint con centésimas exactas
    Given que un Sprint contiene historias con 2.50 y 3.25 horas estimadas
    And tiene registros activos asociados de 1.25 y 5.00 horas
    When se consultan las métricas del Sprint
    Then las horas estimadas son 5.75
    And las horas reales son 6.25
    And la desviación es 0.50

  @BDD-US-21-02 @AC-05
  Scenario: Mostrar desviación negativa cuando el esfuerzo real es menor
    Given que un Sprint tiene 8.00 horas estimadas y 5.50 horas reales
    When se consultan sus métricas
    Then la desviación es -2.50 horas

  @BDD-US-21-03 @AC-02
  Scenario: Agregar proyecto sin duplicar historias trasladadas
    Given que una historia con estimación de 4.00 horas aparece en dos snapshots de Sprint
    And otra historia del Backlog tiene 3.00 horas estimadas
    And los registros activos del proyecto suman 9.25 horas, incluyendo 1.25 sin Sprint
    When se consultan las métricas del proyecto
    Then las horas estimadas son 7.00
    And las horas reales son 9.25
    And la desviación es 2.25

  @BDD-US-21-04 @AC-03
  Scenario: No inventar horas estimadas cuando una historia no tiene estimación
    Given que una historia del Sprint no tiene horas estimadas
    And los registros activos del Sprint suman 2.00 horas
    When se consultan sus métricas
    Then las horas estimadas y la desviación son «No disponible»
    And las horas reales son 2.00

  @BDD-US-21-05 @AC-04
  Scenario: Excluir entradas eliminadas y mantener la consulta de solo lectura
    Given que hay un registro activo de 1.50 horas y otro eliminado de 4.00 horas
    When se consultan las métricas del Sprint y del proyecto
    Then las horas reales solo incluyen 1.50 horas
    And la consulta no cambia los registros almacenados

  @BDD-US-21-06 @AC-06
  Scenario Outline: Informar identificadores inválidos o recursos inexistentes
    Given que existe un proyecto válido
    When se consultan métricas de Sprint con el identificador «<id>»
    Then la respuesta HTTP es <status>

    Examples:
      | id       | status |
      | abc      | 400    |
      | 0        | 400    |
      | 999999   | 404    |
