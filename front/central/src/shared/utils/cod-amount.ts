export interface CodAmountSource {
    cod_total?: number | null;
    cod_carrier_fee?: number | null;
    cod_includes_shipping?: boolean | null;
    cod_checkout_carrier_fee?: number | null;
}

const num = (value: unknown): number => {
    const n = Number(value ?? 0);
    return Number.isFinite(n) ? n : 0;
};

export const codCheckoutFee = (source: CodAmountSource): number => num(source.cod_checkout_carrier_fee);

export const codEffectiveCarrierFee = (source: CodAmountSource): number => {
    const realFee = num(source.cod_carrier_fee);
    return realFee > 0 ? realFee : codCheckoutFee(source);
};

export const codCustomerCharge = (source: CodAmountSource): number => {
    const total = num(source.cod_total);
    if (total <= 0) return 0;
    const checkoutFee = codCheckoutFee(source);
    if (checkoutFee > 0) return total + checkoutFee;
    return source.cod_includes_shipping ? total : total + num(source.cod_carrier_fee);
};

export const codBusinessNet = (source: CodAmountSource): number => {
    const charge = codCustomerCharge(source);
    if (charge <= 0) return 0;
    return charge - codEffectiveCarrierFee(source);
};
