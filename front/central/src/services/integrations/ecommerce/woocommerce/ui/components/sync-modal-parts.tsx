'use client';

export interface Brief {
    sku: string;
    name: string;
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
