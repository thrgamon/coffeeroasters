'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { Suspense, useEffect, useState } from 'react';
import { ErrorBanner } from '@/components/ErrorBanner';
import { useAuth } from '@/lib/auth-context';

function LoginForm() {
	const [email, setEmail] = useState('');
	const [error, setError] = useState('');
	const [sent, setSent] = useState(false);
	const [token, setToken] = useState('');
	const { sendMagicLink, verifyMagicLink, user } = useAuth();
	const router = useRouter();
	const searchParams = useSearchParams();

	// If user is already logged in, redirect
	useEffect(() => {
		if (user) {
			router.push(searchParams.get('redirect') ?? '/my-coffees');
		}
	}, [user, router, searchParams]);

	// Check for token in URL (magic link click)
	useEffect(() => {
		const urlToken = searchParams.get('token');
		if (urlToken) {
			verifyMagicLink(urlToken)
				.then(() => router.push(searchParams.get('redirect') ?? '/my-coffees'))
				.catch((err) => setError(err instanceof Error ? err.message : 'Verification failed'));
		}
	}, [searchParams, verifyMagicLink, router]);

	async function handleSendLink(e: React.FormEvent) {
		e.preventDefault();
		setError('');
		try {
			const result = await sendMagicLink(email);
			setSent(true);
			// In development, auto-fill the token for easy testing
			if (result.token) {
				setToken(result.token);
			}
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Failed to send link');
		}
	}

	async function handleVerify(e: React.FormEvent) {
		e.preventDefault();
		setError('');
		try {
			await verifyMagicLink(token);
			router.push(searchParams.get('redirect') ?? '/my-coffees');
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Verification failed');
		}
	}

	if (sent) {
		return (
			<main className="flex min-h-screen items-center justify-center">
				<form onSubmit={handleVerify} className="w-full max-w-sm space-y-4">
					<h1 className="text-2xl font-bold">Check your email</h1>
					<p className="text-sm text-muted-foreground">
						We sent a login link to <strong>{email}</strong>. Click the link in your email, or paste the
						code below.
					</p>
					{error && <ErrorBanner message={error} />}
					<input
						type="text"
						placeholder="Paste your login code"
						value={token}
						onChange={(e) => setToken(e.target.value)}
						required
						className="w-full rounded border border-input bg-background px-3 py-2 font-mono text-sm"
					/>
					<button type="submit" className="w-full rounded bg-primary px-4 py-2 text-primary-foreground">
						Verify
					</button>
					<button
						type="button"
						onClick={() => {
							setSent(false);
							setToken('');
						}}
						className="w-full text-sm text-muted-foreground underline"
					>
						Use a different email
					</button>
				</form>
			</main>
		);
	}

	return (
		<main className="flex min-h-screen items-center justify-center">
			<form onSubmit={handleSendLink} className="w-full max-w-sm space-y-4">
				<h1 className="text-2xl font-bold">Sign in</h1>
				<p className="text-sm text-muted-foreground">
					Enter your email and we'll send you a login link. No password needed.
				</p>
				{error && <ErrorBanner message={error} />}
				<input
					type="email"
					placeholder="Email"
					value={email}
					onChange={(e) => setEmail(e.target.value)}
					required
					className="w-full rounded border border-input bg-background px-3 py-2"
				/>
				<button type="submit" className="w-full rounded bg-primary px-4 py-2 text-primary-foreground">
					Send login link
				</button>
			</form>
		</main>
	);
}

export default function LoginPage() {
	return (
		<Suspense>
			<LoginForm />
		</Suspense>
	);
}
