"use client";

import { useRouter } from "next/navigation";
import { clearAuth, getUser } from "@/shared/lib/auth";

export default function Header() {
  const router = useRouter();
  const user = getUser();

  const handleLogout = () => {
    clearAuth();
    router.push("/login");
  };

  return (
    <header className="h-14 bg-white border-b border-gray-200 flex items-center justify-between px-6">
      <div className="text-sm text-gray-500">
        Super Admin Testing Environment
      </div>
      <div className="flex items-center gap-4">
        {user && (
          <span className="text-sm text-gray-600">{user.email || user.name}</span>
        )}
        <button
          onClick={handleLogout}
          className="text-sm text-red-600 hover:text-red-800 font-medium"
        >
          Logout
        </button>
      </div>
    </header>
  );
}
