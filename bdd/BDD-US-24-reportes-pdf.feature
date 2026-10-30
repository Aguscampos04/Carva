# Aprobado para implementar el 2026-10-30 por el Agile Enabler. Issue #24 / EP-09 / enunciado §2.9.
Feature: Consultar reportes y generar el PDF final
  Como integrante del equipo
  Quiero consultar reportes de proyecto y de Sprint y descargar sus resultados
  Para revisar y entregar información concordante con las fuentes

  @BDD-01 @AC-01
  Scenario: Consultar el proyecto completo
    Given un proyecto con historias, Sprints, esfuerzo y defectos
    When se genera su reporte de proyecto
    Then se muestran sus datos generales y estado
    And sus historias, estimaciones, planificación, resultados, métricas y defectos
    And el esfuerzo sin Sprint sigue incluido una sola vez en el total de proyecto

  @BDD-02 @AC-01 @AC-03
  Scenario: Conservar los resultados del Sprint original después de un traslado
    Given un Sprint cerrado con una historia incompleta trasladada y completada después
    When se genera el reporte del Sprint original
    Then se conserva su planificación y resultado originales
    And las horas estimadas vigentes se identifican según la excepción de US-16

  @BDD-03 @AC-02
  Scenario: Descargar un PDF de proyecto
    Given un proyecto consultable con datos persistidos
    When se solicita el PDF del proyecto
    Then se descarga un archivo PDF legible con todas las secciones exigidas
    And contiene los mismos valores que el modelo de reporte
    And no finaliza ni modifica el proyecto

  @BDD-04 @AC-01 @AC-03
  Scenario: Distinguir detección y resolución de un defecto
    Given un defecto detectado en un Sprint y resuelto en otro
    When se consulta el reporte de cada Sprint
    Then el defecto aparece con ambas asociaciones identificadas
    And cada métrica conserva la regla de su fuente sin duplicar el defecto

  @BDD-05 @AC-03
  Scenario: Distinguir datos insuficientes y ceros válidos
    Given un reporte con una métrica no disponible y otra cuyo valor es cero
    When se consulta o exporta
    Then la primera muestra "No disponible"
    And la segunda conserva cero
    And una lista vacía no provoca un error ni inventa registros

  @BDD-06 @AC-02
  Scenario: Generar varias páginas con texto extenso y acentos
    Given historias y defectos con textos extensos y acentuados
    When se descarga el PDF
    Then el contenido puede leerse sin recortes ni superposiciones
    And se distribuye en páginas con encabezados comprensibles

  @BDD-07
  Scenario: Rechazar recursos fuera del proyecto o inexistentes
    Given un identificador inválido, inexistente o un Sprint de otro proyecto
    When se solicita un reporte o PDF
    Then se informa el error correspondiente sin datos ajenos

  @BDD-08
  Scenario: Cambiar de selección durante la consulta y reintentar un fallo
    Given una consulta de reporte en curso
    When cambia el proyecto o el Sprint seleccionado
    Then una respuesta tardía no reemplaza el contexto vigente
    And los errores son visibles y se permite reintentar sin mostrar éxito ficticio
