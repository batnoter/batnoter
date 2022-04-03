
import ContentCopyOutlinedIcon from '@mui/icons-material/ContentCopyOutlined';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import { Autocomplete, Breadcrumbs, Button, Container, Link, TextField } from '@mui/material';
import React, { FormEvent, ReactElement, useEffect, useState } from 'react';
import ReactMarkdown from "react-markdown";
import MDEditor from 'react-markdown-editor-lite';
import 'react-markdown-editor-lite/lib/index.css';
import { useLocation, useNavigate, useSearchParams } from 'react-router-dom';
import remarkGfm from 'remark-gfm';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { getNoteAsync, saveNoteAsync, selectNotesTree, TreeUtil } from '../reducer/noteSlice';
import { appendPath, getDecodedPath, getFilenameFromTitle, getTitleFromFilename, splitPath } from '../util/util';
import './Editor.scss';

const VALID_DIR_PATH_REGEX = /^[^/.]([/a-zA-Z0-9-]|[^\S\r\n])+([^/])$/gm;
const VALID_FILENAME_REGEX = /^([a-zA-Z0-9-]|[^\S\r\n])+(\.md)$/gm;

const Editor: React.FC = (): ReactElement => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const { pathname } = useLocation();
  const [searchParams] = useSearchParams();
  const editMode = pathname.startsWith('/edit');
  const path = getDecodedPath(searchParams.get('path'));
  const tree = useAppSelector(selectNotesTree);

  const [sha, setSHA] = useState('');
  const [title, setTitle] = useState('');
  const [titleError, setTitleError] = useState(false);
  const [content, setContent] = useState('');
  const [contentError, setContentError] = useState(false);
  const [endDir, setEndDir] = useState('');
  const [dirPathArray, setDirPathArray] = useState([] as string[]);
  const [dirPathError, setDirPathError] = useState(false);
  const [pathAutoCompleteOptions, setPathAutoCompleteOptions] = useState(TreeUtil.getChildDirs(tree, path));

  useEffect(() => {
    const treeNode = TreeUtil.searchNode(tree, path);
    const dirPathArray = splitPath(path);
    editMode && dirPathArray.pop(); // remove the filename from path
    setDirPathArray(dirPathArray);
    setPathAutoCompleteOptions(TreeUtil.getChildDirs(tree, path));

    if (treeNode == null || treeNode.is_dir) {
      return;
    }
    if (!treeNode.cached) {
      dispatch(getNoteAsync(treeNode.path));
      return;
    }

    setSHA(treeNode?.sha || '');
    setTitle(getTitleFromFilename(treeNode.name));
    setContent(treeNode?.content || '');
  }, [tree, path, editMode])

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setDirPathError(false);
    setTitleError(false);
    setContentError(false);

    const autoSelectedDirPath = dirPathArray.join('/');
    const dirPath = appendPath(autoSelectedDirPath, endDir);
    if (dirPath !== "" && !dirPath.match(VALID_DIR_PATH_REGEX)) {
      setDirPathError(true);
      return;
    }

    const filename = getFilenameFromTitle(title);
    if (!filename.match(VALID_FILENAME_REGEX)) {
      setTitleError(true);
      return;
    }

    if (content === "") {
      setContentError(true);
      return;
    }

    const fullPath = appendPath(dirPath, filename);
    await dispatch(saveNoteAsync({ path: fullPath, content: content, sha: sha }));
    navigate("/?path=" + encodeURIComponent(dirPath));
  }

  return (
    <Container maxWidth="md">
      <form noValidate autoComplete="off" onSubmit={handleSubmit}>
        <Autocomplete freeSolo fullWidth multiple openOnFocus value={dirPathArray} options={pathAutoCompleteOptions}
          disabled={editMode}
          onChange={(e, newPath) => {
            setDirPathArray([...newPath]);
            setPathAutoCompleteOptions(TreeUtil.getChildDirs(tree, newPath.join("/")));
          }}

          renderTags={(tagValue) => (
            <Breadcrumbs itemsAfterCollapse={2}>
              {tagValue.map((option) => (<Link key={option} underline="hover" color="inherit"> {option} </Link>))}
              <span>{/* just a placeholder to show a / at the end */}</span>
            </Breadcrumbs>
          )}

          inputValue={endDir}
          onInputChange={(e, newInputValue) => {
            setDirPathError(false);
            if (newInputValue.indexOf('/') > -1) {
              const trimmedPath = newInputValue.trim().replace(/^\/+|\/+$/g, '');
              const pathArray = [...dirPathArray, ...splitPath(trimmedPath)];
              if (trimmedPath) {
                setDirPathArray(pathArray);
                setPathAutoCompleteOptions(TreeUtil.getChildDirs(tree, pathArray.join("/")));
              }
              setEndDir('');
              return;
            }
            setEndDir(newInputValue);
          }}

          renderInput={(params) => (
            <TextField {...params}
              helperText="Only alphanumeric characters, space, hyphen (-) and forward slash (/) are allowed."
              label="Path" variant="outlined" fullWidth error={dirPathError} placeholder="Select Path..." sx={{ my: 2, display: "block" }} />
          )}
        />

        <TextField sx={{ my: 2, display: "block" }}
          helperText="Only alphanumeric characters, space and hyphen (-) are allowed."
          value={title} disabled={editMode}
          onChange={(e) => { setTitleError(false); setTitle(e.target.value) }} label="Note Title"
          variant="outlined" fullWidth required error={titleError}
        />

        <MDEditor view={{ menu: true, md: true, html: false }} canView={{ menu: true, md: true, html: true, fullScreen: false, hideMenu: false, both: true }}
          value={content}
          renderHTML={text => <ReactMarkdown components={{
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
            remarkPlugins={[remarkGfm]} >{text}</ReactMarkdown>}
          placeholder="Note Content*" className={"batnoter-md-editor " + (contentError ? "error" : "")}
          onChange={({ text }) => { setContentError(false); setContent(text) }} />

        <Button type="submit" variant="contained" endIcon={<KeyboardArrowRightIcon />} sx={{ float: 'right' }}> SAVE </Button>
      </form>
    </Container>
  )
}

export default Editor;
