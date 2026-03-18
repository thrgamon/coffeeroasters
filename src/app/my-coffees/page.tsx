'use client';

import { useQueries } from '@tanstack/react-query';
import CoffeeCard from '@/components/CoffeeCard';
import { Skeleton } from '@/components/ui/skeleton';
import { getGetApiCoffeesIdQueryOptions } from '@/lib/api/generated/coffees/coffees';
import type { DomainCoffeeResponse } from '@/lib/api/generated/models';
import { useCoffeeTracker } from '@/lib/coffee-tracker';

function CoffeeGrid({ title, ids }: { title: string; ids: number[] }) {
	const queries = useQueries({
		queries: ids.map((id) => getGetApiCoffeesIdQueryOptions(id)),
	});

	const loading = queries.some((q) => q.isLoading);
	const coffees = queries
		.map((q) => q.data as DomainCoffeeResponse | undefined)
		.filter((c): c is DomainCoffeeResponse => c != null);

	return (
		<section className="space-y-3">
			<h2 className="text-xl font-semibold">{title}</h2>
			{loading ? (
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{ids.map((id) => (
						<Skeleton key={id} className="h-48 rounded-lg" />
					))}
				</div>
			) : coffees.length === 0 ? (
				<p className="text-sm text-muted-foreground">None yet.</p>
			) : (
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{coffees.map((coffee) => (
						<CoffeeCard key={coffee.id} coffee={coffee} />
					))}
				</div>
			)}
		</section>
	);
}

export default function MyCoffeesPage() {
	const { likedIds, triedIds } = useCoffeeTracker();

	const empty = likedIds.length === 0 && triedIds.length === 0;

	return (
		<div className="space-y-8">
			<h1 className="text-3xl font-bold">My coffees</h1>

			{empty ? (
				<p className="text-muted-foreground">
					You haven't liked or tried any coffees yet. Use the heart icon on coffee cards to start tracking.
				</p>
			) : (
				<>
					{likedIds.length > 0 && <CoffeeGrid title="Liked" ids={likedIds} />}
					{triedIds.length > 0 && <CoffeeGrid title="Tried" ids={triedIds} />}
				</>
			)}
		</div>
	);
}
