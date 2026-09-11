'use client';

import { useState, useEffect, forwardRef, useImperativeHandle, Fragment, ReactNode } from 'react';
import { ProductFamily, Product, GetFamiliesParams } from '../../domain/types';
import { getProductFamiliesAction, deleteProductFamilyAction, getFamilyVariantsAction } from '../../infra/actions';
import { ChevronRightIcon, ChevronDownIcon, XMarkIcon, PlusIcon, ScaleIcon, ArrowUpTrayIcon } from '@heroicons/react/24/outline';
import AssociateProductModal from './AssociateProductModal';
import FamilyDimensionsModal from './FamilyDimensionsModal';
import ImportFamilyVariantsModal from './ImportFamilyVariantsModal';

interface ProductFamilyListProps {
    onEdit?: (family: ProductFamily) => void;
    selectedBusinessId?: number;
    searchName?: string;
    advancedFilters?: Record<string, string | boolean>;
}

export interface ProductFamilyListHandle {
    refresh: () => void;
}

const MODAL_PAGE_SIZE = 10;

const ProductFamilyList = forwardRef<ProductFamilyListHandle, ProductFamilyListProps>(
    ({ onEdit, selectedBusinessId, searchName = '', advancedFilters = {} }, ref) => {
        const [families, setFamilies] = useState<ProductFamily[]>([]);
        const [loading, setLoading] = useState(true);
        const [error, setError] = useState<string | null>(null);
        const [page, setPage] = useState(1);
        const [totalPages, setTotalPages] = useState(1);
        const [total, setTotal] = useState(0);
        const [filters, setFilters] = useState<GetFamiliesParams>({ page: 1, page_size: 10 });

        const [modalFamily, setModalFamily] = useState<ProductFamily | null>(null);
        const [modalVariants, setModalVariants] = useState<Product[]>([]);
        const [modalLoading, setModalLoading] = useState(false);
        const [modalPage, setModalPage] = useState(1);
        const [showAssociateModal, setShowAssociateModal] = useState(false);
        const [showDimensionsModal, setShowDimensionsModal] = useState(false);
        const [showImportVariantsModal, setShowImportVariantsModal] = useState(false);
        const [expandedVariantId, setExpandedVariantId] = useState<string | null>(null);
        const [expandedGroupKey, setExpandedGroupKey] = useState<string | null>(null);
        const [collapsedMegaFamilies, setCollapsedMegaFamilies] = useState<Set<number>>(new Set());

        useImperativeHandle(ref, () => ({ refresh: fetchFamilies }));

        const fetchFamilies = async () => {
            setLoading(true);
            setError(null);
            try {
                const params: GetFamiliesParams = { ...filters };
                if (selectedBusinessId) params.business_id = selectedBusinessId;
                const res = await getProductFamiliesAction(params);
                if (res.success === false) {
                    setError(res.message || 'Error al cargar familias');
                    setFamilies([]);
                } else {
                    setFamilies(res.data || []);
                    setTotal(res.total || 0);
                    setTotalPages(res.total_pages || 1);
                    setPage(res.page || 1);
                }
            } catch (e: any) {
                setError(e.message || 'Error inesperado');
            } finally {
                setLoading(false);
            }
        };

        useEffect(() => {
            const avanzados: Record<string, string> = {};
            for (const [key, value] of Object.entries(advancedFilters)) {
                avanzados[key] = String(value);
            }
            setFilters(prev => ({
                ...prev,
                name: searchName || undefined,
                category: undefined,
                brand: undefined,
                status: undefined,
                has_variants: undefined,
                ...avanzados,
                page: 1,
            }));
        }, [searchName, advancedFilters]);

        useEffect(() => {
            setFilters(prev => ({ ...prev, page: 1 }));
        }, [selectedBusinessId]);

        useEffect(() => { fetchFamilies(); }, [filters, selectedBusinessId]);

        const handleDelete = async (family: ProductFamily) => {
            if (!confirm(`Eliminar la familia "${family.name}"? Esta acci\u00f3n no se puede deshacer.`)) return;
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

        const fetchModalVariants = async (familyId: number) => {
            setModalLoading(true);
            try {
                const res = await getFamilyVariantsAction(familyId, selectedBusinessId);
                setModalVariants(res.data || []);
            } finally {
                setModalLoading(false);
            }
        };

        const handleOpenModal = async (family: ProductFamily) => {
            setModalFamily(family);
            setModalPage(1);
            setModalVariants([]);
            setExpandedVariantId(null);
            setExpandedGroupKey(null);
            await fetchModalVariants(family.id);
        };

        const definedAxes: { key: string; label: string }[] = modalFamily?.variant_axes ?? [];
        const derivedAxes: { key: string; label: string }[] = (() => {
            const keys: string[] = [];
            for (const variant of modalVariants) {
                const attrs = variant.variant_attributes;
                if (attrs && typeof attrs === 'object') {
                    for (const key of Object.keys(attrs)) {
                        if (!keys.includes(key)) keys.push(key);
                    }
                }
            }
            return keys.map(key => ({ key, label: key.charAt(0).toUpperCase() + key.slice(1) }));
        })();
        const familyAxes = definedAxes.length > 0 ? definedAxes : derivedAxes;

        const Field = ({ label, value, children }: { label: string; value?: string | null; children?: ReactNode }) => (
            <div className="min-w-0">
                <p className="text-[10.5px] font-semibold uppercase tracking-wide text-gray-400 mb-0.5">{label}</p>
                {children ?? <p className="text-sm text-gray-800 dark:text-gray-200 truncate">{value || '-'}</p>}
            </div>
        );

        const statusBadge = (status?: string) => {
            const isActive = (status || '').toLowerCase() === 'active';
            return (
                <span className={`px-2 py-0.5 rounded-full text-xs font-semibold ${isActive ? 'bg-green-100 text-green-700' : 'bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300'}`}>
                    {isActive ? 'Activo' : status || 'Inactivo'}
                </span>
            );
        };

        const stockBadge = (qty: number) => {
            if (qty === 0) return <span className="px-2 py-0.5 rounded-full text-xs font-semibold bg-red-100 text-red-700">Agotado</span>;
            if (qty < 5) return <span className="px-2 py-0.5 rounded-full text-xs font-semibold bg-yellow-100 text-yellow-700">{qty}</span>;
            return <span className="px-2 py-0.5 rounded-full text-xs font-semibold bg-green-100 text-green-700">{qty}</span>;
        };

        const groupAxisKey = familyAxes.find(ax => ax.key.toLowerCase() === 'color')?.key ?? familyAxes[0]?.key;
        const restAxes = familyAxes.filter(ax => ax.key !== groupAxisKey);

        type VariantGroup = { key: string; label: string; variants: Product[] };
        const variantGroups: VariantGroup[] = (() => {
            if (!groupAxisKey) return [];
            const order: string[] = [];
            const byKey = new Map<string, VariantGroup>();
            for (const variant of modalVariants) {
                const value = variant.variant_attributes?.[groupAxisKey] || variant.variant_label || 'Sin variante';
                if (!byKey.has(value)) {
                    byKey.set(value, { key: value, label: value, variants: [] });
                    order.push(value);
                }
                byKey.get(value)!.variants.push(variant);
            }
            return order.map(key => byKey.get(key)!);
        })();

        const pagedVariants = modalVariants.slice((modalPage - 1) * MODAL_PAGE_SIZE, modalPage * MODAL_PAGE_SIZE);
        const modalTotalPages = Math.ceil(modalVariants.length / MODAL_PAGE_SIZE);
        const pagedGroups = variantGroups.slice((modalPage - 1) * MODAL_PAGE_SIZE, modalPage * MODAL_PAGE_SIZE);
        const groupTotalPages = Math.ceil(variantGroups.length / MODAL_PAGE_SIZE);

        const familyDescription = modalVariants.find(v => v.description)?.description;

        const pluralizeLabel = (label: string, count: number) => {
            if (count === 1) return label.toLowerCase();
            const lower = label.toLowerCase();
            return ['a', 'e', 'i', 'o', 'u'].includes(lower.slice(-1)) ? `${lower}s` : `${lower}es`;
        };

        const renderVariantRow = (variant: Product, axesForRow: { key: string; label: string }[], showLabelFallback: boolean, showMedia: boolean = true, showName: boolean = true) => {
            const isExpanded = expandedVariantId === variant.id;
            const colCount = 5 + (showName ? 1 : 0) + (axesForRow.length > 0 ? axesForRow.length : (showLabelFallback ? 1 : 0));
            return (
                <Fragment key={variant.id}>
                    <tr
                        onClick={() => setExpandedVariantId(isExpanded ? null : variant.id)}
                        className="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer"
                    >
                        <td>
                            <button
                                type="button"
                                onClick={(e) => { e.stopPropagation(); setExpandedVariantId(isExpanded ? null : variant.id); }}
                                className="p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 rounded transition-colors"
                                title={isExpanded ? 'Ocultar información' : 'Ver información'}
                            >
                                {isExpanded ? <ChevronDownIcon className="w-4 h-4" /> : <ChevronRightIcon className="w-4 h-4" />}
                            </button>
                        </td>
                        <td className="font-mono text-xs text-gray-600 dark:text-gray-300">{variant.sku}</td>
                        {showName && <td className="text-sm text-gray-900 dark:text-white">{variant.name}</td>}
                        {axesForRow.length > 0 ? (
                            axesForRow.map(ax => {
                                const value = variant.variant_attributes?.[ax.key];
                                return (
                                    <td key={ax.key} className="hidden sm:table-cell">
                                        {value ? (
                                            <span className="px-2 py-0.5 rounded text-xs font-medium bg-indigo-100 text-indigo-700">{value}</span>
                                        ) : (
                                            <span className="text-gray-400 text-xs">&mdash;</span>
                                        )}
                                    </td>
                                );
                            })
                        ) : showLabelFallback ? (
                            <td className="hidden sm:table-cell">
                                {variant.variant_label ? (
                                    <span className="px-2 py-0.5 rounded text-xs font-medium bg-indigo-100 text-indigo-700">{variant.variant_label}</span>
                                ) : (
                                    <span className="text-gray-400 text-xs">&mdash;</span>
                                )}
                            </td>
                        ) : null}
                        <td className="hidden sm:table-cell font-mono text-xs text-gray-500 dark:text-gray-400">{variant.barcode || '-'}</td>
                        <td className="text-right text-sm font-semibold text-gray-900 dark:text-white">
                            {new Intl.NumberFormat('es-CO', { style: 'currency', currency: variant.currency || 'COP', maximumFractionDigits: 0 }).format(variant.price)}
                        </td>
                        <td className="text-center">{stockBadge(variant.stock_quantity ?? variant.stock ?? 0)}</td>
                    </tr>
                    {isExpanded && (
                        <tr className="bg-gray-50 dark:bg-gray-900/30">
                            <td colSpan={colCount} className="px-5 py-4">
                                <div className="flex flex-col sm:flex-row gap-5">
                                    {showMedia && variant.image_url && (
                                        <img src={variant.image_url} alt={variant.name} className="w-24 h-24 rounded-lg object-cover flex-shrink-0 border border-gray-200 dark:border-gray-700" />
                                    )}
                                    <div className="flex-1 min-w-0">
                                        <div className="grid grid-cols-2 sm:grid-cols-4 gap-x-6 gap-y-4">
                                            <Field label="Categoría" value={variant.category} />
                                            <Field label="Marca" value={variant.brand} />
                                            <Field label="Estado">{statusBadge(variant.status)}</Field>
                                            <Field label="Gestiona inventario" value={variant.manage_stock ? 'Si' : 'No'} />
                                            <Field label="Peso" value={variant.weight != null ? `${variant.weight} kg` : undefined} />
                                            <Field
                                                label="Dimensiones (LxAxA)"
                                                value={
                                                    variant.length != null && variant.width != null && variant.height != null
                                                        ? `${variant.length} x ${variant.width} x ${variant.height} cm`
                                                        : undefined
                                                }
                                            />
                                            <Field label="Integración" value={variant.integration_type} />
                                            <Field label="Creado" value={variant.created_at ? new Date(variant.created_at).toLocaleDateString('es-CO') : undefined} />
                                        </div>
                                        {showMedia && variant.description && (
                                            <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
                                                <p className="text-[10.5px] font-semibold uppercase tracking-wide text-gray-400 mb-1">{'Descripción'}</p>
                                                <p className="text-sm text-gray-700 dark:text-gray-200 leading-relaxed">{variant.description}</p>
                                            </div>
                                        )}
                                    </div>
                                </div>
                            </td>
                        </tr>
                    )}
                </Fragment>
            );
        };

        const childrenByParent = new Map<number, ProductFamily[]>();
        for (const f of families) {
            if (f.parent_family_id) {
                if (!childrenByParent.has(f.parent_family_id)) childrenByParent.set(f.parent_family_id, []);
                childrenByParent.get(f.parent_family_id)!.push(f);
            }
        }
        const nestedChildIds = new Set(
            families.filter(f => f.parent_family_id && families.some(p => p.id === f.parent_family_id)).map(f => f.id)
        );
        type FamilyRow = { family: ProductFamily; isChild: boolean; isMegaParent: boolean; childCount: number; totalVariants: number; groupId: number };
        const familyRows: FamilyRow[] = [];
        for (const f of families) {
            if (nestedChildIds.has(f.id)) continue;
            const kids = childrenByParent.get(f.id) || [];
            const isMegaParent = kids.length > 0;
            const totalVariants = kids.reduce((sum, k) => sum + (k.variant_count || 0), 0);
            familyRows.push({ family: f, isChild: false, isMegaParent, childCount: kids.length, totalVariants, groupId: f.id });
            if (isMegaParent && !collapsedMegaFamilies.has(f.id)) {
                for (const kid of kids) {
                    familyRows.push({ family: kid, isChild: true, isMegaParent: false, childCount: 0, totalVariants: 0, groupId: f.id });
                }
            }
        }
        const toggleMegaFamily = (id: number) => {
            setCollapsedMegaFamilies(prev => {
                const next = new Set(prev);
                if (next.has(id)) next.delete(id); else next.add(id);
                return next;
            });
        };

        return (
            <>
                <div className="relative rounded-xl overflow-hidden shadow-sm bg-white dark:bg-gray-800">
                    <div className="overflow-x-auto">
                        <table className="min-w-full">
                            <thead style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-on-primary, white)' }}>
                                <tr>
                                    <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest">
                                        Familia
                                    </th>
                                    <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden md:table-cell">
                                        Categoria
                                    </th>
                                    <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden md:table-cell">
                                        Marca
                                    </th>
                                    <th className="px-3 sm:px-6 py-3 text-center text-xs font-bold text-white uppercase tracking-widest">
                                        Variantes
                                    </th>
                                    <th className="px-3 sm:px-6 py-3 text-center text-xs font-bold text-white uppercase tracking-widest">
                                        Estado
                                    </th>
                                    <th className="px-3 sm:px-6 py-3 text-right text-xs font-bold text-white uppercase tracking-widest">
                                        Acciones
                                    </th>
                                </tr>
                            </thead>
                            <tbody>
                                {loading ? (
                                    <tr>
                                        <td colSpan={6} className="px-4 sm:px-6 py-12 text-center text-gray-500 dark:text-gray-400">
                                            Cargando familias...
                                        </td>
                                    </tr>
                                ) : error ? (
                                    <tr>
                                        <td colSpan={6} className="px-4 sm:px-6 py-12 text-center text-red-500 text-sm">{error}</td>
                                    </tr>
                                ) : families.length === 0 ? (
                                    <tr>
                                        <td colSpan={6} className="px-4 sm:px-6 py-8 text-center text-gray-500 dark:text-gray-400">
                                            No hay familias disponibles
                                        </td>
                                    </tr>
                                ) : (
                                    familyRows.map(({ family, isChild, isMegaParent, childCount, totalVariants }, idx) => {
                                        const isPartOfGroup = isChild || isMegaParent;
                                        const isLastInGroup = isChild && (
                                            idx === familyRows.length - 1 || !familyRows[idx + 1].isChild
                                        );
                                        return (
                                        <tr
                                            key={family.id}
                                            className={`transition-colors ${isPartOfGroup ? 'bg-purple-50/50 dark:bg-purple-900/10 border-l-4 border-l-purple-300 dark:border-l-purple-700' : 'bg-white dark:bg-gray-800 border-l-4 border-l-transparent'} ${isLastInGroup ? 'border-b-2 border-b-purple-200 dark:border-b-purple-800' : 'border-b border-gray-100 dark:border-gray-700'} hover:bg-purple-100/60 dark:hover:bg-gray-700`}
                                        >
                                            <td className="px-3 sm:px-6 py-4">
                                                <div className={`flex items-center ${isChild ? 'pl-6' : ''}`}>
                                                    {isChild && (
                                                        <span className="text-purple-300 dark:text-purple-600 mr-2 select-none font-bold">&#8627;</span>
                                                    )}
                                                    {isMegaParent && (
                                                        <button
                                                            type="button"
                                                            onClick={() => toggleMegaFamily(family.id)}
                                                            className="p-1 mr-1 text-purple-500 hover:text-purple-700 dark:hover:text-purple-300 rounded transition-colors"
                                                            title={collapsedMegaFamilies.has(family.id) ? 'Ver subfamilias' : 'Ocultar subfamilias'}
                                                        >
                                                            {collapsedMegaFamilies.has(family.id) ? <ChevronRightIcon className="w-4 h-4" /> : <ChevronDownIcon className="w-4 h-4" />}
                                                        </button>
                                                    )}
                                                    {family.image_url ? (
                                                        <img src={family.image_url} alt={family.name} className={`h-10 w-10 rounded-full mr-3 object-cover ${isMegaParent ? 'ring-2 ring-purple-300 dark:ring-purple-600' : ''}`} />
                                                    ) : (
                                                        <div className={`h-10 w-10 rounded-full mr-3 bg-gray-100 dark:bg-gray-700 flex items-center justify-center flex-shrink-0 ${isMegaParent ? 'ring-2 ring-purple-300 dark:ring-purple-600' : ''}`}>
                                                            <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                                                            </svg>
                                                        </div>
                                                    )}
                                                    <div className="min-w-0">
                                                        <div className="flex items-center gap-2 flex-wrap">
                                                            <span className="text-sm font-medium text-gray-900 dark:text-white">{family.name}</span>
                                                            {isMegaParent && (
                                                                <span
                                                                    className="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wide bg-purple-600 text-white"
                                                                    title="Estas subfamilias se agrupan como UN SOLO producto al sincronizar con un canal de ventas (ej. WooCommerce)"
                                                                >
                                                                    &#128081; Familia madre &middot; {childCount} {childCount === 1 ? 'subfamilia' : 'subfamilias'}
                                                                </span>
                                                            )}
                                                        </div>
                                                        {family.title && family.title !== family.name && (
                                                            <div className="text-xs text-gray-500 dark:text-gray-400 truncate max-w-[220px]">{family.title}</div>
                                                        )}
                                                        {isMegaParent && (
                                                            <div className="text-[11px] text-purple-500 dark:text-purple-400 mt-0.5">
                                                                Se sincroniza como un solo producto variable ({totalVariants} variantes en total)
                                                            </div>
                                                        )}
                                                    </div>
                                                </div>
                                            </td>
                                            <td className="px-3 sm:px-6 py-4 hidden md:table-cell text-sm text-gray-900 dark:text-white">{family.category || '-'}</td>
                                            <td className="px-3 sm:px-6 py-4 hidden md:table-cell text-sm text-gray-900 dark:text-white">{family.brand || '-'}</td>
                                            <td className="px-3 sm:px-6 py-4 text-center whitespace-nowrap">
                                                {isMegaParent ? (
                                                    <span
                                                        className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-bold bg-purple-100 dark:bg-purple-900/40 text-purple-700 dark:text-purple-300"
                                                        title="Suma de variantes de todas las subfamilias agrupadas"
                                                    >
                                                        {totalVariants} agrupadas
                                                    </span>
                                                ) : (
                                                    <button
                                                        onClick={() => handleOpenModal(family)}
                                                        className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-bold bg-indigo-100 dark:bg-indigo-900/40 text-indigo-700 dark:text-indigo-300 hover:bg-indigo-200 dark:hover:bg-indigo-900/70 transition-colors"
                                                        title="Ver variantes"
                                                    >
                                                        {family.variant_count}
                                                        <ChevronRightIcon className="w-3 h-3" />
                                                    </button>
                                                )}
                                            </td>
                                            <td className="px-3 sm:px-6 py-4 text-center whitespace-nowrap">
                                                <span className={`inline-flex px-3 py-1 text-xs font-semibold rounded-full ${family.is_active
                                                    ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200'
                                                    : 'bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-200'
                                                    }`}>
                                                    {family.is_active ? 'Activa' : 'Inactiva'}
                                                </span>
                                            </td>
                                            <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                                <div className="flex flex-row justify-end gap-1.5">
                                                    <button
                                                        onClick={() => handleOpenModal(family)}
                                                        title="Ver variantes"
                                                        className="p-1.5 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors duration-200"
                                                    >
                                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" /></svg>
                                                    </button>
                                                    {onEdit && (
                                                        <button
                                                            onClick={() => onEdit(family)}
                                                            title="Editar"
                                                            className="p-1.5 bg-amber-500 hover:bg-amber-600 text-white rounded-md transition-colors duration-200"
                                                        >
                                                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                                                        </button>
                                                    )}
                                                    <button
                                                        onClick={() => handleDelete(family)}
                                                        title="Eliminar"
                                                        className="p-1.5 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors duration-200"
                                                    >
                                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                                                    </button>
                                                </div>
                                            </td>
                                        </tr>
                                        );
                                    })
                                )}
                            </tbody>
                        </table>
                    </div>

                    {!loading && !error && total > 0 && (
                        <div className="px-3 py-2 flex flex-col sm:flex-row items-center justify-between gap-2" style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-on-primary, white)' }}>
                            <div className="flex items-center gap-3">
                                <p className="text-[11px] text-white">
                                    Mostrando <span className="font-medium">{(page - 1) * (filters.page_size || 10) + 1}</span> a{' '}
                                    <span className="font-medium">{Math.min(page * (filters.page_size || 10), total)}</span> de{' '}
                                    <span className="font-medium">{total}</span> familias
                                </p>
                                <div className="flex items-center gap-1">
                                    <label className="text-[11px] text-white whitespace-nowrap">Mostrar:</label>
                                    <select
                                        value={filters.page_size || 10}
                                        onChange={(e) => setFilters({ ...filters, page_size: parseInt(e.target.value), page: 1 })}
                                        className="px-1.5 py-1 text-[11px] border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent text-gray-900 dark:text-white bg-white dark:bg-gray-800"
                                    >
                                        <option value="10">10</option>
                                        <option value="20">20</option>
                                        <option value="50">50</option>
                                        <option value="100">100</option>
                                    </select>
                                </div>
                            </div>
                            <div className="flex items-center gap-1">
                                <nav className="relative z-0 inline-flex items-center gap-1 flex-wrap justify-end">
                                    <button
                                        onClick={() => setFilters({ ...filters, page: 1 })}
                                        disabled={page === 1}
                                        className="page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all"
                                        title={'Primera p\u00e1gina'}
                                    >
                                        <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 19l-7-7 7-7m8 14l-7-7 7-7" /></svg>
                                    </button>

                                    <button
                                        onClick={() => setFilters({ ...filters, page: page - 1 })}
                                        disabled={page === 1}
                                        className="page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all"
                                        title="Anterior"
                                    >
                                        <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" /></svg>
                                    </button>

                                    {(() => {
                                        const pages: (number | string)[] = [];
                                        const maxVisible = 8;

                                        if (totalPages <= maxVisible) {
                                            for (let i = 1; i <= totalPages; i++) pages.push(i);
                                        } else {
                                            const windowSize = maxVisible - 2;
                                            let start = Math.max(2, page - Math.floor(windowSize / 2));
                                            let end = start + windowSize - 1;
                                            if (end >= totalPages) {
                                                end = totalPages - 1;
                                                start = Math.max(2, end - windowSize + 1);
                                            }

                                            pages.push(1);
                                            if (start > 2) pages.push('...');
                                            for (let i = start; i <= end; i++) pages.push(i);
                                            if (end < totalPages - 1) pages.push('...');
                                            pages.push(totalPages);
                                        }

                                        return pages.map((p, idx) =>
                                            typeof p === 'string' ? (
                                                <span key={`ellipsis-${idx}`} className="relative inline-flex items-center px-1 py-1 text-[11px] font-bold text-white/70">
                                                    ...
                                                </span>
                                            ) : (
                                                <button
                                                    key={p}
                                                    onClick={() => setFilters({ ...filters, page: p })}
                                                    className={`relative inline-flex items-center justify-center min-w-7 px-2 py-1 rounded-md border text-[11px] font-semibold transition-all ${p === page ? 'page-btn-active z-10 shadow-md scale-105' : 'page-btn'
                                                        }`}
                                                    style={p === page
                                                        ? { backgroundColor: 'var(--color-secondary)', borderColor: 'var(--color-secondary)', color: 'white' }
                                                        : undefined}
                                                >
                                                    {p}
                                                </button>
                                            )
                                        );
                                    })()}

                                    <button
                                        onClick={() => setFilters({ ...filters, page: page + 1 })}
                                        disabled={page === totalPages}
                                        className="page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all"
                                        title="Siguiente"
                                    >
                                        <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" /></svg>
                                    </button>

                                    <button
                                        onClick={() => setFilters({ ...filters, page: totalPages })}
                                        disabled={page === totalPages}
                                        className="page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all"
                                        title={'\u00daltima p\u00e1gina'}
                                    >
                                        <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 5l7 7-7 7M5 5l7 7-7 7" /></svg>
                                    </button>
                                </nav>
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
                                <div className="flex items-center gap-2">
                                    {modalVariants.length > 0 && (
                                        <button
                                            onClick={() => setShowDimensionsModal(true)}
                                            className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                                        >
                                            <ScaleIcon className="w-4 h-4" />
                                            Dimensiones
                                        </button>
                                    )}
                                    <button
                                        onClick={() => setShowImportVariantsModal(true)}
                                        className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                                    >
                                        <ArrowUpTrayIcon className="w-4 h-4" />
                                        Cargar productos
                                    </button>
                                    <button
                                        onClick={() => setShowAssociateModal(true)}
                                        className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium rounded-lg bg-[var(--color-primary)] text-white hover:opacity-90 transition-opacity"
                                    >
                                        <PlusIcon className="w-4 h-4" />
                                        Asociar producto
                                    </button>
                                    <button
                                        onClick={() => setModalFamily(null)}
                                        className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                                    >
                                        <XMarkIcon className="w-5 h-5" />
                                    </button>
                                </div>
                            </div>

                            <div className="overflow-auto flex-1 p-5 flex flex-col gap-4">
                                {!modalLoading && familyDescription && (
                                    <div className="rounded-lg border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/30 p-3">
                                        <p className="text-[10.5px] font-semibold uppercase tracking-wide text-gray-400 mb-1">{'Descripción'}</p>
                                        <p className="text-sm text-gray-700 dark:text-gray-200 leading-relaxed">{familyDescription}</p>
                                    </div>
                                )}
                                {modalLoading ? (
                                    <div className="flex justify-center py-12">
                                        <div className="spinner" />
                                    </div>
                                ) : modalVariants.length === 0 ? (
                                    <p className="text-center text-sm text-gray-400 py-12">Esta familia no tiene variantes aun.</p>
                                ) : (
                                    <>
                                        {groupAxisKey ? (
                                            <div className="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700">
                                                <table className="table w-full">
                                                    <thead>
                                                        <tr>
                                                            <th className="!bg-gray-100 dark:!bg-gray-700 w-8"></th>
                                                            <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-500 dark:!text-gray-400 text-left text-[11px] font-semibold uppercase tracking-wide">{familyAxes.find(ax => ax.key === groupAxisKey)?.label || 'Variante'}</th>
                                                            <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-500 dark:!text-gray-400 text-center text-[11px] font-semibold uppercase tracking-wide">Variantes</th>
                                                            <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-500 dark:!text-gray-400 text-center text-[11px] font-semibold uppercase tracking-wide">Stock total</th>
                                                        </tr>
                                                    </thead>
                                                    <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                                                        {pagedGroups.map((group) => {
                                                            const isGroupExpanded = expandedGroupKey === group.key;
                                                            const totalStock = group.variants.reduce((sum, v) => sum + (v.stock_quantity ?? v.stock ?? 0), 0);
                                                            return (
                                                                <Fragment key={group.key}>
                                                                    <tr
                                                                        onClick={() => setExpandedGroupKey(isGroupExpanded ? null : group.key)}
                                                                        className={`cursor-pointer transition-colors ${isGroupExpanded ? 'bg-indigo-50/60 dark:bg-indigo-900/20' : 'hover:bg-gray-50 dark:hover:bg-gray-700/50'}`}
                                                                    >
                                                                        <td className="py-2.5">
                                                                            <button
                                                                                type="button"
                                                                                onClick={(e) => { e.stopPropagation(); setExpandedGroupKey(isGroupExpanded ? null : group.key); }}
                                                                                className="p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 rounded transition-colors"
                                                                                title={isGroupExpanded ? 'Ocultar variantes' : 'Ver variantes'}
                                                                            >
                                                                                {isGroupExpanded ? <ChevronDownIcon className="w-4 h-4" /> : <ChevronRightIcon className="w-4 h-4" />}
                                                                            </button>
                                                                        </td>
                                                                        <td className="py-2.5 text-sm font-semibold text-gray-900 dark:text-white">
                                                                            <div className="flex items-center gap-3">
                                                                                {group.variants[0]?.image_url ? (
                                                                                    <img src={group.variants[0].image_url} alt={group.label} className="w-10 h-10 rounded-full object-cover flex-shrink-0 border-2 border-white dark:border-gray-800 shadow-sm ring-1 ring-gray-200 dark:ring-gray-600" />
                                                                                ) : (
                                                                                    <span className="w-10 h-10 rounded-full flex-shrink-0 bg-gray-100 dark:bg-gray-700 ring-1 ring-gray-200 dark:ring-gray-600" />
                                                                                )}
                                                                                {group.label}
                                                                            </div>
                                                                        </td>
                                                                        <td className="text-center text-sm text-gray-500 dark:text-gray-400">{group.variants.length}</td>
                                                                        <td className="text-center">{stockBadge(totalStock)}</td>
                                                                    </tr>
                                                                    {isGroupExpanded && (
                                                                        <tr className="bg-gray-50 dark:bg-gray-900/30">
                                                                            <td colSpan={4} className="p-0">
                                                                                <div className="overflow-x-auto">
                                                                                    <table className="table w-full">
                                                                                        <thead>
                                                                                            <tr>
                                                                                                <th className="!bg-gray-50 dark:!bg-gray-800 w-8"></th>
                                                                                                <th className="!bg-gray-50 dark:!bg-gray-800 !text-gray-500 dark:!text-gray-400 text-left">SKU</th>
                                                                                                {restAxes.length > 0 && restAxes.map(ax => (
                                                                                                    <th key={ax.key} className="!bg-gray-50 dark:!bg-gray-800 !text-gray-500 dark:!text-gray-400 text-left hidden sm:table-cell">{ax.label}</th>
                                                                                                ))}
                                                                                                <th className="!bg-gray-50 dark:!bg-gray-800 !text-gray-500 dark:!text-gray-400 text-left hidden sm:table-cell">Barcode</th>
                                                                                                <th className="!bg-gray-50 dark:!bg-gray-800 !text-gray-500 dark:!text-gray-400 text-right">Precio</th>
                                                                                                <th className="!bg-gray-50 dark:!bg-gray-800 !text-gray-500 dark:!text-gray-400 text-center">Stock</th>
                                                                                            </tr>
                                                                                        </thead>
                                                                                        <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                                                                                            {group.variants.map(variant => renderVariantRow(variant, restAxes, false, false, false))}
                                                                                        </tbody>
                                                                                    </table>
                                                                                </div>
                                                                            </td>
                                                                        </tr>
                                                                    )}
                                                                </Fragment>
                                                            );
                                                        })}
                                                    </tbody>
                                                </table>
                                            </div>
                                        ) : (
                                            <div className="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700">
                                                <table className="table w-full">
                                                    <thead>
                                                        <tr>
                                                            <th className="!bg-gray-100 dark:!bg-gray-700 w-8"></th>
                                                            <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-600 dark:!text-gray-300 text-left">SKU</th>
                                                            <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-600 dark:!text-gray-300 text-left">Nombre</th>
                                                            <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-600 dark:!text-gray-300 text-left hidden sm:table-cell">Variante</th>
                                                            <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-600 dark:!text-gray-300 text-left hidden sm:table-cell">Barcode</th>
                                                            <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-600 dark:!text-gray-300 text-right">Precio</th>
                                                            <th className="!bg-gray-100 dark:!bg-gray-700 !text-gray-600 dark:!text-gray-300 text-center">Stock</th>
                                                        </tr>
                                                    </thead>
                                                    <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                                                        {pagedVariants.map((variant) => renderVariantRow(variant, [], true))}
                                                    </tbody>
                                                </table>
                                            </div>
                                        )}

                                        {(groupAxisKey ? groupTotalPages : modalTotalPages) > 1 && (
                                            <div className="flex items-center justify-between px-1">
                                                <span className="text-xs text-gray-500">
                                                    {groupAxisKey ? `${variantGroups.length} ${pluralizeLabel(familyAxes.find(ax => ax.key === groupAxisKey)?.label || 'grupo', variantGroups.length)}` : `${modalVariants.length} variantes totales`}
                                                </span>
                                                <div className="flex items-center gap-1">
                                                    <button
                                                        disabled={modalPage <= 1}
                                                        onClick={() => setModalPage((p) => p - 1)}
                                                        className="px-3 py-1.5 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                                                    >
                                                        Anterior
                                                    </button>
                                                    <span className="text-xs text-gray-600 px-2">{modalPage} / {groupAxisKey ? groupTotalPages : modalTotalPages}</span>
                                                    <button
                                                        disabled={modalPage >= (groupAxisKey ? groupTotalPages : modalTotalPages)}
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

                {modalFamily && showAssociateModal && (
                    <AssociateProductModal
                        businessId={selectedBusinessId}
                        family={modalFamily}
                        onClose={() => setShowAssociateModal(false)}
                        onAssociated={() => fetchModalVariants(modalFamily.id)}
                    />
                )}

                {modalFamily && showDimensionsModal && (
                    <FamilyDimensionsModal
                        businessId={selectedBusinessId}
                        family={modalFamily}
                        variants={modalVariants}
                        onClose={() => setShowDimensionsModal(false)}
                        onApplied={() => fetchModalVariants(modalFamily.id)}
                    />
                )}

                {modalFamily && showImportVariantsModal && (
                    <ImportFamilyVariantsModal
                        businessId={selectedBusinessId}
                        family={modalFamily}
                        onClose={() => setShowImportVariantsModal(false)}
                        onAssociated={() => fetchModalVariants(modalFamily.id)}
                    />
                )}
            </>
        );
    }
);

ProductFamilyList.displayName = 'ProductFamilyList';
export default ProductFamilyList;
