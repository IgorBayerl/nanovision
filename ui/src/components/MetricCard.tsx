import { Card, CardContent } from '@/ui/card'

export default function MetricCard({ label, value }: { label: string; value: number | undefined }) {
    const pct = value !== undefined ? Math.max(0, Math.min(100, Math.round(value))) : undefined

    return (
        <Card className="rounded-md">
            <CardContent className="p-5">
                <div className="flex items-center gap-5">
                    <div className="relative h-28 w-4 overflow-hidden rounded-xs bg-muted">
                        <div
                            className="absolute right-0 bottom-0 left-0 bg-primary"
                            style={{ height: pct !== undefined ? `${pct}%` : '0%' }}
                        />
                    </div>
                    <div className="flex-1">
                        <div className="font-medium text-muted-foreground text-sm">{label}</div>
                        <div className="mt-1 font-extrabold text-4xl text-foreground tabular-nums tracking-tight md:text-5xl">
                            {pct !== undefined ? `${pct}%` : 'N/A'}
                        </div>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
