import { Asterisk, Search, Type } from 'lucide-react'
import { forwardRef } from 'react'
import { Button } from '@/ui/button'
import { Input } from '@/ui/input'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/ui/tooltip'

const SearchBox = forwardRef<
    HTMLInputElement,
    {
        value: string
        onChange: (v: string) => void
        mode: 'glob' | 'normal'
        onModeChange: (mode: 'glob' | 'normal') => void
    }
>(({ value, onChange, mode, onModeChange }, ref) => {
    const isGlob = mode === 'glob'

    return (
        <TooltipProvider delayDuration={200}>
            <div className="relative flex w-full items-center gap-2 md:w-[460px]">
                <Search className="absolute left-2 h-4 w-4 text-muted-foreground" />
                <Input
                    ref={ref}
                    placeholder="Search files or paths (Ctrl/Cmd+F)"
                    className="h-8 rounded-sm pr-10 pl-8"
                    value={value}
                    onChange={(e) => onChange(e.target.value)}
                />
                <div className="absolute right-1.5 flex items-center">
                    <Tooltip>
                        <TooltipTrigger asChild>
                            <Button
                                variant="ghost"
                                size="sm"
                                className={
                                    'h-6 w-6 rounded-sm p-0 text-muted-foreground transition-colors hover:bg-primary/10 hover:text-primary'
                                }
                                onClick={() => onModeChange(isGlob ? 'normal' : 'glob')}
                                aria-label="Toggle search mode"
                            >
                                {isGlob ? <Asterisk className="h-4 w-4" /> : <Type className="h-4 w-4" />}
                            </Button>
                        </TooltipTrigger>
                        <TooltipContent>
                            <p>Click to switch to {isGlob ? 'Normal' : 'Glob'} search</p>
                        </TooltipContent>
                    </Tooltip>
                </div>
            </div>
        </TooltipProvider>
    )
})

export default SearchBox
