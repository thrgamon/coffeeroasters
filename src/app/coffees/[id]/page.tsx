import { Bean, Droplets, Flame, Globe, Grape, Layers, Sprout } from 'lucide-react';
import type { Metadata } from 'next';
import Image from 'next/image';
import Link from 'next/link';
import CoffeeTrackButton from '@/components/CoffeeTrackButton';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent } from '@/components/ui/card';
import type { DomainCoffeeDetailResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

export async function generateMetadata({ params }: { params: Promise<{ id: string }> }): Promise<Metadata> {
	const { id } = await params;
	const coffee = await apiFetch<DomainCoffeeDetailResponse>(`/api/coffees/${id}`);
	return { title: `${coffee.name} | COFFEEROASTERS` };
}

export default async function CoffeeDetailPage({ params }: { params: Promise<{ id: string }> }) {
	const { id } = await params;
	const coffee = await apiFetch<DomainCoffeeDetailResponse>(`/api/coffees/${id}`);

	return (
		<div className="space-y-8">
			<div className="space-y-4">
				<div className="flex items-start gap-6">
					{coffee.image_url && (
						<Image
							src={coffee.image_url}
							alt={coffee.name ?? ''}
							width={192}
							height={192}
							unoptimized
							className="h-48 w-48 rounded-lg object-cover border border-border/50"
						/>
					)}
					<div className="space-y-2">
						<div className="flex items-center gap-2">
							<h1 className="text-3xl font-bold text-snow">{coffee.name}</h1>
							{coffee.id && (
								<div className="flex items-center gap-1">
									<CoffeeTrackButton coffeeId={coffee.id} coffeeName={coffee.name} variant="wishlist" size="md" />
									<CoffeeTrackButton coffeeId={coffee.id} coffeeName={coffee.name} variant="log" size="md" />
								</div>
							)}
						</div>
						{coffee.roaster_name && (
							<div className="flex items-center gap-3">
								<Link
									href={`/roasters/${coffee.roaster_slug}`}
									className="text-lg text-dusty-rose hover:text-gold transition-colors"
								>
									{coffee.roaster_name}
								</Link>
								{coffee.product_url && (
									<a
										href={coffee.product_url}
										target="_blank"
										rel="noopener noreferrer"
										className="inline-block rounded bg-gold px-4 py-2 text-sm font-bold uppercase tracking-wider text-rich-mahogany hover:bg-gold/90 transition-colors"
									>
										View on roaster site
									</a>
								)}
							</div>
						)}
						<div className="flex flex-wrap gap-2">
							{coffee.country_name && (
								<Link href={`/countries/${coffee.country_code}`}>
									<Badge variant="secondary" className="cursor-pointer gap-1 hover:bg-gold/10 hover:text-gold">
										<Globe className="size-3" />
										{coffee.country_name}
									</Badge>
								</Link>
							)}
							{coffee.region_name && coffee.region_id && (
								<Link href={`/regions/${coffee.region_id}`}>
									<Badge variant="secondary" className="cursor-pointer gap-1 hover:bg-gold/10 hover:text-gold">
										<Globe className="size-3" />
										{coffee.region_name}
									</Badge>
								</Link>
							)}
							{coffee.process && (
								<Badge variant="outline" className="gap-1">
									<Droplets className="size-3" />
									{coffee.process}
								</Badge>
							)}
							{coffee.roast_level && (
								<Badge variant="outline" className="gap-1">
									<Flame className="size-3" />
									{coffee.roast_level}
								</Badge>
							)}
							{coffee.variety && (
								<Badge variant="outline" className="gap-1">
									<Bean className="size-3" />
									{coffee.variety}
								</Badge>
							)}
							{coffee.species && (
								<Badge variant="outline" className="gap-1">
									<Sprout className="size-3" />
									{coffee.species}
								</Badge>
							)}
							{coffee.is_blend && (
								<Badge variant="outline" className="gap-1">
									<Layers className="size-3" />
									Blend
								</Badge>
							)}
							{!coffee.in_stock && <Badge variant="destructive">Out of stock</Badge>}
						</div>
					</div>
				</div>

				{coffee.producer_name && coffee.producer_id && (
					<Link
						href={`/producers/${coffee.producer_id}`}
						className="block text-sm text-grey-olive hover:text-gold transition-colors"
					>
						Producer: {coffee.producer_name}
					</Link>
				)}

				{coffee.tasting_notes && coffee.tasting_notes.length > 0 && (
					<div>
						<h3 className="mb-1 text-sm font-bold uppercase tracking-wider text-dusty-rose">Tasting notes</h3>
						<div className="flex flex-wrap gap-1">
							{coffee.tasting_notes.map((note) => (
								<Badge key={note} variant="secondary" className="gap-1">
									<Grape className="size-3" />
									{note}
								</Badge>
							))}
						</div>
					</div>
				)}

				{coffee.description && (
					<div>
						<h3 className="mb-1 text-sm font-bold uppercase tracking-wider text-dusty-rose">About this coffee</h3>
						<p className="text-sm text-snow/70 whitespace-pre-line">{coffee.description}</p>
					</div>
				)}

				{coffee.blend_components && coffee.blend_components.length > 0 && (
					<div>
						<h3 className="mb-1 text-sm font-bold uppercase tracking-wider text-dusty-rose">Blend components</h3>
						<div className="space-y-1">
							{coffee.blend_components.map((comp) => (
								<div
									key={`${comp.country_code}-${comp.region_id}-${comp.variety}`}
									className="flex items-center gap-2 text-sm"
								>
									{comp.country_name && (
										<Link href={`/countries/${comp.country_code}`}>
											<Badge variant="secondary" className="cursor-pointer gap-1 hover:bg-gold/10 hover:text-gold">
												<Globe className="size-3" />
												{comp.country_name}
											</Badge>
										</Link>
									)}
									{comp.region_name && comp.region_id && (
										<Link href={`/regions/${comp.region_id}`}>
											<Badge variant="secondary" className="cursor-pointer gap-1 hover:bg-gold/10 hover:text-gold">
												<Globe className="size-3" />
												{comp.region_name}
											</Badge>
										</Link>
									)}
									{comp.variety && (
										<Badge variant="outline" className="gap-1">
											<Bean className="size-3" />
											{comp.variety}
										</Badge>
									)}
									{comp.percentage ? <span className="text-grey-olive font-mono">{comp.percentage}%</span> : null}
								</div>
							))}
						</div>
					</div>
				)}

				<div className="flex items-center gap-4 text-sm">
					{coffee.price_per_100g_min ? (
						<span className="text-lg font-bold text-gold font-mono">
							{coffee.price_per_100g_min === coffee.price_per_100g_max
								? `$${(coffee.price_per_100g_min / 100).toFixed(2)} / 100g`
								: `$${(coffee.price_per_100g_min / 100).toFixed(2)} - $${((coffee.price_per_100g_max ?? coffee.price_per_100g_min) / 100).toFixed(2)} / 100g`}
						</span>
					) : coffee.price_cents ? (
						<span className="text-lg font-bold text-gold font-mono">${(coffee.price_cents / 100).toFixed(2)}</span>
					) : null}
					{coffee.price_cents ? (
						<span className="text-grey-olive font-mono">
							${(coffee.price_cents / 100).toFixed(2)} / {coffee.weight_grams}g
						</span>
					) : null}
				</div>
			</div>

			{coffee.similar_coffees && coffee.similar_coffees.length > 0 && (
				<div className="space-y-3">
					<h3 className="text-lg font-bold uppercase tracking-wider text-dusty-rose">You might also like</h3>
					<div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
						{coffee.similar_coffees.map((similar) => (
							<Link key={similar.id} href={`/coffees/${similar.id}`}>
								<Card className="h-full border-border/50 transition-all hover:border-gold/30 hover:shadow-[0_0_20px_rgba(255,213,0,0.05)]">
									<CardContent className="p-4 space-y-1.5">
										<p className="font-medium text-sm leading-snug text-snow">{similar.name}</p>
										{similar.roaster_name && <p className="text-xs text-dusty-rose">{similar.roaster_name}</p>}
										{similar.reasons && similar.reasons.length > 0 && (
											<p className="text-xs text-grey-olive">{similar.reasons.join(' \u00b7 ')}</p>
										)}
									</CardContent>
								</Card>
							</Link>
						))}
					</div>
				</div>
			)}
		</div>
	);
}
