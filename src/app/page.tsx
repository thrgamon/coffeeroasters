import Link from 'next/link';
import CoffeeCard from '@/components/CoffeeCard';
import Recommendations from '@/components/Recommendations';
import type { DomainCoffeeListResponse, DomainStatsResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export const dynamic = 'force-dynamic';

export default async function Home() {
	const [stats, coffees] = await Promise.all([
		apiFetch<DomainStatsResponse>('/api/stats'),
		apiFetch<DomainCoffeeListResponse>('/api/coffees?page_size=6'),
	]);

	return (
		<div className="space-y-16">
			<section className="text-center py-12">
				<p className="text-xs font-bold tracking-[0.4em] uppercase text-accent mb-6">Australian Specialty Coffee</p>
				<h1 className="text-5xl font-bold tracking-[0.08em] text-foreground sm:text-8xl leading-none">
					COFFEE
					<br />
					ROASTERS
				</h1>
				<div className="mx-auto mt-8 flex justify-center px-4">
					<div className="bg-primary px-6 py-4 sm:px-10">
						<div className="flex gap-6 sm:gap-12 text-sm font-mono font-bold tracking-wider text-foreground">
							<span className="flex flex-col items-center">
								<span className="text-3xl">{stats.roaster_count}</span>
								<span className="text-xs tracking-[0.2em]">ROASTERS</span>
							</span>
							<span className="flex flex-col items-center">
								<span className="text-3xl">{stats.coffee_count}</span>
								<span className="text-xs tracking-[0.2em]">COFFEES</span>
							</span>
							<span className="flex flex-col items-center">
								<span className="text-3xl">{stats.origins?.length ?? 0}</span>
								<span className="text-xs tracking-[0.2em]">ORIGINS</span>
							</span>
						</div>
					</div>
				</div>
				<div className="flex flex-col sm:flex-row justify-center gap-4 mt-10 px-4">
					<Link
						href="/coffees"
						className="bg-primary px-8 py-3 text-sm font-bold uppercase tracking-[0.2em] text-primary-foreground hover:bg-primary/80 transition-colors"
					>
						Browse Coffees
					</Link>
					<Link
						href="/roasters"
						className="border-2 border-foreground px-8 py-3 text-sm font-bold uppercase tracking-[0.2em] text-foreground hover:bg-foreground hover:text-paper transition-colors"
					>
						View Roasters
					</Link>
				</div>
			</section>

			<Recommendations />

			{coffees.coffees && coffees.coffees.length > 0 && (
				<section className="space-y-6">
					<h2 className="text-3xl font-bold tracking-[0.1em] text-foreground">Latest Coffees</h2>
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
