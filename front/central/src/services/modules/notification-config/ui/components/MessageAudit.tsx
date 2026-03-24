'use client';

import { useState, useEffect, useCallback } from 'react';
import { MessageAuditStatsRow } from './MessageAuditStatsRow';
import {
  getMessageAuditLogsAction,
  getMessageAuditStatsAction,
} from '../../infra/actions';
import type {
  MessageAuditLog,
  MessageAuditStats,
} from '../../domain/types';

interface MessageAuditProps {
  businessId?: number;
}

function formatDate(days: number): string {
  const d = new Date();
  d.setDate(d.getDate() - days);
  return d.toISOString().split('T')[0];
}

function maskPhone(phone: string): string {
  if (phone.length <= 6) return phone;
  return phone.slice(0, 3) + '***' + phone.slice(-4);
}

const statusBadge: Record<string, { bg: string; text: string; label: string }> = {
  sent: { bg: 'bg-purple-100 dark:bg-purple-600', text: 'text-purple-700 dark:text-white', label: 'Enviado' },
  delivered: { bg: 'bg-green-100 dark:bg-green-600', text: 'text-green-700 dark:text-white', label: 'Entregado' },
  read: { bg: 'bg-emerald-100 dark:bg-emerald-600', text: 'text-emerald-700 dark:text-white', label: 'Leido' },
  failed: { bg: 'bg-red-100 dark:bg-red-600', text: 'text-red-700 dark:text-white', label: 'Fallido' },
};

const directionLabel: Record<string, string> = {
  outbound: 'Saliente',
  inbound: 'Entrante',
};

