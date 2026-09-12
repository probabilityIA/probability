'use client';

import { useEffect, useState, type ReactNode } from 'react';
import { createPortal } from 'react-dom';

export const INTEGRATIONS_SUBHEADER_SLOT_ID = 'integrations-subheader-slot';
export const INTEGRATIONS_TOOLBAR_SLOT_ID = 'integrations-toolbar-slot';

function usePanelSlot(slotId: string) {
    const [slot, setSlot] = useState<HTMLElement | null>(null);

    useEffect(() => {
        setSlot(document.getElementById(slotId));
    }, [slotId]);

    return slot;
}

export function PanelSummary({ children }: { children: ReactNode }) {
    const slot = usePanelSlot(INTEGRATIONS_SUBHEADER_SLOT_ID);
    if (!slot) return null;
    return createPortal(children, slot);
}

export function PanelToolbar({ children }: { children: ReactNode }) {
    const slot = usePanelSlot(INTEGRATIONS_TOOLBAR_SLOT_ID);
    if (!slot) return null;
    return createPortal(children, slot);
}
