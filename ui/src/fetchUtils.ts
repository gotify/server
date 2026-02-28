export type ResponseTransformer<T> = (response: Response) => Promise<T>;

export const identityTransform: ResponseTransformer<Response> = (response: Response) =>
    Promise.resolve(response);

export const jsonTransform = <T>(response: Response): Promise<T> => response.json();

export const textTransform = (response: Response): Promise<string> => response.text();

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const jsonBody: (body: any) => RequestInit = (body: any) => ({
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify(body),
});

export const yamlBody: (text: string) => RequestInit = (text: string) => ({
    method: 'POST',
    headers: {
        'Content-Type': 'application/x-yaml',
    },
    body: text,
});

export const multipartBody: (body: FormData) => RequestInit = (body: FormData) => ({
    method: 'POST',
    headers: {'content-type': 'multipart/form-data'},
    body,
});
