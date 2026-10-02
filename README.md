# Mortgage Loan Catalogs

Microservicio Gin para consultar catálogos hipotecarios, adaptado a la arquitectura hexagonal de la plantilla corporativa FIF.

## Endpoints

- `GET /health`: estado del servicio.
- `POST /v1/bfcl/mortgage-loan/catalogs/:catalog`: consulta un catálogo.
- `GET /swagger/index.html`: documentación OpenAPI.

La ruta legacy acepta `X-Channel`, `X-Commerce` y `X-Transaction-ID`, y los propaga al backend Java. Para mantener paridad con el servicio anterior, el servidor no rechaza requests cuando esas cabeceras faltan.

`X-Backend-Env` permite seleccionar `real`, `java` o `dummy` para una request. Si no se envía, se usa `DEFAULT_BACKEND`.

## Configuración

Las variables se leen con `config-fif`; credenciales deben inyectarse como secretos del entorno y no guardarse en el repositorio.

| Variable | Uso | Valor por defecto |
|---|---|---|
| `PORT` | Puerto HTTP | `8080` |
| `DEFAULT_BACKEND` | Backend (`real`, `java`, `dummy`) | `real` |
| `FINNFLOW_URL` | URL base Finnflow; requerida para backend `real` | vacío |
| `FINNFLOW_KEY` | Client ID de Finnflow | vacío |
| `FINNFLOW_SECRET` | Client secret de Finnflow | vacío |
| `JAVA_LEGACY_URL` | URL base del proxy Java | vacío |
| `TIMEOUT` | Timeout HTTP en segundos | `10` |
| `APP_NAME` | Nombre del servicio | `mortgage-api-bfcl-mortgage-loans-catalogs` |
| `APP_ENV` / `APP_VERSION` | Ambiente y versión para Datadog | `dev` / `0.0.0` |
| `GIN_MODE` | Modo Gin | `DEBUG` |
| `LOGGING_LEVEL` | Nivel de log | `info` |
| `DD_AGENT_HOST` / `DD_AGENT_PORT` | Agente Datadog | `localhost` / `8126` |
| `DD_PROFILE_ENABLED` | Habilita profiler Datadog | `false` |

Si `DEFAULT_BACKEND=java`, `JAVA_LEGACY_URL` es obligatoria. El backend dummy no requiere URLs externas.

## 🚀 Guía Rápida de Operación (Local y CI/CD)

Requiere Go instalado y acceso corporativo a las dependencias privadas (`lib/go-lib-*`).

### 1. Ubicarse en el proyecto
```bash
# Navegar a la carpeta raíz del microservicio
cd mortgage-api-bfcl-mortgage-loans-catalogs
```

### 2. Descargar las dependencias
Antes de compilar, debes descargar las dependencias del proyecto. Dado que consumes librerías privadas de tu empresa (`github.com/falabella-regulado/*`), asegúrate de indicarle a Go que no utilice el proxy público para ellas:
```bash
# Configurar Go para que tenga acceso a repositorios privados
export GOPRIVATE="github.com/falabella-regulado/*"
# En Windows (PowerShell): $env:GOPRIVATE="github.com/falabella-regulado/*"

# Descargar módulos
go mod download
go mod tidy
```

### 3. Compilar el proyecto
El proyecto cuenta con un archivo `Makefile` para facilitar las tareas. Para generar el binario de ejecución:
```bash
make build
# El binario compilado quedará listo en la ruta: build/bin/dist
```

### 3. Ejecutar la suite de pruebas
Para ejecutar todos los tests unitarios y la suite E2E de paridad:
```bash
make test
```

### 4. Calcular la Cobertura (Coverage >95%)
Para extraer los metadatos de cobertura requeridos por SonarQube:
```bash
make coverage
```
*(Opcional) Para visualizar qué líneas exactas están cubiertas en tu navegador:*
```bash
go tool cover -html=coverfile_out
```

### 5. Inyección de Variables de Entorno (Arranque)
Para conectarse al backend real (Finnflow) garantizando la paridad con el legado Java, inyecta las siguientes variables al correr el proyecto. Elige los comandos según tu sistema operativo:

**🍎 Mac / Linux (Bash / Zsh):**
```bash
# Variables Obligatorias del Legado
export FINNFLOW_URL="https://api-proveedor.com"
export FINNFLOW_KEY="tu_client_id"
export FINNFLOW_SECRET="tu_client_secret"

# Selector de tráfico (Nuevo en Go)
export DEFAULT_BACKEND="real"

# Levantar el servicio
go run ./cmd/api
```

**🪟 Windows (PowerShell):**
```powershell
# Variables Obligatorias del Legado
$env:FINNFLOW_URL="https://api-proveedor.com"
$env:FINNFLOW_KEY="tu_client_id"
$env:FINNFLOW_SECRET="tu_client_secret"

# Selector de tráfico (Nuevo en Go)
$env:DEFAULT_BACKEND="real"

# Levantar el servicio
go run ./cmd/api
```

**☁️ Nullplatform (Entorno Productivo/CI-CD):**
Cuando despliegues tu servicio en Nullplatform, **no es necesario** ejecutar comandos en consola. Simplemente debes configurar las variables `FINNFLOW_URL`, `FINNFLOW_KEY` y `FINNFLOW_SECRET` como *Parameters* (o *Secrets* en el caso de la contraseña) dentro de la UI de Nullplatform para tu aplicación. 

*(Nota: La variable `DEFAULT_BACKEND` no debe registrarse en Nullplatform, ya que es exclusiva para pruebas locales; el microservicio asumirá inteligentemente el backend "real" en producción por defecto).*

Adicionalmente, ten la tranquilidad de que Nullplatform inyectará automáticamente la variable `PORT`, la cual tu servicio Go ya está configurado para leer dinámicamente y ejecutar el servidor bajo ese puerto (con soporte para *Graceful Shutdown* incorporado en `main.go`).

### 6. ¿Cómo consumir el servicio?
Dado que este microservicio imita el contrato legacy de Java, el cliente (Frontend/Consumer) no requiere enviar un JSON Body, solo el Path Variable con el nombre del catálogo.

**Ejemplo de Request (cURL):**
```bash
curl --location --request POST 'http://localhost:8080/v1/bfcl/mortgage-loan/catalogs/Destino' \
--header 'X-Channel: APP' \
--header 'X-Commerce: FALABELLA'
```
*Nota: También hemos anexado el archivo `postman_collection.json` en la raíz con todos los casos de uso documentados.*
