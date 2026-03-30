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
			<section className="relative overflow-hidden rounded-lg border border-border/50 bg-gradient-to-br from-rich-mahogany via-background to-rich-mahogany px-6 py-20 text-center retro-blur scanlines">
				<div className="relative z-10 space-y-6">
					<p className="text-xs font-medium tracking-[0.3em] uppercase text-dusty-rose">
						Australian Specialty Coffee
					</p>
					<h1 className="text-5xl font-bold tracking-tight text-gold glow-gold sm:text-6xl">
						COFFEEROASTERS
					</h1>
					<p className="mx-auto max-w-md text-lg text-snow/70">
						Discover the best indie roasters across Australia
					</p>
					<div className="flex justify-center gap-10 text-sm font-mono tracking-wider text-dusty-rose">
						<span className="flex flex-col items-center">
							<span className="text-2xl font-bold text-gold">{stats.roaster_count}</span>
							roasters
						</span>
						<span className="flex flex-col items-center">
							<span className="text-2xl font-bold text-gold">{stats.coffee_count}</span>
							coffees
						</span>
						<span className="flex flex-col items-center">
							<span className="text-2xl font-bold text-gold">{stats.origins?.length ?? 0}</span>
							origins
						</span>
					</div>
					<div className="flex justify-center gap-4 pt-2">
						<Link
							href="/coffees"
							className="rounded bg-gold px-6 py-2.5 text-sm font-bold uppercase tracking-wider text-rich-mahogany hover:bg-gold/90 transition-colors"
						>
							Browse Coffees
						</Link>
						<Link
							href="/roasters"
							className="rounded border border-dusty-rose/50 px-6 py-2.5 text-sm font-bold uppercase tracking-wider text-dusty-rose hover:bg-dusty-rose/10 transition-colors"
						>
							View Roasters
						</Link>
					</div>
				</div>
			</section>

			<Recommendations />

			{coffees.coffees && coffees.coffees.length > 0 && (
				<section className="space-y-6">
					<h2 className="text-2xl font-bold uppercase tracking-wider text-snow">Latest Coffees</h2>
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
