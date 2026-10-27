# US-19 / Issue #19 / EP-06 / enunciado §2.6.
# Aprobado para implementar por @Aguscampos04 el 2026-10-27.
# @propuesta identifica el origen en P/PT; aprobadas en la SDD. Automatización por tests Go/Node, no por un runner Gherkin.
Feature: Actualizar el estado y los Sprints de un defecto
  Como integrante del equipo
  Quiero actualizar un defecto existente
  Para registrar su evolución y resolución sin perder su identidad

  @US19-01 @AC-01 @AC-03
  Scenario: Asociar posteriormente el Sprint de detección
    Given un defecto registrado para una historia de Carva sin Sprint asociado
    And existe un Sprint de Carva
    When se registra ese Sprint como detección del defecto
    Then la consulta devuelve el mismo defecto con ese Sprint de detección
    And conserva su identidad, historia, descripción y severidad

  @US19-02 @AC-02 @AC-03
  Scenario: Resolver en un Sprint distinto del de detección
    Given un defecto de Carva detectado en el Sprint 1
    And existe el Sprint 2 de Carva
    When se registra estado Resuelto y resolución en el Sprint 2
    Then la consulta conserva detección en el Sprint 1 y resolución en el Sprint 2
    And el estado registrado es Resuelto

  @US19-03 @AC-03 @propuesta @P-01 @PT-01
  Scenario Outline: Cambiar estado sin alterar asociaciones
    Given un defecto de Carva con ambos Sprints asociados
    When se actualiza únicamente su estado a "<estado>"
    Then la consulta devuelve ese estado
    And conserva ambos Sprints y los datos no editables
    Examples:
      | estado      |
      | Abierto     |
      | En progreso |
      | Resuelto    |
      | Cerrado     |

  @US19-04 @AC-03 @PT-03 @propuesta
  Scenario: Conservar la actualización al reabrir la base
    Given un defecto cuyo estado y Sprints fueron actualizados correctamente
    When se reinicia la conexión a la base y se consulta la historia
    Then se recuperan los valores actualizados del mismo defecto

  @US19-05 @BR-04
  Scenario: Rechazar Sprint de otro proyecto sin actualización parcial
    Given un defecto de Carva con estado Abierto
    And un Sprint que pertenece a otro proyecto
    When se intenta cambiar estado y asociar ese Sprint a la vez
    Then la operación es rechazada
    And el defecto conserva su estado y asociaciones anteriores

  @US19-06 @BR-01 @PT-01 @propuesta
  Scenario: No modificar un defecto a través de otra historia o proyecto
    Given un defecto que pertenece a una historia de Carva
    When se intenta actualizarlo desde otra historia o proyecto
    Then se informa recurso no encontrado
    And el defecto original no cambia

  @US19-07 @PT-01 @propuesta
  Scenario Outline: Rechazar cambios inválidos sin guardar
    Given un defecto existente
    When se intenta guardar "<cambio>"
    Then se informa una solicitud inválida
    And el defecto no cambia
    Examples:
      | cambio                                |
      | estado no aprobado                    |
      | estado nulo                           |
      | Sprint cero, negativo o fraccionario   |
      | cuerpo sin campos editables           |
      | JSON con propiedades desconocidas     |
      | JSON con claves duplicadas            |
      | cambio de identidad o historia        |

  @US19-08 @P-02 @propuesta
  Scenario: Quitar explícitamente una asociación para corregirla
    Given un defecto con detección y resolución asociadas
    When se quita explícitamente la asociación de resolución
    Then resolución queda Sin Sprint
    And detección y estado conservan sus valores

  @US19-09 @P-01 @propuesta
  Scenario: Registrar resolución sin deducir un cambio de estado
    Given un defecto con estado En progreso y sin resolución asociada
    When se informa únicamente el Sprint de resolución
    Then se conserva estado En progreso y se guarda la asociación

  @US19-10 @P-03 @P-04 @propuesta
  Scenario: Registrar datos de resolución de un proyecto finalizado
    Given un defecto en una historia completada de un proyecto finalizado
    And un Sprint cerrado del mismo proyecto
    When se asocia ese Sprint como resolución del defecto
    Then se guarda la asociación sin reabrir proyecto, historia ni Sprint

  @US19-11 @PT-03 @propuesta
  Scenario: Rollback ante una falla de persistencia
    Given un defecto existente y un cambio válido de estado y asociaciones
    And el almacenamiento falla al guardar
    When se intenta actualizar el defecto
    Then se informa un error sin detalles internos
    And no quedan cambios parciales ni otro defecto creado

  @US19-12 @PT-03 @propuesta
  Scenario: Conservar cambios concurrentes de campos distintos
    Given dos consultas del mismo defecto
    When un integrante cambia estado y otro cambia detección sin editar estado
    Then ambas actualizaciones válidas quedan almacenadas
    And no se borra el campo omitido por cada solicitud

  @US19-13 @PT-04 @propuesta
  Scenario: Cancelar la edición en la interfaz
    Given un defecto mostrado con sus valores guardados
    When se edita un borrador y se cancela
    Then no se envía una actualización
    And se siguen mostrando los valores guardados

  @US19-14 @PT-04 @propuesta
  Scenario: Conservar borrador ante error y bloquear doble envío
    Given un defecto en edición con un guardado en curso
    When se vuelve a pulsar Guardar y luego falla la solicitud
    Then solo se envía una actualización
    And se muestra el error y se conserva el borrador para corregirlo

  @US19-15 @PT-04 @propuesta
  Scenario: Descartar una respuesta de un proyecto anterior
    Given una actualización de defecto en curso
    When se cambia el proyecto antes de recibir su respuesta
    Then esa respuesta no modifica los datos ni formularios del proyecto nuevo
