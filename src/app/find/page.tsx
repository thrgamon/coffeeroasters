'use client';

import { useState } from 'react';
import CoffeeCard from '@/components/CoffeeCard';
import type { DomainCoffeeResponse } from '@/lib/api/generated/models';

const QUESTIONS = [
	{
		key: 'sweetness',
		title: 'What kind of sweetness do you prefer?',
		options: [
			{ value: 'fruity', label: 'Fruity', description: 'Berries, stone fruit, tropical' },
			{ value: 'caramel', label: 'Caramel and chocolate', description: 'Toffee, nuts, cocoa' },
			{ value: 'both', label: 'I like both', description: 'A balance of fruit and sweetness' },
		],
	},
	{
		key: 'brightness',
		title: 'How do you feel about brightness or tanginess?',
		options: [
			{ value: 'bright', label: 'Bright and zingy', description: 'Like biting into citrus' },
			{ value: 'smooth', label: 'Smooth and mellow', description: 'Easy-going, no sharp edges' },
			{ value: 'neutral', label: 'Somewhere in between', description: '' },
		],
	},
	{
		key: 'body',
		title: 'What kind of body do you like?',
		options: [
			{ value: 'light', label: 'Light and tea-like', description: 'Delicate, clean finish' },
			{ value: 'full', label: 'Rich and full', description: 'Heavy, syrupy mouthfeel' },
			{ value: 'neutral', label: 'No preference', description: '' },
		],
	},
	{
		key: 'appeal',
		title: 'Which of these sounds most appealing?',
		options: [
			{ value: 'floral', label: 'Jasmine, rose, bergamot', description: 'Perfumed and aromatic' },
			{ value: 'chocolate', label: 'Chocolate, nuts, toffee', description: 'Comforting and rich' },
			{ value: 'berry', label: 'Berries and jam', description: 'Juicy and sweet' },
			{ value: 'earthy', label: 'Earth, wood, spice', description: 'Grounded and savoury' },
		],
	},
	{
		key: 'adventurous',
		title: 'How adventurous are you feeling?',
		options: [
			{ value: 'classic', label: 'Keep it classic', description: 'Traditional processes, familiar flavours' },
			{ value: 'surprise', label: 'Surprise me', description: 'Experimental, unusual, bold' },
		],
	},
] as const;

type Answers = Record<string, string>;

export default function FindCoffeePage() {
	const [step, setStep] = useState(0);
	const [answers, setAnswers] = useState<Answers>({});
	const [results, setResults] = useState<DomainCoffeeResponse[] | null>(null);
	const [loading, setLoading] = useState(false);

	const fetchResults = async (finalAnswers: Answers) => {
		setLoading(true);
		try {
			const params = new URLSearchParams(finalAnswers);
			const apiBase = process.env.NEXT_PUBLIC_API_URL || '';
			const res = await fetch(`${apiBase}/api/coffees/find?${params}`);
			const data = await res.json();
			setResults(data.coffees || []);
		} finally {
			setLoading(false);
		}
	};

	const handleSelect = (key: string, value: string) => {
		const newAnswers = { ...answers, [key]: value };
		setAnswers(newAnswers);

		if (step < QUESTIONS.length - 1) {
			setStep(step + 1);
		} else {
			fetchResults(newAnswers);
		}
	};

	const handleBack = () => {
		if (step > 0) {
			setStep(step - 1);
		}
	};

	const handleStartOver = () => {
		setStep(0);
		setAnswers({});
		setResults(null);
	};

	if (loading) {
		return (
			<div className="container mx-auto px-4 py-16 flex flex-col items-center gap-4">
				<div className="h-8 w-8 animate-spin rounded-full border-4 border-gold border-t-transparent" />
				<p className="text-grey-olive font-mono text-sm uppercase tracking-wider">Finding your perfect coffees...</p>
			</div>
		);
	}

	if (results !== null) {
		return (
			<div className="container mx-auto px-4 py-12 max-w-4xl">
				<h2 className="text-2xl font-bold uppercase tracking-wider text-snow mb-2">Your matches</h2>
				<p className="text-grey-olive mb-8">
					Based on your preferences, here are the coffees we think you&apos;ll enjoy.
				</p>

				{results.length === 0 ? (
					<p className="text-grey-olive">No matches found. Try different answers.</p>
				) : (
					<div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
						{results.map((coffee) => (
							<CoffeeCard key={coffee.id} coffee={coffee} />
						))}
					</div>
				)}

				<div className="mt-10">
					<button
						type="button"
						onClick={handleStartOver}
						className="text-sm text-gold hover:underline transition-colors uppercase tracking-wider font-medium"
					>
						Start over
					</button>
				</div>
			</div>
		);
	}

	const question = QUESTIONS[step];
	const progress = ((step + 1) / QUESTIONS.length) * 100;

	return (
		<div className="container mx-auto px-4 py-12 max-w-2xl">
			<h1 className="text-3xl font-bold uppercase tracking-wider text-snow mb-2">Find your coffee</h1>
			<p className="text-grey-olive mb-8">
				Answer a few questions about what you enjoy and we&apos;ll match you with coffees you&apos;ll love.
			</p>

			<div className="mb-8">
				<div className="flex justify-between text-xs text-grey-olive font-mono uppercase tracking-wider mb-2">
					<span>
						Step {step + 1} of {QUESTIONS.length}
					</span>
				</div>
				<div className="h-1 w-full rounded-full bg-secondary overflow-hidden">
					<div
						className="h-full bg-gold rounded-full transition-all duration-300"
						style={{ width: `${progress}%` }}
					/>
				</div>
			</div>

			<h2 className="text-xl font-bold text-snow mb-6">{question.title}</h2>

			<div className="grid gap-4 sm:grid-cols-2">
				{question.options.map((option) => {
					const selected = answers[question.key] === option.value;
					return (
						<button
							key={option.value}
							type="button"
							onClick={() => handleSelect(question.key, option.value)}
							className={`rounded border p-6 cursor-pointer transition-all text-left ${
								selected
									? 'border-gold bg-gold/10 shadow-[0_0_15px_rgba(255,213,0,0.1)]'
									: 'border-border/50 hover:border-gold/30'
							}`}
						>
							<p className="font-bold text-snow">{option.label}</p>
							{option.description && <p className="text-sm text-grey-olive mt-1">{option.description}</p>}
						</button>
					);
				})}
			</div>

			{step > 0 && (
				<div className="mt-6">
					<button
						type="button"
						onClick={handleBack}
						className="text-sm text-grey-olive hover:text-gold transition-colors uppercase tracking-wider"
					>
						Back
					</button>
				</div>
			)}
		</div>
	);
}
