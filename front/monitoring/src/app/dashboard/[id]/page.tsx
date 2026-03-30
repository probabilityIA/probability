import { cookies } from 'next/headers';
import Link from 'next/link';
import { MonitoringApiRepository } from '@/services/monitoring/infra/repository/api-repository';
import { ContainerDetail } from '@/services/monitoring/ui/components/ContainerDetail';
import { ActionButtons } from '@/services/monitoring/ui/components/ActionButtons';
import { LogViewer } from '@/services/monitoring/ui/components/LogViewer';
import { StatsBar } from '@/services/monitoring/ui/components/StatsBar';
import { Header } from '@/services/monitoring/ui/components/Header';

export const dynamic = 'force-dynamic';

interface PageProps {
    params: Promise<{ id: string }>;
}

export default async function ContainerDetailPage({ params }: PageProps) {
    const { id } = await params;
    const cookieStore = await cookies();
    const token = cookieStore.get('monitoring_token')?.value;
    const userCookie = cookieStore.get('monitoring_user')?.value;
    const user = userCookie ? JSON.parse(userCookie) : null;

    let container = null;
    let error = '';

    if (token) {
        try {
            const repo = new MonitoringApiRepository(token);
            container = await repo.getContainer(id);
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to fetch container';
        }
    }

    return (
        <div className="min-h-screen flex flex-col">
            <Header userName={user?.name} />

            <main className="flex-1 px-6 py-6">
                <div className="max-w-7xl mx-auto space-y-4">
                    {/* Back link */}
                    <Link
                        href="/dashboard"
                        className="inline-flex items-center gap-1.5 text-xs transition-colors"
                        style={{ color: '#8888a0' }}
                    >
                        <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
                        </svg>
                        Back to containers
                    </Link>

                    {error ? (
                        <div
                            className="rounded-xl p-4 text-sm"
                            style={{ background: '#ff336610', border: '1px solid #ff336630', color: '#ff3366' }}
                        >
                            {error}
                        </div>
                    ) : container ? (
                        <>
                            {/* Container info + actions + stats */}
                            <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
                                <div className="lg:col-span-2">
                                    <ContainerDetail container={container} />
                                </div>
                                <div className="space-y-3">
                                    <ActionButtons containerId={container.id} state={container.state} />
                                    <StatsBar containerId={container.id} state={container.state} />
                                </div>
                            </div>

                            {/* Logs */}
                            <div>
                                <h3 className="text-sm font-medium mb-2" style={{ color: '#e4e4ef' }}>
                                    Live Logs
                                </h3>
                                <div style={{ height: 'calc(100vh - 380px)', minHeight: '300px' }}>
                                    <LogViewer containerId={container.id} />
                                </div>
                            </div>
                        </>
                    ) : (
                        <div
                            className="rounded-xl p-8 text-center text-sm"
                            style={{ background: '#12121a', border: '1px solid #1e1e2e', color: '#8888a0' }}
                        >
                            Container not found
                        </div>
                    )}
                </div>
            </main>
        </div>
    );
}
