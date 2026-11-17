import {describe, expect, it} from 'vitest';
import {actionButtons} from './extras';

const button = (http: object) => ({
    label: 'test',
    http: {method: 'POST', url: 'https://x.y/', ...http},
});
const parse = (...buttons: unknown[]) => actionButtons({'client::action': {buttons}});

describe('actionButtons', () => {
    it('returns nothing without client::action', () => {
        expect(actionButtons(undefined)).toEqual({buttons: []});
        expect(actionButtons({})).toEqual({buttons: []});
        expect(actionButtons({'client::action': null})).toEqual({buttons: []});
    });

    it.each([{}, {buttons: 'nope'}])('reports a malformed client::action %j', (action) => {
        const result = actionButtons({'client::action': action});
        expect(result.buttons).toEqual([]);
        expect(result.error).toContain('→ at ["client::action"].buttons');
    });

    it('parses the format from gotify/server#494', () => {
        const result = parse(
            {
                label: 'Power off',
                http: {
                    method: 'post',
                    url: 'https://example.com/power_off?id=5',
                    headers: {Authorization: 'bearer token'},
                },
            },
            {label: 'Dismiss', http: {method: 'get', url: 'http://example.com/dismiss'}}
        );
        expect(result).toEqual({
            buttons: [
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
            ],
        });
    });

    it('reports every invalid button and keeps the valid ones', () => {
        const result = parse(
            button({url: 'ftp://x.y/'}),
            {label: 'ok', http: {method: 'POST', url: 'https://x.y/'}},
            button({method: 1})
        );
        expect(result.buttons.map((b) => b.label)).toEqual(['ok']);
        expect(result.error).toContain('→ at ["client::action"].buttons[0].http.url');
        expect(result.error).toContain('→ at ["client::action"].buttons[2].http.method');
    });

    it.each([
        ['null', null],
        ['a string', 'string'],
        ['an empty label', {label: '', http: {method: 'POST', url: 'https://x.y/'}}],
        ['a blank label', {label: '  ', http: {method: 'POST', url: 'https://x.y/'}}],
        ['a missing http', {label: 'no http'}],
        ['a missing url', button({url: undefined})],
        ['a missing method', button({method: undefined})],
        ['an invalid url', button({url: 'not a url'})],
        ['a javascript url', button({url: 'javascript:alert(1)'})],
        ['a data url', button({url: 'data:text/plain,hi'})],
        ['a non-string header', button({headers: {A: 'b', C: 1}})],
        ['a non-string body', button({body: {a: 1}})],
    ])('rejects a button with %s', (_, invalid) => {
        const result = parse(invalid);
        expect(result.buttons).toEqual([]);
        expect(result.error).toContain('→ at ["client::action"].buttons[0]');
    });
});
