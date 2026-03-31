'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { AdminNav } from '@/components/AdminNav';
import { Button } from '@/components/ui/button';
import { adminFetch } from '@/lib/admin-api';
import { useAuth } from '@/lib/auth-context';

interface Cafe {
	id: number;
	name: string;
	roaster_name: string;
	suburb: string;
	state: string;
	type: string;
	active: boolean;
}

export default function CafesListPage() {
	const { user, loading } = useAuth();
	const router = useRouter();
	const [cafes, setCafes] = useState<Cafe[]>([]);
	const [fetching, setFetching] = useState(true);

	useEffect(() => {
		if (!loading && (!user || !user.is_admin)) router.push('/');
	}, [loading, user, router]);

	useEffect(() => {
		if (!user?.is_admin) return;
		adminFetch<Cafe[]>('/api/admin/cafes')
			.then(setCafes)
			.finally(() => setFetching(false));
	}, [user]);

	if (loading || !user || !user.is_admin) return null;

	return (
		<div className="space-y-6">
			<AdminNav />
			<div className="flex items-center justify-between">
				<h1 className="font-display text-4xl">Cafes</h1>
				<Link href="/admin/cafes/new">
					<Button>New Cafe</Button>
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
								<th className="pb-2 pr-4">Roaster</th>
								<th className="pb-2 pr-4">Suburb</th>
								<th className="pb-2 pr-4">State</th>
								<th className="pb-2 pr-4">Type</th>
								<th className="pb-2">Active</th>
							</tr>
						</thead>
						<tbody>
							{cafes.map((c) => (
								<tr
									key={c.id}
									className="cursor-pointer border-b border-border/50 hover:bg-card"
									onClick={() => router.push(`/admin/cafes/${c.id}`)}
								>
									<td className="py-2 pr-4 font-medium">{c.name}</td>
									<td className="py-2 pr-4 text-muted-foreground">{c.roaster_name || '-'}</td>
									<td className="py-2 pr-4">{c.suburb || '-'}</td>
									<td className="py-2 pr-4">{c.state || '-'}</td>
									<td className="py-2 pr-4">{c.type || '-'}</td>
									<td className="py-2">
										<span
											className={`inline-block h-2 w-2 rounded-full ${c.active ? 'bg-green-600' : 'bg-muted-foreground/40'}`}
										/>
									</td>
								</tr>
							))}
						</tbody>
					</table>
				</div>
			)}
		</div>
	);
}
