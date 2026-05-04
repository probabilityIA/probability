interface IVAIncludedBadgeProps {
  ivaAmount: number;
  currency?: string;
  currencyPresentment?: string;
  amountPresentment?: number;
}

export function IVAIncludedBadge({
  ivaAmount,
  currency = 'USD',
  currencyPresentment = 'COP',
  amountPresentment = 0
}: IVAIncludedBadgeProps) {
  if (ivaAmount <= 0) return null;

  const formatCurrency = (amount: number, curr: string) => {
    return new Intl.NumberFormat('es-CO', {
      style: 'currency',
      currency: curr,
    }).format(amount);
  };

  const displayAmount = amountPresentment > 0 ? amountPresentment : ivaAmount;
  const displayCurrency = amountPresentment > 0 ? currencyPresentment : currency;

  return (
    <div className="flex items-center gap-2 mt-2">
      <div className="inline-flex items-center px-2.5 py-1 rounded-full bg-green-500 text-white text-[8px] font-bold whitespace-nowrap">
        ✓ IVA Included
      </div>
      <span className="text-[8px] text-gray-600 dark:text-gray-300 font-semibold">
        {formatCurrency(displayAmount, displayCurrency)}
      </span>
    </div>
  );
}
