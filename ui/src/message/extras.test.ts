import {describe, expect, it} from 'vitest';
import {actionButtons} from './extras';

const withButtons = (buttons: unknown) => ({'client::action': {buttons}});

describe('actionButtons', () => {
    it('returns nothing without extras', () => {
        expect(actionButtons(undefined)).toEqual([]);
        expect(actionButtons({})).toEqual([]);
        expect(actionButtons({'client::action': {}})).toEqual([]);
        expect(actionButtons(withButtons('nope'))).toEqual([]);
    });

    it('parses the format from gotify/server#494', () => {
        const buttons = actionButtons(
            withButtons([
                {
                    label: 'Power off',
                    http: {
                        method: 'post',
                        url: 'https://example.com/power_off?id=5',
                        headers: {Authorization: 'bearer token'},
                    },
                },
                {label: 'Dismiss', http: {method: 'get', url: 'http://example.com/dismiss'}},
            ])
        );
        expect(buttons).toEqual([
            {
                label: 'Power off',
                http: {
                    method: 'POST',
                    url: 'https://example.com/power_off?id=5',
                    headers: {Authorization: 'bearer token'},
                },
            },
            {
                label: 'Dismiss',
                http: {method: 'GET', url: 'http://example.com/dismiss', headers: {}},
            },
        ]);
    });

    it('defaults to POST', () => {
        const [button] = actionButtons(withButtons([{label: 'a', http: {url: 'https://x.y/'}}]));
        expect(button.http.method).toBe('POST');
    });

    it('skips invalid buttons', () => {
        const buttons = actionButtons(
            withButtons([
                null,
                'string',
                {label: '', http: {url: 'https://x.y/'}},
                {label: 'no http'},
                {label: 'no url', http: {}},
                {label: 'bad url', http: {url: 'not a url'}},
                {label: 'js', http: {url: 'javascript:alert(1)'}},
                {label: 'data', http: {url: 'data:text/plain,hi'}},
                {label: 'method', http: {method: 'CONNECT', url: 'https://x.y/'}},
                {label: 'ok', http: {url: 'https://x.y/'}},
            ])
        );
        expect(buttons.map((b) => b.label)).toEqual(['ok']);
    });

    it('drops non-string headers', () => {
        const [button] = actionButtons(
            withButtons([
                {label: 'a', http: {url: 'https://x.y/', headers: {A: 'b', C: 1, D: null}}},
            ])
        );
        expect(button.http.headers).toEqual({A: 'b'});
    });
});
