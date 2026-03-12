# ⚽ Go Top Scorers API

API REST desarrollada en Go utilizando únicamente la librería estándar, sin frameworks externos. Gestiona un ranking de los 10 máximos goleadores históricos del fútbol con operaciones CRUD, filtros por query parameters, path parameters y persistencia real en JSON.

---

## 📁 Estructura del proyecto

```
go-http/
├── data/
│   └── players.json     # Base de datos en JSON
├── main.go              # Toda la lógica de la API
├── Dockerfile
├── docker-compose.yml
├── .gitignore
└── README.md
└── ----
```

---

## 🛠️ Tecnologías utilizadas

- **Go** — lenguaje principal
- **Librería estándar de Go** — `net/http`, `encoding/json`, `os`
- **Docker** — contenedor para correr el servidor
- **JSON** — almacenamiento y persistencia de datos

---

## 🚀 Cómo correr el proyecto con Docker

### Requisitos
- Tener instalado [Docker Desktop](https://www.docker.com/products/docker-desktop/)

### Pasos

**1. Clonar el repositorio**
```bash
git clone https://github.com/Pablownski/Ejercicio4-API-JSON
cd go-http
```

**2. Buildear y levantar el contenedor**
```bash
docker compose up --build
```

**3. Verificar que el servidor está corriendo**

Deberías ver en la terminal:
```
api-json  | 2026/03/11 19:53:04 API running on port 24374
```

**4. Abrir en el navegador**
```
localhost:24374/api/players
```

### Comandos útiles

```bash
# Bajar el contenedor
docker compose down


# Reconstruir luego de cambios
docker compose down
docker compose up --build
```

---

## 📡 Endpoints disponibles

| Método   | Endpoint              | Descripción                    |
|----------|-----------------------|--------------------------------|
| GET      | `/api/players`        | Listar todos los jugadores     |
| GET      | `/api/players?filtros`| Listar jugadores con filtros   |
| GET      | `/api/players/{id}`   | Obtener un jugador por ID      |
| POST     | `/api/players`        | Crear un nuevo jugador         |
| PUT      | `/api/players/{id}`   | Actualizar un jugador por ID   |
| DELETE   | `/api/players/{id}`   | Eliminar un jugador por ID     |

---

## 📋 Ejemplos de uso

---

### 📄 Listar todos los jugadores — GET /api/players

**Request:**
```
GET localhost:24374/api/players
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Cristiano Ronaldo",
    "nationality": "Portugal",
    "position": "Forward",
    "current_team": "Al Nassr",
    "age": 41,
    "career_goals": 965,
    "active": true
  },
  {
    "id": 2,
    "name": "Lionel Messi",
    "nationality": "Argentina",
    "position": "Forward",
    "current_team": "Inter Miami",
    "age": 38,
    "career_goals": 899,
    "active": true
  }
  ...
]
```

---

### 🔍 Filtrar jugadores — GET /api/players?parametros

Se pueden combinar múltiples filtros en una sola petición.

#### Parámetros disponibles

| Parámetro     | Tipo   | Descripción                          |
|---------------|--------|--------------------------------------|
| `nationality` | string | Filtra por nacionalidad              |
| `position`    | string | Filtra por posición                  |
| `active`      | bool   | Filtra por jugadores activos o no    |
| `min_goals`   | int    | Goles mínimos en la carrera          |
| `max_goals`   | int    | Goles máximos en la carrera          |
| `search`      | string | Busca por nombre (contiene)          |

#### Ejemplos

**Jugadores activos:**
```
GET localhost:24374/api/players?active=true
```

**Jugadores de Brasil:**
```
GET localhost:24374/api/players?nationality=Brazil
```

**Jugadores con más de 700 goles:**
```
GET localhost:24374/api/players?min_goals=700
```

**Buscar por nombre:**
```
GET localhost:24374/api/players?search=messi
```

**Filtros combinados — activos con más de 700 goles:**
```
GET localhost:24374/api/players?active=true&min_goals=700
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Cristiano Ronaldo",
    "nationality": "Portugal",
    "position": "Forward",
    "current_team": "Al Nassr",
    "age": 41,
    "career_goals": 965,
    "active": true
  },
  {
    "id": 2,
    "name": "Lionel Messi",
    "nationality": "Argentina",
    "position": "Forward",
    "current_team": "Inter Miami",
    "age": 38,
    "career_goals": 899,
    "active": true
  }
]
```

---

### 🔎 Obtener jugador por ID — GET /api/players/{id}

**Request:**
```
GET localhost:24374/api/players/1
```

**Response:**
```json
{
  "id": 1,
  "name": "Cristiano Ronaldo",
  "nationality": "Portugal",
  "position": "Forward",
  "current_team": "Al Nassr",
  "age": 41,
  "career_goals": 965,
  "active": true
}
```

---

### ➕ Crear jugador — POST /api/players

**Request:**
```
POST localhost:24374/api/players
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Kylian Mbappé",
  "nationality": "France",
  "position": "Forward",
  "current_team": "Real Madrid",
  "age": 26,
  "career_goals": 370,
  "active": true
}
```

**Response `201 Created`:**
```json
{
  "id": 11,
  "name": "Kylian Mbappé",
  "nationality": "France",
  "position": "Forward",
  "current_team": "Real Madrid",
  "age": 26,
  "career_goals": 370,
  "active": true
}
```

---

## ❌ Manejo de errores

Todos los errores devuelven una respuesta JSON estructurada:

```json
{
  "code": 404,
  "error": "Not Found",
  "message": "Player not found"
}
```

| Código | Causa |
|--------|-------|
| `400`  | JSON inválido, campo requerido faltante, tipo de dato incorrecto |
| `404`  | Jugador no encontrado por ID |
| `405`  | Método HTTP no soportado en el endpoint |

---

## 💾 Persistencia

Cada POST, PUT y DELETE reescribe inmediatamente el archivo `data/players.json`, por lo que todos los cambios se mantienen aunque el servidor se reinicie.

El `docker-compose.yml` monta la carpeta `./data` como volumen, asegurando que los datos persistan incluso al reconstruir la imagen.    