import { cn } from '@/lib/utils'
import type { RiskFilter } from '@/types/summary'
import { Button } from '@/ui/button'
import { StatusIcon } from './MetricCard'

export default function RiskSegment({ value, onChange }: { value: RiskFilter; onChange: (v: RiskFilter) => void }) {
    const opts: RiskFilter[] = ['all', 'danger', 'warning', 'safe']
    return (
        <div className="flex items-center gap-1 rounded-sm border border-border p-1">
            {opts.map((opt) => (
                <Button
                    key={opt}
                    size="sm"
                    variant={value === opt ? 'default' : 'ghost'}
                    className={cn('h-8 rounded-sm px-2', value === opt ? 'bg-primary text-primary-foreground' : '')}
                    onClick={() => onChange(opt)}
                    title={`Show ${opt}`}
                >
                    {opt === 'all' && <span className="text-xs">All</span>}
                    {opt === 'danger' && <StatusIcon status="danger" />}
                    {opt === 'warning' && <StatusIcon status="warning" />}
                    {opt === 'safe' && <StatusIcon status="safe" />}
                </Button>
            ))}
        </div>
    )
}
