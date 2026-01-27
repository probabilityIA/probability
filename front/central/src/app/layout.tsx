import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import Script from "next/script";
import { FooterWrapper } from "./footer-wrapper";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Probability App",
  description: "Menos devoluciones en tu e-commerce gracias a la probabilidad",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const apiKey = process.env.NEXT_PUBLIC_SHOPIFY_API_KEY;

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
        className={`${geistSans.variable} ${geistMono.variable} antialiased min-h-screen flex flex-col`}
      >
        <div className="flex-1">
          {children}
        </div>
        <FooterWrapper />
      </body>
    </html>
  );
}

