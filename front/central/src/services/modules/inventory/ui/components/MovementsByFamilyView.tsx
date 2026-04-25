'use client';

import { useState, useEffect, useCallback } from 'react';
import { getProductFamiliesAction, getProductsAction } from '@/services/modules/products/infra/actions';
import { ProductFamilySummary, Product } from '@/services/modules/products/domain/types';
import { getMovementsAction } from '../../infra/actions';
import { StockMovement } from '../../domain/types';
import { ChevronRightIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { Spinner } from '@/shared/ui';

interface Props {
    businessId?: number;
}

interface FamilyWithCount extends ProductFamilySummary {
    movCount: number;
}

const DIRECTION_STYLES: Record<string, { bg: string; text: string; prefix: string }> = {
    in: { bg: 'bg-green-100', text: 'text-green-800', prefix: '+' },
    out: { bg: 'bg-red-100', text: 'text-red-800', prefix: '-' },
    reserve: { bg: 'bg-blue-100', text: 'text-blue-800', prefix: '' },
    release: { bg: 'bg-amber-100', text: 'text-amber-800', prefix: '' },
    confirm: { bg: 'bg-red-100', text: 'text-red-800', prefix: '-' },
    neutral: { bg: 'bg-gray-100', text: 'text-gray-800', prefix: '' },
};

function getDirectionFromCode(code: string): string {
    if (['inbound', 'return_stock', 'adjustment_in'].includes(code)) return 'in';
    if (['outbound', 'sale', 'adjustment_out'].includes(code)) return 'out';
    if (code === 'reserve') return 'reserve';
    if (code === 'release') return 'release';
    if (code === 'confirm_sale') return 'confirm';
    return 'neutral';
}

function parseNotes(notes: string) {
    const rsvMatch = notes?.match(/Reservado:\s*(-?\d+)/i);
    const liberadoMatch = notes?.match(/Liberado:\s*(-?\d+)/i);
    const dispMatch = notes?.match(/Disponible:\s*(\d+)\s*->\s*(\d+)/i);
    const rsvReleasedMatch = notes?.match(/Reserva liberada:\s*(\d+)/i);
    return {
        reserved: rsvMatch ? parseInt(rsvMatch[1], 10) : null,
        availPrev: dispMatch ? parseInt(dispMatch[1], 10) : null,
        availNew: dispMatch ? parseInt(dispMatch[2], 10) : null,
        rsvReleased: rsvReleasedMatch ? parseInt(rsvReleasedMatch[1], 10) : (liberadoMatch ? parseInt(liberadoMatch[1], 10) : null),
    };
}

const MODAL_PAGE = 15;

interface ModalState {
    family: FamilyWithCount;
    movements: StockMovement[];
    loading: boolean;
    page: number;
}

export default function MovementsByFamilyView({ businessId }: Props) {
    const [families, setFamilies] = useState<FamilyWithCount[]>([]);
    const [loading, setLoading] = useState(true);
    const [modal, setModal] = useState<ModalState | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        try {
            const movParams: any = { page: 1, page_size: 500 };
            if (businessId) movParams.business_id = businessId;
            const movRes = await getMovementsAction(movParams);
            const movements: StockMovement[] = movRes.data ?? [];
            const productIdsWithMov = new Set(movements.map((m) => m.product_id));

            const famParams: any = { page: 1, page_size: 100 };
            if (businessId) famParams.business_id = businessId;
            const famRes = await getProductFamiliesAction(famParams);
            const allFamilies: ProductFamilySummary[] = (famRes as any).data ?? [];

            const productsByFamily = await Promise.all(
                allFamilies.map((f) => {
                    const p: any = { page: 1, page_size: 100, family_id: f.id };
                    if (businessId) p.business_id = businessId;
                    return getProductsAction(p)
                        .then((r) => (r as any).data as Product[] ?? [])
                        .catch(() => [] as Product[]);
                })
            );

            const result: FamilyWithCount[] = allFamilies
                .map((f, i) => {
                    const movCount = productsByFamily[i].filter((p) => productIdsWithMov.has(p.id)).length;
                    return { ...f, movCount };
                })
                .filter((f) => f.movCount > 0)
                .sort((a, b) => b.movCount - a.movCount);

            setFamilies(result);
        } finally {
            setLoading(false);
        }
    }, [businessId]);

    useEffect(() => { load(); }, [load]);

    const openFamily = async (family: FamilyWithCount) => {
        setModal({ family, movements: [], loading: true, page: 1 });
        try {
            const params: any = { page: 1, page_size: 100, family_id: family.id };
            if (businessId) params.business_id = businessId;
            const prodRes = await getProductsAction(params);
            const products: Product[] = (prodRes as any).data ?? [];

            const allMovements = await Promise.all(
                products.map((p) =>
                    getMovementsAction({ product_id: p.id, page: 1, page_size: 200, business_id: businessId })
                        .then((r) => r.data ?? [])
                        .catch(() => [] as StockMovement[])
                )
            );

            const merged = allMovements
                .flat()
                .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());

            setModal({ family, movements: merged, loading: false, page: 1 });
        } catch {
            setModal({ family, movements: [], loading: false, page: 1 });
        }
    };

    const pagedMovements = modal
        ? modal.movements.slice((modal.page - 1) * MODAL_PAGE, modal.page * MODAL_PAGE)
        : [];
    const totalPages = modal ? Math.ceil(modal.movements.length / MODAL_PAGE) : 1;

    return (
        <>
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                <table className="table w-full">
                    <thead>
                        <tr>
                            <th className="text-left">Familia</th>
                            <th className="text-left">Categoria</th>
                            <th className="text-left">Marca</th>
                            <th className="text-center">Productos con mov.</th>
                            <th className="text-center w-12"></th>
                        </tr>
                    </thead>
                    <tbody>
                        {loading ? (
                            <tr><td colSpan={5} className="py-12 text-center"><div className="flex justify-center items-center gap-3"><div className="spinner"></div><span className="text-sm text-gray-500">Cargando...</span></div></td></tr>
                        ) : families.length === 0 ? (
                            <tr><td colSpan={5} className="py-12 text-center text-sm text-gray-500">Sin familias con movimientos.</td></tr>
                        ) : (
                            families.map((f) => (
                                <tr key={f.id} className="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
                                    <td className="font-medium text-gray-900 dark:text-white">{f.name}</td>
                                    <td className="text-sm text-gray-500">{f.category || <span className="text-gray-300">&mdash;</span>}</td>
                                    <td className="text-sm text-gray-500">{f.brand || <span className="text-gray-300">&mdash;</span>}</td>
                                    <td className="text-center">
                                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-700">{f.movCount}</span>
                                    </td>
                                    <td className="text-center">
                                        <button onClick={() => openFamily(f)} className="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors" title="Ver movimientos">
                                            <ChevronRightIcon className="w-4 h-4" />
                                        </button>
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>

            {modal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                    <div className="absolute inset-0 bg-black/50" onClick={() => setModal(null)} />
                    <div className="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-5xl max-h-[90vh] flex flex-col">
                        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
                            <div>
                                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{modal.family.name}</h2>
                                <p className="text-sm text-gray-500 mt-0.5">
                                    {modal.loading ? 'Cargando movimientos...' : `${modal.movements.length} movimientos`}
                                </p>
                            </div>
                            <button onClick={() => setModal(null)} className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-colors">
                                <XMarkIcon className="w-5 h-5" />
                            </button>
                        </div>

                        <div className="overflow-auto flex-1 p-5 flex flex-col gap-3">
                            {modal.loading ? (
                                <div className="flex justify-center py-12"><Spinner size="lg" /></div>
                            ) : modal.movements.length === 0 ? (
                                <p className="text-center text-sm text-gray-400 py-12">Sin movimientos.</p>
                            ) : (
                                <>
                                    <div className="overflow-x-auto">
                                        <table className="w-full text-left border-collapse">
                                            <thead>
                                                <tr className="border-b border-gray-200">
                                                    <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50">Fecha</th>
                                                    <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50">Producto</th>
                                                    <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50">Bodega</th>
                                                    <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50">Tipo</th>
                                                    <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50 text-center">Cantidad</th>
                                                    <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50 text-center">Stock Total</th>
                                                    <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50">Razon</th>
                                                </tr>
                                            </thead>
                                            <tbody className="divide-y divide-gray-100">
                                                {pagedMovements.map((m) => {
                                                    const direction = getDirectionFromCode(m.movement_type_code);
                                                    const style = DIRECTION_STYLES[direction] || DIRECTION_STYLES.neutral;
                                                    const isReservation = m.movement_type_code === 'reserve';
                                                    const isRelease = m.movement_type_code === 'release';
                                                    const parsed = (isReservation || isRelease) ? parseNotes(m.notes || '') : null;
                                                    return (
                                                        <tr key={m.id} className="hover:bg-gray-50 transition-colors">
                                                            <td className="px-3 py-2 text-xs text-gray-500 whitespace-nowrap">
                                                                {new Date(m.created_at).toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })}
                                                            </td>
                                                            <td className="px-3 py-2">
                                                                <span className="text-sm font-medium text-gray-900">{m.product_name || m.product_id}</span>
                                                                {m.product_sku && <span className="block text-xs text-gray-400 font-mono">{m.product_sku}</span>}
                                                            </td>
                                                            <td className="px-3 py-2 text-sm text-gray-600">{m.warehouse_name || '-'}</td>
                                                            <td className="px-3 py-2">
                                                                <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${style.bg} ${style.text}`}>
                                                                    {m.movement_type_name || m.movement_type_code}
                                                                </span>
                                                            </td>
                                                            <td className="px-3 py-2 text-center">
                                                                <span className={`text-sm font-semibold ${direction === 'in' ? 'text-green-700' : direction === 'out' || direction === 'confirm' ? 'text-red-700' : 'text-gray-600'}`}>
                                                                    {isReservation || isRelease ? <span className="text-gray-300 text-xs">&mdash;</span> : <>{style.prefix}{Math.abs(m.quantity)}</>}
                                                                </span>
                                                            </td>
                                                            <td className="px-3 py-2 text-center text-xs text-gray-500">
                                                                {m.previous_qty} &rarr; {m.new_qty}
                                                            </td>
                                                            <td className="px-3 py-2 text-sm text-gray-600 max-w-[160px] truncate">{m.reason}</td>
                                                        </tr>
                                                    );
                                                })}
                                            </tbody>
                                        </table>
                                    </div>

                                    {totalPages > 1 && (
                                        <div className="flex items-center justify-between px-2">
                                            <span className="text-xs text-gray-500">{modal.movements.length} movimientos totales</span>
                                            <div className="flex items-center gap-1">
                                                <button
                                                    disabled={modal.page <= 1}
                                                    onClick={() => setModal((prev) => prev ? { ...prev, page: prev.page - 1 } : prev)}
                                                    className="px-2 py-1 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50"
                                                >
                                                    Ant
                                                </button>
                                                <span className="text-xs text-gray-600 px-2">{modal.page} / {totalPages}</span>
                                                <button
                                                    disabled={modal.page >= totalPages}
                                                    onClick={() => setModal((prev) => prev ? { ...prev, page: prev.page + 1 } : prev)}
                                                    className="px-2 py-1 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50"
                                                >
                                                    Sig
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
