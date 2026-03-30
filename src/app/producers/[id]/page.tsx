import Link from 'next/link';
import CoffeeCard from '@/components/CoffeeCard';
import { Badge } from '@/components/ui/badge';
import type { DomainProducerDetailResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export default async function ProducerDetailPage({ params }: { params: Promise<{ id: string }> }) {
	const { id } = await params;
	const data = await apiFetch<DomainProducerDetailResponse>(`/api/producers/${id}`);

	return (
		<div className="space-y-8">
			<div className="space-y-2">
				{data.country_code && (
					<Link
						href={`/countries/${data.country_code}`}
						className="text-sm text-muted-foreground hover:text-foreground transition-colors"
					>
						&larr; {data.country_name}
					</Link>
				)}
				<h1 className="text-3xl font-bold text-foreground">{data.name}</h1>
				<div className="flex gap-2">
					{data.country_name && (
						<Link href={`/countries/${data.country_code}`}>
							<Badge variant="secondary">{data.country_name}</Badge>
						</Link>
					)}
					{data.region_name && <Badge variant="outline">{data.region_name}</Badge>}
				</div>
			</div>

			<section className="space-y-4">
				<h2 className="text-xl font-bold uppercase tracking-wider text-accent">
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
