const API = process.env.NEXT_PUBLIC_API_URL || '';

export async function adminFetch<T>(path: string, opts?: RequestInit): Promise<T> {
	const res = await fetch(`${API}${path}`, { credentials: 'include', ...opts });
	if (!res.ok) throw new Error(`${res.status}`);
	return res.json();
}

export async function adminPost<T>(path: string, body: unknown): Promise<T> {
	return adminFetch(path, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body),
	});
}

export async function adminPut<T>(path: string, body: unknown): Promise<T> {
	return adminFetch(path, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body),
	});
}
