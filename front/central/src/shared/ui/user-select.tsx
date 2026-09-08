'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import { resolveAvatarUrl, userInitials } from '@/shared/utils/avatar-url';

export interface UserSelectOption {
    id: number;
    name: string;
    email?: string;
    avatar_url?: string;
}

type AvatarSize = 'xs' | 'sm' | 'md';

const AVATAR_SIZE: Record<AvatarSize, string> = {
    xs: 'h-5 w-5 text-[9px]',
    sm: 'h-6 w-6 text-[10px]',
    md: 'h-7 w-7 text-[11px]',
};

interface UserAvatarProps {
    name?: string | null;
    avatarUrl?: string | null;
    size?: AvatarSize;
    title?: string;
    decorative?: boolean;
}

export function UserAvatar({ name, avatarUrl, size = 'sm', title, decorative = false }: UserAvatarProps) {
    const url = resolveAvatarUrl(avatarUrl);
    const label = decorative ? undefined : title ?? (name || undefined);
    if (url) {
        return (
            <img
                src={url}
                alt={decorative ? '' : name || ''}
                title={label}
                aria-hidden={decorative || undefined}
                className={`${AVATAR_SIZE[size]} shrink-0 rounded-full object-cover ring-1 ring-gray-200 dark:ring-gray-600`}
            />
        );
    }
    return (
        <span
            title={label}
            aria-hidden={decorative || undefined}
            className={`${AVATAR_SIZE[size]} flex shrink-0 items-center justify-center rounded-full bg-violet-500 font-bold text-white`}
        >
            {userInitials(name)}
        </span>
    );
}

export function UnassignedAvatar({ size = 'sm', title = 'Sin asignar' }: { size?: AvatarSize; title?: string }) {
    return (
        <span
            title={title}
            className={`${AVATAR_SIZE[size]} flex shrink-0 items-center justify-center rounded-full border border-dashed border-gray-300 text-gray-400 dark:border-gray-600 dark:text-gray-500`}
        >
            <svg className="h-2.5 w-2.5" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.5 20.25a7.5 7.5 0 1115 0v.75H4.5v-.75z" />
            </svg>
        </span>
    );
}

interface UserSelectProps {
    value?: number | null;
    options: UserSelectOption[];
    onChange: (userId: number | null) => void;
    disabled?: boolean;
    emptyLabel?: string;
    showEmail?: boolean;
    avatarSize?: AvatarSize;
    className?: string;
    buttonClassName?: string;
    ariaLabel?: string;
    fallbackName?: string;
    fallbackAvatarUrl?: string;
}

const SEARCH_THRESHOLD = 7;

