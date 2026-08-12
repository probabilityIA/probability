import { ISiigoReferralsRepository } from '../domain/ports';
import { GetSiigoReferralsParams } from '../domain/types';

export class SiigoReferralsUseCases {
    constructor(private repository: ISiigoReferralsRepository) { }

    async getReferrals(params?: GetSiigoReferralsParams) {
        return this.repository.getReferrals(params);
    }
}
