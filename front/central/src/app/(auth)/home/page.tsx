/**
 * Página principal después del login
 */

import { Dashboard } from '@/services/modules/dashboard/ui';

export default function HomePage() {
  return (
    <div className="p-8">
      <Dashboard />
    </div>
  );
}
