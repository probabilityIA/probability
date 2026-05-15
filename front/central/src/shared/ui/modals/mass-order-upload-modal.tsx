'use client';

import { useState } from 'react';
import { Button, Input } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { getActionError } from '@/shared/utils/action-result';
import { uploadBulkOrdersAction, downloadOrderTemplateAction } from '@/services/modules/orders/infra/actions';

interface MassOrderUploadModalProps {
    isOpen: boolean;
    onClose: () => void;
    onUploadComplete?: (count: number) => void;
    selectedBusinessId?: number | null;
}

export default function MassOrderUploadModal({ isOpen, onClose, onUploadComplete, selectedBusinessId }: MassOrderUploadModalProps) {
    const { isSuperAdmin } = usePermissions();
    const [file, setFile] = useState<File | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [uploadStats, setUploadStats] = useState<{ total: number; success: number; failed: number, errors?: string[] } | null>(null);

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const selectedFile = e.target.files?.[0];
        if (selectedFile) {
            const ext = selectedFile.name.split('.').pop()?.toLowerCase();
            if (ext === 'csv' || ext === 'xlsx' || ext === 'xls') {
                setFile(selectedFile);
                setError(null);
            } else {
                setError('Solo se permiten archivos CSV o Excel (.xlsx, .xls)');
                setFile(null);
            }
        }
    };

    const handleUpload = async () => {
        if (!file) {
            setError('Por favor selecciona un archivo');
            return;
        }

        if (isSuperAdmin && !selectedBusinessId) {
            setError('Debes seleccionar un negocio antes de cargar órdenes');
            return;
        }

        setLoading(true);
        setError(null);
        setSuccess(null);
        setUploadStats(null);

        try {
            const result = await uploadBulkOrdersAction(file, isSuperAdmin ? selectedBusinessId ?? undefined : undefined);

            if (result.success) {
                if (result.data.success_count > 0) {
                    setSuccess(`¡Proceso completado! ${result.data.success_count} órdenes creadas.`);
                } else {
                    setError('No se pudo crear ninguna orden. Revisa los errores abajo.');
                }
                setUploadStats({
                    total: result.data.total_rows,
                    success: result.data.success_count,
                    failed: result.data.failed_count,
                    errors: result.data.errors,
                });
                if (onUploadComplete && result.data.success_count > 0) {
                    onUploadComplete(result.data.success_count);
                }
            } else {
                setError(result.message || 'Error al procesar el archivo');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar el archivo'));
        } finally {
            setLoading(false);
        }
    };

    const handleDownloadTemplate = async () => {
        try {
            const result = await downloadOrderTemplateAction();
            if (result.success && result.data) {
                const blob = new Blob([result.data], { type: 'text/csv;charset=utf-8;' });
                const link = document.createElement('a');
                const url = URL.createObjectURL(blob);
                link.setAttribute('href', url);
                link.setAttribute('download', 'plantilla_ordenes.csv');
                link.style.visibility = 'hidden';
                document.body.appendChild(link);
                link.click();
                document.body.removeChild(link);
            }
        } catch (err) {
            console.error('Error downloading template:', err);
        }
    };

    const handleClose = () => {
        setFile(null);
        setError(null);
        setSuccess(null);
        setUploadStats(null);
        onClose();
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full p-6 max-h-[90vh] flex flex-col">
                <div className="flex justify-between items-center mb-6 flex-shrink-0">
                    <h2 className="text-2xl font-bold text-gray-800 dark:text-gray-100">Carga Masiva de Órdenes</h2>
                    <button
                        onClick={handleClose}
                        className="text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:text-gray-200 dark:text-gray-200 text-2xl"
                    >
                        ×
                    </button>
                </div>

                <div className="space-y-6 overflow-y-auto pr-2">
                    {/* Instructions */}
                    <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                        <div className="flex justify-between items-start mb-2">
                            <h3 className="font-semibold text-blue-800">Instrucciones</h3>
                            <button
                                onClick={handleDownloadTemplate}
                                className="text-xs font-semibold text-white bg-blue-600 hover:bg-blue-700 px-3 py-1.5 rounded flex items-center gap-1 transition-colors"
                            >
                                Descargar Plantilla CSV
                            </button>
                        </div>
                        <ul className="text-sm text-blue-700 space-y-1 list-disc list-inside">
                            <li>El archivo debe ser CSV o Excel (.xlsx, .xls)</li>
                            <li>La primera fila debe contener los encabezados de columna tal cual la plantilla</li>
                            <li><strong>Requeridas:</strong> order_number, customer_name, customer_email, customer_phone, shipping_street, shipping_city, shipping_state, total_amount</li>
                            <li><strong>Cliente (opcionales):</strong> customer_first_name, customer_last_name, customer_dni</li>
                            <li><strong>Dirección (opcionales):</strong> shipping_country, shipping_postal_code, shipping_lat, shipping_lng</li>
                            <li><strong>Financiero (opcionales):</strong> subtotal, tax, discount, shipping_cost, shipping_discount, currency</li>
                            <li><strong>Logística (opcionales):</strong> tracking_number, guide_id, warehouse_name, driver_name, platform</li>
                            <li><strong>Estado (opcionales):</strong> status, payment_method_id, is_paid, order_type_name, invoiceable</li>
                            <li><strong>Otros (opcionales):</strong> weight, height, width, length, notes</li>
                        </ul>
                    </div>

                    {/* File Input */}
                    {!success && !uploadStats && (
                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 dark:text-gray-200 mb-2">
                                Seleccionar archivo
                            </label>
                            <input
                                type="file"
                                accept=".csv,.xlsx,.xls"
                                onChange={handleFileChange}
                                className="block w-full text-sm text-gray-500
                                    file:mr-4 file:py-2 file:px-4
                                    file:rounded-md file:border-0
                                    file:text-sm file:font-semibold
                                    file:bg-orange-50 file:text-orange-700
                                    hover:file:bg-orange-100"
                            />
                            {file && (
                                <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">
                                    Archivo seleccionado: <strong>{file.name}</strong>
                                </p>
                            )}
                        </div>
                    )}

                    {/* Error / Results */}
                    {(error || uploadStats) && (
                        <div className={`p-4 rounded-lg border ${error ? 'bg-red-50 border-red-200 text-red-700' : 'bg-green-50 border-green-200 text-green-700'}`}>
                            {error && <p className="font-bold mb-2">{error}</p>}
                            {success && <p className="font-bold mb-2">{success}</p>}

                            {uploadStats && (
                                <div className="mt-2 text-sm space-y-1">
                                    <p>Total de filas: <strong>{uploadStats.total}</strong></p>
                                    <p>Exitosas: <strong className="text-green-600">{uploadStats.success}</strong></p>
                                    <p>Fallidas: <strong className="text-red-600">{uploadStats.failed}</strong></p>

                                    {uploadStats.errors && uploadStats.errors.length > 0 && (
                                        <div className="mt-4">
                                            <p className="font-bold text-red-800 mb-1">Detalle de errores:</p>
                                            <ul className="list-disc list-inside text-xs max-h-40 overflow-y-auto bg-white/50 p-2 rounded">
                                                {uploadStats.errors.map((err, idx) => (
                                                    <li key={idx} className="mb-1">{err}</li>
                                                ))}
                                            </ul>
                                        </div>
                                    )}
                                </div>
                            )}
                        </div>
                    )}

                    {/* Actions */}
                    <div className="flex justify-end space-x-3 pt-4 border-t border-gray-100 flex-shrink-0">
                        {success || uploadStats ? (
                            <Button onClick={handleClose}>
                                Cerrar
                            </Button>
                        ) : (
                            <>
                                <Button
                                    variant="outline"
                                    onClick={handleClose}
                                    disabled={loading}
                                >
                                    Cancelar
                                </Button>
                                <Button
                                    onClick={handleUpload}
                                    disabled={!file || loading}
                                >
                                    {loading ? 'Procesando...' : 'Cargar Órdenes'}
                                </Button>
                            </>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}
