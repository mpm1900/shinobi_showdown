import * as React from 'react'
import { useDebounce } from 'use-debounce'

import { cn } from '#/lib/utils'

function Input({
  className,
  type,
  onChange,
  onValueChange,
  ...props
}: React.ComponentProps<'input'> & {
  onValueChange?: (value: React.ComponentProps<'input'>['value']) => void
}) {
  return (
    <input
      type={type}
      data-slot="input"
      className={cn(
        'h-9 w-full min-w-0 rounded-md border border-input bg-transparent px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none selection:bg-primary selection:text-primary-foreground file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm dark:bg-input/30',
        'focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50',
        'aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40',
        className
      )}
      onChange={(e) => {
        onChange?.(e)
        onValueChange?.(e.target.value)
      }}
      {...props}
    />
  )
}

function DebouncedInput({
  delay,
  value,
  onChange,
  onValueChange,
  ...props
}: React.ComponentProps<'input'> & {
  delay: number
  onValueChange: (v: React.ComponentProps<'input'>['value']) => void
}) {
  const [text, setText] = React.useState(value)
  const [debounced] = useDebounce(text, delay)

  React.useEffect(() => {
    onValueChange(debounced)
  }, [debounced])

  return (
    <Input
      {...props}
      value={text}
      onChange={(e) => {
        setText(e.target.value)
        onChange?.(e)
      }}
    />
  )
}

export { Input, DebouncedInput }
