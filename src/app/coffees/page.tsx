'use client';

import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { Suspense, useState } from 'react';
import CoffeeCard from '@/components/CoffeeCard';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { useGetApiCoffees, useGetApiCoffeesId } from '@/lib/api/generated/coffees/coffees';
import { useGetApiCountries } from '@/lib/api/generated/countries/countries';

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

export default function CoffeesPage() {
	return (
		<Suspense fallback={<p className="text-muted-foreground">Loading...</p>}>
			<CoffeesContent />
		</Suspense>
	);
}

function CoffeesContent() {
	const searchParams = useSearchParams();
	const similarTo = searchParams.get('similar_to');

	const [search, setSearch] = useState('');
	const [origin, setOrigin] = useState('');
	const [process, setProcess] = useState('');
	const [roast, setRoast] = useState('');
	const [variety, setVariety] = useState('');
	const [page, setPage] = useState(1);
	const pageSize = 20;

	const params: Record<string, string | number> = { page, page_size: pageSize };
	if (similarTo) {
		params.similar_to = Number(similarTo);
	} else {
		if (search) params.q = search;
		if (origin) params.origin = origin;
		if (process) params.process = process;
		if (roast) params.roast = roast;
		if (variety) params.variety = variety;
	}

	const { data, isLoading } = useGetApiCoffees(params);
	const { data: countries } = useGetApiCountries();
	const { data: sourceCoffee } = useGetApiCoffeesId(similarTo ? Number(similarTo) : 0, {
		query: { enabled: !!similarTo },
	});

	return (
		<div className="space-y-6">
			{similarTo && sourceCoffee ? (
				<div className="space-y-2">
					<h1 className="text-3xl font-bold">Coffees similar to {sourceCoffee.name}</h1>
					<p className="text-sm text-muted-foreground">
						{sourceCoffee.roaster_name}
						{sourceCoffee.country_name ? ` \u00b7 ${sourceCoffee.country_name}` : ''}
						{sourceCoffee.process ? ` \u00b7 ${sourceCoffee.process}` : ''}
					</p>
					<Link href="/coffees" className="inline-block text-sm text-primary hover:underline">
						Back to all coffees
					</Link>
				</div>
			) : (
				<h1 className="text-3xl font-bold">Coffees</h1>
			)}

			{!similarTo && (
				<div className="flex flex-wrap gap-3">
					<Input
						placeholder="Search coffees..."
						value={search}
						onChange={(e) => {
							setSearch(e.target.value);
							setPage(1);
						}}
						className="w-64"
					/>
					<Select
						value={origin}
						onValueChange={(v) => {
							setOrigin(v === 'all' ? '' : v);
							setPage(1);
						}}
					>
						<SelectTrigger className="w-48">
							<SelectValue placeholder="Country" />
						</SelectTrigger>
						<SelectContent>
							<SelectItem value="all">All countries</SelectItem>
							{countries?.countries?.map((c) => (
								<SelectItem key={c.code} value={c.code ?? ''}>
									{c.name} ({c.coffee_count})
								</SelectItem>
							))}
						</SelectContent>
					</Select>
					<Select
						value={process}
						onValueChange={(v) => {
							setProcess(v === 'all' ? '' : v);
							setPage(1);
						}}
					>
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
					<Select
						value={roast}
						onValueChange={(v) => {
							setRoast(v === 'all' ? '' : v);
							setPage(1);
						}}
					>
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
					<Select
						value={variety}
						onValueChange={(v) => {
							setVariety(v === 'all' ? '' : v);
							setPage(1);
						}}
					>
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
			)}

			{isLoading && <p className="text-muted-foreground">Loading...</p>}

			{data && (
				<>
					<p className="text-sm text-muted-foreground">
						{data.total_count ?? 0} coffee{(data.total_count ?? 0) !== 1 ? 's' : ''} found
					</p>

					<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
						{data.coffees?.map((coffee) => (
							<CoffeeCard key={coffee.id} coffee={coffee} />
						))}
					</div>

					{(data.total_count ?? 0) > pageSize && (
						<div className="flex justify-center gap-2">
							<button
								type="button"
								onClick={() => setPage((p) => Math.max(1, p - 1))}
								disabled={page <= 1}
								className="rounded border border-border px-3 py-1 text-sm disabled:opacity-50"
							>
								Previous
							</button>
							<span className="px-3 py-1 text-sm text-muted-foreground">
								Page {page} of {Math.ceil((data.total_count ?? 0) / pageSize)}
							</span>
							<button
								type="button"
								onClick={() => setPage((p) => p + 1)}
								disabled={page * pageSize >= (data.total_count ?? 0)}
								className="rounded border border-border px-3 py-1 text-sm disabled:opacity-50"
							>
								Next
							</button>
						</div>
					)}
				</>
			)}
		</div>
	);
}
