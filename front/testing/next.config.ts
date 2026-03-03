import type { NextConfig } from "next";

const centralApiUrl = process.env.CENTRAL_API_URL || 'http://localhost:3050';
const testingApiUrl = process.env.TESTING_API_URL || 'http://localhost:9092';

const nextConfig: NextConfig = {
  output: 'standalone',
  basePath: process.env.BASE_PATH || '',

  async rewrites() {
    return [
      {
        source: '/api/central/:path*',
        destination: `${centralApiUrl}/api/:path*`,
      },
      {
        source: '/api/testing/:path*',
        destination: `${testingApiUrl}/api/:path*`,
      },
    ];
  },

  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          { key: 'Access-Control-Allow-Credentials', value: 'true' },
          { key: 'Access-Control-Allow-Methods', value: 'GET,POST,PUT,DELETE,OPTIONS' },
          { key: 'Access-Control-Allow-Headers', value: 'Content-Type, Authorization' },
        ],
      },
    ];
  },
};

export default nextConfig;
