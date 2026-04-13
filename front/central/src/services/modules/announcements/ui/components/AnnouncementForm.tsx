'use client';

import { useState, useEffect } from 'react';
import { TrashIcon, PlusIcon } from '@heroicons/react/24/outline';
import { AnnouncementInfo, CreateAnnouncementDTO, UpdateAnnouncementDTO, AnnouncementCategory, CreateLinkDTO, DisplayType, FrequencyType } from '../../domain/types';
import { createAnnouncementAction, updateAnnouncementAction, listCategoriesAction, uploadImageAction, deleteImageAction } from '../../infra/actions';
import { Button, Alert, Input } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';
import ImageUploader from './ImageUploader';
import BusinessTargetSelector from './BusinessTargetSelector';

interface ImageItem {
    id?: number;
    image_url?: string;
    file?: File;
    preview: string;
    sort_order: number;
}

interface AnnouncementFormProps {
    announcement?: AnnouncementInfo;
    onSuccess: () => void;
    onCancel: () => void;
}

const displayTypeOptions: { value: DisplayType; label: string }[] = [
    { value: 'modal_image', label: 'Modal con imagen' },
    { value: 'modal_text', label: 'Modal de texto' },
    { value: 'ticker', label: 'Ticker (barra superior)' },
];

const frequencyOptions: { value: FrequencyType; label: string }[] = [
    { value: 'once', label: 'Una sola vez' },
    { value: 'daily', label: 'Diario' },
    { value: 'always', label: 'Siempre' },
    { value: 'requires_acceptance', label: 'Requiere aceptacion' },
];

