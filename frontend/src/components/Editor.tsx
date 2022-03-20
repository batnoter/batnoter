
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight'
import { Button, Container, TextField, Typography } from '@mui/material'
import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAppDispatch } from '../app/hooks'
import { saveNoteAsync } from '../reducer/note/noteSlice'


const Editor = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate()
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [titleError, setTitleError] = useState(false)
  const [contentError, setContentError] = useState(false)

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setTitleError(false)
    setContentError(false)

    if (title == '') {
      setTitleError(true)
    }
    if (content == '') {
      setContentError(true)
    }
    if (title && content) {
      dispatch(saveNoteAsync({ title, content }))
      navigate('/')
    }
  }

  return (
    <Container maxWidth="sm">
      <Typography
        variant="h6"
        color="textSecondary"
        component="h2"
        gutterBottom
      >
        Create a New Note
      </Typography>

      <form noValidate autoComplete="off" onSubmit={handleSubmit}>
        <TextField sx={{ my: 2, display: "block" }}
          onChange={(e) => setTitle(e.target.value)}
          label="Note Title"
          variant="outlined"
          color="secondary"
          fullWidth
          required
          error={titleError}
        />
        <TextField sx={{ my: 2, display: "block" }}
          onChange={(e) => setContent(e.target.value)}
          label="Content"
          variant="outlined"
          color="secondary"
          multiline
          rows={4}
          fullWidth
          required
          error={contentError}
        />

        <Button
          type="submit"
          color="secondary"
          variant="contained"
          endIcon={<KeyboardArrowRightIcon />}>
          Submit
        </Button>
      </form>
    </Container>
  )
}

export default Editor;