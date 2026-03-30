import type { Metadata } from 'next';

export const metadata: Metadata = {
	title: 'Find your coffee | COFFEEROASTERS',
	description: 'Answer a few questions about your flavour preferences and find coffees you will love',
};

export default function FindLayout({ children }: { children: React.ReactNode }) {
	return children;
}
