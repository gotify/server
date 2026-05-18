import {Page} from 'puppeteer';
import {newTest, GotifyTest} from './setup';
import {count, innerText, waitForExists, waitToDisappear, clearField, ClientCol} from './utils';
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

const waitForClient =
    (name: string, row: number): (() => Promise<void>) =>
    async () => {
        await waitForExists(page, $table.cell(row, ClientCol.Name), name);
    };

interface ClientFields {
    name?: string;
    expiresAfter?: number;
}

const fillClientDialog =
    (opener: string, submit: string, data: ClientFields): (() => Promise<void>) =>
    async () => {
        await page.click(opener);
        await page.waitForSelector($dialog.selector());
        if (data.name !== undefined) {
            const nameSelector = $dialog.input('.name');
            await clearField(page, nameSelector);
            await page.type(nameSelector, data.name);
        }
        if (data.expiresAfter !== undefined) {
            const expiresSelector = $dialog.input('.expires-after');
            await clearField(page, expiresSelector);
            await page.type(expiresSelector, data.expiresAfter.toString());
        }
        await page.click($dialog.button(submit));
        await waitToDisappear(page, $dialog.selector());
    };

const createClient = (data: ClientFields) => fillClientDialog('#create-client', '.create', data);

const updateClient = (id: number, data: ClientFields) =>
    fillClientDialog($table.cell(id, ClientCol.Edit, '.edit'), '.update', data);

const $table = selector.table('#client-table');
const $dialog = selector.form('#client-dialog');

describe('Client', () => {
    it('does login', async () => await auth.login(page));
    it('navigates to clients', async () => {
        await page.click('#navigate-clients');
        await waitForExists(page, selector.heading(), 'Clients');
    });
    it('has changed url', async () => {
        expect(page.url()).toContain('/clients');
    });
    it('has one client (the current session)', async () => {
        await page.waitForSelector($table.row(1));
        expect(await count(page, $table.rows())).toBe(1);
    });
    describe('create clients', () => {
        it('phone', createClient({name: 'phone'}));
        it('desktop app', createClient({name: 'desktop app', expiresAfter: 60 * 60}));
    });
    it('has created clients', async () => {
        await page.waitForSelector($table.row(3));

        expect(await count(page, $table.rows())).toBe(3);

        expect(await innerText(page, $table.cell(1, ClientCol.Name))).toContain('chrome');
        expect(await innerText(page, $table.cell(2, ClientCol.Name))).toBe('phone');
        expect(await innerText(page, $table.cell(3, ClientCol.Name))).toBe('desktop app');
    });
    it('shows expires after for new clients', async () => {
        expect(await innerText(page, $table.cell(2, ClientCol.ExpiresIn))).toBe('-');
        expect(await innerText(page, $table.cell(3, ClientCol.ExpiresIn))).toBe('59m');
    });
    it('updates client', updateClient(1, {name: 'firefox', expiresAfter: 60 * 60 * 10}));
    it('has updated client name', waitForClient('firefox', 1));
    it('has updated expires after', async () => {
        expect(await innerText(page, $table.cell(1, ClientCol.ExpiresIn))).toBe('9h 59m');
    });
    it('shows token', async () => {
        await page.click($table.cell(3, ClientCol.Token, '.toggle-visibility'));
        expect(
            (await innerText(page, $table.cell(3, ClientCol.Token))).startsWith('C')
        ).toBeTruthy();
    });
    it('shows last seen', async () => {
        expect(await innerText(page, $table.cell(3, ClientCol.LastSeen))).toBeTruthy();
    });
    it('deletes client', async () => {
        await page.click($table.cell(2, ClientCol.Delete, '.delete'));

        await page.waitForSelector(selector.$confirmDialog.selector());
        await page.click(selector.$confirmDialog.button('.confirm'));
    });
    it('has deleted client', async () => {
        await waitToDisappear(page, $table.row(3));

        expect(await count(page, $table.rows())).toBe(2);
    });
    it('deletes own client', async () => {
        await page.click($table.cell(1, ClientCol.Delete, '.delete'));

        // confirm delete
        await page.waitForSelector(selector.$confirmDialog.selector());
        await page.click(selector.$confirmDialog.button('.confirm'));
    });
    it('automatically logs out', async () => {
        await waitForExists(page, selector.heading(), 'Login');
    });
});
