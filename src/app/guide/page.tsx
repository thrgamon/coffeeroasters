import type { Metadata } from 'next';
import Link from 'next/link';

export const metadata: Metadata = {
	title: 'Coffee Guide | Coffeeroasters',
	description: 'Learn how processing, roasting, origin, and variety shape coffee flavour',
};

const processes = [
	{
		name: 'Washed',
		value: 'washed',
		description:
			"The cherry is removed before drying, letting the bean's inherent character shine through. Produces the cleanest, most transparent cups.",
		flavours: ['Citrus', 'Floral', 'Tea-like'],
		body: 'Light to medium',
		acidity: 'Bright, pronounced',
	},
	{
		name: 'Natural',
		value: 'natural',
		description:
			'The bean dries inside the whole cherry, absorbing sugars and fruit compounds. Produces bold fruitiness and heavier body.',
		flavours: ['Blueberry', 'Strawberry', 'Tropical fruit', 'Wine-like'],
		body: 'Full, syrupy',
		acidity: 'Moderate',
	},
	{
		name: 'Honey',
		value: 'honey',
		description:
			'A spectrum between washed and natural. More mucilage left on the bean means more sweetness and fruit character.',
		flavours: ['Caramel', 'Stone fruit', 'Buttery sweetness'],
		body: 'Medium to full',
		acidity: 'Moderate',
	},
	{
		name: 'Anaerobic',
		value: 'anaerobic',
		description:
			'Fermented in oxygen-free tanks, concentrating specific flavour compounds. Amplifies what the origin and variety already bring.',
		flavours: ['Intense tropical fruit', 'Spice', 'Wine-like'],
		body: 'Full',
		acidity: 'Varies',
	},
	{
		name: 'Wet-hulled',
		value: 'wet-hulled',
		description:
			'Unique to Indonesia. Beans are hulled at high moisture, creating a distinctive earthy, heavy profile.',
		flavours: ['Earth', 'Wood', 'Cedar', 'Tobacco', 'Spice'],
		body: 'Very full',
		acidity: 'Very low',
	},
	{
		name: 'Experimental',
		value: 'experimental',
		description: 'Includes carbonic maceration and other novel techniques. Results vary widely.',
		flavours: ['Boozy', 'Tropical', 'Funky'],
		body: 'Varies',
		acidity: 'Varies',
	},
];

const roastLevels = [
	{
		name: 'Light',
		value: 'light',
		description: 'Roasted just past first crack. Origin character dominates.',
		notes: ['Bright acidity', 'Floral', 'Citrus', 'Tea-like'],
	},
	{
		name: 'Medium-light',
		value: 'medium-light',
		description: 'Origin character with emerging caramel sweetness. Acidity starts to round.',
		notes: ['Stone fruit', 'Honey', 'Caramel hints'],
	},
	{
		name: 'Medium',
		value: 'medium',
		description: 'The midpoint. Balanced between origin and roast character.',
		notes: ['Caramel', 'Toasted nuts', 'Milk chocolate'],
	},
	{
		name: 'Medium-dark',
		value: 'medium-dark',
		description: 'Roast character begins to lead. Acidity is muted, body is full.',
		notes: ['Dark chocolate', 'Molasses', 'Dried fruit'],
	},
	{
		name: 'Dark',
		value: 'dark',
		description: 'Roast-dominant. Origin character is largely gone. Oil visible on the surface.',
		notes: ['Smoky', 'Bittersweet', 'Burnt sugar'],
	},
];

