# CITIAPS Backend

API REST para la gestión de tareas con etiquetas, desarrollada en Go.

## Descripción

Backend de la aplicación CITIAPS que proporciona endpoints para:
- Crear, leer, actualizar y eliminar tareas
- Gestionar etiquetas con colores personalizados
- Validación de datos y protección contra inyecciones NoSQL

## Estructura del Proyecto

```
BACKEND/
├── main.go              # Punto de entrada y configuración de rutas
├── go.mod              # Dependencias del proyecto
├── LICENSE             # Licencia
├── README.md           # Este archivo
│
├── controller/         # Controladores de lógica de negocio
│   ├── taskController.go     # Manejo de tareas
│   └── tagController.go      # Manejo de etiquetas
│
├── middleware/         # Middleware de la aplicación
│   └── cors.go         # Configuración CORS
│
├── model/              # Modelos de datos
│   └── task.go         # Estructura de tareas y etiquetas
│
└── util/               # Funciones utilitarias
    └── database.go     # Conexión a MongoDB
```

## Requisitos

- **Go** 1.18 o superior
- **MongoDB** (local o remoto)
- **Gorilla Mux** (dependencia automática)

## Instalación

1. Clona el proyecto o descárgalo
   ```powershell
   git clone [url-repositorio]
   ```

2. Descarga las dependencias en la carpeta donde hayas clonado el repositorio:
   ```powershell
   go mod download
   ```

## Variables de Entorno

**⚠️ IMPORTANTE:** Este proyecto requiere configurar variables de entorno para conectarse a MongoDB y otras funciones.

Crea un archivo `.env` en la raíz del proyecto `BACKEND/` con las siguientes variables:

```env
# Configuración de MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB_NAME=citiaps

# Configuración del Servidor
SERVER_PORT=8080
SERVER_HOST=localhost

# Configuración CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8000
```

### Descripción de Variables

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `MONGODB_URI` | URL de conexión a MongoDB | `mongodb://localhost:27017` |
| `MONGODB_DB_NAME` | Nombre de la base de datos | `citiaps` |
| `SERVER_PORT` | Puerto en el que corre el servidor | `8080` |
| `SERVER_HOST` | Host del servidor | `localhost` o `0.0.0.0` |
| `CORS_ALLOWED_ORIGINS` | Orígenes permitidos (separados por comas) | `http://localhost:3000` |

**Nota:** Si no proporcionas un archivo `.env`, el servidor usará valores por defecto.

## Configuración de Base de Datos

**Prerequisito:** MongoDB debe estar instalado y ejecutándose en tu equipo.

### Crear la Base de Datos

1. Abre una terminal/PowerShell
2. Conecta a MongoDB:
   ```powershell
   mongosh mongodb://localhost:27017
   ```

3. Crea la base de datos y colecciones:
   ```javascript
   // Cambiar a la base de datos Tasks CITIAPS
   use tasks_citiaps

   // Crear colección de tareas
   db.createCollection("tasks")

   // Crear colección de etiquetas
   db.createCollection("tags")

   // Verificar que se crearon correctamente
   show collections
   ```

4. (Opcional) Ver la base de datos creada:
   ```javascript
   show databases
   ```

**Listo:** La base de datos `tasks_citiaps` está lista para usar. Las colecciones se crearán automáticamente con los primeros datos cuando ejecutes la aplicación.

## Ejecución

Ejecuta el servidor:

```powershell
go run main.go
```

El servidor estará disponible en: **http://localhost:8080**

## API Endpoints

### Tareas
- `GET /api/tasks` - Obtener todas las tareas
- `POST /api/tasks` - Crear nueva tarea
- `PUT /api/tasks/:id` - Actualizar tarea (marcar como completada)
- `DELETE /api/tasks/:id` - Eliminar tarea

### Etiquetas
- `GET /api/tags` - Obtener todas las etiquetas
- `POST /api/tags` - Crear nueva etiqueta
- `DELETE /api/tags/:id` - Eliminar etiqueta

## Características Principales

✅ CRUD completo de tareas y etiquetas  
✅ Validación de datos en servidor  
✅ Protección contra inyecciones NoSQL  
✅ Etiquetas con colores personalizados  
✅ CORS habilitado para frontend  
✅ Integración con MongoDB  

---

**Desarrollado para capacitación de CITIAPS**
