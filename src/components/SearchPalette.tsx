'use client';

import { Bean, MapPin, Search, Store, X } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useCallback, useEffect, useRef, useState } from 'react';

interface SearchResult {
	type: 'coffee' | 'roaster' | 'cafe';
	id: number | string;
	title: string;
	subtitle?: string;
	href: string;
}

const API = process.env.NEXT_PUBLIC_API_URL || '';

export function SearchPalette() {
	const [open, setOpen] = useState(false);
	const [query, setQuery] = useState('');
	const [results, setResults] = useState<SearchResult[]>([]);
	const [loading, setLoading] = useState(false);
	const [selected, setSelected] = useState(0);
	const inputRef = useRef<HTMLInputElement>(null);
	const router = useRouter();

	const close = useCallback(() => {
		setOpen(false);
		setQuery('');
		setResults([]);
		setSelected(0);
	}, []);

	// CMD+K / Ctrl+K to open
	useEffect(() => {
		const handler = (e: KeyboardEvent) => {
			if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
				e.preventDefault();
				setOpen((prev) => !prev);
			}
			if (e.key === 'Escape') close();
		};
		document.addEventListener('keydown', handler);
		return () => document.removeEventListener('keydown', handler);
	}, [close]);

	// Focus input when opened
	useEffect(() => {
		if (open) {
			setTimeout(() => inputRef.current?.focus(), 50);
		}
	}, [open]);

	// Search on query change
	useEffect(() => {
		if (!query.trim()) {
			setResults([]);
			return;
		}

		const controller = new AbortController();
		const timer = setTimeout(async () => {
			setLoading(true);
			try {
				const [coffeeRes, roasterRes, cafeRes] = await Promise.allSettled([
					fetch(`${API}/api/coffees?q=${encodeURIComponent(query)}&page_size=6`, {
						signal: controller.signal,
					}).then((r) => r.json()),
					fetch(`${API}/api/roasters`, { signal: controller.signal }).then((r) => r.json()),
					fetch(`${API}/api/cafes`, { signal: controller.signal }).then((r) => r.json()),
				]);

				const items: SearchResult[] = [];

				if (coffeeRes.status === 'fulfilled' && coffeeRes.value.coffees) {
					for (const c of coffeeRes.value.coffees) {
						items.push({
							type: 'coffee',
							id: c.id,
							title: c.name,
							subtitle: [c.roaster_name, c.country_name].filter(Boolean).join(' · '),
							href: `/coffees/${c.id}`,
						});
					}
				}

				const q = query.toLowerCase();
				if (roasterRes.status === 'fulfilled' && roasterRes.value.roasters) {
					for (const r of roasterRes.value.roasters) {
						if (r.name.toLowerCase().includes(q)) {
							items.push({
								type: 'roaster',
								id: r.slug,
								title: r.name,
								subtitle: r.state || undefined,
								href: `/roasters/${r.slug}`,
							});
						}
					}
				}

				if (cafeRes.status === 'fulfilled' && cafeRes.value.cafes) {
					for (const c of cafeRes.value.cafes) {
						if (c.name.toLowerCase().includes(q) || c.suburb?.toLowerCase().includes(q)) {
							items.push({
								type: 'cafe',
								id: c.id,
								title: c.name,
								subtitle: [c.suburb, c.state].filter(Boolean).join(', '),
								href: `/cafes/${c.slug}`,
							});
						}
					}
				}

				setResults(items);
				setSelected(0);
			} catch {
				// aborted or network error
			} finally {
				setLoading(false);
			}
		}, 200);

		return () => {
			clearTimeout(timer);
			controller.abort();
		};
	}, [query]);

	const navigate = useCallback(
		(href: string) => {
			close();
			router.push(href);
		},
		[close, router],
	);

	// Keyboard navigation
	const onKeyDown = useCallback(
		(e: React.KeyboardEvent) => {
			if (e.key === 'ArrowDown') {
				e.preventDefault();
				setSelected((s) => Math.min(s + 1, results.length - 1));
			} else if (e.key === 'ArrowUp') {
				e.preventDefault();
				setSelected((s) => Math.max(s - 1, 0));
			} else if (e.key === 'Enter' && results[selected]) {
				e.preventDefault();
				navigate(results[selected].href);
			}
		},
		[results, selected, navigate],
	);

	const typeIcon = (type: string) => {
		switch (type) {
			case 'coffee':
				return <Bean className="size-4 text-muted-foreground" />;
			case 'roaster':
				return <Store className="size-4 text-muted-foreground" />;
			case 'cafe':
				return <MapPin className="size-4 text-muted-foreground" />;
			default:
				return null;
		}
	};

	if (!open) return null;

	return (
		<div className="fixed inset-0 z-50">
			<button
				type="button"
				className="fixed inset-0 bg-ink/60 cursor-default"
				onClick={close}
				aria-label="Close search"
			/>
			<div className="relative mx-auto mt-[15vh] w-full max-w-lg px-4">
				<div className="border border-border bg-paper shadow-2xl">
					<div className="flex items-center gap-3 border-b border-border px-4 py-3">
						<Search className="size-5 text-muted-foreground" />
						<input
							ref={inputRef}
							type="text"
							placeholder="Search coffees, roasters, cafes..."
							value={query}
							onChange={(e) => setQuery(e.target.value)}
							onKeyDown={onKeyDown}
							className="flex-1 bg-transparent text-foreground placeholder:text-muted-foreground outline-none"
						/>
						<kbd className="hidden sm:inline-block text-xs text-muted-foreground border border-border px-1.5 py-0.5 font-mono">
							ESC
						</kbd>
						<button type="button" onClick={close} className="sm:hidden text-muted-foreground">
							<X className="size-4" />
						</button>
					</div>

					{query.trim() && (
						<div className="max-h-[50vh] overflow-y-auto">
							{loading && results.length === 0 && (
								<p className="px-4 py-8 text-center text-sm text-muted-foreground">Searching...</p>
							)}
							{!loading && results.length === 0 && query.trim() && (
								<p className="px-4 py-8 text-center text-sm text-muted-foreground">No results found</p>
							)}
							{results.map((result, i) => (
								<button
									key={`${result.type}-${result.id}`}
									type="button"
									onClick={() => navigate(result.href)}
									className={`flex w-full items-center gap-3 px-4 py-3 text-left transition-colors ${
										i === selected ? 'bg-primary/10' : 'hover:bg-muted/50'
									}`}
								>
									{typeIcon(result.type)}
									<div className="min-w-0 flex-1">
										<p className="text-sm font-medium text-foreground truncate">{result.title}</p>
										{result.subtitle && <p className="text-xs text-muted-foreground truncate">{result.subtitle}</p>}
									</div>
									<span className="text-xs text-muted-foreground capitalize">{result.type}</span>
								</button>
							))}
						</div>
					)}

					{!query.trim() && (
						<p className="px-4 py-6 text-center text-sm text-muted-foreground">
							Type to search across coffees, roasters, and cafes
						</p>
					)}
				</div>
			</div>
		</div>
	);
}

export function SearchButton() {
	return (
		<button
			type="button"
			onClick={() => document.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', metaKey: true }))}
			className="text-paper/70 hover:text-gold transition-colors"
			aria-label="Search"
		>
			<Search className="size-4" />
		</button>
	);
}
