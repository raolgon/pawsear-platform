import { env } from '$env/dynamic/private';

const DEFAULT_API_BASE_URL = 'http://127.0.0.1:8080';

export type ApiResult<T> = {
	data: T;
	connected: boolean;
	status?: number;
	error?: string;
};

export async function fetchAPI<T>(
	fetchFn: typeof fetch,
	path: string,
	fallback: T
): Promise<ApiResult<T>> {
	const baseURL = env.PAWSEAR_API_BASE_URL ?? DEFAULT_API_BASE_URL;
	try {
		const response = await fetchFn(`${baseURL}${path}`);
		if (!response.ok) {
			return {
				data: fallback,
				connected: false,
				status: response.status,
				error: `API returned ${response.status}`
			};
		}

		return {
			data: (await response.json()) as T,
			connected: true,
			status: response.status
		};
	} catch (error) {
		return {
			data: fallback,
			connected: false,
			error: error instanceof Error ? error.message : 'API unavailable'
		};
	}
}

export async function sendAPI<T>(
	fetchFn: typeof fetch,
	path: string,
	method: 'POST' | 'PATCH',
	body: unknown
): Promise<ApiResult<T>> {
	const baseURL = env.PAWSEAR_API_BASE_URL ?? DEFAULT_API_BASE_URL;
	try {
		const response = await fetchFn(`${baseURL}${path}`, {
			method,
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(body)
		});

		const payload = (await response.json().catch(() => ({}))) as T;
		if (!response.ok) {
			return {
				data: payload,
				connected: true,
				status: response.status,
				error:
					typeof payload === 'object' && payload && 'message' in payload
						? String(payload.message)
						: `API returned ${response.status}`
			};
		}

		return {
			data: payload,
			connected: true,
			status: response.status
		};
	} catch (error) {
		return {
			data: {} as T,
			connected: false,
			error: error instanceof Error ? error.message : 'API unavailable'
		};
	}
}
