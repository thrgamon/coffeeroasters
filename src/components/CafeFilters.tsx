'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { useCallback } from 'react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

const STATES = ['VIC', 'NSW', 'QLD', 'WA', 'SA', 'TAS', 'ACT', 'NT'];

export default function CafeFilters() {
	const router = useRouter();
	const searchParams = useSearchParams();
	const state = searchParams.get('state') ?? '';

	const updateState = useCallback(
		(value: string) => {
			const params = new URLSearchParams(searchParams.toString());
			if (value) {
				params.set('state', value);
			} else {
				params.delete('state');
			}
			router.push(`/cafes?${params.toString()}`);
		},
		[router, searchParams],
	);

	return (
		<Select value={state || 'all'} onValueChange={(v) => updateState(v === 'all' ? '' : v)}>
			<SelectTrigger className="w-40">
				<SelectValue placeholder="State" />
			</SelectTrigger>
			<SelectContent>
				<SelectItem value="all">All states</SelectItem>
				{STATES.map((s) => (
					<SelectItem key={s} value={s}>
						{s}
					</SelectItem>
				))}
			</SelectContent>
		</Select>
	);
}
