# Aprobado para implementar el 2026-10-27 por @Aguscampos04. §2.8 → EP-08 → US-23 / #23 → SDD-US-23.
Feature: Consultar el dashboard del proyecto
  Como integrante del equipo
  Quiero ver el estado y gráficos de métricas del proyecto
  Para comprender su situación

  @US23-01 @AC-01
  Scenario: Consultar el estado oficial
    Given un proyecto cuyo estado consultado es "En progreso"
    When se abre su dashboard
    Then se muestra "En progreso" para ese mismo proyecto
    And no se calcula otro estado desde sus métricas

  @US23-02 @AC-01
  Scenario: Graficar las métricas existentes del proyecto y del Sprint
    Given un proyecto con métricas de puntos, avance, horas y defectos
    When se consulta su dashboard y se selecciona un Sprint
    Then cada gráfico muestra los valores de la consulta correspondiente
    And las series identifican sus unidades y valores textuales
    And el agregado del proyecto conserva los valores de su fuente

  @US23-03 @AC-02
  Scenario: Distinguir una métrica no disponible de un cero válido
    Given una velocidad no disponible y cero defectos detectados
    When se consulta el dashboard
    Then la velocidad muestra "No disponible" sin barra numérica
    And los defectos detectados muestran cero como dato válido

  @US23-04 @AC-01
  Scenario: Representar una desviación negativa
    Given que la consulta de horas devuelve una desviación de menos dos horas
    When se consulta el dashboard
    Then la desviación conserva su signo y su valor
    And el gráfico la representa a la izquierda del origen

  @US23-05
  Scenario: Cambiar de selección mientras hay una consulta pendiente
    Given una consulta pendiente para un proyecto o Sprint anterior
    When el usuario cambia o quita esa selección
    Then se limpian los gráficos correspondientes a la selección anterior
    And una respuesta tardía no reemplaza la selección vigente

  @US23-06 @AC-01
  Scenario: Consultar un proyecto sin Sprints
    Given un proyecto sin Sprints
    When se consulta su dashboard
    Then se muestran el estado y las métricas de proyecto según sus fuentes
    And se informa que no hay Sprints sin inventar métricas para uno

  @US23-07
  Scenario: Reintentar una consulta fallida
    Given que falló la consulta de métricas del contexto seleccionado
    When el usuario vuelve a consultar y la fuente responde correctamente
    Then se reemplaza el error por los resultados de esa selección
    And ninguna consulta modifica los datos del proyecto
