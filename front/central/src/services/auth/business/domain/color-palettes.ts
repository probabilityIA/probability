export interface BusinessPaletteColors {
    primary: string;
    secondary: string;
    tertiary: string;
    quaternary: string;
}

export interface BusinessColorPalette {
    name: string;
    colors: BusinessPaletteColors;
}

export const COLOR_PALETTES: BusinessColorPalette[] = [
    {
        name: 'Corporativo',
        colors: { primary: '#1E3A5F', secondary: '#FFFFFF', tertiary: '#3B82F6', quaternary: '#E5E7EB' },
    },
    {
        name: 'Moderno',
        colors: { primary: '#111827', secondary: '#F9FAFB', tertiary: '#6366F1', quaternary: '#E0E7FF' },
    },
    {
        name: 'Natural',
        colors: { primary: '#166534', secondary: '#FFFFFF', tertiary: '#22C55E', quaternary: '#DCFCE7' },
    },
    {
        name: 'Elegante',
        colors: { primary: '#1F2937', secondary: '#F3F4F6', tertiary: '#9333EA', quaternary: '#F3E8FF' },
    },
    {
        name: 'Cálido',
        colors: { primary: '#92400E', secondary: '#FFFBEB', tertiary: '#F59E0B', quaternary: '#FEF3C7' },
    },
    {
        name: 'Energético',
        colors: { primary: '#DC2626', secondary: '#FFFFFF', tertiary: '#F97316', quaternary: '#FEE2E2' },
    },
    {
        name: 'Oceánico',
        colors: { primary: '#0E7490', secondary: '#ECFEFF', tertiary: '#06B6D4', quaternary: '#CFFAFE' },
    },
    {
        name: 'Minimalista',
        colors: { primary: '#000000', secondary: '#FFFFFF', tertiary: '#737373', quaternary: '#F5F5F5' },
    },
    {
        name: 'Rosado',
        colors: { primary: '#BE185D', secondary: '#FDF2F8', tertiary: '#EC4899', quaternary: '#FCE7F3' },
    },
    {
        name: 'Tech',
        colors: { primary: '#7C3AED', secondary: '#0F172A', tertiary: '#A78BFA', quaternary: '#1E293B' },
    },
    {
        name: 'Bosque',
        colors: { primary: '#14532D', secondary: '#F0FDF4', tertiary: '#65A30D', quaternary: '#D9F99D' },
    },
    {
        name: 'Vino',
        colors: { primary: '#7F1D1D', secondary: '#FEF2F2', tertiary: '#B91C1C', quaternary: '#FCA5A5' },
    },
    {
        name: 'Medianoche',
        colors: { primary: '#0F172A', secondary: '#1E293B', tertiary: '#38BDF8', quaternary: '#334155' },
    },
    {
        name: 'Arena',
        colors: { primary: '#78350F', secondary: '#FFFBEB', tertiary: '#D97706', quaternary: '#FDE68A' },
    },
    {
        name: 'Menta',
        colors: { primary: '#115E59', secondary: '#F0FDFA', tertiary: '#14B8A6', quaternary: '#99F6E4' },
    },
    {
        name: 'Lavanda',
        colors: { primary: '#4C1D95', secondary: '#F5F3FF', tertiary: '#8B5CF6', quaternary: '#DDD6FE' },
    },
    {
        name: 'Coral',
        colors: { primary: '#9F1239', secondary: '#FFF1F2', tertiary: '#FB7185', quaternary: '#FECDD3' },
    },
    {
        name: 'Indigo',
        colors: { primary: '#312E81', secondary: '#EEF2FF', tertiary: '#4F46E5', quaternary: '#C7D2FE' },
    },
    {
        name: 'Grafito',
        colors: { primary: '#27272A', secondary: '#FAFAFA', tertiary: '#71717A', quaternary: '#E4E4E7' },
    },
    {
        name: 'Cobre',
        colors: { primary: '#7C2D12', secondary: '#FFF7ED', tertiary: '#EA580C', quaternary: '#FED7AA' },
    },
    {
        name: 'Cielo',
        colors: { primary: '#075985', secondary: '#F0F9FF', tertiary: '#0EA5E9', quaternary: '#BAE6FD' },
    },
    {
        name: 'Neon',
        colors: { primary: '#18181B', secondary: '#022C22', tertiary: '#22D3EE', quaternary: '#A3E635' },
    },
];
