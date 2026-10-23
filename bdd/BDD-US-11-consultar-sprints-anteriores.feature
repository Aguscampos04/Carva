# US-11 / Issue #11 / SDD-US-11 / enunciado §2.3
Feature: Consultar Sprints anteriores sin alterar el pasado

  @BDD-01 @AC-01
  Scenario: Consultar un Sprint cerrado
    Given un proyecto con un Sprint cerrado y snapshots persistidos
    When se listan sus Sprints cerrados y se selecciona uno
    Then se obtiene su identificación y Sprint Goal
    And se obtienen sus historias planificadas y resultados registrados

  @BDD-02 @AC-02
  Scenario: Conservar el historial después de trasladar una historia
    Given un Sprint cerrado con una historia incompleta trasladada al siguiente
    When se consulta el Sprint original
    Then la historia permanece en su planificación histórica con el resultado original

  @BDD-03 @AC-02
  Scenario: Una edición y finalización posterior no cambian el pasado
    Given una historia incompleta de un Sprint cerrado trasladada a otro Sprint
    When se modifica en el Backlog y se completa en el siguiente Sprint
    And se reabre la base de datos y se consulta el Sprint cerrado
    Then su snapshot conserva título, descripción, criterios, puntos y resultado originales

  @BDD-04 @AC-01
  Scenario: Un proyecto sin Sprints cerrados
    Given un proyecto válido sin Sprints cerrados
    When se consulta su historial
    Then se devuelve una lista vacía
    And la interfaz indica que no hay Sprints cerrados

  @BDD-05
  Scenario: Recursos inexistentes o fuera del proyecto
    Given un identificador de proyecto o Sprint inexistente o ajeno
    When se solicita la consulta
    Then se devuelve el error de recurso no encontrado según el contrato

  @BDD-06
  Scenario: Errores y respuestas tardías en la interfaz
    Given una consulta de historial en curso
    When falla o cambia el proyecto o Sprint seleccionado
    Then se presenta el error correspondiente o se descarta la respuesta anterior
    And no se muestran datos de otra selección
