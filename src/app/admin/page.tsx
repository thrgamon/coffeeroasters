'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { AdminNav } from '@/components/AdminNav';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { adminFetch } from '@/lib/admin-api';
import { useAuth } from '@/lib/auth-context';

interface ScrapeRun {
	id: number;
	roaster_name: string;
	roaster_slug: string;
	started_at: string;
	finished_at: string;
	status: string;
	coffees_found: number;
	coffees_added: number;
	coffees_updated: number;
	coffees_removed: number;
	error_message: string;
	duration_ms: number;
}

export default function AdminPage() {
	const { user, loading } = useAuth();
	const router = useRouter();
	const [roasterCount, setRoasterCount] = useState<number | null>(null);
	const [coffeeCount, setCoffeeCount] = useState<number | null>(null);
	const [cafeCount, setCafeCount] = useState<number | null>(null);
	const [scrapeRuns, setScrapeRuns] = useState<ScrapeRun[]>([]);

	useEffect(() => {
		if (!loading && (!user || !user.is_admin)) router.push('/');
	}, [loading, user, router]);

	useEffect(() => {
		if (!user?.is_admin) return;
		adminFetch<{ id: number }[]>('/api/admin/roasters').then((data) => setRoasterCount(data.length));
		adminFetch<{ total: number; coffees: unknown[] }>('/api/admin/coffees?page=1&page_size=1').then((data) =>
			setCoffeeCount(data.total),
		);
		adminFetch<{ id: number }[]>('/api/admin/cafes').then((data) => setCafeCount(data.length));
		adminFetch<ScrapeRun[]>('/api/admin/scrape-runs').then((data) => setScrapeRuns(data.slice(0, 10)));
	}, [user]);

	if (loading || !user || !user.is_admin) return null;

	return (
		<div className="space-y-8">
			<AdminNav />
			<h1 className="font-display text-4xl">Admin Dashboard</h1>

			<div className="grid gap-4 sm:grid-cols-3">
				<Link href="/admin/roasters">
					<Card className="transition-shadow hover:shadow-md">
						<CardHeader>
							<CardTitle className="text-sm uppercase tracking-widest text-muted-foreground">Roasters</CardTitle>
						</CardHeader>
						<CardContent>
							<span className="font-mono text-3xl">{roasterCount ?? '...'}</span>
						</CardContent>
					</Card>
				</Link>
				<Link href="/admin/coffees">
					<Card className="transition-shadow hover:shadow-md">
						<CardHeader>
							<CardTitle className="text-sm uppercase tracking-widest text-muted-foreground">Coffees</CardTitle>
						</CardHeader>
						<CardContent>
							<span className="font-mono text-3xl">{coffeeCount ?? '...'}</span>
						</CardContent>
					</Card>
				</Link>
				<Link href="/admin/cafes">
					<Card className="transition-shadow hover:shadow-md">
						<CardHeader>
							<CardTitle className="text-sm uppercase tracking-widest text-muted-foreground">Cafes</CardTitle>
						</CardHeader>
						<CardContent>
							<span className="font-mono text-3xl">{cafeCount ?? '...'}</span>
						</CardContent>
					</Card>
				</Link>
			</div>

			<div>
				<h2 className="font-display text-2xl mb-4">Recent Scrape Runs</h2>
				{scrapeRuns.length === 0 ? (
					<p className="text-muted-foreground text-sm">No scrape runs found.</p>
				) : (
					<div className="overflow-x-auto">
						<table className="w-full text-sm">
							<thead>
								<tr className="border-b text-left text-xs uppercase tracking-widest text-muted-foreground">
									<th className="pb-2 pr-4">Roaster</th>
									<th className="pb-2 pr-4">Status</th>
									<th className="pb-2 pr-4 text-right">Found</th>
									<th className="pb-2 pr-4 text-right">Added</th>
									<th className="pb-2 pr-4 text-right">Updated</th>
									<th className="pb-2 pr-4 text-right">Removed</th>
									<th className="pb-2 pr-4 text-right">Duration</th>
									<th className="pb-2">Started</th>
								</tr>
							</thead>
							<tbody>
								{scrapeRuns.map((run) => (
									<tr key={run.id} className="border-b border-border/50">
										<td className="py-2 pr-4 font-medium">{run.roaster_name}</td>
										<td className="py-2 pr-4">
											<span
												className={
													run.status === 'completed'
														? 'text-green-700'
														: run.status === 'failed'
															? 'text-destructive'
															: 'text-muted-foreground'
												}
											>
												{run.status}
											</span>
										</td>
										<td className="py-2 pr-4 text-right font-mono">{run.coffees_found}</td>
										<td className="py-2 pr-4 text-right font-mono">{run.coffees_added}</td>
										<td className="py-2 pr-4 text-right font-mono">{run.coffees_updated}</td>
										<td className="py-2 pr-4 text-right font-mono">{run.coffees_removed}</td>
										<td className="py-2 pr-4 text-right font-mono">
											{run.duration_ms ? `${(run.duration_ms / 1000).toFixed(1)}s` : '-'}
										</td>
										<td className="py-2 text-muted-foreground">
											{run.started_at
												? new Date(run.started_at).toLocaleDateString('en-AU', {
														day: 'numeric',
														month: 'short',
														hour: '2-digit',
														minute: '2-digit',
													})
												: '-'}
										</td>
									</tr>
								))}
							</tbody>
						</table>
					</div>
				)}
			</div>
		</div>
	);
}
