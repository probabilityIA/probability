import type { Container, ContainerStats, LogEntry, ComposeService, LoginRequest, LoginResponse, SystemStats } from './types';

export interface IMonitoringRepository {
    login(data: LoginRequest): Promise<LoginResponse>;
    verifyToken(): Promise<{ valid: boolean; email: string; name: string }>;
    listContainers(): Promise<Container[]>;
    getContainer(id: string): Promise<Container>;
    getContainerStats(id: string): Promise<ContainerStats>;
    getContainerLogs(id: string, tail?: number): Promise<LogEntry[]>;
    restartContainer(id: string): Promise<{ message: string }>;
    stopContainer(id: string): Promise<{ message: string }>;
    startContainer(id: string): Promise<{ message: string }>;
    listComposeServices(): Promise<ComposeService[]>;
    getSystemStats(): Promise<SystemStats>;
}
