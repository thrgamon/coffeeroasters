import type { Metadata } from 'next';
import Image from 'next/image';
import Link from 'next/link';
import CoffeeCard from '@/components/CoffeeCard';
import { Badge } from '@/components/ui/badge';
import type { DomainCafeResponse, DomainCoffeeResponse, DomainRoasterResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

interface RoasterDetailData {
	roaster: DomainRoasterResponse;
	available_coffees?: DomainCoffeeResponse[];
	unavailable_coffees?: DomainCoffeeResponse[];
	coffees?: DomainCoffeeResponse[];
	cafes?: DomainCafeResponse[];
}

function deriveStats(coffees: DomainCoffeeResponse[]) {
	const origins = new Map<string, { code: string; count: number }>();
	const processes = new Set<string>();
	let minPrice = Number.POSITIVE_INFINITY;
	let maxPrice = 0;

	for (const c of coffees) {
		if (c.country_name && c.country_code) {
			const existing = origins.get(c.country_code);
			if (existing) {
				existing.count++;
			} else {
				origins.set(c.country_code, { code: c.country_code, count: 1 });
			}
		}
		if (c.process) processes.add(c.process);
		if (c.price_per_100g_min && c.price_per_100g_min > 0) {
			minPrice = Math.min(minPrice, c.price_per_100g_min);
		}
		if (c.price_per_100g_max && c.price_per_100g_max > 0) {
			maxPrice = Math.max(maxPrice, c.price_per_100g_max);
		}
	}

	const sortedOrigins = [...origins.entries()]
		.map(([code, { count }]) => {
			const coffee = coffees.find((c) => c.country_code === code);
			return { code, name: coffee?.country_name ?? code, count };
		})
		.sort((a, b) => b.count - a.count);

	return {
		origins: sortedOrigins,
		processes: [...processes].sort(),
		priceRange:
			minPrice < Number.POSITIVE_INFINITY
				? { min: (minPrice / 100).toFixed(2), max: (maxPrice / 100).toFixed(2) }
				: null,
	};
}

export async function generateMetadata({ params }: { params: Promise<{ slug: string }> }): Promise<Metadata> {
	const { slug } = await params;
	const data = await apiFetch<RoasterDetailData>(`/api/roasters/${slug}`);
	return { title: `${data.roaster?.name ?? 'Roaster'} | COFFEEROASTERS` };
}

export default async function RoasterDetailPage({ params }: { params: Promise<{ slug: string }> }) {
	const { slug } = await params;
	const data = await apiFetch<RoasterDetailData>(`/api/roasters/${slug}`);

	const roaster = data.roaster;
	const available = data.available_coffees ?? data.coffees ?? [];
	const unavailable = data.unavailable_coffees ?? [];
	const cafes = data.cafes ?? [];
	const stats = deriveStats(available);

	return (
		<div className="space-y-8">
			<div className="space-y-4">
				<Link href="/roasters" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
					&larr; All roasters
				</Link>
				<div className="flex items-center gap-5">
					{roaster?.logo_url && (
						<Image
							src={roaster.logo_url}
							alt={`${roaster.name} logo`}
							width={64}
							height={64}
							className="size-16 object-contain"
							unoptimized
						/>
					)}
					<div className="space-y-1">
						<h1 className="text-3xl font-bold text-foreground">{roaster?.name}</h1>
						<div className="flex items-center gap-3">
							{roaster?.state && <Badge variant="secondary">{roaster.state}</Badge>}
							{roaster?.website && (
								<a
									href={roaster.website}
									target="_blank"
									rel="noopener noreferrer"
									className="text-sm text-accent hover:text-foreground transition-colors"
								>
									{roaster.website.replace(/^https?:\/\//, '')}
								</a>
							)}
						</div>
					</div>
				</div>
			</div>

			{(stats.origins.length > 0 || stats.priceRange || stats.processes.length > 0) && (
				<div className="grid gap-4 sm:grid-cols-3">
					{stats.origins.length > 0 && (
						<div className="border border-border bg-card p-4 space-y-2">
							<h3>Origins</h3>
							<div className="flex flex-wrap gap-2">
								{stats.origins.map((o) => (
									<Link key={o.code} href={`/coffees?origin=${o.code}`}>
										<Badge variant="outline">
											{o.name} ({o.count})
										</Badge>
									</Link>
								))}
							</div>
						</div>
					)}
					{stats.processes.length > 0 && (
						<div className="border border-border bg-card p-4 space-y-2">
							<h3>Processes</h3>
							<div className="flex flex-wrap gap-2">
								{stats.processes.map((p) => (
									<Badge key={p} variant="outline">
										{p}
									</Badge>
								))}
							</div>
						</div>
					)}
					{stats.priceRange && (
						<div className="border border-border bg-card p-4 space-y-2">
							<h3>Price range</h3>
							<p className="font-mono text-lg font-bold">
								${stats.priceRange.min} &ndash; ${stats.priceRange.max}
							</p>
							<p className="text-xs text-muted-foreground">per 100g</p>
						</div>
					)}
				</div>
			)}

			{cafes.length > 0 && (
				<section className="space-y-4">
					<h3>
						{cafes.length} cafe{cafes.length !== 1 ? 's' : ''}
					</h3>
					<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
						{cafes.map((cafe) => (
							<div key={cafe.id} className="border border-border bg-card p-4 space-y-1">
								<p className="font-medium text-foreground">{cafe.name}</p>
								{cafe.address && <p className="text-sm text-muted-foreground">{cafe.address}</p>}
								{cafe.phone && <p className="text-sm text-muted-foreground">{cafe.phone}</p>}
							</div>
						))}
					</div>
				</section>
			)}

			<section className="space-y-4">
				<h3>
					{available.length} coffee{available.length !== 1 ? 's' : ''} in stock
				</h3>
				{available.length > 0 ? (
					<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
						{available.map((coffee) => (
							<CoffeeCard key={coffee.id} coffee={coffee} />
						))}
					</div>
				) : (
					<p className="text-muted-foreground">No coffees currently in stock.</p>
				)}
			</section>

			{unavailable.length > 0 && (
				<section className="space-y-4">
					<h3 className="text-muted-foreground">{unavailable.length} unavailable</h3>
					<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 opacity-60">
						{unavailable.map((coffee) => (
							<CoffeeCard key={coffee.id} coffee={coffee} />
						))}
					</div>
				</section>
			)}
		</div>
	);
}
