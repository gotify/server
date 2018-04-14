import ReactList from 'react-list';

// See also https://github.com/coderiety/react-list/blob/master/react-list.es6
class FixedReactList extends ReactList {
    // deleting a messages or adding a message (per shift) requires invalidating the cache, react-list sucks as it does
    // not provide such functionality, therefore we need to hack it inside there :(
    ignoreNextCacheUpdate = false;

    cacheSizes() {
        if (this.ignoreNextCacheUpdate) {
            this.ignoreNextCacheUpdate = false;
            return;
        }
        super.cacheSizes();
    }

    clearCacheFromIndex(startIndex) {
        this.ignoreNextCacheUpdate = true;

        if (startIndex === 0) {
            this.cache = {};
        } else {
            Object.keys(this.cache).filter((index) => index >= startIndex).forEach((index) => {
                delete this.cache[index];
            });
        }
    };

    componentDidUpdate() {
        const hasCacheForLastRenderedItem = Object.keys(this.cache).length && this.cache[this.getVisibleRange()[1]];
        super.componentDidUpdate();
        if (!hasCacheForLastRenderedItem) {
            // when there is no cache for the last rendered item, then its a new item, react-list doesn't know it size
            // and cant correctly calculate the height of the list, we force a rerender where react-list knows the size.
            this.forceUpdate();
        }
    }
}

export default FixedReactList;
