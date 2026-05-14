export function detectSKUPattern(skus: string[]): { prefix: string; nextNumber: number } | null {
  if (skus.length === 0) return null;

  const patterns = new Map<string, number[]>();

  for (const sku of skus) {
    const match = sku.match(/^([a-z0-9]+-?)(\d+)$/i);
    if (match) {
      const prefix = match[1];
      const num = parseInt(match[2], 10);
      if (!patterns.has(prefix)) patterns.set(prefix, []);
      patterns.get(prefix)!.push(num);
    }
  }

  if (patterns.size === 0) return null;

  const [prefix, numbers] = Array.from(patterns.entries()).reduce((a, b) =>
    b[1].length > a[1].length ? b : a
  );

  const sorted = numbers.sort((a, b) => a - b);

  let nextNum = 1;
  for (const num of sorted) {
    if (num === nextNum) {
      nextNum++;
    } else if (num > nextNum) {
      break;
    }
  }

  return {
    prefix,
    nextNumber: nextNum,
  };
}

export function formatSKUWithNumber(prefix: string, number: number, padLength: number = 3): string {
  return `${prefix}${String(number).padStart(padLength, '0')}`;
}

export function extractSKUPrefix(sku: string): string {
  const match = sku.match(/^([a-z0-9]+-?)(?:\d+)?$/i);
  return match ? match[1] : sku;
}
