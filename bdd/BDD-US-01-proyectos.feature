Feature: Soporte mínimo de proyectos (US-01)

  Scenario: Crear y consultar un proyecto planificado
    Given que no existe el proyecto "Carva Demo"
    When se crea el proyecto con ese nombre
    Then el proyecto puede consultarse con el nombre "Carva Demo"
    And su estado es "Planificado"

  Scenario: Rechazar un nombre vacío
    Given que se intenta crear un proyecto con un nombre compuesto solo por espacios
    When se registra el proyecto
    Then la operación es rechazada por validación
