import { useEffect, useState } from 'react';
import { TokenStorage } from '@/shared/utils/token-storage';

interface BusinessColors {
  primary_color: string;
  secondary_color: string;
  tertiary_color: string;
  quaternary_color: string;
}

interface FormattedColors {
  primary_color: string;
  secondary_color: string;
  tertiary_color: string;
  quaternary_color: string;
}

export function useDynamicBusinessColors(businessId?: number) {
  const [colors, setColors] = useState<FormattedColors | null>(null);

  useEffect(() => {
    if (!businessId || businessId <= 0) {
      setColors(null);
      return;
    }

    const storedColors = TokenStorage.getBusinessColors();
    if (storedColors) {
      setColors({
        primary_color: storedColors.primary || '#a855f7',
        secondary_color: storedColors.secondary || '#9f5cf7',
        tertiary_color: storedColors.tertiary || '#c4b5fd',
        quaternary_color: storedColors.quaternary || '#f3e8ff',
      });
    } else {
      setColors(null);
    }
  }, [businessId]);

  return { colors };
}
