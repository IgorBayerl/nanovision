import { List, ListTree } from 'lucide-react'
import type { RefObject } from 'react'
import ColumnsMenu from '@/components/Toolbar.ColumnsMenu'
import RiskSegment from '@/components/Toolbar.RiskSegment'
import SearchBox from '@/components/Toolbar.SearchBox'
import type { FileExplorerState, Setters } from '@/hooks/useFileExplorerState'
import type { MetricConfig } from '@/types/summary'
import { Button } from '@/ui/button'
import { CardTitle } from '@/ui/card'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/ui/tooltip'

interface ToolbarProps {
    state: FileExplorerState
    setters: Setters
    searchRef: RefObject<HTMLInputElement | null>
    metricConfigs: MetricConfig[]
}

export default function FileExplorerToolbar({ state, setters, searchRef, metricConfigs }: ToolbarProps) {
    return (
        <div className="flex flex-col gap-3">
            <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
                <CardTitle className="text-lg">File Coverage</CardTitle>
                <div className="flex flex-wrap items-center gap-2">
                    <TooltipProvider delayDuration={100}>
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    className="h-8 w-8 rounded-sm p-0"
                                    onClick={() => setters.setViewMode(state.viewMode === 'tree' ? 'flat' : 'tree')}
                                    aria-label={`Switch to ${state.viewMode === 'tree' ? 'flat' : 'tree'} view`}
                                >
                                    {state.viewMode === 'tree' ? (
                                        <ListTree className="h-4 w-4" />
                                    ) : (
                                        <List className="h-4 w-4" />
                                    )}
                                </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                                <p>Switch to {state.viewMode === 'tree' ? 'Flat' : 'Tree'} View</p>
                            </TooltipContent>
                        </Tooltip>
                    </TooltipProvider>
                    <ColumnsMenu metricConfigs={metricConfigs} onToggleMetric={setters.toggleMetric} />
                </div>
            </div>
            <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
                <SearchBox
                    ref={searchRef}
                    value={state.query}
                    onChange={setters.setQuery}
                    mode={state.searchMode}
                    onModeChange={setters.setSearchMode}
                />
                <RiskSegment value={state.riskFilter} onChange={setters.setRiskFilter} />
            </div>
        </div>
    )
}
