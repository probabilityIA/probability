'use client';

import { ReactNode, forwardRef } from 'react';

interface IconActionButtonProps {
    label: string;
    icon: ReactNode;
    onClick?: () => void;
    variant?: 'primary' | 'secondary' | 'tertiary' | 'quaternary';
    disabled?: boolean;
    className?: string;
    tooltipAlign?: 'center' | 'left' | 'right';
}

const TOOLTIP_ALIGN: Record<string, string> = {
    center: 'left-1/2 -translate-x-1/2',
    left: 'left-0',
    right: 'right-0',
};

export const IconActionButton = forwardRef<HTMLButtonElement, IconActionButtonProps>(
    function IconActionButton({ label, icon, onClick, variant = 'primary', disabled, className = '', tooltipAlign = 'center' }, ref) {
        return (
            <span className="relative group inline-flex">
                <button
                    ref={ref}
                    type="button"
                    onClick={onClick}
                    disabled={disabled}
                    aria-label={label}
                    className={`btn-${variant} inline-flex items-center justify-center rounded-lg hover:shadow-md transition-all disabled:opacity-50 ${className}`}
                    style={{ width: '2rem', height: '2rem', padding: 0 }}
                >
                    {icon}
                </button>
                <span className={`pointer-events-none absolute ${TOOLTIP_ALIGN[tooltipAlign]} top-full mt-1.5 z-50 whitespace-nowrap rounded-md bg-gray-900 dark:bg-gray-700 px-2 py-1 text-[10px] font-semibold text-white opacity-0 shadow-lg transition-opacity duration-150 group-hover:opacity-100`}>
                    {label}
                </span>
            </span>
        );
    }
);
