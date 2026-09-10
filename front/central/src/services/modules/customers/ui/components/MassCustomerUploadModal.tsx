'use client';

import { useState } from 'react';
import * as XLSX from 'xlsx';
import { ArrowUpTrayIcon, DocumentArrowDownIcon } from '@heroicons/react/24/outline';
import { Button } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';
import { uploadBulkCustomersAction, downloadCustomerTemplateAction } from '../../infra/actions';
import { BulkCustomerResult } from '../../domain/types';

interface MassCustomerUploadModalProps {
    isOpen: boolean;
    onClose: () => void;
    onUploadComplete?: (count: number) => void;
    selectedBusinessId?: number | null;
}

export default function MassCustomerUploadModal({ isOpen, onClose, onUploadComplete, selectedBusinessId }: MassCustomerUploadModalProps) {
    const [file, setFile] = useState<File | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [result, setResult] = useState<BulkCustomerResult | null>(null);

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const selectedFile = e.target.files?.[0];
        if (!selectedFile) return;
        const ext = selectedFile.name.split('.').pop()?.toLowerCase();
        if (ext === 'csv' || ext === 'xlsx' || ext === 'xls') {
            setFile(selectedFile);
            setError(null);
        } else {
            setError('Solo se permiten archivos CSV o Excel (.xlsx, .xls)');
            setFile(null);
        }
    };

    const handleUpload = async () => {
        if (!file) {
            setError('Selecciona un archivo primero');
            return;
        }
        if (!selectedBusinessId) {
            setError('Debes seleccionar un negocio antes de cargar clientes');
            return;
        }

        setLoading(true);
        setError(null);
        setResult(null);

        try {
            const response = await uploadBulkCustomersAction(file, selectedBusinessId);
            if (response.success && response.data) {
                setResult(response.data);
                if (response.data.success_count > 0) {
                    onUploadComplete?.(response.data.success_count);
                }
            } else {
                setError(response.message || 'Error al procesar el archivo');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar el archivo'));
        } finally {
            setLoading(false);
        }
    };

    const handleDownloadTemplate = async () => {
        const response = await downloadCustomerTemplateAction();
        if (response.success && response.data) {
            const { headers, exampleRows } = response.data;
            const ws = XLSX.utils.aoa_to_sheet([headers, ...exampleRows]);
            const wb = XLSX.utils.book_new();
            XLSX.utils.book_append_sheet(wb, ws, 'Clientes');
            XLSX.writeFile(wb, 'plantilla_clientes.xlsx');
        }
    };

    const handleClose = () => {
        setFile(null);
        setError(null);
        setResult(null);
        onClose();
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-2xl max-h-[90vh] flex flex-col">
                <div className="flex items-center justify-between px-6 py-4 border-b border-gray-100 dark:border-gray-700 flex-shrink-0">
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Carga masiva de clientes</h2>
                    <button
                        onClick={handleClose}
                        className="text-gray-400 hover:text-gray-600 dark:text-gray-300 text-xl leading-none"
                    >
                        x
                    </button>
                </div>

                <div className="p-6 space-y-5 overflow-y-auto">
                    <div className="bg-purple-50 dark:bg-purple-950/30 border border-purple-200 dark:border-purple-800 rounded-lg p-4">
                        <div className="flex items-start justify-between gap-3 mb-2">
                            <h3 className="font-semibold text-purple-800 dark:text-purple-300 text-sm">Instrucciones</h3>
                            <button
                                onClick={handleDownloadTemplate}
                                className="inline-flex items-center gap-1.5 text-xs font-semibold text-white bg-purple-600 hover:bg-purple-700 px-3 py-1.5 rounded-md transition-colors flex-shrink-0"
                            >
                                <DocumentArrowDownIcon className="w-4 h-4" />
                                Descargar plantilla
                            </button>
                        </div>
                        <ul className="text-sm text-purple-700 dark:text-purple-300 space-y-1 list-disc list-inside">
                            <li>El archivo debe ser CSV o Excel (.xlsx, .xls)</li>
                            <li>La primera fila debe tener los encabezados tal cual la plantilla</li>
                            <li><strong>Obligatorias:</strong> nombre, apellido, cedula, correo, telefono</li>
                            <li><strong>Opcionales:</strong> direccion, ciudad, notas</li>
                            <li>Si un cliente ya existe (mismo correo o cedula), esa fila se marca como fallida y el resto continua</li>
                        </ul>
                    </div>

                    {!result && (
                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                                Seleccionar archivo
                            </label>
                            <label className="flex items-center justify-center gap-2 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg py-6 cursor-pointer hover:border-purple-400 dark:hover:border-purple-500 transition-colors">
                                <ArrowUpTrayIcon className="w-5 h-5 text-gray-400" />
                                <span className="text-sm text-gray-500 dark:text-gray-400">
                                    {file ? file.name : 'Haz clic para elegir un archivo CSV o Excel'}
                                </span>
                                <input
                                    type="file"
                                    accept=".csv,.xlsx,.xls"
                                    onChange={handleFileChange}
                                    className="hidden"
                                />
                            </label>
                        </div>
                    )}

                    {error && (
                        <div className="p-4 rounded-lg border bg-red-50 border-red-200 text-red-700 dark:bg-red-950/30 dark:border-red-800 dark:text-red-300 text-sm">
                            {error}
                        </div>
                    )}

                    {result && (
                        <div className="space-y-3">
                            <div className="grid grid-cols-3 gap-3">
                                <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-3 text-center">
                                    <div className="text-2xl font-bold text-gray-900 dark:text-white">{result.total_rows}</div>
                                    <div className="text-xs text-gray-500 dark:text-gray-400">Total filas</div>
                                </div>
                                <div className="bg-green-50 dark:bg-green-950/30 rounded-lg p-3 text-center">
                                    <div className="text-2xl font-bold text-green-600 dark:text-green-400">{result.success_count}</div>
                                    <div className="text-xs text-green-700 dark:text-green-400">Exitosos</div>
                                </div>
                                <div className="bg-red-50 dark:bg-red-950/30 rounded-lg p-3 text-center">
                                    <div className="text-2xl font-bold text-red-600 dark:text-red-400">{result.failed_count}</div>
                                    <div className="text-xs text-red-700 dark:text-red-400">Fallidos</div>
                                </div>
                            </div>

                            {result.results.filter((r) => !r.success).length > 0 && (
                                <div>
                                    <p className="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">Filas con error</p>
                                    <ul className="text-xs space-y-1 max-h-40 overflow-y-auto bg-gray-50 dark:bg-gray-700/50 rounded-lg p-3">
                                        {result.results.filter((r) => !r.success).map((r) => (
                                            <li key={r.row} className="text-red-700 dark:text-red-300">
                                                Fila {r.row} ({r.name || 'sin nombre'}): {r.error}
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                            )}
                        </div>
                    )}
                </div>

                <div className="flex justify-end gap-3 px-6 py-4 border-t border-gray-100 dark:border-gray-700 flex-shrink-0">
                    {result ? (
                        <Button onClick={handleClose}>Cerrar</Button>
                    ) : (
                        <>
                            <Button variant="outline" onClick={handleClose} disabled={loading}>
                                Cancelar
                            </Button>
                            <Button onClick={handleUpload} disabled={!file || loading}>
                                {loading ? 'Procesando...' : 'Cargar clientes'}
                            </Button>
                        </>
                    )}
                </div>
            </div>
        </div>
    );
}
