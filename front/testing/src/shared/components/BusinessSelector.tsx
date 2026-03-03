"use client";

import { useState, useEffect } from "react";
import { testingAPI } from "@/shared/lib/api";

interface Business {
  id: number;
  name: string;
}

interface BusinessSelectorProps {
  value: number | null;
  onChange: (id: number | null) => void;
}

export default function BusinessSelector({ value, onChange }: BusinessSelectorProps) {
  const [businesses, setBusinesses] = useState<Business[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    testingAPI<{ data: Business[] }>("/businesses")
      .then((res) => setBusinesses(res.data || []))
      .catch(() => setBusinesses([]))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="bg-indigo-50 border border-indigo-200 rounded-lg p-3 flex items-center gap-3">
      <label className="text-sm font-medium text-indigo-800 whitespace-nowrap">
        Business:
      </label>
      <select
        value={value?.toString() ?? ""}
        onChange={(e) => onChange(e.target.value ? Number(e.target.value) : null)}
        className="flex-1 px-3 py-1.5 bg-white border border-indigo-300 rounded-md text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-indigo-500"
        disabled={loading}
      >
        <option value="">
          {loading ? "Loading..." : "-- Select a business --"}
        </option>
        {businesses.map((b) => (
          <option key={b.id} value={b.id}>
            {b.name} (ID: {b.id})
          </option>
        ))}
      </select>
    </div>
  );
}
