
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight'
import { Button, Container, TextField, Typography } from '@mui/material'
import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAppDispatch } from '../app/hooks'
import { saveNoteAsync } from '../reducer/noteSlice'

const Editor: React.FC = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate()
  const [content, setContent] = useState('')
  const [path, setPath] = useState('')
  const [pathError, setPathError] = useState(false)

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setPathError(false)
    if (!path) {
      setPathError(true)
    }
    dispatch(saveNoteAsync({ path: path, content: content }))
    navigate('/')
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
          onChange={(e) => setPath(e.target.value)}
          label="Note Title"
          variant="outlined"
          fullWidth
          required
          error={pathError}
        />
        <TextField sx={{ my: 2, display: "block" }}
          onChange={(e) => setContent(e.target.value)}
          label="Content"
          variant="outlined"
          multiline
          rows={4}
          fullWidth
          required
        />

        <Button
          type="submit"
          variant="contained"
          endIcon={<KeyboardArrowRightIcon />}>
          Submit
        </Button>
      </form>
    </Container>
  )
}

export default Editor;