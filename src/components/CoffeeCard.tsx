import { Bean, Droplets, Flame, Globe, Grape } from 'lucide-react';
import Link from 'next/link';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { DomainCoffeeResponse } from '@/lib/api/generated/models';

export default function CoffeeCard({ coffee }: { coffee: DomainCoffeeResponse }) {
	return (
		<Card>
			<CardHeader className="pb-2">
				<CardTitle className="text-base">
					<Link href={`/coffees/${coffee.id}`} className="hover:underline">
						{coffee.name}
					</Link>
				</CardTitle>
				{coffee.roaster_name && (
					<Link
						href={`/roasters/${coffee.roaster_slug}`}
						className="text-sm text-muted-foreground hover:text-foreground"
					>
						{coffee.roaster_name}
					</Link>
				)}
			</CardHeader>
			<CardContent className="space-y-2">
				<div className="flex flex-wrap gap-1">
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
					{!coffee.in_stock && <Badge variant="destructive">Out of stock</Badge>}
				</div>
				{coffee.producer_name && coffee.producer_id && (
					<Link
						href={`/producers/${coffee.producer_id}`}
						className="block text-sm text-muted-foreground hover:text-foreground"
					>
						{coffee.producer_name}
					</Link>
				)}
				{coffee.tasting_notes && coffee.tasting_notes.length > 0 && (
					<div className="flex flex-wrap gap-1">
						{coffee.tasting_notes.map((note) => (
							<Badge key={note} variant="secondary" className="gap-1 font-normal">
								<Grape className="size-3" />
								{note}
							</Badge>
						))}
					</div>
				)}
				<div className="flex items-center justify-between text-sm">
					{coffee.price_cents ? (
						<span className="font-medium">${(coffee.price_cents / 100).toFixed(2)}</span>
					) : (
						<span />
					)}
					{coffee.weight_grams ? <span className="text-muted-foreground">{coffee.weight_grams}g</span> : null}
				</div>
				<div className="flex items-center gap-3">
					{coffee.product_url && (
						<a
							href={coffee.product_url}
							target="_blank"
							rel="noopener noreferrer"
							className="text-sm text-primary hover:underline"
						>
							View on roaster site
						</a>
					)}
					<Link
						href={`/coffees?similar_to=${coffee.id}`}
						className="text-sm text-muted-foreground hover:text-foreground"
					>
						Find similar
					</Link>
				</div>
				{coffee.similarity_score ? (
					<p className="text-xs text-muted-foreground">{Math.round(coffee.similarity_score * 100)}% match</p>
				) : null}
			</CardContent>
		</Card>
	);
}
