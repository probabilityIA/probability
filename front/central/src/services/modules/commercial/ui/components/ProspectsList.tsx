'use client';

import { useMemo, useState } from 'react';
import { Table } from '@/shared/ui/table';
import { Badge } from '@/shared/ui/badge';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui';
import { useProspects } from '../hooks/useProspects';
import { Prospect, ProspectStats } from '../../domain/types';
import { updateProspectSeenAction } from '../../infra/actions';

const PLATFORM_LABELS: Record<string, string> = {
  shopify: 'Shopify',
  woocommerce: 'WooCommerce',
  vtex: 'VTEX',
  magento: 'Magento',
  instagram_only: 'Solo Instagram',
  whatsapp_only: 'Solo WhatsApp',
  desconocida: 'Desconocida',
};

const STATUS_LABELS: Record<string, string> = {
  calificado: 'Calificado',
  descartado: 'Descartado',
  nuevo: 'Nuevo',
  contactado: 'Contactado',
  cliente: 'Cliente',
};

const SCORE_BUCKETS: Array<{ value: string; label: string; min: number; max: number }> = [
  { value: '0-20', label: '0 - 20', min: 0, max: 20 },
  { value: '20-40', label: '20 - 40', min: 20, max: 40 },
  { value: '40-60', label: '40 - 60', min: 40, max: 60 },
  { value: '60-80', label: '60 - 80', min: 60, max: 80 },
  { value: '80-100', label: '80 - 100', min: 80, max: 100 },
];

interface ProspectsListProps {
  stats: ProspectStats | null;
}

