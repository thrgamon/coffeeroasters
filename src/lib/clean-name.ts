const COUNTRIES = [
	'ethiopia',
	'colombia',
	'kenya',
	'brazil',
	'guatemala',
	'costa rica',
	'el salvador',
	'honduras',
	'panama',
	'nicaragua',
	'mexico',
	'peru',
	'bolivia',
	'ecuador',
	'rwanda',
	'burundi',
	'tanzania',
	'uganda',
	'indonesia',
	'india',
	'papua new guinea',
	'yemen',
	'congo',
	'malawi',
];

const PROCESSES = ['natural', 'washed', 'honey', 'anaerobic', 'semi-washed'];

/** Strip redundant country, process, roast, and filler from coffee display names. */
export function cleanCoffeeName(name: string, countryName?: string): string {
	let cleaned = name;

	// "- Single Origin Espresso" -> " (Espresso)", preserving the variant indicator
	// "- Single Origin Filter" -> " (Filter)"
	cleaned = cleaned.replace(/\s*-?\s*Single Origin\s+(Filter|Espresso)\s*$/i, ' ($1)');

	// Standalone "FILTER" or "ESPRESSO" at end (not after Single Origin) - redundant with roast badge
	cleaned = cleaned.replace(/\s+(FILTER|ESPRESSO)\s*$/i, '');

	// Process in parentheses: (Natural), (Washed), etc.
	for (const p of PROCESSES) {
		cleaned = cleaned.replace(new RegExp(`\\s*\\(${p}\\)\\s*`, 'i'), ' ');
	}

	// Country name in various positions
	const countries = countryName ? [countryName, ...COUNTRIES] : COUNTRIES;
	for (const country of countries) {
		const escaped = country.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
		cleaned = cleaned.replace(new RegExp(`\\s*[-|,]\\s*${escaped}\\s*$`, 'i'), '');
		cleaned = cleaned.replace(new RegExp(`^\\s*${escaped}\\s*[-|,]\\s*`, 'i'), '');
		cleaned = cleaned.replace(new RegExp(`,\\s*${escaped}\\b`, 'i'), '');
		cleaned = cleaned.replace(new RegExp(`^${escaped}\\s+`, 'i'), '');
	}

	// "Espresso" or "Filter" prefix (roast is shown as a badge)
	cleaned = cleaned.replace(/^(Espresso|Filter)\s+/i, '');

	// Standalone process words at end: ", Washed" or ", Natural"
	for (const p of PROCESSES) {
		cleaned = cleaned.replace(new RegExp(`,\\s*${p}\\s*$`, 'i'), '');
	}

	// "Single Origin" filler without Espresso/Filter
	cleaned = cleaned.replace(/\s*[-|]\s*Single Origin\b/i, '');
	cleaned = cleaned.replace(/\bSingle Origin\s*[-|]\s*/i, '');

	// Weight suffixes
	cleaned = cleaned.replace(/\s*[|]\s*\d+g\b/i, '');
	cleaned = cleaned.replace(/\s+\d+g\s*$/i, '');

	// Trailing/leading separators
	cleaned = cleaned.replace(/\s*[|,-]\s*$/, '');
	cleaned = cleaned.replace(/^\s*[|,-]\s*/, '');

	return cleaned.trim();
}
