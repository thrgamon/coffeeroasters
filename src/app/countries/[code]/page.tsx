'use client';

import Link from 'next/link';
import { use } from 'react';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useGetApiCountriesCode } from '@/lib/api/generated/countries/countries';
import CoffeeCard from '@/components/CoffeeCard';

export default function CountryDetailPage({ params }: { params: Promise<{ code: string }> }) {
	const { code } = use(params);
	const { data, isLoading, error } = useGetApiCountriesCode(code);

	if (isLoading) return <p className="text-muted-foreground">Loading...</p>;
	if (error || !data) return <p className="text-destructive">Country not found.</p>;

	return (
		<div className="space-y-8">
			<div className="space-y-2">
				<Link href="/countries" className="text-sm text-muted-foreground hover:text-foreground">
					&larr; All countries
				</Link>
				<h1 className="text-3xl font-bold">{data.name}</h1>
				<Badge variant="secondary">{data.code}</Badge>
			</div>

			{data.regions && data.regions.length > 0 && (
				<section className="space-y-4">
					<h2 className="text-xl font-semibold">Regions</h2>
					<div className="flex flex-wrap gap-2">
						{data.regions.map((region) => (
							<Link key={region.id} href={`/regions/${region.id}`}>
								<Badge variant="outline" className="cursor-pointer hover:bg-accent">
									{region.name} ({region.coffee_count})
								</Badge>
							</Link>
						))}
					</div>
				</section>
			)}

			<section className="space-y-4">
				<h2 className="text-xl font-semibold">
					{data.coffees?.length ?? 0} coffee{(data.coffees?.length ?? 0) !== 1 ? 's' : ''}
				</h2>
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{data.coffees?.map((coffee) => (
						<CoffeeCard key={coffee.id} coffee={coffee} />
					))}
				</div>
			</section>
		</div>
	);
}
