Feature: Registrar defectos relacionados con una User Story
  Como integrante del equipo
  Quiero registrar un defecto relacionado con una historia
  Para seguirlo aunque todavía no tenga un Sprint asociado

  # AC-01, AC-02
  Scenario: Registrar un defecto sin Sprint
    Given que existe el proyecto "Carva" con la User Story "A"
    When registro el defecto "No guarda cambios" con severidad "Alta" y estado "Abierto" sin Sprint
    Then el defecto queda asociado a "A" con descripción, severidad y estado indicados
    And detection_sprint_id y resolution_sprint_id quedan nulos
    And no se asigna un Sprint por defecto

  # AC-02
  Scenario: Registrar un defecto con Sprint de detección conocido
    Given que existe el proyecto "Carva" con la User Story "A" y el Sprint 1
    When registro el defecto "Error al exportar" con severidad "Media" y estado "En progreso" detectado en el Sprint 1
    Then el defecto queda asociado a "A"
    And detection_sprint_id conserva el ID del Sprint 1
    And resolution_sprint_id queda nulo

  # AC-03 (DoD)
  Scenario Outline: Aceptar los estados y severidades aprobados
    Given que existe el proyecto "Carva" con la User Story "A"
    When registro un defecto con estado "<estado>" y severidad "<severidad>"
    Then el defecto se registra con esos valores aprobados

    Examples:
      | estado       | severidad |
      | Abierto      | Baja      |
      | Abierto      | Media     |
      | Abierto      | Alta      |
      | Abierto      | Crítica   |
      | En progreso  | Baja      |
      | En progreso  | Media     |
      | En progreso  | Alta      |
      | En progreso  | Crítica   |
      | Resuelto     | Baja      |
      | Resuelto     | Media     |
      | Resuelto     | Alta      |
      | Resuelto     | Crítica   |
      | Cerrado      | Baja      |
      | Cerrado      | Media     |
      | Cerrado      | Alta      |
      | Cerrado      | Crítica   |

  # AC-01/02, condición de error
  Scenario: Rechazar un Sprint que pertenece a otro proyecto
    Given que la User Story "A" pertenece al proyecto "Carva"
    And el Sprint 8 pertenece a otro proyecto
    When registro un defecto para "A" con detection_sprint_id 8
    Then la operación se rechaza como recurso inexistente del proyecto
    And no se crea el defecto
