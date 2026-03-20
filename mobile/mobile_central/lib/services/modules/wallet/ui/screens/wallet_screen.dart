import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/wallet_provider.dart';

class WalletScreen extends StatefulWidget {
  final int? businessId;

  const WalletScreen({super.key, this.businessId});

  @override
  State<WalletScreen> createState() => _WalletScreenState();
}

class _WalletScreenState extends State<WalletScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadData();
    });
  }

  void _loadData() {
    final provider = context.read<WalletProvider>();
    provider.fetchBalance();
    provider.fetchHistory();
  }

  Future<void> _showRechargeDialog() async {
    final amountCtrl = TextEditingController();
    final result = await showDialog<double>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Recargar Wallet'),
        content: TextField(
          controller: amountCtrl,
          keyboardType: TextInputType.number,
          decoration: const InputDecoration(
            labelText: 'Monto',
            prefixText: '\$ ',
            border: OutlineInputBorder(),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cancelar'),
          ),
          FilledButton(
            onPressed: () {
              final amount = double.tryParse(amountCtrl.text);
              if (amount != null && amount > 0) {
                Navigator.pop(ctx, amount);
              }
            },
            child: const Text('Recargar'),
          ),
        ],
      ),
    );

    if (result != null && mounted) {
      final success =
          await context.read<WalletProvider>().rechargeWallet(amount: result);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(success
                ? 'Recarga exitosa'
                : context.read<WalletProvider>().error ?? 'Error al recargar'),
          ),
        );
        if (success) _loadData();
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Wallet')),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: _showRechargeDialog,
        icon: const Icon(Icons.add),
        label: const Text('Recargar'),
      ),
      body: Consumer<WalletProvider>(
        builder: (context, provider, _) {
          if (provider.isLoading) {
            return const Center(child: CircularProgressIndicator());
          }
          if (provider.error != null) {
            return Center(
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(Icons.error_outline,
                        size: 48, color: Colors.red),
                    const SizedBox(height: 16),
                    Text(provider.error!,
                        textAlign: TextAlign.center,
                        style: const TextStyle(color: Colors.red)),
                    const SizedBox(height: 16),
                    FilledButton.icon(
                      onPressed: _loadData,
                      icon: const Icon(Icons.refresh),
                      label: const Text('Reintentar'),
                    ),
                  ],
                ),
              ),
            );
          }
          return RefreshIndicator(
            onRefresh: () async => _loadData(),
            child: SingleChildScrollView(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _BalanceCard(
                    balance: provider.wallet?.balance ?? 0,
                  ),
                  const SizedBox(height: 24),
                  const Text('Historial de Transacciones',
                      style:
                          TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                  const SizedBox(height: 12),
                  if (provider.history.isEmpty)
                    Center(
                      child: Padding(
                        padding: const EdgeInsets.symmetric(vertical: 32),
                        child: Column(
                          children: [
                            Icon(Icons.receipt_long_outlined,
                                size: 48, color: Colors.grey.shade400),
                            const SizedBox(height: 8),
                            Text('Sin transacciones',
                                style: TextStyle(
                                    color: Colors.grey.shade600)),
                          ],
                        ),
                      ),
                    )
                  else
                    ...provider.history.map((tx) {
                      final amount = tx is Map
                          ? (tx['Amount'] ?? tx['amount'] ?? 0).toDouble()
                          : 0.0;
                      final date = tx is Map
                          ? (tx['CreatedAt'] ?? tx['created_at'] ?? '')
                          : '';
                      final isPositive = amount >= 0;
                      return Card(
                        margin: const EdgeInsets.only(bottom: 8),
                        child: ListTile(
                          leading: CircleAvatar(
                            backgroundColor: isPositive
                                ? Colors.green.shade100
                                : Colors.red.shade100,
                            child: Icon(
                              isPositive
                                  ? Icons.arrow_downward
                                  : Icons.arrow_upward,
                              color: isPositive
                                  ? Colors.green.shade700
                                  : Colors.red.shade700,
                            ),
                          ),
                          title: Text(
                            '${isPositive ? '+' : ''}\$${amount.toStringAsFixed(0)}',
                            style: TextStyle(
                              fontWeight: FontWeight.w600,
                              color: isPositive ? Colors.green : Colors.red,
                            ),
                          ),
                          subtitle: Text(
                            date.toString().length > 10
                                ? date.toString().substring(0, 10)
                                : date.toString(),
                            style: TextStyle(
                                fontSize: 12, color: Colors.grey.shade600),
                          ),
                        ),
                      );
                    }),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

class _BalanceCard extends StatelessWidget {
  final double balance;
  const _BalanceCard({required this.balance});

  @override
  Widget build(BuildContext context) {
    return Card(
      color: Theme.of(context).colorScheme.primaryContainer,
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          children: [
            Icon(Icons.account_balance_wallet_outlined,
                size: 40,
                color: Theme.of(context).colorScheme.onPrimaryContainer),
            const SizedBox(height: 12),
            Text(
              'Saldo Disponible',
              style: TextStyle(
                fontSize: 14,
                color: Theme.of(context)
                    .colorScheme
                    .onPrimaryContainer
                    .withValues(alpha: 0.7),
              ),
            ),
            const SizedBox(height: 4),
            Text(
              '\$${balance.toStringAsFixed(0)}',
              style: TextStyle(
                fontSize: 36,
                fontWeight: FontWeight.bold,
                color: Theme.of(context).colorScheme.onPrimaryContainer,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
