import { useId } from 'react'
import InfoTooltip from '@/components/InfoTooltip'
import type { ReportSelectionState } from '@/hooks/useReportSelection'
import { cn } from '@/lib/utils'
import { Button } from '@/ui/button'
import { Checkbox } from '@/ui/checkbox'
import { Label } from '@/ui/label'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/ui/tooltip'

interface ReportsSelectorProps {
    state: ReportSelectionState
    /**
     * Metrics that cannot be attributed to a single report and keep their
     * merged value, named so the numbers are not silently misleading.
     */
    frozenMetricLabels?: string[]
}

/**
 * A two-level checkbox tree: one parent governing every parsed report. The
 * parent shows a dash while only some are ticked, and ticks or clears the whole
 * list when clicked.
 *
 * Reports are named by file name alone, so the full path lives on hover.
 */
export default function ReportsSelector({ state, frozenMetricLabels = [] }: ReportsSelectorProps) {
    const { reports, isSelected, isNoneSelected, groupCheckState, toggle, toggleAll, selectOnly } = state
    const checkboxId = useId()

    if (reports.length < 2) return null

    const selectedCount = reports.filter((_, index) => isSelected(index)).length

    return (
        <div className="flex flex-col gap-1.5">
            <div className="flex items-center gap-2">
                <Checkbox id={`${checkboxId}-all`} checked={groupCheckState} onCheckedChange={toggleAll} />
                <Label
                    htmlFor={`${checkboxId}-all`}
                    className="font-medium text-muted-foreground text-xs uppercase tracking-wide"
                >
                    Reports
                </Label>
                <InfoTooltip label="About the report selection" tone={isNoneSelected ? 'warning' : 'muted'}>
                    {isNoneSelected ? (
                        <p>No report is selected, so nothing counts as covered. Tick at least one.</p>
                    ) : (
                        <div className="flex flex-col gap-1.5">
                            <p>Every metric on this page is recomputed from the reports ticked here.</p>
                            {frozenMetricLabels.length > 0 && (
                                <p>
                                    {frozenMetricLabels.join(' and ')} {frozenMetricLabels.length === 1 ? 'is' : 'are'}{' '}
                                    not tracked per report and always shows the merged value.
                                </p>
                            )}
                        </div>
                    )}
                </InfoTooltip>
                <span className="ml-auto shrink-0 text-muted-foreground text-xs tabular-nums">
                    {selectedCount}/{reports.length}
                </span>
            </div>

            {/* Indented, with a rule standing in for the tree's trunk. */}
            <div className="ml-2 flex flex-col gap-1.5 border-border border-l pl-3">
                {reports.map((report, index) => (
                    <div key={`${index}-${report.path}`} className="group/report flex items-center gap-2">
                        <Checkbox
                            id={`${checkboxId}-${index}`}
                            checked={isSelected(index)}
                            onCheckedChange={() => toggle(index)}
                        />
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Label
                                    htmlFor={`${checkboxId}-${index}`}
                                    className={cn(
                                        'min-w-0 flex-1 truncate text-xs',
                                        // Only details pages set `relevant`; a
                                        // report that never touched this file is
                                        // dimmed rather than hidden, so the list
                                        // stays the same everywhere.
                                        report.relevant === false && 'text-muted-foreground/60',
                                    )}
                                >
                                    {report.name}
                                </Label>
                            </TooltipTrigger>
                            {/* To the side: above, it would cover the rest of
                                the list while picking through it. */}
                            <TooltipContent side="right" sideOffset={6} className="max-w-sm break-all font-mono">
                                {report.path}
                            </TooltipContent>
                        </Tooltip>
                        {/* Named group: the sidebar is itself a plain `group`,
                            so an unnamed group-hover would reveal every row at
                            once. The fade lives on this span rather than the
                            button, whose opacity Tailwind's preflight sets. */}
                        <span className="shrink-0 opacity-0 transition-opacity focus-within:opacity-100 group-hover/report:opacity-100">
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-5 px-1.5 text-muted-foreground text-xs"
                                onClick={() => selectOnly(index)}
                                title={`Show ${report.name} only`}
                            >
                                only
                            </Button>
                        </span>
                    </div>
                ))}
            </div>
        </div>
    )
}
