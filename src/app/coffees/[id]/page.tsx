import { Bean, Droplets, Flame, Globe, Grape, Layers, Sprout } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import CoffeeTrackButton from '@/components/CoffeeTrackButton';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent } from '@/components/ui/card';
import type { DomainCoffeeDetailResponse } from '@/lib/api/generated/models';
import { apiFetch } from '@/lib/api/server';

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
							className="h-48 w-48 rounded-lg object-cover"
						/>
					)}
					<div className="space-y-2">
						<div className="flex items-center gap-2">
							<h1 className="text-3xl font-bold">{coffee.name}</h1>
							{coffee.id && (
								<div className="flex items-center gap-1">
									<CoffeeTrackButton coffeeId={coffee.id} variant="like" size="md" />
									<CoffeeTrackButton coffeeId={coffee.id} variant="tried" size="md" />
								</div>
							)}
						</div>
						{coffee.roaster_name && (
							<Link
								href={`/roasters/${coffee.roaster_slug}`}
								className="text-lg text-muted-foreground hover:text-foreground"
							>
								{coffee.roaster_name}
							</Link>
						)}
						<div className="flex flex-wrap gap-2">
							{coffee.country_name && (
								<Link href={`/countries/${coffee.country_code}`}>
									<Badge variant="secondary" className="cursor-pointer gap-1 hover:bg-accent">
										<Globe className="size-3" />
										{coffee.country_name}
									</Badge>
								</Link>
							)}
							{coffee.region_name && coffee.region_id && (
								<Link href={`/regions/${coffee.region_id}`}>
									<Badge variant="secondary" className="cursor-pointer gap-1 hover:bg-accent">
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
						className="block text-sm text-muted-foreground hover:text-foreground"
					>
						Producer: {coffee.producer_name}
					</Link>
				)}

				{coffee.tasting_notes && coffee.tasting_notes.length > 0 && (
					<div>
						<h3 className="mb-1 text-sm font-medium">Tasting notes</h3>
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
						<h3 className="mb-1 text-sm font-medium">About this coffee</h3>
						<p className="text-sm text-muted-foreground whitespace-pre-line">{coffee.description}</p>
					</div>
				)}

				{coffee.blend_components && coffee.blend_components.length > 0 && (
					<div>
						<h3 className="mb-1 text-sm font-medium">Blend components</h3>
						<div className="space-y-1">
							{coffee.blend_components.map((comp) => (
								<div
									key={`${comp.country_code}-${comp.region_id}-${comp.variety}`}
									className="flex items-center gap-2 text-sm"
								>
									{comp.country_name && (
										<Link href={`/countries/${comp.country_code}`}>
											<Badge variant="secondary" className="cursor-pointer gap-1 hover:bg-accent">
												<Globe className="size-3" />
												{comp.country_name}
											</Badge>
										</Link>
									)}
									{comp.region_name && comp.region_id && (
										<Link href={`/regions/${comp.region_id}`}>
											<Badge variant="secondary" className="cursor-pointer gap-1 hover:bg-accent">
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
									{comp.percentage ? <span className="text-muted-foreground">{comp.percentage}%</span> : null}
								</div>
							))}
						</div>
					</div>
				)}

				<div className="flex items-center gap-4 text-sm">
					{coffee.price_per_100g_min ? (
						<span className="text-lg font-medium">
							{coffee.price_per_100g_min === coffee.price_per_100g_max
								? `$${(coffee.price_per_100g_min / 100).toFixed(2)} / 100g`
								: `$${(coffee.price_per_100g_min / 100).toFixed(2)} - $${((coffee.price_per_100g_max ?? coffee.price_per_100g_min) / 100).toFixed(2)} / 100g`}
						</span>
					) : coffee.price_cents ? (
						<span className="text-lg font-medium">${(coffee.price_cents / 100).toFixed(2)}</span>
					) : null}
					{coffee.price_cents ? (
						<span className="text-muted-foreground">
							${(coffee.price_cents / 100).toFixed(2)} / {coffee.weight_grams}g
						</span>
					) : null}
				</div>

				{coffee.product_url && (
					<a
						href={coffee.product_url}
						target="_blank"
						rel="noopener noreferrer"
						className="inline-block rounded-md bg-primary px-4 py-2 text-sm text-primary-foreground hover:bg-primary/90"
					>
						View on roaster site
					</a>
				)}
			</div>

			{coffee.similar_coffees && coffee.similar_coffees.length > 0 && (
				<div className="space-y-3">
					<h3 className="text-lg font-medium text-muted-foreground">You might also like</h3>
					<div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
						{coffee.similar_coffees.map((similar) => (
							<Link key={similar.id} href={`/coffees/${similar.id}`}>
								<Card className="h-full shadow-sm transition-all hover:shadow-md hover:bg-accent/50">
									<CardContent className="p-4 space-y-1.5">
										<p className="font-medium text-sm leading-snug">{similar.name}</p>
										{similar.roaster_name && <p className="text-xs text-muted-foreground">{similar.roaster_name}</p>}
										{similar.reasons && similar.reasons.length > 0 && (
											<p className="text-xs text-muted-foreground">{similar.reasons.join(' · ')}</p>
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
