import { cn } from '@/lib/utils'
import type { DiffStatus, LineDetails, LineStatus } from '@/types/summary'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/ui/tooltip'
import { GitBranch } from 'lucide-react'

const squareBgClasses: Record<LineStatus, string> = {
    covered: 'bg-covered',
    uncovered: 'bg-uncovered',
    partial: 'bg-partial',
    'not-coverable': 'bg-transparent',
}

const lineBgClasses: Record<LineStatus, string> = {
    covered: 'bg-covered/20',
    uncovered: 'bg-uncovered/20',
    partial: 'bg-partial/20',
    'not-coverable': 'bg-transparent',
}

const diffSymbols: Partial<Record<DiffStatus, string>> = {
    added: '+',
    removed: '-',
}

const diffSymbolClasses: Partial<Record<DiffStatus, string>> = {
    added: 'text-covered',
    removed: 'text-uncovered',
}

const gridTemplateColumns = '1.5rem 4rem 4rem 1.5rem 1.5rem 1fr'

interface CodeLineProps extends Omit<LineDetails, 'hits'> {
    hits?: number
}

export default function CodeLine({ lineNumber, content, status, hits, branchInfo, diffStatus }: CodeLineProps) {
    const hasHitCount = typeof hits === 'number' && hits > 0
    const diffSymbol = diffStatus ? diffSymbols[diffStatus] : undefined
    const diffClass = diffStatus ? diffSymbolClasses[diffStatus] : undefined

    return (
        <TooltipProvider delayDuration={100}>
            <div
                className={cn('group grid cursor-default items-stretch font-mono text-sm', lineBgClasses[status])}
                style={{ gridTemplateColumns }}
                data-line-number={lineNumber}
            >
                {/* 1. Status Color Square */}
                <div className="flex items-center justify-center">
                    <div className={cn('h-full w-2', squareBgClasses[status])} />
                </div>

                {/* 2. Line Number */}
                <div className="select-none border-border/30 border-r py-0.5 pr-4 text-right font-extrabold text-muted-foreground">
                    {lineNumber}
                </div>

                {/* 3. Hit Count Badge */}
                <div className="flex items-center justify-center border-border/30 border-r px-4 py-0.5">
                    {hasHitCount && (
                        <span className="select-none rounded-md bg-background px-1.5 text-center font-medium font-sans text-muted-foreground text-xs">
                            {hits}
                        </span>
                    )}
                </div>

                {/* 4. Branch Indicator */}
                <div className="flex items-center justify-center py-0.5">
                    {branchInfo && (
                        <Tooltip>
                            <TooltipTrigger>
                                <GitBranch className="h-4 w-4 text-partial" />
                            </TooltipTrigger>
                            <TooltipContent>
                                <p>
                                    Branch Coverage: {branchInfo.covered} / {branchInfo.total}
                                </p>
                            </TooltipContent>
                        </Tooltip>
                    )}
                </div>

                {/* 5. Diff Indicator*/}
                <div
                    className={cn('select-none border-border/30 border-r px-1 py-0.5 text-center font-bold', diffClass)}
                >
                    {diffSymbol}
                </div>

                {/* 6. Source Code */}
                <div className="py-0.5 pl-4">
                    <pre className="whitespace-pre">{content || ' '}</pre>
                </div>
            </div>
        </TooltipProvider>
    )
}
