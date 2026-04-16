'use client';

import { useEffect, useState } from 'react';

export const ThemeToggle = () => {
  const [isDark, setIsDark] = useState(false);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    // Detectar tema actual del documento
    const htmlElement = document.documentElement;
    const isDarkMode = htmlElement.classList.contains('dark');
    setIsDark(isDarkMode);
  }, []);

  const toggleTheme = () => {
    const htmlElement = document.documentElement;

    if (isDark) {
      htmlElement.classList.remove('dark');
      localStorage.setItem('theme', 'light');
      setIsDark(false);
    } else {
      htmlElement.classList.add('dark');
      localStorage.setItem('theme', 'dark');
      setIsDark(true);
    }
  };

  if (!mounted) {
    return null; // Evitar hidratación mismatch
  }

  return (
    <button
      onClick={toggleTheme}
      className={isDark ? 'theme-toggle-dark' : 'theme-toggle-light'}
      aria-label={isDark ? 'Cambiar a modo claro' : 'Cambiar a modo oscuro'}
      title={isDark ? 'Modo claro' : 'Modo oscuro'}
    >
      <span className="theme-toggle-emoji">
        {isDark ? '🌙' : '☀️'}
      </span>
    </button>
  );
};
