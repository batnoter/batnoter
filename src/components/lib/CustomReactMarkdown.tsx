import ContentCopyOutlinedIcon from '@mui/icons-material/ContentCopyOutlined';
import React, { ReactElement } from 'react';
import ReactMarkdown from 'react-markdown';
import { ReactMarkdownOptions } from 'react-markdown/lib/react-markdown';
import remarkGfm from 'remark-gfm';
import './CustomReactMarkdown.scss';

const CustomReactMarkdown: React.FC<ReactMarkdownOptions> = (props: ReactMarkdownOptions): ReactElement => {
  return (
    <ReactMarkdown {...props}
      components={{
        code({ inline, className, children, ...props }) {
          return (
            <>
              <code className={className} {...props}>{children}</code>
              {!inline && <ContentCopyOutlinedIcon style={{ right: 5, position: "absolute", cursor: 'pointer' }}
                onClick={() => { navigator.clipboard.writeText(String(children)) }} />}
            </>
          )
        }
      }}
      remarkPlugins={[remarkGfm]}>
      {props.children}
    </ReactMarkdown>
  )
}

export default CustomReactMarkdown;
