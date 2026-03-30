import type { Metadata } from 'next';
import Link from 'next/link';
import CoffeeCard from '@/components/CoffeeCard';
import { Badge } from '@/components/ui/badge';
import type { DomainCountryDetailResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export async function generateMetadata({ params }: { params: Promise<{ code: string }> }): Promise<Metadata> {
	const { code } = await params;
	const data = await apiFetch<DomainCountryDetailResponse>(`/api/countries/${code}`);
	return { title: `${data.name} | COFFEEROASTERS` };
}

export default async function CountryDetailPage({ params }: { params: Promise<{ code: string }> }) {
	const { code } = await params;
	const data = await apiFetch<DomainCountryDetailResponse>(`/api/countries/${code}`);

	return (
		<div className="space-y-8">
			<div className="space-y-2">
				<Link href="/countries" className="text-sm text-grey-olive hover:text-gold transition-colors">
					&larr; All countries
				</Link>
				<h1 className="text-3xl font-bold text-snow">{data.name}</h1>
				<Badge variant="secondary">{data.code}</Badge>
			</div>

			{data.regions && data.regions.length > 0 && (
				<section className="space-y-4">
					<h2 className="text-xl font-bold uppercase tracking-wider text-dusty-rose">Regions</h2>
					<div className="flex flex-wrap gap-2">
						{data.regions.map((region) =>
							region.coffee_count && region.coffee_count > 0 ? (
								<Link key={region.id} href={`/regions/${region.id}`}>
									<Badge variant="outline" className="cursor-pointer hover:border-gold/50 hover:text-gold">
										{region.name} ({region.coffee_count})
									</Badge>
								</Link>
							) : (
								<Badge key={region.id} variant="outline" className="text-grey-olive">
									{region.name} (0)
								</Badge>
							),
						)}
					</div>
				</section>
			)}

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
		</div>
	);
}
