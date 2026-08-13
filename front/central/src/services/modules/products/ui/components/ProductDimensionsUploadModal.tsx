'use client';

import { useState } from 'react';
import { Button } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';
import { exportProductDimensionsAction, importProductDimensionsAction } from '../../infra/actions';

interface ProductDimensionsUploadModalProps {
    isOpen: boolean;
    onClose: () => void;
    onUploadComplete?: (count: number) => void;
    businessId?: number | null;
}

export default function ProductDimensionsUploadModal({ isOpen, onClose, onUploadComplete, businessId }: ProductDimensionsUploadModalProps) {
    const [file, setFile] = useState<File | null>(null);
    const [loading, setLoading] = useState(false);
    const [downloading, setDownloading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [uploadStats, setUploadStats] = useState<{ total: number; success: number; failed: number; errors?: string[] } | null>(null);

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const selectedFile = e.target.files?.[0];
        if (selectedFile) {
            const ext = selectedFile.name.split('.').pop()?.toLowerCase();
            if (ext === 'xlsx' || ext === 'xls') {
                setFile(selectedFile);
                setError(null);
            } else {
                setError('Solo se permiten archivos Excel (.xlsx, .xls)');
                setFile(null);
            }
        }
    };

    const handleDownload = async () => {
        setDownloading(true);
        setError(null);
        try {
            const result = await exportProductDimensionsAction(businessId ?? undefined);
            if (result.success && result.blob) {
                const a = document.createElement('a');
                a.href = result.blob;
                a.download = result.filename || 'dimensiones-productos.xlsx';
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
            } else {
                setError(result.message || 'Error al descargar el archivo');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al descargar el archivo'));
        } finally {
            setDownloading(false);
        }
    };

    const handleUpload = async () => {
        if (!file) {
            setError('Por favor selecciona un archivo');
            return;
        }

        setLoading(true);
        setError(null);
        setSuccess(null);
        setUploadStats(null);

        try {
            const result = await importProductDimensionsAction(file, businessId ?? undefined);

            if (result.success) {
                if (result.data.success_count > 0) {
                    setSuccess(`¡Proceso completado! ${result.data.success_count} productos actualizados.`);
                } else {
                    setError('No se actualizó ningún producto. Revisa los errores abajo.');
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

    const handleClose = () => {
        setFile(null);
        setError(null);
        setSuccess(null);
        setUploadStats(null);
        onClose();
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full p-6 max-h-[90vh] flex flex-col">
                <div className="flex justify-between items-center mb-6 flex-shrink-0">
                    <h2 className="text-2xl font-bold text-gray-800 dark:text-gray-100">Peso y dimensiones por Excel</h2>
                    <button
                        onClick={handleClose}
                        className="text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:text-gray-200 text-2xl"
                    >
                        ×
                    </button>
                </div>

                <div className="space-y-6 overflow-y-auto pr-2">
                    <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                        <div className="flex justify-between items-start mb-2">
                            <h3 className="font-semibold text-blue-800">Instrucciones</h3>
                            <Button
                                onClick={handleDownload}
                                disabled={downloading}
                                className="text-xs font-semibold text-white bg-blue-600 hover:bg-blue-700 px-3 py-1.5 rounded flex items-center gap-1 transition-colors"
                            >
                                {downloading ? 'Generando...' : 'Descargar Excel'}
                            </Button>
                        </div>
                        <ul className="text-sm text-blue-700 space-y-1 list-disc list-inside">
                            <li>El Excel trae todos tus productos, con columnas sku, id, nombre, peso_kg, largo_cm, ancho_cm, alto_cm</li>
                            <li>Solo llena las columnas de peso y dimensiones de los productos que quieras actualizar, el resto puede quedar vacío</li>
                            <li>No borres ni cambies la columna sku: es la que identifica el producto a actualizar</li>
                            <li>Solo se actualizan productos que ya existen; si un SKU no se encuentra, esa fila sale como error y no se crea nada nuevo</li>
                        </ul>
                    </div>

                    {!success && !uploadStats && (
                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                                Seleccionar archivo
                            </label>
                            <input
                                type="file"
                                accept=".xlsx,.xls"
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

                    {(error || uploadStats) && (
                        <div className={`p-4 rounded-lg border ${error ? 'bg-red-50 border-red-200 text-red-700' : 'bg-green-50 border-green-200 text-green-700'}`}>
                            {error && <p className="font-bold mb-2">{error}</p>}
                            {success && <p className="font-bold mb-2">{success}</p>}

                            {uploadStats && (
                                <div className="mt-2 text-sm space-y-1">
                                    <p>Total de filas con datos: <strong>{uploadStats.total}</strong></p>
                                    <p>Actualizados: <strong className="text-green-600">{uploadStats.success}</strong></p>
                                    <p>Fallidos: <strong className="text-red-600">{uploadStats.failed}</strong></p>

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
                                    {loading ? 'Procesando...' : 'Cargar Dimensiones'}
                                </Button>
                            </>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}
