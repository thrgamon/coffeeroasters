'use client';

import CoffeeCard from '@/components/CoffeeCard';
import { Skeleton } from '@/components/ui/skeleton';
import { useGetApiCoffees } from '@/lib/api/generated/coffees/coffees';
import { useCoffeeTracker } from '@/lib/coffee-tracker';

export default function Recommendations() {
	const { hydrated, wishlistIds } = useCoffeeTracker();

	const { data, isLoading } = useGetApiCoffees(
		{ liked: wishlistIds.join(','), page_size: 6 },
		{ query: { enabled: hydrated && wishlistIds.length > 0 } },
	);

	if (!hydrated || wishlistIds.length === 0) {
		return null;
	}

	if (isLoading) {
		return (
			<section className="space-y-4">
				<h2 className="text-2xl font-semibold">Recommended for you</h2>
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{['rec-skel-1', 'rec-skel-2', 'rec-skel-3', 'rec-skel-4', 'rec-skel-5', 'rec-skel-6'].map((id) => (
						<Skeleton key={id} className="h-48 rounded-lg" />
					))}
				</div>
			</section>
		);
	}

	if (!data?.coffees || data.coffees.length === 0) {
		return null;
	}

	return (
		<section className="space-y-4">
			<h2 className="text-2xl font-semibold">Recommended for you</h2>
			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{data.coffees.map((coffee) => (
					<CoffeeCard key={coffee.id} coffee={coffee} />
				))}
			</div>
		</section>
	);
}
