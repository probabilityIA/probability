export interface CodAmountSource {
    cod_total?: number | null;
    cod_carrier_fee?: number | null;
    cod_includes_shipping?: boolean | null;
}

const num = (value: unknown): number => {
    const n = Number(value ?? 0);
    return Number.isFinite(n) ? n : 0;
};

export const codCustomerCharge = (source: CodAmountSource): number => {
    const total = num(source.cod_total);
    if (total <= 0) return 0;
    return source.cod_includes_shipping ? total : total + num(source.cod_carrier_fee);
};

export const codBusinessNet = (source: CodAmountSource): number => {
    const total = num(source.cod_total);
    if (total <= 0) return 0;
    return source.cod_includes_shipping ? total - num(source.cod_carrier_fee) : total;
};
