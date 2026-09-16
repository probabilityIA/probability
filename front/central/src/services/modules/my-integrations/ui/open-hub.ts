import type { SyncEnvironment } from './sync-activity-context';

export const OPEN_INTEGRATIONS_HUB_EVENT = 'probability:open-integrations-hub';

const PENDING_TTL_MS = 15000;

export interface IntegrationsHubIntent {
    environment: SyncEnvironment | null;
}

let pending: { intent: IntegrationsHubIntent; at: number } | null = null;

export function queueIntegrationsHub(intent: IntegrationsHubIntent): void {
    pending = { intent, at: Date.now() };
}

export function requestIntegrationsHub(intent: IntegrationsHubIntent): void {
    queueIntegrationsHub(intent);
    if (typeof window !== 'undefined') window.dispatchEvent(new Event(OPEN_INTEGRATIONS_HUB_EVENT));
}

export function consumePendingIntegrationsHub(): IntegrationsHubIntent | null {
    if (!pending) return null;
    const { intent, at } = pending;
    pending = null;
    return Date.now() - at <= PENDING_TTL_MS ? intent : null;
}
