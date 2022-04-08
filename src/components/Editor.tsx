
import SaveIcon from '@mui/icons-material/Save';
import { LoadingButton } from '@mui/lab';
import { Autocomplete, Breadcrumbs, Button, Container, Link, styled, TextField, Theme } from '@mui/material';
import React, { FormEvent, ReactElement, useEffect, useState } from 'react';
import MDEditor from 'react-markdown-editor-lite';
import 'react-markdown-editor-lite/lib/index.css';
import { useLocation, useNavigate, useSearchParams } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { getNoteAsync, NoteStatus, saveNoteAsync, selectNoteStatus, selectNotesTree, TreeUtil } from '../reducer/noteSlice';
import { appendPath, getDecodedPath, getFilenameFromTitle, getTitleFromFilename, splitPath } from '../util/util';
import CustomReactMarkdown from './lib/CustomReactMarkdown';

const VALID_DIR_PATH_REGEX = /^((?!\/)([a-zA-Z0-9-]([/]|[^\S\r\n])?)*)([a-zA-Z0-9-])$/gm;
const VALID_FILENAME_REGEX = /^([a-zA-Z0-9-]|[^\S\r\n])+(\.md)$/gm;

const StyledMDEditor = styled(MDEditor)(
  ({ theme }: { theme: Theme }) => `
  margin: 16px 0;
  height: 275px;
  border-color: rgba(0, 0, 0, 0.23);
  border-radius: ${theme.shape.borderRadius}px;
  & > .rc-md-navigation {
    min-height: 56px;
    border-radius: ${theme.shape.borderRadius}px ${theme.shape.borderRadius}px 0 0;

    .button-wrap {
      .button {
        margin: 0 5px;
      }
      .rmel-iconfont {
        font-size: ${theme.typography.fontSize + 8}px;
      }
    } 
  }

  &.error {
    border-color: ${theme.palette.error.main};
  }
  `,
);

const Editor: React.FC = (): ReactElement => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const { pathname } = useLocation();
  const [searchParams] = useSearchParams();
  const editMode = pathname.startsWith('/edit');
  const path = getDecodedPath(searchParams.get('path'));
  const tree = useAppSelector(selectNotesTree);
  const status = useAppSelector(selectNoteStatus);

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

        <StyledMDEditor view={{ menu: true, md: true, html: false }} canView={{ menu: true, md: true, html: true, fullScreen: false, hideMenu: false, both: true }}
          value={content}
          renderHTML={(text: string) => <CustomReactMarkdown>{text}</CustomReactMarkdown>}
          placeholder="Note Content*" className={"batnoter-md-editor " + (contentError ? "error" : "")}
          onChange={({ text }: { text: string }) => { setContentError(false); setContent(text) }} />

        <LoadingButton loading={status === NoteStatus.LOADING} type="submit" variant="contained" startIcon={<SaveIcon />} sx={{ float: 'right' }}>SAVE</LoadingButton>
        <Button onClick={() => navigate('/')} variant="outlined" sx={{ float: 'right', mx: 1 }} >CANCEL</Button>
      </form>
    </Container>
  )
}

export default Editor;
