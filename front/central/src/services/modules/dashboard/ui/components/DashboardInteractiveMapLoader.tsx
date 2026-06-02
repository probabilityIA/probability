import dynamic from 'next/dynamic';

const DashboardInteractiveMap = dynamic(() => import('./DashboardInteractiveMap'), {
    ssr: false,
    loading: () => (
        <div className="w-full h-96 flex items-center justify-center bg-gray-100 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
            <span className="text-gray-500 dark:text-gray-400">Cargando mapa...</span>
        </div>
    ),
});

export { DashboardInteractiveMap as default };
export { DashboardInteractiveMap };
