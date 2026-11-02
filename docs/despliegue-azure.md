# Despliegue de Carva en Azure

## Estado verificado

Demo desplegada el 2 de noviembre de 2026. Web y GET /api/projects: HTTP 200.
URL: https://grupo10-carva.braveocean-c591f3f5.switzerlandnorth.azurecontainerapps.io

- Aplicación: grupo10-carva.
- Resource group confirmado por API: frsrutn.
- Entorno compartido: utnapp-env, Express, switzerlandnorth.
- Registro existente: ingcalidad2026.azurecr.io.
- Imagen inicial: grupo10-carva:initial-20261102.
- Digest inicial: sha256:e9b8c8878b57aa9a931a1ec59725f1218c9bea7290d8f5306f913f49b46ad751.
- Issue de trazabilidad: https://github.com/Aguscampos04/Carva/issues/58.

## Decisión explícita: demo con datos temporales

El equipo autorizó publicar la demo con SQLite temporal después de conocer la
limitación. Los datos se pierden al reiniciar, desplegar o escalar a cero.
No usar esta demo como registro permanente ni asumir conservación entre sesiones.
No se migró la base local. El contenedor comienza con una base vacía.

La propuesta de Azure Files no pudo implementarse: Microsoft.Storage no está
registrado y la credencial no tiene Microsoft.Storage/register/action.
Además, el entorno Express no admite Azure Files ni múltiples revisiones.
Registrar Storage por sí solo no resuelve la persistencia en este entorno.
No se crearon almacenamiento ni entornos alternativos facturables.
Para persistencia se necesita acordar con el administrador un entorno compatible
y validar una solución para SQLite; no está implementada ni validada aquí.

## CI/CD

.github/workflows/ci.yml valida formato, tests/cobertura Go, vet, build y frontend.
Después de validar el mismo commit de main, construye en Ubuntu, publica la imagen
con etiqueta SHA y ejecuta scripts/deploy-azure.ps1. Los PR y las ramas codex no
despliegan. Los despliegues están serializados y también existe disparo manual
sobre main. Las acciones oficiales están fijadas a SHA. Permisos: contents: read.

Secretos cifrados configurados en GitHub Actions: AZURE_CREDENTIALS y ACR_PASSWORD.
Variable configurada: CARVA_TEMPORARY_DATA=true, correspondiente a la aprobación
explícita del equipo. Sin esa autorización el script falla, sin crear recursos.
No hay credenciales en archivos. Las credenciales compartidas en el chat requieren
rotación coordinada con el administrador; no se rotaron credenciales de otros grupos.

## Recursos y consumo

0.25 vCPU, 0.5 GiB, mínimo 0 y máximo 1 réplica; HTTPS y puerto interno 8080.
La imagen ejecuta UID/GID 10001 sin shell ni herramientas adicionales.
El escalado a cero reduce cómputo inactivo y también elimina la base temporal.
Esto no garantiza costo cero: registro, imágenes, tráfico y servicios compartidos
pueden facturar. Solo se publica en el repositorio de imágenes grupo10-carva.
No se eliminan imágenes ni recursos de otros equipos. Revisar periódicamente el
consumo y las imágenes propias antiguas con el administrador.

## Evidencia y límites

Build Docker Linux amd64 y push a ACR confirmados. Smoke HTTP de web/API confirmado.
Go test ./..., go vet ./... y 71 tests Node: PASS; gofmt -l . sin salida.
No se atribuye TDD a esta configuración de infraestructura. Persistencia: pendiente.
La comprobación inicial de HTTP no sustituye una prueba funcional completa en Azure.

## Fuentes

- https://learn.microsoft.com/en-us/azure/container-apps/express-overview
- https://learn.microsoft.com/en-us/azure/container-apps/storage-mounts
- https://www.sqlite.org/useovernet.html
- https://docs.github.com/en/actions/reference/security/secure-use
