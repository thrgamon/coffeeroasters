'use client';

import { useEffect } from 'react';

export default function ErrorPage({ error, reset }: { error: Error & { digest?: string }; reset: () => void }) {
	useEffect(() => {
		console.error(error);
	}, [error]);

	return (
		<div className="flex flex-col items-center justify-center gap-4 py-20">
			<h2 className="text-xl font-bold uppercase tracking-wider text-snow">Something went wrong</h2>
			<p className="text-grey-olive">An unexpected error occurred. Please try again.</p>
			<button
				type="button"
				onClick={() => reset()}
				className="rounded bg-gold px-4 py-2 text-sm font-bold uppercase tracking-wider text-rich-mahogany hover:bg-gold/90 transition-colors"
			>
				Try again
			</button>
		</div>
	);
}
