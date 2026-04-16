/**
 * HTTP Logger para servidor
 * Imprime peticiones HTTP en la consola con colores
 */

import { inspect } from 'util';

interface LogRequestOptions {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  url: string;
  token?: string;
  body?: unknown;
}

interface LogResponseOptions {
  status: number;
  statusText: string;
  duration: number;
  data?: unknown;
  summary?: string;
}

/**
 * Logea una petición HTTP
 */
export function logHttpRequest(options: LogRequestOptions): void {
  const { method, url, token, body } = options;
  
  console.log('\n🌐 HTTP', method, url);
  
  if (token) {
    console.log('🔑 Token:', token.substring(0, 20) + '...');
  }
  
  if (body) {
    console.log('📤 Body:');
    console.log(inspect(body, { colors: true, depth: null }));
  }
}

/**
 * Logea una respuesta HTTP exitosa
 */
export function logHttpSuccess(options: LogResponseOptions): void {
  const { status, statusText, duration, data, summary } = options;
  
  console.log(`✅ ${status} ${statusText} (${duration}ms)`);
  
  if (summary) {
    console.log('📊', summary);
  }
  
  if (data) {
    console.log('📥 Response:');
    console.log(inspect(data, { colors: true, depth: null }));
  }
  
  console.log('');
}

/**
 * Logea una respuesta HTTP con error
 */
export function logHttpError(options: LogResponseOptions): void {
  const { status, statusText, duration, data } = options;
  
  console.log(`❌ ${status} ${statusText} (${duration}ms)`);
  
  if (data) {
    console.log('📥 Error:');
    console.log(inspect(data, { colors: true, depth: null }));
  }
  
  console.log('');
}

