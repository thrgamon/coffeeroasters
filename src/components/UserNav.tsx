'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/auth-context';

export function UserNav() {
	const { user, loading, logout } = useAuth();
	const router = useRouter();

	if (loading) return null;

	if (user) {
		return (
			<div className="flex items-center gap-3">
				<span className="text-sm text-muted-foreground hidden lg:inline">{user.email}</span>
				<button
					type="button"
					onClick={async () => {
						await logout();
						router.push('/');
					}}
					className="text-sm text-muted-foreground hover:text-foreground"
				>
					Sign out
				</button>
			</div>
		);
	}

	return (
		<Link href="/login" className="text-sm text-muted-foreground hover:text-foreground">
			Sign in
		</Link>
	);
}
