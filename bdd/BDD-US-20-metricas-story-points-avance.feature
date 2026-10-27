Feature: Consultar métricas de Story Points y avance
  Como integrante del equipo
  Quiero consultar las métricas por Sprint y proyecto
  Para medir el avance registrado.

  Scenario: BDD-US-20-01 Sprint con puntos completos
    Given un Sprint con 3 historias planificadas con 2, 3 y 5 Story Points definitivos
    And 2 de las 3 historias están completadas
    When consulto las métricas del Sprint
    Then los Story Points planificados son 10
    And los Story Points completados son la suma de las 2 historias completadas
    And el porcentaje completado es 66.67 por ciento
    And la velocidad solo está disponible si el Sprint está cerrado

  Scenario: BDD-US-20-02 Sprint con estimación definitiva faltante
    Given un Sprint con una historia planificada sin Story Points definitivos
    When consulto sus métricas
    Then los Story Points planificados son «No disponible»
    And la métrica de historias completadas sigue usando los estados registrados

  Scenario: BDD-US-20-03 Sprint vacío y valor cero
    Given un Sprint vacío y otro Sprint con una historia estimada en 0 puntos
    When consulto sus métricas
    Then las sumas de puntos de cada Sprint son 0
    And el porcentaje del Sprint vacío es «No disponible»
    And el valor 0 de la historia estimada se considera un Story Point válido

  Scenario: BDD-US-20-04 agregado de proyecto usa snapshots
    Given dos Sprints con snapshots que suman 15 Story Points planificados
    And una historia incompleta se trasladó y tiene un snapshot en ambos Sprints
    When consulto las métricas del proyecto
    Then la historia trasladada se cuenta en cada Sprint donde fue planificada
    And la suma de puntos del proyecto se obtiene de esos snapshots
    And las ediciones posteriores del Backlog no cambian el resultado

  Scenario: BDD-US-20-05 proyecto sin Sprints
    Given un proyecto sin Sprints ni historias planificadas
    When consulto sus métricas
    Then las métricas del proyecto son «No disponible»
    And la lista de Sprints está vacía
    And la consulta no modifica los datos del proyecto

  Scenario: Sprint de otro proyecto
    Given un Sprint que pertenece a otro proyecto
    When consulto sus métricas indicando el proyecto equivocado
    Then la API responde que el recurso no existe
