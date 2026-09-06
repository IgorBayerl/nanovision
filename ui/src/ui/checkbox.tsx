import * as CheckboxPrimitive from '@radix-ui/react-checkbox'
import { CheckIcon, MinusIcon } from 'lucide-react'
import type * as React from 'react'
import { cn } from '@/lib/utils'

/**
 * Pass `checked="indeterminate"` for a parent that governs a group of boxes
 * where only some children are ticked; it renders a dash instead of a check.
 */
function Checkbox({ className, ...props }: React.ComponentProps<typeof CheckboxPrimitive.Root>) {
    return (
        <CheckboxPrimitive.Root
            data-slot="checkbox"
            className={cn(
                'group/checkbox peer size-4 shrink-0 rounded-sm border border-input shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 data-[state=checked]:border-primary data-[state=indeterminate]:border-primary data-[state=checked]:bg-primary data-[state=indeterminate]:bg-primary data-[state=checked]:text-primary-foreground data-[state=indeterminate]:text-primary-foreground dark:bg-input/30 dark:data-[state=checked]:bg-primary dark:data-[state=indeterminate]:bg-primary dark:aria-invalid:ring-destructive/40',
                className,
            )}
            {...props}
        >
            <CheckboxPrimitive.Indicator
                data-slot="checkbox-indicator"
                className="flex items-center justify-center text-current"
            >
                <CheckIcon className="size-3.5 group-data-[state=indeterminate]/checkbox:hidden" />
                <MinusIcon className="hidden size-3.5 group-data-[state=indeterminate]/checkbox:block" />
            </CheckboxPrimitive.Indicator>
        </CheckboxPrimitive.Root>
    )
}

export { Checkbox }
