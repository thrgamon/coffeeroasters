import type { Metadata } from 'next';
import Link from 'next/link';
import { Suspense } from 'react';
import CafeFilters from '@/components/CafeFilters';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { DomainCafeListResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export const metadata: Metadata = { title: 'Cafes | Coffeeroasters' };

interface CafesPageProps {
	searchParams: Promise<Record<string, string | string[] | undefined>>;
}

export default async function CafesPage({ searchParams }: CafesPageProps) {
	const sp = await searchParams;
	const state = typeof sp.state === 'string' ? sp.state : undefined;

	const params = new URLSearchParams();
	if (state) params.set('state', state);
	const qs = params.toString();

	const data = await apiFetch<DomainCafeListResponse>(`/api/cafes${qs ? `?${qs}` : ''}`);

	// Group cafes by state
	const cafesByState = new Map<string, typeof data.cafes>();
	for (const cafe of data.cafes ?? []) {
		const s = cafe.state ?? 'Other';
		if (!cafesByState.has(s)) cafesByState.set(s, []);
		cafesByState.get(s)!.push(cafe);
	}

	return (
		<div className="space-y-6">
			<h1 className="text-3xl font-bold">Cafes</h1>

			<Suspense>
				<CafeFilters />
			</Suspense>

			{data.cafes && data.cafes.length > 0 ? (
				state ? (
					<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
						{data.cafes.map((cafe) => (
							<CafeCard key={cafe.id} cafe={cafe} />
						))}
					</div>
				) : (
					Array.from(cafesByState.entries()).map(([stateKey, cafes]) => (
						<section key={stateKey} className="space-y-3">
							<h2 className="text-xl font-semibold">{stateKey}</h2>
							<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
								{cafes!.map((cafe) => (
									<CafeCard key={cafe.id} cafe={cafe} />
								))}
							</div>
						</section>
					))
				)
			) : (
				<p className="text-muted-foreground">No cafes found.</p>
			)}
		</div>
	);
}

function CafeCard({ cafe }: { cafe: NonNullable<DomainCafeListResponse['cafes']>[number] }) {
	return (
		<Card className="shadow-sm transition-all hover:shadow-md hover:bg-muted/50">
			<CardHeader>
				<CardTitle className="text-lg">{cafe.name}</CardTitle>
			</CardHeader>
			<CardContent className="space-y-2">
				<div className="flex items-center gap-2">
					{cafe.state && <Badge variant="secondary">{cafe.state}</Badge>}
					{cafe.suburb && <span className="text-sm text-muted-foreground">{cafe.suburb}</span>}
				</div>
				{cafe.address && <p className="text-sm text-muted-foreground">{cafe.address}</p>}
				{cafe.roaster_name && (
					<p className="text-sm">
						<Link href={`/roasters/${cafe.roaster_slug}`} className="text-primary hover:underline">
							{cafe.roaster_name}
						</Link>
					</p>
				)}
			</CardContent>
		</Card>
	);
}
