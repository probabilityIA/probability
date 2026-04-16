import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/actions/domain/entities.dart';

void main() {
  group('ActionEntity', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'name': 'read',
        'description': 'Read access',
      };

      final action = ActionEntity.fromJson(json);

      expect(action.id, 1);
      expect(action.name, 'read');
      expect(action.description, 'Read access');
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final action = ActionEntity.fromJson(json);

      expect(action.id, 0);
      expect(action.name, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 5,
        'name': 'write',
      };

      final action = ActionEntity.fromJson(json);

      expect(action.description, isNull);
    });
  });

  group('GetActionsParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetActionsParams(
        page: 2,
        pageSize: 15,
        name: 'read',
      );

      final queryParams = params.toQueryParams();

      expect(queryParams['page'], 2);
      expect(queryParams['page_size'], 15);
      expect(queryParams['name'], 'read');
    });

    test('toQueryParams excludes null fields', () {
      final params = GetActionsParams(page: 1);

      final queryParams = params.toQueryParams();

      expect(queryParams.length, 1);
      expect(queryParams['page'], 1);
      expect(queryParams.containsKey('page_size'), false);
      expect(queryParams.containsKey('name'), false);
    });

    test('toQueryParams excludes empty name', () {
      final params = GetActionsParams(name: '');

      final queryParams = params.toQueryParams();

      expect(queryParams.containsKey('name'), false);
    });

    test('toQueryParams returns empty map when all null', () {
      final params = GetActionsParams();

      final queryParams = params.toQueryParams();

      expect(queryParams, isEmpty);
    });
  });

  group('CreateActionDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = CreateActionDTO(
        name: 'delete',
        description: 'Delete permission',
      );

      final json = dto.toJson();

      expect(json['name'], 'delete');
      expect(json['description'], 'Delete permission');
    });

    test('toJson includes only name when description is null', () {
      final dto = CreateActionDTO(name: 'read');

      final json = dto.toJson();

      expect(json.length, 1);
      expect(json['name'], 'read');
      expect(json.containsKey('description'), false);
    });
  });

  group('UpdateActionDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = UpdateActionDTO(
        name: 'updated',
        description: 'Updated description',
      );

      final json = dto.toJson();

      expect(json['name'], 'updated');
      expect(json['description'], 'Updated description');
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateActionDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateActionDTO(name: 'newName');

      final json = dto.toJson();

      expect(json.length, 1);
      expect(json['name'], 'newName');
    });
  });
}
