'use client';

import { Modal } from '@/shared/ui';
import { Download, ExternalLink } from 'lucide-react';

interface Props {
    isOpen: boolean;
    onClose: () => void;
    shipmentId: number | null;
    orderLabel?: string;
    carrierLabel?: string;
}

export function GuidePreviewModal({ isOpen, onClose, shipmentId, orderLabel, carrierLabel }: Props) {
    if (!shipmentId) return null;

    const previewUrl = `/internal/shipment-guide/${shipmentId}`;
    const downloadUrl = `/internal/shipment-guide/${shipmentId}?download=1`;

    return (
        <Modal isOpen={isOpen} onClose={onClose} size="4xl" title={`Guia de transporte${orderLabel ? ` - ${orderLabel}` : ''}`}>
            <div className="flex flex-col gap-3" style={{ height: '80vh' }}>
                <div className="flex items-center justify-between gap-2">
                    <span className="text-sm text-gray-500 dark:text-gray-400">
                        {carrierLabel ? `Transportadora: ${carrierLabel}` : 'Vista previa del PDF'}
                    </span>
                    <div className="flex items-center gap-2">
                        <a
                            href={previewUrl}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-semibold text-gray-700 dark:text-gray-200 border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                        >
                            <ExternalLink size={13} /> Abrir en pestana
                        </a>
                        <a
                            href={downloadUrl}
                            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-semibold text-white bg-emerald-600 hover:bg-emerald-700 transition-colors"
                        >
                            <Download size={13} /> Descargar
                        </a>
                    </div>
                </div>
                <iframe
                    src={previewUrl}
                    title="Guia de transporte"
                    className="flex-1 w-full rounded-lg border border-gray-200 dark:border-gray-700 bg-white"
                />
            </div>
        </Modal>
    );
}
