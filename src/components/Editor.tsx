
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight'
import { Autocomplete, Breadcrumbs, Button, Container, Link, TextField, Typography } from '@mui/material'
import React, { FormEvent, ReactElement, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAppDispatch, useAppSelector } from '../app/hooks'
import { saveNoteAsync, selectNotesTree, TreeUtil } from '../reducer/noteSlice'

const Editor: React.FC = (): ReactElement => {
  const VALID_DIR_PATH_REGEX = /^[^/.]([/a-zA-Z0-9-]|[^\S\r\n])+([^/])$/gm;
  const VALID_FILENAME_REGEX = /^([a-zA-Z0-9-]|[^\S\r\n])+(\.md)$/gm;
  const dispatch = useAppDispatch();
  const navigate = useNavigate()
  const tree = useAppSelector(selectNotesTree);
  const [dirPathArray, setDirPathArray] = useState([] as string[])
  const [endDir, setEndDir] = useState("")
  const [dirPathError, setDirPathError] = useState(false)

  const [content, setContent] = useState('')
  const [contentError, setContentError] = useState(false)
  const [title, setTitle] = useState('')
  const [titleError, setTitleError] = useState(false)
  const defaultPathOptions = TreeUtil.getChildDirs(tree, "")
  const [pathOptions, setPathOptions] = useState(defaultPathOptions)
  useEffect(() => {
    setPathOptions(defaultPathOptions);
  }, [tree])

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setDirPathError(false)
    setTitleError(false)
    setContentError(false)

    const autoSelectedDirPath = dirPathArray.join('/')
    const dirPath = autoSelectedDirPath + (endDir === "" ? "" : (autoSelectedDirPath === "" ? endDir : '/' + endDir))
    if (dirPath != "" && !VALID_DIR_PATH_REGEX.test(dirPath)) {
      setDirPathError(true)
      return
    }

    const filename = title + '.md';
    if (!VALID_FILENAME_REGEX.test(filename)) {
      setTitleError(true)
      return
    }

    if (content === "") {
      setContentError(true)
      return
    }

    const fullPath = dirPath !== "" ? (dirPath + '/' + filename) : filename;
    await dispatch(saveNoteAsync({ path: fullPath, content: content }))
    navigate("/?path=" + encodeURIComponent(dirPath))
  }

  return (
    <Container maxWidth="md">
      <Typography variant="h6" color="textSecondary" component="h2" gutterBottom >
        Create a New Note
      </Typography>

      <form noValidate autoComplete="off" onSubmit={handleSubmit}>
        <Autocomplete freeSolo fullWidth multiple openOnFocus value={dirPathArray} options={pathOptions}
          onChange={(e, newPath) => {
            setDirPathArray([...newPath]);
            setPathOptions(TreeUtil.getChildDirs(tree, newPath.join("/")))
          }}

          renderTags={(tagValue) => (
            <Breadcrumbs itemsAfterCollapse={2}>
              {tagValue.map((option) => (<Link key={option} underline="hover" color="inherit"> {option} </Link>))}
              <span>{/* just a placeholder to show a / at the end */}</span>
            </Breadcrumbs>
          )}

          inputValue={endDir}
          onInputChange={(e, newInputValue) => {
            setDirPathError(false)
            if (newInputValue.indexOf('/') > -1) {
              const trimmedPath = newInputValue.trim().replace(/^\/+|\/+$/g, '')
              const newPath = [...dirPathArray, ...trimmedPath.split('/')]
              if (trimmedPath) {
                setDirPathArray(newPath)
                setPathOptions(TreeUtil.getChildDirs(tree, newPath.join("/")))
              }
              setEndDir('')
              return
            }
            setEndDir(newInputValue)
          }}

          renderInput={(params) => (
            <TextField {...params}
              helperText="Only alphanumeric characters, space, hyphen (-) and forward slash (/) are allowed."
              label="Path" variant="outlined" fullWidth error={dirPathError} placeholder="Select Path..." sx={{ my: 2, display: "block" }} />
          )}
        />

        <TextField sx={{ my: 2, display: "block" }}
          helperText="Only alphanumeric characters, space and hyphen (-) are allowed."
          onChange={(e) => { setTitleError(false); setTitle(e.target.value) }} label="Note Title"
          variant="outlined" fullWidth required error={titleError}
        />
        <TextField sx={{ my: 2, display: "block" }}
          onChange={(e) => { setContentError(false); setContent(e.target.value) }} error={contentError} label="Content" variant="outlined" multiline rows={10} fullWidth required
        />
        <Button type="submit" variant="contained" endIcon={<KeyboardArrowRightIcon />} sx={{ float: 'right' }}> SAVE </Button>
      </form>

    </Container>
  )
}

export default Editor;