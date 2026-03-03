'use client';

import { useRef } from 'react';
import { InvoiceList } from '@/services/modules/invoicing/ui/components/InvoiceList';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInvoicingBusiness } from '@/shared/contexts/invoicing-business-context';

export default function InvoicesPage() {
  const { permissions } = usePermissions();
  const { selectedBusinessId } = useInvoicingBusiness();
  const businessId = permissions?.business_id || 0;
  const invoiceListRef = useRef<any>(null);

  return (
    <div className="p-8">
      <div ref={invoiceListRef} />
      <div className="bg-white rounded-lg shadow p-6">
        <InvoiceList ref={invoiceListRef} businessId={businessId} selectedBusinessId={selectedBusinessId} />
      </div>
    </div>
  );
}
