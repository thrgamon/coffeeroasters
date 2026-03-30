import type { Metadata } from 'next';
import Link from 'next/link';
import { Suspense } from 'react';
import RoasterFilters from '@/components/RoasterFilters';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { DomainRoasterListResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export const metadata: Metadata = { title: 'Roasters | COFFEEROASTERS' };

interface RoastersPageProps {
	searchParams: Promise<Record<string, string | string[] | undefined>>;
}

export default async function RoastersPage({ searchParams }: RoastersPageProps) {
	const sp = await searchParams;
	const state = typeof sp.state === 'string' ? sp.state : undefined;

	const params = new URLSearchParams();
	if (state) params.set('state', state);
	const qs = params.toString();

	const data = await apiFetch<DomainRoasterListResponse>(`/api/roasters${qs ? `?${qs}` : ''}`);

	return (
		<div className="space-y-6">
			<h1 className="text-3xl font-bold uppercase tracking-wider text-foreground">Roasters</h1>

			<Suspense>
				<RoasterFilters />
			</Suspense>

			{data.roasters && data.roasters.length > 0 ? (
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{data.roasters.map((roaster) => (
						<Link key={roaster.id} href={`/roasters/${roaster.slug}`}>
							<Card className="border-2 border-border bg-card transition-all hover:border-accent">
								<CardHeader>
									<CardTitle className="text-lg text-foreground">{roaster.name}</CardTitle>
								</CardHeader>
								<CardContent className="space-y-2">
									<div className="flex items-center gap-2">
										{roaster.state && <Badge variant="secondary">{roaster.state}</Badge>}
										{roaster.website && (
											<span className="text-sm text-muted-foreground">{new URL(roaster.website).hostname}</span>
										)}
									</div>
									{roaster.coffee_count ? (
										<p className="text-sm text-accent font-mono">
											{roaster.coffee_count} coffee{roaster.coffee_count !== 1 ? 's' : ''} in stock
										</p>
									) : null}
								</CardContent>
							</Card>
						</Link>
					))}
				</div>
			) : (
				<p className="text-muted-foreground">No roasters found.</p>
			)}
		</div>
	);
}
