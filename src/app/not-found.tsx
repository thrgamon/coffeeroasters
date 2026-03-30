import Link from 'next/link';

export default function NotFound() {
	return (
		<div className="flex flex-col items-center justify-center gap-6 py-20">
			<p className="text-8xl font-bold text-primary font-mono">404</p>
			<h2 className="text-xl font-bold tracking-[0.15em] text-foreground">Page Not Found</h2>
			<p className="text-muted-foreground">The page you were looking for does not exist.</p>
			<Link
				href="/"
				className="bg-foreground px-8 py-3 text-sm font-bold uppercase tracking-[0.2em] text-primary hover:bg-foreground/80 transition-colors"
			>
				Go home
			</Link>
		</div>
	);
}
