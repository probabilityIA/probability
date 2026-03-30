import { cookies } from 'next/headers';
import { MonitoringApiRepository } from '@/services/monitoring/infra/repository/api-repository';
import { ContainerGrid } from '@/services/monitoring/ui/components/ContainerGrid';
import { Header } from '@/services/monitoring/ui/components/Header';

export const dynamic = 'force-dynamic';

export default async function DashboardPage() {
    const cookieStore = await cookies();
    const token = cookieStore.get('monitoring_token')?.value;
    const userCookie = cookieStore.get('monitoring_user')?.value;
    const user = userCookie ? JSON.parse(userCookie) : null;

    let containers: Awaited<ReturnType<MonitoringApiRepository['listContainers']>> = [];
    let error = '';

    if (token) {
        try {
            const repo = new MonitoringApiRepository(token);
            containers = await repo.listContainers();
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to fetch containers';
        }
    }

    return (
        <div className="min-h-screen flex flex-col">
            <Header userName={user?.name} />

            <main className="flex-1 px-4 sm:px-6 py-6">
                <div className="max-w-7xl mx-auto">
                    <div className="flex items-center justify-between mb-5">
                        <h2 className="text-base font-semibold tracking-tight" style={{ color: '#e4e4ef' }}>
                            Containers
                        </h2>
                        <a
                            href="/dashboard"
                            className="text-xs px-3 py-1.5 rounded-md transition-all"
                            style={{ color: '#00f0ff', border: '1px solid #00f0ff25', background: '#00f0ff06' }}
                        >
                            Refresh
                        </a>
                    </div>

                    {error ? (
                        <div
                            className="rounded-xl p-4 text-sm"
                            style={{ background: '#ff336610', border: '1px solid #ff336630', color: '#ff3366' }}
                        >
                            {error}
                        </div>
                    ) : containers.length === 0 ? (
                        <div
                            className="rounded-xl p-8 text-center text-sm"
                            style={{ background: '#12121a', border: '1px solid #1e1e2e', color: '#8888a0' }}
                        >
                            No containers found
                        </div>
                    ) : (
                        <ContainerGrid containers={containers} />
                    )}
                </div>
            </main>
        </div>
    );
}
