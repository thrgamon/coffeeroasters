'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { AdminNav } from '@/components/AdminNav';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { adminFetch, adminPost } from '@/lib/admin-api';
import { useAuth } from '@/lib/auth-context';

const PROCESSES = ['washed', 'natural', 'honey', 'anaerobic', 'wet-hulled', 'experimental'];
const ROAST_LEVELS = ['light', 'medium-light', 'medium', 'medium-dark', 'dark'];

interface Roaster {
	id: number;
	name: string;
}

interface CoffeeForm {
	roaster_id: number;
	name: string;
	product_url: string;
	image_url: string;
	country_code: string;
	region_id: number;
	producer_id: number;
	process: string;
	roast_level: string;
	tasting_notes: string[];
	price_cents: number;
	weight_grams: number;
	price_per_100g_min: number;
	price_per_100g_max: number;
	variety: string;
	species: string;
	is_blend: boolean;
	is_decaf: boolean;
	in_stock: boolean;
	description: string;
}

const emptyCoffee: CoffeeForm = {
	roaster_id: 0,
	name: '',
	product_url: '',
	image_url: '',
	country_code: '',
	region_id: 0,
	producer_id: 0,
	process: '',
	roast_level: '',
	tasting_notes: [],
	price_cents: 0,
	weight_grams: 0,
	price_per_100g_min: 0,
	price_per_100g_max: 0,
	variety: '',
	species: '',
	is_blend: false,
	is_decaf: false,
	in_stock: true,
	description: '',
};

