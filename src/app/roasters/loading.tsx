import { Skeleton } from '@/components/ui/skeleton';

export default function RoastersLoading() {
	return (
		<div className="space-y-6">
			<Skeleton className="h-9 w-48" />
			<Skeleton className="h-10 w-40" />
			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{Array.from({ length: 6 }, (_, i) => `skeleton-${i}`).map((key) => (
					<div key={key} className="space-y-3 rounded-lg border p-4">
						<Skeleton className="h-6 w-3/4" />
						<div className="flex gap-2">
							<Skeleton className="h-6 w-12 rounded-full" />
							<Skeleton className="h-4 w-32" />
						</div>
					</div>
				))}
			</div>
		</div>
	);
}
