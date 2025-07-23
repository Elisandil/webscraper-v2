# WebScraper App

Una aplicación de webscraping construida con Go 1.24 siguiendo principios de Clean Architecture.

## Características

- Clean Architecture (Domain, Use Cases, Infrastructure)
<<<<<<< HEAD
- Interfaz web con HTML, CSS (Tailwind) y JavaScript (SPA)
=======
- Interfaz web moderna con HTML, CSS (Tailwind) y JavaScript
>>>>>>> master
- Persistencia SQLite sin CGO
- Configuración mediante archivo YAML
- API REST para operaciones CRUD

## Estructura del Proyecto

```
webscraper/
├── main.go                          # Punto de entrada
├── config.yaml                      # Configuración
├── config/                          # Carga de configuración
<<<<<<< HEAD
├── domain/                          # Entidades e interfaces
=======
├── domain/                          # Entidades y interfaces
>>>>>>> master
├── usecase/                         # Lógica de negocio
├── infrastructure/                  # Implementaciones concretas
├── interface/
    ├── static/
        └── js/                      # Js para el HTML
    └── /templates/                  # Templates HTML
```

## Requisitos

- Go 1.24+
- Puerto 8080 disponible

## Instalación y Uso

1. **Configurar:**
```bash
cd webscraper
make setup
```

2. **Ejecutar: (Testing)**
```bash
make run
```

3. **Acceder:**
Abre http://localhost:8080 en tu navegador


## API Endpoints

- `GET /` - Interfaz web
- `POST /api/scrape` - Extraer datos de URL
- `GET /api/results` - Listar todos los resultados
- `GET /api/results/{id}` - Obtener resultado específico
- `DELETE /api/results/{id}` - Eliminar resultado

## Configuración

El archivo `config.yaml` permite configurar:

```yaml
server:
  port: "8080"

database:
  path: "./data/scraper.db"

scraping:
  user_agent: "WebScraper/1.0"
  timeout: 30
```

## Arquitectura

### Domain Layer
- `entity/scraping.go`: Entidad ScrapingResult
- `repository/scraping.go`: Interface del repositorio

### Use Case Layer
- `usecase/scraping.go`: Lógica de negocio para scraping

### Infrastructure Layer
- `database/sqlite.go`: Conexión SQLite
- `repository/scraping_repository.go`: Implementación del repositorio
- `web/server.go`: Servidor HTTP y handlers

## Dependencias

- `github.com/gorilla/mux`: Router HTTP
- `golang.org/x/net`: Parsing HTML
- `gopkg.in/yaml.v3`: Configuración YAML
- `modernc.org/sqlite`: Driver SQLite sin CGO

## Comandos Make

- `make build` - Compilar aplicación
- `make run` - Ejecutar aplicación
- `make test` - Ejecutar tests
- `make clean` - Limpiar archivos generados
- `make deps` - Instalar dependencias
- `make setup` - Configuración inicial

## Desarrollo

<<<<<<< HEAD
La aplicación sigue un patrón de Clean Architecture:
=======
La aplicación sigue Clean Architecture:
>>>>>>> master

1. **Domain**: Entidades y reglas de negocio
2. **Use Cases**: Lógica de aplicación
3. **Infrastructure**: Detalles técnicos (DB, Web, etc.)

Los datos extraídos incluyen:
- Título de la página
- Meta descripción
- Enlaces encontrados
<<<<<<< HEAD
- Timestamp de extracción
=======
- Timestamp de extracción
>>>>>>> master
