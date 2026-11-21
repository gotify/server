import {Page} from 'puppeteer';
import {newTest, GotifyTest} from './setup';
import {clickByText, waitForExists} from './utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import * as auth from './authentication';
import axios from 'axios';
import http from 'http';
import {AddressInfo} from 'net';
import {IApplication} from '../types';

interface ReceivedRequest {
    method: string;
    url: string;
    authorization?: string;
}

let page: Page;
let gotify: GotifyTest;
let target: http.Server;
let targetUrl: string;
const received: ReceivedRequest[] = [];

beforeAll(async () => {
    gotify = await newTest();
    page = gotify.page;

    target = http.createServer((req, res) => {
        res.setHeader('Access-Control-Allow-Origin', '*');
        res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, PATCH, DELETE');
        res.setHeader('Access-Control-Allow-Headers', 'Authorization');
        if (req.method === 'OPTIONS') {
            res.writeHead(204).end();
            return;
        }
        received.push({
            method: req.method ?? '',
            url: req.url ?? '',
            authorization: req.headers.authorization,
        });
        res.writeHead(req.url?.startsWith('/fail') ? 500 : 200).end();
    });
    await new Promise<void>((resolve) => target.listen(0, '127.0.0.1', resolve));
    targetUrl = `http://127.0.0.1:${(target.address() as AddressInfo).port}`;
});

afterAll(async () => {
    await gotify.close();
    await new Promise((resolve) => target.close(resolve));
});

const axiosAuth = {auth: {username: 'admin', password: 'admin'}};

const snackbarWithText = (text: string) =>
    page.waitForFunction(
        (t) =>
            Array.from(document.querySelectorAll('#notistack-snackbar')).some((e) =>
                (e.textContent ?? '').includes(t)
            ),
        {},
        text
    );

describe('Message actions', () => {
    it('does login', async () => await auth.login(page));

    it('shows action buttons of a message', async () => {
        const token = await axios
            .post<IApplication>(`${gotify.url}/application`, {name: 'Heater'}, axiosAuth)
            .then((resp) => resp.data.token);
        await axios.post(
            `${gotify.url}/message`,
            {
                title: 'Heater',
                message: 'The heater is still powered on.',
                extras: {
                    'client::action': {
                        buttons: [
                            {
                                label: 'Power off',
                                http: {
                                    method: 'post',
                                    url: `${targetUrl}/power_off?id=5`,
                                    headers: {Authorization: 'bearer token'},
                                },
                            },
                            {label: 'Broken', http: {method: 'post', url: `${targetUrl}/fail`}},
                            {label: 'Ignored', http: {url: 'javascript:alert(1)'}},
                        ],
                    },
                },
            },
            {headers: {'X-Gotify-Key': token}}
        );
        await waitForExists(page, '.message .title', 'Heater');
        const labels = await page.$$eval('.message .actions .action', (els) =>
            els.map((e) => e.textContent)
        );
        expect(labels).toEqual(['Power off', 'Broken']);
    });

    it('sends the request on click', async () => {
        await clickByText(page, '.message .actions .action', 'Power off');
        await snackbarWithText("Action 'Power off' succeeded.");
        expect(received).toEqual([
            {method: 'POST', url: '/power_off?id=5', authorization: 'bearer token'},
        ]);
    });

    it('reports failed requests', async () => {
        await clickByText(page, '.message .actions .action', 'Broken');
        await snackbarWithText("Action 'Broken' failed: 500");
        expect(received.map((r) => r.url)).toEqual(['/power_off?id=5', '/fail']);
    });

    it('does logout', async () => await auth.logout(page));
});
