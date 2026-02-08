const BASE_URL = import.meta.env.DEV ? 'http://localhost:8080/api/v1' : '/api/v1';

export interface ApiClientOptions extends RequestInit {
	/**
	 * Timeout in ms
	 * @default 15000
	 */
	timeout?: number;

	/**
	 * Expected response type:
	 * - 'json' (default)
	 * - 'text' => for plain text responses
	 * - 'blob' => for files, pdf, images
	 * - 'arrayBuffer' => for binary data
	 */
	responseType?: 'json' | 'text' | 'blob' | 'arrayBuffer';
}

export async function apiClient<T>(
	kitFetch: typeof fetch,
	path: string,
	options: ApiClientOptions = {}
): Promise<T> {
	const url = `${BASE_URL}${path}`;

	// destructure custom timeout with option
	const { timeout = 15000, responseType = 'json', ...fetchOptions } = options;

	const controller = new AbortController();
	const timeoutId = setTimeout(() => controller.abort(), timeout);

	const mergedOptions: RequestInit = {
		...fetchOptions,
		signal: controller.signal,
		credentials: 'include' // Important: include cookies for auth
	};

	let response: Response;
	try {
		response = await kitFetch(url, mergedOptions);
	} catch (error) {
		if (error instanceof DOMException && error.name === 'AbortError') {
			throw new Error(`API Error: Request timed out after ${timeout}ms for ${url}`);
		}
		// unknown errors
		throw error;
	} finally {
		clearTimeout(timeoutId);
	}

	// CRITICAL: HTTP Error (400, 500)
	if (!response.ok) {
		let errBody = '';
		let contentType: string | null = null;

		try {
			contentType = response.headers.get('Content-Type');
		} catch (headerError) {
			// In SSR context, if header isn't allowed by filterSerializedResponseHeaders,
			// error will occur. Skip content-type checking in that case.
			console.warn('Unable to access Content-Type header:', headerError);
		}

		if (contentType) {
			if (contentType.includes('application/json')) {
				try {
					const errJson = await response.json();
					errBody = JSON.stringify(errJson.error || errJson.message || errJson);
				} catch {
					errBody = '(failed to parse JSON response)';
				}
			} else if (contentType.includes('text/plain')) {
				const text = await response.text();
				// Truncate long text responses (e.g., HTML error pages)
				errBody = text.length > 200 ? text.substring(0, 200) + '...' : text;
			}
		} else {
			// If we can't get content-type, try to parse as text but truncate it
			try {
				const text = await response.text();
				// Check if it looks like HTML (common for error pages)
				if (text.trim().startsWith('<!doctype') || text.trim().startsWith('<html')) {
					errBody = '(HTML error page returned - check server logs)';
				} else {
					// Truncate long responses
					errBody = text.length > 200 ? text.substring(0, 200) + '...' : text;
				}
			} catch {
				// If even text parsing fails, leave errBody empty
			}
		}

		const statusText = `[${response.status} ${response.statusText}]`;
		const finalMessage = `API Error ${statusText} for ${url}${errBody ? `: ${errBody}` : ''}`;
		throw new Error(finalMessage);
	}

	// 204: No Content
	// .json will fail so return an empty object
	if (response.status == 204) {
		return {} as T;
	}

	switch (responseType) {
		case 'text':
			return (await response.text()) as T;
		case 'blob':
			return (await response.blob()) as T;
		case 'arrayBuffer':
			return (await response.arrayBuffer()) as T;
		case 'json':
		default: {
			// Backend wraps responses in { code, message, data } structure
			// We need to unwrap the data field
			const json = await response.json();
			if (json && typeof json === 'object' && 'data' in json) {
				return json.data as T;
			}
			return json as T;
		}
	}
}

type KitFetch = typeof window.fetch;

export const http = {
	/**
	 * Perform a `GET` request.
	 */
	get: <T>(fetch: KitFetch, path: string, options: ApiClientOptions = {}) => {
		return apiClient<T>(fetch, path, options); // no need to do anything, GET request by default
	},

	/**
	 * Perform a `POST` request
	 */
	post: <T>(fetch: KitFetch, path: string, data: unknown, options: ApiClientOptions = {}) => {
		return apiClient<T>(fetch, path, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data),
			...options
		});
	},

	/**
	 * Perform a `PUT` request
	 */
	put: <T>(fetch: KitFetch, path: string, data: unknown, options: ApiClientOptions = {}) => {
		return apiClient<T>(fetch, path, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data),
			...options
		});
	},

	/**
	 * Perform a `DELETE` request
	 */
	delete: <T = Record<string, never>>(
		fetch: KitFetch,
		path: string,
		options: ApiClientOptions = {}
	) => {
		return apiClient<T>(fetch, path, {
			method: 'DELETE',
			...options
		});
	}
};
