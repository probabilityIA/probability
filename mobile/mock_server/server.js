const http = require('http');
const { URL } = require('url');
const fixtures = require('./fixtures');

const PORT = process.env.MOCK_PORT ? Number(process.env.MOCK_PORT) : 5199;
const PREFIX = '/api/v1';

function send(res, status, payload) {
  const body = JSON.stringify(payload);
  res.writeHead(status, {
    'Content-Type': 'application/json; charset=utf-8',
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Headers': '*',
    'Access-Control-Allow-Methods': 'GET,POST,PUT,PATCH,DELETE,OPTIONS',
    'Content-Length': Buffer.byteLength(body),
  });
  res.end(body);
}

function paginate(items, url) {
  const page = Number(url.searchParams.get('page') || 1);
  const pageSize = Number(url.searchParams.get('page_size') || url.searchParams.get('per_page') || 10);
  const search = (url.searchParams.get('search') || '').toLowerCase();
  let rows = items;
  if (search) {
    rows = rows.filter((r) => JSON.stringify(r).toLowerCase().includes(search));
  }
  const total = rows.length;
  const lastPage = Math.max(1, Math.ceil(total / pageSize));
  const slice = rows.slice((page - 1) * pageSize, page * pageSize);
  return {
    success: true,
    data: slice,
    total,
    page,
    page_size: pageSize,
    total_pages: lastPage,
    pagination: {
      current_page: page,
      per_page: pageSize,
      total,
      last_page: lastPage,
      has_next: page < lastPage,
      has_prev: page > 1,
    },
  };
}

function readBody(req) {
  return new Promise((resolve) => {
    let raw = '';
    req.on('data', (c) => (raw += c));
    req.on('end', () => {
      try {
        resolve(raw ? JSON.parse(raw) : {});
      } catch (_) {
        resolve({});
      }
    });
  });
}

const server = http.createServer(async (req, res) => {
  if (req.method === 'OPTIONS') return send(res, 204, {});

  const url = new URL(req.url, `http://localhost:${PORT}`);
  const path = url.pathname.startsWith(PREFIX)
    ? url.pathname.slice(PREFIX.length)
    : url.pathname;
  const body = ['POST', 'PUT', 'PATCH'].includes(req.method) ? await readBody(req) : {};

  const handler = fixtures.resolve(req.method, path);
  if (!handler) {
    console.log(`[mock] 404 ${req.method} ${path}`);
    return send(res, 404, { success: false, error: `mock sin ruta: ${req.method} ${path}` });
  }

  console.log(`[mock] 200 ${req.method} ${path}`);
  const result = handler({ url, body, params: handler.params || {}, paginate });
  return send(res, result.status || 200, result.payload);
});

server.listen(PORT, () => {
  console.log(`[mock] Probability mock API en http://localhost:${PORT}${PREFIX}`);
});
