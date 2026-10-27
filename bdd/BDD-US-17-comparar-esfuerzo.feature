Feature: Comparar esfuerzo estimado y real de una User Story
  Como integrante del equipo
  Quiero comparar horas estimadas con horas reales registradas
  Para identificar ambos valores de esfuerzo por historia

  # AC-01, AC-02
  Scenario: Mostrar estimación y suma real por separado sin convertir Story Points
    Given que la User Story "A" tiene 2.50 horas estimadas y 5 Story Points
    And "A" tiene registros activos de esfuerzo de 1.25 y 2.50 horas
    And otra User Story "B" tiene un registro activo de 7 horas
    When el integrante consulta la comparación de "A"
    Then se muestran por separado 2.50 horas estimadas y 3.75 horas reales
    And las 7 horas de "B" no forman parte del total de "A"
    And los 5 Story Points de "A" se conservan como dato independiente
    And ninguna estimación ni registro se modifica

  # AC-01, AC-02; caso límite de la fuente de datos US-16
  Scenario: Excluir un registro eliminado de las horas reales
    Given que la User Story "A" tiene 2.50 horas estimadas
    And tiene un registro activo de 1.25 horas y otro registro eliminado de 4 horas
    When el integrante consulta la comparación de "A"
    Then se muestran 2.50 horas estimadas y 1.25 horas reales por separado

  # AC-01; historia histórica permitida por US-16
  Scenario: Mostrar horas reales cuando falta la estimación histórica
    Given que la User Story "A" está pendiente de regularizar sus horas estimadas
    And tiene un registro activo de 1.25 horas
    When el integrante consulta la comparación de "A"
    Then se muestra "Pendiente de regularización" para las horas estimadas
    And se muestran 1.25 horas reales por separado
