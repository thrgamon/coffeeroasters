import Link from 'next/link';

export default function NotFound() {
	return (
		<div className="flex flex-col items-center justify-center gap-6 py-20">
			<p className="text-6xl font-bold text-gold font-mono glow-gold">404</p>
			<h2 className="text-xl font-bold uppercase tracking-wider text-snow">Page not found</h2>
			<p className="text-grey-olive">The page you were looking for does not exist.</p>
			<Link
				href="/"
				className="rounded bg-gold px-6 py-2.5 text-sm font-bold uppercase tracking-wider text-rich-mahogany hover:bg-gold/90 transition-colors"
			>
				Go home
			</Link>
		</div>
	);
}
