'use client';

import { useState, useEffect, useCallback } from 'react';
import { PencilIcon, TrashIcon, EyeIcon } from '@heroicons/react/24/outline';
import { getCustomersAction, deleteCustomerAction } from '../../infra/actions';
import { CustomerInfo, GetCustomersParams } from '../../domain/types';
import { Alert, Table, Spinner } from '@/shared/ui';

interface CustomerListProps {
    onView?: (customer: CustomerInfo) => void;
    onEdit?: (customer: CustomerInfo) => void;
    onRefreshRef?: (ref: () => void) => void;
    selectedBusinessId?: number;
}

export default function CustomerList({ onView, onEdit, onRefreshRef, selectedBusinessId }: CustomerListProps) {
    const [customers, setCustomers] = useState<CustomerInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);

    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');

    const fetchCustomers = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const params: GetCustomersParams = { page, page_size: pageSize };
            if (search) params.search = search;
            if (selectedBusinessId) params.business_id = selectedBusinessId;

            const response = await getCustomersAction(params);
            setCustomers(response.data || []);
            setTotal(response.total || 0);
            setTotalPages(response.total_pages || 1);
            setPage(response.page || page);
        } catch (err: any) {
            setError(err.message || 'Error al cargar clientes');
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, search, selectedBusinessId]);

    useEffect(() => {
        fetchCustomers();
    }, [fetchCustomers]);

    useEffect(() => {
        onRefreshRef?.(fetchCustomers);
    }, [fetchCustomers, onRefreshRef]);

    // Resetear a página 1 cuando cambia el negocio seleccionado
    useEffect(() => {
        setPage(1);
        setSearch('');
        setSearchInput('');
    }, [selectedBusinessId]);

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        setSearch(searchInput);
        setPage(1);
    };

    const handleClearSearch = () => {
        setSearchInput('');
        setSearch('');
        setPage(1);
    };

    const handleDelete = async (customer: CustomerInfo) => {
        if (!confirm(`¿Eliminar al cliente "${customer.name}"? Esta acción no se puede deshacer.`)) return;
        try {
            await deleteCustomerAction(customer.id, selectedBusinessId);
            fetchCustomers();
        } catch (err: any) {
            setError(err.message || 'Error al eliminar el cliente');
        }
    };

    const columns = [
        { key: 'name', label: 'Nombre' },
        { key: 'email', label: 'Email' },
        { key: 'phone', label: 'Teléfono' },
        { key: 'dni', label: 'Documento' },
        { key: 'created_at', label: 'Creado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (customer: CustomerInfo) => ({
        name: (
            <span className="font-medium text-gray-900">{customer.name}</span>
        ),
        email: (
            <span className="text-sm text-gray-600">{customer.email || <span className="text-gray-300">—</span>}</span>
        ),
        phone: (
            <span className="text-sm text-gray-600">{customer.phone || <span className="text-gray-300">—</span>}</span>
        ),
        dni: (
            <span className="text-sm text-gray-600">{customer.dni || <span className="text-gray-300">—</span>}</span>
        ),
        created_at: (
            <span className="text-sm text-gray-500">
                {new Date(customer.created_at).toLocaleDateString('es-CO')}
            </span>
        ),
        actions: (
            <div className="flex justify-end gap-2">
                {onView && (
                    <button
                        onClick={() => onView(customer)}
                        className="p-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors"
                        title="Ver detalle"
                    >
                        <EyeIcon className="w-4 h-4" />
                    </button>
                )}
                {onEdit && (
                    <button
                        onClick={() => onEdit(customer)}
                        className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors"
                        title="Editar"
                    >
                        <PencilIcon className="w-4 h-4" />
                    </button>
                )}
                <button
                    onClick={() => handleDelete(customer)}
                    className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors"
                    title="Eliminar"
                >
                    <TrashIcon className="w-4 h-4" />
                </button>
            </div>
        ),
    });

    if (loading && customers.length === 0) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
    }

    return (
        <div className="space-y-4">
            {/* Buscador */}
            <form onSubmit={handleSearch} className="flex gap-2">
                <input
                    type="text"
                    value={searchInput}
                    onChange={(e) => setSearchInput(e.target.value)}
                    placeholder="Buscar por nombre, email o teléfono..."
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                <button
                    type="submit"
                    className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700 transition-colors"
                >
                    Buscar
                </button>
                {search && (
                    <button
                        type="button"
                        onClick={handleClearSearch}
                        className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg text-sm hover:bg-gray-200 transition-colors"
                    >
                        Limpiar
                    </button>
                )}
            </form>

            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                <Table
                    columns={columns}
                    data={customers.map(renderRow)}
                    keyExtractor={(_, index) => String(customers[index]?.id || index)}
                    emptyMessage="No hay clientes registrados"
                    loading={loading}
                    pagination={{
                        currentPage: page,
                        totalPages: totalPages,
                        totalItems: total,
                        itemsPerPage: pageSize,
                        onPageChange: (newPage) => setPage(newPage),
                        onItemsPerPageChange: (newSize) => {
                            setPageSize(newSize);
                            setPage(1);
                        },
                    }}
                />
            </div>
        </div>
    );
}
