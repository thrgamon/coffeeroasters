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
	cleaned = cleaned.replace(/\s*(FILTER|ESPRESSO)\s*$/i, '');
	for (const p of PROCESSES) {
		cleaned = cleaned.replace(new RegExp(`\\s*\\(${p}\\)\\s*`, 'i'), ' ');
	}
	const countries = countryName ? [countryName, ...COUNTRIES] : COUNTRIES;
	for (const country of countries) {
		const escaped = country.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
		cleaned = cleaned.replace(new RegExp(`\\s*[-|,]\\s*${escaped}\\s*$`, 'i'), '');
		cleaned = cleaned.replace(new RegExp(`^\\s*${escaped}\\s*[-|,]\\s*`, 'i'), '');
		cleaned = cleaned.replace(new RegExp(`,\\s*${escaped}\\b`, 'i'), '');
		cleaned = cleaned.replace(new RegExp(`^${escaped}\\s+`, 'i'), '');
	}
	cleaned = cleaned.replace(/^(Espresso|Filter)\s+/i, '');
	for (const p of PROCESSES) {
		cleaned = cleaned.replace(new RegExp(`,\\s*${p}\\s*$`, 'i'), '');
	}
	cleaned = cleaned.replace(/\s*[-|]\s*Single Origin\b/i, '');
	cleaned = cleaned.replace(/\bSingle Origin\s*[-|]\s*/i, '');
	cleaned = cleaned.replace(/\s*[|]\s*\d+g\b/i, '');
	cleaned = cleaned.replace(/\s+\d+g\s*$/i, '');
	cleaned = cleaned.replace(/\s*[|,-]\s*$/, '');
	cleaned = cleaned.replace(/^\s*[|,-]\s*/, '');
	return cleaned.trim();
}
