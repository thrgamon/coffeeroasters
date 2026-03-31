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

const STATES = ['VIC', 'NSW', 'QLD', 'WA', 'SA', 'TAS', 'ACT', 'NT'];
const TYPES = ['owned', 'stockist'];

interface Roaster {
	id: number;
	name: string;
}

interface CafeForm {
	roaster_id: number;
	slug: string;
	name: string;
	type: string;
	address: string;
	suburb: string;
	state: string;
	postcode: string;
	latitude: number;
	longitude: number;
	phone: string;
	instagram: string;
	website_url: string;
	image_url: string;
	active: boolean;
}

const emptyForm: CafeForm = {
	roaster_id: 0,
	slug: '',
	name: '',
	type: '',
	address: '',
	suburb: '',
	state: '',
	postcode: '',
	latitude: 0,
	longitude: 0,
	phone: '',
	instagram: '',
	website_url: '',
	image_url: '',
	active: true,
};

export default function NewCafePage() {
	const { user, loading } = useAuth();
	const router = useRouter();
	const [form, setForm] = useState<CafeForm>(emptyForm);
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
		try {
			await adminPost('/api/admin/cafes', form);
			router.push('/admin/cafes');
		} catch {
			setError('Failed to create cafe.');
		} finally {
			setSaving(false);
		}
	};

	const update = (field: keyof CafeForm, value: string | number | boolean) => {
		setForm((prev) => ({ ...prev, [field]: value }));
	};

	return (
		<div className="space-y-6">
			<AdminNav />
			<div className="flex items-center justify-between">
				<h1 className="font-display text-4xl">New Cafe</h1>
				<Link href="/admin/cafes" className="text-sm text-muted-foreground hover:text-foreground">
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
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Slug</span>
					<Input value={form.slug} onChange={(e) => update('slug', e.target.value)} />
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
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Type</span>
					<Select value={form.type} onValueChange={(v) => update('type', v)}>
						<SelectTrigger className="w-full">
							<SelectValue placeholder="Select type" />
						</SelectTrigger>
						<SelectContent>
							{TYPES.map((t) => (
								<SelectItem key={t} value={t}>
									{t}
								</SelectItem>
							))}
						</SelectContent>
					</Select>
				</div>

				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Address</span>
					<Input value={form.address} onChange={(e) => update('address', e.target.value)} />
				</div>

				<div className="grid grid-cols-3 gap-4">
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Suburb</span>
						<Input value={form.suburb} onChange={(e) => update('suburb', e.target.value)} />
					</div>
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">State</span>
						<Select value={form.state} onValueChange={(v) => update('state', v)}>
							<SelectTrigger className="w-full">
								<SelectValue placeholder="State" />
							</SelectTrigger>
							<SelectContent>
								{STATES.map((s) => (
									<SelectItem key={s} value={s}>
										{s}
									</SelectItem>
								))}
							</SelectContent>
						</Select>
					</div>
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Postcode</span>
						<Input value={form.postcode} onChange={(e) => update('postcode', e.target.value)} />
					</div>
				</div>

				<div className="grid grid-cols-2 gap-4">
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Latitude</span>
						<Input
							type="number"
							step="any"
							value={form.latitude || ''}
							onChange={(e) => update('latitude', Number(e.target.value))}
						/>
					</div>
					<div className="space-y-1">
						<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Longitude</span>
						<Input
							type="number"
							step="any"
							value={form.longitude || ''}
							onChange={(e) => update('longitude', Number(e.target.value))}
						/>
					</div>
				</div>

				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Phone</span>
					<Input value={form.phone} onChange={(e) => update('phone', e.target.value)} />
				</div>
				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Instagram</span>
					<Input value={form.instagram} onChange={(e) => update('instagram', e.target.value)} />
				</div>
				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Website URL</span>
					<Input value={form.website_url} onChange={(e) => update('website_url', e.target.value)} />
				</div>
				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Image URL</span>
					<Input value={form.image_url} onChange={(e) => update('image_url', e.target.value)} />
				</div>

				<div className="flex gap-6">
					<label className="flex items-center gap-2 text-sm">
						<input
							type="checkbox"
							checked={form.active}
							onChange={(e) => update('active', e.target.checked)}
							className="accent-primary"
						/>
						Active
					</label>
				</div>

				<div className="flex gap-3 pt-2">
					<Button onClick={handleSave} disabled={saving}>
						{saving ? 'Creating...' : 'Create'}
					</Button>
					<Link href="/admin/cafes">
						<Button variant="outline">Cancel</Button>
					</Link>
				</div>
			</div>
		</div>
	);
}
