'use client';

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { useAuth } from '@/lib/auth-context';

interface CoffeeStatus {
	status: 'wishlist' | 'tried';
	enjoyed?: boolean;
}

interface TrackerData {
	coffees: Record<number, CoffeeStatus>;
}

interface CoffeeTrackerContextType {
	hydrated: boolean;
	isWishlisted: (id: number) => boolean;
	isTried: (id: number) => boolean;
	getEnjoyed: (id: number) => boolean | undefined;
	toggleWishlist: (id: number) => void;
	markTried: (id: number, enjoyed: boolean) => void;
	removeCoffee: (id: number) => void;
	wishlistIds: number[];
	triedIds: number[];
	getCoffeeStatus: (id: number) => CoffeeStatus | undefined;
}

const STORAGE_KEY = 'coffee-tracker';
const EMPTY: TrackerData = { coffees: {} };

function readStorage(): TrackerData {
	if (typeof window === 'undefined') return EMPTY;
	try {
		const raw = localStorage.getItem(STORAGE_KEY);
		if (!raw) return EMPTY;
		const parsed = JSON.parse(raw);
		// Migrate old format (liked/tried arrays) to new format
		if (Array.isArray(parsed.liked) || Array.isArray(parsed.tried)) {
			const coffees: Record<number, CoffeeStatus> = {};
			if (Array.isArray(parsed.liked)) {
				for (const id of parsed.liked) {
					if (typeof id === 'number') coffees[id] = { status: 'wishlist' };
				}
			}
			if (Array.isArray(parsed.tried)) {
				for (const id of parsed.tried) {
					if (typeof id === 'number') coffees[id] = { status: 'tried' };
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

	// Load data from server when authenticated, otherwise from localStorage
	useEffect(() => {
		if (user) {
			fetch('/api/user/coffee-ids', { credentials: 'include' })
				.then(async (res) => {
					if (res.ok) {
						const items: { coffee_id: number; status: string; enjoyed?: boolean }[] = await res.json();
						const coffees: Record<number, CoffeeStatus> = {};
						for (const item of items) {
							coffees[item.coffee_id] = {
								status: item.status as 'wishlist' | 'tried',
								enjoyed: item.enjoyed,
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

	const triedIds = useMemo(
		() =>
			Object.entries(data.coffees)
				.filter(([, v]) => v.status === 'tried')
				.map(([k]) => Number(k)),
		[data.coffees],
	);

	const isWishlisted = useCallback((id: number) => data.coffees[id]?.status === 'wishlist', [data.coffees]);
	const isTried = useCallback((id: number) => data.coffees[id]?.status === 'tried', [data.coffees]);
	const getEnjoyed = useCallback((id: number) => data.coffees[id]?.enjoyed, [data.coffees]);
	const getCoffeeStatus = useCallback((id: number) => data.coffees[id], [data.coffees]);

	const syncToServer = useCallback(
		async (coffeeId: number, status: string, enjoyed?: boolean) => {
			if (!user) return;
			await fetch('/api/user/coffees', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'include',
				body: JSON.stringify({ coffee_id: coffeeId, status, enjoyed }),
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

	const toggleWishlist = useCallback(
		(id: number) => {
			setData((prev) => {
				const next = { ...prev, coffees: { ...prev.coffees } };
				if (prev.coffees[id]?.status === 'wishlist') {
					delete next.coffees[id];
					deleteFromServer(id);
				} else {
					next.coffees[id] = { status: 'wishlist' };
					syncToServer(id, 'wishlist');
				}
				if (!user) writeStorage(next);
				return next;
			});
		},
		[user, syncToServer, deleteFromServer],
	);

	const markTried = useCallback(
		(id: number, enjoyed: boolean) => {
			setData((prev) => {
				const next = { ...prev, coffees: { ...prev.coffees } };
				next.coffees[id] = { status: 'tried', enjoyed };
				if (!user) writeStorage(next);
				syncToServer(id, 'tried', enjoyed);
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
				isTried,
				getEnjoyed,
				toggleWishlist,
				markTried,
				removeCoffee,
				wishlistIds,
				triedIds,
				getCoffeeStatus,
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
