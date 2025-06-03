# 📊 Plataforma de Recomendaciones Bursátiles

Esta es una plataforma web desarrollada como parte de un reto técnico. Permite consultar y visualizar información bursátil con un enfoque en recomendaciones y análisis técnico. Su arquitectura está diseñada para escalabilidad, mantenibilidad y claridad, dividida entre backend y frontend desacoplados.

## 🚀 Funcionalidades

- Inicio de sesión simulado.
- Consulta de acciones con filtros por símbolo y calificación.
- Recomendaciones generadas a partir de análisis técnico (como RSI).
- Almacenamiento persistente en base de datos distribuida.
- API RESTful para exponer los datos a la interfaz.
- Interfaz moderna con carga dinámica y manejo de errores.
- Visualización de detalles y cambio porcentual por acción.
- Carga progresiva (paginación).
- Pruebas unitarias para backend y frontend.
- Infraestructura como código (opcional con Terraform).

---

## ⚙️ Instalación

### 1. Clonar el repositorio

```bash
git clone https://github.com/usuario/proyecto.git
cd proyecto
```

### 2. Configuración del Backend

```bash
cd backend
go run main.go
```

Variables necesarias en `.env`:

```
API_KEY=tu_clave
DB_URL=tu_conexion
PORT=8080
```

### 3. Configuración del Frontend

```bash
cd frontend
npm install
npm run dev
```

### 4. Despliegue con Terraform (opcional)

```bash
cd infrastructure
terraform init
terraform apply
```

---

## 🧪 Pruebas

```bash
# Backend
cd backend
go test ./...

# Frontend
cd frontend
npm run test
```

---

## 📁 Estructura del Proyecto

```plaintext
backend/
├── api/              # Handlers, cliente externo, enrutador, lógica de recomendaciones
├── config/           # Configuración y carga de variables de entorno
├── db/               # Conexión y acceso a la base de datos
├── docs/             # Documentación Swagger
├── models/           # Definiciones de estructuras de datos
├── pkg/              # Logging y métricas
├── test/             # Pruebas unitarias
└── main.go           # Entrada principal del servidor

frontend/
├── public/           
├── src/
│   ├── api/          # Funciones para comunicarse con el backend
│   ├── assets/       # Estilos y recursos estáticos
│   ├── components/   # Componentes reutilizables (cards, filtros, headers)
│   ├── router/       # Configuración de rutas
│   ├── stores/       # Gestión de estado global con pinia
│   ├── types/        # Tipos TypeScript
│   ├── views/        # Vistas de cada página
│   └── main.ts       # Punto de entrada de la app
└── vite.config.ts    # Configuración de Vite
```

---

## 🧠 Justificación de la Arquitectura

La estructura fue elegida para separar claramente responsabilidades:

- **Modularidad**: Separar cada dominio (API, DB, UI) permite mayor claridad, pruebas independientes y fácil escalabilidad.
- **Escalabilidad**: Se pueden añadir nuevas fuentes de datos o componentes UI sin afectar el resto del sistema.
- **Mantenibilidad**: Una estructura clara y predecible acelera el desarrollo colaborativo y la depuración.
- **Performance**: Go permite alto rendimiento y concurrencia eficiente para el backend. Vue 3 con Vite ofrece recarga instantánea y excelente rendimiento en el navegador.
- **Responsividad**: La UI está diseñada con Tailwind CSS para adaptarse bien a distintos dispositivos.
- **Infraestructura**: Terraform permite desplegar la infraestructura reproduciblemente y mantener control sobre entornos.

---

## 🧠 Lógica de Recomendación

Las recomendaciones se generan con base en:

- Calificaciones recientes (ej. Buy, Hold, Strong Sell).
- Indicadores técnicos como el RSI.
- Datos actualizados periódicamente desde el backend.
- Filtrado dinámico por el usuario.

---

## 📜 Licencia

Este proyecto fue desarrollado con fines educativos y como solución a un reto técnico. No representa asesoría financiera ni está pensado para uso productivo real.