export default function NewCoffeePage() {
	const { user, loading } = useAuth();
	const router = useRouter();
	const [form, setForm] = useState<CoffeeForm>(emptyCoffee);
	const [tastingNotesText, setTastingNotesText] = useState('');
	const [roasters, setRoasters] = useState<Roaster[]>([]);
	const [saving, setSaving] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		if (!loading && (!user || !user.is_admin)) router.push('/');
	}, [loading, user, router]);

	useEffect(() => {
		if (!user?.is_admin) return;
		adminFetch<Roaster[]>('/api/admin/roasters').then(setRoasters);
	}, [user]);

	if (loading || !user || !user.is_admin) return null;

	const handleSave = async () => {
		setSaving(true);
		setError(null);
		const notes = tastingNotesText
			.split(',')
			.map((n) => n.trim())
			.filter(Boolean);
		try {
			await adminPost('/api/admin/coffees', { ...form, tasting_notes: notes });
			router.push('/admin/coffees');
		} catch {
			setError('Failed to create coffee.');
		} finally {
			setSaving(false);
		}
	};

	const update = (field: keyof CoffeeForm, value: string | number | boolean) => {
		setForm((prev) => ({ ...prev, [field]: value }));
	};

	return (
		<div className="space-y-6">
			<AdminNav />
			<div className="flex items-center justify-between">
				<h1 className="font-display text-4xl">New Coffee</h1>
				<Link href="/admin/coffees" className="text-sm text-muted-foreground hover:text-foreground">
					Back to list
				</Link>
			</div>

			{error && (
				<div className="rounded border border-destructive/30 bg-red-50 px-4 py-2 text-sm text-destructive">{error}</div>
			)}

			<div className="grid max-w-2xl gap-4">
				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Name</span>
					<Input value={form.name} onChange={(e) => update('name', e.target.value)} />
				</div>

				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Roaster</span>
					<Select
						value={form.roaster_id ? String(form.roaster_id) : ''}
						onValueChange={(v) => update('roaster_id', Number(v))}
					>
						<SelectTrigger className="w-full">
							<SelectValue placeholder="Select roaster" />
						</SelectTrigger>
						<SelectContent>
							{roasters.map((r) => (
								<SelectItem key={r.id} value={String(r.id)}>
									{r.name}
								</SelectItem>
							))}
						</SelectContent>
					</Select>
				</div>

				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Product URL</span>
					<Input value={form.product_url} onChange={(e) => update('product_url', e.target.value)} />
				</div>

				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Image URL</span>
					<Input value={form.image_url} onChange={(e) => update('image_url', e.target.value)} />
				</div>

				<div className="grid grid-cols-2 gap-4">
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Country Code</span>
						<Input value={form.country_code} onChange={(e) => update('country_code', e.target.value)} />
					</div>
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Region ID</span>
						<Input
							type="number"
							value={form.region_id || ''}
							onChange={(e) => update('region_id', Number(e.target.value))}
						/>
					</div>
				</div>

				<div className="grid grid-cols-2 gap-4">
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Producer ID</span>
						<Input
							type="number"
							value={form.producer_id || ''}
							onChange={(e) => update('producer_id', Number(e.target.value))}
						/>
					</div>
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Variety</span>
						<Input value={form.variety} onChange={(e) => update('variety', e.target.value)} />
					</div>
				</div>

				<div className="grid grid-cols-2 gap-4">
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Process</span>
						<Select value={form.process} onValueChange={(v) => update('process', v)}>
							<SelectTrigger className="w-full">
								<SelectValue placeholder="Select process" />
							</SelectTrigger>
							<SelectContent>
								{PROCESSES.map((p) => (
									<SelectItem key={p} value={p}>
										{p}
									</SelectItem>
								))}
							</SelectContent>
						</Select>
					</div>
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Roast Level</span>
						<Select value={form.roast_level} onValueChange={(v) => update('roast_level', v)}>
							<SelectTrigger className="w-full">
								<SelectValue placeholder="Select roast" />
							</SelectTrigger>
							<SelectContent>
								{ROAST_LEVELS.map((r) => (
									<SelectItem key={r} value={r}>
										{r}
									</SelectItem>
								))}
							</SelectContent>
						</Select>
					</div>
				</div>

				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Species</span>
					<Input value={form.species} onChange={(e) => update('species', e.target.value)} />
				</div>

				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
						Tasting Notes (comma-separated)
					</span>
					<Input value={tastingNotesText} onChange={(e) => setTastingNotesText(e.target.value)} />
				</div>

				<div className="grid grid-cols-2 gap-4">
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Price (cents)</span>
						<Input
							type="number"
							value={form.price_cents || ''}
							onChange={(e) => update('price_cents', Number(e.target.value))}
						/>
					</div>
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
							Weight (grams)
						</span>
						<Input
							type="number"
							value={form.weight_grams || ''}
							onChange={(e) => update('weight_grams', Number(e.target.value))}
						/>
					</div>
				</div>

				<div className="grid grid-cols-2 gap-4">
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
							Price/100g Min
						</span>
						<Input
							type="number"
							value={form.price_per_100g_min || ''}
							onChange={(e) => update('price_per_100g_min', Number(e.target.value))}
						/>
					</div>
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
							Price/100g Max
						</span>
						<Input
							type="number"
							value={form.price_per_100g_max || ''}
							onChange={(e) => update('price_per_100g_max', Number(e.target.value))}
						/>
					</div>
				</div>

				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Description</span>
					<textarea
						className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
						rows={4}
						value={form.description}
						onChange={(e) => update('description', e.target.value)}
					/>
				</div>

				<div className="flex gap-6">
					<label className="flex items-center gap-2 text-sm">
						<input
							type="checkbox"
							checked={form.is_blend}
							onChange={(e) => update('is_blend', e.target.checked)}
							className="accent-primary"
						/>
						Blend
					</label>
					<label className="flex items-center gap-2 text-sm">
						<input
							type="checkbox"
							checked={form.is_decaf}
							onChange={(e) => update('is_decaf', e.target.checked)}
							className="accent-primary"
						/>
						Decaf
					</label>
					<label className="flex items-center gap-2 text-sm">
						<input
							type="checkbox"
							checked={form.in_stock}
							onChange={(e) => update('in_stock', e.target.checked)}
							className="accent-primary"
						/>
						In Stock
					</label>
				</div>

				<div className="flex gap-3 pt-2">
					<Button onClick={handleSave} disabled={saving}>
						{saving ? 'Creating...' : 'Create'}
					</Button>
					<Link href="/admin/coffees">
						<Button variant="outline">Cancel</Button>
					</Link>
				</div>
			</div>
		</div>
	);
}
