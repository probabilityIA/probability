'use client';

import { useCallback, useEffect, useState } from 'react';
import { getProspectsAction } from '../../infra/actions';
import { GetProspectsParams, Prospect } from '../../domain/types';

const PAGE_SIZE_DEFAULT = 10;

export function useProspects() {
  const [prospects, setProspects] = useState<Prospect[]>([]);
  const [loading, setLoading] = useState(true);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(PAGE_SIZE_DEFAULT);
  const [filters, setFilters] = useState<Omit<GetProspectsParams, 'page' | 'page_size'>>({});

  const load = useCallback(async (page: number, size: number, currentFilters: typeof filters) => {
    try {
      setLoading(true);
      const response = await getProspectsAction({ ...currentFilters, page, page_size: size });
      setProspects(response.data || []);
      setTotalCount(response.total || 0);
    } catch (error) {
      console.error('Error al cargar prospectos:', error);
      setProspects([]);
      setTotalCount(0);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    setCurrentPage(1);
    load(1, pageSize, filters);
  }, [JSON.stringify(filters)]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    load(page, pageSize, filters);
  };

  const handlePageSizeChange = (size: number) => {
    setPageSize(size);
    setCurrentPage(1);
    load(1, size, filters);
  };

  const updateFilters = (next: typeof filters) => {
    setFilters(next);
  };

  const totalPages = Math.max(1, Math.ceil(totalCount / pageSize));

  return {
    prospects,
    loading,
    totalCount,
    currentPage,
    pageSize,
    totalPages,
    filters,
    updateFilters,
    handlePageChange,
    handlePageSizeChange,
  };
}
