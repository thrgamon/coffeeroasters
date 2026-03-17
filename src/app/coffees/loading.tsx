import { Skeleton } from '@/components/ui/skeleton';

export default function CoffeesLoading() {
	return (
		<div className="space-y-6">
			<Skeleton className="h-9 w-48" />
			<div className="flex flex-wrap gap-3">
				<Skeleton className="h-10 w-64" />
				<Skeleton className="h-10 w-48" />
				<Skeleton className="h-10 w-40" />
				<Skeleton className="h-10 w-40" />
				<Skeleton className="h-10 w-40" />
			</div>
			<Skeleton className="h-5 w-32" />
			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{Array.from({ length: 6 }, (_, i) => `skeleton-${i}`).map((key) => (
					<div key={key} className="space-y-3 rounded-lg border p-4">
						<Skeleton className="h-5 w-3/4" />
						<Skeleton className="h-4 w-1/2" />
						<div className="flex gap-1">
							<Skeleton className="h-6 w-16 rounded-full" />
							<Skeleton className="h-6 w-20 rounded-full" />
						</div>
						<Skeleton className="h-4 w-1/3" />
					</div>
				))}
			</div>
		</div>
	);
}
