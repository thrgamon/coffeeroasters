import type { Metadata } from 'next';
import './globals.css';
import Link from 'next/link';
import { MobileNav } from '@/components/MobileNav';
import { AuthProvider } from '@/lib/auth-context';
import { CoffeeTrackerProvider } from '@/lib/coffee-tracker';
import { QueryProvider } from '@/lib/query-provider';

export const metadata: Metadata = {
	title: 'Coffeeroasters',
	description: 'Discover Australian specialty coffee roasters',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
	return (
		<html lang="en">
			<body className="flex min-h-screen flex-col bg-background font-sans antialiased">
				<QueryProvider>
					<AuthProvider>
						<CoffeeTrackerProvider>
							<nav className="border-b border-border shadow-sm">
								<div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-4">
									<Link href="/" className="text-lg font-semibold text-primary">
										Coffeeroasters
									</Link>
									<div className="hidden gap-6 text-sm text-muted-foreground sm:flex">
										<Link href="/coffees" className="hover:text-foreground">
											Coffees
										</Link>
										<Link href="/roasters" className="hover:text-foreground">
											Roasters
										</Link>
										<Link href="/cafes" className="hover:text-foreground">
											Cafes
										</Link>
										<Link href="/countries" className="hover:text-foreground">
											Countries
										</Link>
										<Link href="/my-coffees" className="hover:text-foreground">
											My coffees
										</Link>
									</div>
									<div className="sm:hidden">
										<MobileNav />
									</div>
								</div>
							</nav>
							<main className="mx-auto flex-1 max-w-6xl px-4 py-8">{children}</main>
							<footer className="border-t border-border">
								<div className="mx-auto max-w-6xl px-4 py-4 text-center text-sm text-muted-foreground">
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
