"use client";

import { useState, useEffect, useCallback } from "react";
import { OrdersApiRepository } from "../../infra/repository/api-repository";
import type { ReferenceData } from "../../domain/types";

interface Props {
  businessId: number;
}

export default function ReferenceDataPanel({ businessId }: Props) {
  const [data, setData] = useState<ReferenceData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<"products" | "integrations" | "payments" | "statuses">("products");

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const repo = new OrdersApiRepository();
      const result = await repo.getReferenceData(businessId);
      setData(result);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to load reference data");
    } finally {
      setLoading(false);
    }
  }, [businessId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const tabs = [
    { key: "products" as const, label: "Products", count: data?.products?.length || 0 },
    { key: "integrations" as const, label: "Integrations", count: data?.integrations?.length || 0 },
    { key: "payments" as const, label: "Payments", count: data?.payment_methods?.length || 0 },
    { key: "statuses" as const, label: "Statuses", count: data?.order_statuses?.length || 0 },
  ];

  return (
    <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
      <div className="px-4 py-3 border-b border-gray-200 flex items-center justify-between">
        <h2 className="font-semibold text-gray-900">Reference Data</h2>
        <button
          onClick={fetchData}
          disabled={loading}
          className="text-xs text-indigo-600 hover:text-indigo-800 font-medium"
        >
          {loading ? "Loading..." : "Refresh"}
        </button>
      </div>

      {error && (
        <div className="m-3 p-2 bg-red-50 border border-red-200 rounded text-sm text-red-700">
          {error}
        </div>
      )}

      <div className="flex border-b border-gray-200">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`flex-1 px-2 py-2 text-xs font-medium transition-colors ${
              activeTab === tab.key
                ? "text-indigo-600 border-b-2 border-indigo-600 bg-indigo-50"
                : "text-gray-500 hover:text-gray-700"
            }`}
          >
            {tab.label} ({tab.count})
          </button>
        ))}
      </div>

      <div className="p-3 max-h-96 overflow-y-auto">
        {loading && !data ? (
          <div className="flex justify-center py-8">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600" />
          </div>
        ) : !data ? (
          <p className="text-sm text-gray-400 text-center py-4">No data</p>
        ) : (
          <>
            {activeTab === "products" && (
              <table className="w-full text-sm">
                <thead>
                  <tr className="text-left text-gray-500 text-xs">
                    <th className="pb-2">Name</th>
                    <th className="pb-2">SKU</th>
                    <th className="pb-2 text-right">Price</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {data.products?.map((p) => (
                    <tr key={p.id}>
                      <td className="py-1.5 text-gray-900">{p.name}</td>
                      <td className="py-1.5 text-gray-500 font-mono text-xs">{p.sku}</td>
                      <td className="py-1.5 text-right text-gray-900">
                        {p.price > 0 ? `$${p.price.toLocaleString()}` : "-"}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}

            {activeTab === "integrations" && (
              <table className="w-full text-sm">
                <thead>
                  <tr className="text-left text-gray-500 text-xs">
                    <th className="pb-2">ID</th>
                    <th className="pb-2">Name</th>
                    <th className="pb-2">Category</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {data.integrations?.map((i) => (
                    <tr key={i.id}>
                      <td className="py-1.5 text-gray-500">{i.id}</td>
                      <td className="py-1.5 text-gray-900">{i.name}</td>
                      <td className="py-1.5">
                        <span className="px-2 py-0.5 bg-gray-100 rounded text-xs text-gray-600">
                          {i.category || "other"}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}

            {activeTab === "payments" && (
              <table className="w-full text-sm">
                <thead>
                  <tr className="text-left text-gray-500 text-xs">
                    <th className="pb-2">ID</th>
                    <th className="pb-2">Code</th>
                    <th className="pb-2">Name</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {data.payment_methods?.map((pm) => (
                    <tr key={pm.id}>
                      <td className="py-1.5 text-gray-500">{pm.id}</td>
                      <td className="py-1.5 font-mono text-xs text-gray-500">{pm.code}</td>
                      <td className="py-1.5 text-gray-900">{pm.name}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}

            {activeTab === "statuses" && (
              <table className="w-full text-sm">
                <thead>
                  <tr className="text-left text-gray-500 text-xs">
                    <th className="pb-2">ID</th>
                    <th className="pb-2">Code</th>
                    <th className="pb-2">Name</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {data.order_statuses?.map((s) => (
                    <tr key={s.id}>
                      <td className="py-1.5 text-gray-500">{s.id}</td>
                      <td className="py-1.5 font-mono text-xs text-gray-500">{s.code}</td>
                      <td className="py-1.5 text-gray-900">{s.name}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </>
        )}
      </div>
    </div>
  );
}
