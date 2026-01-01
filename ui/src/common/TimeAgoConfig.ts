import {Formatter} from 'react-timeago';
import {makeIntlFormatter} from 'react-timeago/defaultFormatter';

export const TimeAgoFormatter: Map<'long' | 'narrow', Formatter> = new Map([
    ['long', makeIntlFormatter({style: 'long', locale: 'en'})],
    ['narrow', makeIntlFormatter({style: 'narrow', locale: 'en'})],
]);
