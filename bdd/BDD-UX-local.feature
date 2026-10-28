Feature: Espacio de trabajo por proyecto - Issue 53
  Scenario: UX-01 Abrir proyectos existentes desde Inicio
    Given proyectos guardados en Carva
    When se abre Inicio
    Then Proyectos creados aparece debajo de Abrir por identificador y Crear proyecto
    And cada proyecto puede abrirse sin introducir su identificador

  Scenario: UX-01 Inicio sin proyectos
    Given una base sin proyectos
    When se abre Inicio
    Then se indica que todavía no hay proyectos y cómo crear uno

  Scenario: UX-02 Recuperar navegación y descartar respuestas antiguas
    Given una consulta pendiente del proyecto A
    When se navega al proyecto B y luego responde A
    Then solo se acepta el contexto B
    And recargar el enlace recupera proyecto y sección

  Scenario: UX-03 Filtrar sin modificar historias
    Given historias con diferentes títulos prioridades y estados
    When se aplica búsqueda y filtros
    Then aparecen las coincidencias sin modificar los datos originales
    And las completadas quedan después de las pendientes

  Scenario: UX-04 Plegar Sprint e historias
    Given un Sprint abierto en pantalla
    When se abre una HU
    Then aparecen sus detalles y sus criterios como lista sin otro desplegable
    When se vuelve a plegar la HU
    Then solo aparece su título
    And el detalle del Sprint también puede cerrarse y reabrirse

  Scenario: UX-05 Completar con verificación explícita
    Given una historia sin verificar en En proceso
    When no se confirma la verificación de criterios
    Then no se envía la operación de completitud
    When se confirma y completa la historia
    Then sale de En proceso y aparece en Finalizadas
    And se conserva la pestaña consultada

  Scenario: UX-06 Guardar sin perder la lectura
    Given un Sprint con una HU desplegada
    When se guarda una operación válida
    Then se conservan los despliegues pertinentes y la posición de lectura dentro del alto disponible
    And no se recarga toda la página

  Scenario: UX-07 Responsable estable entre Sprints
    Given una HU con un miembro responsable del mismo proyecto
    When se traslada por quedar incompleta
    Then conserva su responsable
    And puede desasignarse usando Sin responsable

  Scenario: UX-07 Rechazar responsable ajeno
    Given un miembro de otro proyecto
    When se intenta asignarlo a la HU
    Then se rechaza la operación sin cambiar el responsable

  Scenario: UX-08 Registrar horas dentro del Sprint
    Given una HU visible en un Sprint
    When se pulsa Cargar horas trabajadas
    Then se abre un formulario de alta dentro de la HU con el Sprint compatible preseleccionado
    When se guardan horas válidas
    Then se confirma el resultado y vuelve el botón de carga sin navegar ni duplicar el registro

  Scenario: UX-08 Cerrar y reabrir la carga de horas
    Given el formulario de horas abierto sin cambios
    When se pulsa Cerrar
    Then se oculta y puede abrirse nuevamente

  Scenario: UX-09 Conservar un borrador al cancelar su descarte
    Given un formulario de horas modificado
    When se pulsa Cerrar y se cancela la confirmación de descarte
    Then se conserva el formulario y sus valores sin guardar

  Scenario: UX-09 Bloquear doble envío
    Given una escritura pendiente
    When se intenta enviar de nuevo o navegar a otro contexto
    Then no se duplica la escritura ni se abandona el contexto durante la operación

  Scenario: UX-10 Operar en móvil y con teclado
    Given una pantalla móvil o navegación por teclado
    When se utilizan menú formularios y controles
    Then las etiquetas y el foco permiten identificar las acciones y completarlas

  Scenario: UX-11 Datos insuficientes
    Given métricas sin datos suficientes
    When se abre Resumen
    Then se identifica el alcance y se muestra No disponible sin inventar un cero
