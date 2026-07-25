'use client';

interface CodIncludesShippingToggleProps {
    checked: boolean;
    onToggle: () => void;
}

export function CodIncludesShippingToggle({ checked, onToggle }: CodIncludesShippingToggleProps) {
    return (
        <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-3 mb-4">
            <div className="flex items-center justify-between gap-3">
                <div>
                    <span className="block text-[12px] font-semibold text-gray-900 dark:text-gray-100">
                        El total del pedido ya incluye el envio
                    </span>
                    <span className="block text-[11px] text-gray-500 dark:text-gray-400">
                        Activado: el canal ya le cobro el envio al cliente, asi que en contra entrega se recauda el total tal cual.
                        Desactivalo cuando las guias de contra entrega se generan en Probability y no llegan del canal de ventas:
                        ahi el total trae solo los productos y hay que sumarle el flete de la guia.
                    </span>
                </div>
                <button
                    type="button"
                    role="switch"
                    aria-checked={checked}
                    onClick={onToggle}
                    className={`relative inline-flex h-6 w-11 flex-shrink-0 items-center rounded-full transition-colors ${checked ? 'bg-[var(--color-primary)]' : 'bg-gray-300 dark:bg-gray-600'}`}
                >
                    <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${checked ? 'translate-x-6' : 'translate-x-1'}`} />
                </button>
            </div>
        </div>
    );
}
