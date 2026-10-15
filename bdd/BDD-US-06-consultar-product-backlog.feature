# US-06 — Consultar Product Backlog
# Fuente: Issue #6, docs/enunciado.md §2.2 y SDD-US-06.

Feature: Consultar el Product Backlog de un proyecto
  Como integrante del equipo
  Quiero consultar los elementos del Product Backlog
  Para conocer su prioridad, estado y datos antes de planificar

  Scenario: Consultar los elementos y sus campos vigentes
    Given que existe un proyecto con User Stories en su Product Backlog
    When se consulta el Backlog de ese proyecto
    Then se devuelven sus identificadores, títulos, descripciones, prioridades, estados, Story Points y criterios de aceptación
    And cada valor coincide con el elemento vigente

  Scenario: Consultar un Backlog vacío
    Given que existe un proyecto sin User Stories
    When se consulta su Product Backlog
    Then se devuelve una lista vacía

  Scenario: No mezclar elementos de otros proyectos
    Given que dos proyectos tienen User Stories distintas
    When se consulta el Product Backlog de uno de ellos
    Then no se devuelven elementos del otro proyecto

  Scenario: Diferenciar Story Points provisionales y definitivos
    Given que una User Story tiene Story Points provisionales y no tiene estimación definitiva
    When se consulta el Product Backlog
    Then la respuesta mantiene ambos campos separados
    And la interfaz presenta el valor provisional como provisional

  Scenario: Informar que el proyecto no existe
    Given que no existe el proyecto indicado
    When se consulta su Product Backlog
    Then se informa que el proyecto no fue encontrado

  Scenario: No informar una lista completa si falla el almacenamiento
    Given que el almacenamiento falla durante la consulta
    When se consulta el Product Backlog
    Then se informa un error controlado
    And no se devuelve una lista parcial como exitosa
