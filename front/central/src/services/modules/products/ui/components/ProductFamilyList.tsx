'use client';

import { useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import { ProductFamily, Product } from '../../domain/types';
import { getProductFamiliesAction, deleteProductFamilyAction, getFamilyVariantsAction } from '../../infra/actions';
import { ChevronRightIcon, XMarkIcon } from '@heroicons/react/24/outline';

interface ProductFamilyListProps {
    onEdit?: (family: ProductFamily) => void;
    selectedBusinessId?: number;
}

export interface ProductFamilyListHandle {
    refresh: () => void;
}

const MODAL_PAGE_SIZE = 10;

const ProductFamilyList = forwardRef<ProductFamilyListHandle, ProductFamilyListProps>(
    ({ onEdit, selectedBusinessId }, ref) => {
        const [families, setFamilies] = useState<ProductFamily[]>([]);
        const [loading, setLoading] = useState(true);
        const [error, setError] = useState<string | null>(null);
        const [page, setPage] = useState(1);
        const [totalPages, setTotalPages] = useState(1);
        const [total, setTotal] = useState(0);
        const [pageSize, setPageSize] = useState(10);
        const [searchName, setSearchName] = useState('');

        const [modalFamily, setModalFamily] = useState<ProductFamily | null>(null);
        const [modalVariants, setModalVariants] = useState<Product[]>([]);
        const [modalLoading, setModalLoading] = useState(false);
        const [modalPage, setModalPage] = useState(1);

        useImperativeHandle(ref, () => ({ refresh: fetchFamilies }));

        const fetchFamilies = async () => {
            setLoading(true);
            setError(null);
            try {
                const params: any = { page, page_size: pageSize };
                if (selectedBusinessId) params.business_id = selectedBusinessId;
                if (searchName) params.name = searchName;
                const res = await getProductFamiliesAction(params);
                if (res.success === false) {
                    setError(res.message || 'Error al cargar familias');
                    setFamilies([]);
                } else {
                    setFamilies(res.data || []);
                    setTotal(res.total || 0);
                    setTotalPages(res.total_pages || 1);
                }
            } catch (e: any) {
                setError(e.message || 'Error inesperado');
            } finally {
                setLoading(false);
            }
        };

        useEffect(() => { fetchFamilies(); }, [page, pageSize, selectedBusinessId, searchName]);

        const handleDelete = async (family: ProductFamily) => {
            if (!confirm(`Eliminar la familia "${family.name}"? Esta accion no se puede deshacer.`)) return;
            const res = await deleteProductFamilyAction(family.id, selectedBusinessId);
            if (res.success) {
                fetchFamilies();
            } else {
                const msg = res.message || res.error || 'Error al eliminar';
                if (msg.includes('variantes activas') || msg.includes('active variants')) {
                    alert('No se puede eliminar: la familia tiene variantes activas.');
                } else {
                    alert(msg);
                }
            }
        };

        const handleOpenModal = async (family: ProductFamily) => {
            setModalFamily(family);
            setModalPage(1);
            setModalVariants([]);
            setModalLoading(true);
            try {
                const res = await getFamilyVariantsAction(family.id, selectedBusinessId);
                setModalVariants(res.data || []);
            } finally {
                setModalLoading(false);
            }
        };

        const stockBadge = (qty: number) => {
            if (qty === 0) return <span className="px-2 py-0.5 rounded-full text-xs font-semibold bg-red-100 text-red-700">Agotado</span>;
            if (qty < 5) return <span className="px-2 py-0.5 rounded-full text-xs font-semibold bg-yellow-100 text-yellow-700">{qty}</span>;
            return <span className="px-2 py-0.5 rounded-full text-xs font-semibold bg-green-100 text-green-700">{qty}</span>;
        };

        const pagedVariants = modalVariants.slice((modalPage - 1) * MODAL_PAGE_SIZE, modalPage * MODAL_PAGE_SIZE);
        const modalTotalPages = Math.ceil(modalVariants.length / MODAL_PAGE_SIZE);

        return (
            <>
                <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                    <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between gap-4 flex-wrap">
                        <span className="text-sm font-medium text-gray-500 dark:text-gray-400">
                            {total} {total === 1 ? 'familia' : 'familias'}
                        </span>
                        <div className="flex items-center gap-3">
                            <input
                                type="text"
                                placeholder="Buscar familia..."
                                value={searchName}
                                onChange={(e) => { setSearchName(e.target.value); setPage(1); }}
                                className="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder:text-gray-400 focus:ring-2 focus:ring-[var(--color-primary)] focus:border-transparent"
                            />
                            <select
                                value={pageSize}
                                onChange={(e) => { setPageSize(Number(e.target.value)); setPage(1); }}
                                className="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                            >
                                {[10, 20, 50].map((n) => <option key={n} value={n}>{n} por pagina</option>)}
                            </select>
                        </div>
                    </div>

                    {loading ? (
                        <div className="flex justify-center items-center py-20">
                            <div className="spinner" />
                        </div>
                    ) : error ? (
                        <div className="text-center py-16 text-red-500 text-sm">{error}</div>
                    ) : families.length === 0 ? (
                        <div className="text-center py-16 text-gray-400 text-sm">No hay familias registradas</div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="table w-full">
                                <thead>
                                    <tr>
                                        <th className="text-left">Familia</th>
                                        <th className="text-left hidden md:table-cell">Categoria</th>
                                        <th className="text-left hidden md:table-cell">Marca</th>
                                        <th className="text-center">Variantes</th>
                                        <th className="text-center">Estado</th>
                                        <th className="text-right">Acciones</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                                    {families.map((family) => (
                                        <tr key={family.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                                            <td>
                                                <div className="flex items-center gap-3">
                                                    {family.image_url ? (
                                                        <img src={family.image_url} alt={family.name} className="w-8 h-8 rounded-lg object-cover flex-shrink-0" />
                                                    ) : (
                                                        <div className="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center flex-shrink-0">
                                                            <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                                                            </svg>
                                                        </div>
                                                    )}
                                                    <div>
                                                        <p className="text-sm font-semibold text-gray-900 dark:text-white">{family.name}</p>
                                                        {family.title && family.title !== family.name && (
                                                            <p className="text-xs text-gray-400 truncate max-w-[180px]">{family.title}</p>
                                                        )}
                                                    </div>
                                                </div>
                                            </td>
                                            <td className="hidden md:table-cell text-sm text-gray-600 dark:text-gray-300">{family.category || '-'}</td>
                                            <td className="hidden md:table-cell text-sm text-gray-600 dark:text-gray-300">{family.brand || '-'}</td>
                                            <td className="text-center">
                                                <button
                                                    onClick={() => handleOpenModal(family)}
                                                    className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-bold bg-indigo-100 text-indigo-700 hover:bg-indigo-200 transition-colors"
                                                    title="Ver variantes"
                                                >
                                                    {family.variant_count}
                                                    <ChevronRightIcon className="w-3 h-3" />
                                                </button>
                                            </td>
                                            <td className="text-center">
                                                <span className={`px-2 py-0.5 rounded-full text-xs font-semibold ${family.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'}`}>
                                                    {family.is_active ? 'Activa' : 'Inactiva'}
                                                </span>
                                            </td>
                                            <td>
                                                <div className="flex items-center justify-end gap-2">
                                                    {onEdit && (
                                                        <button
                                                            onClick={() => onEdit(family)}
                                                            className="p-1.5 text-blue-600 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-colors"
                                                            title="Editar"
                                                        >
                                                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                                                        </button>
                                                    )}
                                                    <button
                                                        onClick={() => handleDelete(family)}
                                                        className="p-1.5 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
                                                        title="Eliminar"
                                                    >
                                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                                                    </button>
                                                </div>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    )}

                    {!loading && totalPages > 1 && (
                        <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between">
                            <span className="text-sm text-gray-500 dark:text-gray-400">Pagina {page} de {totalPages}</span>
                            <div className="flex gap-2">
                                <button
                                    onClick={() => setPage((p) => Math.max(1, p - 1))}
                                    disabled={page === 1}
                                    className="px-3 py-1.5 text-sm rounded-lg border border-gray-300 dark:border-gray-600 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                                >
                                    Anterior
                                </button>
                                <button
                                    onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                                    disabled={page === totalPages}
                                    className="px-3 py-1.5 text-sm rounded-lg border border-gray-300 dark:border-gray-600 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                                >
                                    Siguiente
                                </button>
                            </div>
                        </div>
                    )}
                </div>

                {modalFamily && (
                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                        <div className="absolute inset-0 bg-black/50" onClick={() => setModalFamily(null)} />
                        <div className="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-4xl max-h-[90vh] flex flex-col">
                            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
                                <div className="flex items-center gap-3">
                                    {modalFamily.image_url ? (
                                        <img src={modalFamily.image_url} alt={modalFamily.name} className="w-9 h-9 rounded-lg object-cover" />
                                    ) : (
                                        <div className="w-9 h-9 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                                            <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                                            </svg>
                                        </div>
                                    )}
                                    <div>
                                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{modalFamily.name}</h2>
                                        <p className="text-sm text-gray-500 mt-0.5">
                                            {modalLoading ? 'Cargando...' : `${modalVariants.length} variantes`}
                                            {modalFamily.category ? ` · ${modalFamily.category}` : ''}
                                            {modalFamily.brand ? ` · ${modalFamily.brand}` : ''}
                                        </p>
                                    </div>
                                </div>
                                <button
                                    onClick={() => setModalFamily(null)}
                                    className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                                >
                                    <XMarkIcon className="w-5 h-5" />
                                </button>
                            </div>

                            <div className="overflow-auto flex-1 p-5 flex flex-col gap-4">
                                {modalLoading ? (
                                    <div className="flex justify-center py-12">
                                        <div className="spinner" />
                                    </div>
                                ) : modalVariants.length === 0 ? (
                                    <p className="text-center text-sm text-gray-400 py-12">Esta familia no tiene variantes aun.</p>
                                ) : (
                                    <>
                                        <div className="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700">
                                            <table className="table w-full">
                                                <thead>
                                                    <tr>
                                                        <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-600 dark:!text-gray-300 text-left">SKU</th>
                                                        <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-600 dark:!text-gray-300 text-left">Nombre</th>
                                                        <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-600 dark:!text-gray-300 text-left hidden sm:table-cell">Variante</th>
                                                        <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-600 dark:!text-gray-300 text-left hidden sm:table-cell">Barcode</th>
                                                        <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-600 dark:!text-gray-300 text-right">Precio</th>
                                                        <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-600 dark:!text-gray-300 text-center">Stock</th>
                                                    </tr>
                                                </thead>
                                                <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                                                    {pagedVariants.map((variant) => (
                                                        <tr key={variant.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                                                            <td className="font-mono text-xs text-gray-600 dark:text-gray-300">{variant.sku}</td>
                                                            <td className="text-sm text-gray-900 dark:text-white">{variant.name}</td>
                                                            <td className="hidden sm:table-cell">
                                                                {variant.variant_label ? (
                                                                    <span className="px-2 py-0.5 rounded text-xs font-medium bg-indigo-100 text-indigo-700">{variant.variant_label}</span>
                                                                ) : (
                                                                    <span className="text-gray-400 text-xs">&mdash;</span>
                                                                )}
                                                            </td>
                                                            <td className="hidden sm:table-cell font-mono text-xs text-gray-500 dark:text-gray-400">{variant.barcode || '-'}</td>
                                                            <td className="text-right text-sm font-semibold text-gray-900 dark:text-white">
                                                                {new Intl.NumberFormat('es-CO', { style: 'currency', currency: variant.currency || 'COP', maximumFractionDigits: 0 }).format(variant.price)}
                                                            </td>
                                                            <td className="text-center">{stockBadge(variant.stock_quantity ?? variant.stock ?? 0)}</td>
                                                        </tr>
                                                    ))}
                                                </tbody>
                                            </table>
                                        </div>

                                        {modalTotalPages > 1 && (
                                            <div className="flex items-center justify-between px-1">
                                                <span className="text-xs text-gray-500">{modalVariants.length} variantes totales</span>
                                                <div className="flex items-center gap-1">
                                                    <button
                                                        disabled={modalPage <= 1}
                                                        onClick={() => setModalPage((p) => p - 1)}
                                                        className="px-3 py-1.5 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                                                    >
                                                        Anterior
                                                    </button>
                                                    <span className="text-xs text-gray-600 px-2">{modalPage} / {modalTotalPages}</span>
                                                    <button
                                                        disabled={modalPage >= modalTotalPages}
                                                        onClick={() => setModalPage((p) => p + 1)}
                                                        className="px-3 py-1.5 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                                                    >
                                                        Siguiente
                                                    </button>
                                                </div>
                                            </div>
                                        )}
                                    </>
                                )}
                            </div>
                        </div>
                    </div>
                )}
            </>
        );
    }
);

ProductFamilyList.displayName = 'ProductFamilyList';
export default ProductFamilyList;
