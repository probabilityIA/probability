/**
 * Página principal después del login
 */

'use client';

import dynamic from 'next/dynamic';
import { Spinner } from '@/shared/ui';

// Importar Dashboard sin SSR (solo cliente)
// Esto previene errores de Server Components en iframe de Shopify
const Dashboard = dynamic(
  () => import('@/services/modules/dashboard/ui').then(mod => ({ default: mod.Dashboard })),
  {
    ssr: false,
    loading: () => (
      <div className="min-h-screen flex items-center justify-center">
        <Spinner size="xl" color="primary" text="Cargando dashboard..." />
      </div>
    ),
  }
);

export default function HomePage() {
  return (
    <div className="p-8">
      <Dashboard />
    </div>
  );
}
