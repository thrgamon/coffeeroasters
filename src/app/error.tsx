'use client';

import { useEffect } from 'react';

export default function ErrorPage({ error, reset }: { error: Error & { digest?: string }; reset: () => void }) {
	useEffect(() => {
		console.error(error);
	}, [error]);

	return (
		<div className="flex flex-col items-center justify-center gap-4 py-20">
			<h2 className="text-xl font-bold tracking-[0.15em] text-foreground">Something Went Wrong</h2>
			<p className="text-muted-foreground">An unexpected error occurred. Please try again.</p>
			<button
				type="button"
				onClick={() => reset()}
				className="bg-foreground px-6 py-2.5 text-sm font-bold uppercase tracking-[0.15em] text-primary hover:bg-foreground/80 transition-colors"
			>
				Try again
			</button>
		</div>
	);
}
