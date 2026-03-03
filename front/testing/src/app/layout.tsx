import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Probability - Testing Platform",
  description: "Testing platform for Probability",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="es">
      <body className="min-h-screen bg-gray-50">{children}</body>
    </html>
  );
}
