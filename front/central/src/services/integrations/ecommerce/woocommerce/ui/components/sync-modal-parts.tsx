'use client';

export interface Brief {
    sku: string;
    name: string;
    family_ref?: string;
    family_name?: string;
    variant_label?: string;
    image_url?: string;
}

export interface ChannelFamilyGroup {
    key: string;
    name: string;
    image_url?: string;
    items: Brief[];
}

export function groupByFamily(items: Brief[]): ChannelFamilyGroup[] {
    const order: string[] = [];
    const byKey = new Map<string, ChannelFamilyGroup>();
    for (const item of items) {
        const key = item.family_ref || `sku:${item.sku}`;
        if (!byKey.has(key)) {
            byKey.set(key, { key, name: item.family_name || item.name, image_url: item.image_url, items: [] });
            order.push(key);
        }
        const group = byKey.get(key)!;
        group.items.push(item);
        if (!group.image_url && item.image_url) group.image_url = item.image_url;
    }
    return order.map(key => byKey.get(key)!);
}

export type SyncAction = 'created' | 'updated' | 'failed' | 'skipped';

export function ActionBadge({ action }: { action: SyncAction }) {
    const map = {
        created: { label: 'Creado', cls: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300' },
        updated: { label: 'Actualizado', cls: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300' },
        failed: { label: 'Fallido', cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300' },
        skipped: { label: 'Omitido', cls: 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300' },
    };
    const { label, cls } = map[action];
    return <span className={`px-1.5 py-0.5 rounded font-semibold ${cls}`}>{label}</span>;
}

export function ProductList({ items }: { items: Brief[] }) {
    if (items.length === 0) return null;
    return (
        <div className="mt-2 max-h-32 overflow-y-auto rounded-md bg-gray-50 dark:bg-gray-800/60 divide-y divide-gray-100 dark:divide-gray-700">
            {items.slice(0, 100).map((p, i) => (
                <div key={i} className="flex items-center justify-between px-2.5 py-1.5 text-[11px]">
                    <span className="text-gray-700 dark:text-gray-200 truncate">{p.name || '(sin nombre)'}</span>
                    <span className="text-gray-400 font-mono ml-2 flex-shrink-0">{p.sku}</span>
                </div>
            ))}
            {items.length > 100 && <div className="px-2.5 py-1.5 text-[11px] text-gray-400">y {items.length - 100} más...</div>}
        </div>
    );
}

export function ChannelFamilyList({ groups, onImport }: { groups: ChannelFamilyGroup[]; onImport: (group: ChannelFamilyGroup) => void }) {
    if (groups.length === 0) return null;
    return (
        <div className="mt-2 max-h-56 overflow-y-auto rounded-md bg-gray-50 dark:bg-gray-800/60 divide-y divide-gray-100 dark:divide-gray-700">
            {groups.slice(0, 100).map((g) => (
                <div key={g.key} className="flex items-center justify-between gap-2 px-2.5 py-2 text-[11px]">
                    <div className="flex items-center gap-2 min-w-0">
                        {g.image_url ? (
                            <img src={g.image_url} alt={g.name} className="w-8 h-8 rounded-full object-cover flex-shrink-0 border border-gray-200 dark:border-gray-700" />
                        ) : (
                            <span className="w-8 h-8 rounded-full flex-shrink-0 bg-gray-200 dark:bg-gray-700" />
                        )}
                        <div className="min-w-0">
                            <p className="text-gray-700 dark:text-gray-200 truncate font-medium">{g.name || '(sin nombre)'}</p>
                            <p className="text-gray-400">{g.items.length} variante{g.items.length !== 1 ? 's' : ''}</p>
                        </div>
                    </div>
                    <button
                        onClick={() => onImport(g)}
                        className="inline-flex items-center gap-1 whitespace-nowrap rounded-md bg-blue-600 hover:bg-blue-700 px-2.5 py-1 text-[11px] font-semibold text-white transition-colors flex-shrink-0"
                    >
                        Traer a Probability
                    </button>
                </div>
            ))}
            {groups.length > 100 && <div className="px-2.5 py-1.5 text-[11px] text-gray-400">y {groups.length - 100} más...</div>}
        </div>
    );
}

export function SelectableProductList({ items, selected, onToggle }: { items: Brief[]; selected: Set<string>; onToggle: (sku: string) => void }) {
    if (items.length === 0) return null;
    return (
        <div className="mt-2 max-h-40 overflow-y-auto rounded-md bg-white dark:bg-gray-800/60 border border-amber-100 dark:border-amber-900/40 divide-y divide-gray-100 dark:divide-gray-700">
            {items.slice(0, 200).map((p, i) => (
                <label key={i} className="flex items-center gap-2 px-2.5 py-1.5 text-[11px] cursor-pointer hover:bg-amber-50/50 dark:hover:bg-amber-900/10">
                    <input
                        type="checkbox"
                        checked={selected.has(p.sku)}
                        onChange={() => onToggle(p.sku)}
                        className="h-3.5 w-3.5 rounded border-gray-300 text-amber-600 focus:ring-amber-500"
                    />
                    <span className="text-gray-700 dark:text-gray-200 truncate flex-1">{p.name || '(sin nombre)'}</span>
                    <span className="text-gray-400 font-mono ml-2 flex-shrink-0">{p.sku}</span>
                </label>
            ))}
            {items.length > 200 && <div className="px-2.5 py-1.5 text-[11px] text-gray-400">y {items.length - 200} mas (usa "Asociar todos")...</div>}
        </div>
    );
}
