// todo before all tests jest start puppeteer
import {Page} from 'puppeteer';
import {newTest, GotifyTest} from './setup';
import {clickByText, count, innerText, waitForCount, waitForExists} from './utils';
import * as auth from './authentication';
import * as selector from './selector';
import axios from 'axios';
import {IApplication, IMessage, IMessageExtras} from '../types';

let page: Page;
let gotify: GotifyTest;
beforeAll(async () => {
    gotify = await newTest();
    page = gotify.page;
});

afterAll(async () => await gotify.close());

// eslint-disable-next-line
const axiosAuth = {auth: {username: 'admin', password: 'admin'}};

let windowsServerToken: string;
let linuxServerToken: string;
let backupServerToken: string;

const naviId = '#message-navigation';

interface Msg {
    message: string;
    title: string;
}

const navigate = async (appName: string) => {
    await clickByText(page, 'a', appName);
    await waitForExists(page, selector.heading(), appName);
};

// eslint-disable-next-line
describe('Messages', () => {
    it('does login', async () => await auth.login(page));
    it('is on messages', async () => {
        await waitForExists(page, selector.heading(), 'All Messages');
    });
    it('has url', async () => {
        expect(page.url()).toContain('/');
    });
    const createApp = (name: string) =>
        axios
            .post<IApplication>(`${gotify.url}/application`, {name}, axiosAuth)
            .then((resp) => resp.data.token);
    it('shows navigation', async () => {
        await page.waitForSelector(naviId);
    });
    it('has all messages button', async () => {
        await page.waitForSelector(`${naviId} .all`);
    });
    it('has no applications', async () => {
        expect(await count(page, `${naviId} .item`)).toBe(0);
    });
    describe('create apps', () => {
        it('Windows', async () => {
            windowsServerToken = await createApp('Windows');
            await page.reload();
            await waitForExists(page, 'a', 'Windows');
        });
        it('Backup', async () => {
            backupServerToken = await createApp('Backup');
            await page.reload();
            await waitForExists(page, 'a', 'Backup');
        });
        it('Linux', async () => {
            linuxServerToken = await createApp('Linux');
            await page.reload();
            await waitForExists(page, 'a', 'Linux');
        });
    });
    it('has three applications', async () => {
        expect(await count(page, `${naviId} .item`)).toBe(3);
    });
    it('changes url when navigating to application', async () => {
        await navigate('Windows');
        expect(page.url()).toContain('/messages/1');
        await navigate('All Messages');
    });
    it('has no messages', async () => {
        expect(await count(page, '#messages')).toBe(0);
    });
    it('has no messages in app', async () => {
        await navigate('Windows');
        expect(await count(page, '#messages')).toBe(0);
        await navigate('All Messages');
    });

    const extractMessages = async (expectCount: number) => {
        await waitForCount(page, '#messages .message', expectCount);
        const messages = await page.$$(`#messages .message`);
        const result: Msg[] = [];
        for (const item of messages) {
            const message = await innerText(item, '.content');
            const title = await innerText(item, '.title');
            result.push({message, title});
        }
        return result;
    };
    const m = (title: string, message: string, extras?: IMessageExtras) => ({
        title,
        message,
        extras,
    });

    const windows1 = m('Login', 'User jmattheis logged in.');
    const windows2 = m('Shutdown', 'Windows will be shut down.');
    const windows3 = m('Login', 'User nicories logged in.');

    const linux1 = m('SSH-Login', 'root@127.0.0.1 did a ssh login.');
    const linux2 = m('Reboot', 'Linux server just rebooted.');
    const linux3 = m('SSH-Login', 'jmattheis@localhost did a ssh login.');

    const backup1 = m('Backup done', 'Linux Server Backup finished (1.6GB).');
    const backup2 = m('Backup done', 'Windows Server Backup finished (6.2GB).');
    const backup3 = m('Backup done', 'Gotify Backup finished (0.1MB).');

    const createMessage = (msg: Partial<IMessage>, token: string) =>
        axios.post<IMessage>(`${gotify.url}/message`, msg, {
            headers: {'X-Gotify-Key': token},
        });

    const expectMessages = async (toCheck: {
        all: Msg[];
        windows: Msg[];
        linux: Msg[];
        backup: Msg[];
    }) => {
        await navigate('All Messages');
        expect(await extractMessages(toCheck.all.length)).toEqual(toCheck.all);
        await navigate('Windows');
        expect(await extractMessages(toCheck.windows.length)).toEqual(toCheck.windows);
        await navigate('Linux');
        expect(await extractMessages(toCheck.linux.length)).toEqual(toCheck.linux);
        await navigate('Backup');
        expect(await extractMessages(toCheck.backup.length)).toEqual(toCheck.backup);
        await navigate('All Messages');
    };

    it('create a message', async () => {
        await createMessage(windows1, windowsServerToken);
        expect(await extractMessages(1)).toEqual([windows1]);
    });
    it('has one message in windows app', async () => {
        await navigate('Windows');
        expect(await extractMessages(1)).toEqual([windows1]);
    });
    it('has no message in linux app', async () => {
        await navigate('Linux');
        expect(await extractMessages(0)).toEqual([]);
        await navigate('All Messages');
    });
    describe('add some messages', () => {
        it('1', async () => {
            await createMessage(windows2, windowsServerToken);
            await expectMessages({
                all: [windows2, windows1],
                windows: [windows2, windows1],
                linux: [],
                backup: [],
            });
        });
        it('2', async () => {
            await createMessage(linux1, linuxServerToken);
            await expectMessages({
                all: [linux1, windows2, windows1],
                windows: [windows2, windows1],
                linux: [linux1],
                backup: [],
            });
        });
        it('3', async () => {
            await createMessage(backup1, backupServerToken);
            await expectMessages({
                all: [backup1, linux1, windows2, windows1],
                windows: [windows2, windows1],
                linux: [linux1],
                backup: [backup1],
            });
        });
        it('4', async () => {
            await createMessage(windows3, windowsServerToken);
            await expectMessages({
                all: [windows3, backup1, linux1, windows2, windows1],
                windows: [windows3, windows2, windows1],
                linux: [linux1],
                backup: [backup1],
            });
        });
        it('5', async () => {
            await createMessage(linux2, linuxServerToken);
            await expectMessages({
                all: [linux2, windows3, backup1, linux1, windows2, windows1],
                windows: [windows3, windows2, windows1],
                linux: [linux2, linux1],
                backup: [backup1],
            });
        });
    });
    it('deletes a windows message', async () => {
        await navigate('Windows');
        await page.evaluate(() =>
            (
                document.querySelectorAll('#messages .message .delete')[1] as HTMLButtonElement
            ).click()
        );
        await expectMessages({
            all: [linux2, windows3, backup1, linux1, windows1],
            windows: [windows3, windows1],
            linux: [linux2, linux1],
            backup: [backup1],
        });
    });
    it('deletes all linux messages', async () => {
        await navigate('Linux');
        await page.click('#delete-all');
        await page.waitForSelector(selector.$confirmDialog.selector());
        await page.click(selector.$confirmDialog.button('.confirm'));
        await page.waitForSelector('#delete-all:disabled');
        await expectMessages({
            all: [windows3, backup1, windows1],
            windows: [windows3, windows1],
            linux: [],
            backup: [backup1],
        });
    });
    describe('add some more messages', () => {
        it('1', async () => {
            await createMessage(linux3, linuxServerToken);
            await expectMessages({
                all: [linux3, windows3, backup1, windows1],
                windows: [windows3, windows1],
                linux: [linux3],
                backup: [backup1],
            });
        });
        it('2', async () => {
            await createMessage(backup2, backupServerToken);
            await expectMessages({
                all: [backup2, linux3, windows3, backup1, windows1],
                windows: [windows3, windows1],
                linux: [linux3],
                backup: [backup2, backup1],
            });
        });
    });
    it('deletes all messages', async () => {
        await navigate('All Messages');
        await page.click('#delete-all');
        await page.waitForSelector(selector.$confirmDialog.selector());
        await page.click(selector.$confirmDialog.button('.confirm'));
        await page.waitForSelector('#delete-all:disabled');
        await expectMessages({
            all: [],
            windows: [],
            linux: [],
            backup: [],
        });
    });
    it('adds one last message', async () => {
        await createMessage(backup3, backupServerToken);
        await expectMessages({
            all: [backup3],
            windows: [],
            linux: [],
            backup: [backup3],
        });
    });
    it('deletes all backup messages and navigates to all messages', async () => {
        await navigate('Backup');
        await page.click('#delete-all');
        await page.waitForSelector(selector.$confirmDialog.selector());
        await page.click(selector.$confirmDialog.button('.confirm'));
        await page.waitForSelector('#delete-all:disabled');
        await navigate('All Messages');
        await createMessage(backup3, backupServerToken);
        await waitForExists(page, '.message .title', backup3.title);
        expect(await extractMessages(1)).toEqual([backup3]);
    });
    it('does logout', async () => await auth.logout(page));
});
