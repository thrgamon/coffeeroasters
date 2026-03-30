import type { Metadata } from 'next';
import Link from 'next/link';
import CoffeeCard from '@/components/CoffeeCard';
import { Badge } from '@/components/ui/badge';
import type { DomainRoasterDetailResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export async function generateMetadata({ params }: { params: Promise<{ slug: string }> }): Promise<Metadata> {
	const { slug } = await params;
	const data = await apiFetch<DomainRoasterDetailResponse>(`/api/roasters/${slug}`);
	return { title: `${data.roaster?.name ?? 'Roaster'} | COFFEEROASTERS` };
}

export default async function RoasterDetailPage({ params }: { params: Promise<{ slug: string }> }) {
	const { slug } = await params;
	const data = await apiFetch<DomainRoasterDetailResponse>(`/api/roasters/${slug}`);

	const roaster = data.roaster;
	const coffees = data.coffees ?? [];
	const cafes = data.cafes ?? [];

	return (
		<div className="space-y-8">
			<div className="space-y-2">
				<Link href="/roasters" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
					&larr; All roasters
				</Link>
				<h1 className="text-3xl font-bold text-foreground">{roaster?.name}</h1>
				<div className="flex items-center gap-3">
					{roaster?.state && <Badge variant="secondary">{roaster.state}</Badge>}
					{roaster?.website && (
						<a
							href={roaster.website}
							target="_blank"
							rel="noopener noreferrer"
							className="text-sm text-accent hover:text-foreground"
						>
							{roaster.website}
						</a>
					)}
				</div>
			</div>

			{cafes.length > 0 && (
				<section className="space-y-4">
					<h2 className="text-xl font-bold uppercase tracking-wider text-accent">
						{cafes.length} cafe{cafes.length !== 1 ? 's' : ''}
					</h2>
					<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
						{cafes.map((cafe) => (
							<div key={cafe.id} className="rounded-lg border-2 border-border bg-card p-4 space-y-1">
								<p className="font-medium text-foreground">{cafe.name}</p>
								{cafe.address && <p className="text-sm text-muted-foreground">{cafe.address}</p>}
								{cafe.phone && <p className="text-sm text-muted-foreground">{cafe.phone}</p>}
							</div>
						))}
					</div>
				</section>
			)}

			<section className="space-y-4">
				<h2 className="text-xl font-bold uppercase tracking-wider text-accent">
					{coffees.length} coffee{coffees.length !== 1 ? 's' : ''}
				</h2>
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{coffees.map((coffee) => (
						<CoffeeCard key={coffee.id} coffee={coffee} />
					))}
				</div>
			</section>
		</div>
	);
}
