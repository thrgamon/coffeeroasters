'use client';

import Link from 'next/link';
import { useGetApiCoffees } from '@/lib/api/generated/coffees/coffees';
import { useGetApiStats } from '@/lib/api/generated/stats/stats';
import CoffeeCard from '@/components/CoffeeCard';

export default function Home() {
	const { data: stats } = useGetApiStats();
	const { data: coffees } = useGetApiCoffees({ page_size: 6 });

	return (
		<div className="space-y-12">
			<section className="space-y-4 text-center">
				<h1 className="text-4xl font-bold">Coffeeroasters</h1>
				<p className="text-lg text-muted-foreground">Discover specialty coffee from Australian indie roasters</p>
				{stats && (
					<div className="flex justify-center gap-8 text-sm text-muted-foreground">
						<span>{stats.roaster_count} roasters</span>
						<span>{stats.coffee_count} coffees</span>
						<span>{stats.origins?.length ?? 0} origins</span>
					</div>
				)}
				<div className="flex justify-center gap-4">
					<Link
						href="/coffees"
						className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
					>
						Browse Coffees
					</Link>
					<Link
						href="/roasters"
						className="rounded-md border border-border px-4 py-2 text-sm font-medium hover:bg-accent"
					>
						View Roasters
					</Link>
				</div>
			</section>

			{coffees?.coffees && coffees.coffees.length > 0 && (
				<section className="space-y-4">
					<h2 className="text-2xl font-semibold">Latest Coffees</h2>
					<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
						{coffees.coffees.map((coffee) => (
							<CoffeeCard key={coffee.id} coffee={coffee} />
						))}
					</div>
				</section>
			)}
		</div>
	);
}
