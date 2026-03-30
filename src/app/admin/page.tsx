'use client';

import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import { useAuth } from '@/lib/auth-context';

export default function AdminPage() {
	const { user, loading } = useAuth();
	const router = useRouter();

	useEffect(() => {
		if (!loading && (!user || !user.is_admin)) {
			router.push('/');
		}
	}, [loading, user, router]);

	if (loading || !user || !user.is_admin) return null;

	return (
		<div className="space-y-6">
			<h1 className="text-3xl font-bold">Admin</h1>
			<p className="text-muted-foreground">Welcome, {user.email}. Admin features coming soon.</p>
		</div>
	);
}
