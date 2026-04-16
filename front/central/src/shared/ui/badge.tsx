/**
 * Componente Badge reutilizable
 * Usa clases globales definidas en globals.css
 */

'use client';

import { ReactNode, HTMLAttributes } from 'react';

type BadgeType = 'primary' | 'success' | 'error' | 'warning' | 'secondary';

interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  type?: BadgeType;
  children: ReactNode;
  className?: string;
}

export function Badge({ type = 'primary', children, className = '', ...props }: BadgeProps) {
  const baseClass = type ? `badge badge-${type}` : 'badge';
  
  return (
    <span className={`${baseClass} ${className}`} {...props}>
      {children}
    </span>
  );
}

