export function isMacOs() {
    const navigator = window.navigator;

    if ('userAgentData' in navigator) {
        /** @type {{platform: string | undefined} | undefined} */
        const userAgentData = navigator.userAgentData;
        if (userAgentData?.platform?.toLowerCase().includes("mac")) {
            return true;
        }
    }

    return navigator.userAgent.toLowerCase().includes("mac");
}