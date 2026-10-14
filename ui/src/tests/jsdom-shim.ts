// Minimal DOM shim so unit tests can import UI source modules (which touch
// `window`/`document`) without pulling in a full jsdom dependency.
if (typeof (globalThis as {window?: unknown}).window === 'undefined') {
    const win = {
        config: {},
        location: {
            href: 'http://localhost/',
            protocol: 'http:',
            host: 'localhost',
            origin: 'http://localhost',
        },
        navigator: {userAgent: 'vitest'},
        addEventListener: () => undefined,
        removeEventListener: () => undefined,
        dispatchEvent: () => true,
        setTimeout,
        clearTimeout,
        setInterval,
        clearInterval,
        localStorage: {
            getItem: () => null,
            setItem: () => undefined,
            removeItem: () => undefined,
        },
    };
    const doc = {
        createElement: () => ({
            style: {},
            setAttribute: () => undefined,
            appendChild: () => undefined,
        }),
        addEventListener: () => undefined,
        removeEventListener: () => undefined,
        querySelector: () => null,
        body: {appendChild: () => undefined, removeChild: () => undefined},
    };
    (globalThis as {window: unknown}).window = win;
    (globalThis as {document: unknown}).document = doc;
}
