import { GitBranch } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { LineDetails, LineStatus } from '@/types/summary'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/ui/tooltip'

// Define classes for the left-side border/block color
const borderClasses: Record<LineStatus, string> = {
    covered: 'border-covered',
    uncovered: 'border-uncovered',
    partial: 'border-partial',
    'not-coverable': 'border-transparent',
}

// Define classes for the subtle line background color
const bgClasses: Record<LineStatus, string> = {
    covered: 'bg-covered/10',
    uncovered: 'bg-uncovered/10',
    partial: 'bg-partial/10',
    'not-coverable': 'bg-transparent',
}

export default function CodeLine({ lineNumber, content, status, hits, branchInfo }: LineDetails) {
    const hasHitCount = typeof hits === 'number'

    return (
        <TooltipProvider delayDuration={100}>
            <div
                className={cn(
                    'flex w-full items-start font-mono text-sm',
                    bgClasses[status],
                    'hover:bg-black/5 dark:hover:bg-white/5',
                )}
                data-line-number={lineNumber}
            >
                <div
                    className={cn(
                        'sticky left-0 z-10 flex select-none items-center border-l-4 pr-4',
                        borderClasses[status],
                        bgClasses[status],
                    )}
                >
                    <span className="w-10 text-right font-mono text-muted-foreground text-xs">
                        {hasHitCount ? hits : ''}
                    </span>
                    <span className="w-12 text-right text-muted-foreground">{lineNumber}</span>
                </div>

                <div className="flex flex-1 items-center py-0.5 pl-4">
                    <div className="w-6">
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
                    <pre className="whitespace-pre-wrap break-words">{content || ' '}</pre>
                </div>
            </div>
        </TooltipProvider>
    )
}
