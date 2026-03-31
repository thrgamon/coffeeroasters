'use client';

import { Menu } from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { useAuth } from '@/lib/auth-context';

export function UserNav() {
	const { user, loading, logout } = useAuth();
	const router = useRouter();
	const [open, setOpen] = useState(false);

	if (loading) return null;

	if (!user) {
		return (
			<Link href="/login" className="text-paper/70 hover:text-gold transition-colors">
				Sign in
			</Link>
		);
	}

	return (
		<div className="relative">
			<button
				type="button"
				onClick={() => setOpen(!open)}
				className="text-paper/70 hover:text-gold transition-colors"
				aria-label="Account menu"
			>
				<Menu className="size-5" />
			</button>
			{open && (
				<>
					<button
						type="button"
						className="fixed inset-0 z-40 cursor-default"
						onClick={() => setOpen(false)}
						aria-label="Close menu"
					/>
					<div className="absolute right-0 top-full z-50 mt-2 w-48 border border-soft-khaki bg-paper shadow-lg">
						<div className="px-4 py-3 border-b border-soft-khaki">
							<p className="text-xs text-warm-grey truncate">{user.email}</p>
						</div>
						<div className="py-1">
							{user.is_admin && (
								<Link
									href="/admin"
									onClick={() => setOpen(false)}
									className="block px-4 py-2 text-sm text-deep-indigo hover:bg-soft-khaki transition-colors"
								>
									Admin
								</Link>
							)}
							<button
								type="button"
								onClick={async () => {
									setOpen(false);
									await logout();
									router.push('/');
								}}
								className="block w-full px-4 py-2 text-left text-sm text-deep-indigo hover:bg-soft-khaki transition-colors"
							>
								Sign out
							</button>
						</div>
					</div>
				</>
			)}
		</div>
	);
}
