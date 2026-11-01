Feature: Corregir identificación de historias, horas al completar y Planning Poker
  Como integrante del equipo
  Quiero reconocer cada HU, registrar el esfuerzo requerido y votar desde una computadora compartida
  Para completar historias correctamente y terminar las rondas sin perder información

  # AC-56-01
  Scenario: Distinguir historias homónimas por su UUID en Backlog y Sprint
    Given que el Backlog contiene dos historias con el mismo título y distintos UUID
    When consulto el listado del Backlog y las historias planificadas en un Sprint
    Then cada fila muestra junto al título el UUID completo correspondiente
    And al editar o trasladar una historia se conserva el mismo UUID

  # AC-56-02, AC-56-03
  Scenario: Completar una historia que ya tiene esfuerzo activo
    Given que una historia asignada al Sprint abierto tiene una entrada de esfuerzo activa
    And confirmo que verifiqué todos sus criterios de aceptación
    When la registro como completada sin enviar otra entrada de horas
    Then la historia queda completada en el Backlog y en el Sprint
    And el esfuerzo previo se conserva sin cambios ni duplicados

  # AC-56-02, AC-56-03
  Scenario: Pedir horas al completar una historia sin esfuerzo activo
    Given que una historia asignada al Sprint abierto no tiene entradas activas de esfuerzo
    When confirmo que verifiqué sus criterios y solicito completarla
    Then el sistema mantiene la historia incompleta y muestra el formulario de horas
    When registro un miembro del proyecto, fecha, actividad y horas válidas desde ese formulario
    Then se crea una única entrada asociada a la historia y al Sprint actual
    And la historia queda completada con sus criterios verificados

  # AC-56-04
  Scenario: Conservar el formulario y no completar ante un error al registrar horas
    Given que una historia sin esfuerzo activo está en el flujo de finalización
    When la fecha, el miembro, la actividad o las horas no son válidos o falla la escritura
    Then la historia permanece incompleta
    And no se guarda una entrada parcial
    And el formulario conserva todos los valores ingresados para corregir o reintentar

  # AC-56-05
  Scenario: Reintentar una finalización cuya respuesta se perdió
    Given que la petición de finalización guardó horas y completó la historia
    But la interfaz no recibió la respuesta
    When se repite la misma petición con los mismos datos de esfuerzo
    Then la historia continúa completada
    And sigue existiendo exactamente una entrada con esos datos

  # AC-56-05
  Scenario: No reemplazar ni duplicar una entrada activa distinta
    Given que la historia recibió esfuerzo activo desde otra operación mientras el formulario estaba abierto
    When se intenta finalizar con otros datos de horas
    Then la operación no inserta ni modifica ninguna entrada
    And mantiene la historia en su estado vigente y comunica el conflicto

  # AC-56-06
  Scenario: Votar por turnos desde una sola computadora
    Given que una ronda abierta incluye a Ana y Bea
    And Ana ya guardó su voto, que sigue oculto para los demás
    When selecciono a Bea en la misma computadora
    Then el selector conserva a Bea después de cargar su contexto
    And el formulario consulta y guarda solamente el voto de Bea
    And el voto de Ana permanece oculto y sin cambios

  # AC-56-06, AC-56-07
  Scenario: Revelar cuando votaron todos los participantes
    Given que todos los integrantes congelados en la ronda guardaron un voto
    When un participante solicita revelar los votos
    Then se muestran todos los votos y se informa si existen diferencias

  # AC-56-07
  Scenario: Informar quiénes faltan sin revelar estimaciones
    Given que una ronda abierta tiene participantes que todavía no votaron
    When se solicita revelar los votos
    Then el sistema rechaza la revelación e identifica los participantes faltantes
    And no devuelve ni muestra sus estimaciones ni las de los demás

  # AC-56-08
  Scenario: Conservar resultados y abrir una ronda posterior
    Given que una ronda finalizada reveló estimaciones diferentes
    When se inicia la ronda siguiente, se registran los votos y se alcanza un acuerdo unánime
    Then el historial conserva el resultado revelado de la ronda anterior
    And se registra la estimación acordada en Story Points
