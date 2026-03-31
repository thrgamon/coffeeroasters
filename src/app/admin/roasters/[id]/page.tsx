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

interface RoasterForm {
	name: string;
	slug: string;
	website: string;
	state: string;
	description: string;
	logo_url: string;
	active: boolean;
	opted_out: boolean;
}

const emptyForm: RoasterForm = {
	name: '',
	slug: '',
	website: '',
	state: '',
	description: '',
	logo_url: '',
	active: true,
	opted_out: false,
};

export default function EditRoasterPage() {
	const { user, loading } = useAuth();
	const router = useRouter();
	const params = useParams();
	const id = params.id as string;
	const [form, setForm] = useState<RoasterForm>(emptyForm);
	const [fetching, setFetching] = useState(true);
	const [saving, setSaving] = useState(false);
	const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

	useEffect(() => {
		if (!loading && (!user || !user.is_admin)) router.push('/');
	}, [loading, user, router]);

	useEffect(() => {
		if (!user?.is_admin) return;
		adminFetch<RoasterForm & { id: number }>(`/api/admin/roasters/${id}`)
			.then((data) => {
				setForm({
					name: data.name || '',
					slug: data.slug || '',
					website: data.website || '',
					state: data.state || '',
					description: data.description || '',
					logo_url: data.logo_url || '',
					active: data.active ?? true,
					opted_out: data.opted_out ?? false,
				});
			})
			.catch(() => setMessage({ type: 'error', text: 'Failed to load roaster.' }))
			.finally(() => setFetching(false));
	}, [user, id]);

	if (loading || !user || !user.is_admin) return null;

	const handleSave = async () => {
		setSaving(true);
		setMessage(null);
		try {
			await adminPut(`/api/admin/roasters/${id}`, form);
			setMessage({ type: 'success', text: 'Roaster saved.' });
		} catch {
			setMessage({ type: 'error', text: 'Failed to save roaster.' });
		} finally {
			setSaving(false);
		}
	};

	const update = (field: keyof RoasterForm, value: string | boolean) => {
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
				<h1 className="font-display text-4xl">Edit Roaster</h1>
				<Link href="/admin/roasters" className="text-sm text-muted-foreground hover:text-foreground">
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
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Website</span>
					<Input value={form.website} onChange={(e) => update('website', e.target.value)} />
				</div>
				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">State</span>
					<Select value={form.state} onValueChange={(v) => update('state', v)}>
						<SelectTrigger className="w-full">
							<SelectValue placeholder="Select state" />
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
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Description</span>
					<textarea
						className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
						rows={4}
						value={form.description}
						onChange={(e) => update('description', e.target.value)}
					/>
				</div>
				<div className="space-y-1">
					<span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Logo URL</span>
					<Input value={form.logo_url} onChange={(e) => update('logo_url', e.target.value)} />
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
					<label className="flex items-center gap-2 text-sm">
						<input
							type="checkbox"
							checked={form.opted_out}
							onChange={(e) => update('opted_out', e.target.checked)}
							className="accent-primary"
						/>
						Opted Out
					</label>
				</div>
				<div className="flex gap-3 pt-2">
					<Button onClick={handleSave} disabled={saving}>
						{saving ? 'Saving...' : 'Save'}
					</Button>
					<Link href="/admin/roasters">
						<Button variant="outline">Cancel</Button>
					</Link>
				</div>
			</div>
		</div>
	);
}
