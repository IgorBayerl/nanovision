import { AlertCircle, AlertTriangle, Shield } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { RiskFilter } from '@/types/summary'
import { Button } from '@/ui/button'

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
                    {opt === 'danger' && <AlertCircle className="h-4 w-4 text-destructive" />}
                    {opt === 'warning' && <AlertTriangle className="h-4 w-4 text-accent" />}
                    {opt === 'safe' && <Shield className="h-4 w-4 text-primary" />}
                </Button>
            ))}
        </div>
    )
}
