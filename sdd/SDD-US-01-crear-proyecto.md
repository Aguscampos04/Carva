# SDD-US-01 — Crear proyecto

**Estado:** Borrador en revisión; no habilita todavía la implementación.  
**Trazabilidad:** requisito §2.1 → EP-01 → US-01 → Issue #1 → este SDD → BDD-US-01 → tests y código.

## 1. Objetivo

Permitir que un integrante registre un proyecto, lo consulte después y obtenga su estado inicial **Planificado**.

## 2. Entradas

El conjunto de campos del alta y sus reglas de obligatoriedad todavía no están definidos en la documentación local. La propuesta pendiente de aprobación es:

| Campo | Obligatorio | Validación propuesta |
| --- | --- | --- |
| Nombre | Sí | Debe contener al menos un carácter que no sea espacio en blanco, luego de recortar espacios exteriores. Sin límite de longitud ni regla de unicidad. |

El identificador lo genera el servidor; la propuesta pendiente es UUID. El cliente no puede enviar ni elegir el identificador o el estado.

## 3. Salidas esperadas

- El alta válida devuelve el proyecto persistido, con identificador estable y estado **Planificado**.
- La consulta posterior por identificador devuelve los mismos datos almacenados.
- Una consulta cuyo identificador no existe informa el resultado como no encontrado.
- Un error de validación o persistencia no se presenta como alta exitosa.

## 4. Reglas de negocio

- **BR-01:** todo proyecto creado por US-01 comienza en estado **Planificado**.
- **BR-02:** el estado inicial lo determina el dominio en Go; el cliente no puede seleccionar otro estado.
- **BR-03:** el identificador del proyecto es estable y lo genera el servidor (formato pendiente de aprobación).
- **BR-04:** la validación de los campos del proyecto se ejecuta en Go, independientemente de la interfaz.
- **BR-05:** dos proyectos distintos conservan sus datos y pueden consultarse sin sobrescribirse.

## 5. Restricciones

- La historia cubre alta y consulta mínima. No cubre modificación, eliminación, finalización, miembros, backlog, Sprints ni permisos complejos.
- Se reutilizan Go, API REST, SQLite y la separación presentación/aplicación/dominio/persistencia ya establecidas.
- La migración de esquema será versionada y embebida según la infraestructura existente.
- Los detalles del contrato REST se documentarán junto con la implementación; la arquitectura no ha establecido aún rutas, JSON ni códigos HTTP.
- La especificación de campos marcada como propuesta no se considera decisión del equipo hasta su aprobación.

## 6. Casos límite

- Dos altas consecutivas deben producir registros independientes.
- El estado no puede alterarse enviando un campo `status` en la solicitud.
- La consulta de un identificador desconocido no debe crear ni modificar registros.
- Una falla al guardar no debe producir una respuesta de éxito ni dejar una creación parcial.

## 7. Condiciones de error

- Datos inválidos: rechazo explícito y sin persistencia.
- Proyecto inexistente: resultado de no encontrado.
- Fallo de almacenamiento: error controlado; no se informa un alta exitosa.
- Error interno: se informa de forma segura, sin filtrar detalles SQL al cliente.

## 8. Criterios de aceptación

- **AC-01:** Dado un alta válida, cuando se registra, entonces el proyecto puede consultarse con los datos ingresados.
- **AC-02:** Dado un proyecto recién creado y sin Sprint, cuando se consulta su estado, entonces figura **Planificado**.

## 9. Decisiones pendientes para cerrar este SDD

1. Confirmar campos, obligatoriedad y validaciones del alta (incluida la regla propuesta para nombre).
2. Confirmar el formato del identificador estable (propuesta: UUID generado por el servidor).
3. Verificar el contrato final de la Issue #1 cuando haya acceso a GitHub; por ahora se usa el borrador local y el texto facilitado por el Developer.
