/** @jsxImportSource react */
import { useEffect, useState } from 'react';

interface TrackingSearchInputProps {
  onSearch: (query: string) => void;
  isLoading?: boolean;
  initialValue?: string;
}

export default function TrackingSearchInput({
  onSearch,
  isLoading = false,
  initialValue = '',
}: TrackingSearchInputProps) {
  const [query, setQuery] = useState(initialValue);

  useEffect(() => {
    if (initialValue && initialValue !== query) {
      setQuery(initialValue);
    }
  }, [initialValue]);

  const handleSubmit = (e: Event) => {
    e.preventDefault();
    if (query.trim()) {
      onSearch(query.trim());
    }
  };

  return (
    <form onSubmit={handleSubmit} class="w-full">
      <div class="relative flex items-center gap-2">
        <div class="flex-1 relative">
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery((e.target as HTMLInputElement).value)}
            placeholder="Ingresa número de tracking o de orden..."
            disabled={isLoading}
            class={`
              w-full px-4 py-3 pr-12 rounded-xl border-2 border-[#EFEAFB]
              focus:border-[#A855F7] focus:outline-none transition-all
              placeholder:text-[#B0A9C2] text-[#181225]
              ${isLoading ? 'bg-[#F6F3FD] cursor-not-allowed' : 'bg-white'}
            `}
          />
          <svg class="absolute right-3 top-1/2 transform -translate-y-1/2 text-[#B0A9C2] w-5 h-5 pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
          </svg>
        </div>
        <button
          type="submit"
          disabled={isLoading || !query.trim()}
          class={`
            px-6 py-3 rounded-xl font-semibold transition-all flex items-center gap-2
            ${
              isLoading || !query.trim()
                ? 'bg-[#F6F3FD] text-[#C7BEDB] cursor-not-allowed'
                : 'bg-gradient-to-r from-[#A855F7] to-[#7C3AED] text-white shadow-[0_10px_20px_-10px_rgba(124,58,237,0.5)] hover:shadow-[0_14px_28px_-10px_rgba(124,58,237,0.55)]'
            }
          `}
        >
          {isLoading ? (
            <>
              <svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
              </svg>
              Buscando...
            </>
          ) : (
            'Rastrear'
          )}
        </button>
      </div>
    </form>
  );
}
