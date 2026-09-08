const S3_BASE_URL = process.env.NEXT_PUBLIC_S3_BASE_URL || 'https://probability-media-assets.s3.us-east-1.amazonaws.com';

export const resolveAvatarUrl = (avatarUrl?: string | null): string => {
    if (!avatarUrl) return '';
    return avatarUrl.startsWith('http') ? avatarUrl : `${S3_BASE_URL}/${avatarUrl.replace(/^\//, '')}`;
};

export const userInitials = (name?: string | null): string => {
    if (!name || !name.trim()) return '?';
    const parts = name.trim().split(/\s+/);
    const first = parts[0]?.[0] || '';
    const last = parts.length > 1 ? parts[parts.length - 1][0] : '';
    return (first + last).toUpperCase();
};
