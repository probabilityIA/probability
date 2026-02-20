'use client';

import React, { useState, useEffect } from 'react';
import { Button } from '@/shared/ui/button';
import { Alert } from '@/shared/ui/alert';
import { Spinner } from '@/shared/ui/spinner';
import { createPermissionsBulkAction } from '../../infra/actions';
import { BulkCreateResult, CreatePermissionDTO } from '../../domain/types';
import { getResourcesAction } from '@/services/auth/resources/infra/actions';
import { getActionsAction } from '@/services/auth/actions/infra/actions';
import { Resource } from '@/services/auth/resources/domain/types';
import { Action } from '@/services/auth/actions/domain/types';
import { SCOPE_OPTIONS } from '../hooks/usePermissionForm';

interface ActionRow {
    action: Action;
    checked: boolean;
    name: string;
}

interface BulkPermissionFormProps {
    onSuccess: () => void;
    onCancel: () => void;
}

export const BulkPermissionForm: React.FC<BulkPermissionFormProps> = ({ onSuccess, onCancel }) => {
    const [loadingData, setLoadingData] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const [resources, setResources] = useState<Resource[]>([]);
    const [allActions, setAllActions] = useState<Action[]>([]);

    const [selectedResourceId, setSelectedResourceId] = useState<string>('');
    const [selectedScopeId, setSelectedScopeId] = useState<string>('2'); // Business por defecto
    const [actionRows, setActionRows] = useState<ActionRow[]>([]);

    const [results, setResults] = useState<BulkCreateResult[] | null>(null);
    const [resultMessage, setResultMessage] = useState('');

    // Cargar recursos y acciones
    useEffect(() => {
        const load = async () => {
            setLoadingData(true);
            try {
                const [resRes, actRes] = await Promise.all([
                    getResourcesAction({ page_size: 100 }),
                    getActionsAction({ page_size: 100 }),
                ]);
                if (resRes.success && resRes.data?.resources) setResources(resRes.data.resources);
                if (actRes.success && actRes.data?.actions) setAllActions(actRes.data.actions);
            } catch (e) {
                console.error(e);
            } finally {
                setLoadingData(false);
            }
        };
        load();
    }, []);

    // Cuando cambia el recurso, regenerar nombres
    const handleResourceChange = (resourceId: string) => {
        setSelectedResourceId(resourceId);
        const resource = resources.find(r => r.id === Number(resourceId));
        setActionRows(
            allActions.map(action => ({
                action,
                checked: false,
                name: resource ? `${action.name} ${resource.name}` : action.name,
            }))
        );
    };

    const handleToggleAction = (index: number) => {
        setActionRows(prev =>
            prev.map((row, i) => i === index ? { ...row, checked: !row.checked } : row)
        );
    };

    const handleToggleAll = () => {
        const allChecked = actionRows.every(r => r.checked);
        setActionRows(prev => prev.map(row => ({ ...row, checked: !allChecked })));
    };

    const handleNameChange = (index: number, value: string) => {
        setActionRows(prev =>
            prev.map((row, i) => i === index ? { ...row, name: value } : row)
        );
    };

    const selectedCount = actionRows.filter(r => r.checked).length;

    const handleSubmit = async () => {
        if (!selectedResourceId) { setError('Selecciona un recurso'); return; }
        if (!selectedScopeId) { setError('Selecciona un scope'); return; }
        if (selectedCount === 0) { setError('Selecciona al menos una acción'); return; }

        const permissions: CreatePermissionDTO[] = actionRows
            .filter(r => r.checked)
            .map(r => ({
                name: r.name,
                resource_id: Number(selectedResourceId),
                action_id: r.action.id,
                scope_id: Number(selectedScopeId),
                business_type_id: 1,
            }));

        setSubmitting(true);
        setError(null);
        try {
            const res = await createPermissionsBulkAction(permissions);
            setResults(res.results);
            setResultMessage(res.message);
        } catch (e: any) {
            setError(e.message || 'Error al crear permisos');
        } finally {
            setSubmitting(false);
        }
    };

    // Vista de resultados
    if (results) {
        const successCount = results.filter(r => r.success).length;
        return (
            <div className="space-y-4">
                <p className="text-sm font-medium text-gray-700">{resultMessage}</p>
                <div className="divide-y divide-gray-200 border rounded-lg overflow-hidden max-h-80 overflow-y-auto">
                    {results.map((r, i) => (
                        <div key={i} className="flex items-center gap-3 px-4 py-3">
                            {r.success ? (
                                <span className="w-5 h-5 flex items-center justify-center rounded-full bg-green-100 text-green-600">
                                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                                    </svg>
                                </span>
                            ) : (
                                <span className="w-5 h-5 flex items-center justify-center rounded-full bg-red-100 text-red-600">
                                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M6 18L18 6M6 6l12 12" />
                                    </svg>
                                </span>
                            )}
                            <div className="flex-1 min-w-0">
                                <p className="text-sm font-medium text-gray-900 truncate">{r.name}</p>
                                {!r.success && <p className="text-xs text-red-600 truncate">{r.error}</p>}
                            </div>
                        </div>
                    ))}
                </div>
                <div className="flex justify-end gap-2 pt-2 border-t">
                    {successCount > 0 && (
                        <Button variant="primary" onClick={onSuccess}>Cerrar</Button>
                    )}
                    {successCount === 0 && (
                        <Button variant="secondary" onClick={() => setResults(null)}>Volver</Button>
                    )}
                </div>
            </div>
        );
    }

    if (loadingData) {
        return <div className="flex justify-center py-8"><Spinner /></div>;
    }

    return (
        <div className="space-y-5">
            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            {/* Recurso + Scope */}
            <div className="grid grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Recurso *</label>
                    <select
                        value={selectedResourceId}
                        onChange={e => handleResourceChange(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    >
                        <option value="">Seleccionar recurso...</option>
                        {resources.map(r => (
                            <option key={r.id} value={r.id}>{r.name}</option>
                        ))}
                    </select>
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Scope *</label>
                    <select
                        value={selectedScopeId}
                        onChange={e => setSelectedScopeId(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    >
                        {SCOPE_OPTIONS.map(s => (
                            <option key={s.id} value={s.id}>{s.name}</option>
                        ))}
                    </select>
                </div>
            </div>

            {/* Lista de acciones */}
            {actionRows.length > 0 ? (
                <div>
                    <div className="flex items-center justify-between mb-2">
                        <label className="block text-sm font-medium text-gray-700">
                            Acciones — selecciona y edita el nombre del permiso
                        </label>
                        <button
                            type="button"
                            onClick={handleToggleAll}
                            className="text-xs text-blue-600 hover:underline"
                        >
                            {actionRows.every(r => r.checked) ? 'Deseleccionar todo' : 'Seleccionar todo'}
                        </button>
                    </div>
                    <div className="border rounded-lg divide-y divide-gray-200 max-h-72 overflow-y-auto">
                        {actionRows.map((row, i) => (
                            <div key={row.action.id} className={`flex items-center gap-3 px-4 py-3 ${row.checked ? 'bg-blue-50' : 'bg-white'}`}>
                                <input
                                    type="checkbox"
                                    checked={row.checked}
                                    onChange={() => handleToggleAction(i)}
                                    className="h-4 w-4 text-blue-600 rounded border-gray-300"
                                />
                                <span className="w-20 text-xs font-medium text-gray-500 shrink-0">{row.action.name}</span>
                                <input
                                    type="text"
                                    value={row.name}
                                    onChange={e => handleNameChange(i, e.target.value)}
                                    disabled={!row.checked}
                                    className="flex-1 px-2 py-1 text-sm border border-gray-300 rounded focus:ring-1 focus:ring-blue-500 disabled:bg-gray-50 disabled:text-gray-400"
                                />
                            </div>
                        ))}
                    </div>
                    <p className="mt-1 text-xs text-gray-500">{selectedCount} acción(es) seleccionada(s)</p>
                </div>
            ) : (
                <p className="text-sm text-gray-400 text-center py-4">
                    {selectedResourceId ? 'No hay acciones disponibles' : 'Selecciona un recurso para ver las acciones'}
                </p>
            )}

            <div className="flex justify-end gap-2 pt-2 border-t">
                <Button variant="secondary" onClick={onCancel} disabled={submitting}>Cancelar</Button>
                <Button
                    variant="primary"
                    onClick={handleSubmit}
                    disabled={submitting || selectedCount === 0}
                >
                    {submitting ? 'Creando...' : `Crear ${selectedCount} permiso${selectedCount !== 1 ? 's' : ''}`}
                </Button>
            </div>
        </div>
    );
};
