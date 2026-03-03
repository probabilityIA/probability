"use client";

import { useState } from "react";
import AppLayout from "@/shared/components/AppLayout";
import BusinessSelector from "@/shared/components/BusinessSelector";
import ReferenceDataPanel from "@/services/orders/ui/components/ReferenceDataPanel";
import OrderGenerator from "@/services/orders/ui/components/OrderGenerator";
import APIConsole from "@/services/orders/ui/components/APIConsole";
import type { APICallLog } from "@/services/orders/domain/types";

export default function OrdersPage() {
  const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
  const [apiLogs, setApiLogs] = useState<APICallLog[]>([]);

  return (
    <AppLayout>
      <div className="space-y-4">
        <h1 className="text-xl font-bold text-gray-900">Order Generator</h1>

        <BusinessSelector
          value={selectedBusinessId}
          onChange={setSelectedBusinessId}
        />

        {selectedBusinessId === null ? (
          <div className="text-center py-16 text-gray-400 bg-white rounded-lg border border-gray-200">
            Select a business to continue
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <ReferenceDataPanel businessId={selectedBusinessId} />
              <OrderGenerator businessId={selectedBusinessId} onApiLogs={setApiLogs} />
            </div>

            {apiLogs.length > 0 && (
              <APIConsole logs={apiLogs} />
            )}
          </>
        )}
      </div>
    </AppLayout>
  );
}
