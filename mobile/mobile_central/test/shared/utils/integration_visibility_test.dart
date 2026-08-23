import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/shared/utils/integration_visibility.dart';

void main() {
  group('visibilidad de integraciones', () {
    test('los modulos internos no se muestran', () {
      expect(
        IntegrationVisibility.isVisible(category: 'internal', name: 'Inventario'),
        isFalse,
      );
      expect(IntegrationVisibility.isHiddenCategory('INTERNAL'), isTrue);
      expect(IntegrationVisibility.isHiddenCategory(' internal '), isTrue);
    });

    test('envioclick no se muestra: para el cliente no existe', () {
      expect(
        IntegrationVisibility.isVisible(category: 'shipping', type: 'envioclick'),
        isFalse,
      );
      expect(
        IntegrationVisibility.isVisible(category: 'shipping', name: 'EnvioClick'),
        isFalse,
      );
      expect(IntegrationVisibility.isHiddenType('Envio Click'), isTrue);
    });

    test('lo que si ve el cliente pasa', () {
      expect(
        IntegrationVisibility.isVisible(category: 'ecommerce', type: 'shopify'),
        isTrue,
      );
      expect(
        IntegrationVisibility.isVisible(category: 'invoicing', type: 'siigo'),
        isTrue,
      );
      expect(
        IntegrationVisibility.isVisible(category: 'messaging', type: 'whatsapp'),
        isTrue,
      );
    });

    test('la plataforma tampoco: es Probability, no un canal del cliente', () {
      expect(
        IntegrationVisibility.isVisible(category: 'platform', name: 'Plataforma'),
        isFalse,
      );
    });

    test('sin datos no se oculta nada', () {
      expect(IntegrationVisibility.isVisible(), isTrue);
      expect(IntegrationVisibility.isHiddenCategory(null), isFalse);
      expect(IntegrationVisibility.isHiddenType(null), isFalse);
    });
  });
}
