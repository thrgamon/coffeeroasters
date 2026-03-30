import type { Metadata } from 'next';
import './globals.css';
import Link from 'next/link';
import { MobileNav } from '@/components/MobileNav';
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
		<html lang="en">
			<body className="flex min-h-screen flex-col bg-background font-sans antialiased">
				<QueryProvider>
					<AuthProvider>
						<CoffeeTrackerProvider>
							<nav className="border-b border-border/60 bg-rich-mahogany/80 backdrop-blur-md sticky top-0 z-50">
								<div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-4">
									<Link href="/" className="text-lg font-bold tracking-widest uppercase text-gold glow-gold">
										Coffeeroasters
									</Link>
									<div className="hidden gap-8 text-sm font-medium tracking-wide uppercase sm:flex">
										<Link href="/coffees" className="text-snow/60 hover:text-gold transition-colors">
											Coffees
										</Link>
										<Link href="/roasters" className="text-snow/60 hover:text-gold transition-colors">
											Roasters
										</Link>
										<Link href="/cafes" className="text-snow/60 hover:text-gold transition-colors">
											Cafes
										</Link>
										<Link href="/countries" className="text-snow/60 hover:text-gold transition-colors">
											Origins
										</Link>
										<Link href="/guide" className="text-snow/60 hover:text-gold transition-colors">
											Guide
										</Link>
										<Link href="/find" className="text-snow/60 hover:text-gold transition-colors">
											Find
										</Link>
										<Link href="/my-coffees" className="text-snow/60 hover:text-gold transition-colors">
											My Coffees
										</Link>
										<UserNav />
									</div>
									<div className="sm:hidden">
										<MobileNav />
									</div>
								</div>
							</nav>
							<main className="mx-auto flex-1 max-w-6xl px-4 py-8">{children}</main>
							<footer className="border-t border-border/40">
								<div className="mx-auto max-w-6xl px-4 py-6 text-center text-xs tracking-widest uppercase text-grey-olive">
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
