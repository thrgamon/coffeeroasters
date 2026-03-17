'use client';

import Link from 'next/link';
import { use } from 'react';
import CoffeeCard from '@/components/CoffeeCard';
import { Badge } from '@/components/ui/badge';
import { useGetApiRegionsId } from '@/lib/api/generated/regions/regions';

export default function RegionDetailPage({ params }: { params: Promise<{ id: string }> }) {
	const { id } = use(params);
	const { data, isLoading, error } = useGetApiRegionsId(Number(id));

	if (isLoading) return <p className="text-muted-foreground">Loading...</p>;
	if (error || !data) return <p className="text-destructive">Region not found.</p>;

	return (
		<div className="space-y-8">
			<div className="space-y-2">
				<Link href={`/countries/${data.country_code}`} className="text-sm text-muted-foreground hover:text-foreground">
					&larr; {data.country_name}
				</Link>
				<h1 className="text-3xl font-bold">{data.name}</h1>
				<div className="flex gap-2">
					<Link href={`/countries/${data.country_code}`}>
						<Badge variant="secondary">{data.country_name}</Badge>
					</Link>
				</div>
			</div>

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

			{data.nearby_regions && data.nearby_regions.length > 0 && (
				<section className="space-y-4">
					<h2 className="text-xl font-semibold">Nearby Regions</h2>
					<div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
						{data.nearby_regions.map((region) => (
							<Link
								key={region.id}
								href={`/regions/${region.id}`}
								className="flex items-center justify-between rounded-lg border p-4 hover:bg-muted/50 transition-colors"
							>
								<div>
									<p className="font-medium">{region.name}</p>
									<p className="text-sm text-muted-foreground">{region.country_name}</p>
									<p className="text-sm text-muted-foreground">
										{region.coffee_count} coffee{region.coffee_count !== 1 ? 's' : ''}
									</p>
								</div>
								<Badge variant="outline">{region.distance_km} km</Badge>
							</Link>
						))}
					</div>
				</section>
			)}
		</div>
	);
}
