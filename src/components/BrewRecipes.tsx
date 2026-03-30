'use client';

import { Clock, Coffee, Droplets, Flame, Plus, Thermometer, Trash2 } from 'lucide-react';
import { useCallback, useEffect, useState } from 'react';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent } from '@/components/ui/card';
import type { DomainBrewRecipeListResponse, DomainBrewRecipeResponse } from '@/lib/api/generated/models';

const BREW_METHODS: Record<string, string> = {
	espresso: 'Espresso',
	pourover: 'Pour Over',
	aeropress: 'AeroPress',
	french_press: 'French Press',
	cold_brew: 'Cold Brew',
	filter: 'Filter',
	moka_pot: 'Moka Pot',
	other: 'Other',
};

function formatBrewTime(seconds: number): string {
	const mins = Math.floor(seconds / 60);
	const secs = seconds % 60;
	if (mins === 0) return `${secs}s`;
	if (secs === 0) return `${mins}m`;
	return `${mins}m ${secs}s`;
}

interface BrewRecipesProps {
	coffeeId: number;
	brewRecipeRaw?: string;
}

export default function BrewRecipes({ coffeeId, brewRecipeRaw }: BrewRecipesProps) {
	const [recipes, setRecipes] = useState<DomainBrewRecipeResponse[]>([]);
	const [showForm, setShowForm] = useState(false);
	const [isLoggedIn, setIsLoggedIn] = useState(false);
	const [currentUserId, setCurrentUserId] = useState<number | null>(null);
	const [loading, setLoading] = useState(true);

	const fetchRecipes = useCallback(async () => {
		try {
			const res = await fetch(`/api/coffees/${coffeeId}/recipes`, { credentials: 'include' });
			if (res.ok) {
				const data: DomainBrewRecipeListResponse = await res.json();
				setRecipes(data.recipes ?? []);
			}
		} finally {
			setLoading(false);
		}
	}, [coffeeId]);

	useEffect(() => {
		fetchRecipes();
		fetch('/api/auth/me', { credentials: 'include' })
			.then((r) => r.json())
			.then((data) => {
				if (data.user) {
					setIsLoggedIn(true);
					setCurrentUserId(data.user.id);
				}
			})
			.catch(() => {});
	}, [fetchRecipes]);

	const handleDelete = async (id: number) => {
		const res = await fetch(`/api/recipes/${id}`, { method: 'DELETE', credentials: 'include' });
		if (res.ok) {
			setRecipes((prev) => prev.filter((r) => r.id !== id));
		}
	};

	const hasContent = brewRecipeRaw || recipes.length > 0;

	if (loading && !brewRecipeRaw) return null;
	if (!hasContent && !isLoggedIn) return null;

	return (
		<div className="space-y-4">
			<div className="flex items-center justify-between">
				<h3 className="text-lg font-medium text-muted-foreground">Brew Recipes</h3>
				{isLoggedIn && !showForm && (
					<button
						onClick={() => setShowForm(true)}
						className="inline-flex items-center gap-1 rounded-md bg-primary px-3 py-1.5 text-xs text-primary-foreground hover:bg-primary/90"
					>
						<Plus className="size-3" />
						Add Recipe
					</button>
				)}
			</div>

			{brewRecipeRaw && (
				<Card className="border-dashed">
					<CardContent className="p-4 space-y-2">
						<div className="flex items-center gap-2">
							<Badge variant="secondary" className="gap-1">
								<Coffee className="size-3" />
								Roaster Recommendation
							</Badge>
						</div>
						<p className="text-sm text-muted-foreground whitespace-pre-line">{brewRecipeRaw}</p>
					</CardContent>
				</Card>
			)}

			{recipes.map((recipe) => (
				<Card key={recipe.id} className="shadow-sm">
					<CardContent className="p-4 space-y-3">
						<div className="flex items-start justify-between">
							<div className="space-y-1">
								<p className="font-medium text-sm">{recipe.title}</p>
								<div className="flex items-center gap-2 text-xs text-muted-foreground">
									<Badge variant="outline" className="gap-1 text-xs">
										<Coffee className="size-3" />
										{BREW_METHODS[recipe.brew_method ?? ''] ?? recipe.brew_method}
									</Badge>
									{recipe.user_email && <span>by {recipe.user_email}</span>}
								</div>
							</div>
							{currentUserId === recipe.user_id && (
								<button
									onClick={() => recipe.id && handleDelete(recipe.id)}
									className="text-muted-foreground hover:text-destructive"
								>
									<Trash2 className="size-4" />
								</button>
							)}
						</div>
						<div className="flex flex-wrap gap-3 text-xs text-muted-foreground">
							{recipe.dose_grams && (
								<span className="flex items-center gap-1">
									<Coffee className="size-3" />
									{recipe.dose_grams}g dose
								</span>
							)}
							{recipe.water_ml && (
								<span className="flex items-center gap-1">
									<Droplets className="size-3" />
									{recipe.water_ml}ml water
								</span>
							)}
							{recipe.water_temp_c && (
								<span className="flex items-center gap-1">
									<Thermometer className="size-3" />
									{recipe.water_temp_c}&deg;C
								</span>
							)}
							{recipe.grind_size && (
								<span className="flex items-center gap-1">
									<Flame className="size-3" />
									{recipe.grind_size}
								</span>
							)}
							{recipe.brew_time_seconds && (
								<span className="flex items-center gap-1">
									<Clock className="size-3" />
									{formatBrewTime(recipe.brew_time_seconds)}
								</span>
							)}
							{recipe.dose_grams && recipe.water_ml && (
								<span className="font-medium">1:{(recipe.water_ml / recipe.dose_grams).toFixed(1)} ratio</span>
							)}
						</div>
						{recipe.notes && <p className="text-sm text-muted-foreground whitespace-pre-line">{recipe.notes}</p>}
					</CardContent>
				</Card>
			))}

			{showForm && (
				<BrewRecipeForm
					coffeeId={coffeeId}
					onSaved={() => {
						setShowForm(false);
						fetchRecipes();
					}}
					onCancel={() => setShowForm(false)}
				/>
			)}
		</div>
	);
}

