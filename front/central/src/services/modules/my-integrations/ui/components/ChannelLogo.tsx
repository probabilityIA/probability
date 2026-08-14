'use client';

import { channelBrand } from '../../domain/types';
import { CARD_BORDER } from '../panel-theme';

interface ChannelLogoProps {
    url?: string;
    code?: string;
    size?: number;
}

export function ChannelLogo({ url, code, size = 20 }: ChannelLogoProps) {
    if (url) {
        return (
            <img
                src={url}
                alt=""
                className="flex-shrink-0 rounded border bg-white object-contain p-0.5"
                style={{ borderColor: CARD_BORDER, width: size, height: size }}
            />
        );
    }
    return (
        <span
            className="flex-shrink-0 rounded-full"
            style={{ backgroundColor: channelBrand(code).dot, width: 8, height: 8 }}
        />
    );
}
