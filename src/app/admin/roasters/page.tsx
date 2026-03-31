'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { AdminNav } from '@/components/AdminNav';
import { Button } from '@/components/ui/button';
import { adminFetch } from '@/lib/admin-api';
import { useAuth } from '@/lib/auth-context';

interface Roaster {
	id: number;
	slug: string;
	name: string;
	website: string;
	state: string;
	active: boolean;
	opted_out: boolean;
}

export default function RoastersListPage() {
	const { user, loading } = useAuth();
	const router = useRouter();
	const [roasters, setRoasters] = useState<Roaster[]>([]);
	const [fetching, setFetching] = useState(true);

	useEffect(() => {
		if (!loading && (!user || !user.is_admin)) router.push('/');
	}, [loading, user, router]);

	useEffect(() => {
		if (!user?.is_admin) return;
		adminFetch<Roaster[]>('/api/admin/roasters')
			.then(setRoasters)
			.finally(() => setFetching(false));
	}, [user]);

	if (loading || !user || !user.is_admin) return null;

	return (
		<div className="space-y-6">
			<AdminNav />
			<div className="flex items-center justify-between">
				<h1 className="font-display text-4xl">Roasters</h1>
				<Link href="/admin/roasters/new">
					<Button>New Roaster</Button>
				</Link>
			</div>

			{fetching ? (
				<p className="text-muted-foreground text-sm">Loading...</p>
			) : (
				<div className="overflow-x-auto">
					<table className="w-full text-sm">
						<thead>
							<tr className="border-b text-left text-xs uppercase tracking-widest text-muted-foreground">
								<th className="pb-2 pr-4">Name</th>
								<th className="pb-2 pr-4">Slug</th>
								<th className="pb-2 pr-4">State</th>
								<th className="pb-2 pr-4">Active</th>
								<th className="pb-2 pr-4">Opted Out</th>
								<th className="pb-2">Website</th>
							</tr>
						</thead>
						<tbody>
							{roasters.map((r) => (
								<tr
									key={r.id}
									className="cursor-pointer border-b border-border/50 hover:bg-card"
									onClick={() => router.push(`/admin/roasters/${r.id}`)}
								>
									<td className="py-2 pr-4 font-medium">{r.name}</td>
									<td className="py-2 pr-4 font-mono text-muted-foreground">{r.slug}</td>
									<td className="py-2 pr-4">{r.state}</td>
									<td className="py-2 pr-4">
										<span
											className={`inline-block h-2 w-2 rounded-full ${r.active ? 'bg-green-600' : 'bg-muted-foreground/40'}`}
										/>
									</td>
									<td className="py-2 pr-4">{r.opted_out ? 'Yes' : '-'}</td>
									<td className="py-2 truncate max-w-48 text-muted-foreground">{r.website || '-'}</td>
								</tr>
							))}
						</tbody>
					</table>
				</div>
			)}
		</div>
	);
}
