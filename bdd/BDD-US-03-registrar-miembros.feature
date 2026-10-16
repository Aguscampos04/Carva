# Estado: autorizado para implementar por el responsable el 2026-10-16.
# Aprobación final de la historia por el equipo pendiente.
# Fuente: docs/enunciado.md §2.1; soporte futuro para §2.4 y §2.5.
# EP-01 -> US-03 -> https://github.com/Aguscampos04/Carva/issues/3
# SDD: sdd/SDD-US-03-registrar-miembros.md
# Automatización y evidencia se registrarán en docs/evidencia-ia-us-03.md.

@US-03 @Issue-3
Feature: Registrar e identificar miembros del proyecto
  Como integrante del equipo
  Quiero registrar e identificar miembros del proyecto
  Para atribuirles votos de Planning Poker y registros de esfuerzo cuando esas funciones existan

  @BDD-US-03-01 @AC-01 @AC-08
  Scenario: Registrar un miembro con solo su nombre
    Given que existe el proyecto "Carva Demo" sin miembros
    When se registra a "Ana Pérez" en ese proyecto ingresando únicamente su nombre
    Then se confirma el alta con un UUID asignado automáticamente
    And la lista del proyecto contiene un miembro con ese UUID y el nombre "Ana Pérez"
    And el proyecto conserva su estado y sus demás datos

  @BDD-US-03-02 @AC-03
  Scenario Outline: Recortar solo los extremos del nombre
    Given que existe el proyecto "Carva Demo"
    When se registra un miembro con el nombre <entrada>
    Then el alta se confirma con el nombre <guardado>
    And la lista conserva exactamente ese nombre

    Examples:
      | entrada           | guardado        |
      | "  Ana  Pérez  "  | "Ana  Pérez"    |
      | "  a  "          | "a"             |
      | "  Éva O'Neill " | "Éva O'Neill"   |

  @BDD-US-03-03 @AC-03
  Scenario Outline: Rechazar un nombre obligatorio inválido
    Given que existe el proyecto "Carva Demo" sin miembros
    When se solicita el alta con el cuerpo JSON <cuerpo>
    Then se informa que el nombre es obligatorio con HTTP 400
    And la lista del proyecto sigue vacía

    Examples:
      | cuerpo          |
      | {}              |
      | {"name":null}   |
      | {"name":""}     |
      | {"name":"   "}  |

  @BDD-US-03-04 @AC-02 @AC-04 @AC-06
  Scenario: Distinguir homónimos en un mismo proyecto
    Given que existe el proyecto "Carva Demo" con un miembro llamado "Ana Pérez"
    When se registra otro miembro con el nombre "  Ana Pérez  " en ese proyecto
    Then el alta se confirma sin reemplazar al miembro anterior
    And la lista contiene dos miembros llamados "Ana Pérez" con UUID distintos
    And la sección muestra el nombre como dato principal y el UUID completo como dato secundario de cada uno
    And cada entrada conserva un UUID distinto para la selección en funciones futuras

  @BDD-US-03-05 @AC-02 @AC-04 @AC-05
  Scenario: Registrar a la misma persona en proyectos independientes
    Given que existen los proyectos "Carva Demo" y "Carva Alternativo"
    And "Ana Pérez" está registrada en "Carva Demo"
    When se registra a "Ana Pérez" en "Carva Alternativo"
    Then el nuevo registro tiene un UUID distinto al de "Carva Demo"
    And la lista de cada proyecto contiene únicamente su propio registro de "Ana Pérez"
    And el registro de "Carva Demo" conserva sus datos

  @BDD-US-03-06 @AC-05
  Scenario: Consultar un proyecto sin miembros
    Given que existe el proyecto "Carva Demo" sin miembros
    When se consulta su lista de miembros
    Then se responde HTTP 200 con el arreglo JSON vacío []
    And no se crea ningún registro

  @BDD-US-03-07 @AC-01 @AC-06
  Scenario: Registrar desde la sección de miembros del proyecto
    Given que se está consultando el proyecto "Carva Demo" y su sección de miembros
    And se muestran el formulario con nombre como único dato y la lista del proyecto
    And no se ofrecen edición, eliminación, autenticación ni permisos
    When se envía el nombre "Bruno" y el servidor confirma el alta
    Then la lista muestra a "Bruno" con su UUID sin volver a registrarlo

  @BDD-US-03-08 @AC-07
  Scenario Outline: Rechazar cuerpos que incumplen el contrato de alta
    Given que existe el proyecto "Carva Demo" sin miembros
    When se solicita el alta con <entrada>
    Then se responde HTTP 400 con un error JSON controlado
    And no se registra ningún miembro

    Examples:
      | entrada                              |
      | un cuerpo vacío                      |
      | el JSON incompleto {"name":           |
      | dos objetos JSON consecutivos {} {}   |
      | el cuerpo JSON null                  |
      | el cuerpo JSON []                    |
      | el cuerpo JSON "Ana"                 |
      | el cuerpo JSON {"name":42}           |
      | el cuerpo JSON {"name":true}         |
      | el cuerpo JSON {"name":[]}           |
      | el cuerpo JSON {"name":{}}           |

  @BDD-US-03-09 @AC-01 @AC-08
  Scenario: Conservar identidad y datos al reiniciar
    Given que "Ana Pérez" fue registrada en "Carva Demo" y se conoce su UUID
    When se reinicia la aplicación usando la misma base y se vuelve a consultar el proyecto
    Then su lista conserva a "Ana Pérez" con el mismo UUID y proyecto
    And la consulta no crea un nuevo registro

  @BDD-US-03-10 @AC-05 @AC-07
  Scenario: Rechazar un alta para un proyecto inexistente
    Given que no existe el proyecto "00000000-0000-0000-0000-000000000999"
    When se intenta registrar a "Ana Pérez" en ese proyecto
    Then se informa que el proyecto no existe con HTTP 404
    And no se crea ningún proyecto ni miembro

  @BDD-US-03-11 @AC-05 @AC-07
  Scenario: Diferenciar proyecto inexistente de lista vacía
    Given que no existe el proyecto "00000000-0000-0000-0000-000000000999"
    When se consulta su lista de miembros
    Then se informa que el proyecto no existe con HTTP 404
    And no se presenta una lista vacía como resultado exitoso

  @BDD-US-03-12 @AC-07
  Scenario: Informar una falla de escritura sin alta aparente
    Given que existe el proyecto "Carva Demo" sin miembros
    And el almacenamiento falla antes de confirmar la escritura
    When se intenta registrar a "Ana Pérez"
    Then se informa un error controlado con HTTP 500 sin detalles internos
    And no queda un miembro parcial
    And la sección no agrega a "Ana Pérez" como si estuviera registrada

  @BDD-US-03-13 @AC-07
  Scenario: Informar una falla al listar sin presentar información incompleta
    Given que existe el proyecto "Carva Demo" con miembros registrados
    And el almacenamiento falla al consultar la lista
    When se solicita la lista del proyecto
    Then se informa un error controlado con HTTP 500 sin detalles internos
    And no se presenta una lista vacía o parcial como consulta exitosa

  @BDD-US-03-14 @AC-07 @AC-08
  Scenario: El cliente no controla UUID ni pertenencia mediante campos adicionales
    Given que existen los proyectos "Carva Demo" y "Carva Alternativo"
    When se registra a "Ana Pérez" en la ruta de "Carva Demo" enviando además un id elegido, el project_id de "Carva Alternativo" y un email
    Then se confirma el alta con HTTP 201 y un UUID nuevo generado por el servidor
    And la respuesta contiene únicamente id, project_id y name
    And el miembro pertenece a "Carva Demo" y no aparece en "Carva Alternativo"
    And no se guarda el email ni el id elegido

  @BDD-US-03-15 @AC-06
  Scenario: Mostrar como texto un nombre que contiene marcas HTML
    Given que existe el proyecto "Carva Demo"
    When se registra el nombre "<b>Ana</b>"
    Then la lista muestra literalmente "<b>Ana</b>" como texto
    And no interpreta el nombre como HTML

  @BDD-US-03-16 @AC-06
  Scenario: No mezclar listas al cambiar de proyecto
    Given que se solicitó la lista de "Carva Demo" y su respuesta está pendiente
    And se cambia al proyecto "Carva Alternativo"
    When llega la respuesta pendiente de "Carva Demo"
    Then esa respuesta no se muestra como lista de "Carva Alternativo"

  @BDD-US-03-17 @AC-03
  Scenario Outline: Aplicar el recorte a espacios en blanco exteriores
    Given que existe el proyecto "Carva Demo"
    When se solicita el alta con el cuerpo JSON <cuerpo>
    Then el resultado es <resultado>

    Examples:
      | cuerpo                                      | resultado                                                   |
      | {"name":"\t\n\u00a0Ana  Pérez\u00a0\n\t"} | HTTP 201 y nombre guardado exactamente "Ana  Pérez"           |
      | {"name":"\t\n\u00a0"}                     | HTTP 400 sin crear ningún miembro                           |

  @BDD-US-03-18 @AC-04
  Scenario: Conservar altas simultáneas de homónimos
    Given que existe el proyecto "Carva Demo" sin miembros
    When se procesan simultáneamente dos altas válidas con el nombre "Ana Pérez" y ambas se confirman
    Then la lista contiene dos miembros con ese nombre y UUID distintos
    And ninguna alta sobrescribe a la otra

  @BDD-US-03-19 @AC-07
  Scenario: Validar el nombre antes de buscar el proyecto en el alta
    Given que no existe el proyecto "00000000-0000-0000-0000-000000000999"
    When se solicita el alta con el cuerpo JSON {"name":"   "}
    Then se informa que el nombre es obligatorio con HTTP 400
    And no se crea ningún proyecto ni miembro

# AC-02: estos escenarios especifican la disponibilidad de identidades distintas.
# La atribución efectiva de votos y esfuerzo debe especificarse y validarse en
# US-12 y US-16 cuando existan. No se afirma cobertura ejecutada de esas funciones.
