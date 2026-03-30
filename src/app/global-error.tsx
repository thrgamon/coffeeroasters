'use client';

export default function GlobalError({ reset }: { error: Error & { digest?: string }; reset: () => void }) {
	return (
		<html lang="en">
			<body className="flex min-h-screen flex-col items-center justify-center bg-background font-sans antialiased">
				<div className="flex flex-col items-center gap-4">
					<h2 className="text-xl font-semibold">Something went wrong</h2>
					<p className="text-muted-foreground">An unexpected error occurred.</p>
					<button
						type="button"
						onClick={() => reset()}
						className="rounded-md bg-neutral-900 px-4 py-2 text-sm font-medium text-white hover:bg-neutral-800"
					>
						Try again
					</button>
				</div>
			</body>
		</html>
	);
}