export default function AnnouncementForm({ announcement, onSuccess, onCancel }: AnnouncementFormProps) {
    const [categories, setCategories] = useState<AnnouncementCategory[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const [title, setTitle] = useState(announcement?.title || '');
    const [message, setMessage] = useState(announcement?.message || '');
    const [categoryId, setCategoryId] = useState(announcement?.category_id || 0);
    const [displayType, setDisplayType] = useState<DisplayType>(announcement?.display_type || 'modal_text');
    const [frequencyType, setFrequencyType] = useState<FrequencyType>(announcement?.frequency_type || 'once');
    const [priority, setPriority] = useState(announcement?.priority || 0);
    const [isGlobal, setIsGlobal] = useState(announcement?.is_global ?? true);
    const [startsAt, setStartsAt] = useState(announcement?.starts_at?.slice(0, 16) || '');
    const [endsAt, setEndsAt] = useState(announcement?.ends_at?.slice(0, 16) || '');
    const [links, setLinks] = useState<CreateLinkDTO[]>(
        announcement?.links?.map(l => ({ label: l.label, url: l.url, sort_order: l.sort_order })) || []
    );
    const [targetIds, setTargetIds] = useState<number[]>(
        announcement?.targets?.map(t => t.business_id) || []
    );

    const [images, setImages] = useState<ImageItem[]>(
        (announcement?.images || []).map(img => ({
            id: img.id,
            image_url: img.image_url,
            preview: img.image_url,
            sort_order: img.sort_order,
        }))
    );
    const [removedImageIds, setRemovedImageIds] = useState<number[]>([]);

    const initialImageIds = (announcement?.images || []).map(img => img.id);

    useEffect(() => {
        listCategoriesAction()
            .then(setCategories)
            .catch(() => {});
    }, []);

    const handleImagesChange = (updated: ImageItem[]) => {
        const currentIds = updated.filter(i => i.id).map(i => i.id!);
        const newlyRemoved = initialImageIds.filter(id => !currentIds.includes(id));
        setRemovedImageIds(prev => {
            const combined = new Set([...prev, ...newlyRemoved]);
            return Array.from(combined);
        });
        setImages(updated);
    };

    const addLink = () => {
        setLinks([...links, { label: '', url: '', sort_order: links.length }]);
    };

    const removeLink = (index: number) => {
        setLinks(links.filter((_, i) => i !== index));
    };

    const updateLink = (index: number, field: keyof CreateLinkDTO, value: string | number) => {
        const updated = [...links];
        updated[index] = { ...updated[index], [field]: value };
        setLinks(updated);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        const data: CreateAnnouncementDTO = {
            category_id: categoryId,
            title,
            message,
            display_type: displayType,
            frequency_type: frequencyType,
            priority,
            is_global: isGlobal,
            starts_at: startsAt ? new Date(startsAt).toISOString() : undefined,
            ends_at: endsAt ? new Date(endsAt).toISOString() : undefined,
            links,
            target_ids: isGlobal ? [] : targetIds,
        };

        try {
            let savedAnnouncement: AnnouncementInfo;

            if (announcement) {
                savedAnnouncement = await updateAnnouncementAction(announcement.id, data as UpdateAnnouncementDTO);
            } else {
                savedAnnouncement = await createAnnouncementAction(data);
            }

            const announcementId = savedAnnouncement.id;

            for (const imageId of removedImageIds) {
                try {
                    await deleteImageAction(announcementId, imageId);
                } catch {}
            }

            const newImages = images.filter(img => img.file);
            for (const img of newImages) {
                const fd = new FormData();
                fd.append('image', img.file!);
                fd.append('sort_order', String(img.sort_order));
                await uploadImageAction(announcementId, fd);
            }

            setSuccess(announcement ? 'Anuncio actualizado' : 'Anuncio creado');
            setTimeout(() => onSuccess(), 800);
        } catch (err: any) {
            setError(getActionError(err, 'Error al guardar el anuncio'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-6">
            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}
            {success && <Alert type="success" onClose={() => setSuccess(null)}>{success}</Alert>}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Titulo <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        placeholder="Titulo del anuncio"
                        required
                        maxLength={255}
                    />
                </div>

                <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Mensaje
                    </label>
                    <textarea
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        placeholder="Contenido del anuncio..."
                        rows={4}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm text-gray-900 dark:text-white bg-white dark:bg-gray-700 placeholder-gray-500 dark:placeholder-white focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Categoria <span className="text-red-500">*</span>
                    </label>
                    <select
                        value={categoryId}
                        onChange={(e) => setCategoryId(Number(e.target.value))}
                        required
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm text-gray-900 dark:text-white bg-white dark:bg-gray-700 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    >
                        <option value={0} disabled>Seleccionar...</option>
                        {categories.map(cat => (
                            <option key={cat.id} value={cat.id}>{cat.name}</option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Tipo de display <span className="text-red-500">*</span>
                    </label>
                    <select
                        value={displayType}
                        onChange={(e) => setDisplayType(e.target.value as DisplayType)}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm text-gray-900 dark:text-white bg-white dark:bg-gray-700 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    >
                        {displayTypeOptions.map(opt => (
                            <option key={opt.value} value={opt.value}>{opt.label}</option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Frecuencia <span className="text-red-500">*</span>
                    </label>
                    <select
                        value={frequencyType}
                        onChange={(e) => setFrequencyType(e.target.value as FrequencyType)}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm text-gray-900 dark:text-white bg-white dark:bg-gray-700 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    >
                        {frequencyOptions.map(opt => (
                            <option key={opt.value} value={opt.value}>{opt.label}</option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Prioridad
                    </label>
                    <Input
                        type="number"
                        value={String(priority)}
                        onChange={(e) => setPriority(Number(e.target.value))}
                        placeholder="0"
                        min={0}
                        max={100}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Inicio
                    </label>
                    <Input
                        type="datetime-local"
                        value={startsAt}
                        onChange={(e) => setStartsAt(e.target.value)}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Fin
                    </label>
                    <Input
                        type="datetime-local"
                        value={endsAt}
                        onChange={(e) => setEndsAt(e.target.value)}
                    />
                </div>

                <div className="md:col-span-2">
                    <label className="flex items-center gap-2 cursor-pointer">
                        <input
                            type="checkbox"
                            checked={isGlobal}
                            onChange={(e) => setIsGlobal(e.target.checked)}
                            className="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500"
                        />
                        <span className="text-sm font-medium text-gray-700 dark:text-gray-200">
                            Anuncio global (visible para todos los negocios)
                        </span>
                    </label>
                </div>
            </div>

            {!isGlobal && (
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                        Negocios objetivo
                    </label>
                    <BusinessTargetSelector
                        selectedIds={targetIds}
                        onChange={setTargetIds}
                    />
                </div>
            )}

            <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                    Imagenes
                </label>
                <ImageUploader
                    images={images}
                    onChange={handleImagesChange}
                />
            </div>

            <div>
                <div className="flex items-center justify-between mb-2">
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200">
                        Links
                    </label>
                    <button
                        type="button"
                        onClick={addLink}
                        className="inline-flex items-center gap-1 text-sm text-purple-600 hover:text-purple-700"
                    >
                        <PlusIcon className="w-4 h-4" />
                        Agregar link
                    </button>
                </div>
                {links.map((link, index) => (
                    <div key={index} className="flex gap-2 mb-2">
                        <Input
                            type="text"
                            value={link.label}
                            onChange={(e) => updateLink(index, 'label', e.target.value)}
                            placeholder="Etiqueta"
                            className="flex-1"
                        />
                        <Input
                            type="url"
                            value={link.url}
                            onChange={(e) => updateLink(index, 'url', e.target.value)}
                            placeholder="https://..."
                            className="flex-1"
                        />
                        <button
                            type="button"
                            onClick={() => removeLink(index)}
                            className="p-2 text-red-500 hover:text-red-700"
                        >
                            <TrashIcon className="w-4 h-4" />
                        </button>
                    </div>
                ))}
            </div>

            <div className="flex justify-end gap-3 pt-4 border-t">
                <Button type="button" variant="outline" onClick={onCancel} disabled={loading}>
                    Cancelar
                </Button>
                <Button type="submit" variant="primary" disabled={loading}>
                    {loading ? 'Guardando...' : announcement ? 'Actualizar' : 'Crear anuncio'}
                </Button>
            </div>
        </form>
    );
}
