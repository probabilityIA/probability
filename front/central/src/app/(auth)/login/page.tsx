'use client';

import { LoginForm } from '@/services/auth/login/ui';
import { useSearchParams } from 'next/navigation';
import { useEffect, Suspense, useState } from 'react';
import { CookieStorage } from '@/shared/utils';
import { useShopifyAuth } from '@/providers/ShopifyAuthProvider';
import { useRouter } from 'next/navigation';
import { ThemeToggle } from '@/shared/ui/theme-toggle';
import { LoginHeroPanel } from '@/shared/ui/login-hero-panel';

function LoginContent() {
  const searchParams = useSearchParams();
  const { isShopifyEmbedded, sessionToken, isLoading: isShopifyLoading } = useShopifyAuth();
  const router = useRouter();
  const [isAuthenticating, setIsAuthenticating] = useState(false);
  const [isDark, setIsDark] = useState(false);

  useEffect(() => {
    const error = searchParams.get('error');
    if (error === 'no_business') {
      console.warn('Usuario no tiene negocio asignado. Contacte al administrador.');
    }
  }, [searchParams]);

  useEffect(() => {
    const htmlElement = document.documentElement;
    setIsDark(htmlElement.classList.contains('dark'));

    const observer = new MutationObserver(() => {
      setIsDark(htmlElement.classList.contains('dark'));
    });

    observer.observe(htmlElement, { attributes: true });
    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    const authenticateWithShopify = async () => {
      if (isShopifyEmbedded && sessionToken) {
        setIsAuthenticating(true);
        try {
          const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'https://www.probabilityia.com.co/api/v1';
          const response = await fetch(`${baseUrl}/integrations/shopify/auth/login`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ session_token: sessionToken }),
            credentials: 'include',
          });

          if (response.ok) {
            const data = await response.json();
            if (data.user) {
              CookieStorage.setUser(data.user);
            }

            console.log('✅ Login con Shopify exitoso, redirigiendo...');
            router.push('/home');
          } else {
            console.error('Fallo login con Shopify', response.status);
            setIsAuthenticating(false);
          }
        } catch (error) {
          console.error('Error autenticando con Shopify', error);
          setIsAuthenticating(false);
        }
      }
    };

    authenticateWithShopify();
  }, [isShopifyEmbedded, sessionToken, router]);

  if (isShopifyEmbedded && (isShopifyLoading || isAuthenticating)) {
    return (
      <div className="h-screen w-screen flex flex-col items-center justify-center bg-white dark:bg-navy-900">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mb-4"></div>
        <p className="text-gray-600 dark:text-gray-300">Autenticando con Shopify...</p>
      </div>
    );
  }

  return (
    <div className={`h-screen w-screen flex overflow-hidden ${isDark ? 'bg-[#0f0a1e]' : 'bg-white'}`}>
      {/* Theme Toggle */}
      <div className="fixed top-5 left-5 z-50">
        <ThemeToggle />
      </div>

      {/* Left Side - Login Form */}
      <div className={`w-full lg:w-[40%] h-full flex flex-col items-center justify-center px-6 sm:px-12 md:px-24 py-8 overflow-y-auto ${isDark ? 'bg-[#0f0a1e]' : 'bg-white'}`}>
        <div className="w-full max-w-sm">
          <LoginForm />
        </div>
      </div>

      {/* Right Side - Dark Mode SVG */}
      {isDark && (
        <div className="hidden lg:flex lg:w-[60%] h-full relative overflow-hidden items-center justify-center bg-gradient-to-br from-[#1e1145] via-[#2d1a6e] to-[#14093a]">
          {/* Decorative Circles/Glows */}
          <div className="absolute rounded-full blur-3xl opacity-60 bg-[#8b5cf6]/20" style={{ width: '400px', height: '400px', top: '80px', right: '150px' }}></div>
          <div className="absolute rounded-full blur-3xl opacity-50 bg-[#14f5a0]/15" style={{ width: '300px', height: '300px', bottom: '250px', right: '350px' }}></div>
          <div className="absolute rounded-full blur-3xl opacity-40 bg-[#f472b6]/15" style={{ width: '280px', height: '280px', bottom: '300px', left: '50px' }}></div>

          {/* SVG Canvas - DARK MODE ONLY */}
          <svg className="w-full h-full" viewBox="0 0 720 900" fill="none" preserveAspectRatio="xMidYMid slice">
            {/* Background Gradient */}
            <defs>
              <linearGradient id="darkGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stopColor="#1e1145" />
                <stop offset="50%" stopColor="#2d1a6e" />
                <stop offset="100%" stopColor="#14093a" />
              </linearGradient>
            </defs>
            <rect width="720" height="900" fill="url(#darkGrad)" />

            {/* Glow Circles */}
            <circle cx="520" cy="290" r="200" fill="#8b5cf6" opacity="0.15" filter="blur(40px)" />
            <circle cx="550" cy="500" r="150" fill="#14f5a0" opacity="0.1" filter="blur(35px)" />
            <circle cx="100" cy="500" r="140" fill="#f472b6" opacity="0.1" filter="blur(30px)" />

            {/* Delivery Truck */}
            <rect x="140" y="320" width="180" height="100" rx="12" fill="#8b5cf6" opacity="0.95" />
            <rect x="300" y="320" width="80" height="100" rx="12" fill="#7c3aed" />
            <rect x="320" y="335" width="45" height="35" rx="4" fill="#c4b5fd" opacity="0.4" />

            {/* Wheels */}
            <circle cx="180" cy="420" r="20" fill="#1a1430" stroke="#4c1d95" strokeWidth="5" />
            <circle cx="180" cy="420" r="9" fill="#8b5cf6" />
            <circle cx="280" cy="420" r="20" fill="#1a1430" stroke="#4c1d95" strokeWidth="5" />
            <circle cx="280" cy="420" r="9" fill="#8b5cf6" />
            <circle cx="350" cy="420" r="20" fill="#1a1430" stroke="#4c1d95" strokeWidth="5" />
            <circle cx="350" cy="420" r="9" fill="#8b5cf6" />


            {/* Packages */}
            <rect x="70" y="260" width="70" height="80" rx="10" fill="#14f5a0" opacity="0.85" />
            <text x="105" y="310" fontSize="32" fontWeight="700" fill="white" textAnchor="middle" opacity="0.25">📦</text>

            <rect x="450" y="220" width="65" height="70" rx="8" fill="#f472b6" opacity="0.8" />
            <text x="482" y="265" fontSize="28" fontWeight="700" fill="white" textAnchor="middle" opacity="0.25">📮</text>

            <rect x="520" y="360" width="60" height="60" rx="8" fill="#c4b5fd" opacity="0.75" />
            <text x="550" y="397" fontSize="28" fontWeight="700" fill="white" textAnchor="middle" opacity="0.25">📦</text>

            {/* Route Line */}
            <path d="M 140 300 Q 240 280 360 360" stroke="#14f5a0" strokeWidth="3" fill="none" opacity="0.4" strokeLinecap="round" />

            {/* Location Pins */}
            <circle cx="140" cy="300" r="14" fill="#14f5a0" />
            <circle cx="140" cy="300" r="8" fill="#0f0a1e" />
            <circle cx="360" cy="360" r="14" fill="#f472b6" />
            <circle cx="360" cy="360" r="8" fill="#0f0a1e" />

            {/* Tracking Card - Moved up */}
            <rect x="70" y="200" width="240" height="80" rx="14" fill="white" opacity="0.12" stroke="white" strokeWidth="1" strokeOpacity="0.15" />
            <circle cx="98" cy="232" r="16" fill="#14f5a0" opacity="0.3" />
            <text x="120" y="225" fontSize="11" fontWeight="600" fill="white" opacity="0.8">En camino</text>
            <text x="120" y="242" fontSize="13" fontWeight="700" fill="white">Bogotá → Medellín</text>

            {/* Sparkles */}
            <circle cx="480" cy="210" r="2" fill="#14f5a0" opacity="0.6" />
            <circle cx="520" cy="240" r="1.5" fill="#f472b6" opacity="0.5" />
            <circle cx="300" cy="220" r="1.5" fill="#c4b5fd" opacity="0.5" />
            <circle cx="650" cy="390" r="1.5" fill="#14f5a0" opacity="0.4" />

            {/* Floating Dots */}
            <circle cx="350" cy="500" r="5" fill="#8b5cf6" opacity="0.5" />
            <circle cx="650" cy="340" r="4" fill="#14f5a0" opacity="0.4" />
            <circle cx="100" cy="440" r="3" fill="#f472b6" opacity="0.5" />

            {/* Hero Text */}
            <text x="50" y="560" fontSize="22" fontWeight="800" fill="white">Haz tu mejor</text>
            <text x="50" y="590" fontSize="22" fontWeight="800" fill="white">movimiento logístico</text>

            {/* Ecommerce Badges */}
            <g>
              <rect x="50" y="640" width="100" height="20" rx="10" fill="white" opacity="0.1" stroke="white" strokeWidth="0.5" strokeOpacity="0.1" />
              <text x="60" y="653" fontSize="10" fontWeight="600" fill="white" opacity="0.9">🛍️ Shopify</text>

              <rect x="160" y="640" width="100" height="20" rx="10" fill="white" opacity="0.1" stroke="white" strokeWidth="0.5" strokeOpacity="0.1" />
              <text x="170" y="653" fontSize="10" fontWeight="600" fill="white" opacity="0.9">💬 WhatsApp</text>

              <rect x="270" y="640" width="110" height="20" rx="10" fill="white" opacity="0.1" stroke="white" strokeWidth="0.5" strokeOpacity="0.1" />
              <text x="280" y="653" fontSize="10" fontWeight="600" fill="white" opacity="0.9">🛒 Woocommerce</text>

              <rect x="390" y="640" width="120" height="20" rx="10" fill="white" opacity="0.1" stroke="white" strokeWidth="0.5" strokeOpacity="0.1" />
              <text x="400" y="653" fontSize="10" fontWeight="600" fill="white" opacity="0.9">📦 Mercado Libre</text>

              <rect x="50" y="670" width="90" height="20" rx="10" fill="white" opacity="0.1" stroke="white" strokeWidth="0.5" strokeOpacity="0.1" />
              <text x="62" y="683" fontSize="10" fontWeight="600" fill="white" opacity="0.9">🏬 Falabella</text>
            </g>

            {/* Connection Lines */}
            <path d="M 200 460 Q 250 440 300 470" stroke="white" strokeWidth="2" fill="none" opacity="0.1" strokeDasharray="5,5" />
          </svg>
        </div>
      )}

      {/* Right Side - Light Mode LoginHeroPanel */}
      {!isDark && <LoginHeroPanel />}
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center bg-white">
        <div className="text-gray-500 dark:text-gray-400">Cargando...</div>
      </div>
    }>
      <LoginContent />
    </Suspense>
  );
}
