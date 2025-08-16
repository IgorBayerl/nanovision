import { Filter } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { FilterRange, MetricConfig } from '@/types/summary'
import { Button } from '@/ui/button'
import { Input } from '@/ui/input'
import { Label } from '@/ui/label'
import { Popover, PopoverContent, PopoverTrigger } from '@/ui/popover'
import { Slider } from '@/ui/slider'

const displayRange = (r: FilterRange) => (r.min === 0 && r.max === 100 ? 'All' : `${r.min}–${r.max}%`)

export default function MetricRangeChip({
    cfg,
    value,
    onChange,
}: {
    cfg: MetricConfig
    value: FilterRange
    onChange: (vals: [number, number]) => void
}) {
    const isAll = value.min === 0 && value.max === 100

    return (
        <Popover>
            <PopoverTrigger asChild>
                <button
                    type="button"
                    className={cn(
                        'inline-flex items-center gap-1 rounded-sm border px-2 py-0.5 text-[11px] leading-5',
                        isAll
                            ? 'border-border bg-transparent text-muted-foreground hover:border-primary hover:bg-primary/10'
                            : 'border-primary bg-primary text-primary-foreground hover:bg-primary/90',
                    )}
                    title={`${cfg.label} range`}
                >
                    <Filter className="h-3.5 w-3.5" />
                    <span className="tabular-nums">{displayRange(value)}</span>
                </button>
            </PopoverTrigger>
            <PopoverContent className="w-64 rounded-md" align="center">
                <div className="space-y-3">
                    <div className="flex items-center justify-between">
                        <div className="font-medium text-sm">{cfg.label}</div>
                        {!isAll && (
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-7 rounded-sm px-2"
                                onClick={() => onChange([0, 100])}
                            >
                                Reset
                            </Button>
                        )}
                    </div>

                    <div className="flex items-center justify-between text-muted-foreground text-xs">
                        <span>0%</span>
                        <span className="font-medium text-foreground tabular-nums">
                            {value.min}% – {value.max}%
                        </span>
                        <span>100%</span>
                    </div>

                    <Slider
                        value={[value.min, value.max]}
                        onValueChange={(vals) => onChange([vals[0] ?? 0, vals[1] ?? 100])}
                        max={100}
                        min={0}
                        step={1}
                    />

                    <div className="flex items-center gap-2">
                        <Label htmlFor={`min-${cfg.id}`} className="text-muted-foreground text-xs">
                            Min
                        </Label>
                        <Input
                            id={`min-${cfg.id}`}
                            type="number"
                            min={0}
                            max={value.max}
                            value={value.min}
                            onChange={(e) =>
                                onChange([Math.max(0, Math.min(100, Number(e.target.value) || 0)), value.max])
                            }
                            className="h-8 w-16 rounded-sm text-xs"
                        />
                        <Label htmlFor={`max-${cfg.id}`} className="ml-2 text-muted-foreground text-xs">
                            Max
                        </Label>
                        <Input
                            id={`max-${cfg.id}`}
                            type="number"
                            min={value.min}
                            max={100}
                            value={value.max}
                            onChange={(e) =>
                                onChange([value.min, Math.max(0, Math.min(100, Number(e.target.value) || 0))])
                            }
                            className="h-8 w-16 rounded-sm text-xs"
                        />
                    </div>
                </div>
            </PopoverContent>
        </Popover>
    )
}
