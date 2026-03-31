'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';

const links = [
	{ href: '/admin', label: 'Dashboard' },
	{ href: '/admin/roasters', label: 'Roasters' },
	{ href: '/admin/coffees', label: 'Coffees' },
	{ href: '/admin/cafes', label: 'Cafes' },
];

export function AdminNav() {
	const pathname = usePathname();

	return (
		<nav className="mb-8 flex gap-4 border-b border-border pb-4">
			{links.map((link) => {
				const active = link.href === '/admin' ? pathname === '/admin' : pathname.startsWith(link.href);
				return (
					<Link
						key={link.href}
						href={link.href}
						className={cn(
							'text-sm font-semibold uppercase tracking-widest transition-colors',
							active ? 'text-foreground' : 'text-muted-foreground hover:text-foreground',
						)}
					>
						{link.label}
					</Link>
				);
			})}
		</nav>
	);
}
