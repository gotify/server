import {Page} from 'puppeteer';
import {newTest, GotifyTest} from './setup';
import {clickByText, ClientCol, count, waitForExists, waitToDisappear} from './utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import * as auth from './authentication';
import * as selector from './selector';

let page: Page;
let gotify: GotifyTest;
beforeAll(async () => {
    gotify = await newTest();
    page = gotify.page;
});

afterAll(async () => await gotify.close());

const $clientTable = selector.table('#client-table');
const $clientDialog = selector.form('#client-dialog');

// This expects the session to be already elevated.
const cancelElevationViaUI = async (row: number) => {
    await page.goto(gotify.url + '/#/clients');
    await waitForExists(page, selector.heading(), 'Clients');

    await page.click($clientTable.cell(row, ClientCol.Elevate, '.elevate'));
    await page.waitForSelector('.elevate-client-dialog');

    await page.click('.elevate-client-dialog .elevate-duration [role=combobox]');
    await clickByText(page, '[role="option"]', 'Cancel elevation');
    await waitToDisappear(page, '[role="listbox"]');

    await page.click('.elevate-client-dialog .elevate-confirm');
    await waitToDisappear(page, '.elevate-client-dialog');
};

const elevateViaForm = async (password: string) => {
    const passwordInput = '.elevation-password input';
    await page.waitForSelector(passwordInput);
    await page.type(passwordInput, password);
    await page.click('.elevation-submit');
};

describe('Elevation', () => {
    it('does login', async () => await auth.login(page));

    describe('setup', () => {
        it('navigates to clients', async () => {
            await page.click('#navigate-clients');
            await waitForExists(page, selector.heading(), 'Clients');
        });
        it('creates a test client', async () => {
            await page.click('#create-client');
            await page.waitForSelector($clientDialog.selector());
            await page.type($clientDialog.input('.name'), 'test-client');
            await page.click($clientDialog.button('.create'));
            await waitToDisappear(page, $clientDialog.selector());
            await page.waitForSelector($clientTable.row(2));
            expect(await count(page, $clientTable.rows())).toBe(2);
        });
    });

    describe('Users page requires elevation', () => {
        it('de-elevates the current client via UI', () => cancelElevationViaUI(1));
        it('navigates to users and sees elevation form', async () => {
            await page.goto(gotify.url + '/#/users');
            await waitForExists(page, selector.heading(), 'Authentication Required');
            await page.waitForSelector('.elevation-password input');
        });
        it('elevates via password and sees users page', async () => {
            await elevateViaForm('admin');
            await waitForExists(page, selector.heading(), 'Users');
            expect(page.url()).toContain('/users');
        });
    });

    describe('Client delete requires elevation', () => {
        it('de-elevates the current client via UI', () => cancelElevationViaUI(1));
        it('navigates to clients', async () => {
            await page.goto(gotify.url + '/#/clients');
            await waitForExists(page, selector.heading(), 'Clients');
        });
        it('clicks delete and sees elevation form in dialog', async () => {
            await page.click($clientTable.cell(2, ClientCol.Delete, '.delete'));
            await page.waitForSelector(selector.$confirmDialog.selector());
            await page.waitForSelector('.confirm-dialog .elevation-password input');
        });
        it('elevates', () => elevateViaForm('admin'));
        it('confirms deletion', async () => {
            await page.waitForSelector(selector.$confirmDialog.button('.confirm'));
            await page.click(selector.$confirmDialog.button('.confirm'));
        });
        it('has deleted the client', async () => {
            await waitToDisappear(page, $clientTable.row(2));
            expect(await count(page, $clientTable.rows())).toBe(1);
        });
    });
});
