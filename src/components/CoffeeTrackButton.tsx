'use client';

import { BookmarkPlus, Check, Eye, Heart, Star, X } from 'lucide-react';
import { useState } from 'react';
import LogCoffeeModal from '@/components/LogCoffeeModal';
import { useCoffeeTracker } from '@/lib/coffee-tracker';

type Variant = 'wishlist' | 'log';

interface CoffeeTrackButtonProps {
	coffeeId: number;
	coffeeName?: string;
	variant: Variant;
	size?: 'sm' | 'md';
}

export default function CoffeeTrackButton({
	coffeeId,
	coffeeName,
	variant,
	size = 'sm',
}: CoffeeTrackButtonProps) {
	const { hydrated, isWishlisted, isLogged, getCoffeeStatus, addToWishlist, removeFromWishlist, removeCoffee } =
		useCoffeeTracker();
	const [showLogModal, setShowLogModal] = useState(false);

	const iconSize = size === 'sm' ? 'size-4' : 'size-5';

	if (variant === 'wishlist') {
		const active = hydrated && isWishlisted(coffeeId);

		return (
			<button
				type="button"
				onClick={(e) => {
					e.preventDefault();
					e.stopPropagation();
					if (active) {
						removeFromWishlist(coffeeId);
					} else {
						addToWishlist(coffeeId);
					}
				}}
				aria-label={active ? 'Remove from watchlist' : 'Add to watchlist'}
				title={active ? 'Remove from watchlist' : 'Want to try'}
				className={`inline-flex items-center gap-1 rounded-md p-1.5 text-xs transition-colors hover:bg-secondary ${
					active ? 'text-green-500' : 'text-muted-foreground'
				}`}
			>
				<BookmarkPlus className={`${iconSize} ${active ? 'fill-current' : ''}`} />
				{size === 'md' && <span>{active ? 'On watchlist' : 'Want to try'}</span>}
			</button>
		);
	}

	// Log variant
	const logged = hydrated && isLogged(coffeeId);
	const status = hydrated ? getCoffeeStatus(coffeeId) : undefined;

	if (logged && status) {
		return (
			<>
				<div className="inline-flex items-center gap-0.5">
					<button
						type="button"
						onClick={(e) => {
							e.preventDefault();
							e.stopPropagation();
							setShowLogModal(true);
						}}
						title="Edit log"
						className="inline-flex items-center gap-1 rounded-md p-1.5 text-xs text-green-500 transition-colors hover:bg-secondary"
					>
						<Eye className={`${iconSize} fill-current`} />
						{size === 'md' && (
							<span className="flex items-center gap-1">
								Logged
								{status.liked && <Heart className="size-3 fill-red-500 text-red-500" />}
								{status.rating && (
									<span className="flex items-center text-yellow-500">
										<Star className="size-3 fill-current" />
										{status.rating}
									</span>
								)}
							</span>
						)}
					</button>
					{size === 'md' && (
						<button
							type="button"
							onClick={(e) => {
								e.preventDefault();
								e.stopPropagation();
								removeCoffee(coffeeId);
							}}
							aria-label="Remove log"
							className="rounded-md p-1 text-xs text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground"
						>
							<X className="size-3" />
						</button>
					)}
				</div>
				<LogCoffeeModal
					coffeeId={coffeeId}
					coffeeName={coffeeName}
					open={showLogModal}
					onOpenChange={setShowLogModal}
				/>
			</>
		);
	}

	return (
		<>
			<button
				type="button"
				onClick={(e) => {
					e.preventDefault();
					e.stopPropagation();
					setShowLogModal(true);
				}}
				aria-label="Log this coffee"
				title="Log this coffee"
				className="inline-flex items-center gap-1 rounded-md p-1.5 text-xs text-muted-foreground transition-colors hover:bg-secondary"
			>
				<Eye className={iconSize} />
				{size === 'md' && <span>Log</span>}
			</button>
			<LogCoffeeModal
				coffeeId={coffeeId}
				coffeeName={coffeeName}
				open={showLogModal}
				onOpenChange={setShowLogModal}
			/>
		</>
	);
}
