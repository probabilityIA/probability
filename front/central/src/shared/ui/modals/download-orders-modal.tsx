'use client';

import { useState } from 'react';
import { Button } from '@/shared/ui';

interface DownloadOrdersModalProps {
    isOpen: boolean;
    onClose: () => void;
    onDownload: (startDate: string, endDate: string) => Promise<void>;
}

export default function DownloadOrdersModal({ isOpen, onClose, onDownload }: DownloadOrdersModalProps) {
    const [startDate, setStartDate] = useState('');
    const [endDate, setEndDate] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    if (!isOpen) return null;

    const handleDownload = async () => {
        setError(null);

        if (!startDate || !endDate) {
            setError('Debes seleccionar una fecha de inicio y fin');
            return;
        }

        if (new Date(startDate) > new Date(endDate)) {
            setError('La fecha de inicio no puede ser mayor que la fecha de fin');
            return;
        }

        const maxDays = 365;
        const diffTime = Math.abs(new Date(endDate).getTime() - new Date(startDate).getTime());
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

        if (diffDays > maxDays) {
            setError(`El rango no puede exceder ${maxDays} días`);
            return;
        }

        setIsLoading(true);
        try {
            await onDownload(startDate, endDate);
            onClose();
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Error al descargar órdenes');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="fixed inset-0 backdrop-blur-sm flex items-center justify-center z-50" style={{ backgroundColor: 'rgba(255, 255, 255, 0.02)' }}>
            <div className="bg-white rounded-lg shadow-lg p-6 max-w-md w-full mx-4">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Descargar Órdenes</h2>

                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 dark:text-gray-200 mb-1">
                            Fecha de Inicio
                        </label>
                        <input
                            type="date"
                            value={startDate}
                            onChange={(e) => setStartDate(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                            disabled={isLoading}
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 dark:text-gray-200 mb-1">
                            Fecha de Fin
                        </label>
                        <input
                            type="date"
                            value={endDate}
                            onChange={(e) => setEndDate(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                            disabled={isLoading}
                        />
                    </div>

                    {error && (
                        <div className="p-3 bg-red-50 border border-red-200 rounded-lg">
                            <p className="text-sm text-red-700">{error}</p>
                        </div>
                    )}

                    <p className="text-xs text-gray-500 dark:text-gray-400">
                        Máximo 1 año entre fechas
                    </p>
                </div>

                <div className="flex gap-3 mt-6">
                    <button
                        onClick={onClose}
                        disabled={isLoading}
                        className="flex-1 px-4 py-2 border border-gray-300 rounded-lg text-gray-700 dark:text-gray-200 dark:text-gray-200 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                        Cancelar
                    </button>
                    <button
                        onClick={handleDownload}
                        disabled={isLoading}
                        className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
                    >
                        {isLoading ? (
                            <>
                                <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                                Descargando...
                            </>
                        ) : (
                            '↓ Descargar'
                        )}
                    </button>
                </div>
            </div>
        </div>
    );
}
