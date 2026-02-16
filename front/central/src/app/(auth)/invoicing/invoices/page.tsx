/**
 * Página de Listado de Facturas
 * Muestra todas las facturas generadas del sistema
 */

'use client';

import { useRef } from 'react';
import { InvoiceList } from '@/services/modules/invoicing/ui/components/InvoiceList';
import { InvoicingHeader } from '@/services/modules/invoicing/ui/components/InvoicingHeader';
import { usePermissions } from '@/shared/contexts/permissions-context';

export default function InvoicesPage() {
  const { permissions } = usePermissions();
  const businessId = permissions?.business_id || 0;
  const invoiceListRef = useRef<any>(null);

  return (
    <div className="p-8">
      <InvoicingHeader
        title="Facturas"
        description="Gestiona las facturas electrónicas generadas automáticamente"
      />

      {/* Pasar ref del botón al InvoiceList */}
      <div ref={invoiceListRef} />

      <div className="bg-white rounded-lg shadow p-6">
        <InvoiceList ref={invoiceListRef} businessId={businessId} />
      </div>
    </div>
  );
}
