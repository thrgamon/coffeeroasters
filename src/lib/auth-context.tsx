'use client';

import { createContext, useCallback, useContext, useEffect, useState } from 'react';
import type { User } from '@/lib/types';

interface AuthContextType {
	user: User | null;
	loading: boolean;
	sendMagicLink: (email: string) => Promise<{ token?: string }>;
	verifyMagicLink: (token: string) => Promise<void>;
	logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

function extractError(data: unknown, fallback: string): string {
	if (
		data &&
		typeof data === 'object' &&
		'error' in data &&
		typeof (data as Record<string, unknown>).error === 'string'
	) {
		return (data as Record<string, string>).error;
	}
	return fallback;
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
	const [user, setUser] = useState<User | null>(null);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		fetch('/api/auth/me', { credentials: 'include' })
			.then(async (res) => {
				if (res.ok) {
					const data = await res.json();
					setUser(data.user ?? null);
				}
			})
			.catch(() => {})
			.finally(() => setLoading(false));
	}, []);

	const sendMagicLink = useCallback(async (email: string) => {
		const res = await fetch('/api/auth/magic-link', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ email }),
		});
		const data = await res.json();
		if (!res.ok) throw new Error(extractError(data, 'Failed to send magic link'));
		return { token: data.token as string | undefined };
	}, []);

	const verifyMagicLink = useCallback(async (token: string) => {
		const res = await fetch('/api/auth/verify', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ token }),
		});
		const data = await res.json();
		if (!res.ok) throw new Error(extractError(data, 'Verification failed'));
		setUser(data.user ?? data);
	}, []);

	const logout = useCallback(async () => {
		await fetch('/api/auth/logout', { method: 'POST', credentials: 'include' });
		setUser(null);
	}, []);

	return (
		<AuthContext value={{ user, loading, sendMagicLink, verifyMagicLink, logout }}>{children}</AuthContext>
	);
}

export function useAuth(): AuthContextType {
	const ctx = useContext(AuthContext);
	if (!ctx) throw new Error('useAuth must be used within AuthProvider');
	return ctx;
}
