import { cn } from '@/lib/utils'
import type { DiffFilter } from '@/types/summary'
import { Button } from '@/ui/button'

export default function DiffFilterSegment({
    value,
    onChange,
}: {
    value: DiffFilter
    onChange: (v: DiffFilter) => void
}) {
    const opts: { value: DiffFilter; label: string }[] = [
        { value: 'all', label: 'All Files' },
        { value: 'changed', label: 'Changed Files' },
    ]
    return (
        <div className="flex items-center gap-1 rounded-sm border border-border bg-subtle p-1">
            {opts.map((opt) => (
                <Button
                    key={opt.value}
                    size="sm"
                    variant={value === opt.value ? 'default' : 'ghost'}
                    className={cn(
                        'h-8 rounded-sm px-3 text-xs',
                        value === opt.value ? 'bg-background text-foreground shadow-sm' : '',
                    )}
                    onClick={() => onChange(opt.value)}
                    title={`Show ${opt.label}`}
                >
                    {opt.label}
                </Button>
            ))}
        </div>
    )
}
