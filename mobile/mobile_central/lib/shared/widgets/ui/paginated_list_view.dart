import 'package:flutter/material.dart';

import '../../pagination/paged_list_controller.dart';
import '../../theme/app_colors.dart';
import '../../theme/app_tokens.dart';
import 'state_views.dart';

class PaginatedListView<T> extends StatefulWidget {
  const PaginatedListView({
    super.key,
    required this.controller,
    required this.itemBuilder,
    required this.emptyTitle,
    required this.emptyMessage,
    required this.emptyIcon,
    required this.unitLabel,
    this.placeholderHeight = 92,
    this.separatorHeight = 10,
    this.padding = const EdgeInsets.fromLTRB(16, 4, 16, 88),
    this.header,
    this.onRefresh,
  });

  final PagedListController<T> controller;
  final Widget Function(BuildContext context, T item, int index) itemBuilder;
  final String emptyTitle;
  final String emptyMessage;
  final IconData emptyIcon;
  final String unitLabel;
  final double placeholderHeight;
  final double separatorHeight;
  final EdgeInsets padding;
  final Widget? header;
  final Future<void> Function()? onRefresh;

  @override
  State<PaginatedListView<T>> createState() => _PaginatedListViewState<T>();
}

class _PaginatedListViewState<T> extends State<PaginatedListView<T>> {
  final _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _scrollController.removeListener(_onScroll);
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (!_scrollController.hasClients) return;
    final position = _scrollController.position;
    if (position.pixels >= position.maxScrollExtent - 320) {
      widget.controller.loadMore();
    }
  }

  Future<void> _refresh() async {
    if (widget.onRefresh != null) {
      await widget.onRefresh!();
      return;
    }
    await widget.controller.refresh();
  }

  @override
  Widget build(BuildContext context) {
    return ListenableBuilder(
      listenable: widget.controller,
      builder: (context, _) {
        final controller = widget.controller;

        if (controller.isLoading && controller.isEmpty) {
          return const AppListSkeleton();
        }

        if (controller.error != null && controller.isEmpty) {
          return AppErrorState(message: controller.error!, onRetry: _refresh);
        }

        if (controller.isEmpty) {
          return AppEmptyState(
            icon: widget.emptyIcon,
            title: widget.emptyTitle,
            message: widget.emptyMessage,
            actionLabel: 'Actualizar',
            onAction: _refresh,
          );
        }

        final headerSlots = widget.header == null ? 0 : 1;

        return RefreshIndicator(
          onRefresh: _refresh,
          color: Theme.of(context).colorScheme.primary,
          child: ListView.separated(
            controller: _scrollController,
            physics: const AlwaysScrollableScrollPhysics(),
            padding: widget.padding,
            cacheExtent: 600,
            addAutomaticKeepAlives: false,
            itemCount: controller.itemCount + headerSlots + 1,
            separatorBuilder: (context, index) =>
                SizedBox(height: widget.separatorHeight),
            itemBuilder: (context, rawIndex) {
              if (headerSlots == 1 && rawIndex == 0) return widget.header!;
              final index = rawIndex - headerSlots;

              if (index == controller.itemCount) {
                return _ListFooter(
                  loading: controller.isLoadingMore,
                  hasMore: controller.hasMore,
                  total: controller.total,
                  shown: controller.itemCount,
                  unitLabel: widget.unitLabel,
                );
              }

              final item = controller.itemAt(index);
              if (item == null) {
                return _SlotPlaceholder(height: widget.placeholderHeight);
              }
              return widget.itemBuilder(context, item, index);
            },
          ),
        );
      },
    );
  }
}

class _SlotPlaceholder extends StatelessWidget {
  const _SlotPlaceholder({required this.height});

  final double height;

  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      decoration: BoxDecoration(
        color: AppColors.surfaceMuted,
        borderRadius: BorderRadius.circular(AppRadius.md),
      ),
    );
  }
}

class _ListFooter extends StatelessWidget {
  const _ListFooter({
    required this.loading,
    required this.hasMore,
    required this.total,
    required this.shown,
    required this.unitLabel,
  });

  final bool loading;
  final bool hasMore;
  final int total;
  final int shown;
  final String unitLabel;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 18),
      child: Center(
        child: loading
            ? const SizedBox(
                height: 22,
                width: 22,
                child: CircularProgressIndicator(strokeWidth: 2.2),
              )
            : Text(
                hasMore
                    ? '$shown de $total $unitLabel'
                    : '$shown $unitLabel',
                style: Theme.of(context).textTheme.labelSmall,
              ),
      ),
    );
  }
}
