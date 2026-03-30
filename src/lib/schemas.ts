import { z } from 'zod';

// Base schemas - extend generated schemas here when Orval produces zod.unknown() for nested $ref
export const UserSchema = z.object({
	id: z.number(),
	email: z.string(),
	is_admin: z.boolean().optional().default(false),
});

export const AuthResponseSchema = z.object({
	user: UserSchema,
});
