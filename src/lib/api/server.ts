const API_BASE = process.env.API_URL || 'http://localhost:8080';

export async function apiFetch<T>(path: string, init?: { revalidate?: number }): Promise<T> {
	const res = await fetch(`${API_BASE}${path}`, {
		next: { revalidate: init?.revalidate ?? 3600 },
	});
	if (!res.ok) throw new Error(`API ${res.status}`);
	return res.json();
}
