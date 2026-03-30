'use client';

import { Check, Grape, Plus, Users } from 'lucide-react';
import { useCallback, useEffect, useState } from 'react';
import { Badge } from '@/components/ui/badge';
import { useAuth } from '@/lib/auth-context';

interface CrowdsourcedNote {
	note: string;
	vote_count: number;
}

interface TastingNoteVotesResponse {
	crowdsourced_notes: CrowdsourcedNote[];
	user_votes: string[];
}

interface TastingNotesVoterProps {
	coffeeId: number;
	roasterNotes: string[];
}

const CROWDSOURCE_THRESHOLD = 3;

export default function TastingNotesVoter({ coffeeId, roasterNotes }: TastingNotesVoterProps) {
	const { user } = useAuth();
	const [crowdsourcedNotes, setCrowdsourcedNotes] = useState<CrowdsourcedNote[]>([]);
	const [userVotes, setUserVotes] = useState<Set<string>>(new Set());
	const [newNote, setNewNote] = useState('');
	const [showInput, setShowInput] = useState(false);
	const [loading, setLoading] = useState(true);

	const fetchVotes = useCallback(async () => {
		try {
			const res = await fetch(`/api/coffees/${coffeeId}/tasting-notes`, {
				credentials: 'include',
			});
			if (res.ok) {
				const data: TastingNoteVotesResponse = await res.json();
				setCrowdsourcedNotes(data.crowdsourced_notes ?? []);
				setUserVotes(new Set(data.user_votes ?? []));
			}
		} catch {
			// silently fail
		} finally {
			setLoading(false);
		}
	}, [coffeeId]);

	useEffect(() => {
		fetchVotes();
	}, [fetchVotes]);

	async function toggleVote(note: string) {
		if (!user) return;

		const wasVoted = userVotes.has(note);
		const method = wasVoted ? 'DELETE' : 'POST';

		// Optimistic update
		setUserVotes((prev) => {
			const next = new Set(prev);
			if (wasVoted) {
				next.delete(note);
			} else {
				next.add(note);
			}
			return next;
		});

		setCrowdsourcedNotes((prev) => {
			const existing = prev.find((n) => n.note === note);
			if (existing) {
				return prev
					.map((n) =>
						n.note === note
							? { ...n, vote_count: n.vote_count + (wasVoted ? -1 : 1) }
							: n
					)
					.filter((n) => n.vote_count > 0);
			}
			// New note being added
			return [...prev, { note, vote_count: 1 }];
		});

		try {
			await fetch(`/api/user/tasting-notes`, {
				method,
				headers: { 'Content-Type': 'application/json' },
				credentials: 'include',
				body: JSON.stringify({ coffee_id: coffeeId, tasting_note: note }),
			});
		} catch {
			// Revert on failure
			fetchVotes();
		}
	}

	async function handleAddNote(e: React.FormEvent) {
		e.preventDefault();
		const note = newNote.trim().toLowerCase();
		if (!note || !user) return;
		setNewNote('');
		setShowInput(false);
		await toggleVote(note);
	}

	// Build the combined view: roaster notes + crowdsourced notes that meet the threshold
	const crowdsourcedAboveThreshold = crowdsourcedNotes.filter(
		(cn) => cn.vote_count >= CROWDSOURCE_THRESHOLD && !roasterNotes.includes(cn.note)
	);

	// All crowdsourced notes (for the voting section), including ones that overlap with roaster notes
	const allVotableNotes = new Set([
		...roasterNotes.map((n) => n.toLowerCase()),
		...crowdsourcedNotes.map((n) => n.note),
	]);

	const voteCountMap = new Map(crowdsourcedNotes.map((n) => [n.note, n.vote_count]));

	return (
		<div className="space-y-4">
			{/* Roaster tasting notes */}
			{roasterNotes.length > 0 && (
				<div>
					<h3 className="mb-1 text-sm font-medium">Tasting notes</h3>
					<div className="flex flex-wrap gap-1">
						{roasterNotes.map((note) => (
							<Badge key={note} variant="secondary" className="gap-1">
								<Grape className="size-3" />
								{note}
							</Badge>
						))}
					</div>
				</div>
			)}

			{/* Crowdsourced notes that have reached the threshold */}
			{crowdsourcedAboveThreshold.length > 0 && (
				<div>
					<h3 className="mb-1 text-sm font-medium flex items-center gap-1">
						<Users className="size-3" />
						Community tasting notes
					</h3>
					<div className="flex flex-wrap gap-1">
						{crowdsourcedAboveThreshold.map((cn) => (
							<Badge key={cn.note} variant="secondary" className="gap-1">
								<Grape className="size-3" />
								{cn.note}
								<span className="ml-0.5 text-xs text-muted-foreground">({cn.vote_count})</span>
							</Badge>
						))}
					</div>
				</div>
			)}

			{/* Voting section (only for logged-in users) */}
			{user && !loading && (
				<div>
					<h3 className="mb-2 text-sm font-medium">
						What do you taste?
					</h3>
					<p className="mb-2 text-xs text-muted-foreground">
						Select the notes you agree with, or add your own
					</p>
					<div className="flex flex-wrap gap-1.5">
						{Array.from(allVotableNotes)
							.sort((a, b) => {
								const aCount = voteCountMap.get(a) ?? 0;
								const bCount = voteCountMap.get(b) ?? 0;
								return bCount - aCount || a.localeCompare(b);
							})
							.map((note) => {
								const voted = userVotes.has(note);
								const count = voteCountMap.get(note) ?? 0;
								return (
									<button
										key={note}
										type="button"
										onClick={() => toggleVote(note)}
										className={`inline-flex items-center gap-1 rounded-full border px-2.5 py-1 text-xs transition-colors ${
											voted
												? 'border-primary bg-primary/10 text-primary'
												: 'border-input text-muted-foreground hover:bg-accent hover:text-foreground'
										}`}
									>
										{voted && <Check className="size-3" />}
										<Grape className="size-3" />
										{note}
										{count > 0 && (
											<span className="text-xs opacity-60">{count}</span>
										)}
									</button>
								);
							})}

						{/* Add new note button */}
						{!showInput && (
							<button
								type="button"
								onClick={() => setShowInput(true)}
								className="inline-flex items-center gap-1 rounded-full border border-dashed border-input px-2.5 py-1 text-xs text-muted-foreground hover:bg-accent hover:text-foreground"
							>
								<Plus className="size-3" />
								Add note
							</button>
						)}
					</div>

					{/* New note input */}
					{showInput && (
						<form onSubmit={handleAddNote} className="mt-2 flex gap-2">
							<input
								type="text"
								value={newNote}
								onChange={(e) => setNewNote(e.target.value)}
								placeholder="e.g. chocolate, berry, floral..."
								maxLength={100}
								autoFocus
								className="flex-1 rounded border border-input bg-background px-3 py-1.5 text-sm"
							/>
							<button
								type="submit"
								disabled={!newNote.trim()}
								className="rounded bg-primary px-3 py-1.5 text-sm text-primary-foreground disabled:opacity-50"
							>
								Add
							</button>
							<button
								type="button"
								onClick={() => {
									setShowInput(false);
									setNewNote('');
								}}
								className="rounded px-3 py-1.5 text-sm text-muted-foreground hover:bg-accent"
							>
								Cancel
							</button>
						</form>
					)}
				</div>
			)}
		</div>
	);
}
