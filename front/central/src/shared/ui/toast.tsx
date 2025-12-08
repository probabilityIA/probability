import React, { useEffect } from 'react';

export type ToastType = 'success' | 'info' | 'warning' | 'error';

export interface ToastProps {
    id: string;
    message: string;
    type: ToastType;
    duration?: number;
    onClose: (id: string) => void;
}

export const Toast: React.FC<ToastProps> = ({ id, message, type, duration = 6000, onClose }) => {
    useEffect(() => {
        const timer = setTimeout(() => {
            onClose(id);
        }, duration);

        return () => clearTimeout(timer);
    }, [id, duration, onClose]);

    const bgColors = {
        success: 'bg-white border-l-4 border-green-500 text-gray-800',
        info: 'bg-white border-l-4 border-blue-500 text-gray-800',
        warning: 'bg-white border-l-4 border-yellow-500 text-gray-800',
        error: 'bg-white border-l-4 border-red-500 text-gray-800',
    };

    const icons = {
        success: (
            <div className="flex items-center justify-center w-8 h-8 rounded-full bg-green-100">
                <svg className="w-5 h-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                </svg>
            </div>
        ),
        info: (
            <div className="flex items-center justify-center w-8 h-8 rounded-full bg-blue-100">
                <svg className="w-5 h-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
            </div>
        ),
        warning: (
            <div className="flex items-center justify-center w-8 h-8 rounded-full bg-yellow-100">
                <svg className="w-5 h-5 text-yellow-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
            </div>
        ),
        error: (
            <div className="flex items-center justify-center w-8 h-8 rounded-full bg-red-100">
                <svg className="w-5 h-5 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
            </div>
        ),
    };

    return (
        <div className={`flex items-start w-full max-w-md p-4 mb-4 rounded-lg shadow-lg border ${bgColors[type]} transition-all duration-300 ease-in-out transform translate-y-0 opacity-100`}>
            <div className="flex-shrink-0 pt-0.5">
                {icons[type]}
            </div>
            <div className="ml-3 flex-1">
                <p className="text-sm font-medium text-gray-900">
                    Notificación
                </p>
                <p className="mt-1 text-sm text-gray-500">
                    {message}
                </p>
            </div>
            <button
                type="button"
                className={`ml-4 -mx-1.5 -my-1.5 rounded-lg focus:ring-2 p-1.5 inline-flex h-8 w-8 text-gray-400 hover:text-gray-900 hover:bg-gray-100 focus:ring-gray-300`}
                onClick={() => onClose(id)}
                aria-label="Close"
            >
                <span className="sr-only">Close</span>
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                    <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd"></path>
                </svg>
            </button>
        </div>
    );
};