export function ProspectsList({ stats }: ProspectsListProps) {
  const {
    prospects, loading, totalCount, currentPage, pageSize, totalPages,
    filters, updateFilters, handlePageChange, handlePageSizeChange,
  } = useProspects();

  const [seenOverrides, setSeenOverrides] = useState<Record<number, boolean>>({});

  const cityOptions = useMemo(
    () => (stats?.by_city || []).map((c) => ({ value: c.key, label: c.key })),
    [stats]
  );

  const availableFilters: FilterOption[] = useMemo(() => [
    {
      key: 'search',
      label: 'Buscar',
      type: 'text',
      placeholder: 'Dominio, nombre o @instagram...',
    },
    {
      key: 'platform',
      label: 'Plataforma',
      type: 'select',
      options: Object.entries(PLATFORM_LABELS).map(([value, label]) => ({ value, label })),
    },
    {
      key: 'status',
      label: 'Estado',
      type: 'select',
      options: Object.entries(STATUS_LABELS).map(([value, label]) => ({ value, label })),
    },
    {
      key: 'city',
      label: 'Ciudad',
      type: 'select',
      options: cityOptions,
    },
    {
      key: 'score_bucket',
      label: 'Score',
      type: 'select',
      options: SCORE_BUCKETS.map((b) => ({ value: b.value, label: b.label })),
    },
    {
      key: 'has_cod',
      label: 'Pago contraentrega',
      type: 'boolean',
    },
  ], [cityOptions]);

  const scoreBucketValue = useMemo(() => {
    if (filters.min_score === undefined || filters.max_score === undefined) return undefined;
    return SCORE_BUCKETS.find((b) => b.min === filters.min_score && b.max === filters.max_score)?.value;
  }, [filters.min_score, filters.max_score]);

  const activeFilters: ActiveFilter[] = useMemo(() => {
    const active: ActiveFilter[] = [];
    if (filters.search) active.push({ key: 'search', label: 'Buscar', value: filters.search, type: 'text' });
    if (filters.platform) active.push({ key: 'platform', label: 'Plataforma', value: PLATFORM_LABELS[filters.platform] || filters.platform, type: 'select' });
    if (filters.status) active.push({ key: 'status', label: 'Estado', value: STATUS_LABELS[filters.status] || filters.status, type: 'select' });
    if (filters.city) active.push({ key: 'city', label: 'Ciudad', value: filters.city, type: 'select' });
    if (scoreBucketValue) active.push({ key: 'score_bucket', label: 'Score', value: scoreBucketValue, type: 'select' });
    if (filters.has_cod !== undefined) active.push({ key: 'has_cod', label: 'Pago contraentrega', value: filters.has_cod, type: 'boolean' });
    return active;
  }, [filters, scoreBucketValue]);

  const handleAddFilter = (key: string, value: any) => {
    if (key === 'score_bucket') {
      const bucket = SCORE_BUCKETS.find((b) => b.value === value);
      if (!bucket) return;
      updateFilters({ ...filters, min_score: bucket.min, max_score: bucket.max });
      return;
    }
    updateFilters({ ...filters, [key]: value });
  };

  const handleRemoveFilter = (key: string) => {
    const next = { ...filters };
    if (key === 'score_bucket') {
      delete (next as any).min_score;
      delete (next as any).max_score;
    } else {
      delete (next as any)[key];
    }
    updateFilters(next);
  };

  const handleToggleSeen = async (prospect: Prospect, nextSeen: boolean) => {
    setSeenOverrides((prev) => ({ ...prev, [prospect.id]: nextSeen }));
    try {
      await updateProspectSeenAction(prospect.id, nextSeen);
    } catch (error) {
      console.error('Error al actualizar visto:', error);
      setSeenOverrides((prev) => ({ ...prev, [prospect.id]: prospect.seen }));
    }
  };

  const columns = [
    {
      key: 'seen',
      label: 'Visto',
      render: (_: unknown, p: Prospect) => (
        <input
          type="checkbox"
          checked={seenOverrides[p.id] ?? p.seen}
          onChange={(e) => handleToggleSeen(p, e.target.checked)}
          className="h-4 w-4 rounded border-gray-300 text-[var(--color-primary)] focus:ring-[var(--color-primary)]"
          aria-label={`Marcar ${p.business_name || p.domain} como visto`}
        />
      ),
    },
    {
      key: 'domain',
      label: 'Negocio',
      render: (_: unknown, p: Prospect) => (
        <div>
          <a href={p.profile_url} target="_blank" rel="noopener noreferrer" className="font-medium text-[var(--color-primary)] hover:underline">
            {p.business_name || p.domain}
          </a>
          {p.business_name && <div className="text-xs text-gray-400">{p.domain}</div>}
        </div>
      ),
    },
    {
      key: 'platform',
      label: 'Plataforma',
      render: (_: unknown, p: Prospect) => (
        <span className="text-xs text-gray-600 dark:text-gray-300">{PLATFORM_LABELS[p.platform] || p.platform}</span>
      ),
    },
    {
      key: 'city',
      label: 'Ciudad',
      render: (_: unknown, p: Prospect) => <span className="text-sm">{p.city || '-'}</span>,
    },
    {
      key: 'contact',
      label: 'Contacto',
      render: (_: unknown, p: Prospect) => (
        <div className="flex items-center gap-2 text-xs">
          {p.whatsapp_url && (
            <a href={p.whatsapp_url} target="_blank" rel="noopener noreferrer" className="text-emerald-600 hover:underline">
              WhatsApp
            </a>
          )}
          {p.instagram_url && (
            <a href={p.instagram_url} target="_blank" rel="noopener noreferrer" className="text-pink-600 hover:underline">
              Instagram
            </a>
          )}
          {!p.whatsapp_url && !p.instagram_url && <span className="text-gray-400">-</span>}
        </div>
      ),
    },
    {
      key: 'product_count',
      label: 'Productos',
      render: (_: unknown, p: Prospect) => <span className="text-sm">{p.product_count || '-'}</span>,
    },
    {
      key: 'has_cod',
      label: 'COD',
      render: (_: unknown, p: Prospect) => (
        p.has_cod
          ? <Badge type="success">Si</Badge>
          : <span className="text-gray-300">-</span>
      ),
    },
    {
      key: 'score',
      label: 'Score',
      render: (_: unknown, p: Prospect) => (
        <span className="font-semibold text-gray-900 dark:text-white">{p.score.toFixed(0)}</span>
      ),
    },
    {
      key: 'status',
      label: 'Estado',
      render: (_: unknown, p: Prospect) => (
        <Badge type={p.status === 'calificado' ? 'success' : 'secondary'}>
          {STATUS_LABELS[p.status] || p.status}
        </Badge>
      ),
    },
  ];

  return (
    <div className="space-y-4">
      <DynamicFilters
        availableFilters={availableFilters}
        activeFilters={activeFilters}
        onAddFilter={handleAddFilter}
        onRemoveFilter={handleRemoveFilter}
      />

      <Table
        data={prospects}
        columns={columns}
        loading={loading}
        emptyMessage="No hay prospectos para mostrar"
        keyExtractor={(p: Prospect) => p.id}
        pagination={{
          currentPage,
          totalPages,
          totalItems: totalCount,
          itemsPerPage: pageSize,
          onPageChange: handlePageChange,
          onItemsPerPageChange: handlePageSizeChange,
          showItemsPerPageSelector: true,
          itemsPerPageOptions: [10, 20, 50],
        }}
      />
    </div>
  );
}
