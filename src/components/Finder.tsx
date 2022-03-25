import { Masonry } from '@mui/lab';
import { Container } from '@mui/material';
import React, { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { deleteNoteAsync, Note, searchNotesAsync, selectNotesPage } from '../reducer/noteSlice';
import NoteCard from './NoteCard';


const Finder = () => {
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(searchNotesAsync())
  }, []);

  const page = useAppSelector(selectNotesPage)

  const handleDelete = (note: Note) => {
    dispatch(deleteNoteAsync(note))
  }

  return (
    <Container>
      <Masonry columns={{ xs: 1, md: 3, xl: 4 }} spacing={2}>
        {page.notes.map(note => (
          <div key={note.path}>
            <NoteCard note={note} handleDelete={handleDelete} />
          </div>
        ))}
      </Masonry>
    </Container>
  )
}

export default Finder