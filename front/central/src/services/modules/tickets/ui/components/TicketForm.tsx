'use client';

import { useState } from 'react';
import { Button, Input } from '@/shared/ui';
import {
    CreateTicketDTO,
    TICKET_TYPES,
    TICKET_PRIORITIES,
    TYPE_META,
    PRIORITY_META,
    TicketType,
    TicketPriority,
} from '../../domain/types';

interface Props {
    isSuperAdmin: boolean;
    selectedBusinessId?: number | null;
    onSubmit: (data: CreateTicketDTO) => Promise<void>;
    onCancel: () => void;
    submitting?: boolean;
}

export default function TicketForm({ isSuperAdmin, selectedBusinessId, onSubmit, onCancel, submitting }: Props) {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [type, setType] = useState<TicketType>('support');
    const [priority, setPriority] = useState<TicketPriority>('medium');
    const [category, setCategory] = useState('');
    const [dueDate, setDueDate] = useState('');
    const [error, setError] = useState<string | null>(null);

    const submit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        if (!title.trim() || !description.trim()) {
            setError('Titulo y descripcion son obligatorios');
            return;
        }
        try {
            await onSubmit({
                title: title.trim(),
                description: description.trim(),
                type,
                priority,
                category: category.trim() || undefined,
                due_date: dueDate || undefined,
                business_id: isSuperAdmin ? (selectedBusinessId ?? null) : undefined,
                source: isSuperAdmin ? 'internal' : 'business',
            });
        } catch (err: any) {
            setError(err?.message || 'Error al crear ticket');
        }
    };

    return (
        <form onSubmit={submit} className="space-y-4">
            <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Titulo *</label>
                <Input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Resumen breve" />
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Descripcion *</label>
                <textarea
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    rows={5}
                    className="block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Detalla el problema, mejora o solicitud"
                />
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Tipo</label>
                    <select
                        value={type}
                        onChange={(e) => setType(e.target.value as TicketType)}
                        className="block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 px-3 py-2 text-sm"
                    >
                        {TICKET_TYPES.map((t) => (
                            <option key={t} value={t}>{TYPE_META[t].label}</option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Prioridad</label>
                    <select
                        value={priority}
                        onChange={(e) => setPriority(e.target.value as TicketPriority)}
                        className="block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 px-3 py-2 text-sm"
                    >
                        {TICKET_PRIORITIES.map((p) => (
                            <option key={p} value={p}>{PRIORITY_META[p].label}</option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Categoria</label>
                    <Input value={category} onChange={(e) => setCategory(e.target.value)} placeholder="Ej: envios, facturacion, frontend" />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Fecha objetivo</label>
                    <Input type="date" value={dueDate} onChange={(e) => setDueDate(e.target.value)} />
                </div>
            </div>

            {error && <div className="text-sm text-red-600 dark:text-red-400">{error}</div>}

            <div className="flex justify-end gap-2 pt-2">
                <Button type="button" variant="outline" onClick={onCancel} disabled={submitting}>Cancelar</Button>
                <Button type="submit" variant="primary" disabled={submitting}>{submitting ? 'Creando...' : 'Crear ticket'}</Button>
            </div>
        </form>
    );
}
