import {Formatter} from 'react-timeago';
import {makeIntlFormatter} from 'react-timeago/defaultFormatter';

export const TimeAgoFormatter: Record<'long' | 'narrow', Formatter> = {
    long: makeIntlFormatter({style: 'long', locale: 'en'}),
    narrow: makeIntlFormatter({style: 'narrow', locale: 'en'}),
};
