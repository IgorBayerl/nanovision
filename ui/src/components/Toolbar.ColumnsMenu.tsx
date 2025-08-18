import { Columns } from 'lucide-react'
import type { MetricConfig } from '@/types/summary'
import { Button } from '@/ui/button'
import { Checkbox } from '@/ui/checkbox'
import { Label } from '@/ui/label'
import { Popover, PopoverContent, PopoverTrigger } from '@/ui/popover'

export default function ColumnsMenu({
    metricConfigs,
    onToggleMetric,
}: {
    metricConfigs: MetricConfig[]
    onToggleMetric: (id: MetricConfig['id']) => void
}) {
    return (
        <Popover>
            <PopoverTrigger asChild>
                <Button variant="outline" className="gap-2 rounded-sm font-normal" size="sm">
                    <Columns className="h-4 w-4" />
                    Columns
                </Button>
            </PopoverTrigger>
            <PopoverContent align="end" className="w-56 rounded-md">
                <div className="space-y-2">
                    <div className="font-medium text-sm">Visible Metrics</div>
                    <div className="space-y-2">
                        {metricConfigs.map((cfg) => {
                            const checkboxId = `metric-checkbox-${cfg.id}`
                            return (
                                <div key={cfg.id} className="flex items-center gap-2">
                                    <Checkbox
                                        id={checkboxId}
                                        checked={cfg.enabled}
                                        onCheckedChange={() => onToggleMetric(cfg.id)}
                                    />
                                    <Label htmlFor={checkboxId} className="cursor-pointer text-sm">
                                        {cfg.label}
                                    </Label>
                                </div>
                            )
                        })}
                    </div>
                </div>
            </PopoverContent>
        </Popover>
    )
}
