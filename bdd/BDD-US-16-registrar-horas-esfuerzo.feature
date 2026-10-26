# Estado: documentación en revisión; 31 decisiones D-01...D-31 aprobadas, incluida D-26 corregida.
# PT son propuestas técnicas, no aprobaciones de implementación.
# Enunciado §2.5/§2.7 -> EP-05 -> US-16 / Issue #16 -> SDD-US-16 -> AC -> BDD.
# Automatización, TEST-US-16-NN y evidencia TDD pendientes; estos escenarios no se ejecutaron.
# Convención existente: keywords Gherkin en inglés y comportamiento en español.
@US-16 @Issue-16 @documentacion-en-revision
Feature: Registrar horas estimadas y esfuerzo real
  Como integrante del equipo
  Quiero mantener estimaciones y registros de esfuerzo independientes por historia
  Para consultar las horas estimadas y el esfuerzo efectivamente registrado

  Background:
    Given que el instante del servidor está fijado en "2026-10-26T15:00:00Z"
    And el día de referencia se calcula en "America/Argentina/Buenos_Aires"
    And existe el proyecto "Carva Demo" en estado "En progreso", con inicio "2026-10-01" y fin ausente
    And existen las historias "A" y "B" del Backlog de "Carva Demo", con estimación de 8 horas cada una y sin Sprint
    And están registrados "Ana" con UUID "b1d168da-325a-40d1-90b7-b5a1ba006389" y "Bruno" con UUID "c2d168da-325a-40d1-90b7-b5a1ba006389" en ese proyecto
    And no hay registros de esfuerzo salvo los indicados en cada escenario

  @BDD-US-16-01 @AC-01 @D-01 @D-06 @D-07
  Scenario Outline: Crear historia con horas estimadas sin requerir Sprint ni puntos
    When se crea una historia con todos los campos requeridos de US-05 y <horas> horas estimadas, sin Story Points ni Sprint
    Then se guarda con UUID nuevo, estado "Pendiente" y <horas> horas estimadas
    And su total real es 0 y no se crea una asignación a Sprint
    And sus Story Points provisionales y definitivos permanecen ausentes

    Examples:
      | horas |
      | 0.01  |
      | 1     |
      | 1.50  |
      | 40.25 |

  @BDD-US-16-02 @AC-01 @AC-17 @D-02 @PT-04
  Scenario: Modificar estimación con esfuerzo y puntos definitivos
    Given que "A" tiene 5 Story Points definitivos y un registro activo de 3.25 horas
    When se edita "A" con campos válidos y estimación de 2.50 horas conservando su valor provisional
    Then "A" muestra estimación vigente de 2.50 y total real de 3.25 horas
    And conserva sus 5 Story Points definitivos, estado y registro de esfuerzo

  @BDD-US-16-03 @AC-01 @D-01 @D-02
  Scenario: Cambiar estimación sin Sprint y sin esfuerzo
    When se guarda una edición válida de "A" con estimación de 12.75 horas
    Then se consulta "A" con 12.75 horas estimadas y total real 0
    And no se asigna a Sprint ni se generan registros reales

  @BDD-US-16-04 @AC-02 @D-06 @D-07 @PT-01
  Scenario Outline: Rechazar estimación inválida sin guardar otros campos
    Given que se preparó una solicitud de "<operacion>" con los demás campos de US-05 válidos
    When se envía estimación <estimacion> junto al título "No debe guardarse"
    Then la operación se rechaza por estimación inválida sin redondearla
    And no se crea una historia ni se modifica ningún dato de las existentes

    Examples:
      | operacion   | estimacion |
      | alta        | ausente    |
      | alta        | null       |
      | alta        | 0          |
      | alta        | -0.01      |
      | alta        | 1.001      |
      | editar A    | ausente    |
      | editar A    | null       |
      | editar A    | 0          |
      | editar A    | 1.001      |
      | editar A    | "1.50"     |
      | editar A    | 1.500      |
      | editar A    | 1e1        |

  @BDD-US-16-05 @AC-03 @AC-12 @D-08 @D-28
  Scenario: Consultar historia histórica sin estimación y registrar esfuerzo
    Given que "A" proviene de datos históricos y no tiene estimación
    When se consulta su detalle y se registra para Ana 1.25 horas el "2026-10-25" con actividad "Revisar criterios", sin Sprint
    Then se conserva la identidad de "A" y su estimación ausente
    And el detalle muestra "Pendiente de regularización", total real 1.25 y el registro activo
    And consultar y registrar esfuerzo no inventan una estimación

  @BDD-US-16-06 @AC-03 @AC-16 @D-08 @D-21
  Scenario: Exigir regularización al editar una historia histórica
    Given que "A" no tiene estimación y está asignada a un Sprint abierto y uno cerrado con horas estimadas ausentes
    When se intenta cambiar solo su título con los demás campos editables válidos pero sin estimación
    Then se rechaza la edición completa y conserva su título y estimación ausente
    When se guarda esa edición con estimación de 4.25 horas
    Then "A" conserva su UUID y muestra el título nuevo y 4.25 horas estimadas
    And ambos snapshots muestran 4.25 horas estimadas y conservan todos los demás datos históricos

  @BDD-US-16-07 @AC-04 @AC-12 @AC-14 @D-03 @D-04 @D-19 @D-20 @D-22
  Scenario: Registrar esfuerzo válido desde el detalle de historia sin Sprint
    Given que el usuario abrió el detalle de "A" desde Product Backlog
    When selecciona el UUID de Ana de la lista y guarda fecha "2026-10-25", actividad "Implementar validación" y 2.50 horas sin Sprint
    Then se crea una entrada con UUID nuevo asociada a Ana, "A" y "Carva Demo"
    And conserva fecha, actividad y 2.50 horas con Sprint nulo
    And el detalle muestra esa entrada activa, 8 horas estimadas y total real 2.50

  @BDD-US-16-08 @AC-04 @AC-05 @D-03 @D-05 @D-10 @PT-02
  Scenario Outline: Rechazar esfuerzo con campo obligatorio ausente o inválido
    Given que se preparó un alta válida de 1 hora de Ana para "A" el "2026-10-25" con actividad "Revisión"
    When se reemplaza "<campo>" por <valor> en la solicitud
    Then se rechaza la solicitud sin registrar esfuerzo ni cambiar el total de "A"

    Examples:
      | campo       | valor   |
      | member_id   | ausente |
      | member_id   | null    |
      | member_id   | ""      |
      | work_date   | ausente |
      | work_date   | null    |
      | activity    | ausente |
      | activity    | null    |
      | activity    | ""      |
      | activity    | "   "   |
      | hours       | ausente |
      | hours       | null    |
      | hours       | 0       |
      | hours       | -0.01   |
      | hours       | 0.001   |
      | hours       | "1.25"  |

  @BDD-US-16-09 @AC-05 @D-09 @PT-01
  Scenario Outline: Rechazar fechas de trabajo mal formadas o imposibles
    When se intenta registrar esfuerzo para "A" con fecha <fecha> y los otros campos válidos
    Then se rechaza por fecha inválida sin crear registros

    Examples:
      | fecha                  |
      | "2026-02-29"           |
      | "2026-04-31"           |
      | "26/10/2026"           |
      | "2026-1-2"             |
      | "2026-10-26T00:00:00Z" |
      | "0000-01-01"           |
      | ""                     |

  @BDD-US-16-10 @AC-05 @PT-02
  Scenario Outline: Rechazar contratos inválidos de esfuerzo
    When se intenta administrar esfuerzo con <solicitud>
    Then se informa HTTP 400 con error JSON controlado
    And no cambian registros, estimaciones ni totales

    Examples:
      | solicitud                                                        |
      | alta con backlog_item_id en el cuerpo además de la historia de ruta |
      | alta con raíz JSON null                                          |
      | alta con dos objetos JSON consecutivos                           |
      | alta con dos claves hours en el mismo objeto                     |
      | alta con id o deleted_at elegidos por el cliente                  |
      | edición con sprint_id de tipo string en lugar de entero o null  |
      | baja con confirmed ausente                                       |
      | baja con confirmed false                                         |
      | baja con confirmed null                                          |
      | baja con confirmed "true"                                        |

  @BDD-US-16-11 @AC-06 @D-03
  Scenario: Distinguir miembros homónimos por UUID
    Given que hay otra Ana en "Carva Demo" con UUID "d3d168da-325a-40d1-90b7-b5a1ba006389"
    When se selecciona la segunda Ana por su UUID y se registra 1 hora en "A"
    Then la entrada pertenece exclusivamente al UUID "d3d168da-325a-40d1-90b7-b5a1ba006389"
    And el selector muestra nombre y UUID completo de ambas para distinguirlas
    And no se atribuye esfuerzo a la primera Ana por coincidencia de nombre

  @BDD-US-16-12 @AC-06 @AC-14 @D-03 @D-22 @D-27
  Scenario Outline: Validar existencia y pertenencia de todos los recursos
    Given que "A" está asignada al Sprint 1 de "Carva Demo"
    And existe "Proyecto vecino" con sus propios miembros, historias y Sprint 2
    When se intenta registrar esfuerzo para "A" con <recurso>
    Then se informa HTTP <codigo> sin cambios en ningún proyecto

    Examples:
      | recurso                                        | codigo |
      | UUID de miembro inexistente                    | 404    |
      | nombre Ana en lugar de UUID                    | 404    |
      | UUID de miembro de Proyecto vecino             | 404    |
      | historia de ruta inexistente                   | 404    |
      | historia de Proyecto vecino en la ruta de Demo | 404    |
      | Sprint inexistente                             | 404    |
      | Sprint 2 de Proyecto vecino                    | 404    |
      | Sprint de Demo al que A no está asignada        | 409    |

  @BDD-US-16-13 @AC-07 @AC-12 @D-09 @D-29 @PT-01
  Scenario Outline: Completar por defecto el día de la zona aprobada
    Given que el instante del servidor está fijado en "<instante>"
    And el navegador informa una fecha distinta
    When se abre el formulario de esfuerzo de "A" desde su detalle
    Then la fecha predeterminada es "<dia>" y puede editarse
    And no se guarda esfuerzo por abrir el formulario

    Examples:
      | instante             | dia        |
      | 2026-10-26T02:30:00Z | 2026-10-25 |
      | 2026-10-26T03:00:00Z | 2026-10-26 |

  @BDD-US-16-14 @AC-07 @D-09 @D-18
  Scenario Outline: Admitir límites de fechas y rechazar días fuera del rango
    When se intenta registrar 0.01 horas de Ana en "A" con actividad "Revisión" el "<fecha>"
    Then el resultado es "<resultado>"
    And el total real de "A" es <total>

    Examples:
      | fecha      | resultado                  | total |
      | 2026-10-01 | registro aceptado          | 0.01  |
      | 2026-10-25 | registro aceptado          | 0.01  |
      | 2026-10-26 | registro aceptado          | 0.01  |
      | 2026-09-30 | rechazo por anterior a inicio | 0  |
      | 2026-10-27 | rechazo por fecha futura   | 0     |

  @BDD-US-16-15 @AC-07 @D-18 @D-31 @PT-04
  Scenario: Proteger esfuerzo existente al cambiar inicio y permitir quitar el límite
    Given que Ana tiene 1 hora activa el "2026-10-05" en "A"
    When se intenta guardar nombre "No guardar" e inicio "2026-10-06" en el proyecto
    Then se rechaza todo el cambio y conserva nombre "Carva Demo" e inicio "2026-10-01"
    And el esfuerzo conserva su fecha "2026-10-05"
    When se cambia solo el inicio a "2026-10-05"
    Then se acepta la igualdad sin modificar registros
    When se elimina explícitamente la fecha de inicio y se registra otra hora el "2026-09-30"
    Then se admite el esfuerzo pasado sin inferir otra fecha de inicio

  @BDD-US-16-16 @AC-07 @D-18 @D-31
  Scenario: Los eliminados no bloquean el cambio de inicio y conservan su historia
    Given que existe un esfuerzo de 2 horas del "2026-10-05" eliminado lógicamente
    And existe otro esfuerzo activo de 1 hora del "2026-10-10"
    When se guarda el inicio del proyecto "2026-10-06"
    Then el cambio se acepta y el inicio consultado es "2026-10-06"
    And el eliminado conserva identificador, miembro, historia, Sprint, fecha, actividad, horas y marca de baja anteriores
    And solo el esfuerzo activo participa del total de 1 hora y del control de inicio

  @BDD-US-16-17 @AC-08 @AC-13 @D-14 @D-15
  Scenario: Administrar esfuerzo de historia Completada en proyecto Finalizado
    Given que "A" está "Completada" y el proyecto está "Finalizado" con fin "2026-10-20"
    When se registra 1 hora de Ana en "A" el "2026-10-20" con actividad "Revisión posterior"
    And se edita ese registro a 1.25 horas y actividad "Revisión corregida" manteniendo fecha y Sprint nulo
    And se confirma su eliminación lógica
    Then el registro mantiene su UUID, última actividad y horas y queda eliminado
    And "A" conserva estado "Completada" y estimación 8, con total real 0
    And el proyecto conserva estado "Finalizado", inicio "2026-10-01" y fin "2026-10-20"

  @BDD-US-16-18 @AC-08 @D-15 @D-17
  Scenario Outline: Validar fin real conocido aunque la carga sea posterior
    Given que el proyecto está "Finalizado" con fin "2026-10-20"
    When se intenta registrar esfuerzo de Ana para "A" el "<fecha>" con los otros campos válidos
    Then el resultado es "<resultado>"
    And el fin sigue siendo "2026-10-20" y el estado sigue "Finalizado"

    Examples:
      | fecha      | resultado                     |
      | 2026-10-19 | registro aceptado             |
      | 2026-10-20 | registro aceptado             |
      | 2026-10-21 | rechazo por posterior al fin  |

  @BDD-US-16-19 @AC-09 @AC-17 @D-16 @D-25
  Scenario: Acumular múltiples registros independientes y superar estimación
    Given que "A" tiene 2 horas estimadas y 5 Story Points definitivos
    When Ana registra para "A" el "2026-10-25" dos esfuerzos de 1.25 y 2.50 horas con actividad válida
    Then se conservan dos entradas con UUID distintos y total real 3.75
    And "A" mantiene 2 horas estimadas y 5 Story Points definitivos

  @BDD-US-16-20 @AC-09 @D-24
  Scenario: Sumar por proyecto atravesando historias y Sprints
    Given que Ana tiene 10 horas en "A" asociadas a un Sprint cerrado y 13.99 en "B" asociadas a uno abierto el "2026-10-25"
    And ambas historias pertenecen a esos Sprints de "Carva Demo"
    When Ana registra 0.01 horas sin Sprint en "A" el "2026-10-25"
    Then se acepta el registro y el acumulado de Ana en el proyecto para ese día es exactamente 24
    When intenta registrar otras 0.01 horas en "B" ese día
    Then se rechaza por superar 24 horas y se conserva el acumulado 24

  @BDD-US-16-21 @AC-09 @D-24
  Scenario Outline: Separar cupo por miembro, día y proyecto
    Given que Ana tiene 24 horas activas el "2026-10-25" en "Carva Demo"
    And Ana también está registrada independientemente en "Proyecto vecino" con UUID propio e historia con estimación válida
    When se registra 1 hora con los demás campos válidos para "<miembro>" el "<dia>" en "<proyecto>"
    Then se acepta sin alterar las 24 horas previas de Ana del "2026-10-25" en "Carva Demo"

    Examples:
      | miembro                    | dia        | proyecto        |
      | Bruno                      | 2026-10-25 | Carva Demo      |
      | Ana                        | 2026-10-26 | Carva Demo      |
      | Ana con UUID del vecino    | 2026-10-25 | Proyecto vecino |

  @BDD-US-16-22 @AC-09 @AC-10 @D-11 @D-24
  Scenario: Editar horas sin contar dos veces el registro actual
    Given que Ana tiene el mismo día en "A" un registro de 10 horas y otro de 13.50
    When se edita el registro de 10 a 10.50 horas con los demás datos válidos
    Then se conserva su UUID y el total diario queda exactamente en 24
    When se intenta cambiarlo a 10.51 horas
    Then se rechaza la edición y conserva 10.50 horas con total diario 24

  @BDD-US-16-23 @AC-09 @AC-10 @D-11 @D-24
  Scenario: Cambiar miembro y fecha revalida cupo destino y conserva origen si falla
    Given que Ana tiene un registro de 2 horas el "2026-10-24" en "A"
    And Bruno tiene 23 horas el "2026-10-25" en "B"
    When se intenta editar el registro de Ana a Bruno el "2026-10-25" conservando 2 horas
    Then se rechaza sin modificar miembro, fecha ni totales anteriores
    When se edita ese mismo UUID a Bruno el "2026-10-25" con 1 hora y actividad "Corrección", historia "B" y Sprint nulo explícito
    Then se conserva el UUID con todos los campos corregidos y total real de "A" 0
    And Bruno acumula 24 horas el "2026-10-25" y Ana 0 el "2026-10-24"

  @BDD-US-16-24 @AC-11 @AC-12 @D-12 @D-13 @PT-02
  Scenario: Confirmar baja lógica libera cupo sin borrar trazabilidad
    Given que Ana tiene dos entradas activas de 12 horas cada una el "2026-10-25" en "A"
    When desde el detalle se elimina una entrada y el usuario confirma el diálogo que la identifica
    Then se conserva almacenada con sus datos e identificador y marca de baja del servidor
    And las consultas normales no la muestran y el total de "A" pasa a 12
    When Ana registra otra entrada de 12 horas ese día y se repite confirmada la baja del registro eliminado
    Then el total vuelve a 24 sin duplicar el efecto de la baja ni modificar su primera marca

  @BDD-US-16-25 @AC-11 @D-12
  Scenario Outline: Cancelar eliminación no envía escritura ni guarda borradores
    Given que existe una entrada activa de 2 horas en "A" con un borrador de edición sin guardar
    When se solicita eliminarla y se elige "<accion>" en la confirmación
    Then no se envía solicitud de baja ni de guardado
    And la entrada y el total activo permanecen en 2 horas

    Examples:
      | accion             |
      | Cancelar           |
      | cerrar el diálogo  |

  @BDD-US-16-26 @AC-10 @AC-14 @D-11 @D-26 @PT-02
  Scenario Outline: Conservar Sprint compatible y resolver explícitamente el incompatible al cambiar historia
    Given que "A" está asignada al Sprint 1 y "B" a <pertenencias> de "Carva Demo"
    And existe el Sprint 3 del proyecto sin asignación de "B"
    And existe una entrada de Ana de 2 horas en "A" asociada al Sprint 1
    When se cambia la historia de la entrada a "B" con <eleccion> y todos los demás campos válidos
    Then el resultado es "<resultado>"
    And la entrada queda en historia "<historia>" con Sprint <sprint> conservando su UUID
    And los totales reales de "A" y "B" son respectivamente <totalA> y <totalB>

    Examples:
      | pertenencias    | eleccion                             | resultado        | historia | sprint | totalA | totalB |
      | Sprints 1 y 2   | sprint_id omitido                    | edición aceptada | B        | 1      | 0      | 2      |
      | Sprints 1 y 2   | selección explícita del Sprint 1     | edición aceptada | B        | 1      | 0      | 2      |
      | Sprints 1 y 2   | selección explícita del Sprint 2     | edición aceptada | B        | 2      | 0      | 2      |
      | Sprints 1 y 2   | selección explícita Sin Sprint       | edición aceptada | B        | nulo   | 0      | 2      |
      | Sprint 2        | sprint_id omitido                    | rechazo          | A        | 1      | 2      | 0      |
      | Sprint 2        | selección explícita del Sprint 1     | rechazo          | A        | 1      | 2      | 0      |
      | Sprint 2        | selección explícita del Sprint 2     | edición aceptada | B        | 2      | 0      | 2      |
      | Sprint 2        | selección explícita Sin Sprint       | edición aceptada | B        | nulo   | 0      | 2      |
      | Sprint 2        | Sprint 3 existente sin pertenencia B | rechazo          | A        | 1      | 2      | 0      |

  @BDD-US-16-27 @AC-06 @AC-10 @D-27
  Scenario: Aislar registros y corregir proyecto mediante baja y alta independientes
    Given que existe una entrada de 2 horas en "A" y "Proyecto vecino" tiene su historia y miembro propios
    When se intenta reasignar esa entrada a la historia del vecino desde la ruta de "Carva Demo"
    Then se rechaza sin modificar registro, historia ni proyecto originales
    When se consulta o intenta editar o eliminar ese UUID desde la ruta de "Proyecto vecino"
    Then se informa recurso no encontrado sin cambios
    When se confirma la baja desde "Carva Demo" y se crea esfuerzo válido en el vecino con su propio miembro
    Then el registro original sigue eliminado en "Carva Demo" y la nueva entrada tiene otro UUID en el vecino

  @BDD-US-16-28 @AC-14 @AC-15 @D-22 @D-23
  Scenario Outline: Administrar esfuerzo con asociación opcional incluso a Sprint cerrado
    Given que "A" pertenece al Sprint 1 cerrado y al Sprint 2 abierto
    When se registra 1 hora en "A" con <asociacion>, se edita a 1.25 horas y se confirma su baja lógica
    Then las tres operaciones son aceptadas y el total activo final es 0
    And los Sprints conservan estado, historias planificadas y resultados

    Examples:
      | asociacion                  |
      | sprint_id omitido en el alta |
      | Sprint nulo explícito       |
      | Sprint 1 cerrado            |
      | Sprint 2 abierto            |

  @BDD-US-16-29 @AC-15 @AC-16 @D-21 @D-22 @D-23 @PT-04
  Scenario: Trasladar historia no mueve esfuerzo ni duplica totales y estimar propaga solo horas
    Given que "A" está incompleta en el Sprint 1 y tiene una entrada de 2 horas asociada a él
    And existe el Sprint 2 abierto posterior y no contiene "A"
    When se cierra el Sprint 1 y se traslada "A" al Sprint 2
    And se modifica el título de "A" y su estimación a 6.25 horas mediante edición válida
    Then la entrada conserva UUID y asociación al Sprint 1 y el total real de "A" es 2
    And ambos snapshots muestran estimación vigente 6.25 y conservan título, criterios, prioridad, puntos y resultados originales
    And el Sprint 1 sigue cerrado y el Sprint 2 conserva su única asignación de "A"

  @BDD-US-16-30 @AC-12 @AC-18 @PT-05
  Scenario: Migrar y reabrir conserva legado y persistencia nueva
    Given una base temporal del esquema anterior con proyecto, miembros, historias sin horas y snapshots abiertos y cerrados
    When se aplica la futura migración de US-16 y se reabre la aplicación con esa base
    Then se conservan identificadores, textos, fechas, estados, puntos, asignaciones y resultados
    And las horas estimadas ausentes siguen nulas con indicador pendiente y los totales reales son 0
    When se registra esfuerzo en una historia histórica, se regulariza su estimación y se elimina lógicamente esa entrada
    And se reinicia usando la misma base y se vuelve a aplicar el mecanismo de migraciones
    Then se conserva la estimación regularizada en historia y snapshots y el registro eliminado sigue almacenado
    And el detalle no muestra esa entrada y el total sigue 0 sin duplicaciones ni reinicializaciones

  @BDD-US-16-31 @AC-11 @AC-18 @D-13 @PT-04
  Scenario Outline: Fallos previos al commit no dejan esfuerzo parcial
    Given que existe una entrada activa de 2 horas en "A" y se inyectará un fallo antes de confirmar la escritura
    When se intenta "<operacion>" con datos válidos
    Then se informa error controlado sin éxito aparente ni detalles SQL
    And al reabrir se conserva exactamente el registro activo original de 2 horas, su UUID y el total de "A"
    And no quedan registros adicionales, reasignaciones parciales ni marcas de baja nuevas

    Examples:
      | operacion                                  |
      | registrar otra entrada                     |
      | editar historia, miembro, fecha y horas     |
      | eliminar lógicamente con confirmación       |

  @BDD-US-16-32 @AC-09 @AC-18 @D-24 @PT-04
  Scenario: Dos altas concurrentes no superan el cupo diario
    Given que Ana tiene 23 horas activas el "2026-10-25" en "A"
    And dos clientes con conexiones independientes prepararon registrar 1 hora adicional en historias diferentes de ese proyecto y día
    When ambas altas se procesan concurrentemente sin agotar la espera de escritura
    Then exactamente una se acepta y la otra se rechaza por superar el cupo
    And el acumulado persistido de Ana para ese proyecto y día es exactamente 24

  @BDD-US-16-33 @AC-07 @AC-18 @D-18 @D-31 @PT-04
  Scenario Outline: Serializar cambio de inicio y alta de esfuerzo
    Given que un cliente prepara inicio "2026-10-10" y otro esfuerzo del "2026-10-05" para "A"
    When se procesan concurrentemente y confirma primero "<primera>"
    Then la segunda operación se rechaza por incompatibilidad temporal
    And el inicio persistido es "<inicio>" y hay <cantidad> registros de ese día
    And nunca se cambia automáticamente la fecha de trabajo

    Examples:
      | primera        | inicio     | cantidad |
      | cambiar inicio | 2026-10-10 | 0        |
      | cargar esfuerzo| 2026-10-01 | 1        |

  @BDD-US-16-34 @AC-15 @AC-16 @AC-18 @D-21 @PT-04
  Scenario Outline: Coordinar estimación con asignación y traslado concurrentes
    Given que se prepara modificar la estimación de "A" a 9.75 horas
    And se prepara "<operacion>" válida de esa misma historia
    When se procesan concurrentemente las operaciones por conexiones independientes y ambas confirman
    Then "A" y todos sus snapshots existentes y recién creados muestran 9.75 horas estimadas
    And solo se crean las asignaciones previstas y no cambia ningún otro dato de planificación por actualizar horas
    And el esfuerzo preexistente conserva sus registros y asociaciones sin duplicarse

    Examples:
      | operacion                                                       |
      | asignarla a un Sprint abierto que aún no la contiene             |
      | cerrar su Sprint origen y trasladarla a un destino abierto vacío |
      | cerrar su Sprint origen cuando el destino ya contiene la historia|

  @BDD-US-16-35 @AC-11 @AC-18 @D-11 @D-13 @PT-04
  Scenario Outline: Baja y edición concurrentes no resucitan registros
    Given que existe una entrada activa de 2 horas en "A"
    And se preparan editarla a 3 horas y eliminarla con confirmación
    When se procesan concurrentemente y confirma primero "<primera>"
    Then la entrada queda eliminada conservando su UUID y <horas> horas almacenadas
    And el resultado de la edición es "<resultado>" y el total activo es 0
    And las consultas y ediciones normales posteriores de esa entrada responden no encontrado

    Examples:
      | primera | horas | resultado      |
      | baja    | 2     | no encontrado  |
      | edición | 3     | edición válida |

  @BDD-US-16-36 @AC-16 @AC-18 @D-21 @PT-04
  Scenario: Propagación de estimación es atómica también en cerrado
    Given que "A" tiene snapshots en dos Sprints, uno de ellos cerrado, ambos con estimación 8
    And se inyectará un fallo al actualizar uno de los snapshots antes de commit
    When se intenta guardar estimación 7.25 y un nuevo título válido en "A"
    Then se informa error sin guardar título ni estimación parcial
    And al reabrir "A" y todos sus snapshots conservan estimación 8 y sus demás datos originales

  @BDD-US-16-37 @AC-17 @AC-18 @PT-04
  Scenario Outline: Edición de horas no sobrescribe puntos o compleción concurrentes
    Given que se preparó una edición válida de "A" para cambiar su estimación a 9.25 horas conservando sus campos de puntos provisionales
    And se preparó "<accion>" con todas sus precondiciones satisfechas
    When ambas operaciones se ejecutan concurrentemente y confirman
    Then "A" conserva estimación vigente 9.25 y "<resultado>"
    And sus registros de esfuerzo no cambian y los demás datos históricos de sus Sprints no se reescriben

    Examples:
      | accion                                                     | resultado                       |
      | registrar acuerdo unánime de 5 puntos en ronda finalizada    | 5 Story Points definitivos      |
      | completar A verificando sus criterios en un Sprint abierto  | estado Completada y verificación |

  @BDD-US-16-38 @AC-18 @PT-04
  Scenario: Agotar contención sin éxito falso ni reintento ciego
    Given que otra conexión mantiene la escritura ocupada más allá de cinco segundos
    When se solicita un registro válido de esfuerzo
    Then se informa indisponibilidad temporal sin registrar la entrada ni mostrar éxito
    And no se presenta un total nuevo ni una lista vacía como respuesta exitosa

  @BDD-US-16-39 @AC-12 @D-19 @D-20 @PT-06
  Scenario: Manejar errores y respuestas tardías sin mezclar historias
    Given que está abierta la administración de esfuerzo de "A" con un borrador sin guardar
    When falla una solicitud de guardado antes de persistir
    Then se muestra el error y se conserva el borrador sin agregar una entrada aparente
    When se cambia a "B" y llega una respuesta de consulta pendiente de "A"
    Then se descarta esa respuesta y no se muestran entradas ni totales de "A" en "B"

  @BDD-US-16-40 @AC-05 @AC-07 @AC-08 @AC-10 @D-09 @D-11 @D-17 @D-18
  Scenario Outline: Revalidar fecha al editar sin guardar otros cambios inválidos
    Given que el proyecto está "Finalizado" con fin "2026-10-20"
    And existe una entrada de Ana de 1 hora el "2026-10-15" en "A"
    When se intenta editar su fecha a "<fecha>" y su actividad a "No debe guardarse"
    Then se rechaza la edición por "<motivo>" y conserva fecha, actividad, horas e identificador anteriores

    Examples:
      | fecha      | motivo           |
      | 2026-09-30 | anterior a inicio|
      | 2026-10-21 | posterior a fin  |
      | 2026-10-27 | fecha futura     |
      | 2026-02-30 | fecha inválida   |

  @BDD-US-16-41 @AC-08 @AC-18 @D-15 @D-17 @PT-04
  Scenario Outline: Finalizar y cargar esfuerzo concurrentemente conserva ambas operaciones válidas
    Given que el proyecto aún no está finalizado
    And dos conexiones prepararon finalizarlo con confirmación y registrar 1 hora para "A" el "2026-10-26"
    When se procesan concurrentemente y confirma primero "<primera>"
    Then ambas operaciones se aceptan y el proyecto queda "Finalizado" con fin "2026-10-26"
    And existe una sola entrada activa de 1 hora con fecha igual al fin, sin alterar inicio ni Sprints

    Examples:
      | primera         |
      | finalizar       |
      | cargar esfuerzo |

  @BDD-US-16-42 @AC-09 @AC-10 @AC-18 @D-11 @D-24 @PT-04
  Scenario: Dos ediciones concurrentes hacia un mismo cupo revalidan el destino
    Given que Ana tiene 20 horas activas el "2026-10-25" en "Carva Demo"
    And tiene dos entradas de 1 hora cada una el "2026-10-23" y el "2026-10-24"
    And dos conexiones prepararon mover respectivamente esas entradas al "2026-10-25" con 3 horas cada una
    When se procesan concurrentemente las dos ediciones sin agotar la espera de escritura
    Then exactamente una edición se acepta y la otra se rechaza por superar el cupo
    And Ana acumula 23 horas el "2026-10-25" y el registro rechazado conserva su fecha original y 1 hora
    And ambas entradas conservan sus UUID sin movimientos parciales

  @BDD-US-16-43 @AC-19 @D-30 @PT-02 @PT-06
  Scenario Outline: Bloquear toda alta y edición de esfuerzo en Finalizado sin fin
    Given que "Carva Demo" está "Finalizado" con fecha de fin NULL conservada de datos históricos
    And existe una entrada activa de Ana de 1 hora del "2026-10-15" en "A" con actividad "Revisión"
    When se intenta por la API "<operacion>" con datos y relaciones válidos
    Then se responde HTTP 409 indicando fecha de fin desconocida
    And no se crean entradas ni se modifica ningún campo de la existente ni el total de 1 hora
    And el proyecto conserva estado "Finalizado", inicio "2026-10-01" y fin NULL
    And la interfaz mantiene alta y edición de esfuerzo deshabilitadas con un aviso explicativo

    Examples:
      | operacion                                                |
      | registrar esfuerzo para A                                 |
      | editar solo actividad conservando fecha y demás datos     |
      | cambiar miembro a Bruno conservando fecha y demás datos   |
      | cambiar fecha de trabajo a 2026-10-16                     |
      | cambiar horas a 2                                        |
      | reasignar entrada a B conservando fecha y Sprint nulo      |

  @BDD-US-16-44 @AC-19 @D-30
  Scenario: Permitir consultas y baja lógica con fin desconocido sin inventar fechas
    Given que "Carva Demo" está "Finalizado" con fin NULL
    And existe una entrada activa de Ana de 1 hora del "2026-10-15" en "A"
    When se consulta el detalle de "A", su esfuerzo y esa entrada individual
    Then se muestra la entrada activa y el total de 1 hora sin exigir fecha de fin
    And la interfaz permite consultar y solicitar la baja con confirmación
    When se confirma su eliminación lógica y se reinicia usando la misma base
    Then el registro permanece almacenado con sus datos y marca de baja y el total activo es 0
    And el proyecto conserva estado "Finalizado", inicio "2026-10-01" y fin NULL

  @BDD-US-16-45 @AC-07 @AC-18 @D-18 @D-31 @PT-04
  Scenario Outline: Cambio de inicio y baja concurrentes consideran solo activos de forma atómica
    Given que existe una única entrada activa del "2026-10-05" y no hay otras entradas anteriores al "2026-10-10"
    And dos conexiones preparan cambiar el inicio a "2026-10-10" y eliminar esa entrada con confirmación
    When se procesan concurrentemente y se serializa primero "<primera>"
    Then la baja lógica se acepta y el resultado del cambio de inicio es "<resultado>"
    And el inicio persistido es "<inicio>" y el total activo es 0
    And el eliminado conserva fecha "2026-10-05" y los demás datos históricos sin ajustes automáticos

    Examples:
      | primera        | resultado | inicio     |
      | baja lógica    | aceptado  | 2026-10-10 |
      | cambio de inicio | rechazado | 2026-10-01 |

  @BDD-US-16-46 @AC-10 @AC-14 @D-26 @PT-02 @PT-06
  Scenario: Cambiar historia conservando ausencia de Sprint sin forzar una selección
    Given que existe una entrada de Ana de 2 horas en "A" sin Sprint
    And "B" pertenece a un Sprint abierto del proyecto
    When se cambia la historia de la entrada a "B" omitiendo sprint_id y conservando los demás campos válidos
    Then se acepta con el mismo UUID y Sprint nulo, sin seleccionar el Sprint de "B" automáticamente
    And el total de "A" pasa a 0 y el de "B" a 2
