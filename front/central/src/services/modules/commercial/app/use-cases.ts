import { ICommercialRepository } from '../domain/ports';
import { GetProspectsParams } from '../domain/types';

export class CommercialUseCases {
  constructor(private repository: ICommercialRepository) { }

  async getProspects(params?: GetProspectsParams) {
    return this.repository.getProspects(params);
  }

  async getProspectsStats() {
    return this.repository.getProspectsStats();
  }
}
