/**
 * ConfigListTable - Lista de configuraciones de notificación usando Table global
 */
'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { Table, TableColumn, Badge, Button } from '@/shared/ui';
import type { FilterOption, ActiveFilter } from '@/shared/ui';
import { NotificationConfig } from '../../domain/types';
import { getConfigsAction, deleteConfigAction } from '../../infra/actions';
import { useToast } from '@/shared/providers/toast-provider';

interface ConfigListTableProps {
  onEdit: (config: NotificationConfig) => void;
  onCreate: () => void;
  refreshKey?: number;
}

export function ConfigListTable({ onEdit, onCreate, refreshKey = 0 }: ConfigListTableProps) {
  const { showToast } = useToast();

  // Estado
  const [configs, setConfigs] = useState<NotificationConfig[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(1);

  // Filtros
  const [activeFilters, setActiveFilters] = useState<ActiveFilter[]>([]);
  const [sortBy, setSortBy] = useState('created_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  // Filtros disponibles
  const availableFilters: FilterOption[] = useMemo(() => [
    {
      key: 'event_type',
      label: 'Tipo de Evento',
      type: 'text',
      placeholder: 'Ej: order.created',
    },
    {
      key: 'enabled',
      label: 'Solo Activas',
      type: 'boolean',
    },
  ], []);

  // Opciones de ordenamiento
  const sortOptions = useMemo(() => [
    { value: 'created_at', label: 'Ordenar por fecha de creación' },
    { value: 'updated_at', label: 'Ordenar por última actualización' },
    { value: 'event_type', label: 'Ordenar por tipo de evento' },
  ], []);

  // Cargar configs
  const fetchConfigs = useCallback(async () => {
    setLoading(true);
    try {
      // Construir filtros desde activeFilters
      const filters: Record<string, any> = {};
      activeFilters.forEach((filter) => {
        if (filter.type === 'boolean') {
          filters[filter.key] = true;
        } else {
          filters[filter.key] = filter.value;
        }
      });

      const response = await getConfigsAction({
        ...filters,
        page,
        page_size: pageSize,
        sort_by: sortBy,
        sort_order: sortOrder,
      });

      setConfigs(response.data || []);
      setTotal(response.total || 0);
      setTotalPages(response.total_pages || 1);
    } catch (error) {
      console.error('Error al cargar configuraciones:', error);
      showToast('Error al cargar configuraciones', 'error');
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, activeFilters, sortBy, sortOrder, showToast]);

  // Efecto para cargar datos
  useEffect(() => {
    fetchConfigs();
  }, [fetchConfigs, refreshKey]);

  // Handlers de filtros
  const handleAddFilter = useCallback((filterKey: string, value: any) => {
    const filter = availableFilters.find(f => f.key === filterKey);
    if (!filter) return;

    setActiveFilters(prev => [
      ...prev.filter(f => f.key !== filterKey),
      {
        key: filterKey,
        label: filter.label,
        value,
        type: filter.type,
      },
    ]);
    setPage(1); // Reset a página 1 cuando se agrega filtro
  }, [availableFilters]);

  const handleRemoveFilter = useCallback((filterKey: string) => {
    setActiveFilters(prev => prev.filter(f => f.key !== filterKey));
    setPage(1);
  }, []);

  const handleSortChange = useCallback((newSortBy: string, newSortOrder: 'asc' | 'desc') => {
    setSortBy(newSortBy);
    setSortOrder(newSortOrder);
    setPage(1);
  }, []);

  // Handler para eliminar
  const handleDelete = useCallback(async (id: number) => {
    if (!confirm('¿Estás seguro de eliminar esta configuración?')) return;

    try {
      await deleteConfigAction(id);
      showToast('Configuración eliminada correctamente', 'success');
      fetchConfigs();
    } catch (error) {
      console.error('Error al eliminar configuración:', error);
      showToast('Error al eliminar configuración', 'error');
    }
  }, [showToast, fetchConfigs]);

  // Definir columnas de la tabla
  const columns: TableColumn<NotificationConfig>[] = useMemo(() => [
    {
      key: 'id',
      label: 'ID',
      width: '80px',
      render: (value) => <span className="text-gray-600">#{value as number}</span>,
    },
    {
      key: 'event_type',
      label: 'Tipo de Evento',
      render: (_, row) => (
        <span className="font-medium text-sm text-gray-900">
          {row.notification_event_name || row.event_type || '-'}
        </span>
      ),
    },
    {
      key: 'channels',
      label: 'Canales',
      render: (_, row) => (
        <div className="flex gap-1 flex-wrap">
          {(row.channels || ['sse']).map((channel) => (
            <Badge key={channel} type="secondary">
              {channel}
            </Badge>
          ))}
        </div>
      ),
    },
    {
      key: 'enabled',
      label: 'Estado',
      width: '120px',
      align: 'center',
      render: (value) => (
        <Badge type={value ? 'success' : 'secondary'}>
          {value ? 'Activo' : 'Inactivo'}
        </Badge>
      ),
    },
    {
      key: 'description',
      label: 'Descripción',
      render: (value) => (
        <span className="text-gray-600 line-clamp-2">{value as string || '-'}</span>
      ),
    },
    {
      key: 'actions',
      label: 'Acciones',
      width: '150px',
      align: 'right',
      render: (_, row) => (
        <div className="flex gap-2 justify-end">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onEdit(row)}
            title="Editar"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => handleDelete(row.id)}
            className="text-red-500 hover:text-red-700"
            title="Eliminar"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </Button>
        </div>
      ),
    },
  ], [onEdit, handleDelete]);

  return (
    <Table<NotificationConfig>
      columns={columns}
      data={configs}
      keyExtractor={(row) => row.id}
      loading={loading}
      emptyMessage="No hay configuraciones de notificación"
      pagination={{
        currentPage: page,
        totalPages,
        totalItems: total,
        itemsPerPage: pageSize,
        onPageChange: setPage,
        onItemsPerPageChange: setPageSize,
        showItemsPerPageSelector: true,
        itemsPerPageOptions: [10, 20, 50, 100],
      }}
      filters={{
        availableFilters,
        activeFilters,
        onAddFilter: handleAddFilter,
        onRemoveFilter: handleRemoveFilter,
        onCreate,
        createButtonIconOnly: true, // Solo mostrar ícono +
        sortBy,
        sortOrder,
        onSortChange: handleSortChange,
        sortOptions,
      }}
    />
  );
}
