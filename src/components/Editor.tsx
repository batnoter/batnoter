
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import { Autocomplete, Breadcrumbs, Button, Container, Link, TextField } from '@mui/material';
import React, { FormEvent, ReactElement, useEffect, useState } from 'react';
import ReactMarkdown from "react-markdown";
import MDEditor from 'react-markdown-editor-lite';
import 'react-markdown-editor-lite/lib/index.css';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import remarkGfm from 'remark-gfm';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { getNoteAsync, saveNoteAsync, selectNotesTree, TreeUtil } from '../reducer/noteSlice';
import { getTitleFromFilename } from '../util/util';
import './Editor.scss';

const Editor: React.FC = (): ReactElement => {
  const VALID_DIR_PATH_REGEX = /^[^/.]([/a-zA-Z0-9-]|[^\S\r\n])+([^/])$/gm;
  const VALID_FILENAME_REGEX = /^([a-zA-Z0-9-]|[^\S\r\n])+(\.md)$/gm;
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const tree = useAppSelector(selectNotesTree);

  const [editMode, setEditMode] = useState(false);
  const [sha, setSHA] = useState('');
  const [title, setTitle] = useState('');
  const [titleError, setTitleError] = useState(false);
  const [content, setContent] = useState('');
  const [contentError, setContentError] = useState(false);
  const [endDir, setEndDir] = useState('');
  const [dirPathArray, setDirPathArray] = useState([] as string[]);
  const [dirPathError, setDirPathError] = useState(false);

  const defaultPathOptions = TreeUtil.getChildDirs(tree, '');
  const [pathAutoCompleteOptions, setPathAutoCompleteOptions] = useState(defaultPathOptions);

  const { pathname } = useLocation();
  const { path } = useParams();

  useEffect(() => {
    if (path != null && pathname.startsWith('/edit/')) {
      setEditMode(true);
      const decodedPath = decodeURIComponent(path);
      const note = TreeUtil.searchNode(tree, decodedPath);
      if (note == null) {
        return;
      }
      if (!note.cached) {
        dispatch(getNoteAsync(note.path));
        return;
      }
      const pathArray = path.split('/');
      pathArray.pop(); // remove the filename from path
      setDirPathArray(pathArray)
      setSHA(note?.sha || '');
      setTitle(getTitleFromFilename(note.name));
      setContent(note?.content || '');
    }
    setPathAutoCompleteOptions(defaultPathOptions);
  }, [tree, path, pathname])

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setDirPathError(false);
    setTitleError(false);
    setContentError(false);

    const autoSelectedDirPath = dirPathArray.join('/');
    const dirPath = autoSelectedDirPath + (endDir === "" ? "" : (autoSelectedDirPath === "" ? endDir : '/' + endDir));
    if (dirPath != "" && !dirPath.match(VALID_DIR_PATH_REGEX)) {
      setDirPathError(true);
      return;
    }

    const filename = title + '.md';
    if (!filename.match(VALID_FILENAME_REGEX)) {
      setTitleError(true);
      return;
    }

    if (content === "") {
      setContentError(true);
      return;
    }

    const fullPath = dirPath !== "" ? (dirPath + '/' + filename) : filename;
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
              const newPath = [...dirPathArray, ...trimmedPath.split('/')];
              if (trimmedPath) {
                setDirPathArray(newPath);
                setPathAutoCompleteOptions(TreeUtil.getChildDirs(tree, newPath.join("/")));
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
          renderHTML={text => <ReactMarkdown remarkPlugins={[remarkGfm]} >{text}</ReactMarkdown>}
          placeholder="Note Content*" className={"batnoter-md-editor " + (contentError ? "error" : "")}
          onChange={({ text }) => { setContentError(false); setContent(text) }} />

        <Button type="submit" variant="contained" endIcon={<KeyboardArrowRightIcon />} sx={{ float: 'right' }}> SAVE </Button>
      </form>

    </Container>
  )
}

export default Editor;