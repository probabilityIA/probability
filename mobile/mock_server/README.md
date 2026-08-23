# Mock API - app mobile

Servidor HTTP sin dependencias que imita el contrato de `back/central` para
desarrollar la app Flutter sin levantar Go, Postgres, Redis ni RabbitMQ.

```bash
node server.js            # escucha en :5199, prefijo /api/v1
MOCK_PORT=6000 node server.js
```

- `server.js` - HTTP, CORS, paginacion y logging de rutas no cubiertas.
- `fixtures.js` - tabla de rutas (`add(metodo, patron, handler)`).
- `data.js` - datos deterministas (ordenes, productos, clientes, guias, etc).

Cuando la app pida una ruta que no existe el servidor responde 404 con el
mensaje `mock sin ruta: METODO /path` y lo escribe en consola: esa es la lista
de trabajo para completar el modulo en curso.
