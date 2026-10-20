# Estado: aprobado para implementar por el responsable el 2026-10-20,
# con ajustes previos de BDD-07/18; automatización y resultados en evidencia-ia-us-02.md.
# Fuente: enunciado §2.1 -> EP-01 -> US-02 / Issue #2 -> D-01...D-10 de US-02
# -> sdd/SDD-US-02-modificar-proyecto.md -> AC -> escenarios de este archivo.
# Las etiquetas @PT identifican precisiones técnicas aprobadas para esta implementación.
# Fechas y reloj controlados: no depender de la fecha real de ejecución.
@US-02 @Issue-2 @aprobado-para-implementar
Feature: Modificar proyecto y registrar fechas
  Como integrante del equipo
  Quiero modificar los datos del proyecto y finalizarlo manualmente
  Para mantener sus datos y consultar su estado y fechas persistidos

  Background:
    Given que el servidor usa la zona propuesta "America/Argentina/Buenos_Aires"
    And el instante del servidor está fijado en "2026-10-20T15:00:00Z"
    And existe el proyecto "Carva Demo" con identificador "0828dc84-050f-421b-b842-420f63851a2b"
    And el proyecto está "Planificado", sin fecha de inicio ni de finalización y sin Sprints

  @AC-01 @D-02
  Scenario: BDD-US-02-01 Modificar y persistir los datos del proyecto correcto
    Given que existe otro proyecto llamado "Proyecto vecino"
    When se guarda el nombre "Carva II" y la fecha de inicio "2026-10-01" en "Carva Demo"
    And se reinicia la aplicación conservando la misma base
    And se consulta el proyecto por su identificador original
    Then se obtiene el nombre "Carva II" y la fecha de inicio "2026-10-01"
    And conserva su identificador y el estado "Planificado"
    And no tiene fecha de finalización
    And "Proyecto vecino" conserva todos sus datos

  @AC-01 @AC-02 @D-02 @PT-02
  Scenario: BDD-US-02-02 Modificar solo el nombre conserva el inicio registrado
    Given que el proyecto tiene fecha de inicio "2026-10-01"
    When se guarda únicamente el nombre "  Carva II  " omitiendo la fecha de inicio
    And se consulta de nuevo el proyecto
    Then el nombre es exactamente "  Carva II  "
    And la fecha de inicio sigue siendo "2026-10-01"
    And conserva su estado y su fecha de finalización ausente

  @AC-02 @D-02 @PT-02
  Scenario: BDD-US-02-03 Reemplazar y eliminar la fecha opcional sin cambiar el nombre
    Given que el proyecto tiene fecha de inicio "2026-10-01"
    When se guarda únicamente la fecha de inicio "2026-10-02"
    Then al consultar la fecha de inicio es "2026-10-02"
    When se guarda explícitamente una fecha de inicio nula
    Then al consultar no tiene fecha de inicio
    And mantiene el nombre "Carva Demo" y el estado "Planificado"

  @AC-03 @PT-02
  Scenario Outline: BDD-US-02-04 Rechazar nombres inválidos sin guardar la fecha válida
    When se intenta guardar el cuerpo de edición <cuerpo>
    Then se rechaza la edición con error de validación
    And al consultar el proyecto conserva exactamente todos sus datos anteriores

    Examples:
      | cuerpo                                        |
      | {"name":"","start_date":"2026-10-01"}        |
      | {"name":"   ","start_date":"2026-10-01"}     |
      | {"name":null,"start_date":"2026-10-01"}      |

  @AC-03 @D-06 @PT-01
  Scenario Outline: BDD-US-02-05 Rechazar fechas no canónicas o imposibles sin edición parcial
    When se intenta guardar el nombre "No debe guardarse" y la fecha de inicio <fecha>
    Then se rechaza la edición por fecha inválida
    And el proyecto mantiene el nombre "Carva Demo"
    And conserva sus fechas ausentes y su estado "Planificado"

    Examples:
      | fecha                  |
      | "2026-02-29"           |
      | "2026-04-31"           |
      | "2026-13-01"           |
      | "2026-00-10"           |
      | "2026-10-00"           |
      | "20/10/2026"           |
      | "2026-1-2"             |
      | "2026-10-20T00:00:00Z" |
      | " 2026-10-01 "         |
      | "0000-01-01"           |
      | "10000-01-01"          |
      | ""                     |

  @AC-03 @PT-01
  Scenario: BDD-US-02-06 Aceptar un día bisiesto válido
    When se guarda la fecha de inicio "2024-02-29"
    Then la consulta devuelve exactamente "2024-02-29" como inicio
    And el proyecto sigue "Planificado" sin fecha de finalización

  @AC-04 @D-01 @D-03 @D-05 @D-06 @PT-01
  Scenario Outline: BDD-US-02-07 Finalizar desde ambos estados incluso sin inicio ni Sprints
    Given que el proyecto está "<estado>" con <cantidad> Sprints existentes
    And su fecha de inicio es <inicio>
    When el usuario pulsa "Finalizar proyecto"
    And confirma la finalización solicitada
    And se reinicia la aplicación conservando la misma base
    And se consulta de nuevo el proyecto
    Then figura en estado "Finalizado" con fecha de finalización "2026-10-20"
    And conserva su nombre, identificador, inicio y los <cantidad> Sprints previos

    Examples:
      | estado      | cantidad | inicio       |
      | Planificado | 0        | ausente      |
      | Planificado | 0        | "2026-10-01" |
      | En progreso | 1        | ausente      |
      | En progreso | 1        | "2026-10-01" |

  @AC-05 @D-07 @PT-05
  Scenario: BDD-US-02-08 Solicitar confirmación y cancelar sin guardar ningún cambio
    Given que el usuario tiene un borrador sin guardar con nombre "Carva borrador"
    When pulsa "Finalizar proyecto"
    Then se solicita confirmación identificando el proyecto y aclarando que se usan los datos guardados
    And todavía no se envía una solicitud de escritura
    When el usuario cancela o cierra la confirmación
    Then no se envía una solicitud de escritura
    And el proyecto conserva nombre "Carva Demo", estado "Planificado" y ambas fechas ausentes
    And conserva sus Sprints, miembros y Product Backlog

  @AC-02 @AC-06 @D-02 @D-04 @PT-02
  Scenario: BDD-US-02-09 Editar nombre e inicio después de finalizar
    Given que el proyecto está "Finalizado" con inicio "2026-10-01" y fin "2026-10-20"
    When se guarda el nombre "Carva corregido" y el inicio "2026-10-02"
    Then al consultar se obtienen el nombre y el inicio corregidos
    And sigue "Finalizado" con fin "2026-10-20"
    When se elimina explícitamente la fecha de inicio
    Then al consultar el inicio es nulo
    And sigue "Finalizado" con nombre "Carva corregido" y fin "2026-10-20"

  @AC-03 @AC-06 @D-04 @D-06
  Scenario: BDD-US-02-10 Permitir inicio igual al fin y rechazar inicio posterior
    Given que el proyecto está "Finalizado" con inicio ausente y fin "2026-10-20"
    When se guarda el inicio "2026-10-20"
    Then la edición es aceptada y la consulta muestra ambas fechas iguales
    When se intenta guardar el nombre "No debe guardarse" y el inicio "2026-10-21"
    Then se rechaza la edición sin cambios parciales
    And conserva el nombre "Carva Demo", estado "Finalizado" y ambas fechas "2026-10-20"

  @AC-07 @D-01 @D-02
  Scenario: BDD-US-02-11 El paso del tiempo y la fecha de inicio no cambian el estado
    Given que el proyecto tiene inicio "2026-10-20"
    When el reloj avanza al día "2026-10-21" en la zona del servidor
    And se consulta el proyecto sin solicitar finalización
    Then conserva el estado "Planificado" y el inicio "2026-10-20"
    And no tiene fecha de finalización

  @AC-07 @D-02 @US-01 @US-07
  Scenario: BDD-US-02-12 Preservar alta y transiciones de creación de Sprints
    Given que se consulta el proyecto recién creado
    Then su estado es "Planificado" y sus fechas son nulas
    When se crea el primer Sprint con objetivo "Entregar MVP"
    Then su estado pasa a "En progreso" y sus fechas siguen nulas
    When se crea otro Sprint con objetivo "Mejorar interfaz"
    Then conserva el estado "En progreso" y ambas fechas nulas
    And ambos Sprints se pueden consultar con sus objetivos

  @AC-08 @D-03 @PT-01
  Scenario Outline: BDD-US-02-13 Registrar el día de la zona acordada con reloj del servidor
    Given que el instante del servidor está fijado en "<instante>"
    And el navegador y el host muestran un día distinto al del proyecto
    When se solicita y confirma la finalización del proyecto
    Then la consulta muestra fecha de finalización "<dia>"

    Examples:
      | instante             | dia        |
      | 2026-10-21T02:30:00Z | 2026-10-20 |
      | 2026-10-21T03:00:00Z | 2026-10-21 |

  @AC-08 @D-03 @D-06 @PT-01
  Scenario: BDD-US-02-14 Inicio futuro impide finalizar hasta corregirlo o eliminarlo
    When se guarda el inicio "2026-10-21" en el proyecto sin fin
    Then se acepta la edición y el proyecto sigue "Planificado"
    When se confirma finalizar en el día "2026-10-20" del servidor
    Then se rechaza la finalización por inicio posterior al fin que se registraría
    And conserva el inicio "2026-10-21", el estado "Planificado" y el fin ausente
    When se elimina el inicio y se vuelve a solicitar y confirmar la finalización
    Then queda "Finalizado" con inicio ausente y fin "2026-10-20"

  @AC-09 @D-04 @PT-03
  Scenario: BDD-US-02-15 Repetir finalización otro día y edición idéntica
    Given que la finalización confirmada del proyecto quedó persistida el "2026-10-20"
    And el cliente perdió la respuesta de esa finalización
    When el reloj avanza al día "2026-10-21" y se repite la solicitud confirmada
    Then la respuesta es exitosa y muestra estado "Finalizado" y fin "2026-10-20"
    And no se producen cambios adicionales en el proyecto ni en sus Sprints
    When se guarda dos veces el nombre "Carva II" y el inicio "2026-10-01"
    Then las dos respuestas son exitosas y la consulta conserva una sola identidad con esos datos
    And sigue "Finalizado" con fin "2026-10-20"

  @AC-06 @AC-14 @D-03 @D-04 @PT-02
  Scenario Outline: BDD-US-02-16 Rechazar edición de campos controlados por el sistema
    Given que el proyecto está "<estado>" con inicio ausente y fin <fin>
    When se intenta guardar el cuerpo de edición <cuerpo>
    Then se rechaza la solicitud por campo no editable
    And el proyecto conserva exactamente su identificador, nombre, estado y fechas anteriores

    Examples:
      | estado      | fin          | cuerpo                                        |
      | Planificado | ausente      | {"status":"Finalizado"}                       |
      | Planificado | ausente      | {"end_date":"2026-10-20"}                     |
      | Finalizado  | "2026-10-20" | {"status":"En progreso"}                      |
      | Finalizado  | "2026-10-20" | {"end_date":"2026-10-21"}                     |
      | Finalizado  | "2026-10-20" | {"end_date":null}                             |
      | Finalizado  | "2026-10-20" | {"name":"No guardar","status":"Finalizado"}    |
      | Planificado | ausente      | {"id":"otro-identificador"}                    |

  @AC-05 @AC-14 @D-07 @PT-02
  Scenario Outline: BDD-US-02-17 La API requiere confirmación verdadera sin aceptar fecha elegida
    When se intenta finalizar con el cuerpo <cuerpo>
    Then se rechaza la solicitud con error de validación
    And el proyecto conserva el estado "Planificado" y el fin ausente

    Examples:
      | cuerpo                                             |
      | {}                                                 |
      | {"confirmed":false}                                |
      | {"confirmed":null}                                 |
      | {"confirmed":"true"}                               |
      | {"confirmed":true,"end_date":"2026-10-19"}         |
      | {"confirmed":true,"start_date":"2026-10-01"}       |

  @AC-10 @D-08 @D-09
  Scenario: BDD-US-02-18 Finalizar conserva Sprints existentes y permite consultarlos
    Given que el proyecto está "En progreso" y tiene dos Sprints existentes
    And se conocen los identificadores, objetivos e historias planificadas de ambos
    When se solicita y confirma la finalización
    Then el proyecto queda "Finalizado"
    And los dos Sprints conservan exactamente sus datos y sus historias planificadas
    And no se eliminan Sprints ni se trasladan sus historias como efecto de finalizar
    And ambos pueden consultarse dentro del proyecto

  # D-09 se verifica conservando todas las filas y snapshots existentes.
  # No se modelan estados ni reglas de cierre de Sprints de US-11.

  @AC-11 @D-08 @US-07
  Scenario: BDD-US-02-19 Rechazar un nuevo Sprint en un proyecto Finalizado
    Given que el proyecto está "Finalizado" con fin "2026-10-20"
    When se intenta crear un Sprint con objetivo válido "Otro incremento" por la API
    Then se rechaza la operación por proyecto finalizado
    And no se crea ningún Sprint
    And el proyecto conserva el estado "Finalizado" y el fin "2026-10-20"

  @AC-12 @D-10 @US-03 @US-05 @US-06
  Scenario: BDD-US-02-20 Continuar registrando miembros y manteniendo el Product Backlog
    Given que el proyecto está "Finalizado" con fin "2026-10-20"
    When se registra el miembro "Ana Pérez" desde su sección de miembros
    And se crea desde Product Backlog la historia "Exportar informe" con descripción "Obtener un PDF", prioridad "Alta" y criterio "Permite descargar el informe"
    And se cambia el título de esa historia por "Exportar informe de proyecto" conservando los demás campos editables válidos
    Then al consultar miembros figura "Ana Pérez" con su UUID
    And al consultar Product Backlog figura la historia con el título actualizado y estado "Pendiente"
    And el proyecto conserva el estado "Finalizado" y fin "2026-10-20"
    And las secciones de miembros y Backlog continúan disponibles para esas operaciones

  @AC-09 @D-03 @D-04 @PT-03
  Scenario: BDD-US-02-21 Dos finalizaciones concurrentes tienen un único efecto
    Given que dos clientes prepararon la confirmación del mismo proyecto
    When ambas solicitudes confirmadas se procesan concurrentemente
    Then ambas obtienen una respuesta exitosa con el proyecto "Finalizado"
    And existe un solo proyecto con su identificador original
    And ambas muestran el mismo fin "2026-10-20"
    And las consultas posteriores conservan ese fin

  @AC-11 @D-08 @D-09 @PT-03 @US-07
  Scenario Outline: BDD-US-02-22 Serializar creación de Sprint y finalización
    Given que dos clientes preparan crear el primer Sprint y finalizar el mismo proyecto
    When las operaciones se procesan concurrentemente y se confirma primero "<primera>"
    Then el proyecto termina "Finalizado" con fin "2026-10-20"
    And el resultado de crear el Sprint es "<resultado>"
    And existen <cantidad> Sprints con el objetivo solicitado
    And ninguna operación devuelve el proyecto de Finalizado a En progreso

    Examples:
      | primera     | resultado                        | cantidad |
      | finalizar   | rechazo por proyecto finalizado  | 0        |
      | crear Sprint| creación exitosa                 | 1        |

  @AC-13 @PT-02
  Scenario Outline: BDD-US-02-23 Proyecto inexistente no se crea al editar o finalizar
    Given que no existe el proyecto "00000000-0000-0000-0000-000000000999"
    When se solicita "<operacion>" sobre ese identificador con datos válidos
    Then se informa proyecto no encontrado
    And no se crea ningún proyecto ni se modifica "Carva Demo"

    Examples:
      | operacion                    |
      | editar nombre e inicio       |
      | finalizar con confirmación   |

  @AC-13 @PT-03
  Scenario Outline: BDD-US-02-24 Fallo antes de commit no deja escrituras parciales
    Given que se producirá un fallo de persistencia antes de confirmar los cambios
    When se solicita "<operacion>" con datos válidos
    Then se informa un error controlado sin detalles de almacenamiento
    And al consultar conserva nombre "Carva Demo", estado "Planificado" y fechas ausentes
    And no hay Sprints nuevos ni cambios parciales de esa operación

    Examples:
      | operacion                      |
      | editar nombre e inicio         |
      | finalizar con confirmación     |
      | crear el primer Sprint         |

  @AC-13 @PT-03
  Scenario: BDD-US-02-25 Informar contención sin mostrar éxito
    Given que otra operación mantiene ocupado el acceso de escritura más allá del límite propuesto
    When se solicita finalizar y la espera permitida se agota antes de obtener ese acceso
    Then se informa indisponibilidad temporal y la posibilidad de reintentar
    And la solicitud no altera el proyecto ni muestra finalización exitosa

  @AC-14 @PT-02
  Scenario Outline: BDD-US-02-26 Rechazar cuerpos de edición inválidos
    When se intenta editar con <cuerpo>
    Then se rechaza la solicitud con error de validación
    And el proyecto conserva exactamente sus datos anteriores

    Examples:
      | cuerpo                                   |
      | un cuerpo vacío                          |
      | JSON malformado                          |
      | null                                     |
      | []                                       |
      | "texto"                                  |
      | {}                                       |
      | {"name":123}                             |
      | {"start_date":true}                      |
      | {"start_date":20261001}                  |
      | {"Name":"Carva II"}                      |
      | {"name":"Carva II","extra":true}         |
      | {"name":"Uno","name":"Dos"}              |
      | {"name":"Uno"} {"name":"Dos"}             |
      | un objeto con bytes UTF-8 inválidos       |

  @AC-15 @D-04 @D-06 @PT-03
  Scenario Outline: BDD-US-02-27 Edición de inicio y finalización concurrentes no violan el rango
    Given que un cliente prepara guardar inicio "2026-10-21" y otro finalizar el mismo proyecto
    When ambas operaciones se procesan concurrentemente y confirma primero "<primera>"
    Then el resultado persistido es estado "<estado>", inicio <inicio> y fin <fin>
    And la otra operación se rechaza por incompatibilidad entre inicio y fin
    And nunca se confirman ambas escrituras dejando inicio posterior al fin

    Examples:
      | primera      | estado      | inicio       | fin          |
      | editar inicio| Planificado | "2026-10-21" | ausente      |
      | finalizar    | Finalizado  | ausente      | "2026-10-20" |

  @AC-15 @D-04 @PT-03
  Scenario: BDD-US-02-28 Actualizaciones parciales no sobrescriben otros campos concurrentes
    Given que tres clientes prepararon respectivamente editar nombre, editar inicio y finalizar
    When se procesan concurrentemente el nombre "Carva II", el inicio "2026-10-01" y la finalización confirmada
    Then las tres operaciones son exitosas
    And la consulta final muestra nombre "Carva II", inicio "2026-10-01", fin "2026-10-20" y estado "Finalizado"
    And conserva su identificador original
