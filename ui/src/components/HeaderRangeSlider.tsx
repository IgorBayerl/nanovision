import type { FilterRange } from '@/types/summary'
import { Slider } from '@/ui/slider'

export default function HeaderRangeSlider({
    range,
    onRangeUpdate,
}: {
    range: FilterRange
    onRangeUpdate: (vals: [number, number]) => void
}) {
    // This component is now fully controlled. Its value is derived directly
    // from props, and it reports every change back to the parent.
    const currentRange: [number, number] = [range.min, range.max]

    return (
        <div className="space-y-1">
            <div className="flex items-center justify-between">
                <span className="font-medium text-foreground tabular-nums">
                    {currentRange[0]}% – {currentRange[1]}%
                </span>
            </div>
            <Slider
                value={currentRange}
                // onValueChange now directly calls the parent's update function in real-time.
                onValueChange={(vals) => onRangeUpdate([vals[0] ?? 0, vals[1] ?? 100])}
                max={100}
                min={0}
                step={1}
            />
        </div>
    )
}
