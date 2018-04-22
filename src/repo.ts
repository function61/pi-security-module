import {FolderResponse, Account, Secret} from 'model';

function getJson<T>(url: string): Promise<T> {
	return fetch(url)
		.then(httpMustBeOk)
		.then((response) => response.json());
}

export function httpMustBeOk(response: Response): Promise<Response> {
	if (!response.ok) {
		return new Promise<Response>(function(_resolve: any, reject: any) {
			response.text().then((text: string) => {
				if (response.headers.get('content-type') === 'application/json') {
					const jsonError = JSON.parse(text) as {};

					if (isStructuredErrorResponse(jsonError)) {
						reject(jsonError);
						return;
					}
				}

				reject(new Error('HTTP response failure: ' + text));
			}, (err: Error) => {
				reject(new Error('HTTP response failure. Also, error fetching response body: ' + err.toString()));
			})
		});
	}

	return Promise.resolve(response);
}

export function getFolder(folderId: string): Promise<FolderResponse> {
	return getJson<FolderResponse>(`/folder/${folderId}`);
}

export function getAccount(id: string): Promise<Account> {
	return getJson<Account>(`/accounts/${id}`);
}

export function getSecrets(accountId: string): Promise<Secret[]> {
	return getJson<Secret[]>(`/accounts/${accountId}/secrets`);
}

export function searchAccounts(query: string): Promise<Account[]> {
	const searchEscaped = encodeURIComponent(query);

	return getJson<Account[]>(`/accounts?search=${searchEscaped}`);
}

export function defaultErrorHandler(err: Error) {
	alert('Error: ' + err.toString());
	console.error(err);
}

interface StructuredErrorResponse {
	error_code: string;
	error_description: string;
}

function isStructuredErrorResponse(err: StructuredErrorResponse | {}): err is StructuredErrorResponse {
	return 'error_code' in (<StructuredErrorResponse>err);
}
