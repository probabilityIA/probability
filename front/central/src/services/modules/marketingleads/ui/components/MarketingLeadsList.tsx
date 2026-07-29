'use client';

import React from 'react';
import { Alert } from '@/shared/ui/alert';
import { useMarketingLeads } from '../hooks/useMarketingLeads';

const SURVEY_LABELS: Record<string, string> = {
    'diagnostico-logistico': 'Diagnóstico Logístico',
    'escanea-tu-logistica': 'Diagnóstico de Operación E-commerce',
};

export const MarketingLeadsList: React.FC = () => {
    const {
        leads,
        loading,
        error,
        page,
        setPage,
        pageSize,
        setPageSize,
        totalPages,
        total,
        setError,
    } = useMarketingLeads();

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Leads de encuestas</h1>
            </div>

            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm dark:shadow-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="table min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                        <thead style={{ backgroundColor: 'var(--color-primary)', color: 'white' }}>
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Nombre</th>
                                <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Correo</th>
                                <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Teléfono</th>
                                <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Encuesta</th>
                                <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Puntaje</th>
                                <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Nivel</th>
                                <th className="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider">WhatsApp</th>
                                <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Fecha</th>
                            </tr>
                        </thead>
                        <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                            {loading ? (
                                <tr>
                                    <td colSpan={8} className="px-6 py-4 text-center text-sm text-gray-500 dark:text-gray-400">
                                        Cargando leads...
                                    </td>
                                </tr>
                            ) : leads.length === 0 ? (
                                <tr>
                                    <td colSpan={8} className="px-6 py-4 text-center text-sm text-gray-500 dark:text-gray-400">
                                        No hay leads capturados todavía
                                    </td>
                                </tr>
                            ) : (
                                leads.map((lead) => (
                                    <tr key={lead.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                                            {lead.name}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                                            {lead.email}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                                            {lead.phone}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                                            {SURVEY_LABELS[lead.survey_slug] || lead.survey_slug}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                                            {lead.score_total} / {lead.score_max}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                                            {lead.level}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-center text-sm">
                                            {lead.whats_app_message_id ? (
                                                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">Enviado</span>
                                            ) : (
                                                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-600">Pendiente</span>
                                            )}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                                            {new Date(lead.created_at).toLocaleString('es-CO')}
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>

                {!loading && leads.length > 0 && (
                    <div className="bg-white dark:bg-gray-800 px-4 py-3 border-t border-gray-200 dark:border-gray-700 sm:px-6">
                        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
                            <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                                <p className="text-sm text-gray-700 dark:text-gray-200">
                                    Mostrando <span className="font-medium">{(page - 1) * pageSize + 1}</span> a{' '}
                                    <span className="font-medium">{Math.min(page * pageSize, total)}</span> de{' '}
                                    <span className="font-medium">{total}</span> resultados
                                </p>
                                <nav className="flex items-center gap-2">
                                    <button
                                        onClick={() => setPage(page - 1)}
                                        disabled={page === 1}
                                        className="btn btn-secondary rounded-l-md rounded-r-none disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        Anterior
                                    </button>
                                    <span
                                        className="relative inline-flex items-center px-3 sm:px-4 py-2 border text-xs sm:text-sm font-medium text-gray-700 dark:text-gray-200"
                                        style={{ borderColor: 'var(--color-secondary-500)' }}
                                    >
                                        Página {page} de {totalPages}
                                    </span>
                                    <button
                                        onClick={() => setPage(page + 1)}
                                        disabled={page === totalPages}
                                        className="btn btn-secondary rounded-r-md rounded-l-none disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        Siguiente
                                    </button>
                                </nav>
                            </div>

                            <div className="flex items-center justify-between w-full sm:hidden pt-2 border-t border-gray-200 dark:border-gray-700">
                                <div className="flex items-center gap-2">
                                    <label className="text-xs text-gray-700 dark:text-gray-200 whitespace-nowrap">Mostrar:</label>
                                    <select
                                        value={pageSize}
                                        onChange={(e) => { setPageSize(parseInt(e.target.value)); setPage(1); }}
                                        className="px-2 py-1.5 text-xs border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 dark:text-white bg-white dark:bg-gray-700"
                                    >
                                        <option value="10">10</option>
                                        <option value="20">20</option>
                                        <option value="50">50</option>
                                        <option value="100">100</option>
                                    </select>
                                </div>
                                <p className="text-xs text-gray-500 dark:text-gray-400">Página {page} de {totalPages}</p>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};
