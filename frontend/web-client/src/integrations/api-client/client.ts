import { env } from "@/env";
import { ApiError } from "@/lib/api-error";

const apiBase = () => `${env.NEXT_PUBLIC_API_URL}/api/v1`;

type Schema<T> = { parse: (input: unknown) => T };

async function parseResponse<T>(res: Response, schema: Schema<T>): Promise<T> {
	const json: unknown = await res.json().catch(() => null);
	if (!res.ok) {
		throw ApiError.fromEnvelope(json);
	}
	if (!json || typeof json !== "object" || !("data" in json)) {
		throw new ApiError("error.invalid_response");
	}
	return schema.parse((json as { data: unknown }).data);
}

export const api = {
	async get<T>(path: string, schema: Schema<T>): Promise<T> {
		const res = await fetch(`${apiBase()}${path}`, {
			credentials: "include",
		});
		return parseResponse(res, schema);
	},

	async post<T>(path: string, schema: Schema<T>, body: unknown): Promise<T> {
		const res = await fetch(`${apiBase()}${path}`, {
			method: "POST",
			credentials: "include",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(body),
		});
		return parseResponse(res, schema);
	},
};
