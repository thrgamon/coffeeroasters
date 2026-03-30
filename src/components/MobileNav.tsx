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
			<SheetContent side="right" className="bg-ink border-2 border-gold">
				<SheetHeader>
					<SheetTitle className="text-gold tracking-[0.2em] uppercase text-sm font-bold">Navigation</SheetTitle>
				</SheetHeader>
				<nav className="flex flex-col gap-1 px-4">
					{navLinks.map((link) => (
						<Link
							key={link.href}
							href={link.href}
							onClick={() => setOpen(false)}
							className={cn(
								'px-3 py-2.5 text-sm font-bold uppercase tracking-[0.15em] transition-colors',
								pathname === link.href ? 'bg-gold text-ink' : 'text-paper/70 hover:text-gold',
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