export function UserSelect({
    value,
    options,
    onChange,
    disabled = false,
    emptyLabel = 'Sin asignar',
    showEmail = false,
    avatarSize = 'sm',
    className = '',
    buttonClassName = '',
    ariaLabel = 'Asignado a',
    fallbackName,
    fallbackAvatarUrl,
}: UserSelectProps) {
    const [open, setOpen] = useState(false);
    const [query, setQuery] = useState('');
    const wrapRef = useRef<HTMLDivElement>(null);

    const selected = useMemo(
        () => (value ? options.find((u) => u.id === value) : undefined),
        [options, value]
    );

    const selectedName = selected?.name || (value ? fallbackName || 'Usuario ' + value : '');
    const selectedAvatar = selected?.avatar_url || (value ? fallbackAvatarUrl : undefined);

    const filtered = useMemo(() => {
        const q = query.trim().toLowerCase();
        if (!q) return options;
        return options.filter(
            (u) => u.name.toLowerCase().includes(q) || (u.email || '').toLowerCase().includes(q)
        );
    }, [options, query]);

    useEffect(() => {
        if (!open) return;
        const onDocClick = (e: MouseEvent) => {
            if (!wrapRef.current?.contains(e.target as Node)) setOpen(false);
        };
        const onKey = (e: KeyboardEvent) => {
            if (e.key === 'Escape') setOpen(false);
        };
        document.addEventListener('mousedown', onDocClick);
        document.addEventListener('keydown', onKey);
        return () => {
            document.removeEventListener('mousedown', onDocClick);
            document.removeEventListener('keydown', onKey);
        };
    }, [open]);

    useEffect(() => {
        if (!open) setQuery('');
    }, [open]);

    const pick = (id: number | null) => {
        setOpen(false);
        if (id !== (value ?? null)) onChange(id);
    };

    return (
        <div ref={wrapRef} className={`relative ${className}`}>
            <button
                type="button"
                disabled={disabled}
                aria-label={ariaLabel}
                aria-haspopup="listbox"
                aria-expanded={open}
                onClick={(e) => { e.stopPropagation(); setOpen((v) => !v); }}
                className={`flex w-full items-center gap-2 rounded-lg border border-gray-300 bg-white px-2.5 py-1.5 text-left text-gray-800 transition hover:border-violet-400 disabled:cursor-not-allowed disabled:opacity-60 dark:border-white/10 dark:bg-gray-950/60 dark:text-gray-100 dark:hover:border-violet-500 ${buttonClassName}`}
            >
                {selectedName ? (
                    <UserAvatar name={selectedName} avatarUrl={selectedAvatar} size={avatarSize} decorative />
                ) : (
                    <UnassignedAvatar size={avatarSize} title={emptyLabel} />
                )}
                <span className={`min-w-0 flex-1 truncate ${selectedName ? '' : 'text-red-500 dark:text-red-400'}`}>
                    {selectedName || emptyLabel}
                </span>
                <svg className="h-3 w-3 shrink-0 text-gray-400" fill="none" stroke="currentColor" strokeWidth={2.5} viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                </svg>
            </button>

            {open && (
                <div
                    role="listbox"
                    onClick={(e) => e.stopPropagation()}
                    className="absolute left-0 right-0 z-50 mt-1 overflow-hidden rounded-lg border border-gray-200 bg-white shadow-xl dark:border-white/12 dark:bg-gray-900"
                >
                    {options.length >= SEARCH_THRESHOLD && (
                        <div className="border-b border-gray-100 p-1.5 dark:border-white/8">
                            <input
                                autoFocus
                                value={query}
                                onChange={(e) => setQuery(e.target.value)}
                                placeholder="Buscar persona..."
                                className="block w-full rounded-md bg-gray-100 px-2 py-1.5 text-xs text-gray-800 outline-none placeholder:text-gray-400 dark:bg-gray-950/70 dark:text-gray-100"
                            />
                        </div>
                    )}

                    <div className="max-h-60 overflow-y-auto p-1">
                        <button
                            type="button"
                            role="option"
                            aria-selected={!value}
                            onClick={() => pick(null)}
                            className={`flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-xs transition hover:bg-gray-100 dark:hover:bg-white/5 ${!value ? 'bg-violet-50 dark:bg-violet-500/10' : ''}`}
                        >
                            <UnassignedAvatar size={avatarSize} title="" />
                            <span className="text-gray-500 dark:text-gray-400">{emptyLabel}</span>
                        </button>

                        {filtered.map((u) => (
                            <button
                                key={u.id}
                                type="button"
                                role="option"
                                aria-selected={u.id === value}
                                onClick={() => pick(u.id)}
                                className={`flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-xs transition hover:bg-gray-100 dark:hover:bg-white/5 ${u.id === value ? 'bg-violet-50 dark:bg-violet-500/10' : ''}`}
                            >
                                <UserAvatar name={u.name} avatarUrl={u.avatar_url} size={avatarSize} decorative />
                                <span className="min-w-0 flex-1">
                                    <span className="block truncate font-medium text-gray-800 dark:text-gray-100">{u.name}</span>
                                    {showEmail && u.email && (
                                        <span className="block truncate text-[10.5px] text-gray-400 dark:text-gray-500">{u.email}</span>
                                    )}
                                </span>
                            </button>
                        ))}

                        {filtered.length === 0 && (
                            <p className="px-2 py-2 text-xs text-gray-400 dark:text-gray-500">Sin resultados</p>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}
