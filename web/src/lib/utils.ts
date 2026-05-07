import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function keys<T extends {}>(obj: T): (keyof T)[] {
  return Object.keys(obj) as (keyof T)[]
}

export function entries<T extends {}>(obj: T): [(keyof T), T[keyof T]][] {
  return Object.entries(obj) as [(keyof T), T[keyof T]][]
}
