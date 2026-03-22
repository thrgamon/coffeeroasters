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
	{ href: '/countries', label: 'Countries' },
	{ href: '/my-coffees', label: 'My coffees' },
] as const;

export function MobileNav() {
	const [open, setOpen] = useState(false);
	const pathname = usePathname();

	return (
		<Sheet open={open} onOpenChange={setOpen}>
			<SheetTrigger aria-label="Open menu">
				<Menu className="size-6 text-muted-foreground" />
			</SheetTrigger>
			<SheetContent side="right">
				<SheetHeader>
					<SheetTitle>Navigation</SheetTitle>
				</SheetHeader>
				<nav className="flex flex-col gap-2 px-4">
					{navLinks.map((link) => (
						<Link
							key={link.href}
							href={link.href}
							onClick={() => setOpen(false)}
							className={cn(
								'rounded-md px-3 py-2 text-sm transition-colors',
								pathname === link.href
									? 'bg-accent font-medium text-foreground'
									: 'text-muted-foreground hover:text-foreground',
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
