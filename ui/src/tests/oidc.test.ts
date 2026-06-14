import axios from 'axios';
import {Browser, Page} from 'puppeteer';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import {newTest, GotifyTest} from './setup';
import {DEX_PASSWORD, DexUser} from './dex';
import {waitForExists} from './utils';
import * as selector from './selector';
import * as auth from './authentication';

const linkUser: DexUser = {email: 'link@gotify.net', username: 'linkuser', userID: 'id-link'};
const dupUser1: DexUser = {email: 'dup1@gotify.net', username: 'dupuser', userID: 'id-dup-1'};
const dupUser2: DexUser = {email: 'dup2@gotify.net', username: 'dupuser', userID: 'id-dup-2'};

const createLocalUser = async (url: string, name: string, pass: string): Promise<void> =>
    axios.post(
        `${url}/user`,
        {name, pass, admin: false},
        {auth: {username: 'admin', password: 'admin'}}
    );

const loginWithOIDC = async (page: Page, user: DexUser): Promise<void> => {
    await waitForExists(page, selector.heading(), 'Login');
    const href = await page.$eval('#oidc-login', (a) => (a as HTMLAnchorElement).href);
    await page.goto(href);

    await page.waitForSelector('#login');
    await page.type('#login', user.email);
    await page.type('#password', DEX_PASSWORD);
    await page.click('#submit-login');
};

const expectLoggedIn = async (page: Page): Promise<void> => {
    await waitForExists(page, selector.heading(), 'All Messages');
    await page.waitForSelector('#logout');
};

const oidcError = async (page: Page): Promise<string> => {
    await page.waitForFunction(() => document.location.pathname.includes('/auth/oidc/callback'));
    return page.evaluate(() => document.body.innerText);
};

const clearSession = async (browser: Browser): Promise<void> => {
    await browser.deleteCookie(...(await browser.cookies()));
};

describe('OIDC login of an existing local user without link-by-username', () => {
    let gotify: GotifyTest;
    let page: Page;
    beforeAll(async () => {
        gotify = await newTest('', {
            oidc: {autoRegister: true, linkByUsername: false, users: [linkUser]},
        });
        page = gotify.page;
        await createLocalUser(gotify.url, linkUser.username, 'localpass');
    });
    afterAll(async () => await gotify.close());

    it('rejects the oidc login because linking is disabled', async () => {
        await loginWithOIDC(page, linkUser);
        expect(await oidcError(page)).toContain(
            `a local user with the username ${linkUser.username} already exists and linking by username is disabled`
        );
    });

    it('still allows the local user to log in with a password', async () => {
        await page.goto(gotify.url);
        await auth.login(page, 'linkuser', 'localpass');
        await auth.logout(page);
    });
});

describe('OIDC login of an existing local user with link-by-username', () => {
    let gotify: GotifyTest;
    let page: Page;
    beforeAll(async () => {
        gotify = await newTest('', {
            oidc: {autoRegister: true, linkByUsername: true, users: [linkUser]},
        });
        page = gotify.page;
        await createLocalUser(gotify.url, linkUser.username, 'localpass');
    });
    afterAll(async () => await gotify.close());

    it('links the existing local user and logs in', async () => {
        await loginWithOIDC(page, linkUser);
        await expectLoggedIn(page);
    });
});

describe('OIDC login with two identities sharing the same username', () => {
    let gotify: GotifyTest;
    let page: Page;
    beforeAll(async () => {
        gotify = await newTest('', {
            oidc: {autoRegister: true, linkByUsername: true, users: [dupUser1, dupUser2]},
        });
        page = gotify.page;
    });
    afterAll(async () => await gotify.close());

    it('auto-registers the first identity', async () => {
        await loginWithOIDC(page, dupUser1);
        await expectLoggedIn(page);
    });

    it('clears session', () => clearSession(gotify.browser));
    it('rejects the second identity with same username', async () => {
        await page.goto(gotify.url);
        await loginWithOIDC(page, dupUser2);
        expect(await oidcError(page)).toContain(
            `the user ${dupUser2.username} is already bound to a different OIDC identity`
        );
    });
});
