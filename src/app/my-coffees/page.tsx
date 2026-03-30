'use client';

import { Heart, Star } from 'lucide-react';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import CoffeeCard from '@/components/CoffeeCard';
import Recommendations from '@/components/Recommendations';
import { Skeleton } from '@/components/ui/skeleton';
import type { DomainCoffeeResponse } from '@/lib/api/generated/models';
import { useAuth } from '@/lib/auth-context';
import { useCoffeeTracker } from '@/lib/coffee-tracker';

interface UserCoffeeDetail extends DomainCoffeeResponse {
	status: string;
	liked?: boolean;
	rating?: number;
	review?: string;
	drunk_at?: string;
}

function ServerCoffeeList() {
	const [coffees, setCoffees] = useState<UserCoffeeDetail[]>([]);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		fetch('/api/user/coffees', { credentials: 'include' })
			.then(async (res) => {
				if (res.ok) {
					const data = await res.json();
					setCoffees(data.coffees ?? []);
				}
			})
			.catch(() => {})
			.finally(() => setLoading(false));
	}, []);

	const wishlist = coffees.filter((c) => c.status === 'wishlist');
	const logged = coffees.filter((c) => c.status === 'logged');

	if (loading) {
		return (
			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{[1, 2, 3].map((i) => (
					<Skeleton key={i} className="h-48 rounded-lg" />
				))}
			</div>
		);
	}

	if (coffees.length === 0) {
		return (
			<p className="text-grey-olive">
				You haven't saved any coffees yet. Use the bookmark icon to add coffees to your watchlist, or
				the eye icon to log coffees you've tried.
			</p>
		);
	}

	return (
		<>
			{wishlist.length > 0 && <CoffeeSection title="Watchlist" coffees={wishlist} />}
			{logged.length > 0 && <LoggedCoffeeSection coffees={logged} />}
			<Recommendations />
		</>
	);
}

function CoffeeSection({ title, coffees }: { title: string; coffees: DomainCoffeeResponse[] }) {
	return (
		<section className="space-y-3">
			<h2 className="text-xl font-bold uppercase tracking-wider text-dusty-rose">
				{title} <span className="text-sm font-normal text-grey-olive font-mono">({coffees.length})</span>
			</h2>
			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{coffees.map((coffee) => (
					<CoffeeCard key={coffee.id} coffee={coffee} />
				))}
			</div>
		</section>
	);
}

function LoggedCoffeeSection({ coffees }: { coffees: UserCoffeeDetail[] }) {
	return (
		<section className="space-y-3">
			<h2 className="text-xl font-bold uppercase tracking-wider text-dusty-rose">
				Diary <span className="text-sm font-normal text-grey-olive font-mono">({coffees.length})</span>
			</h2>
			<div className="space-y-3">
				{coffees.map((coffee) => (
					<LoggedCoffeeEntry key={coffee.id} coffee={coffee} />
				))}
			</div>
		</section>
	);
}

function LoggedCoffeeEntry({ coffee }: { coffee: UserCoffeeDetail }) {
	return (
		<Link
			href={`/coffees/${coffee.id}`}
			className="flex items-center gap-4 rounded border border-border/50 bg-card/80 p-4 transition-all hover:border-gold/30"
		>
			<div className="flex-1 min-w-0">
				<div className="flex items-center gap-2">
					<span className="font-medium text-snow truncate">{coffee.name}</span>
					{coffee.liked && <Heart className="size-4 shrink-0 fill-red-500 text-red-500" />}
				</div>
				{coffee.roaster_name && (
					<p className="text-sm text-dusty-rose">{coffee.roaster_name}</p>
				)}
				{coffee.review && (
					<p className="mt-1 text-sm text-grey-olive line-clamp-2">{coffee.review}</p>
				)}
			</div>
			<div className="flex flex-col items-end gap-1 shrink-0">
				{coffee.rating != null && coffee.rating > 0 && (
					<div className="flex items-center gap-0.5">
						{[1, 2, 3, 4, 5].map((n) => (
							<Star
								key={n}
								className={`size-4 ${n <= coffee.rating! ? 'fill-yellow-400 text-yellow-400' : 'text-grey-olive/30'}`}
							/>
						))}
					</div>
				)}
				{coffee.drunk_at && (
					<span className="text-xs text-grey-olive font-mono">{coffee.drunk_at}</span>
				)}
			</div>
		</Link>
	);
}

function LocalCoffeeList() {
	const { wishlistIds, loggedIds } = useCoffeeTracker();
	const empty = wishlistIds.length === 0 && loggedIds.length === 0;

	if (empty) {
		return (
			<div className="space-y-4">
				<p className="text-grey-olive">
					You haven't saved any coffees yet. Use the bookmark icon to add coffees to your watchlist,
					or the eye icon to log coffees you've tried.
				</p>
				<p className="text-sm text-grey-olive">
					<Link href="/login" className="text-gold underline">
						Sign in
					</Link>{' '}
					to save your coffees across devices.
				</p>
			</div>
		);
	}

	return (
		<>
			<p className="text-sm text-grey-olive">
				<Link href="/login" className="text-gold underline">
					Sign in
				</Link>{' '}
				to save your coffees across devices.
			</p>
			<Recommendations />
		</>
	);
}

export default function MyCoffeesPage() {
	const { user, loading } = useAuth();

	if (loading) return null;

	return (
		<div className="space-y-8">
			<h1 className="text-3xl font-bold uppercase tracking-wider text-snow">My coffees</h1>
			{user ? <ServerCoffeeList /> : <LocalCoffeeList />}
		</div>
	);
}
