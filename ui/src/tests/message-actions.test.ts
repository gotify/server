import {Page} from 'puppeteer';
import {newTest, GotifyTest} from './setup';
import {clickByText, waitForExists} from './utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import * as auth from './authentication';
import {IApplication} from '../types';

let page: Page;
let gotify: GotifyTest;
let app: IApplication;

const adminHeaders = {Authorization: `Basic ${btoa('admin:admin')}`};

beforeAll(async () => {
    gotify = await newTest();
    page = gotify.page;
});

afterAll(async () => await gotify.close());

describe('Message actions', () => {
    it('creates an app', async () => {
        app = await fetch(`${gotify.url}/application`, {
            method: 'POST',
            headers: {...adminHeaders, 'Content-Type': 'application/json'},
            body: JSON.stringify({name: 'Heater'}),
        }).then((resp) => resp.json());
    });

    it('does login', async () => await auth.login(page));

    it('creates a message with action buttons', async () => {
        await fetch(`${gotify.url}/message`, {
            method: 'POST',
            headers: {'X-Gotify-Key': app.token, 'Content-Type': 'application/json'},
            body: JSON.stringify({
                title: 'Heater',
                message: 'The heater is still powered on.',
                extras: {
                    'client::action': {
                        buttons: [
                            {
                                label: 'Power off',
                                http: {
                                    method: 'post',
                                    url: `${gotify.url}/message`,
                                    headers: {
                                        'X-Gotify-Key': app.token,
                                        'Content-Type': 'application/json',
                                    },
                                    body: JSON.stringify({
                                        title: 'Powered off',
                                        message: 'The heater is off.',
                                    }),
                                },
                            },
                            {label: 'Invalid', http: {method: 'post', url: 'javascript:alert(1)'}},
                            {label: 'Broken', http: {method: 'post', url: `${gotify.url}/message`}},
                            {label: 'Version', open: {url: `${gotify.url}/version`}},
                        ],
                    },
                },
            }),
        });
    });

    it('shows the valid action buttons of a message', async () => {
        await waitForExists(page, '.message .title', 'Heater');
        const labels = await page.$$eval('.message .actions .action', (els) =>
            els.map((e) => e.textContent)
        );
        expect(labels).toEqual(['Power off', 'Broken', 'Version']);
    });

    it('shows the parsing error of invalid buttons', async () => {
        await clickByText(page, '.message .action-error-toggle', 'Details');
        const details = await page.waitForSelector('.message .action-error-details');
        expect(await details?.evaluate((e) => e.textContent)).toContain(
            '→ at ["client::action"].buttons[1].http.url'
        );
    });

    it('sends the request on click', async () => {
        await clickByText(page, '.message .actions .action', 'Power off');
        await waitForExists(page, '#notistack-snackbar', "Action 'Power off' succeeded.");
        await waitForExists(page, '.message .title', 'Powered off');
    });

    it('reports failed requests', async () => {
        await clickByText(page, '.message .actions .action', 'Broken');
        await waitForExists(
            page,
            '#notistack-snackbar',
            "Action 'Broken' failed: 401 Unauthorized"
        );
    });

    it('opens the url of an open button in a new tab', async () => {
        const url = `${gotify.url}/version`;
        await clickByText(page, '.message .actions .action', 'Version');
        const target = await gotify.browser.waitForTarget((t) => t.url() === url);
        const tab = await target.page();
        expect(await tab?.evaluate(() => window.opener)).toBeNull();
        await tab?.close();
    });

    it('does logout', () => auth.logout(page));

    it('deletes the app', async () => {
        const resp = await fetch(`${gotify.url}/application/${app.id}`, {
            method: 'DELETE',
            headers: adminHeaders,
        });
        expect(resp.ok).toBe(true);
    });
});
