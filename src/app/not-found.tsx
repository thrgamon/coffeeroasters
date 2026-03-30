import Link from 'next/link';

export default function NotFound() {
	return (
		<div className="flex flex-col items-center justify-center gap-4 py-20">
			<h2 className="text-xl font-semibold">Page not found</h2>
			<p className="text-muted-foreground">The page you were looking for does not exist.</p>
			<Link
				href="/"
				className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
			>
				Go home
			</Link>
		</div>
	);
}
