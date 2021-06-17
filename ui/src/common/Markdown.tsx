import React from 'react';
import ReactMarkdown from 'react-markdown';
import gfm from 'remark-gfm';

export const Markdown = ({children}: {children: string}) => (
    <ReactMarkdown plugins={[gfm]}>{children}</ReactMarkdown>
);
