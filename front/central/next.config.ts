import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone',

  // Headers para soportar iframes de Shopify
  async headers() {
    return [
      {
        // Aplicar a todas las rutas
        source: '/:path*',
        headers: [
          {
            key: 'Access-Control-Allow-Credentials',
            value: 'true',
          },
          {
            key: 'Access-Control-Allow-Methods',
            value: 'GET,POST,PUT,DELETE,OPTIONS',
          },
          {
            key: 'Access-Control-Allow-Headers',
            value: 'X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version, Authorization',
          },
        ],
      },
    ];
  },

  // Configuración experimental para mejorar el manejo de cookies en iframes
  experimental: {
    // Permitir que las cookies funcionen en contextos de terceros
    // @ts-ignore - Esta opción puede no estar en los tipos pero funciona
    serverActions: {
      bodySizeLimit: '2mb',
    },
  },
};

export default nextConfig;
