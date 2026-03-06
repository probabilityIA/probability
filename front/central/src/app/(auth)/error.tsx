'use client';

import { useEffect } from 'react';

export default function AuthErrorBoundary({
    error,
    reset,
}: {
    error: Error & { digest?: string };
    reset: () => void;
}) {
    useEffect(() => {
        console.error('Auth layout error:', error);
    }, [error]);

    return (
        <div className="min-h-[50vh] flex items-center justify-center">
            <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-8 max-w-md w-full text-center space-y-4">
                <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto">
                    <svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                </div>
                <h2 className="text-lg font-semibold text-gray-900">
                    Algo salio mal
                </h2>
                <p className="text-sm text-gray-500">
                    Ocurrio un error inesperado. Por favor intenta de nuevo.
                </p>
                {error.digest && (
                    <p className="text-xs text-gray-400 font-mono">
                        Ref: {error.digest}
                    </p>
                )}
                <button
                    onClick={reset}
                    className="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors"
                >
                    Intentar de nuevo
                </button>
            </div>
        </div>
    );
}
