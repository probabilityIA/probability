'use client';

import { useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import { ProductFamily, Product } from '../../domain/types';
import {
    getProductFamiliesAction,
    deleteProductFamilyAction,
    getFamilyVariantsAction
} from '../../infra/actions';

interface ProductFamilyListProps {
    onEdit?: (family: ProductFamily) => void;
    selectedBusinessId?: number;
}

export interface ProductFamilyListHandle {
    refresh: () => void;
}

const ProductFamilyList = forwardRef<ProductFamilyListHandle, ProductFamilyListProps>(
    ({ onEdit, selectedBusinessId }, ref) => {
        const [families, setFamilies] = useState<ProductFamily[]>([]);
        const [loading, setLoading] = useState(true);
        const [error, setError] = useState<string | null>(null);
        const [page, setPage] = useState(1);
        const [totalPages, setTotalPages] = useState(1);
        const [total, setTotal] = useState(0);
        const [pageSize, setPageSize] = useState(10);
        const [expandedFamily, setExpandedFamily] = useState<number | null>(null);
        const [variantsMap, setVariantsMap] = useState<Record<number, Product[]>>({});
        const [loadingVariants, setLoadingVariants] = useState<number | null>(null);
        const [searchName, setSearchName] = useState('');

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

        useEffect(() => {
            fetchFamilies();
        }, [page, pageSize, selectedBusinessId, searchName]);

        const handleDelete = async (family: ProductFamily) => {
            if (!confirm(`¿Eliminar la familia "${family.name}"? Esta accion no se puede deshacer.`)) return;
            const res = await deleteProductFamilyAction(family.id, selectedBusinessId);
            if (res.success) {
                fetchFamilies();
            } else {
                const msg = res.message || res.error || 'Error al eliminar';
                if (msg.includes('variantes activas') || msg.includes('active variants')) {
                    alert('No se puede eliminar: la familia tiene variantes activas. Elimina primero las variantes o desactivalas.');
                } else {
                    alert(msg);
                }
            }
        };

        const handleToggleExpand = async (familyId: number) => {
            if (expandedFamily === familyId) {
                setExpandedFamily(null);
                return;
            }
            setExpandedFamily(familyId);
            if (!variantsMap[familyId]) {
                setLoadingVariants(familyId);
                try {
                    const res = await getFamilyVariantsAction(familyId, selectedBusinessId);
                    setVariantsMap(prev => ({ ...prev, [familyId]: res.data || [] }));
                } finally {
                    setLoadingVariants(null);
                }
            }
        };

        const stockBadge = (qty: number) => {
            if (qty === 0) return <span className="px-2 py-0.5 rounded-full text-xs font-semibold bg-red-100 text-red-700">Agotado</span>;
            if (qty < 5) return <span className="px-2 py-0.5 rounded-full text-xs font-semibold bg-yellow-100 text-yellow-700">{qty}</span>;
            return <span className="px-2 py-0.5 rounded-full text-xs font-semibold bg-green-100 text-green-700">{qty}</span>;
        };

        return (
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
                <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between gap-4">
                    <div className="flex items-center gap-3">
                        <span className="text-sm font-medium text-gray-500 dark:text-gray-400">
                            {total} {total === 1 ? 'familia' : 'familias'}
                        </span>
                    </div>
                    <div className="flex items-center gap-3">
                        <input
                            type="text"
                            placeholder="Buscar familia..."
                            value={searchName}
                            onChange={e => { setSearchName(e.target.value); setPage(1); }}
                            className="px-4 py-2 text-sm border-2 border-slate-200 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-[#7c3aed] focus:border-[#7c3aed] bg-white dark:bg-gray-700 text-slate-900 dark:text-white placeholder:text-slate-400"
                        />
                        <select
                            value={pageSize}
                            onChange={e => { setPageSize(Number(e.target.value)); setPage(1); }}
                            className="px-3 py-2 text-sm border-2 border-slate-200 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-[#7c3aed] bg-white dark:bg-gray-700 text-slate-900 dark:text-white"
                        >
                            {[10, 20, 50].map(n => <option key={n} value={n}>{n} por pagina</option>)}
                        </select>
                    </div>
                </div>

                {loading ? (
                    <div className="flex justify-center items-center py-20">
                        <div className="w-8 h-8 border-4 border-[#7c3aed] border-t-transparent rounded-full animate-spin" />
                    </div>
                ) : error ? (
                    <div className="text-center py-16 text-red-500 text-sm">{error}</div>
                ) : families.length === 0 ? (
                    <div className="text-center py-16 text-gray-400 text-sm">No hay familias registradas</div>
                ) : (
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead>
                                <tr style={{ background: 'linear-gradient(135deg, #7c3aed 0%, #6d28d9 100%)' }}>
                                    <th className="px-4 py-3 text-left text-xs font-bold text-white uppercase tracking-wider w-8"></th>
                                    <th className="px-4 py-3 text-left text-xs font-bold text-white uppercase tracking-wider">Familia</th>
                                    <th className="px-4 py-3 text-left text-xs font-bold text-white uppercase tracking-wider hidden md:table-cell">Categoria</th>
                                    <th className="px-4 py-3 text-left text-xs font-bold text-white uppercase tracking-wider hidden md:table-cell">Marca</th>
                                    <th className="px-4 py-3 text-center text-xs font-bold text-white uppercase tracking-wider">Variantes</th>
                                    <th className="px-4 py-3 text-center text-xs font-bold text-white uppercase tracking-wider">Estado</th>
                                    <th className="px-4 py-3 text-right text-xs font-bold text-white uppercase tracking-wider">Acciones</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                                {families.map(family => (
                                    <>
                                        <tr key={family.id} className="hover:bg-purple-50 dark:hover:bg-gray-700/50 transition-colors">
                                            <td className="px-4 py-3">
                                                <button
                                                    onClick={() => handleToggleExpand(family.id)}
                                                    className="text-[#7c3aed] hover:text-[#6d28d9] transition-colors"
                                                    title="Ver variantes"
                                                >
                                                    {loadingVariants === family.id ? (
                                                        <span className="inline-block w-4 h-4 border-2 border-[#7c3aed] border-t-transparent rounded-full animate-spin" />
                                                    ) : expandedFamily === family.id ? (
                                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
                                                    ) : (
                                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" /></svg>
                                                    )}
                                                </button>
                                            </td>
                                            <td className="px-4 py-3">
                                                <div className="flex items-center gap-3">
                                                    {family.image_url ? (
                                                        <img src={family.image_url} alt={family.name} className="w-8 h-8 rounded-lg object-cover flex-shrink-0" />
                                                    ) : (
                                                        <div className="w-8 h-8 rounded-lg bg-purple-100 dark:bg-purple-900/30 flex items-center justify-center flex-shrink-0">
                                                            <svg className="w-4 h-4 text-[#7c3aed]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" /></svg>
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
                                            <td className="px-4 py-3 hidden md:table-cell">
                                                <span className="text-sm text-gray-600 dark:text-gray-300">{family.category || '-'}</span>
                                            </td>
                                            <td className="px-4 py-3 hidden md:table-cell">
                                                <span className="text-sm text-gray-600 dark:text-gray-300">{family.brand || '-'}</span>
                                            </td>
                                            <td className="px-4 py-3 text-center">
                                                <span className="inline-flex items-center justify-center w-8 h-8 rounded-full bg-purple-100 dark:bg-purple-900/30 text-[#7c3aed] text-sm font-bold">
                                                    {family.variant_count}
                                                </span>
                                            </td>
                                            <td className="px-4 py-3 text-center">
                                                <span className={`px-2 py-0.5 rounded-full text-xs font-semibold ${family.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'}`}>
                                                    {family.is_active ? 'Activa' : 'Inactiva'}
                                                </span>
                                            </td>
                                            <td className="px-4 py-3">
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

                                        {expandedFamily === family.id && (
                                            <tr key={`variants-${family.id}`} className="bg-purple-50/40 dark:bg-gray-700/20">
                                                <td colSpan={7} className="px-8 py-4">
                                                    {loadingVariants === family.id ? (
                                                        <div className="flex justify-center py-4">
                                                            <div className="w-5 h-5 border-2 border-[#7c3aed] border-t-transparent rounded-full animate-spin" />
                                                        </div>
                                                    ) : !variantsMap[family.id] || variantsMap[family.id].length === 0 ? (
                                                        <p className="text-center text-sm text-gray-400 py-2">Esta familia no tiene variantes aun</p>
                                                    ) : (
                                                        <div className="overflow-x-auto rounded-lg border border-purple-200 dark:border-purple-800/40">
                                                            <table className="w-full text-sm">
                                                                <thead>
                                                                    <tr className="bg-purple-100/80 dark:bg-purple-900/20">
                                                                        <th className="px-4 py-2 text-left text-xs font-bold text-purple-700 dark:text-purple-300 uppercase">SKU</th>
                                                                        <th className="px-4 py-2 text-left text-xs font-bold text-purple-700 dark:text-purple-300 uppercase">Nombre</th>
                                                                        <th className="px-4 py-2 text-left text-xs font-bold text-purple-700 dark:text-purple-300 uppercase hidden sm:table-cell">Variante</th>
                                                                        <th className="px-4 py-2 text-left text-xs font-bold text-purple-700 dark:text-purple-300 uppercase hidden sm:table-cell">Barcode</th>
                                                                        <th className="px-4 py-2 text-right text-xs font-bold text-purple-700 dark:text-purple-300 uppercase">Precio</th>
                                                                        <th className="px-4 py-2 text-center text-xs font-bold text-purple-700 dark:text-purple-300 uppercase">Stock</th>
                                                                    </tr>
                                                                </thead>
                                                                <tbody className="divide-y divide-purple-100 dark:divide-purple-800/20">
                                                                    {variantsMap[family.id].map(variant => (
                                                                        <tr key={variant.id} className="hover:bg-purple-50 dark:hover:bg-purple-900/10 transition-colors">
                                                                            <td className="px-4 py-2 font-mono text-xs text-gray-600 dark:text-gray-300">{variant.sku}</td>
                                                                            <td className="px-4 py-2 text-gray-900 dark:text-white">{variant.name}</td>
                                                                            <td className="px-4 py-2 hidden sm:table-cell">
                                                                                {variant.variant_label ? (
                                                                                    <span className="px-2 py-0.5 rounded bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300 text-xs">
                                                                                        {variant.variant_label}
                                                                                    </span>
                                                                                ) : (
                                                                                    <span className="text-gray-400 text-xs">-</span>
                                                                                )}
                                                                            </td>
                                                                            <td className="px-4 py-2 hidden sm:table-cell font-mono text-xs text-gray-500 dark:text-gray-400">
                                                                                {variant.barcode || '-'}
                                                                            </td>
                                                                            <td className="px-4 py-2 text-right text-gray-900 dark:text-white font-semibold">
                                                                                {new Intl.NumberFormat('es-CO', { style: 'currency', currency: variant.currency || 'COP', maximumFractionDigits: 0 }).format(variant.price)}
                                                                            </td>
                                                                            <td className="px-4 py-2 text-center">
                                                                                {stockBadge(variant.stock_quantity ?? variant.stock ?? 0)}
                                                                            </td>
                                                                        </tr>
                                                                    ))}
                                                                </tbody>
                                                            </table>
                                                        </div>
                                                    )}
                                                </td>
                                            </tr>
                                        )}
                                    </>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}

                {!loading && totalPages > 1 && (
                    <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between">
                        <span className="text-sm text-gray-500 dark:text-gray-400">
                            Pagina {page} de {totalPages}
                        </span>
                        <div className="flex gap-2">
                            <button
                                onClick={() => setPage(p => Math.max(1, p - 1))}
                                disabled={page === 1}
                                className="px-4 py-2 text-sm font-medium rounded-lg border-2 border-[#7c3aed]/40 text-[#7c3aed] disabled:opacity-40 disabled:cursor-not-allowed hover:border-[#7c3aed] hover:bg-purple-50 transition-all"
                            >
                                Anterior
                            </button>
                            <button
                                onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                                disabled={page === totalPages}
                                className="px-4 py-2 text-sm font-medium rounded-lg border-2 border-[#7c3aed]/40 text-[#7c3aed] disabled:opacity-40 disabled:cursor-not-allowed hover:border-[#7c3aed] hover:bg-purple-50 transition-all"
                            >
                                Siguiente
                            </button>
                        </div>
                    </div>
                )}
            </div>
        );
    }
);

ProductFamilyList.displayName = 'ProductFamilyList';
export default ProductFamilyList;
