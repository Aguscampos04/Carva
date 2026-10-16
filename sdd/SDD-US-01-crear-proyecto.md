# SDD-US-01 — Crear proyecto

**Estado:** Especificación aprobada para US-01.
**Trazabilidad:** requisito §2.1 → EP-01 → US-01 → Issue #1 → este SDD → BDD-US-01 → tests y código.

## 1. Objetivo

Permitir que un integrante registre un proyecto, lo consulte después y obtenga su estado inicial **Planificado**.

## 2. Entradas

El alta recibe únicamente el nombre del proyecto:

| Campo | Obligatorio | Validación |
| --- | --- | --- |
| Nombre | Sí | Debe contener al menos un carácter que no sea espacio en blanco, luego de recortar espacios exteriores. Se conserva el texto ingresado. Sin límite de longitud ni regla de unicidad. |

El identificador lo genera el servidor como UUID. El cliente no puede enviar ni elegir el identificador o el estado.

## 3. Salidas esperadas

- El alta válida devuelve el proyecto persistido, con identificador estable y estado **Planificado**.
- La consulta posterior por identificador devuelve los mismos datos almacenados.
- Una consulta cuyo identificador no existe informa el resultado como no encontrado.
- Un error de validación o persistencia no se presenta como alta exitosa.

## 4. Reglas de negocio

- **BR-01:** todo proyecto creado por US-01 comienza en estado **Planificado**.
- **BR-02:** el estado inicial lo determina el dominio en Go; el cliente no puede seleccionar otro estado.
- **BR-03:** el identificador del proyecto es un UUID estable generado por el servidor.
- **BR-04:** la validación de los campos del proyecto se ejecuta en Go, independientemente de la interfaz.
- **BR-05:** dos proyectos distintos conservan sus datos y pueden consultarse sin sobrescribirse.

## 5. Restricciones

- La historia cubre alta y consulta mínima. No cubre modificación, eliminación, finalización, miembros, backlog, Sprints ni permisos complejos.
- Se reutilizan Go, API REST, SQLite y la separación presentación/aplicación/dominio/persistencia ya establecidas.
- La migración de esquema será versionada y embebida según la infraestructura existente.
- La API de esta historia utiliza JSON. `POST /api/projects` acepta `{"name":"..."}` y responde `201` con el proyecto y un encabezado `Location`; nombre inválido o JSON malformado responde `400`, y fallo de persistencia responde `500`.
- `GET /api/projects/{id}` responde `200` con el proyecto, `404` si no existe o `500` ante fallo de almacenamiento.
- Los contratos REST de las demás historias se definirán en sus SDD.
- El alta solo consume el nombre; las propiedades `id` y `status` enviadas por el cliente no controlan el recurso creado.

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

## 9. Decisiones y límites de esta especificación

1. No se pudo consultar GitHub Issue #1 desde el entorno. Se usa el texto facilitado por el Developer y el borrador local; los criterios de aceptación coinciden.
2. La UI se embebe en el binario con `embed` y se sirve desde el mismo proceso HTTP, para no requerir una URL de API separada.
