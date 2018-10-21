import {ElementHandle, Page} from 'puppeteer';

export const innerText = async (page: ElementHandle | Page, selector: string): Promise<string> => {
    const element = await page.$(selector);
    const handle = await element!.getProperty('innerText');
    const value = await handle.jsonValue();
    return value.toString().trim();
};

export const clickByText = async (page: Page, selector: string, text: string): Promise<void> => {
    await waitForExists(page, selector, text);
    text = text.toLowerCase();
    await page.evaluate(
        (_selector, _text) => {
            Array.from(document.querySelectorAll(_selector))
                .filter((element) => element.textContent.toLowerCase().trim() === _text)[0]
                .click();
        },
        selector,
        text
    );
};

export const count = async (page: Page, selector: string): Promise<number> => {
    return page.$$(selector).then((elements) => elements.length);
};

export const waitToDisappear = async (page: Page, selector: string): Promise<void> => {
    return page.waitForFunction(
        (_selector: string) => !document.querySelector(_selector),
        {},
        selector
    );
};

export const waitForCount = async (page: Page, selector: string, amount: number): Promise<void> => {
    return page.waitForFunction(
        (_selector: string, _amount: number) =>
            document.querySelectorAll(_selector).length === _amount,
        {},
        selector,
        amount
    );
};

export const waitForExists = async (page: Page, selector: string, text: string): Promise<void> => {
    text = text.toLowerCase();
    await page.waitForFunction(
        (_selector: string, _text: string) => {
            return (
                Array.from(document.querySelectorAll(_selector)).filter(
                    (element) => element.textContent!.toLowerCase().trim() === _text
                ).length > 0
            );
        },
        {},
        selector,
        text
    );
};

export const clearField = async (element: ElementHandle | Page, selector: string) => {
    const elementHandle = await element.$(selector);
    if (!elementHandle) {
        fail();
        return;
    }
    await elementHandle.click();
    await elementHandle.focus();
    // click three times to select all
    await elementHandle.click({clickCount: 3});
    await elementHandle.press('Backspace');
};
