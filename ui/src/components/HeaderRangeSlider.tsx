import type { FilterRange } from '@/types/summary'
import { Slider } from '@/ui/slider'

export default function HeaderRangeSlider({
    range,
    onRangeUpdate,
    max = 100,
    unit = '%',
}: {
    range: FilterRange
    onRangeUpdate: (vals: [number, number]) => void
    /** Upper bound of the slider. 100 for percentages; data-derived for value metrics. */
    max?: number
    /** Suffix shown next to the numbers ('%' for percentages, '' for scalar values). */
    unit?: string
}) {
    // Fully controlled: value derived from props, every change reported to the parent.
    const currentRange: [number, number] = [range.min, range.max]

    return (
        <div className="space-y-1">
            <div className="flex items-center justify-between">
                <span className="font-medium text-foreground tabular-nums">
                    {currentRange[0]}
                    {unit} – {currentRange[1]}
                    {unit}
                </span>
            </div>
            <Slider
                value={currentRange}
                onValueChange={(vals) => onRangeUpdate([vals[0] ?? 0, vals[1] ?? max])}
                max={max}
                min={0}
                step={1}
            />
        </div>
    )
}
