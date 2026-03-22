import type { Metadata } from 'next';
import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { DomainCountryListResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export const metadata: Metadata = { title: 'Countries | Coffeeroasters' };
export const dynamic = 'force-dynamic';

export default async function CountriesPage() {
	const data = await apiFetch<DomainCountryListResponse>('/api/countries');

	return (
		<div className="space-y-6">
			<h1 className="text-3xl font-bold">Countries</h1>

			{data.countries && data.countries.length > 0 ? (
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{data.countries.map((country) => (
						<Link key={country.code} href={`/countries/${country.code}`}>
							<Card className="shadow-sm transition-all hover:shadow-md hover:bg-muted/50">
								<CardHeader>
									<CardTitle className="text-lg">{country.name}</CardTitle>
								</CardHeader>
								<CardContent>
									<p className="text-sm text-muted-foreground">
										{country.coffee_count} coffee{country.coffee_count !== 1 ? 's' : ''}
									</p>
								</CardContent>
							</Card>
						</Link>
					))}
				</div>
			) : (
				<p className="text-muted-foreground">No countries found.</p>
			)}
		</div>
	);
}
