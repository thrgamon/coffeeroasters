import { Skeleton } from '@/components/ui/skeleton';

export default function CoffeeDetailLoading() {
	return (
		<div className="space-y-8">
			<div className="space-y-4">
				<div className="flex items-start gap-6">
					<Skeleton className="h-48 w-48 rounded-lg" />
					<div className="space-y-3">
						<Skeleton className="h-8 w-64" />
						<Skeleton className="h-5 w-40" />
						<div className="flex gap-2">
							<Skeleton className="h-6 w-20 rounded-full" />
							<Skeleton className="h-6 w-24 rounded-full" />
							<Skeleton className="h-6 w-16 rounded-full" />
						</div>
					</div>
				</div>
				<Skeleton className="h-4 w-48" />
				<Skeleton className="h-7 w-32" />
			</div>
		</div>
	);
}
