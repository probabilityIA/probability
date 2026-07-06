'use client';

import React from 'react';

export const GREEN = 'var(--color-primary)';
export const GREEN_DARK = 'color-mix(in srgb, var(--color-primary) 85%, black)';
export const GREEN_SOFT = 'color-mix(in srgb, var(--color-primary) 10%, white)';
export const GREEN_BORDER = 'color-mix(in srgb, var(--color-primary) 25%, white)';
export const CARD_BG = '#fafafd';
export const CARD_BORDER = '#eceaf3';
export const INPUT_BORDER = '#e9e9f0';

export const fieldLabel = 'block text-[13px] font-semibold text-gray-900 dark:text-gray-100 mb-1';
export const fieldHint = 'text-[11px] text-gray-400 dark:text-gray-500 mt-1 flex items-start gap-1';
export const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)]';

interface SectionCardProps {
    icon: React.ReactNode;
    title: string;
    children: React.ReactNode;
}

export function SectionCard({ icon, title, children }: SectionCardProps) {
    return (
        <div
            className="rounded-xl p-4 dark:bg-gray-800/60"
            style={{ backgroundColor: CARD_BG, border: `1px solid ${CARD_BORDER}` }}
        >
            <div className="flex items-center gap-2 mb-3">
                <span className="flex h-7 w-7 items-center justify-center rounded-md" style={{ backgroundColor: GREEN_SOFT }}>
                    {icon}
                </span>
                <h3 className="text-sm font-bold text-gray-900 dark:text-white">{title}</h3>
            </div>
            {children}
        </div>
    );
}

interface ToggleRowProps {
    icon: React.ReactNode;
    title: string;
    subtitle: string;
    checked: boolean;
    onToggle: () => void;
    disabled?: boolean;
}

export function ToggleRow({ icon, title, subtitle, checked, onToggle, disabled }: ToggleRowProps) {
    return (
        <div className="flex items-center justify-between gap-3 px-3 py-2.5">
            <div className="flex items-center gap-2.5 min-w-0">
                <span
                    className="flex h-8 w-8 items-center justify-center rounded-lg shrink-0"
                    style={{ backgroundColor: GREEN_SOFT }}
                >
                    {icon}
                </span>
                <div className="min-w-0">
                    <p className="text-[13px] font-semibold text-gray-800 dark:text-gray-100 leading-tight">{title}</p>
                    <p className="text-[11px] text-gray-500 dark:text-gray-400 leading-tight mt-0.5">{subtitle}</p>
                </div>
            </div>
            <button
                type="button"
                role="switch"
                aria-checked={checked}
                onClick={onToggle}
                disabled={disabled}
                className="relative inline-flex h-7 w-12 items-center rounded-full transition-colors focus:outline-none shrink-0 disabled:opacity-50"
                style={{ backgroundColor: checked ? GREEN : '#e5e7eb' }}
            >
                <span className={`inline-block h-5 w-5 transform rounded-full bg-white shadow-md transition-transform ${checked ? 'translate-x-6' : 'translate-x-1'}`} />
            </button>
        </div>
    );
}

export function Spinner({ className }: { className?: string }) {
    return (
        <svg className={className} xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
    );
}