const originGroups = [
	{
		label: 'Bright and floral',
		origins: [
			{
				name: 'Ethiopia',
				code: 'ET',
				summary:
					'The birthplace of coffee. Wild and cultivated varieties produce extraordinary range, from jasmine-scented naturals to tea-like washed coffees.',
			},
			{
				name: 'Kenya',
				code: 'KE',
				summary:
					'Known for bold blackcurrant and tomato-like acidity. SL28 and SL34 varieties dominate, producing complex, wine-like cups.',
			},
			{
				name: 'Rwanda',
				code: 'RW',
				summary:
					'Produces vibrant, fruit-forward coffees with bright acidity and floral top notes. Bourbon variety is common.',
			},
			{
				name: 'Burundi',
				code: 'BI',
				summary:
					'Similar profile to Rwanda with pronounced fruit acidity. Bourbon-based lots from high-altitude washing stations.',
			},
		],
	},
	{
		label: 'Balanced and sweet',
		origins: [
			{
				name: 'Colombia',
				code: 'CO',
				summary:
					'Year-round harvests from diverse microclimates. Expect caramel sweetness, stone fruit, and clean, approachable acidity.',
			},
			{
				name: 'Guatemala',
				code: 'GT',
				summary:
					'Volcanic soils produce chocolatey, full-bodied coffees with bright acidity and brown sugar sweetness.',
			},
			{
				name: 'Costa Rica',
				code: 'CR',
				summary:
					'Clean, sweet, and consistent. Honey and natural processing are common, adding fruit complexity to the baseline sweetness.',
			},
			{
				name: 'El Salvador',
				code: 'SV',
				summary: 'Bourbon-heavy production with a soft, rounded profile. Stone fruit sweetness and gentle acidity.',
			},
			{
				name: 'Honduras',
				code: 'HN',
				summary: 'Growing reputation for fruit-forward lots. Caramel, stone fruit, and mild acidity at its best.',
			},
			{
				name: 'Panama',
				code: 'PA',
				summary:
					"Home of Gesha. Exceptional terroir in Boquete produces some of the world's most nuanced and complex coffees.",
			},
		],
	},
	{
		label: 'Chocolatey and nutty',
		origins: [
			{
				name: 'Brazil',
				code: 'BR',
				summary:
					"The world's largest producer. Natural processing is dominant, producing low-acid, chocolatey, nutty cups. Common as espresso base.",
			},
			{
				name: 'Peru',
				code: 'PE',
				summary: 'Underrated origin with clean, balanced cups. Chocolate, caramel, and mild fruit in the better lots.',
			},
			{
				name: 'Nicaragua',
				code: 'NI',
				summary: 'Soft body with caramel and chocolate notes. Improving quality from higher-altitude farms.',
			},
			{
				name: 'Mexico',
				code: 'MX',
				summary: 'Light body, mild acidity, and nutty sweetness. Often used as a base for blends.',
			},
		],
	},
	{
		label: 'Earthy and spiced',
		origins: [
			{
				name: 'Indonesia',
				code: 'ID',
				summary:
					'Wet-hulled processing creates distinctive earthy, full-bodied cups. Sumatra, Java, and Sulawesi each have their own character.',
			},
			{
				name: 'India',
				code: 'IN',
				summary:
					'Spiced, earthy notes with low acidity. Monsoon Malabar is a unique process that creates an intensely woody, funky cup.',
			},
			{
				name: 'Papua New Guinea',
				code: 'PG',
				summary:
					'Earthy and herbaceous with some fruit character. Wild growing conditions create inconsistency but also interesting complexity.',
			},
		],
	},
];

const varietyGroups = [
	{
		label: 'Classic',
		varieties: [
			{
				name: 'Bourbon',
				value: 'bourbon',
				tendency: 'Sweet, round, with red fruit and caramel',
				body: 'Medium',
				acidity: 'Moderate',
			},
			{
				name: 'Typica',
				value: 'typica',
				tendency: 'Clean and elegant, with floral and citrus notes',
				body: 'Light to medium',
				acidity: 'Bright',
			},
			{
				name: 'Caturra',
				value: 'caturra',
				tendency: 'Bright and citrusy, a Bourbon mutation',
				body: 'Light',
				acidity: 'High',
			},
			{
				name: 'Catuai',
				value: 'catuai',
				tendency: 'Mild and balanced, bred for yield not complexity',
				body: 'Medium',
				acidity: 'Moderate',
			},
		],
	},
	{
		label: 'High-complexity',
		varieties: [
			{
				name: 'Gesha',
				value: 'gesha',
				tendency: 'Extraordinarily floral, jasmine and bergamot, tea-like',
				body: 'Light',
				acidity: 'Delicate',
			},
			{
				name: 'Pacamara',
				value: 'pacamara',
				tendency: 'Large bean, bold fruit, and wine-like complexity',
				body: 'Full',
				acidity: 'Pronounced',
			},
			{
				name: 'Sidra',
				value: 'sidra',
				tendency: 'Intense tropical fruit and spice, often compared to Gesha',
				body: 'Medium',
				acidity: 'Bright',
			},
		],
	},
	{
		label: 'Kenyan selections',
		varieties: [
			{
				name: 'SL28',
				value: 'sl28',
				tendency: 'Blackcurrant, tomato, and intense fruit acidity',
				body: 'Full',
				acidity: 'Very high',
			},
			{
				name: 'SL34',
				value: 'sl34',
				tendency: 'Similar to SL28 but rounder and less aggressive',
				body: 'Medium to full',
				acidity: 'High',
			},
		],
	},
	{
		label: 'Ethiopian',
		varieties: [
			{
				name: 'Heirloom',
				value: 'heirloom',
				tendency: 'Umbrella term for thousands of wild and landrace varieties. Unpredictably complex.',
				body: 'Varies',
				acidity: 'Varies',
			},
		],
	},
	{
		label: 'Modern hybrids',
		varieties: [
			{
				name: 'Castillo',
				value: 'castillo',
				tendency: 'Disease-resistant hybrid. Clean and consistent, mild fruit',
				body: 'Medium',
				acidity: 'Moderate',
			},
			{
				name: 'Catimor',
				value: 'catimor',
				tendency: 'Robusta hybrid, bred for resistance. Can taste flat or rubbery at lower altitudes',
				body: 'Medium to full',
				acidity: 'Low',
			},
		],
	},
];

