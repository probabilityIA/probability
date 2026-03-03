"use client";

import { useState, useEffect, useCallback } from "react";
import { getToken } from "@/shared/lib/auth";
import { fetchReferenceDataAction, buildPayloadsAction, sendPayloadAction } from "@/shared/lib/server-actions";
import type { Integration, WebhookPayload, APICallLog } from "../../domain/types";

interface Props {
  businessId: number;
  onApiLogs?: (logs: APICallLog[]) => void;
}

const CATEGORY_ECOMMERCE = 1;
const CATEGORY_PLATFORM = 6;

export default function OrderGenerator({ businessId, onApiLogs }: Props) {
  const [integrations, setIntegrations] = useState<Integration[]>([]);
  const [webhookTopics, setWebhookTopics] = useState<Record<string, string[]>>({});
  const [integrationId, setIntegrationId] = useState<number | "">("");
  const [count, setCount] = useState(1);
  const [topic, setTopic] = useState("");
  const [maxItems, setMaxItems] = useState(3);
  const [loading, setLoading] = useState(false);
  const [sending, setSending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Built payloads from backend
  const [payloads, setPayloads] = useState<WebhookPayload[]>([]);
  // API logs from sending payloads
  const [logs, setLogs] = useState<APICallLog[]>([]);
  // Expanded log index
  const [expandedLog, setExpandedLog] = useState<number | null>(null);

  const selectedIntegration = integrations.find((i) => i.id === integrationId);
  const isEcommerce = selectedIntegration?.category_id === CATEGORY_ECOMMERCE;
  const isPlatform = selectedIntegration?.category_id === CATEGORY_PLATFORM;

  // Get topics for the selected integration type
  const availableTopics = selectedIntegration
    ? webhookTopics[selectedIntegration.integration_type_code] || []
    : [];

  const loadIntegrations = useCallback(async () => {
    const token = getToken();
    if (!token) return;

    try {
      const data = await fetchReferenceDataAction(token, businessId);
      const allIntegrations: Integration[] = data.integrations || [];
      setIntegrations(allIntegrations);
      setWebhookTopics(data.webhook_topics || {});
      if (allIntegrations.length > 0) {
        setIntegrationId(allIntegrations[0].id);
      }
    } catch {
      // ignore — reference data panel handles errors
    }
  }, [businessId]);

  useEffect(() => {
    loadIntegrations();
    setPayloads([]);
    setLogs([]);
    setError(null);
  }, [loadIntegrations]);

  // Auto-select first topic when integration changes
  useEffect(() => {
    if (isEcommerce && availableTopics.length > 0) {
      setTopic(availableTopics[0]);
    } else {
      setTopic("");
    }
  }, [integrationId, isEcommerce, availableTopics.length]);

  const handleBuildAndSend = async () => {
    const token = getToken();
    if (!token) return;

    setLoading(true);
    setError(null);
    setPayloads([]);
    setLogs([]);
    onApiLogs?.([]);

    try {
      // Step 1: Build payloads on backend
      const result = await buildPayloadsAction(token, businessId, {
        count,
        integration_id: integrationId ? Number(integrationId) : undefined,
        random_products: true,
        max_items_per_order: maxItems,
        topic: isEcommerce ? topic : "",
      });

      const builtPayloads: WebhookPayload[] = result.payloads || [];
      setPayloads(builtPayloads);

      if (builtPayloads.length === 0) {
        setError("No payloads were built");
        setLoading(false);
        return;
      }

      // Step 2: Send each payload
      setSending(true);
      const newLogs: APICallLog[] = [];

      for (let i = 0; i < builtPayloads.length; i++) {
        const payload = builtPayloads[i];
        const timestamp = new Date().toISOString();

        const response = await sendPayloadAction(payload);

        const log: APICallLog = {
          index: i,
          success: response.status_code >= 200 && response.status_code < 300,
          timestamp,
          duration_ms: response.duration_ms,
          request: {
            method: payload.method,
            url: payload.url,
            body: payload.body,
          },
          response: {
            status_code: response.status_code,
            body: response.body,
          },
        };

        newLogs.push(log);
        setLogs([...newLogs]);
        onApiLogs?.([...newLogs]);
      }
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to build payloads");
    } finally {
      setLoading(false);
      setSending(false);
    }
  };

  const successCount = logs.filter((l) => l.success).length;
  const failedCount = logs.filter((l) => !l.success).length;

  return (
    <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
      <div className="px-4 py-3 border-b border-gray-200">
        <h2 className="font-semibold text-gray-900">
          {isEcommerce ? "Simulate Ecommerce Webhooks" : "Generate Orders"}
        </h2>
      </div>

      <div className="p-4 space-y-4">
        {error && (
          <div className="p-2 bg-red-50 border border-red-200 rounded text-sm text-red-700">
            {error}
          </div>
        )}

        {/* Integration selector */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Integration
          </label>
          <select
            value={integrationId}
            onChange={(e) => setIntegrationId(e.target.value ? Number(e.target.value) : "")}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900"
          >
            <option value="">— Select integration —</option>
            {integrations.map((i) => (
              <option key={i.id} value={i.id}>
                {i.name} ({i.integration_type_code}) (ID: {i.id})
              </option>
            ))}
          </select>
        </div>

        {/* Topic selector — only for ecommerce */}
        {isEcommerce && availableTopics.length > 0 && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Webhook Topic
            </label>
            <select
              value={topic}
              onChange={(e) => setTopic(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900"
            >
              {availableTopics.map((t) => (
                <option key={t} value={t}>
                  {t}
                </option>
              ))}
            </select>
          </div>
        )}

        {/* Count */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Count (1-20)
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

        {/* Max items per order — only for platform */}
        {isPlatform && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Max Items per Order (1-5)
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
        )}

        {/* Submit button */}
        <button
          onClick={handleBuildAndSend}
          disabled={loading || !integrationId}
          className="w-full py-2.5 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {loading && !sending
            ? "Building payloads..."
            : sending
            ? `Sending ${logs.length}/${payloads.length}...`
            : isEcommerce
            ? `Simulate ${count} Webhook${count > 1 ? "s" : ""}${selectedIntegration ? ` → ${selectedIntegration.integration_type_code}` : ""}`
            : `Generate ${count} Order${count > 1 ? "s" : ""}${selectedIntegration ? ` → ${selectedIntegration.integration_type_code}` : ""}`}
        </button>

        {/* Results summary */}
        {logs.length > 0 && (
          <div className="mt-4 space-y-3">
            <div className="flex gap-3">
              <div className="flex-1 p-3 bg-green-50 border border-green-200 rounded-lg text-center">
                <div className="text-2xl font-bold text-green-700">{successCount}</div>
                <div className="text-xs text-green-600">Success</div>
              </div>
              <div className="flex-1 p-3 bg-red-50 border border-red-200 rounded-lg text-center">
                <div className="text-2xl font-bold text-red-700">{failedCount}</div>
                <div className="text-xs text-red-600">Failed</div>
              </div>
              <div className="flex-1 p-3 bg-gray-50 border border-gray-200 rounded-lg text-center">
                <div className="text-2xl font-bold text-gray-700">{logs.length}</div>
                <div className="text-xs text-gray-600">Total</div>
              </div>
            </div>

            {/* API Logs */}
            <div>
              <h3 className="text-sm font-medium text-gray-700 mb-2">Request / Response Logs</h3>
              <div className="space-y-2 max-h-[600px] overflow-y-auto">
                {logs.map((log, idx) => (
                  <div
                    key={idx}
                    className={`border rounded-lg overflow-hidden ${
                      log.success ? "border-green-200" : "border-red-200"
                    }`}
                  >
                    {/* Log header */}
                    <button
                      onClick={() => setExpandedLog(expandedLog === idx ? null : idx)}
                      className={`w-full px-3 py-2 flex items-center justify-between text-sm ${
                        log.success ? "bg-green-50" : "bg-red-50"
                      }`}
                    >
                      <div className="flex items-center gap-2">
                        <span className={`font-mono text-xs px-1.5 py-0.5 rounded ${
                          log.success ? "bg-green-200 text-green-800" : "bg-red-200 text-red-800"
                        }`}>
                          {log.response.status_code || "ERR"}
                        </span>
                        <span className="font-medium text-gray-700">
                          {log.request.method} #{idx + 1}
                        </span>
                        <span className="text-gray-400 text-xs truncate max-w-[200px]">
                          {log.request.url}
                        </span>
                      </div>
                      <div className="flex items-center gap-2 text-xs text-gray-500">
                        <span>{log.duration_ms}ms</span>
                        <span>{expandedLog === idx ? "▲" : "▼"}</span>
                      </div>
                    </button>

                    {/* Expanded details */}
                    {expandedLog === idx && (
                      <div className="border-t border-gray-200">
                        {/* Request */}
                        <div className="p-3 border-b border-gray-100">
                          <div className="text-xs font-medium text-gray-500 mb-1">REQUEST BODY</div>
                          <pre className="text-xs bg-gray-900 text-green-400 p-3 rounded overflow-x-auto max-h-64 overflow-y-auto">
                            {JSON.stringify(log.request.body, null, 2)}
                          </pre>
                        </div>
                        {/* Response */}
                        <div className="p-3">
                          <div className="text-xs font-medium text-gray-500 mb-1">RESPONSE BODY</div>
                          <pre className="text-xs bg-gray-900 text-green-400 p-3 rounded overflow-x-auto max-h-64 overflow-y-auto">
                            {(() => {
                              try {
                                return JSON.stringify(JSON.parse(log.response.body), null, 2);
                              } catch {
                                return log.response.body;
                              }
                            })()}
                          </pre>
                        </div>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
