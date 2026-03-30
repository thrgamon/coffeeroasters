'use client';

import { Heart, Star } from 'lucide-react';
import { useState } from 'react';
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/sheet';
import { useCoffeeTracker } from '@/lib/coffee-tracker';

interface LogCoffeeModalProps {
	coffeeId: number;
	coffeeName?: string;
	open: boolean;
	onOpenChange: (open: boolean) => void;
}

export default function LogCoffeeModal({ coffeeId, coffeeName, open, onOpenChange }: LogCoffeeModalProps) {
	const { getCoffeeStatus, logCoffee } = useCoffeeTracker();
	const existing = getCoffeeStatus(coffeeId);

	const [liked, setLiked] = useState<boolean | undefined>(existing?.liked);
	const [rating, setRating] = useState<number>(existing?.rating ?? 0);
	const [hoverRating, setHoverRating] = useState(0);
	const [review, setReview] = useState(existing?.review ?? '');
	const [drunkAt, setDrunkAt] = useState(existing?.drunkAt ?? new Date().toISOString().slice(0, 10));

	function handleSave() {
		logCoffee(coffeeId, {
			liked: liked,
			rating: rating > 0 ? rating : undefined,
			review: review.trim() || undefined,
			drunkAt: drunkAt || undefined,
		});
		onOpenChange(false);
	}

	return (
		<Sheet open={open} onOpenChange={onOpenChange}>
			<SheetContent side="bottom" className="mx-auto max-w-lg rounded-t-xl">
				<SheetHeader>
					<SheetTitle>{coffeeName ? `Log "${coffeeName}"` : 'Log this coffee'}</SheetTitle>
				</SheetHeader>
				<div className="space-y-5 px-4 pb-6">
					{/* Date */}
					<div>
						<label htmlFor="drunk-at" className="mb-1 block text-sm font-medium">
							Date
						</label>
						<input
							id="drunk-at"
							type="date"
							value={drunkAt}
							onChange={(e) => setDrunkAt(e.target.value)}
							className="w-full rounded border border-input bg-background px-3 py-2 text-sm"
						/>
					</div>

					{/* Rating */}
					<div>
						<span className="mb-2 block text-sm font-medium">Rating</span>
						<div className="flex gap-1">
							{[1, 2, 3, 4, 5].map((n) => (
								<button
									key={n}
									type="button"
									onClick={() => setRating(rating === n ? 0 : n)}
									onMouseEnter={() => setHoverRating(n)}
									onMouseLeave={() => setHoverRating(0)}
									className="p-0.5 transition-colors"
									aria-label={`${n} star${n > 1 ? 's' : ''}`}
								>
									<Star
										className={`size-7 ${
											n <= (hoverRating || rating) ? 'fill-yellow-400 text-yellow-400' : 'text-muted-foreground/40'
										}`}
									/>
								</button>
							))}
							{rating > 0 && <span className="ml-2 self-center text-sm text-muted-foreground">{rating}/5</span>}
						</div>
					</div>

					{/* Liked */}
					<div>
						<span className="mb-2 block text-sm font-medium">Liked?</span>
						<button
							type="button"
							onClick={() => setLiked(liked === true ? undefined : true)}
							className={`inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm transition-colors ${
								liked === true ? 'bg-red-50 text-red-500 dark:bg-red-950' : 'text-muted-foreground hover:bg-accent'
							}`}
						>
							<Heart className={`size-4 ${liked === true ? 'fill-current' : ''}`} />
							{liked === true ? 'Liked' : 'Like'}
						</button>
					</div>

					{/* Review */}
					<div>
						<label htmlFor="review" className="mb-1 block text-sm font-medium">
							Review <span className="font-normal text-muted-foreground">(optional)</span>
						</label>
						<textarea
							id="review"
							value={review}
							onChange={(e) => setReview(e.target.value)}
							placeholder="What did you think?"
							rows={3}
							className="w-full rounded border border-input bg-background px-3 py-2 text-sm"
						/>
					</div>

					{/* Actions */}
					<div className="flex gap-3">
						<button
							type="button"
							onClick={handleSave}
							className="flex-1 rounded bg-primary px-4 py-2 text-sm font-medium text-primary-foreground"
						>
							Save
						</button>
						<button
							type="button"
							onClick={() => onOpenChange(false)}
							className="rounded px-4 py-2 text-sm text-muted-foreground hover:bg-accent"
						>
							Cancel
						</button>
					</div>
				</div>
			</SheetContent>
		</Sheet>
	);
}
