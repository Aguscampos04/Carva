# US-10 / Issue #10 — Cerrar Sprint y trasladar historias incompletas
# Fuente: criterios de aceptación de Issue #10 y SDD-US-10.

Feature: Cerrar un Sprint y trasladar historias incompletas

  @AC-01
  Scenario: Impedir el cierre cuando hay historias incompletas y no hay destino
    Given que existe un Sprint abierto con una historia incompleta
    And no existe un Sprint abierto posterior en el mismo proyecto
    When se intenta cerrar el Sprint
    Then la operación se rechaza indicando que necesita un siguiente Sprint
    And el Sprint permanece abierto
    And la historia continúa planificada únicamente en el Sprint original

  @AC-02 @AC-03
  Scenario: Trasladar solo historias incompletas y preservar el Sprint cerrado
    Given que el Sprint actual contiene una historia incompleta y una completada
    And existe un Sprint abierto posterior en el mismo proyecto
    When se cierra el Sprint actual
    Then la historia incompleta aparece una sola vez en el Sprint posterior
    And la historia completada no aparece en el Sprint posterior
    And el Sprint cerrado conserva sus dos historias y sus resultados originales
    And el Sprint posterior queda abierto

  @AC-02
  Scenario: No duplicar una historia incompleta ya planificada en el destino
    Given que la misma historia incompleta ya aparece en el Sprint posterior
    When se cierra el Sprint que también la tenía planificada
    Then la historia aparece una sola vez en el Sprint posterior
    And ambos Sprints conservan la planificación que les corresponde

  @AC-04
  Scenario: Cerrar un Sprint sin historias incompletas y sin destino
    Given que un Sprint abierto no tiene historias incompletas
    And no existe otro Sprint abierto posterior en el proyecto
    When se cierra el Sprint
    Then el Sprint queda cerrado
    And no se requiere crear otro Sprint

  @AC-05
  Scenario: Consultar el Sprint cerrado y rechazar nuevas mutaciones
    Given que un Sprint fue cerrado correctamente
    When se consulta su planificación
    Then se conservan sus historias y resultados históricos
    And intentar asignar una historia nueva se rechaza
    And intentar registrar una finalización nueva se rechaza
    And intentar cerrarlo de nuevo se rechaza

  @AC-01 @AC-03
  Scenario: Revertir cierre y traslado ante un error de persistencia
    Given que el Sprint tiene una historia incompleta y existe un Sprint posterior
    And ocurre un error de persistencia durante la operación de cierre
    When se intenta cerrar el Sprint
    Then la operación falla
    And el Sprint original permanece abierto
    And no queda una asignación parcial en el destino
