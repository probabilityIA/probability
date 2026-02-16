'use client';

import { useState } from 'react';
import { Button, Input } from '@/shared/ui';

interface MassOrderUploadModalProps {
    isOpen: boolean;
    onClose: () => void;
    onUploadComplete?: (count: number) => void;
}

export default function MassOrderUploadModal({ isOpen, onClose, onUploadComplete }: MassOrderUploadModalProps) {
    const [file, setFile] = useState<File | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [uploadStats, setUploadStats] = useState<{ total: number; success: number; failed: number } | null>(null);

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

        setLoading(true);
        setError(null);
        setSuccess(null);
        setUploadStats(null);

        try {
            const formData = new FormData();
            formData.append('file', file);

            const response = await fetch('/api/v1/orders/upload-bulk', {
                method: 'POST',
                body: formData,
            });

            const result = await response.json();

            if (response.ok && result.success) {
                setSuccess(`¡Carga exitosa! ${result.data.success_count} órdenes creadas.`);
                setUploadStats({
                    total: result.data.total_rows,
                    success: result.data.success_count,
                    failed: result.data.failed_count,
                });
                if (onUploadComplete) {
                    onUploadComplete(result.data.success_count);
                }
            } else {
                setError(result.message || 'Error al procesar el archivo');
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar el archivo');
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
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full p-6">
                <div className="flex justify-between items-center mb-6">
                    <h2 className="text-2xl font-bold text-gray-800">Carga Masiva de Órdenes</h2>
                    <button
                        onClick={handleClose}
                        className="text-gray-500 hover:text-gray-700 text-2xl"
                    >
                        ×
                    </button>
                </div>

                <div className="space-y-6">
                    {/* Instructions */}
                    <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                        <h3 className="font-semibold text-blue-800 mb-2">Instrucciones</h3>
                        <ul className="text-sm text-blue-700 space-y-1 list-disc list-inside">
                            <li>El archivo debe ser CSV o Excel (.xlsx, .xls)</li>
                            <li>La primera fila debe contener los encabezados de columna</li>
                            <li>Columnas requeridas: order_number, customer_name, customer_email, customer_phone, shipping_street, shipping_city, shipping_state, total_amount</li>
                            <li>Columnas opcionales: weight, height, width, length, platform</li>
                        </ul>
                    </div>

                    {/* File Input */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
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
                            <p className="mt-2 text-sm text-gray-600">
                                Archivo seleccionado: <strong>{file.name}</strong>
                            </p>
                        )}
                    </div>

                    {/* Error */}
                    {error && (
                        <div className="p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
                            {error}
                        </div>
                    )}

                    {/* Success */}
                    {success && (
                        <div className="p-4 bg-green-50 border border-green-200 rounded-lg text-green-700">
                            {success}
                            {uploadStats && (
                                <div className="mt-2 text-sm">
                                    <p>Total de filas: {uploadStats.total}</p>
                                    <p>Exitosas: {uploadStats.success}</p>
                                    <p>Fallidas: {uploadStats.failed}</p>
                                </div>
                            )}
                        </div>
                    )}

                    {/* Actions */}
                    <div className="flex justify-end space-x-3">
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
                    </div>
                </div>
            </div>
        </div>
    );
}
