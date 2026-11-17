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

const actionButton = z.object({
    label: z.string().check(z.trim(), z.minLength(1)),
    http: z.object({
        method: z.string().check(z.toUpperCase()),
        url: z.url({protocol: /^https?$/}),
        headers: z._default(z.record(z.string(), z.string()), {}),
        body: z.optional(z.string()),
    }),
});

export type IActionButton = z.infer<typeof actionButton>;

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
