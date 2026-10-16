Feature: Crear Sprint y definir Sprint Goal (US-07)

  Scenario: Crear y consultar un Sprint asociado al proyecto
    Given que existe el proyecto "Carva Demo" en estado "Planificado"
    When se crea un Sprint con el objetivo "Entregar el MVP"
    Then el Sprint queda asociado al proyecto "Carva Demo"
    And puede consultarse con el objetivo "Entregar el MVP"

  Scenario: El primer Sprint inicia el proyecto
    Given que existe el proyecto "Carva Demo" en estado "Planificado" y sin Sprints
    When se crea el primer Sprint con el objetivo "Entregar el MVP"
    Then el estado del proyecto pasa a "En progreso"

  Scenario: Rechazar un Sprint Goal vacío sin cambios parciales
    Given que existe el proyecto "Carva Demo" en estado "Planificado"
    When se intenta crear un Sprint con un objetivo compuesto solo por espacios
    Then la operación es rechazada por validación
    And el proyecto conserva el estado "Planificado"

  Scenario: Rechazar un Sprint para un proyecto inexistente
    Given que no existe el proyecto con identificador 999
    When se intenta crear un Sprint con un objetivo válido
    Then la operación responde que el proyecto no fue encontrado
    And no se crea ningún Sprint

  Scenario: No consultar un Sprint bajo otro proyecto
    Given que un Sprint pertenece al proyecto "Carva Demo"
    When se consulta usando el identificador de un proyecto diferente
    Then el Sprint no se encuentra bajo ese proyecto