export function MessageAudit({ businessId }: MessageAuditProps) {
  const [logs, setLogs] = useState<MessageAuditLog[]>([]);
  const [stats, setStats] = useState<MessageAuditStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [statsLoading, setStatsLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);

  // Filters
  const [dateFrom, setDateFrom] = useState(formatDate(7));
  const [dateTo, setDateTo] = useState(formatDate(0));
  const [status, setStatus] = useState('');
  const [templateName, setTemplateName] = useState('');
  const [page, setPage] = useState(1);
  const pageSize = 20;

  const fetchStats = useCallback(async () => {
    setStatsLoading(true);
    try {
      const res = await getMessageAuditStatsAction(
        businessId ?? 0,
        dateFrom || undefined,
        dateTo || undefined,
      );
      if (res.success && res.data) {
        setStats(res.data);
      }
    } catch {
      // stats are non-critical
    } finally {
      setStatsLoading(false);
    }
  }, [businessId, dateFrom, dateTo]);

  const fetchLogs = useCallback(async () => {
    setLoading(true);
    try {
      const res = await getMessageAuditLogsAction({
        business_id: businessId ?? 0,
        status: status || undefined,
        template_name: templateName || undefined,
        date_from: dateFrom || undefined,
        date_to: dateTo || undefined,
        page,
        page_size: pageSize,
      });
      if (res.success) {
        setLogs(res.data);
        setTotal(res.total);
        setTotalPages(res.total_pages);
      }
    } catch {
      setLogs([]);
    } finally {
      setLoading(false);
    }
  }, [businessId, status, templateName, dateFrom, dateTo, page]);

  useEffect(() => {
    fetchStats();
  }, [fetchStats]);

  useEffect(() => {
    setPage(1);
  }, [status, templateName, dateFrom, dateTo, businessId]);

  useEffect(() => {
    fetchLogs();
  }, [fetchLogs]);

  // businessId is undefined for regular users (backend resolves from JWT)
  // For super admins, businessId must be provided via prop

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
      {/* Header */}
      <div className="p-4 border-b flex flex-col sm:flex-row items-start sm:items-center gap-3">
        <div className="flex items-center gap-2">
          <h3 className="text-sm font-medium text-gray-900 dark:text-white">Auditoria de Mensajes</h3>
        </div>
        <div className="flex items-center gap-2 ml-auto">
          <input
            type="date"
            value={dateFrom}
            onChange={(e) => setDateFrom(e.target.value)}
            className="px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-purple-500 dark:focus:ring-purple-400"
          />
          <span className="text-xs text-gray-400 dark:text-white">a</span>
          <input
            type="date"
            value={dateTo}
            onChange={(e) => setDateTo(e.target.value)}
            className="px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-purple-500 dark:focus:ring-purple-400"
          />
        </div>
      </div>

      <div className="p-4 space-y-4">
        {/* Stats */}
        <MessageAuditStatsRow stats={stats} loading={statsLoading} />

        {/* Filters */}
        <div className="flex flex-wrap gap-2">
          <select
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            className="px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-purple-500 dark:focus:ring-purple-400"
          >
            <option value="">Todos los estados</option>
            <option value="sent">Enviado</option>
            <option value="delivered">Entregado</option>
            <option value="read">Leido</option>
            <option value="failed">Fallido</option>
          </select>
          <input
            type="text"
            value={templateName}
            onChange={(e) => setTemplateName(e.target.value)}
            placeholder="Buscar plantilla..."
            className="px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-purple-500 dark:focus:ring-purple-400 w-40"
          />
        </div>

        {/* Table */}
        <div className="overflow-x-auto">
          <table className="w-full text-xs">
            <thead className="bg-purple-600 dark:bg-purple-700">
              <tr className="border-b border-purple-700 dark:border-purple-800">
                <th className="text-left py-2 px-2 font-medium text-white">Fecha</th>
                <th className="text-left py-2 px-2 font-medium text-white">Destino</th>
                <th className="text-left py-2 px-2 font-medium text-white">Orden</th>
                <th className="text-left py-2 px-2 font-medium text-white">Plantilla</th>
                <th className="text-center py-2 px-2 font-medium text-white">Estado</th>
                <th className="text-center py-2 px-2 font-medium text-white">Direccion</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                Array.from({ length: 5 }).map((_, i) => (
                  <tr key={i} className="border-b border-gray-100 dark:border-gray-700">
                    {Array.from({ length: 6 }).map((_, j) => (
                      <td key={j} className="py-2 px-2">
                        <div className="h-4 bg-gray-200 rounded animate-pulse" />
                      </td>
                    ))}
                  </tr>
                ))
              ) : logs.length === 0 ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-gray-400">
                    No hay mensajes en el periodo seleccionado
                  </td>
                </tr>
              ) : (
                logs.map((log) => {
                  const badge = statusBadge[log.status] || statusBadge.sent;
                  return (
                    <tr key={log.id} className="border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
                      <td className="py-2 px-2 text-gray-600 dark:text-white whitespace-nowrap">
                        {new Date(log.created_at).toLocaleDateString('es-CO', {
                          day: '2-digit',
                          month: 'short',
                          hour: '2-digit',
                          minute: '2-digit',
                        })}
                      </td>
                      <td className="py-2 px-2 text-gray-700 dark:text-white font-mono">
                        {maskPhone(log.phone_number)}
                      </td>
                      <td className="py-2 px-2 text-gray-600 dark:text-white">
                        {log.order_number || '-'}
                      </td>
                      <td className="py-2 px-2 text-gray-600 dark:text-white max-w-[120px] truncate">
                        {log.template_name || '-'}
                      </td>
                      <td className="py-2 px-2 text-center">
                        <span className={`inline-block px-1.5 py-0.5 rounded-full text-[10px] font-medium ${badge.bg} ${badge.text}`}>
                          {badge.label}
                        </span>
                      </td>
                      <td className="py-2 px-2 text-center text-gray-500 dark:text-white">
                        {directionLabel[log.direction] || log.direction}
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between pt-2">
            <p className="text-xs text-gray-500 dark:text-gray-400">
              Pagina {page} de {totalPages} ({total} registros)
            </p>
            <div className="flex gap-1">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1}
                className="px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300 transition-colors"
              >
                Anterior
              </button>
              <button
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
                className="px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300 transition-colors"
              >
                Siguiente
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
