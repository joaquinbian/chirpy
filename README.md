# 🐦 Chirpy

API REST desarrollada en Go que simula una plataforma de microblogging similar a Twitter.
Proyecto con fines educativos.

## 📋 Descripción

Chirpy es un backend completo que permite a los usuarios registrarse, autenticarse y publicar mensajes cortos llamados "chirps". La API implementa autenticación JWT, refresh tokens, filtrado de contenido y un sistema de suscripción premium.

## 🛠️ Tecnologías

- **Go 1.25** - Lenguaje de programación
- **PostgreSQL** - Base de datos relacional
- **JWT** - Autenticación con access y refresh tokens
- **Argon2id** - Hashing seguro de contraseñas
- **SQLC** - Generación de código SQL type-safe
- **Goose** - Migraciones de base de datos

## 🚀 Funcionalidades

- **Autenticación completa**: Registro, login, refresh tokens y revocación
- **CRUD de Chirps**: Crear, leer y eliminar posts (máximo 140 caracteres)
- **Filtro de profanidad**: Censura automática de palabras prohibidas
- **Sistema de suscripción**: Upgrade a "Chirpy Red" vía webhooks
- **Ordenamiento y filtrado**: Chirps por autor y orden cronológico

## 📁 Estructura del Proyecto

```
chirpy/
├── internal/
│   ├── auth/          # Autenticación (JWT, hashing, tokens)
│   ├── database/      # Código generado por SQLC
│   └── utils/         # Utilidades (filtro de profanidad)
├── sql/
│   ├── queries/       # Queries SQL para SQLC
│   └── schema/        # Migraciones de Goose
├── main.go            # Punto de entrada y rutas
└── handle_*.go        # Handlers de la API
```

## 🔌 API Endpoints

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/api/users` | Registrar usuario |
| `POST` | `/api/login` | Iniciar sesión |
| `PUT` | `/api/users` | Actualizar usuario |
| `POST` | `/api/chirps` | Crear chirp |
| `GET` | `/api/chirps` | Listar chirps |
| `GET` | `/api/chirps/{id}` | Obtener chirp por ID |
| `DELETE` | `/api/chirps/{id}` | Eliminar chirp |
| `POST` | `/api/refresh` | Renovar access token |
| `POST` | `/api/revoke` | Revocar refresh token |
| `POST` | `/api/polka/webhooks` | Webhook de suscripción |

## ⚙️ Configuración

1. Clonar el repositorio
2. Crear archivo `.env`:
```env
DB_URL=postgres://usuario:password@localhost:5432/chirpy?sslmode=disable
SECRET_KEY=tu_clave_secreta_jwt
POLKA_KEY=clave_para_webhooks
PLATFORM=dev
```

3. Ejecutar migraciones:
```bash
goose -dir sql/schema postgres "tu_connection_string" up
```

4. Iniciar el servidor:
```bash
go run .
```

El servidor estará disponible en `http://localhost:8080`

## 🧪 Tests

```bash
go test ./...
```

## 📚 Aprendizajes

Este proyecto me permitió practicar:

- Diseño de APIs RESTful en Go
- Autenticación stateless con JWT
- Manejo seguro de contraseñas con Argon2id
- Arquitectura de aplicaciones web
- Migraciones y manejo de base de datos
- Testing en Go

---
