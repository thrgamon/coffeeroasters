import type { Metadata } from 'next';
import Link from 'next/link';
import { Suspense } from 'react';
import CafeFilters from '@/components/CafeFilters';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { DomainCafeListResponse, DomainCafeResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export const metadata: Metadata = { title: 'Cafes | COFFEEROASTERS' };

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

	const allCafes = data.cafes ?? [];
	const ownedCafes = allCafes.filter((c) => c.type !== 'stockist');
	const stockistCafes = allCafes.filter((c) => c.type === 'stockist');

	return (
		<div className="space-y-6">
			<h1 className="text-3xl font-bold uppercase tracking-wider text-foreground">Cafes</h1>

			<Suspense>
				<CafeFilters />
			</Suspense>

			{allCafes.length > 0 ? (
				<>
					{ownedCafes.length > 0 && (
						<CafeSection title="Roaster-owned cafes" cafes={ownedCafes} groupByState={!state} />
					)}
					{stockistCafes.length > 0 && (
						<CafeSection
							title="Stockists"
							description="Independent cafes serving specialty roasters' coffee"
							cafes={stockistCafes}
							groupByState={!state}
						/>
					)}
				</>
			) : (
				<p className="text-muted-foreground">No cafes found.</p>
			)}
		</div>
	);
}

function CafeSection({
	title,
	description,
	cafes,
	groupByState,
}: {
	title: string;
	description?: string;
	cafes: DomainCafeResponse[];
	groupByState: boolean;
}) {
	if (groupByState) {
		const byState = new Map<string, DomainCafeResponse[]>();
		for (const cafe of cafes) {
			const s = cafe.state ?? 'Other';
			if (!byState.has(s)) byState.set(s, []);
			byState.get(s)?.push(cafe);
		}

		return (
			<section className="space-y-4">
				<div>
					<h2 className="text-xl font-bold uppercase tracking-wider text-accent">{title}</h2>
					{description && <p className="text-sm text-muted-foreground">{description}</p>}
				</div>
				{Array.from(byState.entries()).map(([stateKey, stateCafes]) => (
					<div key={stateKey} className="space-y-3">
						<h3 className="text-lg font-medium text-foreground">{stateKey}</h3>
						<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
							{stateCafes.map((cafe) => (
								<CafeCard key={cafe.id} cafe={cafe} />
							))}
						</div>
					</div>
				))}
			</section>
		);
	}

	return (
		<section className="space-y-4">
			<div>
				<h2 className="text-xl font-bold uppercase tracking-wider text-accent">{title}</h2>
				{description && <p className="text-sm text-muted-foreground">{description}</p>}
			</div>
			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{cafes.map((cafe) => (
					<CafeCard key={cafe.id} cafe={cafe} />
				))}
			</div>
		</section>
	);
}

function CafeCard({ cafe }: { cafe: DomainCafeResponse }) {
	return (
		<Card className="border-2 border-border bg-card transition-all hover:border-accent">
			<CardHeader>
				<CardTitle className="text-lg text-foreground">{cafe.name}</CardTitle>
			</CardHeader>
			<CardContent className="space-y-2">
				<div className="flex items-center gap-2">
					{cafe.state && <Badge variant="secondary">{cafe.state}</Badge>}
					{cafe.type === 'stockist' && <Badge variant="outline">Stockist</Badge>}
					{cafe.suburb && <span className="text-sm text-muted-foreground">{cafe.suburb}</span>}
				</div>
				{cafe.address && <p className="text-sm text-muted-foreground">{cafe.address}</p>}
				{cafe.roaster_name && (
					<p className="text-sm">
						<Link href={`/roasters/${cafe.roaster_slug}`} className="text-accent hover:text-foreground">
							{cafe.roaster_name}
						</Link>
					</p>
				)}
			</CardContent>
		</Card>
	);
}
