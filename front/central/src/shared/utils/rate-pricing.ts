export interface PriceableRate {
    flete?: number;
    minimumInsurance?: number;
    extraInsurance?: number;
    cod?: boolean;
    codCarrierFee?: number;
    codProbabilityMargin?: number;
}

export interface RatePricingOptions {
    cod?: boolean;
    insured?: boolean;
}

const num = (value: unknown): number => {
    const n = Number(value ?? 0);
    return Number.isFinite(n) ? n : 0;
};

export const rateAppliesCOD = (rate: PriceableRate | null | undefined, opts: RatePricingOptions = {}): boolean =>
    Boolean(opts.cod) && rate?.cod === true;

export const rateGuideCost = (rate: PriceableRate | null | undefined, opts: RatePricingOptions = {}): number => {
    if (!rate) return 0;
    const insurance = num(rate.minimumInsurance) + (opts.insured ? num(rate.extraInsurance) : 0);
    const margin = rateAppliesCOD(rate, opts) ? num(rate.codProbabilityMargin) : 0;
    return num(rate.flete) + insurance + margin;
};

export const rateCarrierFee = (rate: PriceableRate | null | undefined, opts: RatePricingOptions = {}): number =>
    rateAppliesCOD(rate, opts) ? num(rate?.codCarrierFee) : 0;

export const rateTotalCost = (rate: PriceableRate | null | undefined, opts: RatePricingOptions = {}): number =>
    rateGuideCost(rate, opts) + rateCarrierFee(rate, opts);
