"use client";

import { useState, useEffect, useCallback } from "react";
import { OrdersApiRepository } from "../../infra/repository/api-repository";
import type { Integration, GenerateResult, APICallLog } from "../../domain/types";

interface Props {
  businessId: number;
  onApiLogs?: (logs: APICallLog[]) => void;
}

export default function OrderGenerator({ businessId, onApiLogs }: Props) {
  const [integrations, setIntegrations] = useState<Integration[]>([]);
  const [integrationId, setIntegrationId] = useState<number | "">("");
  const [count, setCount] = useState(3);
  const [maxItems, setMaxItems] = useState(3);
  const [randomProducts, setRandomProducts] = useState(true);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<GenerateResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  const loadIntegrations = useCallback(async () => {
    try {
      const repo = new OrdersApiRepository();
      const data = await repo.getReferenceData(businessId);
      setIntegrations(data.integrations || []);
      // Auto-select platform integration
      const platform = data.integrations?.find((i) => i.category === "platform");
      if (platform) setIntegrationId(platform.id);
    } catch {
      // ignore — reference data panel handles errors
    }
  }, [businessId]);

  useEffect(() => {
    loadIntegrations();
    setResult(null);
    setError(null);
  }, [loadIntegrations]);

  const handleGenerate = async () => {
    setLoading(true);
    setError(null);
    setResult(null);
    onApiLogs?.([]);

    try {
      const repo = new OrdersApiRepository();
      const res = await repo.generateOrders(businessId, {
        count,
        integration_id: integrationId ? Number(integrationId) : undefined,
        random_products: randomProducts,
        max_items_per_order: maxItems,
      });
      setResult(res);
      if (res.api_logs) onApiLogs?.(res.api_logs);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to generate orders");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
      <div className="px-4 py-3 border-b border-gray-200">
        <h2 className="font-semibold text-gray-900">Generate Orders</h2>
      </div>

      <div className="p-4 space-y-4">
        {error && (
          <div className="p-2 bg-red-50 border border-red-200 rounded text-sm text-red-700">
            {error}
          </div>
        )}

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Integration
          </label>
          <select
            value={integrationId}
            onChange={(e) => setIntegrationId(e.target.value ? Number(e.target.value) : "")}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900"
          >
            <option value="">Auto-select (platform)</option>
            {integrations.map((i) => (
              <option key={i.id} value={i.id}>
                {i.name} (ID: {i.id})
              </option>
            ))}
          </select>
        </div>

        <div className="grid grid-cols-2 gap-3">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Order Count (1-20)
            </label>
            <input
              type="number"
              min={1}
              max={20}
              value={count}
              onChange={(e) => setCount(Math.min(20, Math.max(1, Number(e.target.value))))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Max Items/Order (1-5)
            </label>
            <input
              type="number"
              min={1}
              max={5}
              value={maxItems}
              onChange={(e) => setMaxItems(Math.min(5, Math.max(1, Number(e.target.value))))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900"
            />
          </div>
        </div>

        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={randomProducts}
            onChange={(e) => setRandomProducts(e.target.checked)}
            className="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
          />
          <span className="text-sm text-gray-700">Random product selection</span>
        </label>

        <button
          onClick={handleGenerate}
          disabled={loading}
          className="w-full py-2.5 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {loading ? "Generating..." : `Generate ${count} Order${count > 1 ? "s" : ""}`}
        </button>

        {result && (
          <div className="mt-4 space-y-3">
            <div className="flex gap-3">
              <div className="flex-1 p-3 bg-green-50 border border-green-200 rounded-lg text-center">
                <div className="text-2xl font-bold text-green-700">{result.created}</div>
                <div className="text-xs text-green-600">Created</div>
              </div>
              <div className="flex-1 p-3 bg-red-50 border border-red-200 rounded-lg text-center">
                <div className="text-2xl font-bold text-red-700">{result.failed}</div>
                <div className="text-xs text-red-600">Failed</div>
              </div>
              <div className="flex-1 p-3 bg-gray-50 border border-gray-200 rounded-lg text-center">
                <div className="text-2xl font-bold text-gray-700">{result.total}</div>
                <div className="text-xs text-gray-600">Total</div>
              </div>
            </div>

            {result.orders && result.orders.length > 0 && (
              <div>
                <h3 className="text-sm font-medium text-gray-700 mb-2">Created Orders</h3>
                <div className="max-h-48 overflow-y-auto border border-gray-200 rounded-lg divide-y divide-gray-100">
                  {result.orders.map((order, idx) => (
                    <div key={idx} className="px-3 py-2 text-sm flex justify-between items-center">
                      <div>
                        <span className="font-mono text-xs text-indigo-600">{order.order_number}</span>
                        <span className="mx-2 text-gray-300">|</span>
                        <span className="text-gray-600">{order.customer_name}</span>
                      </div>
                      <span className="text-gray-900 font-medium">
                        ${order.total.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {result.errors && result.errors.length > 0 && (
              <div>
                <h3 className="text-sm font-medium text-red-700 mb-2">Errors</h3>
                <div className="max-h-32 overflow-y-auto border border-red-200 rounded-lg divide-y divide-red-100">
                  {result.errors.map((err, idx) => (
                    <div key={idx} className="px-3 py-2 text-xs text-red-600">
                      #{err.index + 1}: {err.message}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}

      </div>
    </div>
  );
}
