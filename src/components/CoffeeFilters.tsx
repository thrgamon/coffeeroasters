'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { useCallback } from 'react';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import type { DomainCountryResponse } from '@/lib/api/generated/models';

const PROCESSES = ['washed', 'natural', 'honey', 'anaerobic', 'wet-hulled', 'experimental'];
const ROAST_LEVELS = ['light', 'medium-light', 'medium', 'medium-dark', 'dark'];
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

	const q = searchParams.get('q') ?? '';
	const origin = searchParams.get('origin') ?? '';
	const process = searchParams.get('process') ?? '';
	const roast = searchParams.get('roast') ?? '';
	const variety = searchParams.get('variety') ?? '';

	const updateParams = useCallback(
		(key: string, value: string) => {
			const params = new URLSearchParams(searchParams.toString());
			if (value) {
				params.set(key, value);
			} else {
				params.delete(key);
			}
			params.delete('page');
			router.push(`/coffees?${params.toString()}`);
		},
		[router, searchParams],
	);

	return (
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
		</div>
	);
}
