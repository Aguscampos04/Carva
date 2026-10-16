# Evidencia de IA — US-03

**Fecha:** 2026-10-16. **Estado:** implementada y verificada técnicamente; revisión y aprobación final del equipo pendientes.
**Issue:** [#3 — Registrar e identificar miembros del proyecto](https://github.com/Aguscampos04/Carva/issues/3).
**Rama:** `codex/us-03-especificacion`. Al comenzar, después de `git fetch origin`, HEAD y `origin/main` coincidían en `c14dc58`.
**Cadena:** enunciado §2.1 → EP-01 → US-03 / #3 → [SDD](../sdd/SDD-US-03-registrar-miembros.md) → AC → [BDD](../bdd/BDD-US-03-registrar-miembros.feature) → pruebas → Go/web → commits y PR. AC-02 conserva su vínculo futuro con US-12 (§2.4) y US-16 (§2.5).

## Instrucciones y autorización

Primera etapa: documentar D-01 a D-07, SDD y BDD, sin código ni tests. Segunda etapa: el responsable autorizó implementar completamente US-03, incorporando JSON estricto, integridad referencial efectiva, claridad de BDD 07/17 y UUID completo secundario; registrar documentación por separado; conservar evidencia TDD en commits; validar, hacer push y abrir un único PR hacia main, sin merge, cierre de Issue ni aprobación final del equipo. US-12 y US-16 quedan fuera de alcance.

Se conservaron los cinco documentos locales preparados, sin cambios ajenos al comenzar. Se consultaron AGENTS.md, skills de gestión/SDD/BDD/TDD/Go/calidad/GitHub, enunciado, contexto, reglas del agente, arquitectura, especificaciones y código existentes. La Issue remota estaba abierta y sin comentarios al consultar.

## Resultado y archivos

- [Contexto](contexto-proyecto.md), [arquitectura](arquitectura.md), SDD y BDD: decisiones y autorización de implementación diferenciadas de la aprobación final.
- [Dominio](../internal/domain/member/member.go): identidad por registro, nombre y recorte Unicode en Go.
- [Aplicación](../internal/application/members/service.go): alta/listado, UUID automático, existencia del proyecto y errores explícitos.
- [Persistencia](../internal/persistence/sqlite/member_repository.go): aislamiento por proyecto, homónimos permitidos. Migraciones nuevas 000006 y 000007, sin editar migraciones previas.
- [Conexión](../internal/persistence/sqlite/db.go): `_pragma=foreign_keys(1)` en el DSN, ejecutado en cada conexión del pool; sin dependencias nuevas.
- [REST](../internal/presentation/httpapi/members.go): POST/GET, JSON estricto, extras ignorados sin controlar identidad/proyecto y errores seguros.
- [Web](../web/static/js/members.js): formulario solo nombre, lista, UUID completo secundario, errores y descarte de respuestas obsoletas; conectada desde HTML/app.js y `cmd/carva`.
- README: instrucciones de uso y verificación.

## Evidencia TDD

Los fallos se ejecutaron antes de la implementación correspondiente. Los commits RED son evidencia del proceso; la rama final pasa las pruebas.

| Etapa | Commit | Resultado real |
| --- | --- | --- |
| Documentación previa | `bf7d55a` | SDD/BDD y ajustes registrados antes del código funcional. |
| RED dominio/aplicación | `7ba3461` | `go test ./internal/domain/member ./internal/application/members` falló por ausencia de `member.New`, `Member`, errores y `members.NewService`. |
| GREEN dominio/aplicación | `48b1dcd` | Mismo comando PASS; vet de ambos paquetes PASS. |
| RED SQLite/REST | `669b782` | `go test ./tests/integration -run 'TestMembers\|TestSQLiteConnection' -count=1` falló por ausencia de `NewMemberRepository` y de la entrada de miembros en `NewHandler`. |
| GREEN SQLite/REST | `dac5d1a` | `go test ./...` y `go vet ./...` PASS. |
| RED web | `c9ed0d6` | `node --test tests/web/members.test.cjs`: seis fallos por ausencia de `members.js`. |
| GREEN web | `c12eda1` | Seis tests Node y suite Go PASS. Se corrigió el conteo del test de resets: debía incluir cambio de proyecto y alta; no era una falla del producto. |
| REFACTOR | `e90c434` | Se centralizaron D-06 en `NormalizeName` y existencia en `requireProject`; Go, vet y Node continuaron pasando. |
| RED almacenamiento exacto | `fc3b82b` | Nombre JSON `"\u0000Ana"` falló con HTTP 500 por el CHECK `length(TEXT)`, que introducía una restricción no aprobada. Se agregaron también comprobaciones de upgrade, atomicidad y rutas. |
| GREEN almacenamiento exacto | `c603567` | Migración aditiva 000007 cambia el CHECK a longitud en bytes y conserva los registros. Go, vet y Node PASS, incluida conservación de UUID al migrar. |

No se alteró el historial para simular RED. La primera medición de cobertura tuvo argumentos separados por PowerShell y se repitió correctamente con comillas; ese error de herramienta no se cuenta como evidencia TDD.

## Trazabilidad BDD y pruebas

Pruebas de [dominio](../internal/domain/member/member_test.go), [aplicación](../internal/application/members/service_test.go), [integración](../tests/integration/members_test.go) y [web](../tests/web/members.test.cjs). Go usa su runner estándar, `httptest` y SQLite temporal real; Node usa runner, `vm` y `assert` incluidos, ejecutando el controlador real con un adaptador DOM mínimo y transporte controlado. No se agregaron frameworks Gherkin ni paquetes de navegador: el `.feature` es la especificación y esta matriz vincula sus comportamientos con pruebas.

| BDD-US-03 | Pruebas automatizadas / alcance |
| --- | --- |
| 01 | `TestMembersRegisterAndList`, `TestMembersHTTPRegistrationIsolationAndReopen`: UUID, nombre y estado del proyecto conservado. |
| 02 | `TestMemberName`, `TestMembersHTTPStrictJSONAndNameValidation`: recorte y conservación. |
| 03 | `TestMemberName`, `TestMembersHTTPStrictJSONAndNameValidation`: nombres ausentes, null, vacíos y blancos. |
| 04 | `TestMembersRegisterAndList`, `TestMembersHTTPRegistrationIsolationAndReopen`; test Node `BDD-07/04/15`: homónimos y UUID secundarios. |
| 05 | `TestMembersRegisterAndList`, `TestMembersHTTPRegistrationIsolationAndReopen`: registros/listas independientes. |
| 06 | `TestMembersHTTPRegistrationIsolationAndReopen`: arreglo vacío real, no null. |
| 07 | Test Node `BDD-07/04/15`: nombre como único dato enviado, actualización y un POST. `TestMembersHTTPDoesNotExposeEditOrDelete`: sin edición/borrado. Estructura HTML y presentación completadas con smoke en navegador. |
| 08 | `TestMembersHTTPStrictJSONAndNameValidation`: cuerpos, tipos, contenido posterior, claves superiores duplicadas y UTF-8 inválido. |
| 09 | `TestMembersHTTPRegistrationIsolationAndReopen`: reabrir archivo conserva UUID/datos. |
| 10 | `TestMembersHTTPMissingAndStorageErrors`: alta inexistente devuelve 404. |
| 11 | `TestMembersHTTPMissingAndStorageErrors`: consulta inexistente devuelve 404, no lista vacía. |
| 12 | `TestMembersErrorsAndValidationPrecedence`, `TestMembersHTTPMissingAndStorageErrors`, `TestMembersSQLiteAbortedWriteLeavesNoRegistration`; test Node `BDD-12/13`: sin alta parcial o aparente. |
| 13 | `TestMembersErrorsAndValidationPrecedence`, `TestMembersHTTPMissingAndStorageErrors`; test Node `BDD-12/13`: lectura fallida no presentada como éxito. |
| 14 | `TestMembersHTTPRegistrationIsolationAndReopen`: campos extra no controlan ID/proyecto ni se devuelven. |
| 15 | `TestMemberName`, `TestMembersHTTPStrictJSONAndNameValidation` y test Node `BDD-07/04/15`: HTML conservado y presentado como texto. |
| 16 | Tests Node `BDD-16` de listas/altas tardías y de refrescos superpuestos. |
| 17 | `TestMemberName`, `TestMembersHTTPStrictJSONAndNameValidation`: ejemplos separados de tabulación, salto de línea y U+00A0 válidos/blancos. |
| 18 | `TestMembersConcurrentHomonyms`: solicitudes simultáneas confirmadas conservan dos identidades. |
| 19 | `TestMembersErrorsAndValidationPrecedence`, `TestMembersHTTPStrictJSONAndNameValidation`: validación antes de resolver proyecto. |

Adicionalmente: tres conexiones simultáneas rechazan huérfanos; UUID duplicado/texto vacío no sobrescriben; trigger ABORT no deja alta parcial; upgrade e idempotencia conservan datos; 000007 conserva identidades. Se cubren nombres con NUL y de un carácter sin agregar restricciones.

## Validaciones ejecutadas

- Go `go1.27.1 windows/amd64`, módulo mínimo 1.26.0; Node `20.17.0`.
- Suite base: `go test ./...` PASS antes de cambios.
- `gofmt -l .`: sin resultados después de `gofmt -w .`. Se normalizaron archivos CRLF que gofmt reportaba; no se generaron cambios adicionales de contenido en el índice Git.
- Implementación final: `go test ./...`, `go vet ./...`, `go build ./...` y `go mod verify`: PASS; módulos verificados.
- `node --test tests/web/members.test.cjs`: seis tests PASS. `node --check` para app.js y members.js: PASS.
- `go test '-coverpkg=./internal/...' '-coverprofile=.local/us-03-coverage.out' ./...`: PASS. Perfil combinado interno **82,7%**; dominio de miembros **100%**, aplicación de miembros **93,1%**, rutas/traducción de errores de miembros **100%**, decoder **90,6%**. La integración está incluida; los 0% de paquetes sin tests locales no representan el perfil combinado.
- No se forzó el fallo del generador criptográfico de UUID ni todos los errores de iteración del driver. No se afirma cobertura exhaustiva ni ejecución del detector de carreras.

## Smoke de interfaz

Se inició Carva con `.local/us-03-smoke.db`, separada de los datos normales. El puerto 8080 estaba ocupado; se utilizó una copia temporal e ignorada del punto de entrada en 18083, sin detener el servicio existente ni cambiar el puerto de producción.

En el navegador se creó **US-03 Demo**, se registraron `  Ana  Pérez  ` y `Ana  Pérez` y se consultó nuevamente: dos entradas con nombre recortado, UUID distintos y estado **Planificado** conservado. `.local/us-03-smoke.jpg` muestra el UUID como dato secundario. Este smoke complementa las pruebas del controlador; no es una suite E2E automática. La automatización completa del layout/browser no se agregó para evitar dependencias innecesarias; su presentación se verificó visualmente.

## Estado de criterios y revisión

- AC-01 y AC-03 a AC-08: comportamiento de US-03 verificado con las pruebas y smoke indicados; AC-06 incluye inspección visual del formulario/UUID.
- AC-02: identidades distinguibles y pertenencia verificadas. Atribución efectiva de votos/esfuerzo **pendiente en US-12/US-16**, fuera de alcance; no se declara validada.
- Revisión técnica del agente: SDD/BDD, cuatro capas, errores seguros, SQL parametrizado, ausencia de secretos y migraciones anteriores intactas.
- Revisión y aprobación final del equipo: **pendientes**. No se marca la historia aprobada ni se cierran checks de revisión humana.
- Sin merge, cierre de Issue ni cambios de Sprint. La publicación es un único PR relacionado con #3, sin palabras de cierre automático; sus datos confirmados se informan al entregar.

La responsabilidad final del producto corresponde al equipo; la asistencia de IA no sustituye su revisión.
