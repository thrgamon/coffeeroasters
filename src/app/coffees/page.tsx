import type { Metadata } from 'next';
import Link from 'next/link';
import { Suspense } from 'react';
import CoffeeCard from '@/components/CoffeeCard';
import CoffeeFilters from '@/components/CoffeeFilters';
import type {
	DomainCoffeeDetailResponse,
	DomainCoffeeListResponse,
	DomainCountryListResponse,
} from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export const metadata: Metadata = { title: 'Coffees | COFFEEROASTERS' };

interface CoffeesPageProps {
	searchParams: Promise<Record<string, string | string[] | undefined>>;
}

export default async function CoffeesPage({ searchParams }: CoffeesPageProps) {
	const sp = await searchParams;

	const similarTo = typeof sp.similar_to === 'string' ? sp.similar_to : undefined;
	const q = typeof sp.q === 'string' ? sp.q : undefined;
	const origin = typeof sp.origin === 'string' ? sp.origin : undefined;
	const process = typeof sp.process === 'string' ? sp.process : undefined;
	const roast = typeof sp.roast === 'string' ? sp.roast : undefined;
	const variety = typeof sp.variety === 'string' ? sp.variety : undefined;
	const roasterState = typeof sp.roaster_state === 'string' ? sp.roaster_state : undefined;
	const page = typeof sp.page === 'string' ? Number(sp.page) : 1;
	const pageSize = 20;

	const params = new URLSearchParams();
	params.set('page', String(page));
	params.set('page_size', String(pageSize));
	if (similarTo) {
		params.set('similar_to', similarTo);
	} else {
		if (q) params.set('q', q);
		if (origin) params.set('origin', origin);
		if (process) params.set('process', process);
		if (roast) params.set('roast', roast);
		if (variety) params.set('variety', variety);
		if (roasterState) params.set('roaster_state', roasterState);
	}

	const [data, countries, sourceCoffee] = await Promise.all([
		apiFetch<DomainCoffeeListResponse>(`/api/coffees?${params.toString()}`),
		apiFetch<DomainCountryListResponse>('/api/countries'),
		similarTo ? apiFetch<DomainCoffeeDetailResponse>(`/api/coffees/${similarTo}`) : null,
	]);

	const totalPages = Math.ceil((data.total_count ?? 0) / pageSize);

	function pageUrl(p: number) {
		const next = new URLSearchParams();
		if (q) next.set('q', q);
		if (origin) next.set('origin', origin);
		if (process) next.set('process', process);
		if (roast) next.set('roast', roast);
		if (variety) next.set('variety', variety);
		if (roasterState) next.set('roaster_state', roasterState);
		if (similarTo) next.set('similar_to', similarTo);
		if (p > 1) next.set('page', String(p));
		const qs = next.toString();
		return `/coffees${qs ? `?${qs}` : ''}`;
	}

	return (
		<div className="space-y-6">
			{similarTo && sourceCoffee ? (
				<div className="space-y-2">
					<h1 className="text-3xl font-bold uppercase tracking-wider text-foreground">
						Similar to {sourceCoffee.name}
					</h1>
					<p className="text-sm text-muted-foreground">
						{sourceCoffee.roaster_name}
						{sourceCoffee.country_name ? ` \u00b7 ${sourceCoffee.country_name}` : ''}
						{sourceCoffee.process ? ` \u00b7 ${sourceCoffee.process}` : ''}
					</p>
					<Link href="/coffees" className="inline-block text-sm text-accent hover:text-foreground">
						Back to all coffees
					</Link>
				</div>
			) : (
				<h1 className="text-3xl font-bold uppercase tracking-wider text-foreground">Coffees</h1>
			)}

			{!similarTo && (
				<Suspense>
					<CoffeeFilters countries={countries.countries ?? []} />
				</Suspense>
			)}

			<p className="text-sm text-muted-foreground font-mono">
				{data.total_count ?? 0} coffee{(data.total_count ?? 0) !== 1 ? 's' : ''} found
			</p>

			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{data.coffees?.map((coffee) => (
					<CoffeeCard key={coffee.id} coffee={coffee} />
				))}
			</div>

			{totalPages > 1 && (
				<div className="flex justify-center gap-2">
					{page > 1 ? (
						<Link href={pageUrl(page - 1)} className="rounded border border-border px-3 py-1 text-sm text-foreground hover:border-accent hover:text-accent transition-colors">
							Previous
						</Link>
					) : (
						<span className="rounded border border-border px-3 py-1 text-sm opacity-30 text-foreground">Previous</span>
					)}
					<span className="px-3 py-1 text-sm text-muted-foreground font-mono">
						{page} / {totalPages}
					</span>
					{page < totalPages ? (
						<Link href={pageUrl(page + 1)} className="rounded border border-border px-3 py-1 text-sm text-foreground hover:border-accent hover:text-accent transition-colors">
							Next
						</Link>
					) : (
						<span className="rounded border border-border px-3 py-1 text-sm opacity-30 text-foreground">Next</span>
					)}
				</div>
			)}
		</div>
	);
}
