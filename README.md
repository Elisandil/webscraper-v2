# WebScraper App

Un proyecto Full‑Stack con propósito formativo para extraer metadatos, links, imágenes y contenido de páginas web, con Go en el backend y React + Tailwind CSS en el frontend. Incluye autenticación JWT completa, programación de tareas con cron, almacenamiento de resultados por usuario en SQLite y un **asistente de chat con IA** que permite crear scraping mediante lenguaje natural.



## Índice
- [WebScraper App](#webscraper-app)
    - [Screenshots](#screenshots)
      - [Página de Inicio (Landing Page)](#página-de-inicio-landing-page)
      - [Autenticación](#autenticación-1)
      - [Dashboard Principal](#dashboard-principal)
      - [Programación de Tareas (Schedules)](#programación-de-tareas-schedules)
      - [Asistente de Chat con IA](#asistente-de-chat-con-ia)
  - [El Viaje de Aprendizaje](#el-viaje-de-aprendizaje)
  - [Lo que realmente aprendí](#lo-que-realmente-aprendí)
  - [Características](#características)
  - [Estructura del Proyecto](#estructura-del-proyecto)
  - [Requisitos](#requisitos)
  - [Arquitectura](#arquitectura)
    - [Domain Layer](#domain-layer)
    - [Infrastructure Layer](#infrastructure-layer)
    - [Use Case Layer](#use-case-layer)
    - [Presentation Layer](#presentation-layer)
    - [Frontend Layer](#frontend-layer)
  - [Dependencias](#dependencias)
  - [Instalación y Uso](#instalación-y-uso)
    - [Setup para nuevos desarrolladores](#setup-para-nuevos-desarrolladores)
  - [Endpoints de la API](#endpoints-de-la-api)
    - [Autenticación](#autenticación)
    - [Scraping](#scraping)
    - [Programación](#programación)
    - [Chat con IA](#chat-con-ia)
    - [Administración](#administración)
    - [Otros](#otros)
  - [Flujo de Uso](#flujo-de-uso)
  - [Autoría](#autoría)
  - [Licencia](#licencia)


## Screenshots

### Página de Inicio (Landing Page)
![Landing Page](assets/landing_page.png)

La página de bienvenida muestra las características principales de la plataforma con un diseño moderno y atractivo.

![Landing Page con Rate Limiting](assets/landing_page_ratelimiting.png)

### Autenticación

**Formulario de Registro**

![Registro de Usuario](assets/register.png)

Formulario completo con validación de campos para crear una nueva cuenta.

**Inicio de Sesión**

![Login de Usuario](assets/login.png)

Sistema de autenticación con JWT para acceso seguro a la plataforma.

### Dashboard Principal

**Vista General de Resultados**

![Vista Principal](assets/results.png)

Dashboard completo con listado de resultados de scraping, filtros y acciones disponibles.

### Programación de Tareas (Schedules)

**Formulario de Nuevo Schedule**

![Formulario de Schedule](assets/schedule_form.png)

Modal para crear tareas programadas con expresiones cron y ejemplos comunes de uso.

**Listado de Schedules**

![Lista de Schedules](assets/schedule_list.png)

Visualización de todas las tareas programadas con información de estado y próxima ejecución.

### Asistente de Chat con IA

**Chat Assistant**

![Chat con IA](assets/ia.png)

Asistente inteligente que permite crear scraping y schedules mediante lenguaje natural, powered by HuggingFace.


## El Viaje de Aprendizaje
Este proyecto nació como un Trabajo de Fin Grado, comenzando de forma muy básica: un backend ligero con un frontend SPA hecho con HTML, CSS y JavaScript vanilla. Aunque simple en apariencia, no careció de desafíos iniciales que despertaron mi curiosidad por el mundo del web scraping.

Una vez presentado y evaluado exitosamente, me di cuenta de que tenía entre mis manos algo con potencial a largo plazo. Me había propuesto entender realmente cómo funcionaba la extracción de datos web, y esta era la mejor forma de conseguirlo. 

Los comienzos fueron duros. Decidí cambiar completamente el stack tecnológico, pasando de Java/Kotlin a Go. Este salto implicó enfrentarme a una nueva sintaxis, la gestión de punteros(algo completamente nuevo viniendo de Java), el no uso de modificadores de acceso (al menos de forma explícita), concurrencia con goroutines, etc.

Probablemente la mayor dificultad la encontré al intentar entender e implementar Clean Architecture. Venía acostumbrado al típico MVC de toda la vida, donde todo era más directo y menos abstracto.
Clean Architecture me obligó a pensar en capas de abstracción reales, implementar inversión de dependencias manualmente,
separar lógica de negocio de infraestructura (sigue siendo dificil), o crear interfaces como contratos y no como herencia.
Fue frustrante al principio, pero al final entendí la escalabilidad que proporciona a tus proyectos así como sus beneficios.

Cuando mi frontend se quedó pequeño, llegó el momento de otra nueva decisión: React.js vs Angular.
Tras investigar ambos frameworks y sacar mis conclusiones:
- Angular: Potente pero pesado, mayor curva de aprendizaje.
- React: Más directo, enfoque funcional, ecosistema maduro.

Al final me decanté por React. A pesar de no estar muy bien implementados en el proyecto, comprendí los beneficios de usar hooks para cambiar los estados y ciclos de vida de los componentes.

Con el tiempo, el proyecto siguió evolucionando. Implementé una arquitectura de componentes más organizada, separando la lógica en páginas, features y componentes UI reutilizables. La integración del Chat Assistant con IA mediante HuggingFace fue un nuevo desafío: procesar lenguaje natural, interpretar intenciones y ejecutar acciones automáticamente. Esto me introdujo en el mundo del NLP (Procesamiento de Lenguaje Natural) y cómo las APIs de IA pueden mejorar la experiencia de usuario.

## Lo que realmente aprendí
Más allá de las tecnologías, este proyecto me enseñó:

**Perseverancia técnica:** Cuando algo no funciona, hay que seguir iterando, e iterando, e iterando, etc ....

**Arquitectura pensada:** No es solo hacer que funcione, sino que funcione bien y sea mantenible.

**Toma de decisiones:** Evaluar tecnologías y elegir la más adecuada según el contexto y objetivos.

**Evolución gradual:** Un proyecto puede crecer y transformarse completamente sin perder su esencia.

**Integración de IA:** Cómo las APIs de procesamiento de lenguaje natural pueden mejorar significativamente la UX.


## Características
- **Clean Architecture** (Domain, Use Cases, Infrastructure, Presentation)
- **Interfaz web moderna** con React 19 + Tailwind CSS y arquitectura basada en páginas
- **🤖 Chat Assistant con IA** - Programa scraping mediante lenguaje natural usando HuggingFace
- **Sistema de contextos React** - Gestión de estado global con Context API
- **Persistencia SQLite** sin CGO con migraciones automáticas
- **Configuración** mediante archivo YAML
- **API REST** completa para operaciones CRUD
- **Autenticación JWT** con roles, refresh tokens y logout seguro
- **Scraping programado** con expresiones cron y gestión de jobs
- **Paginación** de resultados de scraping
- **Rate limiting** por IP para protección de endpoints
- **Middleware** de logging, CORS y control de acceso
- **Extracción completa** de metadatos, links, imágenes y contenido

## Estructura del Proyecto

```
/
├── README.md                              # Documentación del proyecto
├── .gitignore                             # Archivos ignorados por Git
├── assets/                                # Imágenes y recursos del README
├── server/                                # Backend Go con Clean Architecture
│   ├── config.yaml                        # Configuración de servidor, BD, scraping y auth
│   ├── main.go                            # Punto de entrada del backend
│   ├── generate_secret.go                 # Generador de secretos JWT
│   ├── Makefile                           # Comandos de build y desarrollo
│   ├── go.mod, go.sum                     # Dependencias Go
│   ├── data/                              # Base de datos SQLite (scraper.db)
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── entity/                    # Modelos para la lógica del negocio
│   │   │   │   ├── chat.go                # Entidades de chat y NLP (IA)
│   │   │   │   ├── pagination.go          # Entidades de paginación
│   │   │   │   ├── schedule.go            # Entidad Schedule y DTOs
│   │   │   │   ├── scraping.go            # Entidad ScrapingResult
│   │   │   │   └── user.go                # Entidad User y DTOs de auth
│   │   │   └── repository/                # Interfaces de repositorios
│   │   │       ├── schedule.go            # Interface ScheduleRepository
│   │   │       ├── scraping.go            # Interface ScrapingRepository
│   │   │       ├── token_repository.go    # Interface TokenRepository
│   │   │       └── user.go                # Interface UserRepository
│   │   ├── infrastructure/
│   │   │   ├── config/                    # Carga de configuración YAML
│   │   │   │   └── config.go              # Parser de config.yaml
│   │   │   ├── database/                  # SQLite + migraciones automáticas
│   │   │   │   └── sqlite.go              # Conexión y creación de tablas
│   │   │   └── persistence/               # Implementación de repositorios
│   │   │       ├── in_memory_token_repository.go  # Blacklist de tokens
│   │   │       ├── schedule_repository.go         # Repositorio de Schedule
│   │   │       ├── scraping_repository.go         # Repositorio de Scraping
│   │   │       └── user_repository.go             # Repositorio de User
│   │   ├── usecase/                       # Lógica de negocio
│   │   │   ├── auth.go                    # Casos de uso de autenticación
│   │   │   ├── chat.go                    # Casos de uso de chat con IA (HuggingFace)
│   │   │   ├── schedule.go                # Casos de uso de programación
│   │   │   └── scraping.go                # Casos de uso de scraping
│   │   └── presentation/
│   │       ├── server/                    # Servidor HTTP principal
│   │       │   └── server.go              # Inicialización y configuración
│   │       ├── routes/                    # Configuración centralizada de rutas
│   │       │   └── routes.go              # Setup de todas las rutas
│   │       ├── handlers/                  # Controladores HTTP
│   │       │   ├── auth.go                # Handlers de autenticación
│   │       │   ├── chat.go                # Handlers de chat con IA
│   │       │   ├── common.go              # Handlers comunes (health, index)
│   │       │   ├── schedule.go            # Handlers de programación
│   │       │   └── scraping.go            # Handlers de scraping
│   │       ├── middleware/                # Middleware de la aplicación
│   │       │   ├── auth.go                # JWT, roles y autenticación
│   │       │   ├── common.go              # Logging, CORS, content-type
│   │       │   └── rate_limiter.go        # Rate limiting por IP
│   │       └── response/                  # Helpers para respuestas HTTP
│   │           └── helpers.go             # Respuestas estandarizadas
│   └── pkg/                               # Paquetes compartidos
│       ├── crypto/                        # Utilidades de criptografía
│       │   └── password.go                # Hash y validación de contraseñas
│       ├── errors/                        # Errores personalizados
│       │   └── errors.go                  # Tipos de error del dominio
│       └── validator/                     # Validaciones
│           └── validator.go               # Validador de entradas
└── client/                                # Frontend React + Vite + Tailwind
    ├── index.html                         # HTML principal
    ├── package.json                       # Dependencias frontend
    ├── pnpm-lock.yaml                     # Lockfile de pnpm
    ├── vite.config.js                     # Configuración de Vite
    ├── tailwind.config.js                 # Configuración de Tailwind CSS
    ├── postcss.config.js                  # Configuración de PostCSS
    ├── eslint.config.js                   # Configuración de ESLint
    ├── public/
    │   └── vite.svg                       # Logo de Vite
    └── src/
        ├── main.jsx                       # Punto de entrada del frontend
        ├── App.jsx                        # Componente principal con routing
        ├── index.css                      # Estilos globales
        ├── api/
        │   └── client.js                  # Cliente HTTP para la API
        ├── pages/                         # Páginas principales de la app
        │   ├── index.js                   # Barrel export de páginas
        │   ├── Dashboard.jsx              # Dashboard principal con scraping
        │   ├── Landing.jsx                # Página de bienvenida
        │   └── Login.jsx                  # Página de login/registro
        ├── components/                    # Componentes React organizados
        │   ├── chat/                      # Componentes de chat con IA
        │   │   ├── ChatMessage.jsx        # Mensaje individual del chat
        │   │   └── ChatWidget.jsx         # Widget flotante de chat
        │   ├── features/                  # Componentes de funcionalidades
        │   │   ├── schedules/             # Módulo de programación
        │   │   │   ├── index.js           # Barrel export
        │   │   │   ├── ScheduleList.jsx   # Lista de schedules
        │   │   │   └── ScheduleSection.jsx # Sección de schedules
        │   │   └── scraping/              # Módulo de scraping
        │   │       ├── index.js           # Barrel export
        │   │       ├── PaginatedResultsList.jsx # Lista con paginación
        │   │       ├── ResultsList.jsx    # Lista simple de resultados
        │   │       └── ScrapeForm.jsx     # Formulario de scraping
        │   ├── modals/                    # Modales reutilizables
        │   │   ├── index.js               # Barrel export
        │   │   ├── DetailModal.jsx        # Modal de detalles de scraping
        │   │   ├── RegisterModal.jsx      # Modal de registro
        │   │   └── ScheduleModal.jsx      # Modal para crear/editar schedules
        │   └── ui/                        # Componentes UI base
        │       ├── index.js               # Barrel export
        │       ├── Alert.jsx              # Sistema de alertas
        │       ├── HealthIndicator.jsx    # Indicador de estado del servidor
        │       ├── LoadingSpinner.jsx     # Spinner de carga
        │       └── Pagination.jsx         # Controles de paginación
        ├── contexts/                      # Contextos de React (Estado global)
        │   ├── index.jsx                  # Barrel export de contextos
        │   ├── AlertContext.jsx           # Contexto de alertas
        │   ├── AuthContext.jsx            # Contexto de autenticación
        │   ├── ChatContext.jsx            # Contexto de chat con IA
        │   ├── ResultsContext.jsx         # Contexto de resultados
        │   └── ScheduleContext.jsx        # Contexto de schedules
        └── hooks/                         # Custom hooks
            ├── usePagination.jsx          # Hook de paginación
            └── useSchedules.js            # Hook de schedules
```

## Requisitos
- Go ≥ 1.24
- Node.js ≥ 18 y pnpm ≥ 10.24
- SQLite (incluido vía `modernc.org/sqlite`, sin CGO)
- API Token de HuggingFace (opcional, para funcionalidad de chat con IA)

## Arquitectura

### Domain Layer
- `entity/chat.go`: Entidades para chat con IA (ChatRequest, ChatIntent, ChatResponse, ChatConfirmation)
- `entity/pagination.go`: Entidades para paginación de resultados
- `entity/schedule.go`: Entidad Schedule con DTOs de creación y actualización
- `entity/scraping.go`: Entidad ScrapingResult con headers estructurados
- `entity/user.go`: Entidad User con DTOs de login, registro y respuesta de auth
- `repository/schedule.go`: Interface del repositorio de Schedule
- `repository/scraping.go`: Interface del repositorio de Scraping con paginación
- `repository/user.go`: Interface del repositorio de User

### Infrastructure Layer
- `database/sqlite.go`: Conexión SQLite con creación automática de tablas, índices y triggers
- `persistence/schedule_repository.go`: Implementación completa del repositorio de Schedule
- `persistence/scraping_repository.go`: Implementación del repositorio de Scraping con paginación
- `persistence/user_repository.go`: Implementación del repositorio de User con validaciones

### Use Case Layer
- `usecase/auth.go`: Lógica completa de autenticación con JWT, refresh tokens, blacklist y limpieza automática
- `usecase/chat.go`: Lógica de chat con IA usando HuggingFace para interpretación de lenguaje natural y ejecución de acciones
- `usecase/schedule.go`: Lógica de programación con cron jobs, gestión de scheduler y ejecución automática
- `usecase/scraping.go`: Lógica de negocio para scraping con paginación

### Presentation Layer
- `presentation/server/server.go`: Servidor HTTP con inicialización de todos los componentes
- `presentation/routes/routes.go`: Configuración centralizada de rutas con middleware aplicado por grupos
- `presentation/handlers/auth.go`: Controladores HTTP para registro, login, logout, refresh y perfil
- `presentation/handlers/chat.go`: Controladores HTTP para chat con IA (parse y execute)
- `presentation/handlers/common.go`: Controladores HTTP para health check, index y 404
- `presentation/handlers/schedule.go`: Controladores HTTP para CRUD completo de programación
- `presentation/handlers/scraping.go`: Controladores HTTP para scraping y gestión de resultados con paginación
- `presentation/middleware/auth.go`: Middleware JWT con soporte para roles, autenticación opcional y blacklist
- `presentation/middleware/common.go`: Middleware de logging, CORS y content-type
- `presentation/middleware/rate_limiter.go`: Rate limiting por IP con configuración por endpoint
- `presentation/response/helpers.go`: Helpers para respuestas HTTP estandarizadas

### Frontend Layer
**Arquitectura basada en páginas y componentes modulares:**
- `pages/`: Páginas principales (Landing, Login, Dashboard) con routing
- `components/chat/`: Widget de chat con IA para crear scraping y schedules mediante lenguaje natural
- `components/features/`: Módulos de funcionalidades (scraping, schedules) organizados por dominio
- `components/modals/`: Modales reutilizables para detalles, registro y programación
- `components/ui/`: Componentes base de interfaz (alerts, loading, pagination)
- `contexts/`: Context API para gestión de estado global (Auth, Chat, Results, Schedules, Alerts)
- `hooks/`: Custom hooks para lógica reutilizable (usePagination, useSchedules)

## Dependencias

**Backend (Go 1.24.3)**
- `github.com/robfig/cron/v3@v3.0.1`: Programación de tareas con expresiones cron y soporte para segundos
- `github.com/gorilla/mux@v1.8.1`: Router HTTP con soporte para variables de ruta y middleware
- `golang.org/x/net@v0.41.0`: Parsing HTML y herramientas de red para scraping
- `golang.org/x/time@v0.14.0`: Rate limiting y time utilities
- `gopkg.in/yaml.v3@v3.0.1`: Configuración YAML con unmarshaling automático
- `modernc.org/sqlite@v1.37.1`: Driver SQLite sin CGO con rendimiento optimizado
- `github.com/golang-jwt/jwt/v5@v5.2.2`: Generación, validación y parsing de tokens JWT
- `golang.org/x/crypto@v0.39.0`: Hash seguro de contraseñas con bcrypt

**Frontend (React 19 + Vite 7)**
- `react@^19.2.0, react-dom@^19.2.0`: Framework React última versión para interfaz de usuario
- `vite@^7.2.4`: Build tool y dev server ultrarrápido de nueva generación
- `tailwindcss@^3.4.18`: Framework de utilidades CSS para diseño responsive
- `@vitejs/plugin-react@^5.1.1`: Plugin oficial de React para Vite con Fast Refresh
- `eslint@^9.39.1`: Linter para calidad de código JavaScript/JSX
- `autoprefixer@^10.4.22`: PostCSS plugin para prefijos CSS automáticos

## Instalación y Uso

### Setup para nuevos desarrolladores

1. **Clonar el repositorio**

```bash
git clone https://github.com/Elisandil/webscraper-v2
cd webscraper-v2
```

2. **Configurar el backend**

```bash
cd server
cp config.yaml.example config.yaml

# Generar un JWT secret seguro
openssl rand -base64 32

# Editar config.yaml y pegar el secret generado en jwt_secret
nano config.yaml
```

El archivo `config.yaml` debe tener esta estructura:

```yaml
server:
  port: "8080"

database:
  path: "./data/scraper.db"

scraping:
  user_agent: "WebScraper/1.0 (Enhanced Edition)"
  timeout: 30
  max_redirects: 10
  extract_images: true
  extract_favicon: true
  extract_headers: true
  max_links: 100
  max_images: 50

features:
  enable_analytics: true
  enable_caching: false
  cache_duration: 3600

auth:
  require_auth: true
  jwt_secret: "PEGAR_AQUI_EL_SECRET_GENERADO"
  token_duration_hours: 24
  default_role: "user"

chat:
  hf_api_token: "YOUR_HUGGINGFACE_API_TOKEN"  # Opcional: Para funcionalidad de chat con IA
  hf_model_id: "google/flan-t5-small"         # Modelo de HuggingFace a usar
```

**IMPORTANTE:** 
- El archivo `config.yaml` **NO** está en Git por seguridad
- Siempre usa `config.yaml.template` como referencia
- Genera un secret único para cada entorno
- La funcionalidad de chat con IA requiere un token de HuggingFace (obtenerlo en https://huggingface.co/settings/tokens)

3. **Instalar dependencias y levantar el backend**

```bash
go mod download

go run main.go # creará data/ y scraper.db automáticamente
```

El backend:
- Lee la configuración desde `config.yaml`
- Crea automáticamente `data/scraper.db` con tablas, índices y triggers
- Inicia el scheduler para tareas programadas
- Sirve en `http://localhost:8080`

4. **Configurar y levantar el frontend**

```bash
cd client

pnpm install

pnpm run dev
```

El frontend:
- Usa proxy a `http://localhost:8080` (configurado en `vite.config.js`)
- Abre automáticamente `http://localhost:5173`

## Endpoints de la API

### Autenticación
- `POST /api/auth/register` - Registro de usuario
- `POST /api/auth/login` - Autenticación de usuario
- `POST /api/auth/refresh` - Renovar token JWT
- `POST /api/auth/logout` - Logout con revocación de token
- `GET /api/profile` - Obtener perfil del usuario autenticado

### Scraping
- `POST /api/scrape` - Realizar scraping de una URL
- `GET /api/results` - Listar resultados (con paginación opcional: `?page=1&per_page=10`)
- `GET /api/results/{id}` - Obtener resultado específico
- `DELETE /api/results/{id}` - Eliminar resultado

### Programación
- `POST /api/schedules` - Crear tarea programada
- `GET /api/schedules` - Listar tareas del usuario
- `GET /api/schedules/{id}` - Obtener tarea específica
- `PUT /api/schedules/{id}` - Actualizar tarea programada
- `DELETE /api/schedules/{id}` - Eliminar tarea programada

### Chat con IA
- `POST /api/chat/parse` - Interpretar mensaje en lenguaje natural y detectar intención
- `POST /api/chat/execute` - Ejecutar acción detectada (crear scraping o schedule)

### Administración
- `GET /api/admin/*` - Endpoints administrativos (requieren rol "admin")

### Otros
- `GET /api/health` - Health check del servicio
- `GET /` - Interfaz web principal
- `GET /static/*` - Archivos estáticos

## Flujo de Uso

1. **Registro de usuario**
   
```bash
POST /api/auth/register
Content-Type: application/json

{
  "username": "usuario",
  "email": "usuario@ejemplo.com", 
  "password": "password123"
}
```

1. **Autenticación**
   
```bash
POST /api/auth/login
Content-Type: application/json

{
  "username": "usuario",
  "password": "password123"
}
```

Respuesta:
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "usuario",
      "email": "usuario@ejemplo.com",
      "role": "user",
      "active": true
    },
    "expires_at": "2024-01-02T15:04:05Z"
  }
}
```

El token se guarda en `localStorage` para futuras peticiones.

2. **Scraping de URLs**
   
```bash
POST /api/scrape
Authorization: Bearer <token>
Content-Type: application/json

{
  "url": "https://ejemplo.com"
}
```

Devuelve el objeto `ScrapingResult` completo y lo persiste asociándolo al usuario autenticado.

3. **Listar resultados con paginación**
   
```bash
GET /api/results?page=1&per_page=10
Authorization: Bearer <token>
```

Solo devuelve los scrapes del usuario autenticado con información de paginación (opcional).



1. **Programar tareas automáticas**
   
```bash
POST /api/schedules
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Scraping diario de noticias",
  "url": "https://ejemplo.com/noticias",
  "cron_expression": "0 0 9 * * *"
}
```

La tarea se ejecutará automáticamente según la expresión cron (diariamente a las 9:00).

5. **Chat con IA - Crear scraping con lenguaje natural**
   
```bash
POST /api/chat/parse
Authorization: Bearer <token>
Content-Type: application/json

{
  "message": "Escanea https://ejemplo.com cada día a las 9 de la mañana"
}
```

El asistente de IA interpreta el mensaje y responde con la intención detectada:

```json
{
  "message": "Entendido. Quieres programar un scraping de https://ejemplo.com diariamente a las 9:00.",
  "intent": {
    "action": "create_schedule",
    "url": "https://ejemplo.com",
    "frequency": "daily",
    "cron_expr": "0 0 9 * * *",
    "confidence": 0.95
  },
  "needs_confirm": true,
  "action": "schedule"
}
```

Luego confirmas y ejecutas la acción:

```bash
POST /api/chat/execute
Authorization: Bearer <token>
Content-Type: application/json

{
  "confirmed": true,
  "intent_id": "..."
}
```

El sistema crea automáticamente el schedule sin necesidad de especificar manualmente los campos.

6. **Detalles y eliminación**
   
- `GET /api/results/{id}` - Ver detalles completos
- `DELETE /api/results/{id}` - Eliminar resultado

7. **Logout seguro**
   
```bash
POST /api/auth/logout
Authorization: Bearer <token>
```

Revoca el token añadiéndolo a la blacklist hasta su expiración natural.

8. **Health check**
   
```bash
GET /api/health
```

Verifica el estado del servicio y configuración de autenticación.

## Autoría
Desarrollado por Antonio Ortega (aog-dev).

## Licencia
Proyecto desarrollado con fines educativos/formativos.
