# US-01 — Crear proyecto
# Fuente: docs/enunciado.md §2.1, borrador local de Issue #1 y SDD-US-01.
# La validez de los datos y el formato del identificador esperan aprobación del contrato de alta.

Feature: Crear y consultar un proyecto
  Como integrante del equipo
  Quiero crear un proyecto
  Para organizar su trabajo

  # AC-01 → BDD-01 → test de aplicación/integración; presentación/API y SQLite.
  Scenario: Crear un proyecto y consultar los datos registrados
    Given que se proporcionan los campos válidos acordados para el alta
    When se registra el proyecto y luego se consulta por su identificador
    Then se recuperan los datos ingresados
    And se recupera el mismo identificador asignado al proyecto

  # AC-02 → BDD-02 → test de dominio/aplicación e integración.
  Scenario: Un proyecto nuevo comienza Planificado
    Given que todavía no existe un Sprint para el proyecto
    When se crea y consulta el proyecto
    Then su estado es "Planificado"

  # BR-05 → BDD-03 → test de persistencia e integración.
  Scenario: Dos proyectos conservan registros independientes
    Given que se proporcionan datos válidos para dos proyectos distintos
    When se crean ambos proyectos y se consulta cada uno por su identificador
    Then cada consulta devuelve los datos del proyecto correspondiente
    And ninguno de los dos proyectos sobrescribe al otro

  # BR-02 → BDD-04 → test de dominio/API.
  Scenario: El cliente no puede escoger un estado distinto al inicial
    Given que la solicitud de alta intenta enviar un estado diferente de "Planificado"
    When se registra el proyecto
    Then el estado del proyecto creado sigue siendo "Planificado"

  # AC-01 y contrato de validación pendiente → BDD-05 → test de dominio/aplicación.
  Scenario: Rechazar datos que no cumplen las validaciones aprobadas
    Given que los datos del alta incumplen una validación aprobada
    When se intenta registrar el proyecto
    Then se informa un error de validación
    And no se crea ningún proyecto

  # Error de consulta → BDD-06 → test de consulta/API.
  Scenario: Consultar un proyecto inexistente
    Given que no existe un proyecto con el identificador consultado
    When se solicita el proyecto
    Then se informa que no fue encontrado

  # Restricción de persistencia → BDD-07 → test de aplicación/integración.
  Scenario: Una falla de persistencia no se informa como alta exitosa
    Given que el almacenamiento rechaza la creación
    When se intenta registrar el proyecto
    Then se informa un error de persistencia
    And no se devuelve una creación exitosa
