'use client';

import { useState } from 'react';
import { LinkIcon, KeyIcon } from '@heroicons/react/24/outline';
import { JumpsellerOAuthForm } from './JumpsellerOAuthForm';
import { JumpsellerConfigForm } from './JumpsellerConfigForm';

interface JumpsellerConnectTabsProps {
    onSuccess?: () => void;
    onCancel?: () => void;
    integrationTypeBaseURLTest?: string;
}

type Tab = 'oauth' | 'manual';

const GREEN = 'var(--color-primary)';
const INPUT_BORDER = '#e9e9f0';

export function JumpsellerConnectTabs({ onSuccess, onCancel, integrationTypeBaseURLTest }: JumpsellerConnectTabsProps) {
    const [tab, setTab] = useState<Tab>('oauth');

    const tabBtn = (active: boolean) =>
        `flex-1 flex items-center justify-center gap-2 rounded-lg py-2 text-[13px] font-semibold transition-colors ${active
            ? 'text-white'
            : 'text-gray-600 dark:text-gray-300 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700'
        }`;

    return (
        <div className="space-y-3 w-full">
            <div
                className="flex gap-2 p-1 rounded-xl"
                style={{ backgroundColor: '#f3f4f6', border: `1px solid ${INPUT_BORDER}` }}
            >
                <button
                    type="button"
                    onClick={() => setTab('oauth')}
                    className={tabBtn(tab === 'oauth')}
                    style={tab === 'oauth' ? { backgroundColor: GREEN } : undefined}
                >
                    <LinkIcon className="w-4 h-4" />
                    OAuth
                </button>
                <button
                    type="button"
                    onClick={() => setTab('manual')}
                    className={tabBtn(tab === 'manual')}
                    style={tab === 'manual' ? { backgroundColor: GREEN } : undefined}
                >
                    <KeyIcon className="w-4 h-4" />
                    Login + Auth Token
                </button>
            </div>

            {tab === 'oauth' ? (
                <JumpsellerOAuthForm onCancel={onCancel} />
            ) : (
                <JumpsellerConfigForm
                    onSuccess={onSuccess}
                    onCancel={onCancel}
                    integrationTypeBaseURLTest={integrationTypeBaseURLTest}
                />
            )}
        </div>
    );
}
