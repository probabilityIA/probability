'use client';

import { useState, useEffect, useCallback } from 'react';
import { PencilIcon, TrashIcon, EyeIcon, ArrowPathIcon, ChartBarIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { getAnnouncementsAction, deleteAnnouncementAction, changeStatusAction, forceRedisplayAction } from '../../infra/actions';
import { AnnouncementInfo, GetAnnouncementsParams, AnnouncementStatus } from '../../domain/types';
import { Alert, Table, Spinner, Badge } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';

interface AnnouncementListProps {
    onEdit?: (announcement: AnnouncementInfo) => void;
    onRefreshRef?: (ref: () => void) => void;
    selectedBusinessId?: number;
}

const statusConfig: Record<AnnouncementStatus, { label: string; type: 'success' | 'warning' | 'error' | 'secondary' | 'primary' }> = {
    active: { label: 'Activo', type: 'success' },
    scheduled: { label: 'Programado', type: 'warning' },
    draft: { label: 'Borrador', type: 'secondary' },
    inactive: { label: 'Inactivo', type: 'error' },
};

const displayTypeLabels: Record<string, string> = {
    modal_image: 'Modal imagen',
    modal_text: 'Modal texto',
    ticker: 'Ticker',
};

export default function AnnouncementList({ onEdit, onRefreshRef, selectedBusinessId }: AnnouncementListProps) {
    const [announcements, setAnnouncements] = useState<AnnouncementInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);

    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');
    const [statusFilter, setStatusFilter] = useState<string>('');

    const fetchAnnouncements = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const params: GetAnnouncementsParams = { page, page_size: pageSize };
            if (search) params.search = search;
            if (statusFilter) params.status = statusFilter as AnnouncementStatus;
            if (selectedBusinessId) params.business_id = selectedBusinessId;

            const response = await getAnnouncementsAction(params);
            setAnnouncements(response.data || []);
            setTotal(response.total || 0);
            setTotalPages(response.total_pages || 1);
            setPage(response.page || page);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar anuncios'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, search, statusFilter, selectedBusinessId]);

    useEffect(() => {
        fetchAnnouncements();
    }, [fetchAnnouncements]);

    useEffect(() => {
        onRefreshRef?.(fetchAnnouncements);
    }, [fetchAnnouncements, onRefreshRef]);

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

    const handleDelete = async (announcement: AnnouncementInfo) => {
        if (!confirm(`Eliminar el anuncio "${announcement.title}"? Esta accion no se puede deshacer.`)) return;
        try {
            await deleteAnnouncementAction(announcement.id);
            fetchAnnouncements();
        } catch (err: any) {
            setError(getActionError(err, 'Error al eliminar el anuncio'));
        }
    };

    const handleToggleStatus = async (announcement: AnnouncementInfo) => {
        const newStatus: AnnouncementStatus = announcement.status === 'active' ? 'inactive' : 'active';
        try {
            await changeStatusAction(announcement.id, { status: newStatus });
            setSuccess(`Anuncio ${newStatus === 'active' ? 'activado' : 'desactivado'}`);
            setTimeout(() => setSuccess(null), 3000);
            fetchAnnouncements();
        } catch (err: any) {
            setError(getActionError(err, 'Error al cambiar estado'));
        }
    };

    const handleForceRedisplay = async (announcement: AnnouncementInfo) => {
        if (!confirm(`Forzar que se muestre de nuevo "${announcement.title}" a todos los usuarios?`)) return;
        try {
            await forceRedisplayAction(announcement.id);
            setSuccess('Se forzara la re-visualizacion del anuncio');
            setTimeout(() => setSuccess(null), 3000);
            fetchAnnouncements();
        } catch (err: any) {
            setError(getActionError(err, 'Error al forzar re-visualizacion'));
        }
    };

    const formatDate = (date: string | null) => {
        if (!date) return '--';
        return new Date(date).toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric' });
    };

    const columns = [
        { key: 'title', label: 'Titulo' },
        { key: 'category', label: 'Categoria' },
        { key: 'display_type', label: 'Tipo' },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'dates', label: 'Vigencia' },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (announcement: AnnouncementInfo) => {
        const config = statusConfig[announcement.status] || statusConfig.draft;
        return {
            title: (
                <div>
                    <span className="font-medium text-gray-900 dark:text-white">{announcement.title}</span>
                    {announcement.is_global && (
                        <span className="ml-2 text-xs bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400 px-1.5 py-0.5 rounded">
                            Global
                        </span>
                    )}
                </div>
            ),
            category: (
                <span className="text-sm text-gray-600 dark:text-gray-300">
                    {announcement.category ? (
                        <span className="inline-flex items-center gap-1">
                            <span className="w-2 h-2 rounded-full" style={{ backgroundColor: announcement.category.color }} />
                            {announcement.category.name}
                        </span>
                    ) : '--'}
                </span>
            ),
            display_type: (
                <span className="text-sm text-gray-600 dark:text-gray-300">
                    {displayTypeLabels[announcement.display_type] || announcement.display_type}
                </span>
            ),
            status: <Badge type={config.type}>{config.label}</Badge>,
            dates: (
                <span className="text-sm text-gray-600 dark:text-gray-300">
                    {formatDate(announcement.starts_at)} - {formatDate(announcement.ends_at)}
                </span>
            ),
            actions: (
                <div className="flex justify-end gap-1.5">
                    <button
                        onClick={() => handleToggleStatus(announcement)}
                        className={`p-1.5 text-white rounded-md transition-colors ${
                            announcement.status === 'active'
                                ? 'bg-orange-500 hover:bg-orange-600'
                                : 'bg-green-500 hover:bg-green-600'
                        }`}
                        title={announcement.status === 'active' ? 'Desactivar' : 'Activar'}
                    >
                        <EyeIcon className="w-4 h-4" />
                    </button>
                    <Link
                        href={`/announcements/${announcement.id}/stats`}
                        className="p-1.5 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors"
                        title="Estadisticas"
                    >
                        <ChartBarIcon className="w-4 h-4" />
                    </Link>
                    <button
                        onClick={() => handleForceRedisplay(announcement)}
                        className="p-1.5 bg-indigo-500 hover:bg-indigo-600 text-white rounded-md transition-colors"
                        title="Forzar re-visualizacion"
                    >
                        <ArrowPathIcon className="w-4 h-4" />
                    </button>
                    {onEdit && (
                        <button
                            onClick={() => onEdit(announcement)}
                            className="p-1.5 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors"
                            title="Editar"
                        >
                            <PencilIcon className="w-4 h-4" />
                        </button>
                    )}
                    <button
                        onClick={() => handleDelete(announcement)}
                        className="p-1.5 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors"
                        title="Eliminar"
                    >
                        <TrashIcon className="w-4 h-4" />
                    </button>
                </div>
            ),
        };
    };

    if (loading && announcements.length === 0) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
    }

    return (
        <div className="space-y-4">
            <div className="flex gap-2 flex-wrap">
                <form onSubmit={handleSearch} className="flex gap-2 flex-1 min-w-[200px]">
                    <input
                        type="text"
                        value={searchInput}
                        onChange={(e) => setSearchInput(e.target.value)}
                        placeholder="Buscar por titulo..."
                        className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm text-gray-900 dark:text-white bg-white dark:bg-gray-700 placeholder-gray-500 dark:placeholder-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
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
                            className="px-4 py-2 bg-gray-100 text-gray-700 dark:text-gray-200 rounded-lg text-sm hover:bg-gray-200 transition-colors"
                        >
                            Limpiar
                        </button>
                    )}
                </form>
                <select
                    value={statusFilter}
                    onChange={(e) => { setStatusFilter(e.target.value); setPage(1); }}
                    className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm text-gray-900 dark:text-white bg-white dark:bg-gray-700 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                    <option value="">Todos los estados</option>
                    <option value="active">Activo</option>
                    <option value="scheduled">Programado</option>
                    <option value="draft">Borrador</option>
                    <option value="inactive">Inactivo</option>
                </select>
            </div>

            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}
            {success && (
                <Alert type="success" onClose={() => setSuccess(null)}>
                    {success}
                </Alert>
            )}

            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                <Table
                    columns={columns}
                    data={announcements.map(renderRow)}
                    keyExtractor={(_, index) => String(announcements[index]?.id || index)}
                    emptyMessage="No hay anuncios registrados"
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
