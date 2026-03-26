'use client';

import { Loader2, X } from 'lucide-react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useCallback, useTransition } from 'react';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import type { DomainCountryResponse } from '@/lib/api/generated/models';

const PROCESSES = ['washed', 'natural', 'honey', 'anaerobic', 'wet-hulled', 'experimental'];
const ROAST_LEVELS = ['light', 'medium-light', 'medium', 'medium-dark', 'dark'];
const ROASTER_STATES = ['VIC', 'NSW', 'QLD', 'WA', 'SA', 'TAS', 'ACT', 'NT'];

const VARIETIES = [
	'bourbon',
	'typica',
	'caturra',
	'catuai',
	'gesha',
	'sl28',
	'sl34',
	'pacamara',
	'heirloom',
	'castillo',
	'catimor',
	'mundo-novo',
	'yellow-bourbon',
	'pink-bourbon',
	'red-bourbon',
	'sidra',
	'wush-wush',
];

interface CoffeeFiltersProps {
	countries: DomainCountryResponse[];
}

export default function CoffeeFilters({ countries }: CoffeeFiltersProps) {
	const router = useRouter();
	const searchParams = useSearchParams();
	const [isPending, startTransition] = useTransition();

	const q = searchParams.get('q') ?? '';
	const origin = searchParams.get('origin') ?? '';
	const process = searchParams.get('process') ?? '';
	const roast = searchParams.get('roast') ?? '';
	const variety = searchParams.get('variety') ?? '';
	const roasterState = searchParams.get('roaster_state') ?? '';

	const updateParams = useCallback(
		(key: string, value: string) => {
			const params = new URLSearchParams(searchParams.toString());
			if (value) {
				params.set(key, value);
			} else {
				params.delete(key);
			}
			params.delete('page');
			startTransition(() => {
				router.push(`/coffees?${params.toString()}`);
			});
		},
		[router, searchParams],
	);

	const activeFilters: { key: string; label: string }[] = [];
	if (q) activeFilters.push({ key: 'q', label: `Search: ${q}` });
	if (origin) {
		const countryName = countries.find((c) => c.code === origin)?.name ?? origin;
		activeFilters.push({ key: 'origin', label: `Country: ${countryName}` });
	}
	if (process) activeFilters.push({ key: 'process', label: `Process: ${process}` });
	if (roast) activeFilters.push({ key: 'roast', label: `Roast: ${roast}` });
	if (variety) activeFilters.push({ key: 'variety', label: `Variety: ${variety}` });
	if (roasterState) activeFilters.push({ key: 'roaster_state', label: `Roaster location: ${roasterState}` });

	return (
		<div className="space-y-2">
			<div className="flex flex-wrap gap-3">
				<Input
					placeholder="Search coffees..."
					defaultValue={q}
					onChange={(e) => updateParams('q', e.target.value)}
					className="w-64"
				/>
				<Select value={origin || 'all'} onValueChange={(v) => updateParams('origin', v === 'all' ? '' : v)}>
					<SelectTrigger className="w-48">
						<SelectValue placeholder="Country" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="all">All countries</SelectItem>
						{countries.map((c) => (
							<SelectItem key={c.code} value={c.code ?? ''}>
								{c.name} ({c.coffee_count})
							</SelectItem>
						))}
					</SelectContent>
				</Select>
				<Select value={process || 'all'} onValueChange={(v) => updateParams('process', v === 'all' ? '' : v)}>
					<SelectTrigger className="w-40">
						<SelectValue placeholder="Process" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="all">All processes</SelectItem>
						{PROCESSES.map((p) => (
							<SelectItem key={p} value={p}>
								{p}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
				<Select value={roast || 'all'} onValueChange={(v) => updateParams('roast', v === 'all' ? '' : v)}>
					<SelectTrigger className="w-40">
						<SelectValue placeholder="Roast level" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="all">All roasts</SelectItem>
						{ROAST_LEVELS.map((r) => (
							<SelectItem key={r} value={r}>
								{r}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
				<Select value={variety || 'all'} onValueChange={(v) => updateParams('variety', v === 'all' ? '' : v)}>
					<SelectTrigger className="w-40">
						<SelectValue placeholder="Variety" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="all">All varieties</SelectItem>
						{VARIETIES.map((v) => (
							<SelectItem key={v} value={v}>
								{v}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
				<Select value={roasterState || 'all'} onValueChange={(v) => updateParams('roaster_state', v === 'all' ? '' : v)}>
					<SelectTrigger className="w-48">
						<SelectValue placeholder="Roaster location" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="all">All locations</SelectItem>
						{ROASTER_STATES.map((s) => (
							<SelectItem key={s} value={s}>
								{s}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
				{isPending && <Loader2 className="size-5 animate-spin self-center text-muted-foreground" />}
			</div>
			{activeFilters.length > 0 && (
				<div className="flex flex-wrap gap-2">
					{activeFilters.map((filter) => (
						<Badge key={filter.key} variant="secondary" className="gap-1">
							{filter.label}
							<button
								type="button"
								onClick={(e) => {
									e.preventDefault();
									updateParams(filter.key, '');
								}}
								className="ml-0.5 rounded-full p-0.5 hover:bg-muted-foreground/20"
							>
								<X className="size-3" />
								<span className="sr-only">Remove {filter.label} filter</span>
							</button>
						</Badge>
					))}
				</div>
			)}
		</div>
	);
}
