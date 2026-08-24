import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/shipments/domain/entities.dart';
import 'package:mobile_central/shared/utils/rate_pricing.dart';

EnvioClickRate buildRate({
  double flete = 10000,
  double? minimumInsurance = 6530,
  double? extraInsurance = 4000,
  bool? cod = true,
  double? codCarrierFee = 3000,
  double? codProbabilityMargin = 1500,
}) {
  return EnvioClickRate(
    idRate: 1,
    idProduct: 1,
    product: 'Estandar',
    idCarrier: 1,
    carrier: 'Interrapidisimo',
    flete: flete,
    deliveryDays: 2,
    quotationType: 'normal',
    minimumInsurance: minimumInsurance,
    extraInsurance: extraInsurance,
    cod: cod,
    codCarrierFee: codCarrierFee,
    codProbabilityMargin: codProbabilityMargin,
  );
}

void main() {
  group('RatePricing replica la formula de rate-pricing.ts', () {
    test('sin contra entrega ni seguro extra suma flete y seguro minimo', () {
      final rate = buildRate();
      const options = RatePricingOptions();
      expect(RatePricing.guideCost(rate, options), 16530);
      expect(RatePricing.carrierFee(rate, options), 0);
      expect(RatePricing.totalCost(rate, options), 16530);
    });

    test('el seguro minimo nunca se omite', () {
      final rate = buildRate(minimumInsurance: 6530);
      expect(RatePricing.guideCost(rate, const RatePricingOptions()), 16530);
    });

    test('asegurar suma el seguro extra', () {
      final rate = buildRate();
      const options = RatePricingOptions(insured: true);
      expect(RatePricing.guideCost(rate, options), 20530);
    });

    test('contra entrega suma el margen a la guia y la comision al total', () {
      final rate = buildRate();
      const options = RatePricingOptions(cod: true);
      expect(RatePricing.guideCost(rate, options), 18030);
      expect(RatePricing.carrierFee(rate, options), 3000);
      expect(RatePricing.totalCost(rate, options), 21030);
    });

    test('si la tarifa no soporta cod no se aplica margen ni comision', () {
      final rate = buildRate(cod: false);
      const options = RatePricingOptions(cod: true);
      expect(RatePricing.appliesCod(rate, options), isFalse);
      expect(RatePricing.guideCost(rate, options), 16530);
      expect(RatePricing.carrierFee(rate, options), 0);
    });

    test('el desglose cuadra con el total', () {
      final rate = buildRate();
      const options = RatePricingOptions(cod: true, insured: true);
      final breakdown = RatePricing.breakdown(rate, options);
      expect(breakdown.flete, 11500);
      expect(breakdown.minimumInsurance, 6530);
      expect(breakdown.extraInsurance, 4000);
      expect(breakdown.carrierFee, 3000);
      expect(breakdown.guideCost, 22030);
      expect(breakdown.total, 25030);
      expect(breakdown.total, RatePricing.totalCost(rate, options));
    });

    test('valores nulos se tratan como cero', () {
      final rate = buildRate(minimumInsurance: null, extraInsurance: null, codProbabilityMargin: null, codCarrierFee: null);
      const options = RatePricingOptions(cod: true, insured: true);
      expect(RatePricing.guideCost(rate, options), 10000);
      expect(RatePricing.totalCost(rate, options), 10000);
    });
  });
}
