Feature: Consultar estado del proyecto
  Como integrante del equipo
  Quiero consultar el estado del proyecto
  Para conocer su situación vigente en el flujo de trabajo

  @AC-01 @AC-02
  Scenario: Consultar un proyecto todavía planificado
    Given que existe el proyecto "Carva Demo" sin Sprints y sin finalizar
    When consulto el proyecto por su identificador
    Then recibo el estado "Planificado"

  @AC-01 @AC-03
  Scenario: Consultar un proyecto después de iniciar su primer Sprint
    Given que existe el proyecto "Carva Demo" sin Sprints y sin finalizar
    When se crea su primer Sprint y consulto el proyecto
    Then recibo el estado "En progreso"

  @AC-01 @AC-04
  Scenario: Consultar un proyecto finalizado
    Given que existe el proyecto "Carva Demo"
    And se finaliza manualmente con confirmación
    When consulto el proyecto por su identificador
    Then recibo el estado "Finalizado"

  @AC-05
  Scenario: Consultar un proyecto inexistente
    Given que no existe un proyecto con el identificador consultado
    When consulto ese proyecto
    Then recibo una respuesta de recurso no encontrado
