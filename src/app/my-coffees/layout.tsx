import type { Metadata } from 'next';

export const metadata: Metadata = { title: 'My Coffees | Coffeeroasters' };

export default function Layout({ children }: { children: React.ReactNode }) {
	return children;
}
