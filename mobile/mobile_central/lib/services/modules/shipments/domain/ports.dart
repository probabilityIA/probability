import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IShipmentRepository {
  Future<PaginatedResponse<Shipment>> getShipments(GetShipmentsParams? params);
  Future<Map<String, dynamic>> quoteShipment(EnvioClickQuoteRequest req);
  Future<Map<String, dynamic>> generateGuide(EnvioClickQuoteRequest req);
  Future<Map<String, dynamic>> trackShipment(String trackingNumber);
  Future<Map<String, dynamic>> cancelShipment(String id);
  Future<Shipment> createShipment(Map<String, dynamic> data);
  Future<List<OriginAddress>> getOriginAddresses({int? businessId});
  Future<OriginAddress> createOriginAddress(CreateOriginAddressDTO data, {int? businessId});
  Future<OriginAddress> updateOriginAddress(int id, Map<String, dynamic> data, {int? businessId});
  Future<void> deleteOriginAddress(int id, {int? businessId});
}
