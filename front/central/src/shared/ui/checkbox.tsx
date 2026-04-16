"use client";

import * as React from "react";

export interface CheckboxProps extends React.InputHTMLAttributes<HTMLInputElement> {
    onCheckedChange?: (checked: boolean) => void;
}

export function Checkbox({ className = "", onCheckedChange, ...props }: CheckboxProps) {
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (onCheckedChange) {
            onCheckedChange(e.target.checked);
        }
        if (props.onChange) {
            props.onChange(e);
        }
    };

    return (
        <input
            type="checkbox"
            className={`h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary ${className}`}
            onChange={handleChange}
            {...props}
        />
    );
}
