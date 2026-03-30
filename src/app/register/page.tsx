'use client';

import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export default function RegisterPage() {
	const router = useRouter();

	useEffect(() => {
		// Passwordless auth - no separate registration needed
		router.replace('/login');
	}, [router]);

	return null;
}
