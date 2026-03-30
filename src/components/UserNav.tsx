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
				{user.is_admin && (
					<Link href="/admin" className="text-paper/70 hover:text-gold transition-colors">
						Admin
					</Link>
				)}
				<span className="text-paper/50 hidden lg:inline">{user.email}</span>
				<button
					type="button"
					onClick={async () => {
						await logout();
						router.push('/');
					}}
					className="text-paper/70 hover:text-gold transition-colors"
				>
					Sign out
				</button>
			</div>
		);
	}

	return (
		<Link href="/login" className="text-paper/70 hover:text-gold transition-colors">
			Sign in
		</Link>
	);
}
