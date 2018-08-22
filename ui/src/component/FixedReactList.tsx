// @ts-ignore
import ReactList from 'react-list';

// See also https://github.com/coderiety/react-list/blob/master/react-list.es6
class FixedReactList extends ReactList {
    // deleting a messages or adding a message (per shift) requires invalidating the cache, react-list sucks as it does
    // not provide such functionality, therefore we need to hack it inside there :(
    public ignoreNextCacheUpdate = false;

    public cacheSizes(): void {
        if (this.ignoreNextCacheUpdate) {
            this.ignoreNextCacheUpdate = false;
            return;
        }
        // @ts-ignore accessing private member
        super.cacheSizes();
    }

    public clearCacheFromIndex(startIndex: number): void {
        this.ignoreNextCacheUpdate = true;

        if (startIndex === 0) {
            // @ts-ignore accessing private member
            this.cache = {};
        } else {
            // @ts-ignore accessing private member
            Object.keys(this.cache)
                .filter((index) => +index >= startIndex)
                .forEach((index) => {
                    // @ts-ignore accessing private member
                    delete this.cache[index];
                });
        }
    }

    public componentDidUpdate() {
        const hasCacheForLastRenderedItem =
            // @ts-ignore accessing private member
            Object.keys(this.cache).length && this.cache[this.getVisibleRange()[1]];
        // @ts-ignore accessing private member
        super.componentDidUpdate();
        if (!hasCacheForLastRenderedItem) {
            // when there is no cache for the last rendered item, then its a new item, react-list doesn't know it size
            // and cant correctly calculate the height of the list, we force a rerender where react-list knows the size.
            this.forceUpdate();
        }
    }
}

export default FixedReactList;
