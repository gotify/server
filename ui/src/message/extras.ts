import {z} from 'zod/mini';
import {IMessageExtras} from '../types';

export enum RenderMode {
    Markdown = 'text/markdown',
    Plain = 'text/plain',
}

const displayExtras = z.object({
    'client::display': z.object({contentType: z.enum(RenderMode)}),
});

export const contentType = (extras?: IMessageExtras): RenderMode => {
    const result = displayExtras.safeParse(extras);
    return result.success ? result.data['client::display'].contentType : RenderMode.Plain;
};

// fetch rejects these methods and a body on GET/HEAD requests without sending anything.
const forbiddenMethods = ['CONNECT', 'TRACE', 'TRACK'];
const bodylessMethods = ['GET', 'HEAD'];

const httpAction = z
    .object({
        method: z.string().check(
            z.toUpperCase(),
            z.refine((method) => !forbiddenMethods.includes(method), 'Unsupported method')
        ),
        url: z.url({protocol: /^https?$/}),
        headers: z._default(z.record(z.string(), z.string()), {}),
        body: z.optional(z.string()),
    })
    .check(
        z.refine(({method, body}) => body === undefined || !bodylessMethods.includes(method), {
            message: 'GET and HEAD requests cannot have a body',
            path: ['body'],
        })
    );

const openAction = z.object({
    url: z.url({protocol: /^https?$/}),
});

// A button either sends a request (http) or opens a page (open).
const actionButton = z
    .object({
        label: z.string().check(z.trim(), z.minLength(1)),
        http: z.optional(httpAction),
        open: z.optional(openAction),
    })
    .check(
        z.refine(
            ({http, open}) => (http === undefined) !== (open === undefined),
            'A button needs either http or open'
        )
    );

export type IActionButton = z.infer<typeof actionButton>;
export type IHttpAction = z.infer<typeof httpAction>;

const actionExtras = z.optional(
    z.object({
        'client::action': z.nullish(z.object({buttons: z.array(actionButton)})),
    })
);

const validButtons = (rawButtons: unknown): IActionButton[] =>
    Array.isArray(rawButtons)
        ? rawButtons.flatMap((rawButton) => actionButton.safeParse(rawButton).data ?? [])
        : [];

export const actionButtons = (
    extras?: IMessageExtras
): {buttons: IActionButton[]; error?: string} => {
    const {error} = actionExtras.safeParse(extras);
    return {
        buttons: validButtons(extras?.['client::action']?.buttons),
        error: error && z.prettifyError(error),
    };
};
