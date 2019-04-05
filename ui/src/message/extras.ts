export enum RenderMode {
    Markdown = 'text/markdown',
    Plain = 'text/plain',
}

export const contentType = (extras?: IMessageExtras): RenderMode => {
    const type = extract(extras, 'client::display', 'contentType');
    const valid = Object.keys(RenderMode)
        .map((k) => RenderMode[k])
        .some((mode) => mode === type);
    return valid ? type : RenderMode.Plain;
};

// tslint:disable-next-line:no-any
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
