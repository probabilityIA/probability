"use client";

import { useState } from "react";
import AppLayout from "@/shared/components/AppLayout";
import BusinessSelector from "@/shared/components/BusinessSelector";
import ReferenceDataPanel from "@/services/orders/ui/components/ReferenceDataPanel";
import OrderGenerator from "@/services/orders/ui/components/OrderGenerator";
import APIConsole from "@/services/orders/ui/components/APIConsole";
import { OrdersApiRepository } from "@/services/orders/infra/repository/api-repository";
import type { APICallLog } from "@/services/orders/domain/types";

export default function OrdersPage() {
  const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
  const [apiLogs, setApiLogs] = useState<APICallLog[]>([]);
  const [showDeleteAllModal, setShowDeleteAllModal] = useState(false);
  const [deleteStep, setDeleteStep] = useState<1 | 2 | 3>(1);
  const [deleteConfirmText, setDeleteConfirmText] = useState("");
  const [isDeletingAll, setIsDeletingAll] = useState(false);
  const [deleteResult, setDeleteResult] = useState<string | null>(null);

  const handleDeleteAllOrders = async () => {
    if (!selectedBusinessId) return;
    setIsDeletingAll(true);
    try {
      const repo = new OrdersApiRepository();
      const result = await repo.deleteAllOrders(selectedBusinessId);
      setDeleteResult(`${result.deleted} orders deleted successfully`);
    } catch (error: any) {
      setDeleteResult(`Error: ${error.message || "Failed to delete orders"}`);
    } finally {
      setIsDeletingAll(false);
      setShowDeleteAllModal(false);
      setDeleteStep(1);
      setDeleteConfirmText("");
    }
  };

  return (
    <AppLayout>
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h1 className="text-xl font-bold text-gray-900">Order Generator</h1>
          {selectedBusinessId !== null && (
            <button
              onClick={() => { setDeleteStep(1); setDeleteConfirmText(""); setDeleteResult(null); setShowDeleteAllModal(true); }}
              className="px-4 py-2 text-sm font-semibold text-white bg-red-600 rounded-lg hover:bg-red-700 hover:shadow-lg transition-all"
            >
              Delete All Orders
            </button>
          )}
        </div>

        <BusinessSelector
          value={selectedBusinessId}
          onChange={setSelectedBusinessId}
        />

        {deleteResult && (
          <div className={`p-3 rounded-lg text-sm ${deleteResult.startsWith("Error") ? "bg-red-50 text-red-700 border border-red-200" : "bg-green-50 text-green-700 border border-green-200"}`}>
            {deleteResult}
            <button onClick={() => setDeleteResult(null)} className="ml-2 underline">dismiss</button>
          </div>
        )}

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

      {/* Delete All Orders Modal - 3-step confirmation */}
      {showDeleteAllModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="bg-white rounded-xl shadow-2xl w-full max-w-md mx-4 p-6">
            {deleteStep === 1 && (
              <>
                <div className="flex items-center gap-3 mb-4">
                  <span className="text-3xl">Warning</span>
                  <h2 className="text-xl font-bold text-gray-900">Delete all orders</h2>
                </div>
                <p className="text-gray-600 mb-6">
                  This will <strong>permanently</strong> delete all orders for the selected business, including invoices, payments, shipments, and all related data. This action <strong>cannot be undone</strong>.
                </p>
                <div className="flex gap-3 justify-end">
                  <button
                    onClick={() => setShowDeleteAllModal(false)}
                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={() => setDeleteStep(2)}
                    className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors"
                  >
                    Continue
                  </button>
                </div>
              </>
            )}

            {deleteStep === 2 && (
              <>
                <div className="flex items-center gap-3 mb-4">
                  <span className="text-3xl">Alert</span>
                  <h2 className="text-xl font-bold text-gray-900">Are you sure?</h2>
                </div>
                <p className="text-gray-600 mb-6">
                  <strong>All</strong> orders for this business will be deleted. This is irreversible and will also affect associated invoices, payments, and shipments.
                </p>
                <div className="flex gap-3 justify-end">
                  <button
                    onClick={() => setShowDeleteAllModal(false)}
                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={() => setDeleteStep(3)}
                    className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors"
                  >
                    Yes, continue
                  </button>
                </div>
              </>
            )}

            {deleteStep === 3 && (
              <>
                <div className="flex items-center gap-3 mb-4">
                  <span className="text-3xl text-red-600 font-bold">!</span>
                  <h2 className="text-xl font-bold text-gray-900">Final confirmation</h2>
                </div>
                <p className="text-gray-600 mb-4">
                  Type <strong className="text-red-600">DELETE</strong> to confirm:
                </p>
                <input
                  type="text"
                  value={deleteConfirmText}
                  onChange={(e) => setDeleteConfirmText(e.target.value)}
                  placeholder="DELETE"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm mb-6 focus:outline-none focus:ring-2 focus:ring-red-500"
                  autoFocus
                />
                <div className="flex gap-3 justify-end">
                  <button
                    onClick={() => setShowDeleteAllModal(false)}
                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
                    disabled={isDeletingAll}
                  >
                    Cancel
                  </button>
                  <button
                    onClick={handleDeleteAllOrders}
                    disabled={deleteConfirmText !== "DELETE" || isDeletingAll}
                    className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {isDeletingAll ? "Deleting..." : "Delete all"}
                  </button>
                </div>
              </>
            )}
          </div>
        </div>
      )}
    </AppLayout>
  );
}
