'use client';

import { PaymentGatewayType } from '../../domain/types';

interface PaymentMethodGridProps {
    gateways: PaymentGatewayType[];
    loading: boolean;
    onSelect: (gateway: PaymentGatewayType) => void;
}

export function PaymentMethodGrid({ gateways, loading, onSelect }: PaymentMethodGridProps) {
    if (loading) {
        return (
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
                {[...Array(6)].map((_, i) => (
                    <div key={i} className="h-24 bg-gray-100 rounded-xl animate-pulse" />
                ))}
            </div>
        );
    }

    if (gateways.length === 0) {
        return (
            <p className="text-center text-gray-500 text-sm py-6">
                No hay métodos de pago disponibles.
            </p>
        );
    }

    return (
        <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
            {gateways.map((gateway) => {
                const disabled = gateway.in_development;

                return (
                    <button
                        key={gateway.id}
                        onClick={() => !disabled && onSelect(gateway)}
                        disabled={disabled}
                        className={`
                            flex flex-col items-center justify-center gap-2 p-4 rounded-xl border transition-all
                            ${disabled
                                ? 'opacity-50 cursor-not-allowed bg-gray-50 border-gray-200'
                                : 'cursor-pointer bg-white border-gray-200 hover:border-purple-400 hover:shadow-md hover:scale-[1.02] active:scale-[0.98]'
                            }
                        `}
                    >
                        {gateway.image_url ? (
                            <img
                                src={gateway.image_url}
                                alt={gateway.name}
                                className="h-10 w-auto object-contain"
                                onError={(e) => {
                                    (e.target as HTMLImageElement).style.display = 'none';
                                }}
                            />
                        ) : (
                            <div className="h-10 w-10 bg-gray-200 rounded-lg flex items-center justify-center">
                                <span className="text-gray-500 text-xs font-bold">
                                    {gateway.name.slice(0, 2).toUpperCase()}
                                </span>
                            </div>
                        )}

                        <span className="text-xs font-semibold text-gray-800 text-center leading-tight">
                            {gateway.name}
                        </span>

                        <span className={`
                            text-[10px] font-medium px-2 py-0.5 rounded-full
                            ${disabled
                                ? 'bg-gray-100 text-gray-500'
                                : 'bg-green-100 text-green-700'
                            }
                        `}>
                            {disabled ? 'Próximamente' : 'Disponible'}
                        </span>
                    </button>
                );
            })}
        </div>
    );
}
