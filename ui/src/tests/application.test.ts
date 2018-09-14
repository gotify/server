import {Page} from 'puppeteer';
import {newTest, GotifyTest} from './setup';
import {count, innerText, waitForExists, waitToDisappear} from './utils';
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
    Name = 2,
    Token = 3,
    Description = 4,
    EditDelete = 5,
}

const hiddenToken = '•••••••••••••••';

const $table = selector.table('#app-table');
const $dialog = selector.form('#app-dialog');

export const createApp = (name: string, description: string): (() => Promise<void>) => {
    return async () => {
        await page.click('#create-app');
        await page.waitForSelector($dialog.selector());
        await page.type($dialog.input('.name'), name);
        await page.type($dialog.textarea('.description'), description);
        await page.click($dialog.button('.create'));
    };
};

describe('Application', () => {
    it('does login', async () => await auth.login(page));
    it('navigates to applications', async () => {
        await page.click('#navigate-apps');
        await waitForExists(page, selector.heading(), 'Applications');
    });
    it('has changed url', async () => {
        expect(page.url()).toContain('/applications');
    });
    it('does not have any applications', async () => {
        expect(await count(page, $table.rows())).toBe(0);
    });
    describe('create apps', () => {
        it('server', createApp('server', '#1'));
        it('desktop', createApp('desktop', '#2'));
        it('raspberry', createApp('raspberry', '#3'));
    });
    describe('has created apps', () => {
        const hasApp = (name: string, description: string, row: number): (() => Promise<void>) => {
            return async () => {
                expect(await innerText(page, $table.cell(row, Col.Name))).toBe(name);
                expect(await innerText(page, $table.cell(row, Col.Token))).toBe(hiddenToken);
                expect(await innerText(page, $table.cell(row, Col.Description))).toBe(description);
            };
        };

        it('has three apps', async () => {
            await page.waitForSelector($table.row(3));
            expect(await count(page, $table.rows())).toBe(3);
        });
        it('has server app', hasApp('server', '#1', 1));
        it('has desktop app', hasApp('desktop', '#2', 2));
        it('has raspberry app', hasApp('raspberry', '#3', 3));
        it('shows token', async () => {
            await page.click($table.cell(3, Col.Token, '.toggle-visibility'));
            const token = await innerText(page, $table.cell(3, Col.Token));
            expect(token.startsWith('A')).toBeTruthy();
        });
    });
    it('deletes application', async () => {
        await page.click($table.cell(2, Col.EditDelete, '.delete'));

        await page.waitForSelector(selector.$confirmDialog.selector());
        await page.click(selector.$confirmDialog.button('.confirm'));
    });
    it('has deleted application', async () => {
        await waitToDisappear(page, $table.row(3));
        expect(await count(page, $table.rows())).toBe(2);
    });
    it('does logout', async () => await auth.logout(page));
});
