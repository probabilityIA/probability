import React from 'react';

// Simple className utility
const cn = (...classes: (string | boolean | undefined)[]) => classes.filter(Boolean).join(' ');

interface Step {
    id: number;
    label: string;
}

interface StepperProps {
    steps: Step[];
    currentStep: number;
    className?: string;
}

export function Stepper({ steps, currentStep, className }: StepperProps) {
    return (
        <div className={cn("w-full py-4", className)}>
            <div className="flex items-center justify-between">
                {steps.map((step, index) => {
                    const isActive = currentStep === step.id;
                    const isCompleted = currentStep > step.id;
                    const isLast = index === steps.length - 1;

                    return (
                        <React.Fragment key={step.id}>
                            <div className="flex flex-col items-center flex-1">
                                {/* Circle */}
                                <div
                                    className={cn(
                                        "w-10 h-10 rounded-full flex items-center justify-center font-semibold transition-colors",
                                        isCompleted && "bg-purple-700 text-white",
                                        isActive && "bg-purple-700 text-white",
                                        !isActive && !isCompleted && "bg-gray-200 text-gray-500"
                                    )}
                                >
                                    {isCompleted ? (
                                        <svg
                                            className="w-6 h-6"
                                            fill="none"
                                            stroke="currentColor"
                                            viewBox="0 0 24 24"
                                        >
                                            <path
                                                strokeLinecap="round"
                                                strokeLinejoin="round"
                                                strokeWidth={2}
                                                d="M5 13l4 4L19 7"
                                            />
                                        </svg>
                                    ) : (
                                        step.id
                                    )}
                                </div>
                                {/* Label */}
                                <span
                                    className={cn(
                                        "mt-2 text-sm font-medium text-center",
                                        isActive && "text-purple-700",
                                        isCompleted && "text-purple-700",
                                        !isActive && !isCompleted && "text-gray-500"
                                    )}
                                >
                                    {step.label}
                                </span>
                            </div>

                            {/* Connector Line */}
                            {!isLast && (
                                <div
                                    className={cn(
                                        "h-1 flex-1 mx-2 rounded transition-colors",
                                        isCompleted ? "bg-purple-700" : "bg-gray-200"
                                    )}
                                    style={{ marginTop: '-2rem' }}
                                />
                            )}
                        </React.Fragment>
                    );
                })}
            </div>
        </div>
    );
}
