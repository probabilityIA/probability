import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone',

  // Proxy para desarrollo local - evita problemas de CORS
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:3050/api/:path*',
      },
    ];
  },

  // Headers para CORS y Shopify iframes
  async headers() {
    return [
      {
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

  // Configuración experimental para mejorar el manejo de cookies
  experimental: {
    // @ts-ignore
    serverActions: {
      bodySizeLimit: '2mb',
    },
  },
};

export default nextConfig;
