import '../../services/modules/shipments/domain/entities.dart';

class RatePricingOptions {
  const RatePricingOptions({this.cod = false, this.insured = false});

  final bool cod;
  final bool insured;
}

class RateBreakdown {
  const RateBreakdown({
    required this.flete,
    required this.minimumInsurance,
    required this.extraInsurance,
    required this.carrierFee,
    required this.guideCost,
    required this.total,
  });

  final double flete;
  final double minimumInsurance;
  final double extraInsurance;
  final double carrierFee;
  final double guideCost;
  final double total;
}

class RatePricing {
  const RatePricing._();

  static double _num(double? value) {
    if (value == null || value.isNaN || value.isInfinite) return 0;
    return value;
  }

  static bool appliesCod(EnvioClickRate? rate, RatePricingOptions options) =>
      options.cod && rate?.cod == true;

  static double guideCost(EnvioClickRate? rate, RatePricingOptions options) {
    if (rate == null) return 0;
    final insurance = _num(rate.minimumInsurance) +
        (options.insured ? _num(rate.extraInsurance) : 0);
    final double margin =
        appliesCod(rate, options) ? _num(rate.codProbabilityMargin) : 0;
    return _num(rate.flete) + insurance + margin;
  }

  static double carrierFee(EnvioClickRate? rate, RatePricingOptions options) =>
      appliesCod(rate, options) ? _num(rate?.codCarrierFee) : 0;

  static double totalCost(EnvioClickRate? rate, RatePricingOptions options) =>
      guideCost(rate, options) + carrierFee(rate, options);

  static RateBreakdown breakdown(EnvioClickRate? rate, RatePricingOptions options) {
    final double margin =
        appliesCod(rate, options) ? _num(rate?.codProbabilityMargin) : 0;
    final flete = _num(rate?.flete) + margin;
    final minimumInsurance = _num(rate?.minimumInsurance);
    final double extraInsurance =
        options.insured ? _num(rate?.extraInsurance) : 0;
    final fee = carrierFee(rate, options);
    final cost = flete + minimumInsurance + extraInsurance;
    return RateBreakdown(
      flete: flete,
      minimumInsurance: minimumInsurance,
      extraInsurance: extraInsurance,
      carrierFee: fee,
      guideCost: cost,
      total: cost + fee,
    );
  }
}
