'use client';

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
	enjoyed?: boolean;
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
	const enjoyed = coffees.filter((c) => c.status === 'tried' && c.enjoyed === true);
	const notEnjoyed = coffees.filter((c) => c.status === 'tried' && c.enjoyed === false);
	const triedNoRating = coffees.filter(
		(c) => c.status === 'tried' && c.enjoyed === undefined && c.enjoyed !== true && c.enjoyed !== false,
	);

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
			<p className="text-muted-foreground">
				You haven't saved any coffees yet. Use the heart icon to add coffees to your wishlist, or mark
				them as tried.
			</p>
		);
	}

	return (
		<>
			{wishlist.length > 0 && (
				<CoffeeSection title="Wishlist" coffees={wishlist} />
			)}
			{enjoyed.length > 0 && (
				<CoffeeSection title="Enjoyed" coffees={enjoyed} />
			)}
			{notEnjoyed.length > 0 && (
				<CoffeeSection title="Not for me" coffees={notEnjoyed} />
			)}
			{triedNoRating.length > 0 && (
				<CoffeeSection title="Tried" coffees={triedNoRating} />
			)}
			<Recommendations />
		</>
	);
}

function CoffeeSection({ title, coffees }: { title: string; coffees: DomainCoffeeResponse[] }) {
	return (
		<section className="space-y-3">
			<h2 className="text-xl font-semibold">
				{title} <span className="text-sm font-normal text-muted-foreground">({coffees.length})</span>
			</h2>
			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{coffees.map((coffee) => (
					<CoffeeCard key={coffee.id} coffee={coffee} />
				))}
			</div>
		</section>
	);
}

function LocalCoffeeList() {
	const { wishlistIds, triedIds } = useCoffeeTracker();
	const empty = wishlistIds.length === 0 && triedIds.length === 0;

	if (empty) {
		return (
			<div className="space-y-4">
				<p className="text-muted-foreground">
					You haven't saved any coffees yet. Use the heart icon to add coffees to your wishlist, or
					mark them as tried.
				</p>
				<p className="text-sm text-muted-foreground">
					<Link href="/login" className="text-primary underline">
						Sign in
					</Link>{' '}
					to save your coffees across devices.
				</p>
			</div>
		);
	}

	return (
		<>
			<p className="text-sm text-muted-foreground">
				<Link href="/login" className="text-primary underline">
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
			<h1 className="text-3xl font-bold">My coffees</h1>
			{user ? <ServerCoffeeList /> : <LocalCoffeeList />}
		</div>
	);
}
