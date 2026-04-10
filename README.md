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
   git clone https://github.com/RGonzalez-23/Proyecto_Capacitacion_CITIAPS_BACKEND
   ```
2. Navega a la carpeta `BACKEND`:
   ```powershell
   cd BACKEND
   ```

3. Descarga las dependencias:
   ```powershell
   go mod download
   ```

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
