import type { IMonitoringRepository } from '../../domain/ports';
import type { Container, ContainerStats, LogEntry, ComposeService, LoginRequest, LoginResponse, SystemStats } from '../../domain/types';
import { apiFetch } from '@/shared/lib/api';

export class MonitoringApiRepository implements IMonitoringRepository {
    constructor(private token?: string) {}

    async login(data: LoginRequest): Promise<LoginResponse> {
        return apiFetch<LoginResponse>('/api/v1/auth/login', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async verifyToken(): Promise<{ valid: boolean; email: string; name: string }> {
        return apiFetch('/api/v1/auth/verify', { token: this.token });
    }

    async listContainers(): Promise<Container[]> {
        return apiFetch<Container[]>('/api/v1/containers', { token: this.token });
    }

    async getContainer(id: string): Promise<Container> {
        return apiFetch<Container>(`/api/v1/containers/${id}`, { token: this.token });
    }

    async getContainerStats(id: string): Promise<ContainerStats> {
        return apiFetch<ContainerStats>(`/api/v1/containers/${id}/stats`, { token: this.token });
    }

    async getContainerLogs(id: string, tail = 100): Promise<LogEntry[]> {
        return apiFetch<LogEntry[]>(`/api/v1/containers/${id}/logs?tail=${tail}`, { token: this.token });
    }

    async restartContainer(id: string): Promise<{ message: string }> {
        return apiFetch(`/api/v1/containers/${id}/restart`, { method: 'POST', token: this.token });
    }

    async stopContainer(id: string): Promise<{ message: string }> {
        return apiFetch(`/api/v1/containers/${id}/stop`, { method: 'POST', token: this.token });
    }

    async startContainer(id: string): Promise<{ message: string }> {
        return apiFetch(`/api/v1/containers/${id}/start`, { method: 'POST', token: this.token });
    }

    async listComposeServices(): Promise<ComposeService[]> {
        return apiFetch<ComposeService[]>('/api/v1/compose/services', { token: this.token });
    }

    async getSystemStats(): Promise<SystemStats> {
        return apiFetch<SystemStats>('/api/v1/system/stats', { token: this.token });
    }
}
