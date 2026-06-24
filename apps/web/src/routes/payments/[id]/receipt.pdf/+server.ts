import { env } from '$env/dynamic/private';
import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ fetch, params, url }) => {
	const baseURL = env.PAWSEAR_API_BASE_URL ?? 'http://127.0.0.1:8080';
	const download = url.searchParams.get('download') === '1' ? '?download=1' : '';
	const response = await fetch(`${baseURL}/api/payments/${params.id}/receipt/pdf${download}`);
	if (!response.ok) error(response.status, 'Receipt PDF is unavailable');
	return new Response(await response.arrayBuffer(), {
		headers: {
			'Content-Type': 'application/pdf',
			'Content-Disposition': response.headers.get('Content-Disposition') ?? 'inline'
		}
	});
};
