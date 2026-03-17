'use client';

import Link from 'next/link';
import { use } from 'react';
import { Badge } from '@/components/ui/badge';
import { useGetApiProducersId } from '@/lib/api/generated/producers/producers';
import CoffeeCard from '@/components/CoffeeCard';

export default function ProducerDetailPage({ params }: { params: Promise<{ id: string }> }) {
	const { id } = use(params);
	const { data, isLoading, error } = useGetApiProducersId(Number(id));

	if (isLoading) return <p className="text-muted-foreground">Loading...</p>;
	if (error || !data) return <p className="text-destructive">Producer not found.</p>;

	return (
		<div className="space-y-8">
			<div className="space-y-2">
				{data.country_code && (
					<Link
						href={`/countries/${data.country_code}`}
						className="text-sm text-muted-foreground hover:text-foreground"
					>
						&larr; {data.country_name}
					</Link>
				)}
				<h1 className="text-3xl font-bold">{data.name}</h1>
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
