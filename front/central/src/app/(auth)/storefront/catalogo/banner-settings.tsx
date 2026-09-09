'use client';

import { useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { PhotoIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { updateCatalogBannerAction, uploadCatalogBannerImageAction } from '@/services/modules/storefront/infra/actions';
import { CatalogBanner } from '@/services/modules/storefront/domain/types';
import { usePermissions } from '@/shared/contexts/permissions-context';

interface BannerSettingsProps {
    banner: CatalogBanner;
    businessId?: number;
}

export function BannerSettings({ banner, businessId }: BannerSettingsProps) {
    const router = useRouter();
    const { permissions } = usePermissions();
    const fileInputRef = useRef<HTMLInputElement>(null);
    const [open, setOpen] = useState(false);
    const [uploading, setUploading] = useState(false);
    const [togglingEnabled, setTogglingEnabled] = useState(false);
    const [error, setError] = useState<string | null>(null);

    if (permissions?.role_name === 'cliente_final') {
        return null;
    }

    const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        setUploading(true);
        setError(null);

        const formData = new FormData();
        formData.append('image', file);

        const result = await uploadCatalogBannerImageAction(formData, businessId);
        setUploading(false);

        if (!result.success) {
            setError(result.message || 'Error al subir la imagen');
            return;
        }
        router.refresh();
    };

    const handleToggle = async () => {
        setTogglingEnabled(true);
        setError(null);
        const result = await updateCatalogBannerAction(!banner.enabled, businessId);
        setTogglingEnabled(false);

        if (!result.success) {
            setError(result.message || 'Error al actualizar el banner');
            return;
        }
        router.refresh();
    };

    return (
        <div className="relative">
            <button
                type="button"
                onClick={() => setOpen(prev => !prev)}
                className="flex items-center gap-2 px-4 py-2 h-[42px] text-sm font-medium border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
                <PhotoIcon className="w-4 h-4" />
                Banner
            </button>

            {open && (
                <>
                    <div className="fixed inset-0 z-30" onClick={() => setOpen(false)} />
                    <div className="absolute right-0 mt-2 w-80 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-xl z-40 p-4">
                        <div className="flex items-center justify-between mb-3">
                            <h3 className="text-sm font-bold text-gray-900 dark:text-white">Banner del catalogo</h3>
                            <button onClick={() => setOpen(false)} className="text-gray-400 hover:text-gray-600">
                                <XMarkIcon className="w-4 h-4" />
                            </button>
                        </div>

                        <p className="text-xs text-gray-500 dark:text-gray-400 mb-3">
                            Muestra una imagen o tu logo arriba de los productos para tus clientes.
                        </p>

                        <div className="mb-3">
                            {banner.image_url ? (
                                <div className="relative w-full aspect-[3/1] rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900">
                                    <img src={banner.image_url} alt="Banner del catalogo" className="w-full h-full object-cover" />
                                </div>
                            ) : (
                                <div className="w-full aspect-[3/1] rounded-lg border border-dashed border-gray-300 dark:border-gray-600 flex items-center justify-center text-xs text-gray-400">
                                    Sin imagen
                                </div>
                            )}
                        </div>

                        <input
                            ref={fileInputRef}
                            type="file"
                            accept="image/png,image/jpeg,image/jpg,image/gif,image/webp"
                            onChange={handleFileChange}
                            className="hidden"
                        />
                        <button
                            type="button"
                            onClick={() => fileInputRef.current?.click()}
                            disabled={uploading}
                            className="w-full py-2 mb-3 border border-gray-300 dark:border-gray-600 rounded-lg text-sm font-medium text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 transition-colors"
                        >
                            {uploading ? 'Subiendo...' : banner.image_url ? 'Cambiar imagen' : 'Subir imagen'}
                        </button>

                        <label className="flex items-center justify-between">
                            <span className="text-sm text-gray-700 dark:text-gray-200">Mostrar banner a mis clientes</span>
                            <button
                                type="button"
                                onClick={handleToggle}
                                disabled={togglingEnabled || !banner.image_url}
                                role="switch"
                                aria-checked={banner.enabled}
                                className={`relative w-10 h-6 rounded-full transition-colors disabled:opacity-40 ${
                                    banner.enabled ? 'bg-indigo-600' : 'bg-gray-300 dark:bg-gray-600'
                                }`}
                            >
                                <span
                                    className={`absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full transition-transform ${
                                        banner.enabled ? 'translate-x-4' : ''
                                    }`}
                                />
                            </button>
                        </label>
                        {!banner.image_url && (
                            <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">Sube una imagen primero para poder activarlo</p>
                        )}

                        {error && <p className="text-xs text-red-500 mt-3">{error}</p>}
                    </div>
                </>
            )}
        </div>
    );
}
