'use client';

import { Check, Heart } from 'lucide-react';
import { useCoffeeTracker } from '@/lib/coffee-tracker';

type Variant = 'like' | 'tried';

interface CoffeeTrackButtonProps {
	coffeeId: number;
	variant: Variant;
	size?: 'sm' | 'md';
}

export default function CoffeeTrackButton({ coffeeId, variant, size = 'sm' }: CoffeeTrackButtonProps) {
	const { hydrated, isLiked, isTried, toggleLike, toggleTried } = useCoffeeTracker();

	const active = hydrated && (variant === 'like' ? isLiked(coffeeId) : isTried(coffeeId));
	const toggle = variant === 'like' ? toggleLike : toggleTried;
	const Icon = variant === 'like' ? Heart : Check;
	const label = variant === 'like' ? 'Like' : 'Tried';

	const iconSize = size === 'sm' ? 'size-4' : 'size-5';

	return (
		<button
			type="button"
			onClick={(e) => {
				e.preventDefault();
				e.stopPropagation();
				toggle(coffeeId);
			}}
			aria-label={`${active ? 'Remove from' : 'Add to'} ${label.toLowerCase()}`}
			className={`inline-flex items-center gap-1 rounded-md p-1.5 text-xs transition-colors hover:bg-accent ${
				active ? (variant === 'like' ? 'text-red-500' : 'text-green-500') : 'text-muted-foreground'
			}`}
		>
			<Icon className={`${iconSize} ${active ? 'fill-current' : ''}`} />
			{size === 'md' && <span>{label}</span>}
		</button>
	);
}
