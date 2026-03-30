'use client';

import { Menu } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useState } from 'react';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '@/components/ui/sheet';
import { cn } from '@/lib/utils';

const navLinks = [
	{ href: '/', label: 'Home' },
	{ href: '/coffees', label: 'Coffees' },
	{ href: '/roasters', label: 'Roasters' },
	{ href: '/cafes', label: 'Cafes' },
	{ href: '/countries', label: 'Origins' },
	{ href: '/guide', label: 'Guide' },
	{ href: '/find', label: 'Find' },
	{ href: '/my-coffees', label: 'My Coffees' },
] as const;

export function MobileNav() {
	const [open, setOpen] = useState(false);
	const pathname = usePathname();

	return (
		<Sheet open={open} onOpenChange={setOpen}>
			<SheetTrigger aria-label="Open menu">
				<Menu className="size-6 text-gold" />
			</SheetTrigger>
			<SheetContent side="right" className="bg-rich-mahogany border-border/50">
				<SheetHeader>
					<SheetTitle className="text-gold tracking-widest uppercase text-sm">Navigation</SheetTitle>
				</SheetHeader>
				<nav className="flex flex-col gap-1 px-4">
					{navLinks.map((link) => (
						<Link
							key={link.href}
							href={link.href}
							onClick={() => setOpen(false)}
							className={cn(
								'rounded px-3 py-2.5 text-sm font-medium uppercase tracking-wider transition-colors',
								pathname === link.href
									? 'bg-gold/10 text-gold'
									: 'text-snow/60 hover:text-gold',
							)}
						>
							{link.label}
						</Link>
					))}
				</nav>
			</SheetContent>
		</Sheet>
	);
}
