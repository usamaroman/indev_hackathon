'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import clsx from 'clsx'

const navItems = [
  { href: '/map', label: 'Карта' },
  { href: '/load', label: 'Загруженность' },
]

export default function Navbar() {
  const pathname = usePathname()

  return (
    <nav className="flex space-x-4 p-4 bg-gray-100 shadow">
      {navItems.map(({ href, label }) => (
        <Link
          key={href}
          href={href}
          className={clsx(
            'px-4 py-2 rounded hover:bg-blue-100',
            pathname === href && 'bg-blue-500 text-white'
          )}
        >
          {label}
        </Link>
      ))}
    </nav>
  )
}
