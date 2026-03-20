import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/fulfillmentstatus/domain/entities.dart';

void main() {
  group('FulfillmentStatusInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'code': 'delivered',
        'name': 'Delivered',
        'description': 'Package has been delivered',
        'category': 'completed',
        'color': '#00FF00',
      };

      final status = FulfillmentStatusInfo.fromJson(json);

      expect(status.id, 1);
      expect(status.code, 'delivered');
      expect(status.name, 'Delivered');
      expect(status.description, 'Package has been delivered');
      expect(status.category, 'completed');
      expect(status.color, '#00FF00');
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final status = FulfillmentStatusInfo.fromJson(json);

      expect(status.id, 0);
      expect(status.code, '');
      expect(status.name, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 5,
        'code': 'pending',
        'name': 'Pending',
      };

      final status = FulfillmentStatusInfo.fromJson(json);

      expect(status.description, isNull);
      expect(status.category, isNull);
      expect(status.color, isNull);
    });

    test('fromJson handles explicitly null optional fields', () {
      final json = {
        'id': 3,
        'code': 'processing',
        'name': 'Processing',
        'description': null,
        'category': null,
        'color': null,
      };

      final status = FulfillmentStatusInfo.fromJson(json);

      expect(status.id, 3);
      expect(status.code, 'processing');
      expect(status.name, 'Processing');
      expect(status.description, isNull);
      expect(status.category, isNull);
      expect(status.color, isNull);
    });

    test('fromJson parses all required fields with optional present', () {
      final json = {
        'id': 10,
        'code': 'shipped',
        'name': 'Shipped',
        'description': 'In transit',
        'category': 'active',
        'color': '#FFA500',
      };

      final status = FulfillmentStatusInfo.fromJson(json);

      expect(status.id, 10);
      expect(status.code, 'shipped');
      expect(status.name, 'Shipped');
      expect(status.description, 'In transit');
      expect(status.category, 'active');
      expect(status.color, '#FFA500');
    });
  });
}
