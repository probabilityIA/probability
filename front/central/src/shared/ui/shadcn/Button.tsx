import React from 'react';
import clsx from 'clsx';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'default' | 'ghost' | 'outline';
}

const Button: React.FC<ButtonProps> = ({ children, className, variant = 'default', ...props }) => {
  return (
    <button
      {...props}
      className={clsx(
        'inline-flex items-center justify-center rounded-md px-4 py-2 text-sm font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2',
        variant === 'default' && 'bg-slate-900 text-white hover:bg-slate-800',
        variant === 'ghost' && 'bg-transparent hover:bg-slate-50 text-slate-900',
        variant === 'outline' && 'border border-slate-200 bg-white text-slate-900',
        className
      )}
    >
      {children}
    </button>
  );
};

export default Button;
