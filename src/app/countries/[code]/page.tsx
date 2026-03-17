import Link from 'next/link';
import CoffeeCard from '@/components/CoffeeCard';
import { Badge } from '@/components/ui/badge';
import type { DomainCountryDetailResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export default async function CountryDetailPage({ params }: { params: Promise<{ code: string }> }) {
	const { code } = await params;
	const data = await apiFetch<DomainCountryDetailResponse>(`/api/countries/${code}`);

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
