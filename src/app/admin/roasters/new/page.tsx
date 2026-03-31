'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { AdminNav } from '@/components/AdminNav';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { adminPost } from '@/lib/admin-api';
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

export default function NewRoasterPage() {
	const { user, loading } = useAuth();
	const router = useRouter();
	const [form, setForm] = useState<RoasterForm>(emptyForm);
	const [saving, setSaving] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		if (!loading && (!user || !user.is_admin)) router.push('/');
	}, [loading, user, router]);

	if (loading || !user || !user.is_admin) return null;

	const handleSave = async () => {
		setSaving(true);
		setError(null);
		try {
			await adminPost('/api/admin/roasters', form);
			router.push('/admin/roasters');
		} catch {
			setError('Failed to create roaster.');
		} finally {
			setSaving(false);
		}
	};

	const update = (field: keyof RoasterForm, value: string | boolean) => {
		setForm((prev) => ({ ...prev, [field]: value }));
	};

	return (
		<div className="space-y-6">
			<AdminNav />
			<div className="flex items-center justify-between">
				<h1 className="font-display text-4xl">New Roaster</h1>
				<Link href="/admin/roasters" className="text-sm text-muted-foreground hover:text-foreground">
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
						{saving ? 'Creating...' : 'Create'}
					</Button>
					<Link href="/admin/roasters">
						<Button variant="outline">Cancel</Button>
					</Link>
				</div>
			</div>
		</div>
	);
}
