import Link from 'next/link';
import { Suspense } from 'react';
import RoasterFilters from '@/components/RoasterFilters';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { DomainRoasterListResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

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
			<h1 className="text-3xl font-bold">Roasters</h1>

			<Suspense>
				<RoasterFilters />
			</Suspense>

			{data.roasters && data.roasters.length > 0 ? (
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{data.roasters.map((roaster) => (
						<Link key={roaster.id} href={`/roasters/${roaster.slug}`}>
							<Card className="shadow-sm transition-all hover:shadow-md hover:bg-muted/50">
								<CardHeader>
									<CardTitle className="text-lg">{roaster.name}</CardTitle>
								</CardHeader>
								<CardContent className="flex items-center gap-2">
									{roaster.state && <Badge variant="secondary">{roaster.state}</Badge>}
									{roaster.website && (
										<span className="text-sm text-muted-foreground">{new URL(roaster.website).hostname}</span>
									)}
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
