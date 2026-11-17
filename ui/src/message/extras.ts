import {IMessageExtras} from '../types';

export enum RenderMode {
    Markdown = 'text/markdown',
    Plain = 'text/plain',
}

export const contentType = (extras?: IMessageExtras): RenderMode => {
    const type = extract(extras, 'client::display', 'contentType');
    const valid = Object.values(RenderMode).includes(type);
    return valid ? type : RenderMode.Plain;
};

export interface IActionButton {
    label: string;
    http: {
        method: string;
        url: string;
        headers: Record<string, string>;
    };
}

const allowedMethods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'];

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const parseActionButton = (button: any): IActionButton | null => {
    if (!button || typeof button !== 'object') {
        return null;
    }
    const label = typeof button.label === 'string' ? button.label.trim() : '';
    const http = button.http;
    if (label === '' || !http || typeof http !== 'object' || typeof http.url !== 'string') {
        return null;
    }

    let url: URL;
    try {
        url = new URL(http.url);
    } catch {
        return null;
    }
    if (url.protocol !== 'http:' && url.protocol !== 'https:') {
        return null;
    }

    const method = typeof http.method === 'string' ? http.method.toUpperCase() : 'POST';
    if (!allowedMethods.includes(method)) {
        return null;
    }

    const headers: Record<string, string> = {};
    if (http.headers && typeof http.headers === 'object' && !Array.isArray(http.headers)) {
        for (const [name, value] of Object.entries(http.headers)) {
            if (typeof value === 'string') {
                headers[name] = value;
            }
        }
    }

    return {label, http: {method, url: url.toString(), headers}};
};

/**
 * Returns the valid buttons of `extras["client::action"].buttons`.
 * Invalid buttons (missing label, non http(s) url, unknown method) are skipped.
 */
export const actionButtons = (extras?: IMessageExtras): IActionButton[] => {
    const buttons = extract(extras, 'client::action', 'buttons');
    if (!Array.isArray(buttons)) {
        return [];
    }
    return buttons
        .map(parseActionButton)
        .filter((button): button is IActionButton => button !== null);
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const extract = (extras: IMessageExtras | undefined, key: string, path: string): any => {
    if (!extras) {
        return null;
    }

    if (!extras[key]) {
        return null;
    }

    if (!extras[key][path]) {
        return null;
    }

    return extras[key][path];
};
