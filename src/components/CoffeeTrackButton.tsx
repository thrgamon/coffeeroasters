'use client';

import { Check, Heart, ThumbsDown, ThumbsUp, X } from 'lucide-react';
import { useState } from 'react';
import { useCoffeeTracker } from '@/lib/coffee-tracker';

type Variant = 'wishlist' | 'tried';

interface CoffeeTrackButtonProps {
	coffeeId: number;
	variant: Variant;
	size?: 'sm' | 'md';
}

export default function CoffeeTrackButton({ coffeeId, variant, size = 'sm' }: CoffeeTrackButtonProps) {
	const { hydrated, isWishlisted, isTried, getEnjoyed, toggleWishlist, markTried, removeCoffee } =
		useCoffeeTracker();
	const [showTriedOptions, setShowTriedOptions] = useState(false);

	if (variant === 'wishlist') {
		const active = hydrated && isWishlisted(coffeeId);
		const iconSize = size === 'sm' ? 'size-4' : 'size-5';

		return (
			<button
				type="button"
				onClick={(e) => {
					e.preventDefault();
					e.stopPropagation();
					toggleWishlist(coffeeId);
				}}
				aria-label={active ? 'Remove from wishlist' : 'Add to wishlist'}
				className={`inline-flex items-center gap-1 rounded-md p-1.5 text-xs transition-colors hover:bg-accent ${
					active ? 'text-red-500' : 'text-muted-foreground'
				}`}
			>
				<Heart className={`${iconSize} ${active ? 'fill-current' : ''}`} />
				{size === 'md' && <span>{active ? 'Wishlisted' : 'Wishlist'}</span>}
			</button>
		);
	}

	// Tried variant
	const tried = hydrated && isTried(coffeeId);
	const enjoyed = hydrated ? getEnjoyed(coffeeId) : undefined;
	const iconSize = size === 'sm' ? 'size-4' : 'size-5';

	if (tried) {
		return (
			<div className="inline-flex items-center gap-0.5">
				<span
					className={`inline-flex items-center gap-1 rounded-md p-1.5 text-xs ${
						enjoyed === true ? 'text-green-500' : enjoyed === false ? 'text-orange-500' : 'text-muted-foreground'
					}`}
				>
					{enjoyed === true ? (
						<ThumbsUp className={`${iconSize} fill-current`} />
					) : enjoyed === false ? (
						<ThumbsDown className={`${iconSize} fill-current`} />
					) : (
						<Check className={`${iconSize}`} />
					)}
					{size === 'md' && <span>{enjoyed === true ? 'Enjoyed' : enjoyed === false ? 'Not for me' : 'Tried'}</span>}
				</span>
				<button
					type="button"
					onClick={(e) => {
						e.preventDefault();
						e.stopPropagation();
						removeCoffee(coffeeId);
					}}
					aria-label="Remove tried status"
					className="rounded-md p-1 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
				>
					<X className="size-3" />
				</button>
			</div>
		);
	}

	if (showTriedOptions) {
		return (
			<div
				className="inline-flex items-center gap-1"
				onClick={(e) => {
					e.preventDefault();
					e.stopPropagation();
				}}
			>
				<button
					type="button"
					onClick={(e) => {
						e.preventDefault();
						e.stopPropagation();
						markTried(coffeeId, true);
						setShowTriedOptions(false);
					}}
					aria-label="Mark as enjoyed"
					className="inline-flex items-center gap-1 rounded-md p-1.5 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-green-500"
				>
					<ThumbsUp className={iconSize} />
					{size === 'md' && <span>Enjoyed</span>}
				</button>
				<button
					type="button"
					onClick={(e) => {
						e.preventDefault();
						e.stopPropagation();
						markTried(coffeeId, false);
						setShowTriedOptions(false);
					}}
					aria-label="Mark as not enjoyed"
					className="inline-flex items-center gap-1 rounded-md p-1.5 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-orange-500"
				>
					<ThumbsDown className={iconSize} />
					{size === 'md' && <span>Not for me</span>}
				</button>
				<button
					type="button"
					onClick={(e) => {
						e.preventDefault();
						e.stopPropagation();
						setShowTriedOptions(false);
					}}
					aria-label="Cancel"
					className="rounded-md p-1 text-xs text-muted-foreground transition-colors hover:bg-accent"
				>
					<X className="size-3" />
				</button>
			</div>
		);
	}

	return (
		<button
			type="button"
			onClick={(e) => {
				e.preventDefault();
				e.stopPropagation();
				setShowTriedOptions(true);
			}}
			aria-label="Mark as tried"
			className="inline-flex items-center gap-1 rounded-md p-1.5 text-xs text-muted-foreground transition-colors hover:bg-accent"
		>
			<Check className={iconSize} />
			{size === 'md' && <span>Tried it?</span>}
		</button>
	);
}
