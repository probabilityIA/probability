import type { Metadata } from "next";

export const dynamic = "force-dynamic";
import { Inter, Roboto_Mono } from "next/font/google";
import Script from "next/script";
import { FooterWrapper } from "./footer-wrapper";
import { ClientProviders } from "@/providers/ClientProviders";
import "./globals.css";

const inter = Inter({
  variable: "--font-inter",
  subsets: ["latin"],
});

const robotoMono = Roboto_Mono({
  variable: "--font-roboto-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Probability App",
  description: "Menos devoluciones en tu e-commerce gracias a la probabilidad",
};

async function getShopifyConfig() {
  try {
    const baseUrl = process.env.API_BASE_URL || 'http://localhost:3050/api/v1';
    const res = await fetch(`${baseUrl}/integrations/shopify/config`, {
      cache: 'no-store',
      next: { tags: ['shopify-config'] }
    });

    if (!res.ok) {
      console.warn(`Error fetching Shopify config: ${res.status}`);
      return null;
    }

    const data = await res.json();
    return data.shopify_client_id;
  } catch (error) {
    console.error("Failed to fetch Shopify config:", error);
    return null;
  }
}

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const apiKey = await getShopifyConfig();

  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        {apiKey && (
          <Script
            src="https://cdn.shopify.com/shopify-cloud/app-bridge.js"
            strategy="beforeInteractive"
            data-api-key={apiKey}
          />
        )}
      </head>
      <body
        className={`${inter.variable} ${robotoMono.variable} antialiased min-h-screen flex flex-col`}
        suppressHydrationWarning
      >
        <div className="flex-1">
          <ClientProviders>
            {children}
          </ClientProviders>
        </div>
        <FooterWrapper />
      </body>
    </html>
  );
}

