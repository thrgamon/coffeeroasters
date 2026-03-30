'use client';

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { useAuth } from '@/lib/auth-context';

export interface CoffeeStatus {
	status: 'wishlist' | 'logged';
	liked?: boolean;
	rating?: number;
	review?: string;
	drunkAt?: string;
}

interface TrackerData {
	coffees: Record<number, CoffeeStatus>;
}

interface CoffeeTrackerContextType {
	hydrated: boolean;
	isWishlisted: (id: number) => boolean;
	isLogged: (id: number) => boolean;
	getCoffeeStatus: (id: number) => CoffeeStatus | undefined;
	addToWishlist: (id: number) => void;
	removeFromWishlist: (id: number) => void;
	logCoffee: (id: number, data: { liked?: boolean; rating?: number; review?: string; drunkAt?: string }) => void;
	removeCoffee: (id: number) => void;
	wishlistIds: number[];
	loggedIds: number[];
}

const STORAGE_KEY = 'coffee-tracker';
const EMPTY: TrackerData = { coffees: {} };

function readStorage(): TrackerData {
	if (typeof window === 'undefined') return EMPTY;
	try {
		const raw = localStorage.getItem(STORAGE_KEY);
		if (!raw) return EMPTY;
		const parsed = JSON.parse(raw);
		// Migrate old format
		if (Array.isArray(parsed.liked) || Array.isArray(parsed.tried)) {
			const coffees: Record<number, CoffeeStatus> = {};
			if (Array.isArray(parsed.liked)) {
				for (const id of parsed.liked) {
					if (typeof id === 'number') coffees[id] = { status: 'wishlist' };
				}
			}
			if (Array.isArray(parsed.tried)) {
				for (const id of parsed.tried) {
					if (typeof id === 'number') coffees[id] = { status: 'logged' };
				}
			}
			return { coffees };
		}
		return { coffees: parsed.coffees ?? {} };
	} catch {
		return EMPTY;
	}
}

function writeStorage(data: TrackerData) {
	try {
		localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
	} catch {
		// quota exceeded
	}
}

const CoffeeTrackerContext = createContext<CoffeeTrackerContextType | null>(null);

export function CoffeeTrackerProvider({ children }: { children: React.ReactNode }) {
	const { user } = useAuth();
	const [data, setData] = useState<TrackerData>(EMPTY);
	const [hydrated, setHydrated] = useState(false);

	useEffect(() => {
		if (user) {
			fetch('/api/user/coffee-ids', { credentials: 'include' })
				.then(async (res) => {
					if (res.ok) {
						const items: { coffee_id: number; status: string; liked?: boolean; rating?: number }[] = await res.json();
						const coffees: Record<number, CoffeeStatus> = {};
						for (const item of items) {
							coffees[item.coffee_id] = {
								status: item.status as 'wishlist' | 'logged',
								liked: item.liked,
								rating: item.rating,
							};
						}
						setData({ coffees });
					}
				})
				.catch(() => {})
				.finally(() => setHydrated(true));
		} else {
			setData(readStorage());
			setHydrated(true);
		}
	}, [user]);

	const wishlistIds = useMemo(
		() =>
			Object.entries(data.coffees)
				.filter(([, v]) => v.status === 'wishlist')
				.map(([k]) => Number(k)),
		[data.coffees],
	);

	const loggedIds = useMemo(
		() =>
			Object.entries(data.coffees)
				.filter(([, v]) => v.status === 'logged')
				.map(([k]) => Number(k)),
		[data.coffees],
	);

	const isWishlisted = useCallback((id: number) => data.coffees[id]?.status === 'wishlist', [data.coffees]);
	const isLogged = useCallback((id: number) => data.coffees[id]?.status === 'logged', [data.coffees]);
	const getCoffeeStatus = useCallback((id: number) => data.coffees[id], [data.coffees]);

	const syncToServer = useCallback(
		async (
			coffeeId: number,
			status: string,
			extra?: { liked?: boolean; rating?: number; review?: string; drunkAt?: string },
		) => {
			if (!user) return;
			await fetch('/api/user/coffees', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'include',
				body: JSON.stringify({
					coffee_id: coffeeId,
					status,
					liked: extra?.liked,
					rating: extra?.rating,
					review: extra?.review,
					drunk_at: extra?.drunkAt,
				}),
			}).catch(() => {});
		},
		[user],
	);

	const deleteFromServer = useCallback(
		async (coffeeId: number) => {
			if (!user) return;
			await fetch(`/api/user/coffees/${coffeeId}`, {
				method: 'DELETE',
				credentials: 'include',
			}).catch(() => {});
		},
		[user],
	);

	const addToWishlist = useCallback(
		(id: number) => {
			setData((prev) => {
				const next = { ...prev, coffees: { ...prev.coffees, [id]: { status: 'wishlist' as const } } };
				if (!user) writeStorage(next);
				syncToServer(id, 'wishlist');
				return next;
			});
		},
		[user, syncToServer],
	);

	const removeFromWishlist = useCallback(
		(id: number) => {
			setData((prev) => {
				const next = { ...prev, coffees: { ...prev.coffees } };
				delete next.coffees[id];
				if (!user) writeStorage(next);
				deleteFromServer(id);
				return next;
			});
		},
		[user, deleteFromServer],
	);

	const logCoffee = useCallback(
		(id: number, logData: { liked?: boolean; rating?: number; review?: string; drunkAt?: string }) => {
			setData((prev) => {
				const next = {
					...prev,
					coffees: {
						...prev.coffees,
						[id]: { status: 'logged' as const, ...logData },
					},
				};
				if (!user) writeStorage(next);
				syncToServer(id, 'logged', logData);
				return next;
			});
		},
		[user, syncToServer],
	);

	const removeCoffee = useCallback(
		(id: number) => {
			setData((prev) => {
				const next = { ...prev, coffees: { ...prev.coffees } };
				delete next.coffees[id];
				if (!user) writeStorage(next);
				deleteFromServer(id);
				return next;
			});
		},
		[user, deleteFromServer],
	);

	return (
		<CoffeeTrackerContext
			value={{
				hydrated,
				isWishlisted,
				isLogged,
				getCoffeeStatus,
				addToWishlist,
				removeFromWishlist,
				logCoffee,
				removeCoffee,
				wishlistIds,
				loggedIds,
			}}
		>
			{children}
		</CoffeeTrackerContext>
	);
}

export function useCoffeeTracker(): CoffeeTrackerContextType {
	const ctx = useContext(CoffeeTrackerContext);
	if (!ctx) throw new Error('useCoffeeTracker must be used within CoffeeTrackerProvider');
	return ctx;
}
