import React, { useState } from 'react';
import { Copy, Check } from 'lucide-react';

const CopyCode = ({ code, singleLine = false }) => {
    const [copied, setCopied] = useState(false);

    const handleCopy = () => {
        if (!navigator.clipboard) {
            return;
        }
        navigator.clipboard
            .writeText(code)
            .then(() => {
                setCopied(true);
                setTimeout(() => setCopied(false), 2000);
            })
            .catch(() => {
                // ignore
            });
    };

    return (
        <div className="relative bg-card border border-border rounded-lg">
            <button
                type="button"
                onClick={handleCopy}
                className="absolute top-2 right-2 inline-flex items-center gap-1 px-2 py-1 rounded-md bg-background/80 text-[11px] text-muted-foreground hover:text-foreground hover:bg-secondary/80 border border-border cursor-pointer"
            >
                {copied ? <Check size={12} /> : <Copy size={12} />}
                <span>{copied ? 'Copied' : 'Copy'}</span>
            </button>
            <pre
                className={`m-0 p-3 font-mono text-sm overflow-x-auto ${singleLine ? 'whitespace-nowrap' : 'whitespace-pre'
                    }`}
                style={{ scrollbarWidth: 'none' }}
            >
                <code>{code}</code>
            </pre>
        </div>
    );
};

export default CopyCode;
