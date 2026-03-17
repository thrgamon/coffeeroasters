'use client';

import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useGetApiCountries } from '@/lib/api/generated/countries/countries';

export default function CountriesPage() {
	const { data, isLoading } = useGetApiCountries();

	return (
		<div className="space-y-6">
			<h1 className="text-3xl font-bold">Countries</h1>

			{isLoading && <p className="text-muted-foreground">Loading...</p>}

			{data && (
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{data.countries?.map((country) => (
						<Link key={country.code} href={`/countries/${country.code}`}>
							<Card className="transition-colors hover:border-primary/50">
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
			)}

			{data && (data.countries?.length ?? 0) === 0 && <p className="text-muted-foreground">No countries found.</p>}
		</div>
	);
}