export default function GuidePage() {
	return (
		<div className="space-y-12">
			<div className="space-y-3">
				<h1 className="text-3xl font-bold">Coffee Guide</h1>
				<p className="text-muted-foreground max-w-2xl">
					Four factors shape what you taste in the cup: how the coffee was processed after harvest, how darkly it was
					roasted, where it was grown, and which variety of plant it came from. Each one leaves a distinct mark on the
					flavour.
				</p>
			</div>

			{/* Process */}
			<section className="space-y-6">
				<h2 className="text-2xl font-semibold">Process</h2>
				<div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
					{processes.map((p) => (
						<div key={p.value} className="rounded-lg border bg-card p-6 space-y-3">
							<p className="font-semibold">{p.name}</p>
							<p className="text-sm text-muted-foreground">{p.description}</p>
							<div className="text-sm space-y-1">
								<p>
									<span className="font-medium">Typical flavours:</span>{' '}
									<span className="text-muted-foreground">{p.flavours.join(', ')}</span>
								</p>
								<p>
									<span className="font-medium">Body:</span> <span className="text-muted-foreground">{p.body}</span>
								</p>
								<p>
									<span className="font-medium">Acidity:</span>{' '}
									<span className="text-muted-foreground">{p.acidity}</span>
								</p>
							</div>
							<Link href={`/coffees?process=${p.value}`} className="text-sm text-primary hover:underline">
								See {p.name} coffees
							</Link>
						</div>
					))}
				</div>
			</section>

			{/* Roast Level */}
			<section className="space-y-6">
				<h2 className="text-2xl font-semibold">Roast Level</h2>
				<div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
					{roastLevels.map((r) => (
						<div key={r.value} className="rounded-lg border bg-card p-6 space-y-3">
							<p className="font-semibold">{r.name}</p>
							<p className="text-sm text-muted-foreground">{r.description}</p>
							<p className="text-sm">
								<span className="font-medium">Typical notes:</span>{' '}
								<span className="text-muted-foreground">{r.notes.join(', ')}</span>
							</p>
							<Link href={`/coffees?roast=${r.value}`} className="text-sm text-primary hover:underline">
								See {r.name} roast coffees
							</Link>
						</div>
					))}
				</div>
			</section>

			{/* Origin */}
			<section className="space-y-8">
				<h2 className="text-2xl font-semibold">Origin</h2>
				{originGroups.map((group) => (
					<div key={group.label} className="space-y-4">
						<h3 className="text-lg font-medium">{group.label}</h3>
						<div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
							{group.origins.map((origin) => (
								<div key={origin.code} className="rounded-lg border bg-card p-6 space-y-3">
									<p className="font-semibold">{origin.name}</p>
									<p className="text-sm text-muted-foreground">{origin.summary}</p>
									<Link href={`/coffees?origin=${origin.code}`} className="text-sm text-primary hover:underline">
										See coffees from {origin.name}
									</Link>
								</div>
							))}
						</div>
					</div>
				))}
			</section>

			{/* Variety */}
			<section className="space-y-8">
				<h2 className="text-2xl font-semibold">Variety</h2>
				{varietyGroups.map((group) => (
					<div key={group.label} className="space-y-4">
						<h3 className="text-lg font-medium">{group.label}</h3>
						<div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
							{group.varieties.map((v) => (
								<div key={v.value} className="rounded-lg border bg-card p-6 space-y-3">
									<p className="font-semibold">{v.name}</p>
									<p className="text-sm text-muted-foreground">{v.tendency}</p>
									<div className="text-sm space-y-1">
										<p>
											<span className="font-medium">Body:</span> <span className="text-muted-foreground">{v.body}</span>
										</p>
										<p>
											<span className="font-medium">Acidity:</span>{' '}
											<span className="text-muted-foreground">{v.acidity}</span>
										</p>
									</div>
									<Link href={`/coffees?variety=${v.value}`} className="text-sm text-primary hover:underline">
										See {v.name} coffees
									</Link>
								</div>
							))}
						</div>
					</div>
				))}
			</section>
		</div>
	);
}
