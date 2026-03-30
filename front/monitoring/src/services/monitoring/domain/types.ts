export interface PortMapping {
    host_port: number;
    container_port: number;
    protocol: string;
}

export interface Container {
    id: string;
    name: string;
    service: string;
    project: string;
    image: string;
    state: string;
    status: string;
    health: string;
    created_at: string;
    started_at: string;
    ports: PortMapping[];
}

export interface ContainerStats {
    container_id: string;
    cpu_percent: number;
    memory_usage: number;
    memory_limit: number;
    memory_percent: number;
    network_rx: number;
    network_tx: number;
}

export interface LogEntry {
    timestamp: string;
    stream: string;
    message: string;
}

export interface ComposeService {
    name: string;
    container_id: string;
    state: string;
    health: string;
    depends_on: string[];
    ports: PortMapping[];
    image: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface LoginResponse {
    token: string;
    name: string;
    email: string;
}

export interface SystemStats {
    cpu_percent: number;
    cpu_cores: number;
    memory_total: number;
    memory_used: number;
    memory_percent: number;
    disk_total: number;
    disk_used: number;
    disk_percent: number;
}

export interface AuthUser {
    name: string;
    email: string;
}

export type ServiceLayer = 'frontend' | 'backend' | 'infra' | 'monitoring';

export const SERVICE_LAYERS: Record<string, ServiceLayer> = {
    'font-central': 'frontend',
    'front-central': 'frontend',
    'font-website': 'frontend',
    'front-website': 'frontend',
    'front-testing': 'frontend',
    'back-central': 'backend',
    'back-testing': 'backend',
    'redis': 'infra',
    'rabbitmq': 'infra',
    'nginx': 'infra',
    'monitoring-api': 'monitoring',
    'monitoring-web': 'monitoring',
};

export const LAYER_CONFIG: Record<ServiceLayer, { label: string; color: string; bg: string; border: string }> = {
    frontend: { label: 'Frontend', color: '#00f0ff', bg: '#00f0ff08', border: '#00f0ff20' },
    backend: { label: 'Backend', color: '#a855f7', bg: '#a855f708', border: '#a855f720' },
    infra: { label: 'Infrastructure', color: '#ffaa00', bg: '#ffaa0008', border: '#ffaa0020' },
    monitoring: { label: 'Monitoring', color: '#00ff88', bg: '#00ff8808', border: '#00ff8820' },
};
