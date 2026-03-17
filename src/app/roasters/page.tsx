'use client';

import Link from 'next/link';
import { useState } from 'react';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { useGetApiRoasters } from '@/lib/api/generated/roasters/roasters';

const STATES = ['VIC', 'NSW', 'QLD', 'WA', 'SA', 'TAS', 'ACT', 'NT'];

export default function RoastersPage() {
	const [state, setState] = useState('');

	const params: Record<string, string> = {};
	if (state) params.state = state;

	const { data, isLoading } = useGetApiRoasters(params);

	return (
		<div className="space-y-6">
			<h1 className="text-3xl font-bold">Roasters</h1>

			<Select value={state} onValueChange={(v) => setState(v === 'all' ? '' : v)}>
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

			{isLoading && <p className="text-muted-foreground">Loading...</p>}

			{data && (
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{data.roasters?.map((roaster) => (
						<Link key={roaster.id} href={`/roasters/${roaster.slug}`}>
							<Card className="transition-colors hover:border-primary/50">
								<CardHeader>
									<CardTitle className="text-lg">{roaster.name}</CardTitle>
								</CardHeader>
								<CardContent className="flex items-center gap-2">
									{roaster.state && <Badge variant="secondary">{roaster.state}</Badge>}
									{roaster.website && (
										<a
											href={roaster.website}
											target="_blank"
											rel="noopener noreferrer"
											className="text-sm text-muted-foreground hover:text-primary"
											onClick={(e) => e.stopPropagation()}
										>
											{new URL(roaster.website).hostname}
										</a>
									)}
								</CardContent>
							</Card>
						</Link>
					))}
				</div>
			)}

			{data && (data.roasters?.length ?? 0) === 0 && <p className="text-muted-foreground">No roasters found.</p>}
		</div>
	);
}
