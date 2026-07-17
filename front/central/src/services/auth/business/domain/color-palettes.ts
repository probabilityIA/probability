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
];
