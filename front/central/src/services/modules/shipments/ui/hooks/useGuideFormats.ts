'use client';

import { useEffect, useState } from 'react';
import { getGuideFormatsAction } from '../../infra/actions';
import { GuideFormat } from '../../domain/types';

const cache = new Map<string, GuideFormat[]>();

export function useGuideFormats(carrier?: string) {
    const [formats, setFormats] = useState<GuideFormat[]>(() => cache.get(carrier?.toUpperCase() || '*') || []);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        const key = carrier?.toUpperCase() || '*';
        if (cache.has(key)) {
            setFormats(cache.get(key)!);
            return;
        }
        let cancelled = false;
        setLoading(true);
        getGuideFormatsAction(carrier)
            .then((list) => {
                if (cancelled) return;
                cache.set(key, list);
                setFormats(list);
            })
            .finally(() => {
                if (!cancelled) setLoading(false);
            });
        return () => {
            cancelled = true;
        };
    }, [carrier]);

    const defaultFormat = formats.find((f) => f.is_default) || formats[0];
    return { formats, defaultFormat, loading };
}
