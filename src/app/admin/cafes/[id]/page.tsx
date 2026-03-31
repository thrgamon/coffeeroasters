'use client';

import Link from 'next/link';
import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { AdminNav } from '@/components/AdminNav';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { adminFetch, adminPut } from '@/lib/admin-api';
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

export default function EditCafePage() {
	const { user, loading } = useAuth();
	const router = useRouter();
	const params = useParams();
	const id = params.id as string;
	const [form, setForm] = useState<CafeForm>(emptyForm);
	const [roasters, setRoasters] = useState<Roaster[]>([]);
	const [fetching, setFetching] = useState(true);
	const [saving, setSaving] = useState(false);
	const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

	useEffect(() => {
		if (!loading && (!user || !user.is_admin)) router.push('/');
	}, [loading, user, router]);

	useEffect(() => {
		if (!user?.is_admin) return;
		Promise.all([
			adminFetch<CafeForm & { id: number }>(`/api/admin/cafes/${id}`),
			adminFetch<Roaster[]>('/api/admin/roasters'),
		])
			.then(([cafe, roasterList]) => {
				setForm({
					roaster_id: cafe.roaster_id || 0,
					slug: cafe.slug || '',
					name: cafe.name || '',
					type: cafe.type || '',
					address: cafe.address || '',
					suburb: cafe.suburb || '',
					state: cafe.state || '',
					postcode: cafe.postcode || '',
					latitude: cafe.latitude || 0,
					longitude: cafe.longitude || 0,
					phone: cafe.phone || '',
					instagram: cafe.instagram || '',
					website_url: cafe.website_url || '',
					image_url: cafe.image_url || '',
					active: cafe.active ?? true,
				});
				setRoasters(roasterList);
			})
			.catch(() => setMessage({ type: 'error', text: 'Failed to load cafe.' }))
			.finally(() => setFetching(false));
	}, [user, id]);

	if (loading || !user || !user.is_admin) return null;

	const handleSave = async () => {
		setSaving(true);
		setMessage(null);
		try {
			await adminPut(`/api/admin/cafes/${id}`, form);
			setMessage({ type: 'success', text: 'Cafe saved.' });
		} catch {
			setMessage({ type: 'error', text: 'Failed to save cafe.' });
		} finally {
			setSaving(false);
		}
	};

	const update = (field: keyof CafeForm, value: string | number | boolean) => {
		setForm((prev) => ({ ...prev, [field]: value }));
	};

	if (fetching) {
		return (
			<div className="space-y-6">
				<AdminNav />
				<p className="text-muted-foreground text-sm">Loading...</p>
			</div>
		);
	}

	return (
		<div className="space-y-6">
			<AdminNav />
			<div className="flex items-center justify-between">
				<h1 className="font-display text-4xl">Edit Cafe</h1>
				<Link href="/admin/cafes" className="text-sm text-muted-foreground hover:text-foreground">
					Back to list
				</Link>
			</div>

			{message && (
				<div
					className={`rounded border px-4 py-2 text-sm ${message.type === 'success' ? 'border-green-600/30 bg-green-50 text-green-800' : 'border-destructive/30 bg-red-50 text-destructive'}`}
				>
					{message.text}
				</div>
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
						{saving ? 'Saving...' : 'Save'}
					</Button>
					<Link href="/admin/cafes">
						<Button variant="outline">Cancel</Button>
					</Link>
				</div>
			</div>
		</div>
	);
}
