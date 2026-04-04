import type { Metadata } from 'next';
import './globals.css';
import Link from 'next/link';
import { MobileNav } from '@/components/MobileNav';
import { SearchButton, SearchPalette } from '@/components/SearchPalette';
import { UserNav } from '@/components/UserNav';
import { AuthProvider } from '@/lib/auth-context';
import { CoffeeTrackerProvider } from '@/lib/coffee-tracker';
import { QueryProvider } from '@/lib/query-provider';

export const metadata: Metadata = {
	title: 'COFFEEROASTERS',
	description: 'Discover Australian specialty coffee roasters',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
	return (
		<html lang="en" suppressHydrationWarning>
			<body className="flex min-h-screen flex-col bg-background font-sans antialiased">
				<QueryProvider>
					<AuthProvider>
						<CoffeeTrackerProvider>
							<nav className="bg-ink">
								<div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-5">
									<Link href="/" className="text-xl font-bold tracking-[0.25em] uppercase text-gold">
										Coffeeroasters
									</Link>
									<div className="hidden gap-8 text-xs font-bold tracking-[0.2em] uppercase sm:flex">
										<Link href="/coffees" className="text-paper/70 hover:text-gold transition-colors">
											Coffees
										</Link>
										<Link href="/roasters" className="text-paper/70 hover:text-gold transition-colors">
											Roasters
										</Link>
										<Link href="/cafes" className="text-paper/70 hover:text-gold transition-colors">
											Cafes
										</Link>
										<Link href="/countries" className="text-paper/70 hover:text-gold transition-colors">
											Origins
										</Link>
										<Link href="/guide" className="text-paper/70 hover:text-gold transition-colors">
											Guide
										</Link>
										<Link href="/find" className="text-paper/70 hover:text-gold transition-colors">
											Find
										</Link>
										<Link href="/my-coffees" className="text-paper/70 hover:text-gold transition-colors">
											My Coffees
										</Link>
										<SearchButton />
										<UserNav />
									</div>
									<div className="sm:hidden">
										<MobileNav />
									</div>
								</div>
							</nav>
							<SearchPalette />
							<main className="mx-auto flex-1 max-w-6xl px-4 py-10">{children}</main>
							<footer className="bg-ink">
								<div className="mx-auto max-w-6xl px-4 py-6 text-center text-xs tracking-[0.25em] uppercase font-bold text-warm-grey">
									&copy; {new Date().getFullYear()} Coffeeroasters
								</div>
							</footer>
						</CoffeeTrackerProvider>
					</AuthProvider>
				</QueryProvider>
			</body>
		</html>
	);
}
