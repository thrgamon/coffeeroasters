import { Bean, Droplets, Flame, Layers } from 'lucide-react';
import Link from 'next/link';
import CoffeeTrackButton from '@/components/CoffeeTrackButton';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { DomainCoffeeResponse } from '@/lib/api/generated/models';

export default function CoffeeCard({ coffee }: { coffee: DomainCoffeeResponse }) {
	const origin = [coffee.country_name, coffee.region_name].filter(Boolean).join(', ');

	return (
		<Link href={`/coffees/${coffee.id}`} className="block">
			<Card className="shadow-sm transition-all hover:shadow-md hover:bg-muted/50">
				<CardHeader className="pb-1">
					<div className="flex items-start justify-between gap-1">
						<CardTitle className="text-base leading-snug">{coffee.name}</CardTitle>
						{coffee.id && <CoffeeTrackButton coffeeId={coffee.id} variant="like" />}
					</div>
					{coffee.roaster_name && <p className="text-xs text-muted-foreground">{coffee.roaster_name}</p>}
				</CardHeader>
				<CardContent className="space-y-2 pt-0">
					{origin && <p className="text-sm font-medium text-foreground/70">{origin}</p>}

					<div className="flex flex-wrap gap-1">
						{coffee.process && (
							<Badge variant="outline" className="gap-1 text-xs">
								<Droplets className="size-3" />
								{coffee.process}
							</Badge>
						)}
						{coffee.roast_level && (
							<Badge variant="outline" className="gap-1 text-xs">
								<Flame className="size-3" />
								{coffee.roast_level}
							</Badge>
						)}
						{coffee.variety && (
							<Badge variant="outline" className="gap-1 text-xs">
								<Bean className="size-3" />
								{coffee.variety}
							</Badge>
						)}
						{coffee.is_blend && (
							<Badge variant="outline" className="gap-1 text-xs">
								<Layers className="size-3" />
								Blend
							</Badge>
						)}
						{!coffee.in_stock && (
							<Badge variant="destructive" className="text-xs">
								Out of stock
							</Badge>
						)}
					</div>

					{coffee.tasting_notes && coffee.tasting_notes.length > 0 && (
						<p className="text-xs text-muted-foreground leading-relaxed">{coffee.tasting_notes.join(', ')}</p>
					)}

					<div className="flex items-center justify-between text-sm">
						<div className="flex items-baseline gap-2">
							{coffee.price_per_100g_min ? (
								<span className="font-medium">
									{coffee.price_per_100g_min === coffee.price_per_100g_max
										? `$${(coffee.price_per_100g_min / 100).toFixed(2)} / 100g`
										: `$${(coffee.price_per_100g_min / 100).toFixed(2)} - $${((coffee.price_per_100g_max ?? coffee.price_per_100g_min) / 100).toFixed(2)} / 100g`}
								</span>
							) : coffee.price_cents ? (
								<span className="font-medium">${(coffee.price_cents / 100).toFixed(2)}</span>
							) : null}
						</div>
						{coffee.weight_grams ? <span className="text-muted-foreground">{coffee.weight_grams}g</span> : null}
					</div>
					{coffee.similarity_score ? (
						<p className="text-xs text-muted-foreground">{Math.round(coffee.similarity_score * 100)}% match</p>
					) : null}
				</CardContent>
			</Card>
		</Link>
	);
}
