import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../core/errors/error_parser.dart';
import '../types/paginated_response.dart';
import 'paged_collection.dart';

typedef PageFetcher<T> = Future<PaginatedResponse<T>> Function(
  int page,
  int pageSize,
);

class PagedListController<T> extends ChangeNotifier {
  PagedListController({
    required PageFetcher<T> fetcher,
    int pageSize = 20,
    int maxPagesInMemory = 8,
  })  : _fetcher = fetcher,
        _collection = PagedCollection<T>(
          pageSize: pageSize,
          maxPagesInMemory: maxPagesInMemory,
        );

  PageFetcher<T> _fetcher;
  final PagedCollection<T> _collection;

  final Set<int> _inFlight = <int>{};
  final Set<int> _framePages = <int>{};

  bool _isLoading = false;
  bool _isLoadingMore = false;
  bool _maintenanceScheduled = false;
  bool _disposed = false;
  String? _error;

  int get pageSize => _collection.pageSize;
  int get maxPagesInMemory => _collection.maxPagesInMemory;
  int get itemCount => _collection.length;
  int get total => _collection.total;
  int get liveItemCount => _collection.liveItemCount;
  int get pagesInMemory => _collection.pagesInMemory;
  bool get isLoading => _isLoading;
  bool get isLoadingMore => _isLoadingMore;
  bool get hasMore => _collection.hasMore;
  bool get isEmpty => _collection.isEmpty;
  String? get error => _error;
  List<T> get loadedItems => _collection.loadedItems;

  set fetcher(PageFetcher<T> value) => _fetcher = value;

  bool isHole(int index) => _collection.isHole(index);

  T? itemAt(int index) {
    final item = _collection.at(index);
    _framePages.add(_collection.pageOf(index));
    _scheduleMaintenance();
    return item;
  }

  Future<void> refresh() async {
    if (_isLoading) return;
    _isLoading = true;
    _error = null;
    _inFlight.clear();
    notifyListeners();

    try {
      final response = await _fetcher(1, pageSize);
      _collection.reset();
      _collection.setPage(1, response.data, response.pagination.total);
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    _safeNotify();
  }

  Future<void> loadMore() async {
    if (_isLoading || _isLoadingMore || !hasMore) return;
    final next = _collection.highestPageLoaded + 1;
    if (_inFlight.contains(next)) return;

    _isLoadingMore = true;
    _inFlight.add(next);
    notifyListeners();

    try {
      final response = await _fetcher(next, pageSize);
      _collection.setPage(next, response.data, response.pagination.total);
      _error = null;
    } catch (e) {
      _error = parseError(e);
    }

    _inFlight.remove(next);
    _isLoadingMore = false;
    _scheduleMaintenance();
    _safeNotify();
  }

  void clear() {
    _collection.reset();
    _inFlight.clear();
    _framePages.clear();
    _error = null;
    _safeNotify();
  }

  void _scheduleMaintenance() {
    if (_maintenanceScheduled || _disposed) return;
    _maintenanceScheduled = true;
    scheduleMicrotask(_runMaintenance);
  }

  void _runMaintenance() {
    _maintenanceScheduled = false;
    if (_disposed) return;

    final visible = Set<int>.from(_framePages);
    _framePages.clear();

    _collection.evictColdPages(protectedPages: visible);

    for (final page in visible) {
      if (_collection.evictedPages.contains(page)) {
        unawaited(_refill(page));
      }
    }
  }

  Future<void> _refill(int page) async {
    if (_inFlight.contains(page) || _disposed) return;
    _inFlight.add(page);

    try {
      final response = await _fetcher(page, pageSize);
      if (_disposed) return;
      _collection.setPage(page, response.data, response.pagination.total);
      _scheduleMaintenance();
      _safeNotify();
    } catch (_) {
    } finally {
      _inFlight.remove(page);
    }
  }

  void _safeNotify() {
    if (_disposed) return;
    notifyListeners();
  }

  @override
  void dispose() {
    _disposed = true;
    _collection.reset();
    _inFlight.clear();
    _framePages.clear();
    super.dispose();
  }
}
