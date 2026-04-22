import {Formatter} from 'react-timeago';
import {makeIntlFormatter} from 'react-timeago/defaultFormatter';

const longFormatter = makeIntlFormatter({style: 'long', locale: 'en'});

const longMinutesFormatter: Formatter = (value, unit, ...rest) =>
    unit === 'second' ? longFormatter(1, 'minute', ...rest) : longFormatter(value, unit, ...rest);

export const TimeAgoFormatter = {
    long: longFormatter,
    narrow: makeIntlFormatter({style: 'narrow', locale: 'en'}),
    longMinutes: longMinutesFormatter,
} as const satisfies Record<string, Formatter>;
