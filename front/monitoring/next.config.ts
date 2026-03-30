import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone',

  async rewrites() {
    return [
      {
        source: '/api/monitoring/:path*',
        destination: `${process.env.MONITORING_API_URL || 'http://localhost:3070'}/api/v1/:path*`,
      },
    ];
  },
};

export default nextConfig;
