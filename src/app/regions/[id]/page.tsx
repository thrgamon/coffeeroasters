import type { Metadata } from 'next';
import Link from 'next/link';
import CoffeeCard from '@/components/CoffeeCard';
import { Badge } from '@/components/ui/badge';
import type { DomainRegionDetailResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export async function generateMetadata({ params }: { params: Promise<{ id: string }> }): Promise<Metadata> {
	const { id } = await params;
	const data = await apiFetch<DomainRegionDetailResponse>(`/api/regions/${id}`);
	return { title: `${data.name} | COFFEEROASTERS` };
}

export default async function RegionDetailPage({ params }: { params: Promise<{ id: string }> }) {
	const { id } = await params;
	const data = await apiFetch<DomainRegionDetailResponse>(`/api/regions/${id}`);

	return (
		<div className="space-y-8">
			<div className="space-y-2">
				<Link href={`/countries/${data.country_code}`} className="text-sm text-grey-olive hover:text-gold transition-colors">
					&larr; {data.country_name}
				</Link>
				<h1 className="text-3xl font-bold text-snow">{data.name}</h1>
				<div className="flex gap-2">
					<Link href={`/countries/${data.country_code}`}>
						<Badge variant="secondary">{data.country_name}</Badge>
					</Link>
				</div>
			</div>

			<section className="space-y-4">
				<h2 className="text-xl font-bold uppercase tracking-wider text-dusty-rose">
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
					<h2 className="text-xl font-bold uppercase tracking-wider text-dusty-rose">Nearby Regions</h2>
					<div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
						{data.nearby_regions.map((region) =>
							region.coffee_count && region.coffee_count > 0 ? (
								<Link
									key={region.id}
									href={`/regions/${region.id}`}
									className="flex items-center justify-between rounded border border-border/50 bg-card/80 p-4 hover:border-gold/30 transition-all"
								>
									<div>
										<p className="font-medium text-snow">{region.name}</p>
										<p className="text-sm text-grey-olive">{region.country_name}</p>
										<p className="text-sm text-dusty-rose font-mono">
											{region.coffee_count} coffee{region.coffee_count !== 1 ? 's' : ''}
										</p>
									</div>
									<Badge variant="outline">{region.distance_km} km</Badge>
								</Link>
							) : (
								<div key={region.id} className="flex items-center justify-between rounded border border-border/50 bg-card/80 p-4">
									<div>
										<p className="font-medium text-grey-olive">{region.name}</p>
										<p className="text-sm text-grey-olive">{region.country_name}</p>
										<p className="text-sm text-grey-olive font-mono">0 coffees</p>
									</div>
									<Badge variant="outline">{region.distance_km} km</Badge>
								</div>
							),
						)}
					</div>
				</section>
			)}
		</div>
	);
}
