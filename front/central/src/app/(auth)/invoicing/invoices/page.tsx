/**
 * Página de Listado de Facturas
 * Muestra todas las facturas generadas del sistema
 */

'use client';

import { InvoiceList } from '@/services/modules/invoicing/ui/components/InvoiceList';
import { usePermissions } from '@/shared/contexts/permissions-context';

export default function InvoicesPage() {
  const { permissions } = usePermissions();
  const businessId = permissions?.business_id || 0;

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Facturas</h1>
        <p className="text-gray-600 mt-2">
          Gestiona las facturas electrónicas generadas automáticamente
        </p>
      </div>

      <div className="bg-white rounded-lg shadow p-6">
        <InvoiceList businessId={businessId} />
      </div>
    </div>
  );
}
