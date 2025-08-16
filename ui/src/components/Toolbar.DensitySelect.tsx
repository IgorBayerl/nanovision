import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/ui/select'

export default function DensitySelect({
    value,
    onValueChange,
}: {
    value: 'comfortable' | 'compact'
    onValueChange: (v: 'comfortable' | 'compact') => void
}) {
    return (
        <Select value={value} onValueChange={onValueChange}>
            <SelectTrigger className="h-8 w-[140px] rounded-sm" size="sm">
                <SelectValue placeholder="Density" />
            </SelectTrigger>
            <SelectContent>
                <SelectItem value="comfortable">Comfortable</SelectItem>
                <SelectItem value="compact">Compact</SelectItem>
            </SelectContent>
        </Select>
    )
}
