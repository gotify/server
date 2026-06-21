export const copyToClipboard = async (value: string) => {
    try {
        await navigator.clipboard.writeText(value);
    } catch (error) {
        console.warn('Failed to copy to clipboard using Clipboard API:', error);
        const elem = document.createElement('textarea');
        elem.value = value;
        document.body.appendChild(elem);
        elem.select();
        document.execCommand('copy');
        document.body.removeChild(elem);
    }
};
