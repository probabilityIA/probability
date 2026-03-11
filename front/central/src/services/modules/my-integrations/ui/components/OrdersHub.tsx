'use client';

import { useState, useEffect } from 'react';

/**
 * Nodo central del diagrama de flujo: Gestión de Órdenes.
 * Círculo con efectos visuales reactivos.
 * A futuro mostrará métricas en tiempo real.
 */
export function OrdersHub() {
    const [pulse, setPulse] = useState(false);

    useEffect(() => {
        const interval = setInterval(() => {
            setPulse(true);
            setTimeout(() => setPulse(false), 1200);
        }, 3500);
        return () => clearInterval(interval);
    }, []);

    return (
        <div className="flex items-center justify-center py-1">
            <div className="relative">
                {/* Glow giratorio */}
                <div className="absolute inset-[-6px] rounded-full bg-conic-gradient opacity-30 blur-sm animate-hub-spin" />

                {/* Anillo de pulso expansivo */}
                <div
                    className={`absolute inset-[-8px] rounded-full border-2 border-indigo-400/50 transition-all duration-1000 ease-out ${
                        pulse ? 'scale-125 opacity-0' : 'scale-100 opacity-60'
                    }`}
                />

                {/* Círculo principal */}
                <div className="relative flex flex-col items-center justify-center w-28 h-28 rounded-full bg-white dark:bg-gray-800 border-[3px] border-indigo-500 dark:border-indigo-400 shadow-xl shadow-indigo-500/20">
                    <span className="text-2xl">📋</span>
                    <span className="text-[10px] font-bold text-gray-700 dark:text-gray-300 mt-1 leading-tight text-center px-2">
                        Gestión de Órdenes
                    </span>

                    {/* Indicador de actividad */}
                    <div className="absolute top-1.5 right-1.5 w-2.5 h-2.5 rounded-full bg-green-500 border-2 border-white dark:border-gray-800">
                        <div className="absolute inset-0 rounded-full bg-green-400 animate-ping" />
                    </div>
                </div>
            </div>

            <style jsx>{`
                @keyframes hub-spin {
                    from { transform: rotate(0deg); }
                    to { transform: rotate(360deg); }
                }
                .animate-hub-spin {
                    animation: hub-spin 6s linear infinite;
                }
                .bg-conic-gradient {
                    background: conic-gradient(
                        from 0deg,
                        #6366f1,
                        #8b5cf6,
                        #a78bfa,
                        #6366f1
                    );
                }
            `}</style>
        </div>
    );
}
