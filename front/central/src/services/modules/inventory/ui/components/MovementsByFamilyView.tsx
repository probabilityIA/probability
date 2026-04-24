'use client';

import { useState, useEffect, useCallback } from 'react';
import { getProductFamiliesAction, getProductsAction } from '@/services/modules/products/infra/actions';
import { ProductFamilySummary, Product } from '@/services/modules/products/domain/types';
import { ChevronRightIcon, XMarkIcon, ArrowLeftIcon } from '@heroicons/react/24/outline';
import MovementsInlineTable from './MovementsInlineTable';

interface Props {
    businessId?: number;
}

const FAM_PAGE_SIZE = 15;

type ModalState =
    | { stage: 'products'; family: ProductFamilySummary; products: Product[]; loading: boolean }
    | { stage: 'movements'; family: ProductFamilySummary; product: Product };

export default function MovementsByFamilyView({ businessId }: Props) {
    const [families, setFamilies] = useState<ProductFamilySummary[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [modal, setModal] = useState<ModalState | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        try {
            const params: any = { page, page_size: FAM_PAGE_SIZE };
            if (businessId) params.business_id = businessId;
            const res = await getProductFamiliesAction(params);
            setFamilies((res as any).data ?? []);
            setTotal((res as any).total ?? 0);
            setTotalPages((res as any).total_pages ?? 1);
        } finally {
            setLoading(false);
        }
    }, [page, businessId]);

    useEffect(() => { load(); }, [load]);
    useEffect(() => { setPage(1); }, [businessId]);

    const openFamily = async (family: ProductFamilySummary) => {
        setModal({ stage: 'products', family, products: [], loading: true });
        try {
            const params: any = { page: 1, page_size: 100, family_id: family.id };
            if (businessId) params.business_id = businessId;
            const res = await getProductsAction(params);
            const products: Product[] = (res as any).data ?? [];
            setModal({ stage: 'products', family, products, loading: false });
        } catch {
            setModal({ stage: 'products', family, products: [], loading: false });
        }
    };

    const openProduct = (product: Product) => {
        if (!modal) return;
        setModal({ stage: 'movements', family: modal.family, product });
    };

    const backToProducts = () => {
        if (!modal) return;
        setModal({ stage: 'products', family: modal.family, products: modal.stage === 'movements' ? [] : (modal as any).products, loading: modal.stage === 'movements' });
        if (modal.stage === 'movements') {
            openFamily(modal.family);
        }
    };

    return (
        <>
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                <table className="table w-full">
                    <thead>
                        <tr>
                            <th className="text-left">Familia</th>
                            <th className="text-left">Categoria</th>
                            <th className="text-left">Marca</th>
                            <th className="text-center">Variantes</th>
                            <th className="text-center w-12"></th>
                        </tr>
                    </thead>
                    <tbody>
                        {loading ? (
                            <tr>
                                <td colSpan={5} className="py-12 text-center">
                                    <div className="flex justify-center items-center gap-3">
                                        <div className="spinner"></div>
                                        <span className="text-sm text-gray-500">Cargando...</span>
                                    </div>
                                </td>
                            </tr>
                        ) : families.length === 0 ? (
                            <tr>
                                <td colSpan={5} className="py-12 text-center text-sm text-gray-500">
                                    Sin familias de productos.
                                </td>
                            </tr>
                        ) : (
                            families.map((f) => (
                                <tr key={f.id} className="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
                                    <td className="font-medium text-gray-900 dark:text-white">{f.name}</td>
                                    <td className="text-sm text-gray-500">{f.category || <span className="text-gray-300">&mdash;</span>}</td>
                                    <td className="text-sm text-gray-500">{f.brand || <span className="text-gray-300">&mdash;</span>}</td>
                                    <td className="text-center">
                                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-700">
                                            {(f as any).variant_count ?? 0}
                                        </span>
                                    </td>
                                    <td className="text-center">
                                        <button
                                            onClick={() => openFamily(f)}
                                            className="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                                            title="Ver movimientos"
                                        >
                                            <ChevronRightIcon className="w-4 h-4" />
                                        </button>
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>

            {totalPages > 1 && (
                <div className="flex items-center justify-between px-1">
                    <span className="text-xs text-gray-500">{total} familias</span>
                    <div className="flex items-center gap-1">
                        <button
                            disabled={page <= 1}
                            onClick={() => setPage((p) => p - 1)}
                            className="px-3 py-1.5 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50"
                        >
                            Anterior
                        </button>
                        <span className="text-xs text-gray-600 px-2">{page} / {totalPages}</span>
                        <button
                            disabled={page >= totalPages}
                            onClick={() => setPage((p) => p + 1)}
                            className="px-3 py-1.5 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50"
                        >
                            Siguiente
                        </button>
                    </div>
                </div>
            )}

            {modal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                    <div className="absolute inset-0 bg-black/50" onClick={() => setModal(null)} />
                    <div className="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-5xl max-h-[90vh] flex flex-col">
                        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
                            <div className="flex items-center gap-3">
                                {modal.stage === 'movements' && (
                                    <button
                                        onClick={backToProducts}
                                        className="p-1.5 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-colors"
                                        title="Volver a productos"
                                    >
                                        <ArrowLeftIcon className="w-4 h-4" />
                                    </button>
                                )}
                                <div>
                                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                                        {modal.stage === 'products' ? modal.family.name : modal.product.name}
                                    </h2>
                                    <p className="text-sm text-gray-500 mt-0.5">
                                        {modal.stage === 'products'
                                            ? `Familia — selecciona un producto`
                                            : <span className="font-mono">{modal.product.sku} &mdash; {modal.family.name}</span>
                                        }
                                    </p>
                                </div>
                            </div>
                            <button
                                onClick={() => setModal(null)}
                                className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-colors"
                            >
                                <XMarkIcon className="w-5 h-5" />
                            </button>
                        </div>

                        <div className="overflow-auto flex-1 p-5">
                            {modal.stage === 'products' ? (
                                modal.loading ? (
                                    <div className="flex justify-center py-12">
                                        <div className="spinner"></div>
                                    </div>
                                ) : modal.products.length === 0 ? (
                                    <p className="text-center text-sm text-gray-400 py-12">Sin productos en esta familia.</p>
                                ) : (
                                    <table className="w-full text-left border-collapse">
                                        <thead>
                                            <tr className="border-b border-gray-200">
                                                <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50">Producto</th>
                                                <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50">SKU</th>
                                                <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50 text-center w-12"></th>
                                            </tr>
                                        </thead>
                                        <tbody className="divide-y divide-gray-100">
                                            {modal.products.map((p) => (
                                                <tr key={p.id} className="hover:bg-gray-50 transition-colors">
                                                    <td className="px-3 py-2 font-medium text-gray-900">{p.name}</td>
                                                    <td className="px-3 py-2 text-sm text-gray-500 font-mono">{p.sku}</td>
                                                    <td className="px-3 py-2 text-center">
                                                        <button
                                                            onClick={() => openProduct(p)}
                                                            className="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                                                        >
                                                            <ChevronRightIcon className="w-4 h-4" />
                                                        </button>
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                )
                            ) : (
                                <MovementsInlineTable productId={modal.product.id} businessId={businessId} />
                            )}
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}
