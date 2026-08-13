export const ACCENT = 'var(--color-primary)';
export const ACCENT_DARK = 'color-mix(in srgb, var(--color-primary) 85%, black)';
export const ACCENT_SOFT = 'color-mix(in srgb, var(--color-primary) 10%, white)';
export const ACCENT_BORDER = 'color-mix(in srgb, var(--color-primary) 25%, white)';
export const ACCENT_TINT = 'color-mix(in srgb, var(--color-primary) 6%, white)';

export const CARD_BG = '#fafafd';
export const CARD_BORDER = '#eceaf3';

export const cardStyle = {
    backgroundColor: CARD_BG,
    border: `1px solid ${CARD_BORDER}`,
};

export const inputCls =
    'w-full rounded-lg border bg-white px-3 py-2 text-[12px] text-gray-900 placeholder-gray-400 focus:border-[var(--color-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 dark:bg-gray-800 dark:text-white';

export const primaryButtonCls =
    'inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-[12px] font-semibold text-[var(--color-on-primary)] transition-colors disabled:cursor-not-allowed disabled:opacity-50';

export const ghostButtonCls =
    'inline-flex items-center gap-1.5 rounded-lg border bg-white px-3 py-1.5 text-[12px] font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-40 dark:bg-gray-800';
