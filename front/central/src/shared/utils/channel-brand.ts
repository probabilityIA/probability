export interface ChannelBrand {
    chip: string;
    dot: string;
}

const NEUTRAL_BRAND: ChannelBrand = {
    chip: 'border-gray-200 bg-white text-gray-600 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300',
    dot: '#9ca3af',
};

const CHANNEL_BRANDS: Record<string, ChannelBrand> = {
    'mercado libre': {
        chip: 'border-amber-300 bg-amber-50 text-amber-800 dark:border-amber-500/50 dark:bg-amber-900/30 dark:text-amber-200',
        dot: '#ffe600',
    },
    woocommerce: {
        chip: 'border-purple-300 bg-purple-50 text-purple-800 dark:border-purple-500/50 dark:bg-purple-900/30 dark:text-purple-200',
        dot: '#7f54b3',
    },
    siigo: {
        chip: 'border-sky-300 bg-sky-50 text-sky-800 dark:border-sky-500/50 dark:bg-sky-900/30 dark:text-sky-200',
        dot: '#0ea5e9',
    },
    shopify: {
        chip: 'border-lime-300 bg-lime-50 text-lime-800 dark:border-lime-500/50 dark:bg-lime-900/30 dark:text-lime-200',
        dot: '#95bf47',
    },
    vtex: {
        chip: 'border-pink-300 bg-pink-50 text-pink-800 dark:border-pink-500/50 dark:bg-pink-900/30 dark:text-pink-200',
        dot: '#f71963',
    },
    jumpseller: {
        chip: 'border-orange-300 bg-orange-50 text-orange-800 dark:border-orange-500/50 dark:bg-orange-900/30 dark:text-orange-200',
        dot: '#f97316',
    },
};

export function channelBrand(channelCode?: string): ChannelBrand {
    if (!channelCode) return NEUTRAL_BRAND;
    return CHANNEL_BRANDS[channelCode.trim().toLowerCase()] ?? NEUTRAL_BRAND;
}
