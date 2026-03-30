'use client';

import { Bean, Droplets, Flame, Layers } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import CoffeeTrackButton from '@/components/CoffeeTrackButton';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { DomainCoffeeResponse } from '@/lib/api/generated/models';

const COUNTRIES = [
	'ethiopia',
	'colombia',
	'kenya',
	'brazil',
	'guatemala',
	'costa rica',
	'el salvador',
	'honduras',
	'panama',
	'nicaragua',
	'mexico',
	'peru',
	'bolivia',
	'ecuador',
	'rwanda',
	'burundi',
	'tanzania',
	'uganda',
	'indonesia',
	'india',
	'papua new guinea',
	'yemen',
	'congo',
	'malawi',
];

const PROCESSES = ['natural', 'washed', 'honey', 'anaerobic', 'semi-washed'];

/** Strip redundant country, process, roast, and filler from coffee display names. */
function cleanName(name: string, countryName?: string): string {
	let cleaned = name;

	// Remove FILTER / ESPRESSO roast suffixes
	cleaned = cleaned.replace(/\s*(FILTER|ESPRESSO)\s*$/i, '');

	// Remove process in parentheses: (Natural), (Washed), etc.
	for (const p of PROCESSES) {
		cleaned = cleaned.replace(new RegExp(`\\s*\\(${p}\\)\\s*`, 'i'), ' ');
	}

	// Remove country name in various positions
	const countries = countryName ? [countryName, ...COUNTRIES] : COUNTRIES;
	for (const country of countries) {
		const escaped = country.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
		// "- Ethiopia" or "| Ethiopia" or ", Ethiopia" at end
		cleaned = cleaned.replace(new RegExp(`\\s*[-|,]\\s*${escaped}\\s*$`, 'i'), '');
		// "Ethiopia -" or "Ethiopia |" or "Ethiopia," at start
		cleaned = cleaned.replace(new RegExp(`^\\s*${escaped}\\s*[-|,]\\s*`, 'i'), '');
		// ", Ethiopia," in the middle
		cleaned = cleaned.replace(new RegExp(`,\\s*${escaped}\\b`, 'i'), '');
		// "Colombia " at the very start (no separator, just a space before the rest)
		cleaned = cleaned.replace(new RegExp(`^${escaped}\\s+`, 'i'), '');
	}

	// Remove "Espresso" or "Filter" prefix (roast is shown as a badge)
	cleaned = cleaned.replace(/^(Espresso|Filter)\s+/i, '');

	// Remove standalone process words at end: ", Washed" or ", Natural"
	for (const p of PROCESSES) {
		cleaned = cleaned.replace(new RegExp(`,\\s*${p}\\s*$`, 'i'), '');
	}

	// Remove "Single Origin" filler
	cleaned = cleaned.replace(/\s*[-|]\s*Single Origin\b/i, '');
	cleaned = cleaned.replace(/\bSingle Origin\s*[-|]\s*/i, '');

	// Remove weight suffixes like "500g", "| 250g", "1kg"
	cleaned = cleaned.replace(/\s*[|]\s*\d+g\b/i, '');
	cleaned = cleaned.replace(/\s+\d+g\s*$/i, '');

	// Remove trailing/leading separators and whitespace
	cleaned = cleaned.replace(/\s*[|,-]\s*$/, '');
	cleaned = cleaned.replace(/^\s*[|,-]\s*/, '');

	return cleaned.trim();
}

export default function CoffeeCard({ coffee }: { coffee: DomainCoffeeResponse }) {
	const origin = [coffee.country_name, coffee.region_name].filter(Boolean).join(', ');
	const displayName = cleanName(coffee.name ?? '', coffee.country_name);

	return (
		<Link href={`/coffees/${coffee.id}`} className="block">
			<Card className="h-full flex flex-col shadow-sm transition-all hover:shadow-md hover:bg-muted/50">
				<CardHeader className="pb-1">
					<div className="flex items-start justify-between gap-1">
						<CardTitle className="text-base leading-snug">{displayName}</CardTitle>
						{coffee.id && (
								<div className="flex items-center gap-0.5">
									<CoffeeTrackButton coffeeId={coffee.id} variant="wishlist" />
									<CoffeeTrackButton coffeeId={coffee.id} variant="tried" />
								</div>
							)}
					</div>
					{coffee.roaster_name && (
						<Link
							href={`/roasters/${coffee.roaster_slug}`}
							onClick={(e) => e.stopPropagation()}
							className="flex items-center gap-1.5 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
						>
							{coffee.roaster_logo_url && (
								<Image
									src={coffee.roaster_logo_url}
									alt=""
									width={16}
									height={16}
									className="size-4 rounded-sm object-contain"
									unoptimized
								/>
							)}
							{coffee.roaster_name}
						</Link>
					)}
				</CardHeader>
				<CardContent className="flex-1 space-y-2 pt-0">
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
