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

