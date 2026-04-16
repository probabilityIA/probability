import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class ShipmentApiRepository implements IShipmentRepository {
  final ApiClient _client;
  ShipmentApiRepository(this._client);

  @override
  Future<PaginatedResponse<Shipment>> getShipments(GetShipmentsParams? params) async {
    final response = await _client.get('/shipments', queryParameters: params?.toQueryParams());
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)?.map((e) => Shipment.fromJson(e)).toList() ?? [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<Map<String, dynamic>> quoteShipment(EnvioClickQuoteRequest req) async {
    final response = await _client.post('/shipments/quote', data: req.toJson());
    return response.data;
  }

  @override
  Future<Map<String, dynamic>> generateGuide(EnvioClickQuoteRequest req) async {
    final response = await _client.post('/shipments/generate', data: req.toJson());
    return response.data;
  }

  @override
  Future<Map<String, dynamic>> trackShipment(String trackingNumber) async {
    final response = await _client.get('/shipments/track/$trackingNumber');
    return response.data;
  }

  @override
  Future<Map<String, dynamic>> cancelShipment(String id) async {
    final response = await _client.post('/shipments/$id/cancel');
    return response.data;
  }

  @override
  Future<Shipment> createShipment(Map<String, dynamic> data) async {
    final response = await _client.post('/shipments', data: data);
    return Shipment.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<List<OriginAddress>> getOriginAddresses({int? businessId}) async {
    final qp = businessId != null ? {'business_id': businessId} : null;
    final response = await _client.get('/shipments/origin-addresses', queryParameters: qp);
    return (response.data['data'] as List<dynamic>?)?.map((e) => OriginAddress.fromJson(e)).toList() ?? [];
  }

  @override
  Future<OriginAddress> createOriginAddress(CreateOriginAddressDTO data, {int? businessId}) async {
    final body = data.toJson();
    if (businessId != null) body['business_id'] = businessId;
    final response = await _client.post('/shipments/origin-addresses', data: body);
    return OriginAddress.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<OriginAddress> updateOriginAddress(int id, Map<String, dynamic> data, {int? businessId}) async {
    final response = await _client.put('/shipments/origin-addresses/$id', data: data);
    return OriginAddress.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteOriginAddress(int id, {int? businessId}) async {
    final qp = businessId != null ? {'business_id': businessId} : null;
    await _client.delete('/shipments/origin-addresses/$id', queryParameters: qp);
  }
}
