import { env } from '@/shared/config/env';
import { ICommercialRepository } from '../../domain/ports';
import { GetProspectsParams, PaginatedProspectsResponse, ProspectStatsResponse, UpdateSeenResponse } from '../../domain/types';

export class CommercialApiRepository implements ICommercialRepository {
  private baseUrl: string;
  private token: string | null;

  constructor(token?: string | null) {
    this.baseUrl = env.API_BASE_URL;
    this.token = token || null;
  }

  private headers(): Record<string, string> {
    const headers: Record<string, string> = {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    };
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }
    return headers;
  }

  async getProspects(params?: GetProspectsParams): Promise<PaginatedProspectsResponse> {
    const searchParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          searchParams.append(key, String(value));
        }
      });
    }

    const res = await fetch(`${this.baseUrl}/commercial/prospects?${searchParams.toString()}`, {
      headers: this.headers(),
      cache: 'no-store',
    });
    const data = await res.json();

    if (!res.ok) {
      throw new Error(data.error || data.message || 'Error al obtener prospectos');
    }

    return data;
  }

  async getProspectsStats(): Promise<ProspectStatsResponse> {
    const res = await fetch(`${this.baseUrl}/commercial/prospects/stats`, {
      headers: this.headers(),
      cache: 'no-store',
    });
    const data = await res.json();

    if (!res.ok) {
      throw new Error(data.error || data.message || 'Error al obtener estadisticas de prospectos');
    }

    return data;
  }

  async updateProspectSeen(id: number, seen: boolean): Promise<UpdateSeenResponse> {
    const res = await fetch(`${this.baseUrl}/commercial/prospects/${id}/seen`, {
      method: 'PATCH',
      headers: this.headers(),
      body: JSON.stringify({ seen }),
    });
    const data = await res.json();

    if (!res.ok) {
      throw new Error(data.error || data.message || 'Error al actualizar visto');
    }

    return data;
  }
}
