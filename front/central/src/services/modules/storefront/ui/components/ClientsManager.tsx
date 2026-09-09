'use client';

import { useState, useCallback, useEffect } from 'react';
import { StorefrontClient, CreateClientDTO } from '../../domain/types';
import { getClientsAction, createClientAction } from '../../infra/actions';

interface ClientsManagerProps {
    businessId?: number;
}

export function ClientsManager({ businessId }: ClientsManagerProps) {
    const [clients, setClients] = useState<StorefrontClient[]>([]);
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [tempPassword, setTempPassword] = useState<string | null>(null);
    const [form, setForm] = useState<CreateClientDTO>({ name: '', email: '', phone: '', password: '' });

    const loadClients = useCallback(async () => {
        setLoading(true);
        const result = await getClientsAction({ page: 1, page_size: 50, business_id: businessId });
        setClients(result.data);
        setLoading(false);
    }, [businessId]);

    useEffect(() => {
        loadClients();
    }, [loadClients]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitting(true);
        setError(null);
        setTempPassword(null);

        const result = await createClientAction(
            { ...form, password: form.password || undefined },
            businessId,
        );

        if (!result.success) {
            setError(result.message || 'Error al crear el cliente');
            setSubmitting(false);
            return;
        }

        if (result.temp_password) {
            setTempPassword(result.temp_password);
        }
        setForm({ name: '', email: '', phone: '', password: '' });
        setSubmitting(false);
        await loadClients();
    };

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-1">
                <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
                    <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-4">Nuevo Cliente</h3>

                    <form onSubmit={handleSubmit} className="space-y-3">
                        <input
                            type="text"
                            placeholder="Nombre completo"
                            value={form.name}
                            onChange={e => setForm(prev => ({ ...prev, name: e.target.value }))}
                            required
                            className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                        />
                        <input
                            type="email"
                            placeholder="Correo electronico"
                            value={form.email}
                            onChange={e => setForm(prev => ({ ...prev, email: e.target.value }))}
                            required
                            className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                        />
                        <input
                            type="text"
                            placeholder="Telefono"
                            value={form.phone}
                            onChange={e => setForm(prev => ({ ...prev, phone: e.target.value }))}
                            className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                        />
                        <input
                            type="text"
                            placeholder="Contrasena (opcional, se genera una si la dejas vacia)"
                            value={form.password}
                            onChange={e => setForm(prev => ({ ...prev, password: e.target.value }))}
                            className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                        />

                        {error && <p className="text-red-500 text-sm">{error}</p>}

                        {tempPassword && (
                            <div className="p-3 bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800 rounded-lg text-sm">
                                <p className="text-emerald-800 dark:text-emerald-300 font-medium">Cliente creado</p>
                                <p className="text-emerald-700 dark:text-emerald-400 mt-1">
                                    Contrasena temporal: <span className="font-mono font-bold">{tempPassword}</span>
                                </p>
                                <p className="text-emerald-600 dark:text-emerald-500 mt-1 text-xs">
                                    Copiala ahora, no se volvera a mostrar.
                                </p>
                            </div>
                        )}

                        <button
                            type="submit"
                            disabled={submitting}
                            className="w-full py-2 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
                        >
                            {submitting ? 'Creando...' : 'Crear cliente'}
                        </button>
                    </form>
                </div>
            </div>

            <div className="lg:col-span-2">
                <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
                    {loading ? (
                        <div className="text-center py-8 text-gray-500 dark:text-gray-400">Cargando clientes...</div>
                    ) : clients.length === 0 ? (
                        <div className="text-center py-8 text-gray-500 dark:text-gray-400">Aun no hay clientes registrados</div>
                    ) : (
                        <table className="w-full">
                            <thead className="bg-gray-50 dark:bg-gray-700">
                                <tr>
                                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Nombre</th>
                                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Correo</th>
                                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Telefono</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                                {clients.map(client => (
                                    <tr key={client.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50">
                                        <td className="px-4 py-3 text-sm font-medium text-gray-900 dark:text-white">{client.name}</td>
                                        <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-300">{client.email || '-'}</td>
                                        <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-300">{client.phone || '-'}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    )}
                </div>
            </div>
        </div>
    );
}
