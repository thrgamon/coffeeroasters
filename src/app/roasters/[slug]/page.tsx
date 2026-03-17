'use client';

import Link from 'next/link';
import { use } from 'react';
import { Badge } from '@/components/ui/badge';
import { useGetApiRoastersSlug } from '@/lib/api/generated/roasters/roasters';
import CoffeeCard from '@/components/CoffeeCard';

export default function RoasterDetailPage({ params }: { params: Promise<{ slug: string }> }) {
	const { slug } = use(params);
	const { data, isLoading, error } = useGetApiRoastersSlug(slug);

	if (isLoading) return <p className="text-muted-foreground">Loading...</p>;
	if (error || !data) return <p className="text-destructive">Roaster not found.</p>;

	const roaster = data?.roaster;
	const coffees = data?.coffees ?? [];

	return (
		<div className="space-y-8">
			<div className="space-y-2">
				<Link href="/roasters" className="text-sm text-muted-foreground hover:text-foreground">
					&larr; All roasters
				</Link>
				<h1 className="text-3xl font-bold">{roaster?.name}</h1>
				<div className="flex items-center gap-3">
					{roaster?.state && <Badge variant="secondary">{roaster.state}</Badge>}
					{roaster?.website && (
						<a
							href={roaster.website}
							target="_blank"
							rel="noopener noreferrer"
							className="text-sm text-primary hover:underline"
						>
							{roaster.website}
						</a>
					)}
				</div>
			</div>

			<section className="space-y-4">
				<h2 className="text-xl font-semibold">
					{coffees.length} coffee{coffees.length !== 1 ? 's' : ''}
				</h2>
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{coffees.map((coffee) => (
						<CoffeeCard key={coffee.id} coffee={coffee} />
					))}
				</div>
			</section>
		</div>
	);
}
