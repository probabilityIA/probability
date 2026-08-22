import 'dart:collection';

class PagedCollection<T> {
  PagedCollection({
    this.pageSize = 20,
    this.maxPagesInMemory = 8,
  })  : assert(pageSize > 0),
        assert(maxPagesInMemory >= 2);

  final int pageSize;
  final int maxPagesInMemory;

  final List<T?> _slots = <T?>[];
  final Map<int, int> _pageUse = <int, int>{};
  final Set<int> _loadedPages = <int>{};
  final Set<int> _evictedPages = <int>{};

  int _total = 0;
  int _highestPageLoaded = 0;
  int _clock = 0;

  int get total => _total;
  int get length => _slots.length;
  bool get isEmpty => _slots.isEmpty;
  bool get isNotEmpty => _slots.isNotEmpty;
  int get highestPageLoaded => _highestPageLoaded;
  int get pagesInMemory => _loadedPages.length;

  int get liveItemCount {
    var count = 0;
    for (final slot in _slots) {
      if (slot != null) count++;
    }
    return count;
  }

  UnmodifiableSetView<int> get evictedPages => UnmodifiableSetView(_evictedPages);

  bool get hasMore => _slots.length < _total;

  int pageOf(int index) => (index ~/ pageSize) + 1;

  T? at(int index) {
    if (index < 0 || index >= _slots.length) return null;
    _touch(pageOf(index));
    return _slots[index];
  }

  T? peek(int index) {
    if (index < 0 || index >= _slots.length) return null;
    return _slots[index];
  }

  bool isHole(int index) => peek(index) == null && index < _slots.length;

  List<T> get loadedItems {
    final out = <T>[];
    for (final slot in _slots) {
      if (slot != null) out.add(slot);
    }
    return out;
  }

  void reset() {
    _slots.clear();
    _pageUse.clear();
    _loadedPages.clear();
    _evictedPages.clear();
    _total = 0;
    _highestPageLoaded = 0;
    _clock = 0;
  }

  void setPage(int page, List<T> items, int total) {
    final start = (page - 1) * pageSize;
    _total = total > 0 ? total : start + items.length;
    final needed = start + items.length;

    if (_slots.length < needed) {
      _slots.length = needed;
    }
    for (var i = 0; i < items.length; i++) {
      _slots[start + i] = items[i];
    }

    if (_total > 0 && _slots.length > _total) {
      _slots.length = _total;
    }

    _loadedPages.add(page);
    _evictedPages.remove(page);
    if (page > _highestPageLoaded) _highestPageLoaded = page;
    _touch(page);
  }

  void replaceAt(int index, T item) {
    if (index < 0 || index >= _slots.length) return;
    _slots[index] = item;
  }

  void _touch(int page) {
    _clock++;
    _pageUse[page] = _clock;
  }

  List<int> evictColdPages({Set<int> protectedPages = const <int>{}}) {
    final evicted = <int>[];
    if (_loadedPages.length <= maxPagesInMemory) return evicted;

    final candidates = _loadedPages
        .where((page) => !protectedPages.contains(page))
        .toList()
      ..sort((a, b) => (_pageUse[a] ?? 0).compareTo(_pageUse[b] ?? 0));

    var over = _loadedPages.length - maxPagesInMemory;
    for (final page in candidates) {
      if (over <= 0) break;
      _clearPage(page);
      evicted.add(page);
      over--;
    }
    return evicted;
  }

  void _clearPage(int page) {
    final start = (page - 1) * pageSize;
    final end = start + pageSize;
    for (var i = start; i < end && i < _slots.length; i++) {
      _slots[i] = null;
    }
    _loadedPages.remove(page);
    _evictedPages.add(page);
    _pageUse.remove(page);
  }
}
