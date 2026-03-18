'use client';

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';

interface TrackerData {
	liked: number[];
	tried: number[];
}

interface CoffeeTrackerContextType {
	hydrated: boolean;
	isLiked: (id: number) => boolean;
	isTried: (id: number) => boolean;
	toggleLike: (id: number) => void;
	toggleTried: (id: number) => void;
	likedIds: number[];
	triedIds: number[];
}

const STORAGE_KEY = 'coffee-tracker';
const EMPTY: TrackerData = { liked: [], tried: [] };

function filterNumbers(arr: unknown[]): number[] {
	return arr.filter((v): v is number => typeof v === 'number' && Number.isFinite(v));
}

function readStorage(): TrackerData {
	if (typeof window === 'undefined') return EMPTY;
	try {
		const raw = localStorage.getItem(STORAGE_KEY);
		if (!raw) return EMPTY;
		const parsed = JSON.parse(raw);
		return {
			liked: Array.isArray(parsed.liked) ? filterNumbers(parsed.liked) : [],
			tried: Array.isArray(parsed.tried) ? filterNumbers(parsed.tried) : [],
		};
	} catch {
		return EMPTY;
	}
}

function writeStorage(data: TrackerData) {
	try {
		localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
	} catch {
		// quota exceeded - state is still in memory, just won't persist
	}
}

function toggle(list: number[], id: number): number[] {
	return list.includes(id) ? list.filter((x) => x !== id) : [...list, id];
}

const CoffeeTrackerContext = createContext<CoffeeTrackerContextType | null>(null);

export function CoffeeTrackerProvider({ children }: { children: React.ReactNode }) {
	const [data, setData] = useState<TrackerData>(EMPTY);
	const [hydrated, setHydrated] = useState(false);

	useEffect(() => {
		setData(readStorage());
		setHydrated(true);
	}, []);

	const likedSet = useMemo(() => new Set(data.liked), [data.liked]);
	const triedSet = useMemo(() => new Set(data.tried), [data.tried]);

	const isLiked = useCallback((id: number) => likedSet.has(id), [likedSet]);
	const isTried = useCallback((id: number) => triedSet.has(id), [triedSet]);

	const toggleLike = useCallback((id: number) => {
		setData((prev) => {
			const next = { ...prev, liked: toggle(prev.liked, id) };
			writeStorage(next);
			return next;
		});
	}, []);

	const toggleTried = useCallback((id: number) => {
		setData((prev) => {
			const next = { ...prev, tried: toggle(prev.tried, id) };
			writeStorage(next);
			return next;
		});
	}, []);

	return (
		<CoffeeTrackerContext
			value={{ hydrated, isLiked, isTried, toggleLike, toggleTried, likedIds: data.liked, triedIds: data.tried }}
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