function BrewRecipeForm({
	coffeeId,
	onSaved,
	onCancel,
}: {
	coffeeId: number;
	onSaved: () => void;
	onCancel: () => void;
}) {
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState('');

	const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		setSubmitting(true);
		setError('');

		const form = new FormData(e.currentTarget);
		const body: Record<string, unknown> = {
			coffee_id: coffeeId,
			title: form.get('title'),
			brew_method: form.get('brew_method'),
			is_public: true,
		};

		const doseGrams = form.get('dose_grams');
		if (doseGrams) body.dose_grams = parseFloat(doseGrams as string);
		const waterMl = form.get('water_ml');
		if (waterMl) body.water_ml = parseInt(waterMl as string, 10);
		const waterTemp = form.get('water_temp_c');
		if (waterTemp) body.water_temp_c = parseInt(waterTemp as string, 10);
		const grindSize = form.get('grind_size');
		if (grindSize) body.grind_size = grindSize;
		const brewTime = form.get('brew_time_seconds');
		if (brewTime) body.brew_time_seconds = parseInt(brewTime as string, 10);
		const notes = form.get('notes');
		if (notes) body.notes = notes;

		try {
			const res = await fetch('/api/recipes', {
				method: 'POST',
				credentials: 'include',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(body),
			});
			if (!res.ok) {
				const data = await res.json().catch(() => ({}));
				setError(data.error || 'Failed to save recipe');
				return;
			}
			onSaved();
		} finally {
			setSubmitting(false);
		}
	};

	const inputClass = 'w-full rounded-md border border-input bg-background px-3 py-2 text-sm';
	const labelClass = 'block text-xs font-medium text-muted-foreground mb-1';

	return (
		<Card>
			<CardContent className="p-4">
				<form onSubmit={handleSubmit} className="space-y-4">
					<div>
						<label className={labelClass}>Title *</label>
						<input name="title" required placeholder="My go-to recipe" className={inputClass} />
					</div>

					<div>
						<label className={labelClass}>Brew Method *</label>
						<select name="brew_method" required className={inputClass}>
							{Object.entries(BREW_METHODS).map(([value, label]) => (
								<option key={value} value={value}>
									{label}
								</option>
							))}
						</select>
					</div>

					<div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
						<div>
							<label className={labelClass}>Dose (g)</label>
							<input name="dose_grams" type="number" step="0.1" min="0" placeholder="18" className={inputClass} />
						</div>
						<div>
							<label className={labelClass}>Water (ml)</label>
							<input name="water_ml" type="number" min="0" placeholder="250" className={inputClass} />
						</div>
						<div>
							<label className={labelClass}>Temp (&deg;C)</label>
							<input name="water_temp_c" type="number" min="0" max="100" placeholder="93" className={inputClass} />
						</div>
						<div>
							<label className={labelClass}>Grind Size</label>
							<select name="grind_size" className={inputClass}>
								<option value="">-</option>
								<option value="fine">Fine</option>
								<option value="medium-fine">Medium-Fine</option>
								<option value="medium">Medium</option>
								<option value="medium-coarse">Medium-Coarse</option>
								<option value="coarse">Coarse</option>
							</select>
						</div>
						<div>
							<label className={labelClass}>Brew Time (s)</label>
							<input name="brew_time_seconds" type="number" min="0" placeholder="210" className={inputClass} />
						</div>
					</div>

					<div>
						<label className={labelClass}>Notes</label>
						<textarea
							name="notes"
							rows={3}
							placeholder="Steps, tips, or anything else..."
							className={inputClass}
						/>
					</div>

					{error && <p className="text-sm text-destructive">{error}</p>}

					<div className="flex gap-2">
						<button
							type="submit"
							disabled={submitting}
							className="rounded-md bg-primary px-4 py-2 text-sm text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
						>
							{submitting ? 'Saving...' : 'Save Recipe'}
						</button>
						<button type="button" onClick={onCancel} className="rounded-md px-4 py-2 text-sm hover:bg-accent">
							Cancel
						</button>
					</div>
				</form>
			</CardContent>
		</Card>
	);
}
