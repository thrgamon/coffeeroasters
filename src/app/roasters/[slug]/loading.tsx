import { Skeleton } from '@/components/ui/skeleton';

export default function RoasterDetailLoading() {
	return (
		<div className="space-y-8">
			<div className="space-y-2">
				<Skeleton className="h-4 w-24" />
				<Skeleton className="h-8 w-64" />
				<div className="flex gap-3">
					<Skeleton className="h-6 w-12 rounded-full" />
					<Skeleton className="h-4 w-48" />
				</div>
			</div>
			<div className="space-y-4">
				<Skeleton className="h-6 w-32" />
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{Array.from({ length: 3 }, (_, i) => `skeleton-${i}`).map((key) => (
						<div key={key} className="space-y-3 rounded-lg border p-4">
							<Skeleton className="h-5 w-3/4" />
							<Skeleton className="h-4 w-1/2" />
							<div className="flex gap-1">
								<Skeleton className="h-6 w-16 rounded-full" />
								<Skeleton className="h-6 w-20 rounded-full" />
							</div>
						</div>
					))}
				</div>
			</div>
		</div>
	);
}
