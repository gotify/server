export const copyToClipboard = async (value: string) => {
    try {
        await navigator.clipboard.writeText(value);
    } catch (error) {
        console.error('Failed to copy to clipboard:', error);
        try {
            const elem = document.createElement('textarea');
            elem.value = value;
            document.body.appendChild(elem);
            elem.select();
            document.execCommand('copy');
            document.body.removeChild(elem);
        } catch (error) {
            console.error('Failed to copy to clipboard (fallback):', error);
        }
    }
};
