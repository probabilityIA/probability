"use client";

import { useState } from "react";
import type { APICallLog } from "../../domain/types";

interface Props {
  logs: APICallLog[];
}

function JsonSyntaxHighlight({ json }: { json: string }) {
  const highlighted = json.replace(
    /("(?:\\.|[^"\\])*")\s*:/g,
    '<span class="text-purple-400">$1</span>:'
  ).replace(
    /:\s*("(?:\\.|[^"\\])*")/g,
    ': <span class="text-green-400">$1</span>'
  ).replace(
    /:\s*(\d+\.?\d*)/g,
    ': <span class="text-yellow-300">$1</span>'
  ).replace(
    /:\s*(true|false)/g,
    ': <span class="text-cyan-400">$1</span>'
  ).replace(
    /:\s*(null)/g,
    ': <span class="text-gray-500">$1</span>'
  );

  return (
    <pre
      className="text-sm text-gray-300 whitespace-pre overflow-x-auto"
      dangerouslySetInnerHTML={{ __html: highlighted }}
    />
  );
}

function formatJson(value: unknown): string {
  try {
    if (typeof value === "string") {
      const parsed = JSON.parse(value);
      return JSON.stringify(parsed, null, 2);
    }
    return JSON.stringify(value, null, 2);
  } catch {
    return typeof value === "string" ? value : String(value);
  }
}

function StatusBadge({ code }: { code: number }) {
  const isSuccess = code >= 200 && code < 300;
  const bg = isSuccess ? "bg-green-500/20 text-green-400" : "bg-red-500/20 text-red-400";
  return (
    <span className={`px-2 py-0.5 rounded text-xs font-mono font-bold ${bg}`}>
      {code || "ERR"}
    </span>
  );
}

function LogEntry({ log }: { log: APICallLog }) {
  const [expanded, setExpanded] = useState(false);

  const pathMatch = log.request.url.match(/\/api\/.*/);
  const path = pathMatch ? pathMatch[0] : log.request.url;

  return (
    <div className="border border-gray-700/50 rounded-lg overflow-hidden">
      <button
        onClick={() => setExpanded(!expanded)}
        className="w-full flex items-center gap-3 px-4 py-2.5 hover:bg-gray-800/50 transition-colors text-left"
      >
        <span className="text-gray-500 font-mono text-xs w-6 text-right shrink-0">
          #{log.index + 1}
        </span>
        <span className="font-mono text-xs text-blue-400 font-bold">
          {log.request.method}
        </span>
        <span className="font-mono text-xs text-gray-400 truncate flex-1">
          {path}
        </span>
        <StatusBadge code={log.response.status_code} />
        <span className="text-xs text-gray-500 font-mono shrink-0">
          {log.duration_ms}ms
        </span>
        <svg
          className={`w-4 h-4 text-gray-500 transition-transform ${expanded ? "rotate-180" : ""}`}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {expanded && (
        <div className="border-t border-gray-700/50 divide-y divide-gray-700/30">
          <div className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <span className="text-xs font-bold text-blue-400 uppercase tracking-wider">Request</span>
              <span className="text-xs text-gray-500 font-mono">{log.request.url}</span>
            </div>
            <div className="bg-black/30 rounded-lg p-3 overflow-x-auto">
              <JsonSyntaxHighlight json={formatJson(log.request.body)} />
            </div>
          </div>

          <div className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <span className="text-xs font-bold text-green-400 uppercase tracking-wider">Response</span>
              <StatusBadge code={log.response.status_code} />
            </div>
            <div className="bg-black/30 rounded-lg p-3 overflow-x-auto">
              <JsonSyntaxHighlight json={formatJson(log.response.body)} />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default function APIConsole({ logs }: Props) {
  const [collapsed, setCollapsed] = useState(false);
  const successCount = logs.filter((l) => l.success).length;
  const failCount = logs.length - successCount;

  return (
    <div className="bg-[#0d1117] rounded-lg border border-gray-700/50 overflow-hidden">
      <button
        onClick={() => setCollapsed(!collapsed)}
        className="w-full flex items-center justify-between px-4 py-3 border-b border-gray-700/50 hover:bg-gray-800/30 transition-colors"
      >
        <div className="flex items-center gap-3">
          <svg className="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <span className="text-sm font-semibold text-gray-200">API Console</span>
          <span className="text-xs text-gray-500 bg-gray-800 px-2 py-0.5 rounded-full font-mono">
            {logs.length} call{logs.length !== 1 ? "s" : ""}
          </span>
          {successCount > 0 && (
            <span className="text-xs text-green-400 bg-green-500/10 px-2 py-0.5 rounded-full font-mono">
              {successCount} ok
            </span>
          )}
          {failCount > 0 && (
            <span className="text-xs text-red-400 bg-red-500/10 px-2 py-0.5 rounded-full font-mono">
              {failCount} err
            </span>
          )}
        </div>
        <svg
          className={`w-4 h-4 text-gray-500 transition-transform ${collapsed ? "" : "rotate-180"}`}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {!collapsed && (
        <div className="p-3 space-y-2 max-h-[600px] overflow-y-auto">
          {logs.map((log) => (
            <LogEntry key={log.index} log={log} />
          ))}
        </div>
      )}
    </div>
  );
}
