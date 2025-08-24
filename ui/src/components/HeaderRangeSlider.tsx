import { useState } from 'react'
import type { FilterRange } from '@/types/summary'
import { Slider } from '@/ui/slider'

export default function HeaderRangeSlider({
    range,
    onRangeCommit,
}: {
    range: FilterRange
    onRangeCommit: (vals: [number, number]) => void
}) {
    // Local state provides smooth UI updates during dragging without re-rendering the whole table.
    const [displayRange, setDisplayRange] = useState<[number, number]>([range.min, range.max])

    return (
        <div className="space-y-1">
            <div className="flex items-center justify-between">
                <span className="font-medium text-foreground tabular-nums">
                    {displayRange[0]}% – {displayRange[1]}%
                </span>
            </div>
            <Slider
                value={displayRange}
                // This updates the local UI smoothly while dragging.
                onValueChange={(vals) => setDisplayRange([vals[0] ?? 0, vals[1] ?? 100])}
                // This commits the final value to the parent state, triggering the filter logic.
                onValueCommit={(vals) => onRangeCommit([vals[0] ?? 0, vals[1] ?? 100])}
                max={100}
                min={0}
                step={1}
            />
        </div>
    )
}
