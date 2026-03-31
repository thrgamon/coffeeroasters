'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useCallback, useEffect, useState } from 'react';
import { AdminNav } from '@/components/AdminNav';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { adminFetch } from '@/lib/admin-api';
import { useAuth } from '@/lib/auth-context';

interface Coffee {
	id: number;
	name: string;
	roaster_name: string;
	country_code: string;
	process: string;
	in_stock: boolean;
	price_cents: number;
}

interface CoffeeListResponse {
	coffees: Coffee[];
	total: number;
	page: number;
	page_size: number;
}

const PAGE_SIZE = 50;

export default function CoffeesListPage() {
	const { user, loading } = useAuth();
	const router = useRouter();
	const [data, setData] = useState<CoffeeListResponse | null>(null);
	const [page, setPage] = useState(1);
	const [search, setSearch] = useState('');
	const [fetching, setFetching] = useState(true);

	useEffect(() => {
		if (!loading && (!user || !user.is_admin)) router.push('/');
	}, [loading, user, router]);

	const fetchCoffees = useCallback(
		(p: number) => {
			if (!user?.is_admin) return;
			setFetching(true);
			const params = new URLSearchParams({ page: String(p), page_size: String(PAGE_SIZE) });
			if (search) params.set('search', search);
			adminFetch<CoffeeListResponse>(`/api/admin/coffees?${params}`)
				.then(setData)
				.finally(() => setFetching(false));
		},
		[user, search],
	);

	useEffect(() => {
		fetchCoffees(page);
	}, [fetchCoffees, page]);

	if (loading || !user || !user.is_admin) return null;

	const totalPages = data ? Math.ceil(data.total / PAGE_SIZE) : 0;

	const handleSearch = () => {
		setPage(1);
		fetchCoffees(1);
	};

	return (
		<div className="space-y-6">
			<AdminNav />
			<div className="flex items-center justify-between">
				<h1 className="font-display text-4xl">Coffees</h1>
				<Link href="/admin/coffees/new">
					<Button>New Coffee</Button>
				</Link>
			</div>

			<div className="flex gap-2">
				<Input
					placeholder="Search by roaster name..."
					value={search}
					onChange={(e) => setSearch(e.target.value)}
					onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
					className="max-w-xs"
				/>
				<Button variant="outline" onClick={handleSearch}>
					Search
				</Button>
			</div>

			{fetching ? (
				<p className="text-muted-foreground text-sm">Loading...</p>
			) : !data || data.coffees.length === 0 ? (
				<p className="text-muted-foreground text-sm">No coffees found.</p>
			) : (
				<>
					<div className="overflow-x-auto">
						<table className="w-full text-sm">
							<thead>
								<tr className="border-b text-left text-xs uppercase tracking-widest text-muted-foreground">
									<th className="pb-2 pr-4">Name</th>
									<th className="pb-2 pr-4">Roaster</th>
									<th className="pb-2 pr-4">Origin</th>
									<th className="pb-2 pr-4">Process</th>
									<th className="pb-2 pr-4">In Stock</th>
									<th className="pb-2 text-right">Price</th>
								</tr>
							</thead>
							<tbody>
								{data.coffees.map((c) => (
									<tr
										key={c.id}
										className="cursor-pointer border-b border-border/50 hover:bg-card"
										onClick={() => router.push(`/admin/coffees/${c.id}`)}
									>
										<td className="py-2 pr-4 font-medium">{c.name}</td>
										<td className="py-2 pr-4 text-muted-foreground">{c.roaster_name}</td>
										<td className="py-2 pr-4">{c.country_code || '-'}</td>
										<td className="py-2 pr-4">{c.process || '-'}</td>
										<td className="py-2 pr-4">
											<span
												className={`inline-block h-2 w-2 rounded-full ${c.in_stock ? 'bg-green-600' : 'bg-muted-foreground/40'}`}
											/>
										</td>
										<td className="py-2 text-right font-mono">
											{c.price_cents ? `$${(c.price_cents / 100).toFixed(2)}` : '-'}
										</td>
									</tr>
								))}
							</tbody>
						</table>
					</div>

					{totalPages > 1 && (
						<div className="flex items-center justify-between pt-2">
							<span className="text-sm text-muted-foreground">
								Page {page} of {totalPages} ({data.total} coffees)
							</span>
							<div className="flex gap-2">
								<Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage(page - 1)}>
									Previous
								</Button>
								<Button variant="outline" size="sm" disabled={page >= totalPages} onClick={() => setPage(page + 1)}>
									Next
								</Button>
							</div>
						</div>
					)}
				</>
			)}
		</div>
	);
}
