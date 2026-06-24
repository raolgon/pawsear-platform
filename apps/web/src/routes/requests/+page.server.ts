import { fetchAPI, sendAPI } from '$lib/server/api';
import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

type DetectedRequest = {
	id: string;
	messageId: string;
	householdId?: string | null;
	householdName: string;
	contactId?: string | null;
	contactName: string;
	channel: string;
	body: string;
	sentAt?: string | null;
	senderExternalId?: string | null;
	serviceType?: string | null;
	startAt?: string | null;
	endAt?: string | null;
	confidence: string;
	status: string;
	reviewNotes?: string | null;
	convertedBookingStartAt?: string | null;
	convertedBookingHouseholdId?: string | null;
	createdAt: string;
};

type Contact = {
	id: string;
	displayName: string;
	telegramId?: string | null;
	whatsappId?: string | null;
};

type Household = {
	id: string;
	displayName: string;
};

type OutboundMessage = {
	id: string;
	detectedRequestId: string;
	templateKey: string;
	body: string;
	status: string;
	attempts: number;
	lastError?: string | null;
	createdAt: string;
	sentAt?: string | null;
};

type ContactHouseholdLink = {
	contactId: string;
	householdId: string;
	householdName: string;
};

export const load: PageServerLoad = async ({ fetch, url }) => {
	const requestedView = url.searchParams.get('view') ?? 'pending';
	const view = ['pending', 'history', 'all'].includes(requestedView) ? requestedView : 'pending';
	const [result, contacts, households, outboundMessages, contactHouseholds] = await Promise.all([
		fetchAPI<{ detectedRequests: DetectedRequest[] }>(fetch, '/api/detected-requests', {
			detectedRequests: []
		}),
		fetchAPI<{ contacts: Contact[] }>(fetch, '/api/contacts', { contacts: [] }),
		fetchAPI<{ households: Household[] }>(fetch, '/api/households', { households: [] }),
		fetchAPI<{ outboundMessages: OutboundMessage[] }>(fetch, '/api/outbound-messages?limit=200', {
			outboundMessages: []
		}),
		fetchAPI<{ contactHouseholds: ContactHouseholdLink[] }>(fetch, '/api/contact-household-links', {
			contactHouseholds: []
		})
	]);
	return {
		apiConnected:
			result.connected &&
			contacts.connected &&
			households.connected &&
			outboundMessages.connected &&
			contactHouseholds.connected,
		apiError:
			result.error ??
			contacts.error ??
			households.error ??
			outboundMessages.error ??
			contactHouseholds.error,
		view,
		detectedRequests: result.data.detectedRequests,
		contacts: contacts.data.contacts,
		households: households.data.households,
		outboundMessages: outboundMessages.data.outboundMessages,
		contactHouseholds: contactHouseholds.data.contactHouseholds
	};
};

export const actions: Actions = {
	importManual: async ({ fetch, request }) => {
		const formData = await request.formData();
		const body = String(formData.get('body') ?? '').trim();
		if (!body) return fail(400, { error: 'Paste a message before adding it to the review queue.' });
		const result = await sendAPI(fetch, '/api/message-imports', 'POST', {
			channel: 'manual',
			direction: 'inbound',
			body
		});
		if (result.error) return fail(400, { error: result.error });
		return { success: 'Message added to the review queue.' };
	},
	updateStatus: async ({ fetch, request }) => {
		const formData = await request.formData();
		const id = String(formData.get('id') ?? '');
		const status = String(formData.get('status') ?? '');
		const reviewNotes = String(formData.get('reviewNotes') ?? '').trim();
		if (!id || !['needs_review', 'ignored', 'needs_more_info', 'confirmed'].includes(status)) {
			return fail(400, { error: 'Choose a valid review action.' });
		}
		const result = await sendAPI(fetch, `/api/detected-requests/${id}`, 'PATCH', {
			status,
			reviewNotes: reviewNotes || undefined
		});
		if (result.error) return fail(400, { error: result.error });
		return { success: 'Request review updated.' };
	},
	linkContact: async ({ fetch, request }) => {
		const formData = await request.formData();
		const id = String(formData.get('id') ?? '');
		const contactId = String(formData.get('contactId') ?? '');
		const displayName = String(formData.get('displayName') ?? '').trim();
		const householdId = String(formData.get('householdId') ?? '');
		if (!id || (!contactId && (!displayName || !householdId))) {
			return fail(400, {
				error: 'Choose an existing contact or create the sender with a household.'
			});
		}
		const result = await sendAPI(fetch, `/api/detected-requests/${id}/contact-link`, 'POST', {
			contactId: contactId || undefined,
			displayName: displayName || undefined,
			householdId: householdId || undefined,
			role: 'owner'
		});
		if (result.error) return fail(result.status === 409 ? 409 : 400, { error: result.error });
		return { success: 'Sender linked. Future messages will recognize this contact and household.' };
	},
	queueReply: async ({ fetch, request }) => {
		const formData = await request.formData();
		const id = String(formData.get('id') ?? '');
		const templateKey = String(formData.get('templateKey') ?? '');
		if (
			!id ||
			!['request_details', 'booking_confirmed', 'request_declined'].includes(templateKey)
		) {
			return fail(400, { error: 'Choose a valid Telegram reply.' });
		}
		const result = await sendAPI(fetch, `/api/detected-requests/${id}/replies`, 'POST', {
			templateKey
		});
		if (result.error) return fail(result.status === 409 ? 409 : 400, { error: result.error });
		return { success: 'Telegram reply approved and queued for delivery.' };
	},
	linkHousehold: async ({ fetch, request }) => {
		const formData = await request.formData();
		const id = String(formData.get('id') ?? '');
		const householdId = String(formData.get('householdId') ?? '');
		if (!id || !householdId)
			return fail(400, { error: 'Choose the household for this conversation.' });
		const result = await sendAPI(fetch, `/api/detected-requests/${id}/household-link`, 'POST', {
			householdId
		});
		if (result.error) return fail(result.status === 409 ? 409 : 400, { error: result.error });
		return { success: 'Household selected and remembered for this conversation.' };
	}
};
