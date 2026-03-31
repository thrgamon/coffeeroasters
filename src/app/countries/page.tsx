import type { Metadata } from 'next';
import Link from 'next/link';
import type { DomainCountryListResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export const metadata: Metadata = { title: 'Origins | COFFEEROASTERS' };
export const revalidate = 60;

export default async function CountriesPage() {
	const data = await apiFetch<DomainCountryListResponse>('/api/countries');

	const countries = (data.countries ?? []).sort((a, b) => (b.coffee_count ?? 0) - (a.coffee_count ?? 0));

	return (
		<div className="space-y-6">
			<div>
				<h1 className="text-3xl font-bold uppercase tracking-wider text-foreground">Origins</h1>
				<p className="mt-1 text-muted-foreground">
					{countries.length} coffee-producing countries, sorted by number of coffees available.
				</p>
			</div>

			{countries.length > 0 ? (
				<div className="grid gap-3 grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
					{countries.map((country) => (
						<Link
							key={country.code}
							href={`/countries/${country.code}`}
							className="flex items-baseline justify-between gap-2 border-2 border-border bg-card px-4 py-3 transition-all hover:border-accent"
						>
							<span className="font-medium text-foreground">{country.name}</span>
							<span className="text-sm text-foreground font-mono tabular-nums">{country.coffee_count}</span>
						</Link>
					))}
				</div>
			) : (
				<p className="text-muted-foreground">No countries found.</p>
			)}
		</div>
	);
}
