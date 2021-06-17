import {Page} from 'puppeteer';
import {newTest, GotifyTest} from './setup';
import {count, innerText, waitForExists, waitToDisappear, clearField} from './utils';
import * as auth from './authentication';

import * as selector from './selector';

let page: Page;
let gotify: GotifyTest;
beforeAll(async () => {
    gotify = await newTest();
    page = gotify.page;
});

afterAll(async () => await gotify.close());

enum Col {
    Name = 1,
    Token = 2,
    Edit = 3,
    Delete = 4,
}

const hasClient =
    (name: string, row: number): (() => Promise<void>) =>
    async () => {
        expect(await innerText(page, $table.cell(row, Col.Name))).toBe(name);
    };

const updateClient =
    (id: number, data: {name?: string}): (() => Promise<void>) =>
    async () => {
        await page.click($table.cell(id, Col.Edit, '.edit'));
        await page.waitForSelector($dialog.selector());
        if (data.name) {
            const nameSelector = $dialog.input('.name');
            await clearField(page, nameSelector);
            await page.type(nameSelector, data.name);
        }
        await page.click($dialog.button('.update'));
        await waitToDisappear(page, $dialog.selector());
    };

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
        expect(await count(page, $table.rows())).toBe(1);
    });
    describe('create clients', () => {
        const createClient =
            (name: string): (() => Promise<void>) =>
            async () => {
                await page.click('#create-client');
                await page.waitForSelector($dialog.selector());
                await page.type($dialog.input('.name'), name);
                await page.click($dialog.button('.create'));
            };
        it('phone', createClient('phone'));
        it('desktop app', createClient('desktop app'));
    });
    it('has created clients', async () => {
        await page.waitForSelector($table.row(3));

        expect(await count(page, $table.rows())).toBe(3);

        expect(await innerText(page, $table.cell(1, Col.Name))).toContain('chrome');
        expect(await innerText(page, $table.cell(2, Col.Name))).toBe('phone');
        expect(await innerText(page, $table.cell(3, Col.Name))).toBe('desktop app');
    });
    it('updates client', updateClient(1, {name: 'firefox'}));
    it('has updated client name', hasClient('firefox', 1));
    it('shows token', async () => {
        await page.click($table.cell(3, Col.Token, '.toggle-visibility'));
        expect((await innerText(page, $table.cell(3, Col.Token))).startsWith('C')).toBeTruthy();
    });
    it('deletes client', async () => {
        await page.click($table.cell(2, Col.Delete, '.delete'));

        await page.waitForSelector(selector.$confirmDialog.selector());
        await page.click(selector.$confirmDialog.button('.confirm'));
    });
    it('has deleted client', async () => {
        await waitToDisappear(page, $table.row(3));

        expect(await count(page, $table.rows())).toBe(2);
    });
    // eslint-disable-next-line
    it('deletes own client', async () => {
        await page.click($table.cell(1, Col.Delete, '.delete'));

        // confirm delete
        await page.waitForSelector(selector.$confirmDialog.selector());
        await page.click(selector.$confirmDialog.button('.confirm'));
    });
    it('automatically logs out', async () => {
        await waitForExists(page, selector.heading(), 'Login');
    });
});
