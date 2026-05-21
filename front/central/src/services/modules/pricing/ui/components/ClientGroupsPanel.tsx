'use client';

import { useCallback, useEffect, useState } from 'react';
import { ClientGroup, ClientSummary, SaveClientGroupInput } from '../../domain/types';
import {
    saveClientGroupAction,
    deleteClientGroupAction,
    listGroupMembersAction,
    addGroupMembersAction,
    removeGroupMemberAction,
    listAvailableClientsAction,
} from '../../infra/actions';

interface ClientGroupsPanelProps {
    businessId?: number;
    groups: ClientGroup[];
    loading: boolean;
    onGroupsChanged: () => void | Promise<void>;
}

export const GROUP_COLORS = ['#ef4444', '#f59e0b', '#10b981', '#3b82f6', '#8b5cf6', '#ec4899', '#14b8a6', '#6b7280'];

const emptyForm: SaveClientGroupInput = { name: '', description: '', color: GROUP_COLORS[3], is_active: true };

export function ClientGroupsPanel({ businessId, groups, loading, onGroupsChanged }: ClientGroupsPanelProps) {
    const [selectedGroupId, setSelectedGroupId] = useState<number | null>(null);
    const [form, setForm] = useState<SaveClientGroupInput | null>(null);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    const [members, setMembers] = useState<ClientSummary[]>([]);
    const [membersLoading, setMembersLoading] = useState(false);

    const [addOpen, setAddOpen] = useState(false);
    const [clientSearch, setClientSearch] = useState('');
    const [available, setAvailable] = useState<ClientSummary[]>([]);
    const [availableLoading, setAvailableLoading] = useState(false);
    const [picked, setPicked] = useState<Set<number>>(new Set());

    const selectedGroup = groups.find((g) => g.id === selectedGroupId) || null;

    const loadMembers = useCallback(async (groupId: number) => {
        setMembersLoading(true);
        const result = await listGroupMembersAction(businessId, groupId, '', 1);
        setMembers(result.data);
        setMembersLoading(false);
    }, [businessId]);

    useEffect(() => {
        if (selectedGroupId) {
            loadMembers(selectedGroupId);
            setAddOpen(false);
        } else {
            setMembers([]);
        }
    }, [selectedGroupId, loadMembers]);

    const loadAvailable = useCallback(async (search: string) => {
        setAvailableLoading(true);
        const result = await listAvailableClientsAction(businessId, search, false, 1);
        setAvailable(result.data);
        setAvailableLoading(false);
    }, [businessId]);

    useEffect(() => {
        if (!addOpen) return;
        const handle = setTimeout(() => loadAvailable(clientSearch), 300);
        return () => clearTimeout(handle);
    }, [addOpen, clientSearch, loadAvailable]);

    const startCreate = () => {
        setForm({ ...emptyForm });
        setError('');
    };

    const startEdit = (group: ClientGroup) => {
        setForm({
            id: group.id,
            name: group.name,
            description: group.description,
            color: group.color || GROUP_COLORS[3],
            is_active: group.is_active,
        });
        setError('');
    };

    const handleSave = async () => {
        if (!form) return;
        if (!form.name.trim()) {
            setError('El nombre del grupo es obligatorio');
            return;
        }
        setSaving(true);
        setError('');
        const result = await saveClientGroupAction(businessId, form);
        setSaving(false);
        if (!result.success) {
            setError(result.message || 'No se pudo guardar el grupo');
            return;
        }
        setForm(null);
        await onGroupsChanged();
    };

    const handleDelete = async (group: ClientGroup) => {
        if (!window.confirm(`Eliminar el grupo "${group.name}"? Se borraran sus precios personalizados.`)) return;
        const result = await deleteClientGroupAction(businessId, group.id);
        if (!result.success) {
            setError(result.message || 'No se pudo eliminar el grupo');
            return;
        }
        if (selectedGroupId === group.id) setSelectedGroupId(null);
        await onGroupsChanged();
    };

    const togglePicked = (id: number) => {
        setPicked((prev) => {
            const next = new Set(prev);
            if (next.has(id)) next.delete(id);
            else next.add(id);
            return next;
        });
    };

    const handleAddClients = async () => {
        if (!selectedGroupId || picked.size === 0) return;
        const result = await addGroupMembersAction(businessId, selectedGroupId, [...picked]);
        if (!result.success) {
            setError(result.message || 'No se pudieron agregar los clientes');
            return;
        }
        setPicked(new Set());
        setAddOpen(false);
        await loadMembers(selectedGroupId);
        await onGroupsChanged();
    };

    const handleRemoveMember = async (clientId: number) => {
        if (!selectedGroupId) return;
        const result = await removeGroupMemberAction(businessId, selectedGroupId, clientId);
        if (!result.success) {
            setError(result.message || 'No se pudo quitar el cliente');
            return;
        }
        await loadMembers(selectedGroupId);
        await onGroupsChanged();
    };

    return (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <div className="border border-gray-200 dark:border-gray-700 rounded-xl p-4 space-y-3">
                <div className="flex items-center justify-between">
                    <h3 className="font-bold text-gray-900 dark:text-white">Tipos de cliente</h3>
                    <button
                        onClick={startCreate}
                        className="px-3 py-1.5 btn-business-primary text-white text-sm font-bold rounded-lg"
                    >
                        + Nuevo grupo
                    </button>
                </div>

                {error && <p className="text-sm text-red-600">{error}</p>}

                {form && (
                    <div className="border border-business-primary/40 rounded-lg p-3 space-y-2 bg-gray-50 dark:bg-gray-800">
                        <input
                            type="text"
                            placeholder="Nombre del grupo (ej: Mayorista)"
                            value={form.name}
                            onChange={(e) => setForm({ ...form, name: e.target.value })}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                        />
                        <input
                            type="text"
                            placeholder="Descripcion (opcional)"
                            value={form.description}
                            onChange={(e) => setForm({ ...form, description: e.target.value })}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                        />
                        <div>
                            <p className="text-xs font-bold text-gray-600 dark:text-gray-300 mb-1.5">Color del grupo</p>
                            <div className="flex gap-1.5">
                                {GROUP_COLORS.map((c) => (
                                    <button
                                        key={c}
                                        type="button"
                                        onClick={() => setForm({ ...form, color: c })}
                                        className={`w-6 h-6 rounded-full transition-transform ${form.color === c ? 'ring-2 ring-offset-1 ring-gray-700 dark:ring-white scale-110' : ''}`}
                                        style={{ backgroundColor: c }}
                                        aria-label={`Color ${c}`}
                                    />
                                ))}
                            </div>
                        </div>
                        <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200">
                            <input
                                type="checkbox"
                                checked={form.is_active}
                                onChange={(e) => setForm({ ...form, is_active: e.target.checked })}
                            />
                            Activo
                        </label>
                        <div className="flex gap-2">
                            <button
                                onClick={handleSave}
                                disabled={saving}
                                className="px-3 py-1.5 btn-business-primary text-white text-sm font-bold rounded-lg disabled:opacity-50"
                            >
                                {saving ? 'Guardando...' : 'Guardar'}
                            </button>
                            <button
                                onClick={() => setForm(null)}
                                className="px-3 py-1.5 border border-gray-300 dark:border-gray-600 text-sm font-bold rounded-lg text-gray-700 dark:text-gray-200"
                            >
                                Cancelar
                            </button>
                        </div>
                    </div>
                )}

                {loading ? (
                    <p className="text-sm text-gray-500">Cargando grupos...</p>
                ) : groups.length === 0 ? (
                    <p className="text-sm text-gray-500">Aun no hay grupos. Crea el primero.</p>
                ) : (
                    <div className="space-y-2 max-h-80 overflow-y-auto">
                        {groups.map((group) => (
                            <div
                                key={group.id}
                                onClick={() => setSelectedGroupId(group.id)}
                                className={`flex items-center justify-between p-3 rounded-lg cursor-pointer border transition-all ${
                                    selectedGroupId === group.id
                                        ? 'border-business-primary bg-business-primary/5'
                                        : 'border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800'
                                }`}
                            >
                                <div className="flex items-center gap-2.5">
                                    <span
                                        className="w-3.5 h-3.5 rounded-full flex-shrink-0"
                                        style={{ backgroundColor: group.color || '#6b7280' }}
                                    />
                                    <div>
                                    <p className="font-semibold text-gray-900 dark:text-white text-sm">
                                        {group.name}
                                        {!group.is_active && <span className="ml-2 text-xs text-gray-400">(inactivo)</span>}
                                    </p>
                                    <p className="text-xs text-gray-500">{group.member_count} clientes</p>
                                    </div>
                                </div>
                                <div className="flex gap-1" onClick={(e) => e.stopPropagation()}>
                                    <button
                                        onClick={() => startEdit(group)}
                                        className="px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded text-gray-600 dark:text-gray-300"
                                    >
                                        Editar
                                    </button>
                                    <button
                                        onClick={() => handleDelete(group)}
                                        className="px-2 py-1 text-xs border border-red-300 text-red-600 rounded"
                                    >
                                        Eliminar
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            <div className="border border-gray-200 dark:border-gray-700 rounded-xl p-4 space-y-3">
                {!selectedGroup ? (
                    <p className="text-sm text-gray-500">Selecciona un grupo para administrar sus clientes.</p>
                ) : (
                    <>
                        <div className="flex items-center justify-between">
                            <h3 className="font-bold text-gray-900 dark:text-white">
                                Clientes de {selectedGroup.name}
                            </h3>
                            <button
                                onClick={() => { setAddOpen((v) => !v); setPicked(new Set()); }}
                                className="px-3 py-1.5 btn-business-primary text-white text-sm font-bold rounded-lg"
                            >
                                {addOpen ? 'Cerrar' : '+ Agregar clientes'}
                            </button>
                        </div>

                        {addOpen && (
                            <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-3 space-y-2 bg-gray-50 dark:bg-gray-800">
                                <input
                                    type="text"
                                    placeholder="Buscar cliente por nombre, email o documento"
                                    value={clientSearch}
                                    onChange={(e) => setClientSearch(e.target.value)}
                                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                                />
                                <div className="max-h-52 overflow-y-auto space-y-1">
                                    {availableLoading ? (
                                        <p className="text-sm text-gray-500">Buscando...</p>
                                    ) : available.length === 0 ? (
                                        <p className="text-sm text-gray-500">Sin resultados.</p>
                                    ) : (
                                        available.map((client) => (
                                            <label
                                                key={client.id}
                                                className="flex items-center gap-2 p-2 rounded hover:bg-white dark:hover:bg-gray-700 cursor-pointer"
                                            >
                                                <input
                                                    type="checkbox"
                                                    checked={picked.has(client.id)}
                                                    onChange={() => togglePicked(client.id)}
                                                />
                                                <span className="text-sm text-gray-900 dark:text-white">{client.name}</span>
                                                {client.group_id && client.group_id !== selectedGroup.id && (
                                                    <span className="text-xs text-amber-600">
                                                        (en {client.group_name})
                                                    </span>
                                                )}
                                            </label>
                                        ))
                                    )}
                                </div>
                                <button
                                    onClick={handleAddClients}
                                    disabled={picked.size === 0}
                                    className="px-3 py-1.5 btn-business-primary text-white text-sm font-bold rounded-lg disabled:opacity-50"
                                >
                                    Agregar {picked.size > 0 ? `(${picked.size})` : ''}
                                </button>
                            </div>
                        )}

                        {membersLoading ? (
                            <p className="text-sm text-gray-500">Cargando clientes...</p>
                        ) : members.length === 0 ? (
                            <p className="text-sm text-gray-500">Este grupo aun no tiene clientes.</p>
                        ) : (
                            <div className="space-y-1 max-h-72 overflow-y-auto">
                                {members.map((client) => (
                                    <div
                                        key={client.id}
                                        className="flex items-center justify-between p-2 rounded border border-gray-100 dark:border-gray-700"
                                    >
                                        <div>
                                            <p className="text-sm text-gray-900 dark:text-white">{client.name}</p>
                                            <p className="text-xs text-gray-500">{client.email || client.phone || client.dni}</p>
                                        </div>
                                        <button
                                            onClick={() => handleRemoveMember(client.id)}
                                            className="px-2 py-1 text-xs border border-red-300 text-red-600 rounded"
                                        >
                                            Quitar
                                        </button>
                                    </div>
                                ))}
                            </div>
                        )}
                    </>
                )}
            </div>
        </div>
    );
}
