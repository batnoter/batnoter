import ContentCopyOutlinedIcon from '@mui/icons-material/ContentCopyOutlined';
import { styled } from '@mui/material';
import React, { ReactElement } from 'react';
import ReactMarkdown from 'react-markdown';
import { ReactMarkdownOptions } from 'react-markdown/lib/react-markdown';
import remarkGfm from 'remark-gfm';

const StyledReactMarkdown = styled(ReactMarkdown)(
  () => `
    position: relative;
    pre {
      display: flex;
    }
    svg {
      background-color: #f5f5f5;
      opacity: 0.3;
      &:hover {
        opacity: 1;
      }
    }
`,
);

const CustomReactMarkdown: React.FC<ReactMarkdownOptions> = (props: ReactMarkdownOptions): ReactElement => {
  return (
    <StyledReactMarkdown {...props}
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
    </StyledReactMarkdown>
  )
}

export default CustomReactMarkdown;
